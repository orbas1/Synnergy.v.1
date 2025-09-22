package synnergy

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
)

// RegistryEventType enumerates wallet registry lifecycle notifications.
type RegistryEventType string

const (
	// RegistryEventRegistered indicates a wallet was newly registered.
	RegistryEventRegistered RegistryEventType = "registered"
	// RegistryEventUpdated indicates wallet metadata or info was updated.
	RegistryEventUpdated RegistryEventType = "updated"
	// RegistryEventRemoved indicates a wallet was removed from the registry.
	RegistryEventRemoved RegistryEventType = "removed"
)

// RegistryEvent describes a change to the identity wallet registry.
type RegistryEvent struct {
	Type         RegistryEventType
	Address      string
	RegisteredAt time.Time
	UpdatedAt    time.Time
	Metadata     map[string]string
}

// RegistryEventHandler receives wallet registry events.
type RegistryEventHandler interface {
	HandleRegistryEvent(ctx context.Context, event RegistryEvent) error
}

// RegistryEventHandlerFunc adapts a function into a RegistryEventHandler.
type RegistryEventHandlerFunc func(context.Context, RegistryEvent) error

// HandleRegistryEvent implements RegistryEventHandler.
func (f RegistryEventHandlerFunc) HandleRegistryEvent(ctx context.Context, event RegistryEvent) error {
	return f(ctx, event)
}

// IDRegistry manages on-chain registration of wallets that hold identity tokens.
type IDRegistry struct {
	mu       sync.RWMutex
	wallets  map[string]walletRecord
	sealer   recordSealer
	watchers []RegistryEventHandler
}

type walletRecord struct {
	nonce        []byte
	ciphertext   []byte
	registeredAt time.Time
	updatedAt    time.Time
	metadata     map[string]string
}

type recordSealer interface {
	Seal(plaintext []byte) (nonce, ciphertext []byte, err error)
	Open(nonce, ciphertext []byte) ([]byte, error)
}

type aeadSealer struct {
	aead       cipherAEAD
	randSource io.Reader
}

type cipherAEAD interface {
	NonceSize() int
	Seal(dst, nonce, plaintext, additionalData []byte) []byte
	Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error)
}

func newAEADSealer(key []byte) (*aeadSealer, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("registry: key must be %d bytes", chacha20poly1305.KeySize)
	}
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	return &aeadSealer{aead: aead, randSource: rand.Reader}, nil
}

func (s *aeadSealer) Seal(plaintext []byte) (nonce, ciphertext []byte, err error) {
	nonce = make([]byte, s.aead.NonceSize())
	if _, err = io.ReadFull(s.randSource, nonce); err != nil {
		return nil, nil, err
	}
	sealed := s.aead.Seal(nil, nonce, plaintext, nil)
	return nonce, sealed, nil
}

func (s *aeadSealer) Open(nonce, ciphertext []byte) ([]byte, error) {
	return s.aead.Open(nil, nonce, ciphertext, nil)
}

// IDRegistryOption configures the wallet registry.
type IDRegistryOption func(*IDRegistry) error

// WithRegistryEncryptionKey configures the registry to use the provided encryption key.
func WithRegistryEncryptionKey(key []byte) IDRegistryOption {
	return func(r *IDRegistry) error {
		sealer, err := newAEADSealer(key)
		if err != nil {
			return err
		}
		r.sealer = sealer
		return nil
	}
}

// WithRegistryEventHandler registers an event handler for registry lifecycle events.
func WithRegistryEventHandler(handler RegistryEventHandler) IDRegistryOption {
	return func(r *IDRegistry) error {
		if handler != nil {
			r.watchers = append(r.watchers, handler)
		}
		return nil
	}
}

var (
	// ErrWalletExists signals a wallet is already registered.
	ErrWalletExists = errors.New("registry: wallet already registered")
	// ErrEmptyAddress is returned when addr is empty.
	ErrEmptyAddress = errors.New("registry: address required")
	// ErrWalletNotFound is returned when a wallet lookup fails.
	ErrWalletNotFound = errors.New("registry: wallet not registered")
)

// NewIDRegistry creates a new IDRegistry instance.
func NewIDRegistry(opts ...IDRegistryOption) *IDRegistry {
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(key); err != nil {
		// Fallback to deterministic but time-varying key if entropy source unavailable.
		copy(key, time.Now().UTC().AppendFormat(nil, time.RFC3339Nano))
	}
	sealer, err := newAEADSealer(key)
	if err != nil {
		panic(fmt.Errorf("registry: failed to initialise encryption: %w", err))
	}
	reg := &IDRegistry{
		wallets: make(map[string]walletRecord),
		sealer:  sealer,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(reg); err != nil {
			panic(fmt.Errorf("registry: option error: %w", err))
		}
	}
	return reg
}

// Register adds a wallet with associated info. Returns error if already registered.
func (r *IDRegistry) Register(addr, info string) error {
	return r.RegisterContext(context.Background(), addr, info, nil)
}

