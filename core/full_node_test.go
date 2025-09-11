package core

import (
	"sync"
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

// TestFullNode_ModeConcurrency ensures SetMode is safe under concurrent access.
func TestFullNode_ModeConcurrency(t *testing.T) {
	fn := NewFullNode(nodes.Address("f1"), FullNodeModeArchive)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				fn.SetMode(FullNodeModeArchive)
			} else {
				fn.SetMode(FullNodeModePruned)
			}
		}(i)
	}
	wg.Wait()
	m := fn.CurrentMode()
	if m != FullNodeModeArchive && m != FullNodeModePruned {
		t.Fatalf("unexpected mode after concurrent updates: %v", m)
	}
}
