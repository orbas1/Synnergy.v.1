package core

import "testing"

func TestDifficultyManager(t *testing.T) {
	sc := NewSynnergyConsensus()
	dm := NewDifficultyManager(sc, 3, 1, 10)
	d1 := dm.AddSample(8)
	d2 := dm.AddSample(12)
	if d1 == 0 || d2 == 0 {
		t.Fatalf("difficulty should adjust and remain non-zero")
	}
	if dm.Difficulty() != d2 {
		t.Fatalf("expected latest difficulty to be returned")
	}
}