// RegisterContext registers a wallet with metadata using the provided context.
func (r *IDRegistry) RegisterContext(ctx context.Context, addr, info string, metadata map[string]string) error {
	if addr == "" {
		return ErrEmptyAddress
	}
	nonce, cipher, err := r.sealer.Seal([]byte(info))
	if err != nil {
		return fmt.Errorf("registry: encrypt info: %w", err)
	}
	now := time.Now().UTC()
	record := walletRecord{
		nonce:        nonce,
		ciphertext:   cipher,
		registeredAt: now,
		updatedAt:    now,
		metadata:     cloneMetadata(metadata),
	}

	r.mu.Lock()
	if _, exists := r.wallets[addr]; exists {
		r.mu.Unlock()
		return ErrWalletExists
	}
	r.wallets[addr] = record
	r.mu.Unlock()

	return r.fireEvent(ctx, RegistryEvent{Type: RegistryEventRegistered, Address: addr, RegisteredAt: now, UpdatedAt: now, Metadata: cloneMetadata(metadata)})
}

// Update replaces the stored wallet info and metadata.
func (r *IDRegistry) Update(ctx context.Context, addr, info string, metadata map[string]string) error {
	if addr == "" {
		return ErrEmptyAddress
	}
	nonce, cipher, err := r.sealer.Seal([]byte(info))
	if err != nil {
		return fmt.Errorf("registry: encrypt info: %w", err)
	}
	now := time.Now().UTC()

	r.mu.Lock()
	record, ok := r.wallets[addr]
	if !ok {
		r.mu.Unlock()
		return ErrWalletNotFound
	}
	record.nonce = nonce
	record.ciphertext = cipher
	record.updatedAt = now
	if metadata != nil {
		record.metadata = cloneMetadata(metadata)
	}
	r.wallets[addr] = record
	r.mu.Unlock()

	return r.fireEvent(ctx, RegistryEvent{Type: RegistryEventUpdated, Address: addr, RegisteredAt: record.registeredAt, UpdatedAt: now, Metadata: cloneMetadata(record.metadata)})
}

// Info returns registration info for an address if present.
func (r *IDRegistry) Info(addr string) (string, bool) {
	info, _, ok := r.InfoWithMetadata(addr)
	return info, ok
}

// InfoWithMetadata returns the decrypted wallet info and metadata.
func (r *IDRegistry) InfoWithMetadata(addr string) (string, map[string]string, bool) {
	r.mu.RLock()
	record, ok := r.wallets[addr]
	r.mu.RUnlock()
	if !ok {
		return "", nil, false
	}
	plaintext, err := r.sealer.Open(record.nonce, record.ciphertext)
	if err != nil {
		return "", nil, false
	}
	return string(plaintext), cloneMetadata(record.metadata), true
}

// IsRegistered reports whether the address has been registered.
func (r *IDRegistry) IsRegistered(addr string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.wallets[addr]
	return ok
}

// Unregister removes a wallet from the registry. It returns true if the wallet was present.
func (r *IDRegistry) Unregister(ctx context.Context, addr string) bool {
	r.mu.Lock()
	record, ok := r.wallets[addr]
	if ok {
		delete(r.wallets, addr)
	}
	r.mu.Unlock()
	if !ok {
		return false
	}
	_ = r.fireEvent(ctx, RegistryEvent{Type: RegistryEventRemoved, Address: addr, RegisteredAt: record.registeredAt, UpdatedAt: time.Now().UTC(), Metadata: cloneMetadata(record.metadata)})
	return true
}

// RotateKey re-encrypts all stored wallet data with the provided key.
func (r *IDRegistry) RotateKey(key []byte) error {
	sealer, err := newAEADSealer(key)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for addr, record := range r.wallets {
		plaintext, err := r.sealer.Open(record.nonce, record.ciphertext)
		if err != nil {
			return fmt.Errorf("registry: decrypt during rotate: %w", err)
		}
		nonce, cipher, err := sealer.Seal(plaintext)
		if err != nil {
			return fmt.Errorf("registry: encrypt during rotate: %w", err)
		}
		record.nonce = nonce
		record.ciphertext = cipher
		r.wallets[addr] = record
	}
	r.sealer = sealer
	return nil
}

// Addresses returns a sorted snapshot of registered wallet addresses.
func (r *IDRegistry) Addresses() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.wallets))
	for addr := range r.wallets {
		out = append(out, addr)
	}
	sort.Strings(out)
	return out
}

func (r *IDRegistry) fireEvent(ctx context.Context, event RegistryEvent) error {
	r.mu.RLock()
	watchers := append([]RegistryEventHandler(nil), r.watchers...)
	r.mu.RUnlock()
	var errs []error
	for _, watcher := range watchers {
		if watcher == nil {
			continue
		}
		if err := watcher.HandleRegistryEvent(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
