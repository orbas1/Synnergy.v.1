package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
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
