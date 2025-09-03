package synnergy

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func TestBiometricSecurityNode(t *testing.T) {
	auth := NewBiometricsAuth()
	node := NewBiometricSecurityNode("node1", auth)
	addr := "addr1"
	bio := []byte("biometric")
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	node.Enroll(addr, bio, &key.PublicKey)
	h := sha256.Sum256(bio)
	sig, err := ecdsa.SignASN1(rand.Reader, key, h[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if !node.Authenticate(addr, bio, sig) {
		t.Fatalf("authentication failed")
	}
	executed := false
	err = node.SecureExecute(addr, bio, sig, func() error {
		executed = true
		return nil
	})
	if err != nil || !executed {
		t.Fatalf("secure execute: %v, executed=%v", err, executed)
	}
	node.Remove(addr)
	if node.Authenticate(addr, bio, sig) {
		t.Fatalf("authentication should fail after removal")
	}
}
