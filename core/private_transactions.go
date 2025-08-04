package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"sync"
)

// Encrypt encrypts plaintext using AES-GCM with the provided key.
// The returned slice contains nonce||ciphertext.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherText := gcm.Seal(nonce, nonce, plaintext, nil)
	return cipherText, nil
}

// Decrypt decrypts data produced by Encrypt.
func Decrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, io.ErrUnexpectedEOF
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, cipherText, nil)
}

// PrivateTransaction holds an encrypted payload.
type PrivateTransaction struct {
	Payload []byte
}

// PrivateTxManager manages private transactions.
type PrivateTxManager struct {
	mu  sync.Mutex
	txs []PrivateTransaction
}

// NewPrivateTxManager creates a new PrivateTxManager.
func NewPrivateTxManager() *PrivateTxManager {
	return &PrivateTxManager{}
}

// Send adds a private transaction to the internal pool.
func (m *PrivateTxManager) Send(tx PrivateTransaction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.txs = append(m.txs, tx)
}

// List returns a copy of stored private transactions.
func (m *PrivateTxManager) List() []PrivateTransaction {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]PrivateTransaction, len(m.txs))
	copy(out, m.txs)
	return out
}
