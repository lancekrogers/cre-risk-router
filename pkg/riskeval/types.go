package riskeval

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

// RiskRequest is the input to the risk evaluation pipeline.
type RiskRequest struct {
	AgentID           string  `json:"agent_id"`
	TaskID            string  `json:"task_id"`
	Signal            string  `json:"signal"`             // buy, sell, hold
	SignalConfidence  float64 `json:"signal_confidence"`   // 0.0–1.0
	RiskScore         int     `json:"risk_score"`          // 0–100
	MarketPair        string  `json:"market_pair"`
	RequestedPosition float64 `json:"requested_position"`
	Timestamp         int64   `json:"timestamp"`           // Unix seconds
}

// RiskDecision is the output of the risk evaluation pipeline.
type RiskDecision struct {
	RunID          [32]byte `json:"run_id"`
	DecisionHash   [32]byte `json:"decision_hash"`
	Approved       bool     `json:"approved"`
	MaxPositionUSD uint64   `json:"max_position_usd"`
	MaxSlippageBps uint64   `json:"max_slippage_bps"`
	TTLSeconds     uint64   `json:"ttl_seconds"`
	Reason         string   `json:"reason"`
	ChainlinkPrice uint64   `json:"chainlink_price"` // 8-decimal precision
	Timestamp      int64    `json:"timestamp"`
}

// MarketData holds data from the CoinGecko API response.
type MarketData struct {
	Price         float64 `json:"current_price"`
	Volume24h     float64 `json:"total_volume"`
	Volatility24h float64 `json:"price_change_percentage_24h"`
	MarketCap     float64 `json:"market_cap"`
}

// OracleData holds the 5-tuple from Chainlink latestRoundData().
type OracleData struct {
	RoundID         int64
	Answer          int64
	StartedAt       int64
	UpdatedAt       int64
	AnsweredInRound int64
}
