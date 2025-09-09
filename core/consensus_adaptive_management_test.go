package core

import (
	"math"
	"testing"
)

func TestAdaptiveManagerWindowAndReset(t *testing.T) {
	sc := NewSynnergyConsensus()
	am := NewAdaptiveManager(sc, 2)

	if th := am.Threshold(1, 0); th <= 0 {
		t.Fatalf("expected positive threshold, got %v", th)
	}
	if th := am.Threshold(0, 1); math.Abs(th-0.5) > 0.01 {
		t.Fatalf("expected avg threshold ~0.5, got %v", th)
	}
	if th := am.Threshold(0, 0); math.Abs(th-0.25) > 0.01 {
		t.Fatalf("expected sliding threshold ~0.25, got %v", th)
	}

	am.Reset()
	if th := am.Threshold(0, 0); th != 0 {
		t.Fatalf("expected 0 after reset, got %v", th)
	}

	w1 := am.Adjust(1, 0)
	w2 := am.Adjust(0, 1)
	if w1 == w2 {
		t.Fatalf("expected weights to change with new metrics")
	}
}
