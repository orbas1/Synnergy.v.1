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
	"math/big"
	"os"

	"golang.org/x/crypto/scrypt"
)

const (
	walletFileVersion = "1.1"
	scryptN           = 1 << 15
	scryptR           = 8
	scryptP           = 1
)

// Wallet holds keys used to sign transactions and interact with the virtual
// machine. Public keys are stored alongside the private key to enable hardware
// attestation and out-of-band verification flows.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
	Address    string
}

// NewWallet generates a new wallet with a random key pair.
func NewWallet() (*Wallet, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	addr := deriveAddress(&pk.PublicKey)
	return &Wallet{PrivateKey: pk, PublicKey: pk.PublicKey, Address: addr}, nil
}

// NewWalletFromSeed deterministically derives a wallet from the provided seed.
// It is primarily used for test fixtures and secure key recovery procedures.
func NewWalletFromSeed(seed []byte) (*Wallet, error) {
	if len(seed) < 32 {
		return nil, errors.New("seed length must be at least 32 bytes")
	}
	curve := elliptic.P256()
	h := sha256.Sum256(seed)
	d := new(big.Int).SetBytes(h[:])
	n := new(big.Int).Sub(curve.Params().N, big.NewInt(1))
	d.Mod(d, n)
	d.Add(d, big.NewInt(1))

	priv := &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: curve}, D: d}
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	addr := deriveAddress(&priv.PublicKey)
	return &Wallet{PrivateKey: priv, PublicKey: priv.PublicKey, Address: addr}, nil
}

// Sign signs the transaction hash with the wallet's private key.
func (w *Wallet) Sign(tx *Transaction) ([]byte, error) {
	if w == nil || w.PrivateKey == nil {
		return nil, errors.New("wallet private key not initialised")
	}
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

// SignMessage signs arbitrary data by hashing it with SHA-256. The helper is
// used by wallet CLI diagnostics and cross-chain attestations.
func (w *Wallet) SignMessage(msg []byte) ([]byte, error) {
	if w == nil || w.PrivateKey == nil {
		return nil, errors.New("wallet private key not initialised")
	}
	digest := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, digest[:])
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

// VerifySignature verifies the signature for the transaction using the public
// key provided.
func VerifySignature(tx *Transaction, sig []byte, pub *ecdsa.PublicKey) bool {
	h, err := hex.DecodeString(tx.Hash())
	if err != nil {
		return false
	}
	return verifyDigest(h, sig, pub)
}

// VerifyMessage verifies data signed with SignMessage.
func VerifyMessage(msg, sig []byte, pub *ecdsa.PublicKey) bool {
	digest := sha256.Sum256(msg)
	return verifyDigest(digest[:], sig, pub)
}

func verifyDigest(digest []byte, sig []byte, pub *ecdsa.PublicKey) bool {
	if len(sig)%2 != 0 || pub == nil {
		return false
	}
	mid := len(sig) / 2
	r := new(big.Int).SetBytes(sig[:mid])
	s := new(big.Int).SetBytes(sig[mid:])
	return ecdsa.Verify(pub, digest, r, s)
}

// DeriveSharedSecret computes an ECDH shared secret with the peer public key.
// The secret is hashed to ensure uniform length and provide key material for
// symmetric encryption channels.
func (w *Wallet) DeriveSharedSecret(peer *ecdsa.PublicKey) ([]byte, error) {
	if w == nil || w.PrivateKey == nil {
		return nil, errors.New("wallet private key not initialised")
	}
	if peer == nil || peer.X == nil || peer.Y == nil {
		return nil, errors.New("peer public key invalid")
	}
	x, _ := peer.Curve.ScalarMult(peer.X, peer.Y, w.PrivateKey.D.Bytes())
	if x == nil {
		return nil, errors.New("failed to derive shared secret")
	}
	digest := sha256.Sum256(x.Bytes())
	return digest[:], nil
}

// PublicKeyBytes returns the uncompressed public key encoding.
func (w *Wallet) PublicKeyBytes() []byte {
	if w == nil {
		return nil
	}
	return encodePublicKey(&w.PublicKey)
}

// Save writes the wallet's private key encrypted with the provided password to
// path. AES-256 GCM with an scrypt derived key provides confidentiality and
// integrity.
func (w *Wallet) Save(path, password string) error {
	if w == nil || w.PrivateKey == nil {
		return errors.New("wallet private key not initialised")
	}
	if password == "" {
		return errors.New("password required")
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	key, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, 32)
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
	file := walletFile{
		Version:   walletFileVersion,
		Address:   w.Address,
		PublicKey: encodePublicKey(&w.PrivateKey.PublicKey),
		Salt:      salt,
		Nonce:     nonce,
		Key:       ct,
		KDF:       walletKDF{N: scryptN, R: scryptR, P: scryptP},
	}
	data, err := json.Marshal(file)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadWallet decrypts a wallet file previously written with Save.
func LoadWallet(path, password string) (*Wallet, error) {
	if password == "" {
		return nil, errors.New("password required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file walletFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}

	params := file.KDF
	if params.N == 0 {
		params = walletKDF{N: scryptN, R: scryptR, P: scryptP}
	}
	key, err := scrypt.Key([]byte(password), file.Salt, params.N, params.R, params.P, 32)
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
	priv, err := gcm.Open(nil, file.Nonce, file.Key, nil)
	if err != nil {
		return nil, err
	}
	d := new(big.Int).SetBytes(priv)
	pk := &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: elliptic.P256()}, D: d}
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(priv)

	pub := pk.PublicKey
	if len(file.PublicKey) > 0 {
		if decoded, err := decodePublicKey(file.PublicKey); err == nil {
			pub = *decoded
		}
	}
	addr := file.Address
	if addr == "" {
		addr = deriveAddress(&pub)
	}
	return &Wallet{PrivateKey: pk, PublicKey: pub, Address: addr}, nil
}

type walletFile struct {
	Version   string    `json:"version"`
	Address   string    `json:"address"`
	PublicKey []byte    `json:"public_key,omitempty"`
	Salt      []byte    `json:"salt"`
	Nonce     []byte    `json:"nonce"`
	Key       []byte    `json:"key"`
	KDF       walletKDF `json:"kdf,omitempty"`
}

type walletKDF struct {
	N int `json:"n"`
	R int `json:"r"`
	P int `json:"p"`
}

func deriveAddress(pub *ecdsa.PublicKey) string {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return ""
	}
	encoded := encodePublicKey(pub)
	hash := sha256.Sum256(encoded)
	return hex.EncodeToString(hash[:20])
}

func encodePublicKey(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	xb := pub.X.Bytes()
	yb := pub.Y.Bytes()
	out := make([]byte, 1+len(xb)+len(yb))
	out[0] = 0x04
	copy(out[1:], xb)
	copy(out[1+len(xb):], yb)
	return out
}

func decodePublicKey(data []byte) (*ecdsa.PublicKey, error) {
	if len(data) == 0 || data[0] != 0x04 {
		return nil, errors.New("public key encoding invalid")
	}
	if len(data)%2 != 1 {
		return nil, errors.New("public key length invalid")
	}
	half := (len(data) - 1) / 2
	xb := data[1 : 1+half]
	yb := data[1+half:]
	curve := elliptic.P256()
	x := new(big.Int).SetBytes(xb)
	y := new(big.Int).SetBytes(yb)
	if !curve.IsOnCurve(x, y) {
		return nil, errors.New("public key not on curve")
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}
