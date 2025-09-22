package core

import "testing"

// TestNewConsensusSpecificNode ensures the node locks to a single consensus mode.

func TestNewConsensusSpecificNode(t *testing.T) {
	ledger := NewLedger()
	n := NewConsensusSpecificNode(ModePoS, "n1", "addr", ledger)
	weights := n.Consensus.WeightsSnapshot()
	if !n.Consensus.PoSAvailable || weights.PoS != 1 {
		t.Fatalf("expected PoS only mode")
	}
	if n.Consensus.PoWAvailable || n.Consensus.PoHAvailable {
		t.Fatalf("other modes should be disabled")
	}
}
