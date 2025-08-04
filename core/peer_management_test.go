package core

import "testing"

func TestPeerManager(t *testing.T) {
	pm := NewPeerManager()
	pm.AddPeer("p1", "addr1")
	pm.AddPeer("p2", "addr2")
	if addr, ok := pm.GetPeer("p1"); !ok || addr != "addr1" {
		t.Fatalf("peer lookup failed")
	}
	if len(pm.ListPeers()) != 2 {
		t.Fatalf("expected 2 peers")
	}
	pm.RemovePeer("p1")
	if _, ok := pm.GetPeer("p1"); ok {
		t.Fatalf("remove failed")
	}
}
