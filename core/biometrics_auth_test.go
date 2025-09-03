package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestBiometricsAuth(t *testing.T) {
	auth := NewBiometricsAuth()
	addr := "user1"
	data := []byte("fingerprint")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	hash := sha256.Sum256(data)
	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	if auth.Verify(addr, data, sig) {
		t.Fatal("expected verification to fail before enrollment")
	}

	auth.Enroll(addr, data, &key.PublicKey)
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

	auth.Remove(addr)
	if auth.Verify(addr, data, sig) {
		t.Fatal("expected verification to fail after removal")
	}
	if auth.Enrolled(addr) {
		t.Fatal("expected address to be unenrolled")
	}
}
