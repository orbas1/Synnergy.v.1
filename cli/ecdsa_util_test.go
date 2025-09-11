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
        rBytes := r.Bytes()
        sBytes := s.Bytes()
        rPad := append(make([]byte, 32-len(rBytes)), rBytes...)
        sPad := append(make([]byte, 32-len(sBytes)), sBytes...)
        sig := append(rPad, sPad...)

        xBytes := priv.PublicKey.X.Bytes()
        yBytes := priv.PublicKey.Y.Bytes()
        xPad := append(make([]byte, 32-len(xBytes)), xBytes...)
        yPad := append(make([]byte, 32-len(yBytes)), yBytes...)
        pubBytes := append([]byte{4}, xPad...)
        pubBytes = append(pubBytes, yPad...)
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
