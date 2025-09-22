package p2p

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"sync"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

var (
	errRemoteKeyRequired  = errors.New("p2p: remote public key required")
	errCiphertextTooShort = errors.New("p2p: ciphertext too short")
)

// PFSChannel implements perfect forward secrecy using X25519 and XChaCha20. Each
// message uses an ephemeral key ensuring compromise of long-term secrets does not
// reveal historical ciphertexts.
type PFSChannel struct {
	mu        sync.RWMutex
	curve     ecdh.Curve
	localPriv *ecdh.PrivateKey
	localPub  []byte
	remotePub *ecdh.PublicKey
}

// NewPFSChannel creates a new channel with a freshly generated static key pair.
func NewPFSChannel() *PFSChannel {
	curve := ecdh.X25519()
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		panic(fmt.Errorf("p2p: x25519 keygen failed: %w", err))
	}
	return &PFSChannel{
		curve:     curve,
		localPriv: priv,
		localPub:  append([]byte(nil), priv.PublicKey().Bytes()...),
	}
}

// LocalPublicKey returns the static public key for handshake distribution.
func (c *PFSChannel) LocalPublicKey() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]byte(nil), c.localPub...)
}

// SetRemotePublicKey registers the remote static key.
func (c *PFSChannel) SetRemotePublicKey(pub []byte) error {
	if len(pub) == 0 {
		return errRemoteKeyRequired
	}
	remote, err := c.curve.NewPublicKey(pub)
	if err != nil {
		return fmt.Errorf("p2p: invalid remote key: %w", err)
	}
	c.mu.Lock()
	c.remotePub = remote
	c.mu.Unlock()
	return nil
}

// Encrypt seals the plaintext using an ephemeral key and returns a byte slice
// containing version, ephemeral public key, nonce, and ciphertext.
func (c *PFSChannel) Encrypt(msg []byte, aad []byte) ([]byte, error) {
	c.mu.RLock()
	remote := c.remotePub
	c.mu.RUnlock()
	if remote == nil {
		return nil, errRemoteKeyRequired
	}
	ephemeral, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("p2p: ephemeral keygen failed: %w", err)
	}
	secret, err := ephemeral.ECDH(remote)
	if err != nil {
		return nil, fmt.Errorf("p2p: ecdh failed: %w", err)
	}
	key, err := deriveKey(secret)
	if err != nil {
		return nil, err
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("p2p: aead init failed: %w", err)
	}
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("p2p: nonce generation failed: %w", err)
	}
	ciphertext := aead.Seal(nil, nonce, msg, aad)
	out := make([]byte, 1+len(ephemeral.PublicKey().Bytes())+len(nonce)+len(ciphertext))
	out[0] = 1
	copy(out[1:], ephemeral.PublicKey().Bytes())
	offset := 1 + len(ephemeral.PublicKey().Bytes())
	copy(out[offset:], nonce)
	offset += len(nonce)
	copy(out[offset:], ciphertext)
	return out, nil
}

// Decrypt verifies and opens the ciphertext produced by Encrypt.
func (c *PFSChannel) Decrypt(data []byte, aad []byte) ([]byte, error) {
	if len(data) < 1+32+chacha20poly1305.NonceSizeX {
		return nil, errCiphertextTooShort
	}
	version := data[0]
	if version != 1 {
		return nil, fmt.Errorf("p2p: unsupported pfs version %d", version)
	}
	ephemeralBytes := data[1 : 1+32]
	nonceStart := 1 + 32
	nonceEnd := nonceStart + chacha20poly1305.NonceSizeX
	nonce := data[nonceStart:nonceEnd]
	ciphertext := data[nonceEnd:]

	ephemeral, err := c.curve.NewPublicKey(ephemeralBytes)
	if err != nil {
		return nil, fmt.Errorf("p2p: invalid ephemeral key: %w", err)
	}
	c.mu.RLock()
	local := c.localPriv
	c.mu.RUnlock()
	secret, err := local.ECDH(ephemeral)
	if err != nil {
		return nil, fmt.Errorf("p2p: ecdh failed: %w", err)
	}
	key, err := deriveKey(secret)
	if err != nil {
		return nil, err
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("p2p: aead init failed: %w", err)
	}
	plaintext, err := aead.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("p2p: decrypt failed: %w", err)
	}
	return plaintext, nil
}

// SessionFingerprint returns a stable hash identifying the channel pairing.
func (c *PFSChannel) SessionFingerprint() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	h := sha256.New()
	h.Write(c.localPub)
	if c.remotePub != nil {
		h.Write(c.remotePub.Bytes())
	}
	return h.Sum(nil)
}

func deriveKey(secret []byte) ([]byte, error) {
	reader := hkdf.New(sha256.New, secret, []byte("synnergy-pfs"), nil)
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := io.ReadFull(reader, key); err != nil {
		return nil, fmt.Errorf("p2p: hkdf failed: %w", err)
	}
	return key, nil
}
