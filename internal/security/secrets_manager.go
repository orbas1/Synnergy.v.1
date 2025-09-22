package security

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
)

// ErrSecretNotFound is returned when a secret is missing from the manager.
var ErrSecretNotFound = errors.New("security: secret not found")

// ErrSecretExpired indicates the secret exists but has passed its retention
// window.
var ErrSecretExpired = errors.New("security: secret expired")

// SecretsManagerOption customises the behaviour of the manager.
type SecretsManagerOption func(*SecretsManager)

// WithMasterKey installs a deterministic master key. The key must be
// chacha20poly1305.KeySize bytes in length.
func WithMasterKey(key []byte) SecretsManagerOption {
	return func(s *SecretsManager) {
		if len(key) == chacha20poly1305.KeySize {
			copy(s.masterKey[:], key)
		}
	}
}

// WithSecretsClock overrides the time source, primarily used in tests.
func WithSecretsClock(clock func() time.Time) SecretsManagerOption {
	return func(s *SecretsManager) {
		if clock != nil {
			s.now = clock
		}
	}
}

// StoreOption configures a single secret write operation.
type StoreOption func(*storeOptions)

type storeOptions struct {
	ttl      time.Duration
	metadata map[string]string
}

// WithTTL configures the retention period for a secret. Once elapsed the secret
// is automatically purged.
func WithTTL(ttl time.Duration) StoreOption {
	return func(o *storeOptions) {
		if ttl > 0 {
			o.ttl = ttl
		}
	}
}

// WithMetadata stores additional key/value pairs alongside the encrypted
// secret. Metadata is copied to avoid callers mutating the stored map.
func WithMetadata(metadata map[string]string) StoreOption {
	return func(o *storeOptions) {
		if len(metadata) == 0 {
			return
		}
		o.metadata = make(map[string]string, len(metadata))
		for k, v := range metadata {
			o.metadata[k] = v
		}
	}
}

type secretRecord struct {
	nonce     []byte
	value     []byte
	checksum  [32]byte
	version   int
	createdAt time.Time
	updatedAt time.Time
	expiresAt time.Time
	metadata  map[string]string
}

// SecretMetadata summarises the lifecycle details associated with a secret.
type SecretMetadata struct {
	Key       string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	Metadata  map[string]string
}

// SecretsManager provides envelope encryption, versioning and TTL-based secret
// management for CLI and web integrations. All operations are concurrency safe.
type SecretsManager struct {
	mu        sync.RWMutex
	aead      cipher.AEAD
	masterKey [chacha20poly1305.KeySize]byte
	records   map[string]secretRecord
	now       func() time.Time
}

// NewSecretsManager constructs a manager. When no master key is supplied a
// cryptographically secure random key is generated on startup.
func NewSecretsManager(opts ...SecretsManagerOption) *SecretsManager {
	m := &SecretsManager{
		records: make(map[string]secretRecord),
		now:     time.Now,
	}
	for _, opt := range opts {
		opt(m)
	}
	m.initCipher()
	return m
}

func (s *SecretsManager) initCipher() {
	if s.aead != nil {
		return
	}
	key := s.masterKey
	if key == ([chacha20poly1305.KeySize]byte{}) {
		if _, err := rand.Read(key[:]); err != nil {
			panic(fmt.Errorf("security: failed to initialise master key: %w", err))
		}
		copy(s.masterKey[:], key[:])
	}
	block, err := chacha20poly1305.New(key[:])
	if err != nil {
		panic(fmt.Errorf("security: unable to create cipher: %w", err))
	}
	s.aead = block
}

