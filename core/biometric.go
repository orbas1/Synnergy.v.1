package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"sync"
)

// BiometricService stores biometric hashes for user identities and provides
// enrollment and verification functionality. It is a simple in-memory
// implementation suitable for testing and demonstrates how biometric
// authentication can be integrated into the blockchain.
type BiometricService struct {
	mu   sync.RWMutex
	data map[string]biometricRecord
}

type biometricRecord struct {
	hash [32]byte
	pub  *ecdsa.PublicKey
}

// NewBiometricService creates a new biometric service instance.
func NewBiometricService() *BiometricService {
	return &BiometricService{data: make(map[string]biometricRecord)}
}

// Enroll registers biometric data for a user. The biometric data is hashed and
// stored so raw biometric information is never persisted.
func (b *BiometricService) Enroll(userID string, biometric []byte, pub *ecdsa.PublicKey) {
	b.mu.Lock()
	defer b.mu.Unlock()
	h := sha256.Sum256(biometric)
	b.data[userID] = biometricRecord{hash: h, pub: pub}
}

// Verify checks the provided biometric data against the stored hash for the
// user. It returns true if the biometric matches what was enrolled.
func (b *BiometricService) Verify(userID string, biometric []byte, sig []byte) bool {
	b.mu.RLock()
	rec, ok := b.data[userID]
	b.mu.RUnlock()
	if !ok {
		return false
	}
	h := sha256.Sum256(biometric)
	if h != rec.hash || rec.pub == nil {
		return false
	}
	return ecdsa.VerifyASN1(rec.pub, h[:], sig)
}
