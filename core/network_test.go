package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"testing"
	"time"
)

// TestNetworkBroadcast verifies that transactions are propagated to all nodes
// and relay nodes after biometric authentication.
func TestNetworkBroadcast(t *testing.T) {
	svc := NewBiometricService()
	network := NewNetwork(svc)

	// setup nodes and relay
	n1 := NewNode("n1", "addr1", NewLedger())
	n2 := NewNode("n2", "addr2", NewLedger())
	relay := NewNode("r1", "addrR", NewLedger())

	// credit balances so validation passes
	n1.Ledger.Credit("alice", 100)
	n2.Ledger.Credit("alice", 100)
	relay.Ledger.Credit("alice", 100)

	network.AddNode(n1)
	network.AddNode(n2)
	network.AddRelay(relay)

	bio := []byte("finger")
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	if err := svc.Enroll("alice", bio, pub); err != nil {
		t.Fatalf("enroll: %v", err)
	}
	h := sha256.Sum256(bio)
	sig := ed25519.Sign(priv, h[:])
	tx := NewTransaction("alice", "bob", 1, 1, 1)
	if err := network.Broadcast(tx, "alice", bio, sig); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}

	// allow goroutine to process queue
	time.Sleep(50 * time.Millisecond)

	if len(n1.Mempool) != 1 || len(n2.Mempool) != 1 || len(relay.Mempool) != 1 {
		t.Fatalf("expected transaction to be broadcast to all nodes and relay")
	}
}

// TestNetworkPubSub ensures the lightweight publish/subscribe system delivers
// messages to all subscribed listeners.
func TestNetworkPubSub(t *testing.T) {
	net := NewNetwork(NewBiometricService())
	sub := net.Subscribe("demo")
	msg := []byte("ping")
	net.Publish("demo", msg)
	select {
	case got := <-sub:
		if string(got) != string(msg) {
			t.Fatalf("expected %s, got %s", msg, got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("did not receive published message")
	}
}
