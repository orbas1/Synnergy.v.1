package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// parsePubKey decodes an uncompressed hex encoded P-256 public key.
func parsePubKey(hexStr string) (*ecdsa.PublicKey, error) {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	if len(b) != 65 || b[0] != 4 {
		return nil, fmt.Errorf("invalid public key")
	}
	x := new(big.Int).SetBytes(b[1:33])
	y := new(big.Int).SetBytes(b[33:])
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// decodeSig decodes a hex encoded signature.
func decodeSig(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// VerifySignature verifies a hex encoded ECDSA signature over the provided
// message using an uncompressed hex encoded P-256 public key. The message is
// hashed with SHA-256 prior to verification. It returns true if the signature
// is valid for the given message and key.
func VerifySignature(pubHex, msgHex, sigHex string) (bool, error) {
	pub, err := parsePubKey(pubHex)
	if err != nil {
		return false, err
	}
	msg, err := hex.DecodeString(msgHex)
	if err != nil {
		return false, err
	}
	sig, err := decodeSig(sigHex)
	if err != nil {
		return false, err
	}
	if len(sig)%2 != 0 {
		return false, fmt.Errorf("invalid signature length")
	}
	n := len(sig) / 2
	r := new(big.Int).SetBytes(sig[:n])
	s := new(big.Int).SetBytes(sig[n:])
	digest := sha256.Sum256(msg)
	return ecdsa.Verify(pub, digest[:], r, s), nil
}
