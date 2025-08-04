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
	if !auth.Enrolled(addr) {
		t.Fatal("expected address to be enrolled")
	}
	list := auth.List()
	if len(list) != 1 || list[0] != addr {
		t.Fatalf("unexpected list contents: %#v", list)
	}

	auth.Remove(addr)
	if auth.Verify(addr, data) {
		t.Fatal("expected verification to fail after removal")
	}
	if auth.Enrolled(addr) {
		t.Fatal("expected address to be unenrolled")
	}
}
