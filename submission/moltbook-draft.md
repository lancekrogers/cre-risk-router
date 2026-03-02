# CRE Risk Router -- On-Chain Risk Decisions for Autonomous DeFi Agents

## Summary

CRE Risk Router is a Chainlink Runtime Environment (CRE) workflow that evaluates trade signals from autonomous DeFi agents through 8 sequential risk gates. Every decision -- approved or denied -- is recorded on-chain via a RiskDecisionReceipt smart contract using Chainlink DON consensus. The result is a transparent, auditable, and trustless risk decision layer where no trade executes without on-chain proof of evaluation.

## Problem Statement

Autonomous AI agents are increasingly generating trade signals in DeFi, but they lack a systematic risk evaluation layer. Without guardrails, an agent can execute trades based on stale market data, take positions far beyond safe exposure, or continue trading during oracle outages. Traditional risk management systems require centralized trust and offer no on-chain auditability.

CRE Risk Router solves this by providing a deterministic, configurable risk pipeline that runs on Chainlink's decentralized oracle network. Every decision is verifiable on-chain, creating an immutable audit trail that agents, protocols, and regulators can inspect.

## Architecture

```
Cron Trigger (every 5 min)
    |
    v
Synthetic RiskRequest
  {agent_id, signal, confidence, risk_score, market_pair, position, timestamp}
    |
    v
+-------------------------------------------+
|          8 Sequential Risk Gates           |
|                                            |
|  Gate 7: Hold Signal Filter (fast-path)    |
|  Gate 1: Signal Confidence >= 0.6          |
|  Gate 2: Risk Score <= 75                  |
|  Gate 3: Signal Age <= 300s                |
|  Gate 4: Chainlink Oracle Health           |
|  Gate 5: Price Deviation <= 500 BPS        |
|  Gate 6: Volatility-Adjusted Position      |
|  Gate 8: Agent Heartbeat Circuit Breaker   |
+-------------------------------------------+
    |
    v
RiskDecision {approved, maxPositionUSD, maxSlippageBps, reason}
    |
    v
ABI Encode -> CRE GenerateReport -> DON Consensus
    |
    v
RiskDecisionReceipt.sol::recordDecision()
    |
    v
DecisionRecorded event (on-chain, Sepolia)
```

**CRE Capabilities Used:**
- `cron-trigger@1.0.0` -- Periodic risk sweep every 5 minutes
- `evm-write` -- On-chain receipt via report-based DON consensus
- Go WASM compilation -- Deterministic execution in CRE runtime

## Risk Gates

| # | Gate | Type | Description | Default |
|---|------|------|-------------|---------|
| 7 | Hold Signal Filter | Hard Deny | Rejects `hold` signals immediately (fast-path) | N/A |
| 1 | Signal Confidence | Hard Deny | Requires `signal_confidence >= threshold` | 0.6 |
| 2 | Risk Score Ceiling | Hard Deny | Rejects if `risk_score > max_risk_score` | 75 |
| 3 | Signal Staleness | Hard Deny | Rejects signals older than TTL | 300s |
| 4 | Oracle Health | Hard Deny | Validates Chainlink `latestRoundData()` 5-tuple: positive answer, non-zero updatedAt, answeredInRound >= roundID, freshness | 3600s staleness |
| 5 | Price Deviation | Soft Deny | Rejects if market price diverges from Chainlink oracle by more than threshold | 500 BPS (5%) |
| 6 | Position Sizing | Constrain | Adjusts position down based on volatility and risk score; never denies | 10% fallback vol |
| 8 | Agent Heartbeat | Hard Deny | Circuit breaker if agent has not sent a heartbeat within TTL (optional, disabled by default) | 600s |

## Smart Contract

**RiskDecisionReceipt.sol** on Ethereum Sepolia:

- **Address:** `0xfcA344515D72a05232DF168C1eA13Be22383cCB6`
- **Chain:** Sepolia (Chain ID 11155111)
- **Deploy tx:** `0x36c066ba6a3d29abf6888382d5c44c014c7bff4443895cf6a7c84092c4314b46`

**Interface:**

```solidity
function recordDecision(
    bytes32 runId,
    bytes32 decisionHash,
    bool approved,
    uint256 maxPositionUsd,
    uint256 maxSlippageBps,
    uint256 ttlSeconds,
    uint256 chainlinkPrice
) external;

function isDecisionValid(bytes32 runId) external view returns (bool);
function getRunCount() external view returns (uint256);

event DecisionRecorded(
    bytes32 indexed runId,
    bytes32 indexed decisionHash,
    bool approved,
    uint256 maxPositionUsd,
    uint256 maxSlippageBps,
    uint256 ttlSeconds,
    uint256 chainlinkPrice,
    uint256 timestamp
);
```

