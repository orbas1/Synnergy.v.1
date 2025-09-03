package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// Wallet holds keys used to sign transactions.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	// Address is the hex encoded SHA-256 hash of the public key.
	Address string
}

// NewWallet generates a new wallet with a random key pair.
func NewWallet() (*Wallet, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	// Derive a deterministic address from the uncompressed public key and
	// return it as a hex string. Only the first 20 bytes are used to keep
	// addresses short while preserving collision resistance for this
	// example implementation.
	xBytes := pk.PublicKey.X.Bytes()
	yBytes := pk.PublicKey.Y.Bytes()
	pub := append(append([]byte{0x04}, xBytes...), yBytes...)
	hash := sha256.Sum256(pub)
	addr := hex.EncodeToString(hash[:20])
	return &Wallet{PrivateKey: pk, Address: addr}, nil
}

// Sign signs the transaction hash with the wallet's private key.
func (w *Wallet) Sign(tx *Transaction) ([]byte, error) {
	h, err := hex.DecodeString(tx.Hash())
	if err != nil {
		return nil, err
	}
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, h)
	if err != nil {
		return nil, err
	}
	sig := append(r.Bytes(), s.Bytes()...)
	tx.Signature = sig
	return sig, nil
}

// VerifySignature verifies the signature for the transaction using the public
// key provided.
func VerifySignature(tx *Transaction, sig []byte, pub *ecdsa.PublicKey) bool {
	h, err := hex.DecodeString(tx.Hash())
	if err != nil {
		return false
	}
	r := new(big.Int).SetBytes(sig[:len(sig)/2])
	s := new(big.Int).SetBytes(sig[len(sig)/2:])
	return ecdsa.Verify(pub, h, r, s)
}
