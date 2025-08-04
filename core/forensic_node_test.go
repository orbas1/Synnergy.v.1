package core

import (
	nodes "synnergy/nodes"
	"testing"
	"time"
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
