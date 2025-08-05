package core

import (
	"testing"

	"synnergy/internal/nodes"
)

func TestFullNodeModes(t *testing.T) {
	fn := NewFullNode(nodes.Address("f1"), FullNodeModeArchive)
	if !fn.IsArchive() {
		t.Fatalf("expected archive mode")
	}
	if fn.CurrentMode() != FullNodeModeArchive {
		t.Fatalf("unexpected mode: %v", fn.CurrentMode())
	}
	fn.SetMode(FullNodeModePruned)
	if fn.IsArchive() {
		t.Fatalf("expected pruned mode after SetMode")
	}
	if fn.CurrentMode() != FullNodeModePruned {
		t.Fatalf("mode not updated")
	}
}
