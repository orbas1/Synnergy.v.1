package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"synnergy/internal/nodes"
)

func TestBaseNodeLifecycle(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	bn := NewBaseNode(nodes.Address(hex.EncodeToString(pub)))
	if bn.IsRunning() {
		t.Fatalf("expected node to be stopped initially")
	}
	if err := bn.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !bn.IsRunning() {
		t.Fatalf("expected node to be running after start")
	}
	peerPub, peerPriv, _ := ed25519.GenerateKey(rand.Reader)
	peerAddr := nodes.Address(hex.EncodeToString(peerPub))
	sig := ed25519.Sign(peerPriv, []byte(peerAddr))
	if err := bn.DialSeedSigned(peerAddr, sig, peerPub); err != nil {
		t.Fatalf("dial: %v", err)
	}
	if len(bn.Peers()) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(bn.Peers()))
	}
	snaps := bn.PeerSnapshots()
	if len(snaps) != 1 || snaps[0].Address != peerAddr {
		t.Fatalf("unexpected snapshot %+v", snaps)
	}
	if err := bn.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if bn.IsRunning() {
		t.Fatalf("expected node to be stopped after Stop")
	}
}

func TestBaseNodePeerRetention(t *testing.T) {
	bn := NewBaseNode(nodes.Address("base"))
	if err := bn.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer bn.Stop()

	bn.SetPeerRetention(1, 10*time.Millisecond)
	p1 := nodes.Address("p1")
	if err := bn.DialSeed(p1); err != nil {
		t.Fatalf("dial p1: %v", err)
	}
	if !bn.PromotePeer(p1) {
		t.Fatalf("expected p1 promotion")
	}
	if err := bn.DialSeed(nodes.Address("p2")); err == nil {
		t.Fatalf("expected capacity error when adding p2")
	}

	if !bn.DemotePeer(p1) {
		t.Fatalf("expected demotion to succeed")
	}
	bn.SetFailureThreshold(2)
	if !bn.RecordPeerFailure(p1) {
		t.Fatalf("expected failure to be tracked")
	}
	if bn.RecordPeerFailure(p1) {
		t.Fatalf("peer should have been evicted after second failure")
	}

	if err := bn.DialSeed(p1); err != nil {
		t.Fatalf("redial p1: %v", err)
	}
	bn.SetPeerRetention(2, 5*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	if removed := bn.PruneExpiredPeers(); removed == 0 {
		t.Fatalf("expected expired peer to be pruned")
	}
	if len(bn.Peers()) != 0 {
		t.Fatalf("expected peer list to be empty")
	}
}
