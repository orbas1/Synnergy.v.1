package synnergy

import "testing"

func TestBiometricSecurityNode(t *testing.T) {
	auth := NewBiometricsAuth()
	node := NewBiometricSecurityNode("node1", auth)
	addr := "addr1"
	bio := []byte("biometric")
	node.Enroll(addr, bio)
	if !node.Authenticate(addr, bio) {
		t.Fatalf("authentication failed")
	}
	executed := false
	err := node.SecureExecute(addr, bio, func() error {
		executed = true
		return nil
	})
	if err != nil || !executed {
		t.Fatalf("secure execute: %v, executed=%v", err, executed)
	}
	node.Remove(addr)
	if node.Authenticate(addr, bio) {
		t.Fatalf("authentication should fail after removal")
	}
}
