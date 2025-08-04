package synnergy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"sync"
)

// ContentNode handles storage of large encrypted content pieces.
// Data is encrypted using AES-GCM before being persisted in memory.
type ContentNode struct {
	mu       sync.RWMutex
	key      []byte
	contents map[string][]byte
	metas    map[string]ContentMeta
}

// NewContentNode creates a ContentNode using the provided 32-byte AES key.
func NewContentNode(key []byte) (*ContentNode, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes")
	}
	return &ContentNode{
		key:      append([]byte(nil), key...),
		contents: make(map[string][]byte),
		metas:    make(map[string]ContentMeta),
	}, nil
}

// StoreContent encrypts data, stores it and returns associated metadata.
func (n *ContentNode) StoreContent(name string, data []byte) (ContentMeta, error) {
	block, err := aes.NewCipher(n.key)
	if err != nil {
		return ContentMeta{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ContentMeta{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return ContentMeta{}, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	hash := sha256.Sum256(data)
	id := hex.EncodeToString(hash[:])
	meta := NewContentMeta(id, name, int64(len(data)), id)

	n.mu.Lock()
	n.contents[id] = ciphertext
	n.metas[id] = meta
	n.mu.Unlock()

	return meta, nil
}

// RetrieveContent decrypts and returns the plaintext for the given id.
func (n *ContentNode) RetrieveContent(id string) ([]byte, bool, error) {
	n.mu.RLock()
	enc, ok := n.contents[id]
	n.mu.RUnlock()
	if !ok {
		return nil, false, nil
	}

	block, err := aes.NewCipher(n.key)
	if err != nil {
		return nil, false, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, false, err
	}
	nonceSize := gcm.NonceSize()
	if len(enc) < nonceSize {
		return nil, false, errors.New("ciphertext too short")
	}
	nonce, cipherText := enc[:nonceSize], enc[nonceSize:]
	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, false, err
	}
	return plain, true, nil
}

// Meta returns the ContentMeta for the provided id.
func (n *ContentNode) Meta(id string) (ContentMeta, bool) {
	n.mu.RLock()
	meta, ok := n.metas[id]
	n.mu.RUnlock()
	return meta, ok
}

// DeleteContent removes the stored ciphertext and metadata for the given id.
func (n *ContentNode) DeleteContent(id string) {
	n.mu.Lock()
	delete(n.contents, id)
	delete(n.metas, id)
	n.mu.Unlock()
}
