package p2p

import (
	"testing"
	"time"

	"synnergy/internal/security"
)

func TestManagerLifecycle(t *testing.T) {
	manager := NewManager(nil)
	peer := Peer{ID: "peer1", Address: "127.0.0.1:8000", Capabilities: map[string]bool{"archive": true}}
	added := manager.AddPeer(peer)
	if added.ID == "" {
		t.Fatalf("expected id")
	}
	peers := manager.ListPeers()
	if len(peers) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(peers))
	}
	snapshot := manager.Snapshot()
	if len(snapshot) != 1 || snapshot[0].ID != added.ID {
		t.Fatalf("unexpected snapshot")
	}
	updated := Peer{ID: added.ID, Address: "127.0.0.1:8001", Capabilities: map[string]bool{"archive": true}}
	manager.UpdatePeer(updated, "address change")
	got, ok := manager.GetPeer(added.ID)
	if !ok || got.Address != "127.0.0.1:8001" {
		t.Fatalf("expected updated address")
	}
	manager.RemovePeer(added.ID, "decommission")
	if len(manager.ListPeers()) != 0 {
		t.Fatalf("expected removal")
	}
}

func TestManagerEventsAndFailures(t *testing.T) {
	mitigator := security.NewDDoSMitigator(security.MitigationConfig{Window: time.Second, MaxRequests: 5})
	manager := NewManager(mitigator)
	ch, cancel := manager.Subscribe(0)
	defer cancel()
	peer := manager.AddPeer(Peer{ID: "p", Address: "1.1.1.1"})
	evt := <-ch
	if evt.Type != PeerEventAdded {
		t.Fatalf("expected add event")
	}
	manager.MarkFailure(peer.ID, "timeout")
	evt = <-ch
	if evt.Type != PeerEventQuarantined {
		t.Fatalf("expected quarantine event")
	}
}
