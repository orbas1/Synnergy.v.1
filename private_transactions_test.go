package synnergy

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"
)

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, aes256KeySize)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("rand read: %v", err)
	}
	payload := []byte("hello world")
	cipherText, err := Encrypt(key, payload)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	plain, err := Decrypt(key, cipherText)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(plain, payload) {
		t.Fatalf("unexpected plaintext: %v", plain)
	}
}

func TestEncryptWithAAD(t *testing.T) {
	key := make([]byte, aes256KeySize)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("rand read: %v", err)
	}
	aad := []byte("consensus:epoch=9")
	env, err := EncryptWithAAD(key, []byte("payload"), aad)
	if err != nil {
		t.Fatalf("encrypt with aad: %v", err)
	}
	plain, err := DecryptWithAAD(key, env)
	if err != nil {
		t.Fatalf("decrypt with aad: %v", err)
	}
	if !bytes.Equal(plain, []byte("payload")) {
		t.Fatalf("unexpected plaintext: %v", plain)
	}
	env.Nonce[0] ^= 0xFF
	if _, err := DecryptWithAAD(key, env); err == nil {
		t.Fatalf("expected failure after tampering")
	}
}

func TestEncryptRejectsInvalidKey(t *testing.T) {
	if _, err := Encrypt(make([]byte, 16), []byte("data")); err == nil {
		t.Fatalf("expected invalid key length error")
	}
}

func TestPrivateTxManagerStoreAndGet(t *testing.T) {
	mgr := NewPrivateTxManager()
	tx := PrivateTransaction{Payload: []byte("data"), Nonce: []byte("nonce")}
	if err := mgr.Store(tx); err != nil {
		t.Fatalf("store: %v", err)
	}
	stored, ok := mgr.Get(deriveTransactionID(tx))
	if !ok {
		t.Fatalf("transaction not stored")
	}
	if !bytes.Equal(stored.Payload, tx.Payload) {
		t.Fatalf("payload mismatch")
	}
	if err := mgr.Store(tx); err != ErrDuplicateTransaction {
		t.Fatalf("expected duplicate error, got %v", err)
	}

	tx.ID = "manual"
	tx.Timestamp = time.Now().Add(-time.Hour)
	mgr.Upsert(tx)
	list := mgr.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(list))
	}
	if list[0].Timestamp.After(list[1].Timestamp) {
		t.Fatalf("transactions not ordered by timestamp")
	}
}

func TestSignatureLifecycle(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tx := PrivateTransaction{ID: "tx1", Payload: []byte("cipher"), Nonce: []byte("nonce"), Timestamp: time.Now()}
	if err := SignTransaction(priv, &tx); err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := VerifyTransaction(pub, tx); err != nil {
		t.Fatalf("verify: %v", err)
	}
	tx.Signature[0] ^= 0xFF
	if err := VerifyTransaction(pub, tx); err == nil {
		t.Fatalf("expected signature mismatch")
	}
}
