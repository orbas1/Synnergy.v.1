package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
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
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	svc.Enroll("alice", bio, &key.PublicKey)
	h := sha256.Sum256(bio)
	sig, err := ecdsa.SignASN1(rand.Reader, key, h[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
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
