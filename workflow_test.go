package main

import "testing"

func TestConfigFields(t *testing.T) {
	config := Config{
		MarketDataURL:             "https://api.coingecko.com/api/v3/coins/ethereum",
		SignalConfidenceThreshold: 0.6,
		MaxRiskScore:              75,
		DefaultMaxPositionBPS:     10000,
		DecisionTTLSeconds:        300,
		PriceDeviationMaxBPS:      500,
		VolatilityScaleFactor:     1.0,
		OracleStalenessSeconds:    3600,
		FeedDecimals:              8,
		HeartbeatTTLSeconds:       600,
	}

	if config.SignalConfidenceThreshold != 0.6 {
		t.Errorf("expected 0.6, got %f", config.SignalConfidenceThreshold)
	}
	if config.MaxRiskScore != 75 {
		t.Errorf("expected 75, got %d", config.MaxRiskScore)
	}
}
