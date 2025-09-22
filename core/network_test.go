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

	// ensure telemetry recorded successful deliveries
	want := map[string]string{"n1": "node", "n2": "node", "r1": "relay"}
	deadline := time.After(200 * time.Millisecond)
	for len(want) > 0 {
		select {
		case event := <-network.BroadcastEvents():
			if event.TransactionID != tx.ID {
				continue
			}
			role, ok := want[event.Target]
			if !ok {
				continue
			}
			if role != event.Role {
				t.Fatalf("unexpected role for %s: %s", event.Target, event.Role)
			}
			if !event.Success {
				t.Fatalf("expected success for %s: %s", event.Target, event.Error)
			}
			delete(want, event.Target)
		case <-deadline:
			t.Fatalf("missing telemetry for: %v", want)
		}
	}
}

// TestNetworkPubSub ensures the lightweight publish/subscribe system delivers
// messages to all subscribed listeners.
func TestNetworkPubSub(t *testing.T) {
	net := NewNetwork(NewBiometricService())
	sub, cancel := net.Subscribe("demo")
	defer cancel()
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

// TestNetworkBroadcastTelemetry verifies that validation failures are surfaced
// through broadcast events instead of being silently dropped.
func TestNetworkBroadcastTelemetry(t *testing.T) {
	svc := NewBiometricService()
	network := NewNetwork(svc)

	good := NewNode("good", "addr", NewLedger())
	bad := NewNode("bad", "addr2", NewLedger())
	good.Ledger.Credit("alice", 200)
	// bad node intentionally lacks balance for alice

	network.AddNode(good)
	network.AddNode(bad)

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
	tx := NewTransaction("alice", "bob", 50, 10, 1)
	if err := network.Broadcast(tx, "alice", bio, sig); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}

	// allow goroutine to process queue
	time.Sleep(50 * time.Millisecond)

	if len(good.Mempool) != 1 {
		t.Fatalf("expected transaction to reach good node")
	}
	if len(bad.Mempool) != 0 {
		t.Fatalf("expected failing node to reject transaction")
	}

	var badErr string
	deadline := time.After(200 * time.Millisecond)
	for badErr == "" {
		select {
		case event := <-network.BroadcastEvents():
			if event.TransactionID != tx.ID {
				continue
			}
			if event.Target == "bad" {
				badErr = event.Error
				if event.Success {
					t.Fatalf("expected bad node to report failure")
				}
			}
		case <-deadline:
			t.Fatal("did not observe failure telemetry")
		}
	}
	if badErr == "" {
		t.Fatal("expected non-empty error message")
	}
}
