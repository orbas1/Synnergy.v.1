package core

import (
	"testing"
	"time"

	"synnergy/internal/nodes"
)

func TestReplicator(t *testing.T) {
	l := NewLedger()
	b := &Block{Hash: "b1"}
	if err := l.AddBlock(b); err != nil {
		t.Fatalf("add block: %v", err)
	}

	r := NewReplicator(l)
	r.RegisterPeer(nodes.Address("p1"))
	r.RegisterPeer(nodes.Address("p2"))
	r.SetRetryPolicy(2, 10*time.Millisecond)
	r.SetEnqueueTimeout(20 * time.Millisecond)
	r.Start()
	defer r.Stop()
	if !r.Status() {
		t.Fatalf("expected running")
	}

	if r.ReplicateBlock("missing") {
		t.Fatalf("expected missing block to be rejected")
	}
	if !r.ReplicateBlock("b1") {
		t.Fatalf("replication failed")
	}
	time.Sleep(20 * time.Millisecond)
	if r.Replicated("b1") {
		t.Fatalf("expected pending acknowledgements")
	}
	rec, ok := r.ReplicationStatus("b1")
	if !ok || rec.Attempts == 0 {
		t.Fatalf("expected replication attempts recorded")
	}
	if !r.MarkAcknowledged("b1", nodes.Address("p1")) {
		t.Fatalf("expected acknowledgement to be recorded")
	}
	time.Sleep(5 * time.Millisecond)
	if r.Replicated("b1") {
		t.Fatalf("should still await second acknowledgement")
	}
	if !r.MarkAcknowledged("b1", nodes.Address("p2")) {
		t.Fatalf("expected second acknowledgement")
	}
	time.Sleep(20 * time.Millisecond)
	if !r.Replicated("b1") {
		t.Fatalf("expected block marked replicated")
	}
	metrics := r.Metrics()
	if metrics.Enqueued == 0 || metrics.Acknowledged < 2 || metrics.Pending != 0 {
		t.Fatalf("unexpected metrics %+v", metrics)
	}
}
