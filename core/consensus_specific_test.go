package core

import "testing"

func TestConsensusSwitcher(t *testing.T) {
	sc := NewSynnergyConsensus()
	cs := NewConsensusSwitcher(ModePoW)
	sc.Weights = ConsensusWeights{PoW: 0.2, PoS: 0.6, PoH: 0.2}
	mode := cs.Evaluate(sc)
	if mode != ModePoS {
		t.Fatalf("expected ModePoS, got %s", mode)
	}
	if cs.Mode() != ModePoS {
		t.Fatalf("mode getter mismatch")
	}
}
