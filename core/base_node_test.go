package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"

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
	if err := bn.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if bn.IsRunning() {
		t.Fatalf("expected node to be stopped after Stop")
	}
}
