package core

import (
    "math"
    "testing"
)

func TestThreshold(t *testing.T) {
    sc := NewSynnergyConsensus()
    if sc.Threshold(2, 3) != sc.Alpha*2+sc.Beta*3 {
        t.Fatalf("threshold calculation incorrect")
    }
}

func TestAdjustWeightsAndAvailability(t *testing.T) {
    sc := NewSynnergyConsensus()
    sc.SetAvailability(true, false, true)
    sc.AdjustWeights(0.5, 0.5)
    if sc.Weights.PoS != 0 {
        t.Fatalf("PoS weight should be zero when unavailable")
    }
    total := sc.Weights.PoW + sc.Weights.PoS + sc.Weights.PoH
    if math.Abs(total-1) > 1e-9 {
        t.Fatalf("weights not normalized")
    }
}

func TestTransitionThreshold(t *testing.T) {
    sc := NewSynnergyConsensus()
    tt := sc.TransitionThreshold(1, 2, 3)
    expected := sc.Tload(1) + sc.Tsecurity(2) + sc.Tstake(3)
    if tt != expected {
        t.Fatalf("transition threshold mismatch")
    }
}

func TestDifficultyAdjust(t *testing.T) {
    sc := NewSynnergyConsensus()
    if sc.DifficultyAdjust(1, 20, 10) != 2 {
        t.Fatalf("difficulty adjust incorrect")
    }
}

func TestSelectValidator(t *testing.T) {
    sc := NewSynnergyConsensus()
    stakes := map[string]uint64{"a": 1, "b": 2}
    addr := sc.SelectValidator(stakes)
    if addr != "a" && addr != "b" {
        t.Fatalf("unexpected validator: %s", addr)
    }
    if sc.SelectValidator(map[string]uint64{}) != "" {
        t.Fatalf("expected empty string when no stakes")
    }
}

