package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// TestVerifySignature ensures the helper validates ECDSA signatures generated
// with P-256 keys and rejects tampered data.
func TestVerifySignature(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	msg := []byte("hello dao")
	r, s, err := ecdsa.Sign(rand.Reader, priv, msgHash(msg))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	sig := append(r.Bytes(), s.Bytes()...)

	pubBytes := append([]byte{4}, priv.PublicKey.X.Bytes()...)
	pubBytes = append(pubBytes, priv.PublicKey.Y.Bytes()...)
	pubHex := hex.EncodeToString(pubBytes)
	msgHex := hex.EncodeToString(msg)
	sigHex := hex.EncodeToString(sig)

	ok, err := VerifySignature(pubHex, msgHex, sigHex)
	if err != nil || !ok {
		t.Fatalf("signature should verify: %v", err)
	}

	// Tamper with the signature
	sigHex = hex.EncodeToString(append(sig[:len(sig)-1], sig[len(sig)-1]^0x01))
	ok, err = VerifySignature(pubHex, msgHex, sigHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected signature verification to fail")
	}
}

func msgHash(msg []byte) []byte {
	h := sha256.Sum256(msg)
	return h[:]
}
