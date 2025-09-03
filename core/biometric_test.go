package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestBiometricService(t *testing.T) {
	svc := NewBiometricService()
	user := "alice"
	data := []byte("fingerprint")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	svc.Enroll(user, data, &key.PublicKey)
	hash := sha256.Sum256(data)
	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if !svc.Verify(user, data, sig) {
		t.Fatalf("expected verification to succeed")
	}
	if svc.Verify(user, []byte("wrong"), sig) {
		t.Fatalf("expected verification to fail for wrong data")
	}
}
