package main

// Config is deserialized from config.staging.json or config.production.json.
type Config struct {
	MarketDataURL             string  `json:"market_data_url"`
	PriceFeedAddress          string  `json:"price_feed_address"`
	ReceiptContractAddress    string  `json:"receipt_contract_address"`
	TargetNetwork             string  `json:"target_network"`
	SignalConfidenceThreshold float64 `json:"signal_confidence_threshold"`
	MaxRiskScore              int     `json:"max_risk_score"`
	DefaultMaxPositionBPS     int     `json:"default_max_position_bps"`
	DecisionTTLSeconds        int     `json:"decision_ttl_seconds"`
	PriceDeviationMaxBPS      int     `json:"price_deviation_max_bps"`
	VolatilityScaleFactor     float64 `json:"volatility_scale_factor"`
	OracleStalenessSeconds    int     `json:"oracle_staleness_seconds"`
	FeedDecimals              int     `json:"feed_decimals"`
	EnableHeartbeatGate       bool    `json:"enable_heartbeat_gate"`
	HeartbeatMirrorNodeURL    string  `json:"heartbeat_mirror_node_url"`
	HeartbeatTTLSeconds       int     `json:"heartbeat_ttl_seconds"`
}

// RiskDecision is the output of the risk evaluation pipeline.
type RiskDecision struct {
	Decision string `json:"decision"`
}