**Features:** Duplicate prevention per runId, TTL-based expiry via `isDecisionValid()`, approval/denial counters.

## Simulation

### Dry-Run Output

```
cre workflow simulate . --non-interactive --trigger-index 0 -T staging-settings
```

```
[USER LOG] msg="Scheduled risk sweep triggered"
[USER LOG] msg="Risk decision" approved=true reason=approved maxPositionUSD=810000000 maxSlippageBps=500 chainlinkPrice=200000000000
[USER LOG] msg="Receipt written on-chain" txHash=0x0000000000000000000000000000000000000000000000000000000000000000
```

Result: Approved with $810 constrained position (from $1000 requested), 500 BPS slippage, $2000 Chainlink price.

### Broadcast Output

```
cre workflow simulate . --broadcast --non-interactive --trigger-index 0 -T staging-settings
```

```
[USER LOG] msg="Risk decision" approved=true reason=approved maxPositionUSD=810000000 maxSlippageBps=500 chainlinkPrice=200000000000
[USER LOG] msg="Receipt written on-chain" txHash=0x4cd1d6664747b5e2c53f1e10b819b50d437827d632212d204d941b1130c068f2
```

Real transaction submitted and confirmed on Sepolia.

## Evidence

### On-Chain Transactions

| # | Tx Hash | Block | Status |
|---|---------|-------|--------|
| 1 | [`0xd8505f...0458`](https://sepolia.etherscan.io/tx/0xd8505ff76caa1e2d17b2ee49b625048f353359fabf68f02abedc9fda87360458) | 10367283 | Success |
| 2 | [`0x4cd1d6...68f2`](https://sepolia.etherscan.io/tx/0x4cd1d6664747b5e2c53f1e10b819b50d437827d632212d204d941b1130c068f2) | 10367301 | Success |

Both transactions:
- Called via CRE report-based DON consensus forwarder (`0x15fC6ae953E024d975e77382eEeC56A9101f9F88`)
- Targeted RiskDecisionReceipt at `0xfcA344515D72a05232DF168C1eA13Be22383cCB6`
- Emitted event with approval status and decision parameters

### Contract Deployment

- Deploy tx: [`0x36c066...4b46`](https://sepolia.etherscan.io/tx/0x36c066ba6a3d29abf6888382d5c44c014c7bff4443895cf6a7c84092c4314b46)
- Deployer: `0xC71d8A19422C649fe9bdCbF3ffA536326c82b58b`

## Scenarios

| Scenario | Signal | Confidence | Risk | Expected | Gate |
|----------|--------|------------|------|----------|------|
| `approved_trade` | buy | 0.85 | 10 | Approved ($810 constrained) | All passed |
| `denied_low_confidence` | buy | 0.45 | 35 | Denied | Gate 1: confidence below 0.6 |
| `denied_high_risk` | sell | 0.90 | 82 | Denied | Gate 2: risk score exceeds 75 |
| `denied_stale_signal` | buy | 0.80 | 40 | Denied | Gate 3: signal expired (600s > 300s TTL) |
| `denied_price_deviation` | buy | 0.85 | 20 | Denied | Gate 5: market vs oracle >5% |

Full scenario JSON files available in `scenarios/` directory.

## Integration Path

### Current State (P0 -- Hackathon)

CRE Risk Router runs as a standalone CRE workflow with cron-triggered synthetic requests. The workflow demonstrates the full pipeline: risk evaluation through 8 gates, DON consensus, and on-chain receipt writing.

### Next Phase (P1 -- Agent Integration)

```
agent-inference-hedera
    |
    | produces trade signal
    v
agent-coordinator
    |
    | HTTP POST /evaluate-risk
    v
CRE Risk Router
    |
    | RiskDecision {approved, maxPosition, slippage}
    v
agent-coordinator
    |
    | if approved: assign task with constraints
    | if denied: log reason, skip execution
    v
agent-execution (DeFi trade)
```

**P1 Requirements:**
1. CRE HTTP trigger for real-time evaluation (pending SDK support)
2. Agent-coordinator CRE client package for calling the workflow
3. Signal field mapping from agent-inference output format
4. DeFi guard integration for position limits enforcement
5. HCS (Hedera Consensus Service) message logging for cross-chain audit trail
