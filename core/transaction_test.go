package core

import "testing"

func TestNewTransactionAndHash(t *testing.T) {
    tx := NewTransaction("alice", "bob", 10, 1, 0)
    if tx.ID == "" {
        t.Fatalf("expected ID to be set")
    }
    // ensure hash deterministic
    fixed := &Transaction{From: "alice", To: "bob", Amount: 10, Fee: 1, Nonce: 0, Timestamp: 42}
    h1 := fixed.Hash()
    fixed2 := &Transaction{From: "alice", To: "bob", Amount: 10, Fee: 1, Nonce: 0, Timestamp: 42}
    h2 := fixed2.Hash()
    if h1 != h2 {
        t.Fatalf("expected deterministic hash")
    }
}

func TestAttachBiometric(t *testing.T) {
    svc := NewBiometricService()
    user := "u1"
    bio := []byte("fingerprint")
    svc.Enroll(user, bio)

    tx := NewTransaction("a", "b", 1, 0, 0)
    origID := tx.ID
    if err := tx.AttachBiometric(user, bio, svc); err != nil {
        t.Fatalf("attach biometric failed: %v", err)
    }
    if tx.ID == origID {
        t.Fatalf("transaction ID should change after attaching biometric")
    }
    if len(tx.BiometricHash) == 0 {
        t.Fatalf("biometric hash not set")
    }
    if err := tx.AttachBiometric(user, []byte("wrong"), svc); err == nil {
        t.Fatalf("expected verification failure")
    }
    if err := tx.AttachBiometric(user, bio, nil); err == nil {
        t.Fatalf("expected error when service nil")
    }
}

