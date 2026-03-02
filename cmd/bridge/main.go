// Command bridge runs a standalone HTTP server that wraps the CRE risk
// evaluation pipeline. This bridges the gap between the coordinator's
// creclient (which makes HTTP POST calls) and the CRE risk evaluation
// logic (which normally runs as a CRE WASM workflow).
//
// The coordinator sends POST requests with a RiskRequest JSON body and
// expects a RiskDecision JSON response. This bridge provides that endpoint
// using the same evaluation logic as the CRE workflow.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lancekrogers/cre-risk-router/pkg/riskeval"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := loadConfig()
	addr := envOr("BRIDGE_ADDR", ":8080")

	mux := http.NewServeMux()
	mux.Handle("POST /evaluate", newEvaluateHandler(cfg, log))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Info("CRE bridge starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}

func newEvaluateHandler(cfg riskeval.Config, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req riskeval.RiskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Warn("invalid request body", "error", err)
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		now := time.Now().Unix()

		// Use mock oracle data consistent with the CRE workflow.
		// In production, this would fetch from a Chainlink price feed.
		oracle := riskeval.OracleData{
			RoundID:         1,
			Answer:          200000000000, // $2000 at 8 decimals
			StartedAt:       now - 60,
			UpdatedAt:       now - 60,
			AnsweredInRound: 1,
		}

		decision := riskeval.EvaluateRisk(req, nil, oracle, cfg, now, 0)

		log.Info("risk evaluation",
			"task_id", req.TaskID,
			"agent_id", req.AgentID,
			"approved", decision.Approved,
			"reason", decision.Reason,
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(decision)
	})
}

func loadConfig() riskeval.Config {
	return riskeval.Config{
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
}

func envOr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
