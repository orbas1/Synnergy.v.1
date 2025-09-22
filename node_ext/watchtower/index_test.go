package watchtower

import (
	"testing"
	"time"
)

func TestValidateMetrics(t *testing.T) {
	now := time.Now().UTC()
	metrics := Metrics{
		CPUUsage:        42.1,
		MemoryUsage:     512 << 20,
		PeerCount:       12,
		LastBlockHeight: 1024,
		Timestamp:       now.Add(-time.Minute),
	}
	if err := ValidateMetrics(now, metrics); err != nil {
		t.Fatalf("expected valid metrics, got %v", err)
	}

	badCPU := metrics
	badCPU.CPUUsage = 140
	if err := ValidateMetrics(now, badCPU); err == nil {
		t.Fatalf("expected cpu usage validation error")
	}

	stale := metrics
	stale.Timestamp = now.Add(-6 * time.Minute)
	if err := ValidateMetrics(now, stale); err != ErrStaleMetrics {
		t.Fatalf("expected stale metrics error, got %v", err)
	}

	missingHeight := metrics
	missingHeight.LastBlockHeight = 0
	if err := ValidateMetrics(now, missingHeight); err == nil {
		t.Fatalf("expected block height validation error")
	}
}
