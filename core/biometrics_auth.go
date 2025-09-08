package core

import (
	"crypto/ed25519"
	"crypto/sha256"
	"errors"
	"sync"
)

// BiometricsAuth manages hashed biometric templates for addresses.
type BiometricsAuth struct {
	mu        sync.RWMutex
	templates map[string]biometricTemplate
}

type biometricTemplate struct {
	hash [32]byte
	pub  ed25519.PublicKey
}

var (
	ErrAddressRequired = errors.New("address required")
	ErrAlreadyEnrolled = errors.New("biometric already enrolled")
)

// NewBiometricsAuth creates a new biometrics authentication manager.
func NewBiometricsAuth() *BiometricsAuth {
	return &BiometricsAuth{templates: make(map[string]biometricTemplate)}
}

// Enroll stores a hashed biometric template for the given address.
// Returns an error if input is invalid or the address already exists.
func (b *BiometricsAuth) Enroll(addr string, biometric []byte, pub ed25519.PublicKey) error {
	if addr == "" {
		return ErrAddressRequired
	}
	if len(biometric) == 0 {
		return ErrInvalidBiometric
	}
	if len(pub) != ed25519.PublicKeySize {
		return ErrInvalidPublicKey
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.templates[addr]; exists {
		return ErrAlreadyEnrolled
	}
	b.templates[addr] = biometricTemplate{hash: sha256.Sum256(biometric), pub: pub}
	return nil
}

// Verify compares the provided biometric data with the stored template for the address.
func (b *BiometricsAuth) Verify(addr string, biometric []byte, sig []byte) bool {
	b.mu.RLock()
	tmpl, ok := b.templates[addr]
	b.mu.RUnlock()
	if !ok {
		return false
	}
	h := sha256.Sum256(biometric)
	if h != tmpl.hash || len(tmpl.pub) != ed25519.PublicKeySize || len(sig) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(tmpl.pub, h[:], sig)
}

// Remove deletes the biometric template for the given address.
func (b *BiometricsAuth) Remove(addr string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.templates, addr)
}

// Enrolled returns true if biometric data has been enrolled for the address.
func (b *BiometricsAuth) Enrolled(addr string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.templates[addr]
	return ok
}

// List returns a snapshot of all addresses that have enrolled biometrics.
// The returned slice is a copy and safe for the caller to modify.
func (b *BiometricsAuth) List() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	addrs := make([]string, 0, len(b.templates))
	for addr := range b.templates {
		addrs = append(addrs, addr)
	}
	return addrs
}
