package nodes

import (
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
	peers := n.Peers()
	if len(peers) != 1 || peers[0] != Address("peer1") {
		t.Fatalf("unexpected peers: %#v", peers)
	}
}
