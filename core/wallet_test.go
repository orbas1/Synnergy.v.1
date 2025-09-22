package core

import (
	"bytes"
	"path/filepath"
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
	if len(w1.PublicKeyBytes()) == 0 || len(w2.PublicKeyBytes()) == 0 {
		t.Fatalf("expected public key bytes")
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
	msg := []byte("stage75")
	msgSig, err := w1.SignMessage(msg)
	if err != nil {
		t.Fatalf("sign message: %v", err)
	}
	if !VerifyMessage(msg, msgSig, &w1.PublicKey) {
		t.Fatalf("message verification failed")
	}
	if VerifyMessage(msg, msgSig, &w2.PublicKey) {
		t.Fatalf("message verification should fail with other key")
	}
}

func TestWalletSaveLoad(t *testing.T) {
	w, err := NewWallet()
	if err != nil {
		t.Fatalf("new wallet: %v", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "wallet.json")
	if err := w.Save(path, "pass"); err != nil {
		t.Fatalf("save: %v", err)
	}
	w2, err := LoadWallet(path, "pass")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if w.Address != w2.Address {
		t.Fatalf("addresses differ")
	}
	if !bytes.Equal(w.PublicKeyBytes(), w2.PublicKeyBytes()) {
		t.Fatalf("public keys differ")
	}
}

func TestWalletDeterministicSeed(t *testing.T) {
	seed := bytes.Repeat([]byte{0xAB}, 32)
	w1, err := NewWalletFromSeed(seed)
	if err != nil {
		t.Fatalf("seed wallet: %v", err)
	}
	w2, err := NewWalletFromSeed(seed)
	if err != nil {
		t.Fatalf("seed wallet: %v", err)
	}
	if w1.Address != w2.Address {
		t.Fatalf("deterministic wallets should match")
	}
	other, err := NewWalletFromSeed(bytes.Repeat([]byte{0xAC}, 32))
	if err != nil {
		t.Fatalf("seed wallet: %v", err)
	}
	if other.Address == w1.Address {
		t.Fatalf("different seeds should produce different addresses")
	}
}

func TestWalletSharedSecret(t *testing.T) {
	w1, err := NewWallet()
	if err != nil {
		t.Fatalf("new wallet: %v", err)
	}
	w2, err := NewWallet()
	if err != nil {
		t.Fatalf("new wallet: %v", err)
	}
	s1, err := w1.DeriveSharedSecret(&w2.PublicKey)
	if err != nil {
		t.Fatalf("shared secret: %v", err)
	}
	s2, err := w2.DeriveSharedSecret(&w1.PublicKey)
	if err != nil {
		t.Fatalf("shared secret: %v", err)
	}
	if !bytes.Equal(s1, s2) {
		t.Fatalf("shared secrets differ")
	}
}
