package main

import (
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	cronTrigger := cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"})

	return cre.Workflow[*Config]{
		cre.Handler(cronTrigger, onCronTrigger),
	}, nil
}

func onCronTrigger(config *Config, runtime cre.Runtime, trigger *cron.Payload) (*RiskDecision, error) {
	// TODO: implement risk evaluation pipeline (Phase 005)
	return nil, nil
}
