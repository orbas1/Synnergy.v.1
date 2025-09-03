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
}

func TestForensicNodeConcurrent(t *testing.T) {
	n := NewForensicNode(Address("f2"))
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
