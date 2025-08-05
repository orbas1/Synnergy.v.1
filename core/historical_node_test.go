package core

import (
	nodes "synnergy/internal/nodes"
	"testing"
	"time"
)

func TestHistoricalNode_ArchiveAndRetrieve(t *testing.T) {
	hn := NewHistoricalNode()
	summary := nodes.BlockSummary{Height: 1, Hash: "abc", Timestamp: time.Now()}
	if err := hn.ArchiveBlock(summary); err != nil {
		t.Fatalf("archive failed: %v", err)
	}
	if hn.TotalBlocks() != 1 {
		t.Fatalf("expected 1 block, got %d", hn.TotalBlocks())
	}
	if s, ok := hn.GetBlockByHeight(1); !ok || s.Hash != "abc" {
		t.Fatalf("failed to get by height")
	}
	if s, ok := hn.GetBlockByHash("abc"); !ok || s.Height != 1 {
		t.Fatalf("failed to get by hash")
	}
}

func TestHistoricalNode_Duplicate(t *testing.T) {
	hn := NewHistoricalNode()
	summary := nodes.BlockSummary{Height: 1, Hash: "abc", Timestamp: time.Now()}
	if err := hn.ArchiveBlock(summary); err != nil {
		t.Fatalf("archive failed: %v", err)
	}
	if err := hn.ArchiveBlock(summary); err == nil {
		t.Fatalf("expected duplicate error")
	}
}
