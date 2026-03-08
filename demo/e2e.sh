#!/usr/bin/env bash
set -euo pipefail

# CRE Risk Router - End-to-End Integration Demo
# Demonstrates the full risk evaluation pipeline:
#   Cron trigger -> 8 risk gates -> ABI encode -> DON consensus -> on-chain receipt
#
# Usage:
#   ./demo/e2e.sh              # dry-run simulation
#   ./demo/e2e.sh --broadcast  # write receipt on-chain (requires funded wallet)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CRE_CLI="${CRE_CLI:-cre}"
BROADCAST_FLAG=""

if [[ "${1:-}" == "--broadcast" ]]; then
    BROADCAST_FLAG="--broadcast"
    echo "=== CRE Risk Router E2E Demo (BROADCAST MODE) ==="
    echo "Transactions will be written on-chain to Sepolia testnet."
else
    echo "=== CRE Risk Router E2E Demo (DRY-RUN MODE) ==="
    echo "No on-chain transactions. Use --broadcast to write receipts."
fi

echo ""
echo "Contract: 0x9C7Aa5502ad229c80894E272Be6d697Fd02001d7 (Sepolia)"
echo "Risk Gates: 8 sequential gates (hold, confidence, risk score, staleness,"
echo "            oracle health, price deviation, position sizing, heartbeat)"
echo ""

echo "--- Synthetic RiskRequest ---"
echo "  AgentID:          agent-inference-001"
echo "  Signal:           buy"
echo "  Confidence:       0.85"
echo "  RiskScore:        10"
echo "  RequestedPosition: \$1000 (6-decimal)"
echo "  MarketPair:       ETH/USD"
echo ""

echo "--- Running CRE Simulation ---"
cd "$PROJECT_DIR"
$CRE_CLI workflow simulate . \
    --non-interactive \
    --trigger-index 0 \
    -T staging-settings \
    $BROADCAST_FLAG

echo ""
echo "--- Expected Outcome ---"
echo "  Approved:         true"
echo "  MaxPositionUSD:   810000000 (~\$810, constrained by risk/volatility)"
echo "  MaxSlippageBps:   500 (5%)"
echo "  ChainlinkPrice:   200000000000 (\$2000 at 8 decimals)"
echo "  TTLSeconds:       300 (5 minutes)"
echo ""

if [[ -n "$BROADCAST_FLAG" ]]; then
    echo "--- Verify On-Chain ---"
    echo "  Etherscan: https://sepolia.etherscan.io/address/0x9C7Aa5502ad229c80894E272Be6d697Fd02001d7"
    echo "  Look for DecisionRecorded event in the latest transaction."
fi

echo ""
echo "=== Demo Complete ==="
