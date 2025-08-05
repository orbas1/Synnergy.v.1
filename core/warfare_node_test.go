package core

import (
	militarynodes "synnergy/internal/nodes/military_nodes"
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
	wn.TrackLogistics("asset2", "locC", "standby")
	logs := wn.Logistics()
	if len(logs) != 3 {
		t.Fatalf("expected 3 records, got %d", len(logs))
	}
	// Verify filtering by asset
	asset1 := wn.LogisticsByAsset("asset1")
	if len(asset1) != 2 {
		t.Fatalf("expected 2 records for asset1, got %d", len(asset1))
	}
	asset2 := wn.LogisticsByAsset("asset2")
	if len(asset2) != 1 {
		t.Fatalf("expected 1 record for asset2, got %d", len(asset2))
	}
	// ensure interface satisfaction at compile time
	var _ militarynodes.WarfareNode = wn
}
