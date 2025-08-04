package core

import "testing"

func TestAdaptiveManagerAdjust(t *testing.T) {
	sc := NewSynnergyConsensus()
	am := NewAdaptiveManager(sc)
	weights := am.Adjust(0.5, 0.5)
	if weights.PoW == 0 && weights.PoS == 0 && weights.PoH == 0 {
		t.Fatalf("expected weights to be adjusted")
	}
	if am.Threshold(0.3, 0.4) <= 0 {
		t.Fatalf("expected positive threshold")
	}
}
