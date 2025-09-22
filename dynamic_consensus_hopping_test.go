package synnergy

import (
	"sync"
	"testing"
	"time"
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
	stake := NewStakingNode()
	stake.Stake("validatorA", 120000)
	stake.Stake("validatorB", 60000)
	hopper.AttachStakingNode(stake)

	failover := NewFailoverManager("primary", 200*time.Millisecond)
	failover.RegisterBackup("backup")
	hopper.AttachFailoverManager(failover)

	// High throughput, deep staking and stable failover should favour PoS.
	m1 := NetworkMetrics{TPS: 1800, LatencySec: 0.6, Validators: 80, FinalityLagSec: 1.2, ForkRate: 0.01, QueueDepth: 200}
	if mode := hopper.Evaluate(m1); mode != ConsensusPoS {
		t.Fatalf("expected PoS, got %s", mode)
	}
	if last := hopper.LastMetrics(); last != m1 {
		t.Fatalf("last metrics mismatch: %+v vs %+v", last, m1)
	}

	// Validator scarcity with finality instability should switch to PoH.
	m2 := NetworkMetrics{TPS: 600, LatencySec: 2.7, Validators: 6, FinalityLagSec: 4.5, ForkRate: 0.09, QueueDepth: 2200}
	if mode := hopper.Evaluate(m2); mode != ConsensusPoH {
		t.Fatalf("expected PoH, got %s", mode)
	}
	if last := hopper.LastMetrics(); last != m2 {
		t.Fatalf("last metrics mismatch: %+v vs %+v", last, m2)
	}

	// Simulate a recent failover promotion which should push the hopper back
	// to PoW until the network settles.
	timeout := 50 * time.Millisecond
	unstableFailover := NewFailoverManager("primary", timeout)
	unstableFailover.RegisterBackup("backup")
	hopper.AttachFailoverManager(unstableFailover)
	unstableFailover.mu.Lock()
	unstableFailover.nodes["primary"] = time.Now().Add(-2 * timeout)
	unstableFailover.mu.Unlock()
	if active := unstableFailover.Active(); active != "backup" {
		t.Fatalf("expected failover to promote backup, got %s", active)
	}
	unstableFailover.Heartbeat("backup")

	m3 := NetworkMetrics{TPS: 800, LatencySec: 1.4, Validators: 40, FinalityLagSec: 2, ForkRate: 0.02, QueueDepth: 600}
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
		{TPS: 2000, LatencySec: 0.4, Validators: 50, FinalityLagSec: 1, ForkRate: 0.01, QueueDepth: 300}, // PoS
		{TPS: 100, LatencySec: 3, Validators: 5, FinalityLagSec: 4.2, ForkRate: 0.08, QueueDepth: 2500},  // PoH
		{TPS: 200, LatencySec: 2, Validators: 20, FinalityLagSec: 2.5, ForkRate: 0.03, QueueDepth: 1000}, // PoW
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
