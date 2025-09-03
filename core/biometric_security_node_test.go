package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestBiometricSecurityNode(t *testing.T) {
	ledger := NewLedger()
	ledger.Credit("from", 100)
	base := NewNode("node1", "addr1", ledger)
	bsn := NewBiometricSecurityNode(base, nil)

	admin := "admin"
	bio := []byte("admin-bio")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	bsn.Enroll(admin, bio, &key.PublicKey)
	hash := sha256.Sum256(bio)
	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	tx := NewTransaction("from", "to", 1, 0, 0)
	if err := bsn.SecureAddTransaction(admin, bio, sig, tx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(base.Mempool) != 1 {
		t.Fatal("transaction not added to mempool")
	}

	tx2 := NewTransaction("from", "to", 1, 0, 1)
	wrongHash := sha256.Sum256([]byte("wrong"))
	wrongSig, _ := ecdsa.SignASN1(rand.Reader, key, wrongHash[:])
	if err := bsn.SecureAddTransaction(admin, []byte("wrong"), wrongSig, tx2); err == nil {
		t.Fatal("expected authentication failure")
	}
	if len(base.Mempool) != 1 {
		t.Fatal("unexpected transaction added")
	}

	// Test SecureExecute with correct and incorrect biometrics
	if err := bsn.SecureExecute(admin, bio, sig, func() error {
		bsn.Node.ID = "updated"
		return nil
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bsn.Node.ID != "updated" {
		t.Fatal("secure execute did not run")
	}
	badHash := sha256.Sum256([]byte("bad"))
	badSig, _ := ecdsa.SignASN1(rand.Reader, key, badHash[:])
	if err := bsn.SecureExecute(admin, []byte("bad"), badSig, nil); err == nil {
		t.Fatal("expected verification failure")
	}

	// Test removal of biometrics
	bsn.Remove(admin)
	if bsn.Auth.Enrolled(admin) {
		t.Fatal("expected biometric data to be removed")
	}
}
