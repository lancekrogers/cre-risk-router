# CRE Risk Router

# Show available recipes
@_default:
    just --list --justfile {{source_file()}}

# Install dependencies
install:
    go mod tidy

# Build WASM binary
build:
    GOOS=wasip1 GOARCH=wasm go build -o risk-router.wasm .

# Run local HTTP bridge for coordinator integration
bridge:
    go run ./cmd/bridge

# Run CRE simulation (dry run)
simulate:
    cre workflow simulate . --non-interactive --trigger-index=0 --target=staging-settings

# Run CRE broadcast (on-chain)
broadcast:
    cre workflow simulate . --non-interactive --trigger-index=0 --target=staging-settings --broadcast

# Run all tests (Go + Solidity)
test: test-go test-sol

# Run Go tests
test-go:
    go test ./... -v

# Run Solidity tests
test-sol:
    forge test -vvv

# Run linter
lint:
    go vet ./...

# Format Go code
fmt:
    gofmt -w .

# Generate EVM bindings from ABI files
bindings:
    cre generate-bindings evm

# Build Solidity contracts
forge-build:
    forge build

# Deploy contract to Sepolia
deploy:
    forge create contracts/evm/src/RiskDecisionReceipt.sol:RiskDecisionReceipt \
        --rpc-url $TESTNET_RPC \
        --private-key $DEPLOYER_PRIVATE_KEY

# Run end-to-end demo
demo:
    bash demo/e2e.sh
