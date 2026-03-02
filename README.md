# CRE Risk Router

An on-chain risk decision layer for autonomous DeFi agents, built on the Chainlink Runtime Environment (CRE). The Risk Router evaluates trade signals through 8 sequential risk gates and writes immutable decision receipts to Ethereum via Chainlink DON consensus.

## Architecture

```
Cron Trigger
    |
    v
RiskRequest (synthetic or from agent-coordinator)
    |
    v
+------------------+
| 8 Risk Gates     |
| (sequential)     |
|                  |
| Gate 7: Hold     |  <-- fast-path filter
| Gate 1: Confidence|
| Gate 2: Risk Score|
| Gate 3: Staleness |
| Gate 4: Oracle    |  <-- Chainlink feed health
| Gate 5: Price Dev |  <-- oracle vs market
| Gate 6: Position  |  <-- volatility-adjusted sizing
| Gate 8: Heartbeat |  <-- agent liveness
+------------------+
    |
    v
RiskDecision (approved/denied + constraints)
    |
    v
ABI Encode -> GenerateReport -> WriteReport
    |
    v
RiskDecisionReceipt.sol (on-chain, Sepolia)
```

## Risk Gates

| Gate | Name | Type | Description |
|------|------|------|-------------|
| 7 | Hold Signal Filter | Deny | Rejects `hold` signals immediately |
| 1 | Signal Confidence | Deny | Requires confidence >= threshold (default 0.6) |
| 2 | Risk Score Ceiling | Deny | Rejects scores above maximum (default 75) |
| 3 | Signal Staleness | Deny | Rejects signals older than TTL (default 300s) |
| 4 | Oracle Health | Deny | Validates Chainlink latestRoundData() 5-tuple |
| 5 | Price Deviation | Deny | Rejects if market vs oracle diverge >500 BPS |
| 6 | Position Sizing | Constrain | Adjusts position based on volatility and risk |
| 8 | Heartbeat | Deny | Circuit breaker for agent liveness (optional) |

## Quick Start

```bash
# 1. Clone and enter project
git clone https://github.com/lancekrogers/cre-risk-router.git
cd cre-risk-router

# 2. Install CRE CLI (requires Go 1.25+)
go install github.com/smartcontractkit/cre-cli@latest

# 3. Login to CRE
cre auth login

# 4. Run simulation (dry-run, no on-chain write)
cre workflow simulate . --non-interactive --trigger-index 0 -T staging-settings
```

## Simulation Commands

**Dry-run** (no on-chain transaction):
```bash
cre workflow simulate . --non-interactive --trigger-index 0 -T staging-settings
```

**Broadcast** (writes receipt on-chain, requires `CRE_ETH_PRIVATE_KEY` in `.env`):
```bash
cre workflow simulate . --broadcast --non-interactive --trigger-index 0 -T staging-settings
```

**E2E demo script**:
```bash
./demo/e2e.sh              # dry-run
./demo/e2e.sh --broadcast  # on-chain write
```

### Expected Output

```json
{
  "Approved": true,
  "ChainlinkPrice": 200000000000,
  "MaxPositionUSD": 810000000,
  "MaxSlippageBps": 500,
  "Reason": "approved",
  "TTLSeconds": 300
}
```

## Scenarios

Pre-built simulation scenarios in `scenarios/`:

| Scenario | Signal | Confidence | Risk | Expected |
|----------|--------|------------|------|----------|
| `approved_trade.json` | buy | 0.85 | 10 | Approved, $810 constrained |
| `denied_low_confidence.json` | buy | 0.45 | 35 | Denied: confidence below threshold |
| `denied_high_risk.json` | sell | 0.90 | 82 | Denied: risk score exceeds max |
| `denied_stale_signal.json` | buy | 0.80 | 40 | Denied: signal expired |
| `denied_price_deviation.json` | buy | 0.85 | 20 | Denied: price deviation >5% |

## Configuration

Key fields in `config.staging.json`:

| Field | Default | Description |
|-------|---------|-------------|
| `signal_confidence_threshold` | 0.6 | Minimum signal confidence (0.0-1.0) |
| `max_risk_score` | 75 | Maximum allowed risk score (0-100) |
| `decision_ttl_seconds` | 300 | Signal freshness window |
| `price_deviation_max_bps` | 500 | Max oracle-market divergence (5%) |
| `volatility_scale_factor` | 1.0 | Position scaling sensitivity |
| `oracle_staleness_seconds` | 3600 | Max Chainlink feed age |
| `enable_heartbeat_gate` | false | Agent liveness check |

## Contract

**RiskDecisionReceipt.sol** deployed on Sepolia:
- Address: [`0xfcA344515D72a05232DF168C1eA13Be22383cCB6`](https://sepolia.etherscan.io/address/0xfcA344515D72a05232DF168C1eA13Be22383cCB6)
- Broadcast tx: [`0xd8505ff76caa1e2d17b2ee49b625048f353359fabf68f02abedc9fda87360458`](https://sepolia.etherscan.io/tx/0xd8505ff76caa1e2d17b2ee49b625048f353359fabf68f02abedc9fda87360458)

Features:
- `recordDecision()` with duplicate prevention per runId
- `isDecisionValid()` with TTL-based expiry
- `DecisionRecorded` event for off-chain indexing
- Approval/denial counters

## Project Structure

```
cre-risk-router/
  main.go               # WASM entrypoint
  workflow.go            # CRE handlers (onScheduledSweep, executeRiskPipeline)
  risk.go                # 8 risk gates + evaluateRisk pipeline
  types.go               # Config, RiskRequest, RiskDecision, MarketData
  helpers.go             # keccak256 hashing, slippage, position math
  workflow.yaml           # CRE workflow settings
  project.yaml            # RPC endpoints
  config.staging.json     # Gate thresholds (staging)
  config.production.json  # Gate thresholds (production)
  contracts/evm/src/
    RiskDecisionReceipt.sol
    generated/             # CRE-generated Go bindings
  scenarios/              # Pre-built simulation inputs
  demo/e2e.sh             # End-to-end demo script
  test/                   # Foundry tests
```

## Tech Stack

- **Runtime**: Chainlink Runtime Environment (CRE) v1.2.0
- **Language**: Go (compiled to WASM)
- **Chain**: Ethereum Sepolia testnet
- **Oracle**: Chainlink price feeds (latestRoundData)
- **Contracts**: Solidity 0.8.x, Foundry
- **Consensus**: CRE report-based DON consensus
