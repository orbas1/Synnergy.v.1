package synnergy

import (
	"sync"
	"testing"
)

// TestConsensusHopperInitial verifies that the hopper respects the initial
// mode and allows explicit overrides via SetMode.
func TestConsensusHopperInitial(t *testing.T) {
	hopper := NewConsensusHopper(ConsensusPoW)
	if hopper.Mode() != ConsensusPoW {
		t.Fatalf("expected initial mode %s, got %s", ConsensusPoW, hopper.Mode())
	}
	hopper.SetMode(ConsensusPoS)
	if hopper.Mode() != ConsensusPoS {
		t.Fatalf("expected mode after SetMode %s, got %s", ConsensusPoS, hopper.Mode())
	}
}

// TestConsensusHopperEvaluate ensures that Evaluate selects the correct
// consensus mode for each branch and records the last metrics supplied.
func TestConsensusHopperEvaluate(t *testing.T) {
	hopper := NewConsensusHopper(ConsensusPoW)

	// High TPS and low latency should choose PoS regardless of validators.
	m1 := NetworkMetrics{TPS: 2000, LatencySec: 0.5, Validators: 5}
	if mode := hopper.Evaluate(m1); mode != ConsensusPoS {
		t.Fatalf("expected PoS, got %s", mode)
	}
	if last := hopper.LastMetrics(); last != m1 {
		t.Fatalf("last metrics mismatch: %+v vs %+v", last, m1)
	}

	// Few validators triggers PoH when TPS/latency aren't in the PoS range.
	m2 := NetworkMetrics{TPS: 100, LatencySec: 2, Validators: 3}
	if mode := hopper.Evaluate(m2); mode != ConsensusPoH {
		t.Fatalf("expected PoH, got %s", mode)
	}
	if last := hopper.LastMetrics(); last != m2 {
		t.Fatalf("last metrics mismatch: %+v vs %+v", last, m2)
	}

	// Default case falls back to PoW.
	m3 := NetworkMetrics{TPS: 500, LatencySec: 1.5, Validators: 20}
	if mode := hopper.Evaluate(m3); mode != ConsensusPoW {
		t.Fatalf("expected PoW, got %s", mode)
	}
	if last := hopper.LastMetrics(); last != m3 {
		t.Fatalf("last metrics mismatch: %+v vs %+v", last, m3)
	}
}

// TestConsensusHopperConcurrency performs simple concurrent evaluations to
// ensure the hopper's locking protects internal state. The final mode should be
// one of the evaluated results and the test should run without data races when
// executed with the -race flag.
func TestConsensusHopperConcurrency(t *testing.T) {
	hopper := NewConsensusHopper(ConsensusPoW)
	metrics := []NetworkMetrics{
		{TPS: 2000, LatencySec: 0.4, Validators: 50}, // PoS
		{TPS: 100, LatencySec: 3, Validators: 5},     // PoH
		{TPS: 200, LatencySec: 2, Validators: 20},    // PoW
	}

	var wg sync.WaitGroup
	for _, m := range metrics {
		wg.Add(1)
		go func(m NetworkMetrics) {
			defer wg.Done()
			hopper.Evaluate(m)
		}(m)
	}
	wg.Wait()

	final := hopper.Mode()
	if final != ConsensusPoS && final != ConsensusPoH && final != ConsensusPoW {
		t.Fatalf("unexpected final mode: %s", final)
	}
}
