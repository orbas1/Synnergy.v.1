package p2p

import "testing"

func TestPFSChannel(t *testing.T) {
	alice := NewPFSChannel()
	bob := NewPFSChannel()
	if err := alice.SetRemotePublicKey(bob.LocalPublicKey()); err != nil {
		t.Fatalf("alice set remote: %v", err)
	}
	if err := bob.SetRemotePublicKey(alice.LocalPublicKey()); err != nil {
		t.Fatalf("bob set remote: %v", err)
	}
	ciphertext, err := alice.Encrypt([]byte("hello"), []byte("aad"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	plaintext, err := bob.Decrypt(ciphertext, []byte("aad"))
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(plaintext) != "hello" {
		t.Fatalf("unexpected plaintext: %s", plaintext)
	}
	if len(alice.SessionFingerprint()) == 0 {
		t.Fatalf("fingerprint missing")
	}
}

func TestPFSChannelRemoteRequired(t *testing.T) {
	ch := NewPFSChannel()
	if _, err := ch.Encrypt([]byte("data"), nil); err == nil {
		t.Fatalf("expected remote key error")
	}
}