// Store encrypts and persists the secret value. The operation creates a new
// version and preserves the original creation timestamp when overwriting.
func (s *SecretsManager) Store(key, value string, opts ...StoreOption) error {
	if key == "" {
		return errors.New("security: key required")
	}
	if value == "" {
		return errors.New("security: value required")
	}
	options := storeOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("security: nonce generation failed: %w", err)
	}
	now := s.now().UTC()
	ciphertext := s.aead.Seal(nil, nonce, []byte(value), []byte(key))
	checksum := sha256.Sum256([]byte(value))

	s.mu.Lock()
	defer s.mu.Unlock()
	rec := secretRecord{
		nonce:     nonce,
		value:     ciphertext,
		checksum:  checksum,
		version:   1,
		createdAt: now,
		updatedAt: now,
	}
	if existing, ok := s.records[key]; ok {
		rec.version = existing.version + 1
		rec.createdAt = existing.createdAt
		if options.ttl == 0 {
			rec.expiresAt = existing.expiresAt
		}
		if options.metadata == nil && existing.metadata != nil {
			rec.metadata = make(map[string]string, len(existing.metadata))
			for k, v := range existing.metadata {
				rec.metadata[k] = v
			}
		}
	}
	if options.ttl > 0 {
		rec.expiresAt = now.Add(options.ttl)
	}
	if options.metadata != nil {
		rec.metadata = options.metadata
	}
	s.records[key] = rec
	return nil
}

// Retrieve decrypts the secret value when present and not expired.
func (s *SecretsManager) Retrieve(key string) (string, error) {
	if key == "" {
		return "", errors.New("security: key required")
	}
	s.mu.RLock()
	rec, ok := s.records[key]
	s.mu.RUnlock()
	if !ok {
		return "", ErrSecretNotFound
	}
	if !rec.expiresAt.IsZero() && !s.now().UTC().Before(rec.expiresAt) {
		s.mu.Lock()
		delete(s.records, key)
		s.mu.Unlock()
		return "", ErrSecretExpired
	}
	plain, err := s.aead.Open(nil, rec.nonce, rec.value, []byte(key))
	if err != nil {
		return "", fmt.Errorf("security: decrypt failed: %w", err)
	}
	if sha256.Sum256(plain) != rec.checksum {
		return "", errors.New("security: checksum mismatch")
	}
	return string(plain), nil
}

// Metadata returns lifecycle information associated with the key.
func (s *SecretsManager) Metadata(key string) (SecretMetadata, error) {
	if key == "" {
		return SecretMetadata{}, errors.New("security: key required")
	}
	s.mu.RLock()
	rec, ok := s.records[key]
	s.mu.RUnlock()
	if !ok {
		return SecretMetadata{}, ErrSecretNotFound
	}
	meta := SecretMetadata{
		Key:       key,
		Version:   rec.version,
		CreatedAt: rec.createdAt,
		UpdatedAt: rec.updatedAt,
		ExpiresAt: rec.expiresAt,
	}
	if rec.metadata != nil {
		meta.Metadata = make(map[string]string, len(rec.metadata))
		for k, v := range rec.metadata {
			meta.Metadata[k] = v
		}
	}
	return meta, nil
}

// Delete removes a secret without requiring knowledge of the stored value.
func (s *SecretsManager) Delete(key string) {
	if key == "" {
		return
	}
	s.mu.Lock()
	delete(s.records, key)
	s.mu.Unlock()
}

// RotateMasterKey replaces the encryption key and re-encrypts all stored
// secrets. The operation maintains version counters and updates the updatedAt
// timestamp.
func (s *SecretsManager) RotateMasterKey(newKey []byte) error {
	if len(newKey) != chacha20poly1305.KeySize {
		return errors.New("security: invalid master key length")
	}
	block, err := chacha20poly1305.New(newKey)
	if err != nil {
		return fmt.Errorf("security: unable to create cipher: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, rec := range s.records {
		plain, err := s.aead.Open(nil, rec.nonce, rec.value, []byte(key))
		if err != nil {
			return fmt.Errorf("security: decrypt during rotate failed: %w", err)
		}
		nonce := make([]byte, block.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			return fmt.Errorf("security: nonce generation failed: %w", err)
		}
		rec.value = block.Seal(nil, nonce, plain, []byte(key))
		rec.nonce = nonce
		rec.updatedAt = s.now().UTC()
		s.records[key] = rec
	}
	copy(s.masterKey[:], newKey)
	s.aead = block
	return nil
}

// ListKeys returns the currently stored keys sorted lexicographically.
func (s *SecretsManager) ListKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.records))
	for k := range s.records {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// PurgeExpired removes secrets whose TTL has elapsed.
func (s *SecretsManager) PurgeExpired() {
	now := s.now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, rec := range s.records {
		if rec.expiresAt.IsZero() {
			continue
		}
		if !now.Before(rec.expiresAt) {
			delete(s.records, key)
		}
	}
}
