package core

import (
	"testing"

	"synnergy/nodes"
)

func TestBaseNodeLifecycle(t *testing.T) {
	bn := NewBaseNode(nodes.Address("node1"))
	if err := bn.Start(); err != nil {
		t.Fatalf("start: %v", err)
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
}
