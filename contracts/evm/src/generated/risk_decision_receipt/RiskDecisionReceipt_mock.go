// Code generated — DO NOT EDIT.

//go:build !wasip1

package risk_decision_receipt

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
)

var (
	_ = errors.New
	_ = fmt.Errorf
	_ = big.NewInt
	_ = common.Big1
)

// RiskDecisionReceiptMock is a mock implementation of RiskDecisionReceipt for testing.
type RiskDecisionReceiptMock struct {
	Decisions       func(DecisionsInput) (DecisionsOutput, error)
	GetRunCount     func() (*big.Int, error)
	IsDecisionValid func(IsDecisionValidInput) (bool, error)
	Recorded        func(RecordedInput) (bool, error)
	RunIds          func(RunIdsInput) ([32]byte, error)
	TotalApproved   func() (*big.Int, error)
	TotalDecisions  func() (*big.Int, error)
	TotalDenied     func() (*big.Int, error)
}

// NewRiskDecisionReceiptMock creates a new RiskDecisionReceiptMock for testing.
func NewRiskDecisionReceiptMock(address common.Address, clientMock *evmmock.ClientCapability) *RiskDecisionReceiptMock {
	mock := &RiskDecisionReceiptMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["decisions"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Decisions == nil {
				return nil, errors.New("decisions method not mocked")
			}
			inputs := abi.Methods["decisions"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := DecisionsInput{
				Arg0: values[0].([32]byte),
			}

			result, err := mock.Decisions(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["decisions"].Outputs.Pack(
				result.RunId,
				result.DecisionHash,
				result.Approved,
				result.MaxPositionUsd,
				result.MaxSlippageBps,
				result.TtlSeconds,
				result.ChainlinkPrice,
				result.Timestamp,
				result.Recorder,
			)
		},
		string(abi.Methods["getRunCount"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetRunCount == nil {
				return nil, errors.New("getRunCount method not mocked")
			}
			result, err := mock.GetRunCount()
			if err != nil {
				return nil, err
			}
			return abi.Methods["getRunCount"].Outputs.Pack(result)
		},
		string(abi.Methods["isDecisionValid"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.IsDecisionValid == nil {
				return nil, errors.New("isDecisionValid method not mocked")
			}
			inputs := abi.Methods["isDecisionValid"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := IsDecisionValidInput{
				RunId: values[0].([32]byte),
			}

			result, err := mock.IsDecisionValid(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["isDecisionValid"].Outputs.Pack(result)
		},
		string(abi.Methods["recorded"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Recorded == nil {
				return nil, errors.New("recorded method not mocked")
			}
			inputs := abi.Methods["recorded"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := RecordedInput{
				Arg0: values[0].([32]byte),
			}

			result, err := mock.Recorded(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["recorded"].Outputs.Pack(result)
		},
		string(abi.Methods["runIds"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.RunIds == nil {
				return nil, errors.New("runIds method not mocked")
			}
			inputs := abi.Methods["runIds"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := RunIdsInput{
				Arg0: values[0].(*big.Int),
			}

			result, err := mock.RunIds(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["runIds"].Outputs.Pack(result)
		},
		string(abi.Methods["totalApproved"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TotalApproved == nil {
				return nil, errors.New("totalApproved method not mocked")
			}
			result, err := mock.TotalApproved()
			if err != nil {
				return nil, err
			}
			return abi.Methods["totalApproved"].Outputs.Pack(result)
		},
		string(abi.Methods["totalDecisions"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TotalDecisions == nil {
				return nil, errors.New("totalDecisions method not mocked")
			}
			result, err := mock.TotalDecisions()
			if err != nil {
				return nil, err
			}
			return abi.Methods["totalDecisions"].Outputs.Pack(result)
		},
		string(abi.Methods["totalDenied"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TotalDenied == nil {
				return nil, errors.New("totalDenied method not mocked")
			}
			result, err := mock.TotalDenied()
			if err != nil {
				return nil, err
			}
			return abi.Methods["totalDenied"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
