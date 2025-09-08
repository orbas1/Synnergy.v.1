package core

import (
	"crypto/ed25519"
	"crypto/sha256"
	"errors"
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
	pub  ed25519.PublicKey
}

var (
	ErrEmptyUserID       = errors.New("userID required")
	ErrInvalidBiometric  = errors.New("biometric data required")
	ErrInvalidPublicKey  = errors.New("invalid public key")
	ErrAlreadyRegistered = errors.New("biometric already enrolled")
)

// NewBiometricService creates a new biometric service instance.
func NewBiometricService() *BiometricService {
	return &BiometricService{data: make(map[string]biometricRecord)}
}

// Enroll registers biometric data for a user. The biometric data is hashed and
// stored so raw biometric information is never persisted.
// Returns an error if input validation fails or the user is already enrolled.
func (b *BiometricService) Enroll(userID string, biometric []byte, pub ed25519.PublicKey) error {
	if userID == "" {
		return ErrEmptyUserID
	}
	if len(biometric) == 0 {
		return ErrInvalidBiometric
	}
	if len(pub) != ed25519.PublicKeySize {
		return ErrInvalidPublicKey
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.data[userID]; exists {
		return ErrAlreadyRegistered
	}
	h := sha256.Sum256(biometric)
	b.data[userID] = biometricRecord{hash: h, pub: pub}
	return nil
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
	if h != rec.hash || len(rec.pub) != ed25519.PublicKeySize || len(sig) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(rec.pub, h[:], sig)
}

// Remove deletes stored biometric data for the user.
func (b *BiometricService) Remove(userID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.data, userID)
}

// Enrolled returns true if biometric data has been registered for the user.
func (b *BiometricService) Enrolled(userID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.data[userID]
	return ok
}

// List returns a copy of all user IDs with enrolled biometrics.
func (b *BiometricService) List() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	ids := make([]string, 0, len(b.data))
	for id := range b.data {
		ids = append(ids, id)
	}
	return ids
}
