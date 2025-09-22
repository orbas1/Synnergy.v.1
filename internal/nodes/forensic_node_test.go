package nodes

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestForensicNodeRecords(t *testing.T) {
	n := NewForensicNode(Address("f1"))
	tx := TransactionLite{Hash: "h1", From: "a", To: "b", Value: 10, Timestamp: time.Now()}
	if err := n.RecordTransaction(tx); err != nil {
		t.Fatalf("record transaction: %v", err)
	}
	trace := NetworkTrace{PeerID: "p1", Event: "connect", Timestamp: time.Now()}
	if err := n.RecordNetworkTrace(trace); err != nil {
		t.Fatalf("record trace: %v", err)
	}
	if got := n.Transactions(); len(got) != 1 || got[0].Hash != "h1" {
		t.Fatalf("unexpected transactions: %#v", got)
	}
	if traces := n.NetworkTraces(); len(traces) != 1 || traces[0].PeerID != "p1" {
		t.Fatalf("unexpected traces: %#v", traces)
	}

	stats := n.Stats()
	if stats.TransactionCount != 1 || stats.TraceCount != 1 {
		t.Fatalf("unexpected stats %+v", stats)
	}
}

func TestForensicNodeConcurrent(t *testing.T) {
	n := NewForensicNode(Address("f2"), WithMaxTransactionRecords(1000))
	var wg sync.WaitGroup
	count := 50
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = n.RecordTransaction(TransactionLite{Hash: fmt.Sprintf("h%02d", i), Timestamp: time.Now()})
		}(i)
	}
	wg.Wait()
	if got := len(n.Transactions()); got != count {
		t.Fatalf("expected %d transactions got %d", count, got)
	}
}

func TestForensicNodeRetentionAndValidation(t *testing.T) {
	n := NewForensicNode(Address("f3"), WithMaxTransactionRecords(3), WithMaxTraceRecords(2))

	if err := n.RecordTransaction(TransactionLite{}); err == nil {
		t.Fatalf("expected error for empty hash")
	}
	if err := n.RecordNetworkTrace(NetworkTrace{}); err == nil {
		t.Fatalf("expected error for empty peer id")
	}

	for i := 0; i < 5; i++ {
		if err := n.RecordTransaction(TransactionLite{Hash: fmt.Sprintf("tx-%d", i)}); err != nil {
			t.Fatalf("record transaction %d: %v", i, err)
		}
	}
	if got := n.Transactions(); len(got) != 3 || got[0].Hash != "tx-2" {
		t.Fatalf("unexpected retained transactions: %#v", got)
	}

	for i := 0; i < 4; i++ {
		if err := n.RecordNetworkTrace(NetworkTrace{PeerID: fmt.Sprintf("peer-%d", i)}); err != nil {
			t.Fatalf("record trace %d: %v", i, err)
		}
	}
	if got := n.NetworkTraces(); len(got) != 2 || got[0].PeerID != "peer-2" {
		t.Fatalf("unexpected retained traces: %#v", got)
	}
}
