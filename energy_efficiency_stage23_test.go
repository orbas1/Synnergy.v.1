package synnergy

import (
	"math"
	"testing"
	"time"
)

func TestEnergyEfficiencyTracker(t *testing.T) {
	tracker := NewEnergyEfficiencyTracker()

	// Initially efficiency and stats should report zero/false
	if eff, ok := tracker.Efficiency("v1"); ok || eff != 0 {
		t.Fatalf("expected no metrics initially, got %v %v", eff, ok)
	}
	if avg := tracker.NetworkAverage(); avg != 0 {
		t.Fatalf("expected zero network average, got %v", avg)
	}

	// Record metrics for two validators
	tracker.Record("v1", 100, 10) // 10 tx/kWh
	tracker.Record("v2", 50, 5)   // 10 tx/kWh
	tracker.Record("v1", 50, 5)   // total v1 -> 150 tx, 15 kWh -> 10

	eff, ok := tracker.Efficiency("v1")
	if !ok || math.Abs(eff-10) > 1e-9 {
		t.Fatalf("unexpected efficiency for v1: %v %v", eff, ok)
	}
	eff, ok = tracker.Efficiency("v2")
	if !ok || math.Abs(eff-10) > 1e-9 {
		t.Fatalf("unexpected efficiency for v2: %v %v", eff, ok)
	}

	if avg := tracker.NetworkAverage(); math.Abs(avg-10) > 1e-9 {
		t.Fatalf("network average mismatch: %v", avg)
	}

	// Stats should reflect aggregated values
	if stats, ok := tracker.Stats("v1"); !ok || stats.Transactions != 150 || math.Abs(stats.EnergyKWh-15) > 1e-9 {
		t.Fatalf("stats mismatch for v1: %+v %v", stats, ok)
	}

	// Reset should remove metrics
	tracker.Reset("v1")
	if _, ok := tracker.Stats("v1"); ok {
		t.Fatalf("expected v1 stats to be removed")
	}
	if eff, ok := tracker.Efficiency("v1"); ok || eff != 0 {
		t.Fatalf("expected no efficiency after reset, got %v %v", eff, ok)
	}
}

func TestEnergyEfficientNode(t *testing.T) {
	tracker := NewEnergyEfficiencyTracker()
	node := NewEnergyEfficientNode("nodeA", tracker)

	if node.ID() != "nodeA" {
		t.Fatalf("unexpected node id")
	}

	// Without records ShouldThrottle should report true because no data
	if !node.ShouldThrottle(1) {
		t.Fatalf("expected throttle when no data present")
	}

	node.RecordUsage(100, 20) // efficiency 5
	node.AddOffset(2.5)

	if c := node.OffsetCredits(); math.Abs(c-2.5) > 1e-9 {
		t.Fatalf("offset credits mismatch: %v", c)
	}

	// Certify should produce certificate with matching fields
	cert := node.Certify()
	if cert.Validator != "nodeA" {
		t.Fatalf("cert validator mismatch")
	}
	if math.Abs(cert.Efficiency-5) > 1e-9 {
		t.Fatalf("certificate efficiency mismatch: %v", cert.Efficiency)
	}
	if math.Abs(cert.Offsets-2.5) > 1e-9 {
		t.Fatalf("certificate offsets mismatch: %v", cert.Offsets)
	}
	if time.Since(cert.IssuedAt) > time.Second {
		t.Fatalf("certificate issued time not recent: %v", cert.IssuedAt)
	}
	// Certificate should be stored and retrievable
	stored := node.Certificate()
	if stored != cert {
		t.Fatalf("stored certificate mismatch")
	}

	// After recording usage efficiency is 5 tx/kWh.
	if node.ShouldThrottle(4) {
		t.Fatalf("should not throttle when threshold below efficiency")
	}
	if !node.ShouldThrottle(6) {
		t.Fatalf("should throttle when threshold above efficiency")
	}
}
