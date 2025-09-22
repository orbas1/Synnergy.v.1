package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math"
	"testing"
	"time"
)

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
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	if err := svc.Enroll(user, bio, pub); err != nil {
		t.Fatalf("enroll: %v", err)
	}
	hash := sha256.Sum256(bio)
	sig := ed25519.Sign(priv, hash[:])

	tx := NewTransaction("a", "b", 1, 0, 0)
	origID := tx.ID
	if err := tx.AttachBiometric(user, bio, sig, svc); err != nil {
		t.Fatalf("attach biometric failed: %v", err)
	}
	if tx.ID == origID {
		t.Fatalf("transaction ID should change after attaching biometric")
	}
	if len(tx.BiometricHash) == 0 {
		t.Fatalf("biometric hash not set")
	}
	wh := sha256.Sum256([]byte("wrong"))
	wrongSig := ed25519.Sign(priv, wh[:])
	if err := tx.AttachBiometric(user, []byte("wrong"), wrongSig, svc); err == nil {
		t.Fatalf("expected verification failure")
	}
	if err := tx.AttachBiometric(user, bio, sig, nil); err == nil {
		t.Fatalf("expected error when service nil")
	}
}

func TestTransactionProgramAffectsHash(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 1, 0)
	original := tx.ID
	tx.Program = []Instruction{{Op: OpPush, Value: 5}}
	tx.ID = tx.Hash()
	if tx.ID == original {
		t.Fatal("expected hash to change when program is attached")
	}
	prev := tx.ID
	tx.Program[0].Value = 6
	tx.ID = tx.Hash()
	if tx.ID == prev {
		t.Fatal("expected program mutation to alter hash")
	}
}

func TestTransactionClone(t *testing.T) {
	tx := NewTransaction("from", "to", 5, 1, 2)
	tx.Signature = []byte{1, 2, 3}
	tx.BiometricHash = []byte{4, 5, 6}
	tx.Program = []Instruction{{Op: OpPush, Value: 3}}
	clone := tx.Clone()
	if clone == tx {
		t.Fatal("expected clone to allocate new struct")
	}
	if clone.ID != clone.Hash() {
		t.Fatalf("clone hash mismatch: %s", clone.ID)
	}
	clone.Signature[0] = 9
	if tx.Signature[0] == 9 {
		t.Fatal("signature not deep copied")
	}
	clone.Program[0].Value = 99
	if tx.Program[0].Value == 99 {
		t.Fatal("program not deep copied")
	}
}

func TestTransactionValidateBasic(t *testing.T) {
	tx := NewTransaction("from", "to", 5, 1, 2)
	cfg := DefaultTransactionValidationConfig()
	cfg.Now = time.Unix(tx.Timestamp, 0)
	if err := tx.ValidateBasic(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bad := tx.Clone()
	bad.Amount = 0
	if err := bad.ValidateBasic(cfg); !errors.Is(err, ErrZeroAmount) {
		t.Fatalf("expected zero amount error got %v", err)
	}

	bad = tx.Clone()
	bad.From = ""
	if err := bad.ValidateBasic(cfg); !errors.Is(err, ErrEmptyAddress) {
		t.Fatalf("expected empty address error got %v", err)
	}

	bad = tx.Clone()
	bad.Timestamp = time.Unix(tx.Timestamp, 0).Add(2 * DefaultMaxClockSkew).Unix()
	if err := bad.ValidateBasic(cfg); !errors.Is(err, ErrTransactionClockSkew) {
		t.Fatalf("expected clock skew error got %v", err)
	}

	bad = tx.Clone()
	bad.Program = []Instruction{{Op: OpAdd, Value: 5}}
	bad.ID = bad.Hash()
	if err := bad.ValidateBasic(cfg); err == nil {
		t.Fatal("expected invalid program error")
	}
}

func TestTransactionTotalCostOverflow(t *testing.T) {
	tx := &Transaction{Amount: math.MaxUint64, Fee: 1, ID: "", From: "a", To: "b", Timestamp: time.Now().Unix()}
	if _, err := tx.TotalCost(); !errors.Is(err, ErrTransactionOverflow) {
		t.Fatalf("expected overflow error got %v", err)
	}
}
