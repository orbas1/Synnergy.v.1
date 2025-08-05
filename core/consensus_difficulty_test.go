package core

import "testing"

// TestDifficultyManagerAdjustAndRetrieve ensures that adding samples adjusts the
// difficulty and that the Difficulty method returns the latest value.
func TestDifficultyManagerAdjustAndRetrieve(t *testing.T) {
	sc := NewSynnergyConsensus()
	dm := NewDifficultyManager(sc, 3, 1, 10)

	d1 := dm.AddSample(20) // new difficulty: 1*(20/10) = 2
	d2 := dm.AddSample(10) // avg(20,10)=15 -> 2*(15/10)=3

	if d1 != 2 || d2 != 3 {
		t.Fatalf("unexpected difficulties d1=%v d2=%v", d1, d2)
	}
	if dm.Difficulty() != d2 {
		t.Fatalf("Difficulty()=%v want %v", dm.Difficulty(), d2)
	}
}

// TestDifficultyManagerSlidingWindow verifies that the sliding window drops old
// samples and only considers the most recent ones.
func TestDifficultyManagerSlidingWindow(t *testing.T) {
	sc := NewSynnergyConsensus()
	dm := NewDifficultyManager(sc, 2, 1, 10)

	if diff := dm.AddSample(20); diff != 2 {
		t.Fatalf("diff1=%v want 2", diff)
	}
	if diff := dm.AddSample(10); diff != 3 {
		t.Fatalf("diff2=%v want 3", diff)
	}
	// Third sample should drop the first (20) from the window: avg= (10+10)/2=10
	if diff := dm.AddSample(10); diff != 3 {
		t.Fatalf("diff3=%v want 3", diff)
	}
}

// TestDifficultyManagerNilEngine ensures difficulty remains unchanged when no
// consensus engine is provided.
func TestDifficultyManagerNilEngine(t *testing.T) {
	dm := NewDifficultyManager(nil, 2, 5, 10)

	if diff := dm.AddSample(20); diff != 5 {
		t.Fatalf("expected diff 5 with nil engine, got %v", diff)
	}
	if diff := dm.AddSample(5); diff != 5 {
		t.Fatalf("expected diff 5 with nil engine, got %v", diff)
	}
}

// TestDifficultyManagerWindowFloor verifies that non-positive window sizes are
// floored to one, so only the latest sample affects the difficulty.
func TestDifficultyManagerWindowFloor(t *testing.T) {
	sc := NewSynnergyConsensus()
	dm := NewDifficultyManager(sc, 0, 1, 10) // window should become 1

	if diff := dm.AddSample(20); diff != 2 {
		t.Fatalf("first diff=%v want 2", diff)
	}
	// If window were >1, diff would become 3; with window=1 it stays 2.
	if diff := dm.AddSample(10); diff != 2 {
		t.Fatalf("expected diff 2 with window floor, got %v", diff)
	}
}
