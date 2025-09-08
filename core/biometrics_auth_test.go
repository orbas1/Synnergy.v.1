package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestBiometricsAuth(t *testing.T) {
	auth := NewBiometricsAuth()
	addr := "user1"
	data := []byte("fingerprint")
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	hash := sha256.Sum256(data)
	sig := ed25519.Sign(priv, hash[:])

	if auth.Verify(addr, data, sig) {
		t.Fatal("expected verification to fail before enrollment")
	}

	if err := auth.Enroll(addr, data, pub); err != nil {
		t.Fatalf("enroll: %v", err)
	}

	if err := auth.Enroll(addr, data, pub); err != ErrAlreadyEnrolled {
		t.Fatalf("expected ErrAlreadyEnrolled, got %v", err)
	}

	if !auth.Verify(addr, data, sig) {
		t.Fatal("expected verification to succeed after enrollment")
	}
	if !auth.Enrolled(addr) {
		t.Fatal("expected address to be enrolled")
	}
	list := auth.List()
	if len(list) != 1 || list[0] != addr {
		t.Fatalf("unexpected list contents: %#v", list)
	}

	if err := auth.Enroll("", data, pub); err != ErrAddressRequired {
		t.Fatalf("expected ErrAddressRequired, got %v", err)
	}
	if err := auth.Enroll("other", nil, pub); err != ErrInvalidBiometric {
		t.Fatalf("expected ErrInvalidBiometric, got %v", err)
	}
	if err := auth.Enroll("other", data, ed25519.PublicKey{}); err != ErrInvalidPublicKey {
		t.Fatalf("expected ErrInvalidPublicKey, got %v", err)
	}

	auth.Remove(addr)
	if auth.Verify(addr, data, sig) {
		t.Fatal("expected verification to fail after removal")
	}
	if auth.Enrolled(addr) {
		t.Fatal("expected address to be unenrolled")
	}
}
