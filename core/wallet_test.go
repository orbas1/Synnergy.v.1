package core

import (
	"os"
	"testing"
)

func TestWalletSignAndVerify(t *testing.T) {
	w1, err := NewWallet()
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}
	w2, err := NewWallet()
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}
	if w1.Address == w2.Address {
		t.Fatalf("wallet addresses should be unique")
	}
	if len(w1.Address) != 40 || len(w2.Address) != 40 {
		t.Fatalf("addresses must be 40 hex characters")
	}
	tx := NewTransaction("a", "b", 1, 0, 0)
	sig, err := w1.Sign(tx)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	if !VerifySignature(tx, sig, &w1.PrivateKey.PublicKey) {
		t.Fatalf("signature verification failed")
	}
	if VerifySignature(tx, sig, &w2.PrivateKey.PublicKey) {
		t.Fatalf("verification should fail with wrong public key")
	}
	if !tx.Verify(&w1.PrivateKey.PublicKey) {
		t.Fatalf("tx.Verify should succeed with correct key")
	}
}

func TestWalletSaveLoad(t *testing.T) {
	w, err := NewWallet()
	if err != nil {
		t.Fatalf("new wallet: %v", err)
	}
	path := "testwallet.json"
	if err := w.Save(path, "pass"); err != nil {
		t.Fatalf("save: %v", err)
	}
	defer os.Remove(path)
	w2, err := LoadWallet(path, "pass")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if w.Address != w2.Address {
		t.Fatalf("addresses differ")
	}
}
