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

// OracleData holds the 5-tuple from Chainlink latestRoundData().
type OracleData struct {
	RoundID          int64
	Answer           int64
	StartedAt        int64
	UpdatedAt        int64
	AnsweredInRound  int64
}

// evaluateRisk runs all active gates sequentially and produces a RiskDecision.
// Gate order: 7 (hold), 1 (confidence), 2 (risk score), 3 (staleness),
// 4 (oracle health), 5 (price deviation), 6 (position sizing), 8 (heartbeat).
func evaluateRisk(
	req RiskRequest,
	market *MarketData,
	oracle OracleData,
	cfg Config,
	now int64,
	heartbeatTimestamp int64,
) RiskDecision {
	makeDenied := func(reason string, chainlinkPrice uint64) RiskDecision {
		d := RiskDecision{
			Approved:       false,
			MaxPositionUSD: 0,
			MaxSlippageBps: 0,
			TTLSeconds:     uint64(cfg.DecisionTTLSeconds),
			Reason:         reason,
			ChainlinkPrice: chainlinkPrice,
			Timestamp:      now,
		}
		d.RunID = generateRunID(req.TaskID, req.AgentID, now)
		d.DecisionHash = hashDecision(d)
		return d
	}

	// Gate 7: Hold signal (fast path)
	if ok, reason := checkHoldSignal(req); !ok {
		return makeDenied(reason, 0)
	}

	// Gate 1: Signal confidence
	if ok, reason := checkSignalConfidence(req, cfg); !ok {
		return makeDenied(reason, 0)
	}

	// Gate 2: Risk score
	if ok, reason := checkRiskScore(req, cfg); !ok {
		return makeDenied(reason, 0)
	}

	// Gate 3: Signal staleness
	if ok, reason := checkSignalStaleness(req, cfg, now); !ok {
		return makeDenied(reason, 0)
	}

	// Gate 4: Oracle health
	oracleOk, chainlinkPrice, oracleReason := checkOracleHealth(
		oracle.RoundID, oracle.Answer, oracle.StartedAt, oracle.UpdatedAt, oracle.AnsweredInRound,
		cfg, now,
	)
	if !oracleOk {
		return makeDenied(oracleReason, 0)
	}

	// Gate 5: Price deviation (skip if no market data)
	if market != nil && market.Price > 0 {
		if ok, reason := checkPriceDeviation(chainlinkPrice, market.Price, cfg); !ok {
			return makeDenied(reason, chainlinkPrice)
		}
	}

	// Gate 6: Position sizing (never denies, constrains position)
	volatility := 10.0 // fallback volatility
	if market != nil {
		volatility = market.Volatility24h
	}
	maxPosition := calculatePositionSize(req.RequestedPosition, volatility, req.RiskScore, cfg)
	slippageBps := calculateSlippage(volatility, cfg.VolatilityScaleFactor)

	// Gate 8: Heartbeat (skip if disabled)
	if ok, reason := checkAgentHeartbeat(cfg, heartbeatTimestamp, now); !ok {
		return makeDenied(reason, chainlinkPrice)
	}

	// All gates passed — approved
	d := RiskDecision{
		Approved:       true,
		MaxPositionUSD: maxPosition,
		MaxSlippageBps: slippageBps,
		TTLSeconds:     uint64(cfg.DecisionTTLSeconds),
		Reason:         "approved",
		ChainlinkPrice: chainlinkPrice,
		Timestamp:      now,
	}
	d.RunID = generateRunID(req.TaskID, req.AgentID, now)
	d.DecisionHash = hashDecision(d)
	return d
}
