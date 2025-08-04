package core

import (
	"crypto/sha256"
	"sync"
)

// BiometricService stores biometric hashes for user identities and provides
// enrollment and verification functionality. It is a simple in-memory
// implementation suitable for testing and demonstrates how biometric
// authentication can be integrated into the blockchain.
type BiometricService struct {
	mu   sync.RWMutex
	data map[string][32]byte
}

// NewBiometricService creates a new biometric service instance.
func NewBiometricService() *BiometricService {
	return &BiometricService{data: make(map[string][32]byte)}
}

// Enroll registers biometric data for a user. The biometric data is hashed and
// stored so raw biometric information is never persisted.
func (b *BiometricService) Enroll(userID string, biometric []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	h := sha256.Sum256(biometric)
	b.data[userID] = h
}

// Verify checks the provided biometric data against the stored hash for the
// user. It returns true if the biometric matches what was enrolled.
func (b *BiometricService) Verify(userID string, biometric []byte) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	h, ok := b.data[userID]
	if !ok {
		return false
	}
	return h == sha256.Sum256(biometric)
}
