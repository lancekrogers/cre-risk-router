package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lancekrogers/cre-risk-router/pkg/riskeval"
)

func TestEvaluateHandler_Approved(t *testing.T) {
	cfg := loadConfig()
	log := slog.Default()
	handler := newEvaluateHandler(cfg, log)

	req := riskeval.RiskRequest{
		AgentID:           "agent-1",
		TaskID:            "task-1",
		Signal:            "buy",
		SignalConfidence:  0.85,
		RiskScore:         10,
		MarketPair:        "ETH/USD",
		RequestedPosition: 1000_000000,
		Timestamp:         time.Now().Unix(),
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/evaluate", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var decision riskeval.RiskDecision
	if err := json.NewDecoder(rec.Body).Decode(&decision); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if !decision.Approved {
		t.Fatalf("expected approved, got denied: %s", decision.Reason)
	}
	if decision.MaxPositionUSD == 0 {
		t.Error("MaxPositionUSD should be non-zero")
	}
	if decision.Reason != "approved" {
		t.Errorf("reason = %q, want approved", decision.Reason)
	}
}

func TestEvaluateHandler_Denied(t *testing.T) {
	cfg := loadConfig()
	log := slog.Default()
	handler := newEvaluateHandler(cfg, log)

	req := riskeval.RiskRequest{
		AgentID:           "agent-1",
		TaskID:            "task-1",
		Signal:            "hold",
		SignalConfidence:  0.85,
		RiskScore:         10,
		MarketPair:        "ETH/USD",
		RequestedPosition: 1000_000000,
		Timestamp:         time.Now().Unix(),
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/evaluate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var decision riskeval.RiskDecision
	if err := json.NewDecoder(rec.Body).Decode(&decision); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if decision.Approved {
		t.Fatal("hold signal should be denied")
	}
	if decision.Reason != "hold_signal_no_trade" {
		t.Errorf("reason = %q", decision.Reason)
	}
}

func TestEvaluateHandler_InvalidBody(t *testing.T) {
	cfg := loadConfig()
	log := slog.Default()
	handler := newEvaluateHandler(cfg, log)

	httpReq := httptest.NewRequest(http.MethodPost, "/evaluate", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}
