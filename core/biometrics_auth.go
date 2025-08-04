package core

import (
	"crypto/sha256"
	"sync"
)

// BiometricsAuth manages hashed biometric templates for addresses.
type BiometricsAuth struct {
	mu        sync.RWMutex
	templates map[string][32]byte
}

// NewBiometricsAuth creates a new biometrics authentication manager.
func NewBiometricsAuth() *BiometricsAuth {
	return &BiometricsAuth{templates: make(map[string][32]byte)}
}

// Enroll stores a hashed biometric template for the given address.
func (b *BiometricsAuth) Enroll(addr string, biometric []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.templates[addr] = sha256.Sum256(biometric)
}

// Verify compares the provided biometric data with the stored template for the address.
func (b *BiometricsAuth) Verify(addr string, biometric []byte) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	h, ok := b.templates[addr]
	if !ok {
		return false
	}
	return h == sha256.Sum256(biometric)
}

// Remove deletes the biometric template for the given address.
func (b *BiometricsAuth) Remove(addr string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.templates, addr)
}
