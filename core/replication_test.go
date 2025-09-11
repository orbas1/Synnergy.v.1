package core

import "testing"

func TestReplicator(t *testing.T) {
	l := NewLedger()
	b := &Block{Hash: "b1"}
	if err := l.AddBlock(b); err != nil {
		t.Fatalf("add block: %v", err)
	}

	r := NewReplicator(l)
	r.Start()
	if !r.Status() {
		t.Fatalf("expected running")
	}

	if !r.ReplicateBlock("b1") {
		t.Fatalf("replication failed")
	}
	if !r.Replicated("b1") {
		t.Fatalf("expected block marked replicated")
	}

	r.Stop()
	if r.Status() {
		t.Fatalf("expected stopped")
	}
}
