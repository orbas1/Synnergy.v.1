package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

func TestEnvelopeEncryptorSealAndOpen(t *testing.T) {
	master := make([]byte, 32)
	if _, err := rand.Read(master); err != nil {
		t.Fatalf("entropy: %v", err)
	}
	_, signer, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	enc, err := NewEnvelopeEncryptor(master, signer, "test")
	if err != nil {
		t.Fatalf("encryptor: %v", err)
	}
	aad := []byte("round-1")
	env, err := enc.Seal([]byte("payload"), aad)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	out, err := enc.Open(env, aad)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if string(out) != "payload" {
		t.Fatalf("unexpected plaintext: %s", out)
	}
	if len(enc.Fingerprint()) == 0 {
		t.Fatalf("fingerprint missing")
	}
}

func TestEnvelopeEncryptorRotate(t *testing.T) {
	master := make([]byte, 32)
	if _, err := rand.Read(master); err != nil {
		t.Fatalf("entropy: %v", err)
	}
	enc, err := NewEnvelopeEncryptor(master, nil, "initial")
	if err != nil {
		t.Fatalf("encryptor: %v", err)
	}
	env, err := enc.Seal([]byte("data"), nil)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	if _, err := enc.Open(env, nil); err != nil {
		t.Fatalf("open pre-rotation: %v", err)
	}
	next := make([]byte, 32)
	if _, err := rand.Read(next); err != nil {
		t.Fatalf("entropy: %v", err)
	}
	if err := enc.Rotate(next, nil, "rotated"); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	// old envelope should fail due to key id mismatch
	if _, err := enc.Open(env, nil); err == nil {
		t.Fatalf("expected error for old envelope")
	}
	env2, err := enc.Seal([]byte("fresh"), []byte("aad"))
	if err != nil {
		t.Fatalf("seal post-rotation: %v", err)
	}
	out, err := enc.Open(env2, []byte("aad"))
	if err != nil {
		t.Fatalf("open post-rotation: %v", err)
	}
	if string(out) != "fresh" {
		t.Fatalf("unexpected plaintext: %s", out)
	}
}

func TestEnvelopeEncryptorInvalidSignature(t *testing.T) {
	master := make([]byte, 32)
	rand.Read(master)
	enc, err := NewEnvelopeEncryptor(master, nil, "sig")
	if err != nil {
		t.Fatalf("encryptor: %v", err)
	}
	env, err := enc.Seal([]byte("message"), nil)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	env.Signature[0] ^= 0xFF
	if _, err := enc.Open(env, nil); err == nil {
		t.Fatalf("expected signature error")
	}
}
