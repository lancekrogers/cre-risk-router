package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	receipt "github.com/lancekrogers/cre-risk-router/contracts/evm/src/generated/risk_decision_receipt"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	sweepTrigger := cron.Trigger(&cron.Config{Schedule: "0 */5 * * * *"})

	return cre.Workflow[*Config]{
		cre.Handler(sweepTrigger, onScheduledSweep),
	}, nil
}

// onScheduledSweep is the cron-triggered handler that generates a synthetic
// RiskRequest with realistic parameters and runs the full risk evaluation
// pipeline. This is the primary entry point for CRE simulation.
func onScheduledSweep(config *Config, runtime cre.Runtime, trigger *cron.Payload) (*RiskDecision, error) {
	logger := runtime.Logger()
	now := trigger.ScheduledExecutionTime.AsTime().Unix()
	logger.Info("Scheduled risk sweep triggered", "time", time.Unix(now, 0))

	// Synthetic RiskRequest with realistic parameters
	req := RiskRequest{
		AgentID:           "agent-inference-001",
		TaskID:            fmt.Sprintf("task-sweep-%d", now),
		Signal:            "buy",
		SignalConfidence:  0.85,
		RiskScore:         10,
		MarketPair:        "ETH/USD",
		RequestedPosition: 1000_000000, // $1000 in 6-decimal
		Timestamp:         now,         // uses trigger time so Gate 3 always passes
	}

	decision, err := executeRiskPipeline(config, runtime, req, now)
	if err != nil {
		return nil, fmt.Errorf("scheduled sweep failed: %w", err)
	}
	return decision, nil
}

// executeRiskPipeline runs the full risk evaluation and on-chain receipt
// write pipeline. Shared by all trigger handlers.
func executeRiskPipeline(config *Config, runtime cre.Runtime, req RiskRequest, now int64) (*RiskDecision, error) {
	logger := runtime.Logger()

	chainSelector, err := evm.ChainSelectorFromName(config.TargetNetwork)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain selector: %w", err)
	}
	evmClient := &evm.Client{ChainSelector: chainSelector}

	// Market data placeholder (CRE HTTP fetch not available in SDK v1.2.0)
	var market *MarketData
	// When HTTP capability is available:
	// resp, err := runtime.HTTP().Get(config.MarketDataURL).Await()
	// For now, use nil (triggers fallback: skip Gate 5, 10% fallback volatility)

	// Oracle data: In simulation, Chainlink read returns mock data.
	// In production, this would read latestRoundData() from the price feed.
	oracle := OracleData{
		RoundID:         1,
		Answer:          200000000000, // $2000 at 8 decimals
		StartedAt:       now - 60,
		UpdatedAt:       now - 60,
		AnsweredInRound: 1,
	}

	decision := evaluateRisk(req, market, oracle, *config, now, 0)

	logger.Info("Risk decision",
		"approved", decision.Approved,
		"reason", decision.Reason,
		"maxPositionUSD", decision.MaxPositionUSD,
		"maxSlippageBps", decision.MaxSlippageBps,
		"chainlinkPrice", decision.ChainlinkPrice,
	)

	receiptContract, err := receipt.NewRiskDecisionReceipt(evmClient, common.HexToAddress(config.ReceiptContractAddress), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create receipt binding: %w", err)
	}

	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	boolType, _ := abi.NewType("bool", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)

	args := abi.Arguments{
		{Type: bytes32Type}, // runId
		{Type: bytes32Type}, // decisionHash
		{Type: boolType},    // approved
		{Type: uint256Type}, // maxPositionUsd
		{Type: uint256Type}, // maxSlippageBps
		{Type: uint256Type}, // ttlSeconds
		{Type: uint256Type}, // chainlinkPrice
	}

	encoded, err := args.Pack(
		decision.RunID,
		decision.DecisionHash,
		decision.Approved,
		new(big.Int).SetUint64(decision.MaxPositionUSD),
		new(big.Int).SetUint64(decision.MaxSlippageBps),
		new(big.Int).SetUint64(decision.TTLSeconds),
		new(big.Int).SetUint64(decision.ChainlinkPrice),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode decision: %w", err)
	}

	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}

	resp, err := receiptContract.WriteReport(runtime, report, nil).Await()
	if err != nil {
		return nil, fmt.Errorf("failed to write receipt: %w", err)
	}
	logger.Info("Receipt written on-chain", "txHash", common.BytesToHash(resp.TxHash).Hex())

	decisionJSON, _ := json.Marshal(decision)
	logger.Info("Decision output", "decision", string(decisionJSON))

	return &decision, nil
}
