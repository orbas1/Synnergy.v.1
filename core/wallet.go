package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"

	"golang.org/x/crypto/scrypt"
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

// Save writes the wallet's private key encrypted with the provided password to path.
// AES-256 GCM with an scrypt derived key provides confidentiality and integrity.
func (w *Wallet) Save(path, password string) error {
	if password == "" {
		return errors.New("password required")
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	key, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	priv := w.PrivateKey.D.Bytes()
	ct := gcm.Seal(nil, nonce, priv, nil)
	data, err := json.Marshal(struct {
		Address string `json:"address"`
		Salt    []byte `json:"salt"`
		Nonce   []byte `json:"nonce"`
		Key     []byte `json:"key"`
	}{w.Address, salt, nonce, ct})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0o600)
}

// LoadWallet decrypts a wallet file previously written with Save.
func LoadWallet(path, password string) (*Wallet, error) {
	if password == "" {
		return nil, errors.New("password required")
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var enc struct {
		Address string `json:"address"`
		Salt    []byte `json:"salt"`
		Nonce   []byte `json:"nonce"`
		Key     []byte `json:"key"`
	}
	if err := json.Unmarshal(b, &enc); err != nil {
		return nil, err
	}
	key, err := scrypt.Key([]byte(password), enc.Salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	priv, err := gcm.Open(nil, enc.Nonce, enc.Key, nil)
	if err != nil {
		return nil, err
	}
	d := new(big.Int).SetBytes(priv)
	pk := new(ecdsa.PrivateKey)
	pk.PublicKey.Curve = elliptic.P256()
	pk.D = d
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(priv)
	return &Wallet{PrivateKey: pk, Address: enc.Address}, nil
}
