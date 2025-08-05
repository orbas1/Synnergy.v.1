package nodes

import (
	"fmt"
	"sync"
	"testing"

	"synnergy"
)

func TestHolographicNodeStoreRetrieve(t *testing.T) {
	n := NewHolographicNode(Address("node1"))
	frame := synnergy.HolographicFrame{ID: "f1", Shards: [][]byte{{1, 2, 3}}}
	n.Store(frame)
	got, ok := n.Retrieve("f1")
	if !ok || got.ID != "f1" {
		t.Fatalf("expected to retrieve frame f1")
	}
}

func TestHolographicNodePeers(t *testing.T) {
	n := NewHolographicNode(Address("node1"))
	n.DialSeed(Address("peer1"))
	n.DialSeed(Address("peer1"))
	peers := n.Peers()
	if len(peers) != 1 || peers[0] != Address("peer1") {
		t.Fatalf("unexpected peers: %#v", peers)
	}
}

func TestHolographicNodeIDAndLifecycle(t *testing.T) {
	n := NewHolographicNode(Address("node1"))
	if n.ID() != Address("node1") {
		t.Fatalf("unexpected id: %s", n.ID())
	}
	if err := n.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if err := n.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}

func TestHolographicNodeRetrieveMissing(t *testing.T) {
	n := NewHolographicNode(Address("node1"))
	if _, ok := n.Retrieve("missing"); ok {
		t.Fatalf("expected missing frame to return false")
	}
}

func TestHolographicNodeConcurrentStoreRetrieve(t *testing.T) {
	n := NewHolographicNode(Address("node1"))
	var wg sync.WaitGroup
	count := 50
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			frame := synnergy.HolographicFrame{ID: fmt.Sprintf("f%02d", i), Shards: [][]byte{{byte(i)}}}
			n.Store(frame)
		}(i)
	}
	wg.Wait()
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("f%02d", i)
		if _, ok := n.Retrieve(id); !ok {
			t.Fatalf("expected to find frame %s", id)
		}
	}
}
