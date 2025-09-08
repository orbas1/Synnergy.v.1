package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestBiometricService(t *testing.T) {
	svc := NewBiometricService()
	user := "alice"
	data := []byte("fingerprint")
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	svc.Enroll(user, data, pub)
	hash := sha256.Sum256(data)
	sig := ed25519.Sign(priv, hash[:])
	if !svc.Verify(user, data, sig) {
		t.Fatalf("expected verification to succeed")
	}
	if svc.Verify(user, []byte("wrong"), sig) {
		t.Fatalf("expected verification to fail for wrong data")
	}
}
