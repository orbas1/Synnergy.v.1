package core

import "testing"

func TestBiometricsAuth(t *testing.T) {
	auth := NewBiometricsAuth()
	addr := "user1"
	data := []byte("fingerprint")

	if auth.Verify(addr, data) {
		t.Fatal("expected verification to fail before enrollment")
	}

	auth.Enroll(addr, data)
	if !auth.Verify(addr, data) {
		t.Fatal("expected verification to succeed after enrollment")
	}

	auth.Remove(addr)
	if auth.Verify(addr, data) {
		t.Fatal("expected verification to fail after removal")
	}
}
