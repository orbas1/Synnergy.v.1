package core

import (
	militarynodes "synnergy/Nodes/military_nodes"
	"testing"
)

func TestWarfareNode(t *testing.T) {
	base := NewNode("n1", "addr", NewLedger())
	wn := NewWarfareNode(base)

	if err := wn.SecureCommand(""); err == nil {
		t.Fatal("expected error for empty command")
	}
	if err := wn.SecureCommand("move"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wn.TrackLogistics("asset1", "locA", "idle")
	wn.TrackLogistics("asset1", "locB", "moving")
	logs := wn.Logistics()
	if len(logs) != 2 {
		t.Fatalf("expected 2 records, got %d", len(logs))
	}
	// ensure interface satisfaction at compile time
	var _ militarynodes.WarfareNode = wn
}
