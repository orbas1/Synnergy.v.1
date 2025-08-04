package core

import "testing"

func TestReplicator(t *testing.T) {
	l := NewLedger()
	b := &Block{Hash: "b1"}
	l.AddBlock(b)

	r := NewReplicator(l)
	r.Start()
	if !r.Status() {
		t.Fatalf("expected running")
	}

	if !r.ReplicateBlock("b1") {
		t.Fatalf("replication failed")
	}

	r.Stop()
	if r.Status() {
		t.Fatalf("expected stopped")
	}
}
