package security

import "sync"

// KeyManager stores an encryption key with basic rotation capability.
type KeyManager struct {
	mu  sync.RWMutex
	key byte
}

// NewKeyManager creates a KeyManager.
func NewKeyManager(k byte) *KeyManager {
	return &KeyManager{key: k}
}

// Rotate replaces the current key.
func (k *KeyManager) Rotate(newKey byte) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.key = newKey
}

// Key returns the current key.
func (k *KeyManager) Key() byte {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.key
}
