package p2p

import "testing"

func TestManager(t *testing.T) {
	m := NewManager()
	p := Peer{ID: "peer1", Address: "127.0.0.1:8000"}
	m.AddPeer(p)
	if len(m.ListPeers()) != 1 {
		t.Fatalf("expected 1 peer")
	}
	if _, ok := m.GetPeer("peer1"); !ok {
		t.Fatalf("peer not found")
	}
	m.RemovePeer("peer1")
	if len(m.ListPeers()) != 0 {
		t.Fatalf("peer not removed")
	}
}

func TestDiscoveryService(t *testing.T) {
	m := NewManager()
	bootstrap := []Peer{{ID: "b1", Address: "127.0.0.1:8001"}}
	d := NewDiscoveryService(m, bootstrap)
	peers := d.DiscoverPeers()
	if len(peers) != 1 || peers[0].ID != "b1" {
		t.Fatalf("unexpected discovery result: %+v", peers)
	}
	// Subsequent discovery should return manager peers
	peers = d.DiscoverPeers()
	if len(peers) != 1 || peers[0].ID != "b1" {
		t.Fatalf("unexpected second discovery: %+v", peers)
	}
}
