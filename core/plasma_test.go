package core

import "testing"

// TestPausePreventsOperations ensures paused bridge rejects new deposits and exits.
func TestPausePreventsOperations(t *testing.T) {
	b := NewPlasmaBridge()
	b.Pause()
	if err := b.Deposit("a", "t", 1); err != ErrBridgePaused {
		t.Fatalf("expected ErrBridgePaused, got %v", err)
	}
	if _, err := b.StartExit("a", "t", 1); err != ErrBridgePaused {
		t.Fatalf("expected ErrBridgePaused from StartExit")
	}
}
