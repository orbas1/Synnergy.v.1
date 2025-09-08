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
	hash := sha256.Sum256(data)
	sig := ed25519.Sign(priv, hash[:])

	if svc.Enrolled(user) {
		t.Fatal("expected user to be unenrolled initially")
	}
	if err := svc.Enroll(user, data, pub); err != nil {
		t.Fatalf("enroll: %v", err)
	}
	if err := svc.Enroll(user, data, pub); err != ErrAlreadyRegistered {
		t.Fatalf("expected ErrAlreadyRegistered, got %v", err)
	}
	if !svc.Enrolled(user) {
		t.Fatal("expected user to be enrolled")
	}
	if !svc.Verify(user, data, sig) {
		t.Fatalf("expected verification to succeed")
	}
	if svc.Verify(user, []byte("wrong"), sig) {
		t.Fatalf("expected verification to fail for wrong data")
	}
	if err := svc.Enroll("", data, pub); err != ErrEmptyUserID {
		t.Fatalf("expected ErrEmptyUserID, got %v", err)
	}
	if err := svc.Enroll("bob", nil, pub); err != ErrInvalidBiometric {
		t.Fatalf("expected ErrInvalidBiometric, got %v", err)
	}
	if err := svc.Enroll("bob", data, ed25519.PublicKey{}); err != ErrInvalidPublicKey {
		t.Fatalf("expected ErrInvalidPublicKey, got %v", err)
	}
	ids := svc.List()
	if len(ids) != 1 || ids[0] != user {
		t.Fatalf("unexpected list contents: %#v", ids)
	}
	svc.Remove(user)
	if svc.Enrolled(user) {
		t.Fatal("expected user to be removed")
	}
}
