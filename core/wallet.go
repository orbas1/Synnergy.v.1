package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

// Wallet holds keys used to sign transactions.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

// NewWallet generates a new wallet with a random key pair.
func NewWallet() (*Wallet, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	addr := sha256.Sum256(elliptic.Marshal(elliptic.P256(), pk.PublicKey.X, pk.PublicKey.Y))
	return &Wallet{PrivateKey: pk, Address: string(addr[:])}, nil
}

// Sign signs the transaction data.
func (w *Wallet) Sign(tx *Transaction) ([]byte, error) {
	h := sha256.Sum256([]byte(tx.From + tx.To))
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, h[:])
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

// VerifySignature verifies the signature for the transaction.
func VerifySignature(tx *Transaction, sig []byte, pub *ecdsa.PublicKey) bool {
	h := sha256.Sum256([]byte(tx.From + tx.To))
	r := new(big.Int).SetBytes(sig[:len(sig)/2])
	s := new(big.Int).SetBytes(sig[len(sig)/2:])
	return ecdsa.Verify(pub, h[:], r, s)
}
