package core

import (
	"sync"
	"testing"
	"time"

	nodes "synnergy/internal/nodes"
)

func TestForensicNode_RecordAndRetrieve(t *testing.T) {
	fn := NewForensicNode()
	tx := nodes.TransactionLite{Hash: "tx1", From: "a", To: "b", Value: 10, Timestamp: time.Now()}
	trace := nodes.NetworkTrace{PeerID: "peer1", Event: "connect", Timestamp: time.Now()}
	if err := fn.RecordTransaction(tx); err != nil {
		t.Fatalf("record tx: %v", err)
	}
	if err := fn.RecordNetworkTrace(trace); err != nil {
		t.Fatalf("record trace: %v", err)
	}
	if len(fn.Transactions()) != 1 {
		t.Fatalf("expected 1 tx, got %d", len(fn.Transactions()))
	}
	if len(fn.NetworkTraces()) != 1 {
		t.Fatalf("expected 1 trace, got %d", len(fn.NetworkTraces()))
	}
}

// TestForensicNode_PrunesOldEntries ensures the buffers do not grow without
// bound and older entries are dropped when the limit is exceeded.
func TestForensicNode_PrunesOldEntries(t *testing.T) {
	fn := NewForensicNodeWithLimit(3)
	for i := 0; i < 5; i++ {
		tx := nodes.TransactionLite{Hash: string(rune('a' + i))}
		trace := nodes.NetworkTrace{PeerID: string(rune('a' + i))}
		_ = fn.RecordTransaction(tx)
		_ = fn.RecordNetworkTrace(trace)
	}
	if got := len(fn.Transactions()); got != 3 {
		t.Fatalf("expected 3 transactions after pruning, got %d", got)
	}
	if got := len(fn.NetworkTraces()); got != 3 {
		t.Fatalf("expected 3 traces after pruning, got %d", got)
	}
	if fn.Transactions()[0].Hash != "c" {
		t.Fatalf("expected oldest tx to be 'c', got %s", fn.Transactions()[0].Hash)
	}
}

// TestForensicNode_Concurrency checks that concurrent writes do not race and
// the limit is respected.
func TestForensicNode_Concurrency(t *testing.T) {
	fn := NewForensicNodeWithLimit(100)
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tx := nodes.TransactionLite{Hash: string(rune(i))}
			trace := nodes.NetworkTrace{PeerID: string(rune(i))}
			_ = fn.RecordTransaction(tx)
			_ = fn.RecordNetworkTrace(trace)
		}(i)
	}
	wg.Wait()
	if got := len(fn.Transactions()); got != 100 {
		t.Fatalf("expected 100 transactions after concurrent writes, got %d", got)
	}
	if got := len(fn.NetworkTraces()); got != 100 {
		t.Fatalf("expected 100 traces after concurrent writes, got %d", got)
	}
}
