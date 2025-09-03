package nodes

import (
	"testing"
	"time"
)

func TestHistoricalNodeArchiveRetrieve(t *testing.T) {
	n := NewHistoricalNode(Address("h1"))
	summary := BlockSummary{Height: 1, Hash: "hash1", Timestamp: time.Now()}
	if err := n.ArchiveBlock(summary); err != nil {
		t.Fatalf("archive: %v", err)
	}
	if got, ok := n.GetBlockByHeight(1); !ok || got.Hash != "hash1" {
		t.Fatalf("get by height failed: %#v %v", got, ok)
	}
	if got, ok := n.GetBlockByHash("hash1"); !ok || got.Height != 1 {
		t.Fatalf("get by hash failed: %#v %v", got, ok)
	}
	if n.TotalBlocks() != 1 {
		t.Fatalf("expected total blocks 1 got %d", n.TotalBlocks())
	}
}
