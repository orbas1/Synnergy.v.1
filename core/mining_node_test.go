package core

import "testing"

func TestMiningNode(t *testing.T) {
	mn := NewMiningNode(100)
	if mn.IsMining() {
		t.Fatalf("expected node to be stopped initially")
	}
	mn.Start()
	if !mn.IsMining() {
		t.Fatalf("node should be mining after Start")
	}
	hash, err := mn.Mine([]byte("data"))
	if err != nil || hash == "" {
		t.Fatalf("mine returned invalid hash: %v %s", err, hash)
	}
	mn.Stop()
	if mn.IsMining() {
		t.Fatalf("node should be stopped after Stop")
	}
	if _, err := mn.Mine([]byte("data")); err == nil {
		t.Fatalf("expected error when mining is inactive")
	}
}
