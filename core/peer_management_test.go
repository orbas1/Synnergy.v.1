package core

import (
	"fmt"
	"sync"
	"testing"
)

// TestPeerManagerBasicOperations verifies add, get, list and remove functionality.
func TestPeerManagerBasicOperations(t *testing.T) {
	pm := NewPeerManager()
	pm.AddPeer("p1", "addr1")
	pm.AddPeer("p2", "addr2")

	if addr, ok := pm.GetPeer("p1"); !ok || addr != "addr1" {
		t.Fatalf("peer lookup failed for p1: %v %v", addr, ok)
	}
	if addr, ok := pm.GetPeer("p2"); !ok || addr != "addr2" {
		t.Fatalf("peer lookup failed for p2: %v %v", addr, ok)
	}

	peers := pm.ListPeers()
	if len(peers) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(peers))
	}
	if pm.Count() != 2 {
		t.Fatalf("expected count 2, got %d", pm.Count())
	}

	pm.RemovePeer("p1")
	if _, ok := pm.GetPeer("p1"); ok {
		t.Fatalf("peer p1 should have been removed")
	}
	if len(pm.ListPeers()) != 1 {
		t.Fatalf("expected 1 peer after removal, got %d", len(pm.ListPeers()))
	}
}

// TestPeerManagerConnect ensures Connect records the peer by its address.
func TestPeerManagerConnect(t *testing.T) {
	pm := NewPeerManager()
	id := pm.Connect("addr1")
	if id != "addr1" {
		t.Fatalf("expected id addr1, got %s", id)
	}
	if addr, ok := pm.GetPeer(id); !ok || addr != "addr1" {
		t.Fatalf("connect should add peer with same address")
	}
}

// TestPeerManagerAdvertiseAndDiscover checks topic based discovery.
func TestPeerManagerAdvertiseAndDiscover(t *testing.T) {
	pm := NewPeerManager()
	pm.AddPeer("p1", "addr1")
	pm.AddPeer("p2", "addr2")
	pm.Advertise("p1", "sync")
	pm.Advertise("p2", "sync")
	pm.Advertise("p2", "other")

	ids := pm.Discover("sync")
	if len(ids) != 2 || !contains(ids, "p1") || !contains(ids, "p2") {
		t.Fatalf("discover on sync returned %v", ids)
	}
	ids = pm.Discover("other")
	if len(ids) != 1 || ids[0] != "p2" {
		t.Fatalf("discover on other returned %v", ids)
	}
}

// TestPeerManagerConcurrentAccess performs basic operations concurrently to
// ensure thread safety of the manager.
func TestPeerManagerConcurrentAccess(t *testing.T) {
	pm := NewPeerManager()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("p%d", i)
			addr := fmt.Sprintf("addr%d", i)
			pm.AddPeer(id, addr)
			pm.Advertise(id, "topic")
			if a, ok := pm.GetPeer(id); !ok || a != addr {
				t.Errorf("peer %s lookup failed", id)
			}
		}(i)
	}
	wg.Wait()

	if len(pm.ListPeers()) != 20 {
		t.Fatalf("expected 20 peers, got %d", len(pm.ListPeers()))
	}
	if len(pm.Discover("topic")) != 20 {
		t.Fatalf("expected 20 adverts, got %d", len(pm.Discover("topic")))
	}
}

// helper contains checks if slice contains element
func contains(s []string, v string) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}
