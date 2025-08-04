package core

import "testing"

func TestBiometricService(t *testing.T) {
	svc := NewBiometricService()
	user := "alice"
	data := []byte("fingerprint")
	svc.Enroll(user, data)
	if !svc.Verify(user, data) {
		t.Fatalf("expected verification to succeed")
	}
	if svc.Verify(user, []byte("wrong")) {
		t.Fatalf("expected verification to fail for wrong data")
	}
}
