package main

import "math"

// Gate 1: Signal Confidence Threshold
func checkSignalConfidence(req RiskRequest, cfg Config) (bool, string) {
	if req.SignalConfidence < cfg.SignalConfidenceThreshold {
		return false, "signal_confidence_below_threshold"
	}
	return true, ""
}

// Gate 2: Risk Score Ceiling
func checkRiskScore(req RiskRequest, cfg Config) (bool, string) {
	if req.RiskScore > cfg.MaxRiskScore {
		return false, "risk_score_exceeds_maximum"
	}
	return true, ""
}

// Gate 3: Signal Staleness
func checkSignalStaleness(req RiskRequest, cfg Config, now int64) (bool, string) {
	if (now - req.Timestamp) > int64(cfg.DecisionTTLSeconds) {
		return false, "signal_expired"
	}
	return true, ""
}

// Gate 4: Chainlink Oracle Health
// Validates the 5-tuple from latestRoundData().
// Returns (true, chainlinkPrice, "") on success, or (false, 0, reason) on failure.
func checkOracleHealth(
	roundID, answer, startedAt, updatedAt, answeredInRound int64,
	cfg Config,
	now int64,
) (bool, uint64, string) {
	if answer <= 0 {
		return false, 0, "chainlink_feed_invalid"
	}
	if updatedAt == 0 {
		return false, 0, "chainlink_feed_not_updated"
	}
	if answeredInRound < roundID {
		return false, 0, "chainlink_round_incomplete"
	}
	if (now - updatedAt) > int64(cfg.OracleStalenessSeconds) {
		return false, 0, "chainlink_feed_stale"
	}
	return true, uint64(answer), ""
}

// Gate 5: Price Deviation vs Oracle
func checkPriceDeviation(chainlinkPrice uint64, marketPrice float64, cfg Config) (bool, string) {
	if chainlinkPrice == 0 {
		return true, "" // skip if no oracle price
	}
	marketPrice8d := toFeedDecimals(marketPrice, cfg.FeedDecimals)
	clPrice := int64(chainlinkPrice)
	mkPrice := int64(marketPrice8d)

	diff := clPrice - mkPrice
	if diff < 0 {
		diff = -diff
	}
	deviationBps := diff * 10000 / clPrice

	if deviationBps > int64(cfg.PriceDeviationMaxBPS) {
		return false, "price_deviation_exceeds_threshold"
	}
	return true, ""
}

// Gate 6: Volatility-Adjusted Position Sizing
// Does not deny — constrains position size and returns it.
func calculatePositionSize(requestedPosition float64, volatility float64, riskScore int, cfg Config) uint64 {
	absVol := math.Abs(volatility)
	volatilityFactor := clamp(1.0-(absVol/100.0*cfg.VolatilityScaleFactor), 0.1, 1.0)
	riskFactor := clamp(1.0-(float64(riskScore)/100.0), 0.1, 1.0)

	dynamicPosition := requestedPosition * volatilityFactor * riskFactor
	bpsCap := requestedPosition * float64(cfg.DefaultMaxPositionBPS) / 10000.0
	finalPosition := math.Min(dynamicPosition, math.Min(bpsCap, requestedPosition))

	return uint64(finalPosition)
}

// Gate 7: Hold Signal Filter
func checkHoldSignal(req RiskRequest) (bool, string) {
	if req.Signal == "hold" {
		return false, "hold_signal_no_trade"
	}
	return true, ""
}

// Gate 8: Agent Heartbeat Circuit Breaker
func checkAgentHeartbeat(cfg Config, heartbeatTimestamp int64, now int64) (bool, string) {
	if !cfg.EnableHeartbeatGate {
		return true, "" // skip when disabled
	}
	if heartbeatTimestamp == 0 {
		return false, "agent_heartbeat_stale"
	}
	if (now - heartbeatTimestamp) > int64(cfg.HeartbeatTTLSeconds) {
		return false, "agent_heartbeat_stale"
	}
	return true, ""
}
