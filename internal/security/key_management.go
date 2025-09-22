package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

// KeyPurpose describes what a key is used for. It allows the CLI, consensus,
// VM and wallet subsystems to request key material deterministically without
// sharing the raw key map.
type KeyPurpose string

const (
	PurposeNoiseStatic    KeyPurpose = "noise_static"
	PurposeNoiseEphemeral            = "noise_ephemeral"
	PurposeTLS                       = "tls"
	PurposeEnvelope                  = "envelope"
	PurposeStateSigning              = "state_signing"
)

type symmetricKey struct {
	material []byte
	version  int
	rotated  time.Time
	actor    string
}

type signingKey struct {
	private ed25519.PrivateKey
	public  ed25519.PublicKey
	version int
	rotated time.Time
	actor   string
}

// AuditEntry captures key lifecycle information. The log is intentionally kept
// simple so it can be mirrored in the on-chain governance history.
type AuditEntry struct {
	Purpose   KeyPurpose
	Version   int
	RotatedAt time.Time
	Actor     string
}

// KeyManager manages symmetric and signing keys with strict concurrency
// guarantees. It also keeps an audit trail so operators can prove compliance.
type KeyManager struct {
	mu         sync.RWMutex
	symmetric  map[KeyPurpose]symmetricKey
	signing    map[KeyPurpose]signingKey
	auditLog   []AuditEntry
	entropy    io.Reader
	keyFactory func(io.Reader) (ed25519.PublicKey, ed25519.PrivateKey, error)
}

// NewKeyManager constructs a KeyManager with crypto/rand entropy.
func NewKeyManager() *KeyManager {
	return &KeyManager{
		symmetric: make(map[KeyPurpose]symmetricKey),
		signing:   make(map[KeyPurpose]signingKey),
		entropy:   rand.Reader,
		keyFactory: func(r io.Reader) (ed25519.PublicKey, ed25519.PrivateKey, error) {
			return ed25519.GenerateKey(r)
		},
	}
}

// WithEntropy allows tests to inject deterministic entropy sources.
func (k *KeyManager) WithEntropy(r io.Reader) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if r == nil {
		r = rand.Reader
	}
	k.entropy = r
}

// WithKeyFactory overrides the Ed25519 key factory, primarily used for tests.
func (k *KeyManager) WithKeyFactory(factory func(io.Reader) (ed25519.PublicKey, ed25519.PrivateKey, error)) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if factory == nil {
		factory = func(r io.Reader) (ed25519.PublicKey, ed25519.PrivateKey, error) {
			return ed25519.GenerateKey(r)
		}
	}
	k.keyFactory = factory
}

func (k *KeyManager) recordAudit(purpose KeyPurpose, version int, actor string) {
	k.auditLog = append(k.auditLog, AuditEntry{
		Purpose:   purpose,
		Version:   version,
		RotatedAt: time.Now().UTC(),
		Actor:     actor,
	})
}

// GenerateSymmetricKey produces a new 32-byte key for the provided purpose. The
// actor string identifies the subsystem that requested the rotation.
func (k *KeyManager) GenerateSymmetricKey(purpose KeyPurpose, actor string) (version int, key []byte, err error) {
	key = make([]byte, 32)
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, err = io.ReadFull(k.entropy, key); err != nil {
		return 0, nil, fmt.Errorf("security: entropy failure: %w", err)
	}
	rec := k.symmetric[purpose]
	rec.material = append([]byte(nil), key...)
	rec.version++
	rec.rotated = time.Now().UTC()
	rec.actor = actor
	k.symmetric[purpose] = rec
	k.recordAudit(purpose, rec.version, actor)
	return rec.version, append([]byte(nil), rec.material...), nil
}

// SetSymmetricKey sets key material explicitly. It is intended for recovery
// operations where the key is derived from an external HSM.
func (k *KeyManager) SetSymmetricKey(purpose KeyPurpose, material []byte, actor string, version int) error {
	if len(material) == 0 {
		return errors.New("security: symmetric key material required")
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	if version <= 0 {
		version = k.symmetric[purpose].version + 1
	}
	k.symmetric[purpose] = symmetricKey{
		material: append([]byte(nil), material...),
		version:  version,
		rotated:  time.Now().UTC(),
		actor:    actor,
	}
	k.recordAudit(purpose, version, actor)
	return nil
}

// SymmetricKey returns the material and version for the provided purpose.
func (k *KeyManager) SymmetricKey(purpose KeyPurpose) ([]byte, int, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	rec, ok := k.symmetric[purpose]
	if !ok {
		return nil, 0, fmt.Errorf("security: symmetric key for %s not found", purpose)
	}
	return append([]byte(nil), rec.material...), rec.version, nil
}

// GenerateSigningKey creates a new Ed25519 key pair for the given purpose.
func (k *KeyManager) GenerateSigningKey(purpose KeyPurpose, actor string) (pub ed25519.PublicKey, version int, err error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	pub, priv, err := k.keyFactory(k.entropy)
	if err != nil {
		return nil, 0, fmt.Errorf("security: ed25519 keygen failed: %w", err)
	}
	rec := k.signing[purpose]
	rec.private = priv
	rec.public = append(ed25519.PublicKey(nil), pub...)
	rec.version++
	rec.rotated = time.Now().UTC()
	rec.actor = actor
	k.signing[purpose] = rec
	k.recordAudit(purpose, rec.version, actor)
	return append(ed25519.PublicKey(nil), pub...), rec.version, nil
}

// SigningKey retrieves the private key, public key and version for the provided
// purpose.
func (k *KeyManager) SigningKey(purpose KeyPurpose) (ed25519.PrivateKey, ed25519.PublicKey, int, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	rec, ok := k.signing[purpose]
	if !ok {
		return nil, nil, 0, fmt.Errorf("security: signing key for %s not found", purpose)
	}
	return append(ed25519.PrivateKey(nil), rec.private...), append(ed25519.PublicKey(nil), rec.public...), rec.version, nil
}

// Sign uses the key registered under purpose to sign the provided payload. The
// function returns the signature together with the public key and key version so
// that verifiers can perform strict validation.
func (k *KeyManager) Sign(purpose KeyPurpose, payload []byte) ([]byte, ed25519.PublicKey, int, error) {
	priv, pub, version, err := k.SigningKey(purpose)
	if err != nil {
		return nil, nil, 0, err
	}
	sig := ed25519.Sign(priv, payload)
	return sig, pub, version, nil
}

// Verify checks a signature with either the stored public key or an explicitly
// provided one. When pub is nil the function uses the latest registered key.
func (k *KeyManager) Verify(purpose KeyPurpose, payload, signature []byte, pub ed25519.PublicKey) error {
	k.mu.RLock()
	defer k.mu.RUnlock()
	rec, ok := k.signing[purpose]
	if !ok {
		return fmt.Errorf("security: signing key for %s not found", purpose)
	}
	if pub == nil {
		pub = rec.public
	}
	if !ed25519.Verify(pub, payload, signature) {
		return errors.New("security: signature verification failed")
	}
	return nil
}

// AuditLog returns a copy of the audit entries accumulated since start-up.
func (k *KeyManager) AuditLog() []AuditEntry {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]AuditEntry, len(k.auditLog))
	copy(out, k.auditLog)
	return out
}
