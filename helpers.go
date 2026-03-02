package main

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/crypto"
)

// generateRunID produces a unique run identifier by hashing
// taskID, agentID, and a nonce together with keccak256.
func generateRunID(taskID, agentID string, nonce int64) [32]byte {
	data := []byte(fmt.Sprintf("%s:%s:%d", taskID, agentID, nonce))
	return crypto.Keccak256Hash(data)
}

// hashDecision produces a keccak256 hash of all decision fields
// for on-chain verification.
func hashDecision(d RiskDecision) [32]byte {
	buf := make([]byte, 0, 256)
	buf = append(buf, d.RunID[:]...)
	if d.Approved {
		buf = append(buf, 1)
	} else {
		buf = append(buf, 0)
	}

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, d.MaxPositionUSD)
	buf = append(buf, b...)

	binary.BigEndian.PutUint64(b, d.MaxSlippageBps)
	buf = append(buf, b...)

	binary.BigEndian.PutUint64(b, d.TTLSeconds)
	buf = append(buf, b...)

	buf = append(buf, []byte(d.Reason)...)

	binary.BigEndian.PutUint64(b, d.ChainlinkPrice)
	buf = append(buf, b...)

	binary.BigEndian.PutUint64(b, uint64(d.Timestamp))
	buf = append(buf, b...)

	return crypto.Keccak256Hash(buf)
}

// calculateSlippage returns slippage in basis points based on volatility.
// Higher volatility results in wider slippage tolerance.
func calculateSlippage(volatility float64, scaleFactor float64) uint64 {
	absVol := math.Abs(volatility)
	bps := uint64(math.Round(absVol * scaleFactor * 100))
	if bps < 10 {
		bps = 10 // minimum 10 bps
	}
	if bps > 500 {
		bps = 500 // maximum 500 bps
	}
	return bps
}

// toFeedDecimals converts a float price to integer with specified
// decimal precision. At 8 decimals, uint64 supports prices up to ~184 billion.
func toFeedDecimals(price float64, decimals int) uint64 {
	return uint64(math.Round(price * math.Pow(10, float64(decimals))))
}

// clamp bounds value to [min, max] range.
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
