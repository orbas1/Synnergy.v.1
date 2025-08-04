package synnergy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"sync"
)

// SecureStorage encrypts and stores model bytes.
type SecureStorage struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// NewSecureStorage constructs a new SecureStorage.
func NewSecureStorage() *SecureStorage {
	return &SecureStorage{store: make(map[string][]byte)}
}

// Store encrypts data with key and stores it under hash.
func (s *SecureStorage) Store(hash string, data, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	sealed := gcm.Seal(nonce, nonce, data, nil)
	s.mu.Lock()
	s.store[hash] = sealed
	s.mu.Unlock()
	return nil
}

// Retrieve decrypts and returns stored data.
func (s *SecureStorage) Retrieve(hash string, key []byte) ([]byte, error) {
	s.mu.RLock()
	sealed, ok := s.store[hash]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("model not found")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(sealed) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := sealed[:nonceSize], sealed[nonceSize:]
	data, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}
