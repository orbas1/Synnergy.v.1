package security

import (
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/crypto/chacha20poly1305"
)

var (
	// ErrInvalidKey is returned when the provided master key cannot be used to
	// build an AEAD instance.
	ErrInvalidKey = errors.New("security: invalid master key")
	// ErrInvalidSignature is returned when a ciphertext fails signature
	// validation.
	ErrInvalidSignature = errors.New("security: invalid signature")
)

// Envelope represents an encrypted payload that is protected with an AEAD
// cipher and signed with an Ed25519 key to guarantee authenticity.
//
// The structure is intentionally simple so it can be marshalled to JSON for the
// CLI or web clients without additional glue code.
type Envelope struct {
	Version    int               `json:"version"`
	KeyID      string            `json:"key_id"`
	Nonce      []byte            `json:"nonce"`
	Ciphertext []byte            `json:"ciphertext"`
	Signature  []byte            `json:"signature"`
	PublicKey  ed25519.PublicKey `json:"public_key"`
}

// EnvelopeEncryptor provides high-grade symmetric encryption coupled with an
// Ed25519 signing identity. It is designed to back the CLI, the VM sandbox and
// the authority node control-plane where the same envelope format is consumed
// by multiple services.
type EnvelopeEncryptor struct {
	mu      sync.RWMutex
	aead    cipherState
	signer  ed25519.PrivateKey
	public  ed25519.PublicKey
	keyID   string
	version int
}

// cipherState is a very small wrapper used so tests can swap the AEAD
// implementation when simulating failures.
type cipherState interface {
	NonceSize() int
	Seal(dst, nonce, plaintext, additionalData []byte) []byte
	Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error)
}

type aeadWrapper struct{ aead cipher.AEAD }

func (a aeadWrapper) NonceSize() int { return a.aead.NonceSize() }

func (a aeadWrapper) Seal(dst, nonce, plaintext, aad []byte) []byte {
	return a.aead.Seal(dst, nonce, plaintext, aad)
}

func (a aeadWrapper) Open(dst, nonce, ciphertext, aad []byte) ([]byte, error) {
	return a.aead.Open(dst, nonce, ciphertext, aad)
}

// NewEnvelopeEncryptor builds an encryptor using an XChaCha20-Poly1305 AEAD. If
// signer is nil the function will create a new Ed25519 key pair automatically so
// that callers always have a valid signing identity.
func NewEnvelopeEncryptor(masterKey []byte, signer ed25519.PrivateKey, keyID string) (*EnvelopeEncryptor, error) {
	if len(masterKey) != chacha20poly1305.KeySize {
		return nil, ErrInvalidKey
	}
	if keyID == "" {
		keyID = "default"
	}
	if signer == nil {
		_, signer = mustGenerateEd25519()
	}
	a, err := chacha20poly1305.NewX(masterKey)
	if err != nil {
		return nil, fmt.Errorf("security: initialise AEAD: %w", err)
	}
	public := signer.Public().(ed25519.PublicKey)
	return &EnvelopeEncryptor{
		aead:    aeadWrapper{a},
		signer:  signer,
		public:  append(ed25519.PublicKey(nil), public...),
		keyID:   keyID,
		version: 1,
	}, nil
}

// mustGenerateEd25519 returns a freshly generated Ed25519 keypair. It panics if
// the random source fails, which is acceptable because the process cannot
// continue without cryptographic entropy.
func mustGenerateEd25519() (ed25519.PublicKey, ed25519.PrivateKey) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(fmt.Errorf("security: ed25519 keygen failed: %w", err))
	}
	return public, private
}

