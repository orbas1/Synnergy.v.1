package core

import (
	"testing"

	"synnergy/internal/nodes"
)

// TestNewNodeAdapter verifies that the adapter correctly wraps a Node and
// exposes BaseNode behaviour.
func TestNewNodeAdapter(t *testing.T) {
	ledger := NewLedger()
	n := NewNode("node1", "addr1", ledger)

	adapter := NewNodeAdapter(n)
	if adapter.node != n {
		t.Fatalf("expected adapter to wrap provided node")
	}
	if adapter.BaseNode == nil {
		t.Fatalf("expected BaseNode to be initialized")
	}
	if adapter.ID() != nodes.Address(n.ID) {
		t.Fatalf("adapter ID %s does not match node ID %s", adapter.ID(), n.ID)
	}

	if adapter.IsRunning() {
		t.Fatalf("adapter should not be running initially")
	}
	if err := adapter.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !adapter.IsRunning() {
		t.Fatalf("adapter should be running after Start")
	}

	seed := nodes.Address("peer1")
	if err := adapter.DialSeed(seed); err != nil {
		t.Fatalf("dial seed: %v", err)
	}
	peers := adapter.Peers()
	if len(peers) != 1 || peers[0] != seed {
		t.Fatalf("expected peers to contain %v, got %v", seed, peers)
	}

	if err := adapter.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if adapter.IsRunning() {
		t.Fatalf("adapter should be stopped after Stop")
	}
	if err := adapter.DialSeed(nodes.Address("peer2")); err == nil {
		t.Fatalf("DialSeed should fail when adapter not running")
	}
}

func TestNewNodeAdapterNilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for nil node")
		}
	}()
	_ = NewNodeAdapter(nil)
}
