package watchtower

import (
	"testing"
	"time"
)

func TestValidateMetrics(t *testing.T) {
	m := Metrics{CPUUsage: 50, PeerCount: 5, Timestamp: time.Now()}
	if err := ValidateMetrics(m); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := ValidateMetrics(Metrics{CPUUsage: 150}); err == nil {
		t.Fatalf("expected cpu validation error")
	}
	if err := ValidateMetrics(Metrics{CPUUsage: 10, PeerCount: -1}); err == nil {
		t.Fatalf("expected peer count validation error")
	}
}
