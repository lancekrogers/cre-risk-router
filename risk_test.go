package main

import (
	"testing"
)

func defaultConfig() Config {
	return Config{
		SignalConfidenceThreshold: 0.6,
		MaxRiskScore:              75,
		DecisionTTLSeconds:        300,
		PriceDeviationMaxBPS:      500,
		VolatilityScaleFactor:     1.0,
		DefaultMaxPositionBPS:     10000,
		OracleStalenessSeconds:    3600,
		FeedDecimals:              8,
		HeartbeatTTLSeconds:       600,
	}
}

func defaultOracle(now int64) OracleData {
	return OracleData{
		RoundID:         1,
		Answer:          200000000000, // $2000 at 8 decimals
		StartedAt:       now - 60,
		UpdatedAt:       now - 60,
		AnsweredInRound: 1,
	}
}

// --- Gate 7: Hold Signal Filter ---

func TestGate7_HoldSignalDenied(t *testing.T) {
	ok, reason := checkHoldSignal(RiskRequest{Signal: "hold"})
	if ok {
		t.Fatal("expected hold signal to be denied")
	}
	if reason != "hold_signal_no_trade" {
		t.Errorf("reason = %q, want hold_signal_no_trade", reason)
	}
}

func TestGate7_BuySellPass(t *testing.T) {
	for _, signal := range []string{"buy", "sell"} {
		ok, _ := checkHoldSignal(RiskRequest{Signal: signal})
		if !ok {
			t.Errorf("signal %q should pass gate 7", signal)
		}
	}
}

// --- Gate 1: Signal Confidence ---

