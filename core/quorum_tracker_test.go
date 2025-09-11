package core

import "testing"

func TestQuorumTracker(t *testing.T) {
	qt := NewQuorumTracker(2)
	qt.Join("a")
	if qt.Reached() {
		t.Fatalf("quorum should not be reached with one validator")
	}
	qt.Join("b")
	if !qt.Reached() {
		t.Fatalf("quorum should be reached with two validators")
	}
	qt.Leave("b")
	if qt.Reached() {
		t.Fatalf("quorum should no longer be reached after leave")
	}

	qt.Reset()
	if qt.Count() != 0 {
		t.Fatalf("expected empty quorum after reset")
	}

	qt.SetRequired(1)
	qt.Join("c")
	if qt.Required() != 1 || !qt.Reached() {
		t.Fatalf("failed to update quorum requirement")
	}
}