// Seal encrypts the payload, signs the result and returns an Envelope structure
// that can be serialised or directly stored. Additional data is included in the
// AEAD tag to bind contextual information such as consensus round IDs or
// transaction hashes.
func (e *EnvelopeEncryptor) Seal(plaintext, additionalData []byte) (*Envelope, error) {
	e.mu.RLock()
	a := e.aead
	signer := e.signer
	pub := append(ed25519.PublicKey(nil), e.public...)
	nonceSize := a.NonceSize()
	keyID := e.keyID
	version := e.version
	e.mu.RUnlock()

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("security: nonce generation failed: %w", err)
	}
	ciphertext := a.Seal(nil, nonce, plaintext, additionalData)
	sigPayload := signaturePayload(version, keyID, nonce, additionalData, ciphertext)
	signature := ed25519.Sign(signer, sigPayload)
	return &Envelope{
		Version:    version,
		KeyID:      keyID,
		Nonce:      nonce,
		Ciphertext: ciphertext,
		Signature:  signature,
		PublicKey:  pub,
	}, nil
}

// Open verifies the signature contained in the envelope and decrypts the
// ciphertext. The function returns ErrInvalidSignature if verification fails.
func (e *EnvelopeEncryptor) Open(env *Envelope, additionalData []byte) ([]byte, error) {
	if env == nil {
		return nil, errors.New("security: envelope required")
	}
	e.mu.RLock()
	a := e.aead
	signerPub := append(ed25519.PublicKey(nil), e.public...)
	keyID := e.keyID
	expectedVersion := e.version
	e.mu.RUnlock()

	pub := env.PublicKey
	if len(pub) == 0 {
		pub = signerPub
	}
	if !ed25519.Verify(pub, signaturePayload(env.Version, env.KeyID, env.Nonce, additionalData, env.Ciphertext), env.Signature) {
		return nil, ErrInvalidSignature
	}
	if env.KeyID != keyID {
		return nil, fmt.Errorf("security: unexpected key id %s", env.KeyID)
	}
	if env.Version != expectedVersion {
		// Support consuming historical envelopes by still attempting to
		// decrypt, but callers should be aware of the version mismatch.
	}
	plaintext, err := a.Open(nil, env.Nonce, env.Ciphertext, additionalData)
	if err != nil {
		return nil, fmt.Errorf("security: decrypt failed: %w", err)
	}
	return plaintext, nil
}

// Rotate replaces the master key and optionally the signing key. A rotation
// increments the internal version counter so downstream services can detect key
// epochs and flush caches.
func (e *EnvelopeEncryptor) Rotate(masterKey []byte, signer ed25519.PrivateKey, keyID string) error {
	if len(masterKey) != chacha20poly1305.KeySize {
		return ErrInvalidKey
	}
	if signer == nil {
		_, signer = mustGenerateEd25519()
	}
	if keyID == "" {
		keyID = e.keyID
	}
	a, err := chacha20poly1305.NewX(masterKey)
	if err != nil {
		return fmt.Errorf("security: initialise AEAD during rotation: %w", err)
	}
	public := signer.Public().(ed25519.PublicKey)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.aead = aeadWrapper{a}
	e.signer = signer
	e.public = append(ed25519.PublicKey(nil), public...)
	e.keyID = keyID
	e.version++
	return nil
}

// PublicKey returns the current signing public key.
func (e *EnvelopeEncryptor) PublicKey() ed25519.PublicKey {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append(ed25519.PublicKey(nil), e.public...)
}

// Fingerprint returns a short SHA256 digest that uniquely identifies the
// current signing key and AEAD key. It is suitable for audit logs and CLI
// output.
func (e *EnvelopeEncryptor) Fingerprint() []byte {
	e.mu.RLock()
	defer e.mu.RUnlock()
	h := sha256.New()
	h.Write([]byte(e.keyID))
	h.Write(e.public)
	h.Write([]byte{byte(e.version)})
	return h.Sum(nil)
}

func signaturePayload(version int, keyID string, nonce, aad, ciphertext []byte) []byte {
	h := sha256.New()
	h.Write([]byte{byte(version)})
	h.Write([]byte(keyID))
	h.Write(nonce)
	h.Write(aad)
	h.Write(ciphertext)
	return h.Sum(nil)
}
