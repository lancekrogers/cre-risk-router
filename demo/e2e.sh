#!/usr/bin/env bash
set -euo pipefail

# CRE Risk Router - End-to-End Demo
# Runs simulation with each scenario file and prints results.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "=== CRE Risk Router E2E Demo ==="
echo ""

# Simulate
echo "--- Running CRE Simulation ---"
cd "$PROJECT_DIR"
cre workflow simulate . --non-interactive --trigger-index=0 --target=staging-settings

echo ""
echo "=== Demo Complete ==="
