package core

import (
	"testing"

	"synnergy/internal/nodes"
)

func TestBaseNodeLifecycle(t *testing.T) {
	bn := NewBaseNode(nodes.Address("node1"))
	if bn.IsRunning() {
		t.Fatalf("expected node to be stopped initially")
	}
	if err := bn.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !bn.IsRunning() {
		t.Fatalf("expected node to be running after start")
	}
	if err := bn.DialSeed(nodes.Address("peer1")); err != nil {
		t.Fatalf("dial: %v", err)
	}
	if len(bn.Peers()) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(bn.Peers()))
	}
	if err := bn.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if bn.IsRunning() {
		t.Fatalf("expected node to be stopped after Stop")
	}
}