func TestGate1_BelowThreshold(t *testing.T) {
	cfg := defaultConfig()
	ok, reason := checkSignalConfidence(RiskRequest{SignalConfidence: 0.5}, cfg)
	if ok {
		t.Fatal("expected low confidence to be denied")
	}
	if reason != "signal_confidence_below_threshold" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate1_AtThreshold(t *testing.T) {
	cfg := defaultConfig()
	ok, _ := checkSignalConfidence(RiskRequest{SignalConfidence: 0.6}, cfg)
	if !ok {
		t.Fatal("confidence at threshold should pass")
	}
}

func TestGate1_AboveThreshold(t *testing.T) {
	cfg := defaultConfig()
	ok, _ := checkSignalConfidence(RiskRequest{SignalConfidence: 0.9}, cfg)
	if !ok {
		t.Fatal("high confidence should pass")
	}
}

// --- Gate 2: Risk Score Ceiling ---

func TestGate2_ExceedsMaximum(t *testing.T) {
	cfg := defaultConfig()
	ok, reason := checkRiskScore(RiskRequest{RiskScore: 80}, cfg)
	if ok {
		t.Fatal("expected high risk score to be denied")
	}
	if reason != "risk_score_exceeds_maximum" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate2_AtMaximum(t *testing.T) {
	cfg := defaultConfig()
	ok, _ := checkRiskScore(RiskRequest{RiskScore: 75}, cfg)
	if !ok {
		t.Fatal("risk score at max should pass")
	}
}

func TestGate2_BelowMaximum(t *testing.T) {
	cfg := defaultConfig()
	ok, _ := checkRiskScore(RiskRequest{RiskScore: 10}, cfg)
	if !ok {
		t.Fatal("low risk score should pass")
	}
}

// --- Gate 3: Signal Staleness ---

func TestGate3_FreshSignal(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, _ := checkSignalStaleness(RiskRequest{Timestamp: now - 100}, cfg, now)
	if !ok {
		t.Fatal("100s-old signal should pass with 300s TTL")
	}
}

func TestGate3_StaleSignal(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, reason := checkSignalStaleness(RiskRequest{Timestamp: now - 600}, cfg, now)
	if ok {
		t.Fatal("expected stale signal to be denied")
	}
	if reason != "signal_expired" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate3_ExactTTL(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, _ := checkSignalStaleness(RiskRequest{Timestamp: now - 300}, cfg, now)
	if !ok {
		t.Fatal("signal at exact TTL boundary should pass")
	}
}

// --- Gate 4: Oracle Health ---

func TestGate4_HealthyOracle(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, price, _ := checkOracleHealth(1, 200000000000, now-60, now-60, 1, cfg, now)
	if !ok {
		t.Fatal("healthy oracle should pass")
	}
	if price != 200000000000 {
		t.Errorf("price = %d, want 200000000000", price)
	}
}

func TestGate4_InvalidAnswer(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, _, reason := checkOracleHealth(1, 0, now-60, now-60, 1, cfg, now)
	if ok {
		t.Fatal("zero answer should be denied")
	}
	if reason != "chainlink_feed_invalid" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate4_NegativeAnswer(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, _, reason := checkOracleHealth(1, -1, now-60, now-60, 1, cfg, now)
	if ok {
		t.Fatal("negative answer should be denied")
	}
	if reason != "chainlink_feed_invalid" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate4_NotUpdated(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, _, reason := checkOracleHealth(1, 200000000000, now-60, 0, 1, cfg, now)
	if ok {
		t.Fatal("zero updatedAt should be denied")
	}
	if reason != "chainlink_feed_not_updated" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate4_RoundIncomplete(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, _, reason := checkOracleHealth(5, 200000000000, now-60, now-60, 4, cfg, now)
	if ok {
		t.Fatal("incomplete round should be denied")
	}
	if reason != "chainlink_round_incomplete" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate4_StaleFeed(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	ok, _, reason := checkOracleHealth(1, 200000000000, now-7200, now-7200, 1, cfg, now)
	if ok {
		t.Fatal("stale feed should be denied")
	}
	if reason != "chainlink_feed_stale" {
		t.Errorf("reason = %q", reason)
	}
}

// --- Gate 5: Price Deviation ---

func TestGate5_WithinThreshold(t *testing.T) {
	cfg := defaultConfig()
	ok, _ := checkPriceDeviation(200000000000, 2010.0, cfg)
	if !ok {
		t.Fatal("0.5% deviation should pass with 5% threshold")
	}
}

func TestGate5_ExceedsThreshold(t *testing.T) {
	cfg := defaultConfig()
	ok, reason := checkPriceDeviation(200000000000, 2200.0, cfg)
	if ok {
		t.Fatal("10% deviation should be denied")
	}
	if reason != "price_deviation_exceeds_threshold" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate5_SkippedWhenNoPriceOracle(t *testing.T) {
	cfg := defaultConfig()
	ok, _ := checkPriceDeviation(0, 2000.0, cfg)
	if !ok {
		t.Fatal("should skip when chainlinkPrice is 0")
	}
}

// --- Gate 6: Position Sizing ---

func TestGate6_LowVolatilityLowRisk(t *testing.T) {
	cfg := defaultConfig()
	pos := calculatePositionSize(1000_000000, 3.0, 10, cfg)
	// volatilityFactor = 1.0 - (3.0/100.0 * 1.0) = 0.97
	// riskFactor = 1.0 - (10/100.0) = 0.90
	// dynamic = 1000000000 * 0.97 * 0.90 = 873000000
	if pos != 873000000 {
		t.Errorf("pos = %d, want 873000000", pos)
	}
}

func TestGate6_HighVolatilityHighRisk(t *testing.T) {
	cfg := defaultConfig()
	pos := calculatePositionSize(1000_000000, 50.0, 70, cfg)
	// volatilityFactor = 1.0 - (50.0/100.0 * 1.0) = 0.50
	// riskFactor = 1.0 - (70/100.0) = 0.30
	// dynamic = 1000000000 * 0.50 * 0.30 = 150000000
	if pos != 150000000 {
		t.Errorf("pos = %d, want 150000000", pos)
	}
}

func TestGate6_ExtremeVolatilityClamped(t *testing.T) {
	cfg := defaultConfig()
	pos := calculatePositionSize(1000_000000, 200.0, 99, cfg)
	// volatilityFactor clamped to 0.1, riskFactor clamped to 0.1
	// dynamic = 1000000000 * 0.1 * 0.1 = 10000000
	if pos != 10000000 {
		t.Errorf("pos = %d, want 10000000", pos)
	}
}

// --- Gate 8: Heartbeat ---

func TestGate8_DisabledSkips(t *testing.T) {
	cfg := defaultConfig()
	cfg.EnableHeartbeatGate = false
	ok, _ := checkAgentHeartbeat(cfg, 0, 1000000)
	if !ok {
		t.Fatal("disabled heartbeat gate should pass")
	}
}

func TestGate8_EnabledNoHeartbeat(t *testing.T) {
	cfg := defaultConfig()
	cfg.EnableHeartbeatGate = true
	ok, reason := checkAgentHeartbeat(cfg, 0, 1000000)
	if ok {
		t.Fatal("zero heartbeat with gate enabled should deny")
	}
	if reason != "agent_heartbeat_stale" {
		t.Errorf("reason = %q", reason)
	}
}

func TestGate8_EnabledFreshHeartbeat(t *testing.T) {
	cfg := defaultConfig()
	cfg.EnableHeartbeatGate = true
	now := int64(1000000)
	ok, _ := checkAgentHeartbeat(cfg, now-100, now)
	if !ok {
		t.Fatal("fresh heartbeat should pass")
	}
}

func TestGate8_EnabledStaleHeartbeat(t *testing.T) {
	cfg := defaultConfig()
	cfg.EnableHeartbeatGate = true
	now := int64(1000000)
	ok, reason := checkAgentHeartbeat(cfg, now-1200, now)
	if ok {
		t.Fatal("stale heartbeat should deny")
	}
	if reason != "agent_heartbeat_stale" {
		t.Errorf("reason = %q", reason)
	}
}

// --- evaluateRisk integration ---

func TestEvaluateRisk_ApprovedLowRisk(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	req := RiskRequest{
		AgentID:           "agent-1",
		TaskID:            "task-1",
		Signal:            "buy",
		SignalConfidence:  0.85,
		RiskScore:         10,
		MarketPair:        "ETH/USD",
		RequestedPosition: 1000_000000,
		Timestamp:         now,
	}
	oracle := defaultOracle(now)

	d := evaluateRisk(req, nil, oracle, cfg, now, 0)
	if !d.Approved {
		t.Fatalf("expected approved, got denied: %s", d.Reason)
	}
	if d.MaxPositionUSD == 0 {
		t.Error("MaxPositionUSD should be non-zero for approved decision")
	}
	if d.ChainlinkPrice != 200000000000 {
		t.Errorf("ChainlinkPrice = %d, want 200000000000", d.ChainlinkPrice)
	}
	if d.RunID == [32]byte{} {
		t.Error("RunID should not be zero")
	}
	if d.DecisionHash == [32]byte{} {
		t.Error("DecisionHash should not be zero")
	}
}

func TestEvaluateRisk_DeniedHoldSignal(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	req := RiskRequest{
		Signal:           "hold",
		SignalConfidence: 0.85,
		RiskScore:        10,
		Timestamp:        now,
	}
	oracle := defaultOracle(now)

	d := evaluateRisk(req, nil, oracle, cfg, now, 0)
	if d.Approved {
		t.Fatal("hold signal should be denied")
	}
	if d.Reason != "hold_signal_no_trade" {
		t.Errorf("reason = %q", d.Reason)
	}
}

func TestEvaluateRisk_DeniedLowConfidence(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	req := RiskRequest{
		Signal:           "buy",
		SignalConfidence: 0.45,
		RiskScore:        30,
		Timestamp:        now,
	}
	oracle := defaultOracle(now)

	d := evaluateRisk(req, nil, oracle, cfg, now, 0)
	if d.Approved {
		t.Fatal("low confidence should be denied")
	}
	if d.Reason != "signal_confidence_below_threshold" {
		t.Errorf("reason = %q", d.Reason)
	}
}

func TestEvaluateRisk_DeniedHighRisk(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	req := RiskRequest{
		Signal:           "sell",
		SignalConfidence: 0.9,
		RiskScore:        82,
		Timestamp:        now,
	}
	oracle := defaultOracle(now)

	d := evaluateRisk(req, nil, oracle, cfg, now, 0)
	if d.Approved {
		t.Fatal("high risk score should be denied")
	}
	if d.Reason != "risk_score_exceeds_maximum" {
		t.Errorf("reason = %q", d.Reason)
	}
}

func TestEvaluateRisk_DeniedStaleSignal(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	req := RiskRequest{
		Signal:           "buy",
		SignalConfidence: 0.85,
		RiskScore:        10,
		Timestamp:        now - 600,
	}
	oracle := defaultOracle(now)

	d := evaluateRisk(req, nil, oracle, cfg, now, 0)
	if d.Approved {
		t.Fatal("stale signal should be denied")
	}
	if d.Reason != "signal_expired" {
		t.Errorf("reason = %q", d.Reason)
	}
}

func TestEvaluateRisk_DeniedBadOracle(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	req := RiskRequest{
		Signal:           "buy",
		SignalConfidence: 0.85,
		RiskScore:        10,
		Timestamp:        now,
	}
	oracle := OracleData{
		RoundID:         1,
		Answer:          0, // invalid
		StartedAt:       now - 60,
		UpdatedAt:       now - 60,
		AnsweredInRound: 1,
	}

	d := evaluateRisk(req, nil, oracle, cfg, now, 0)
	if d.Approved {
		t.Fatal("bad oracle should be denied")
	}
	if d.Reason != "chainlink_feed_invalid" {
		t.Errorf("reason = %q", d.Reason)
	}
}

func TestEvaluateRisk_GateOrderFirstDenyWins(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)

	// Request fails Gate 7 (hold), Gate 1 (low confidence), and Gate 2 (high risk).
	// Gate 7 should be the denial reason since it runs first.
	req := RiskRequest{
		Signal:           "hold",
		SignalConfidence: 0.1,
		RiskScore:        99,
		Timestamp:        now,
	}
	oracle := defaultOracle(now)

	d := evaluateRisk(req, nil, oracle, cfg, now, 0)
	if d.Approved {
		t.Fatal("should be denied")
	}
	if d.Reason != "hold_signal_no_trade" {
		t.Errorf("first gate (7) should deny, got reason = %q", d.Reason)
	}
}

func TestEvaluateRisk_UniqueRunIDs(t *testing.T) {
	cfg := defaultConfig()
	now := int64(1000000)
	oracle := defaultOracle(now)

	d1 := evaluateRisk(RiskRequest{
		AgentID: "a1", TaskID: "t1", Signal: "buy",
		SignalConfidence: 0.85, RiskScore: 10, Timestamp: now,
		RequestedPosition: 1000_000000,
	}, nil, oracle, cfg, now, 0)

	d2 := evaluateRisk(RiskRequest{
		AgentID: "a1", TaskID: "t2", Signal: "buy",
		SignalConfidence: 0.85, RiskScore: 10, Timestamp: now,
		RequestedPosition: 1000_000000,
	}, nil, oracle, cfg, now, 0)

	if d1.RunID == d2.RunID {
		t.Error("different tasks should produce different RunIDs")
	}
}
