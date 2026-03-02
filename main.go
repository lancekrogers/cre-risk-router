//go:build wasip1

package main

import (
	"github.com/lancekrogers/cre-risk-router/pkg/riskeval"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

func main() {
	wasm.NewRunner(cre.ParseJSON[riskeval.Config]).Run(InitWorkflow)
}
