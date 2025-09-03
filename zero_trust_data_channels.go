package synnergy

import (
	"crypto/ed25519"
	"errors"
	"sync"
)

// DataChannel represents a secure channel with an encryption key.
type DataChannel struct {
	Key      []byte
	PrivKey  ed25519.PrivateKey
	PubKey   ed25519.PublicKey
	Messages []SignedMessage
	Open     bool
}

// SignedMessage bundles an encrypted payload with its signature.
type SignedMessage struct {
	Cipher    []byte
	Signature []byte
}

// ZeroTrustEngine manages encrypted data channels backed by ledger escrows.
type ZeroTrustEngine struct {
	mu       sync.RWMutex
	channels map[string]*DataChannel
}

// NewZeroTrustEngine creates a new ZeroTrustEngine instance.
func NewZeroTrustEngine() *ZeroTrustEngine {
	return &ZeroTrustEngine{channels: make(map[string]*DataChannel)}
}

// OpenChannel initialises a new channel with the provided key.
func (e *ZeroTrustEngine) OpenChannel(id string, key []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.channels[id]; exists {
		return errors.New("channel already exists")
	}
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	e.channels[id] = &DataChannel{Key: key, PrivKey: priv, PubKey: pub, Open: true}
	return nil
}

// Send encrypts and stores a payload on the channel.
func (e *ZeroTrustEngine) Send(id string, payload []byte) ([]byte, error) {
	e.mu.RLock()
	ch, ok := e.channels[id]
	e.mu.RUnlock()
	if !ok || !ch.Open {
		return nil, errors.New("channel not open")
	}
	cipherText, err := Encrypt(ch.Key, payload)
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(ch.PrivKey, cipherText)
	e.mu.Lock()
	ch.Messages = append(ch.Messages, SignedMessage{Cipher: cipherText, Signature: sig})
	e.mu.Unlock()
	return cipherText, nil
}

// Messages returns encrypted messages for a channel.
func (e *ZeroTrustEngine) Messages(id string) []SignedMessage {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ch, ok := e.channels[id]
	if !ok {
		return nil
	}
	out := make([]SignedMessage, len(ch.Messages))
	copy(out, ch.Messages)
	return out
}

// Receive verifies and decrypts a stored message by index.
func (e *ZeroTrustEngine) Receive(id string, index int) ([]byte, error) {
	e.mu.RLock()
	ch, ok := e.channels[id]
	e.mu.RUnlock()
	if !ok {
		return nil, errors.New("channel not found")
	}
	if index < 0 || index >= len(ch.Messages) {
		return nil, errors.New("message index out of range")
	}
	msg := ch.Messages[index]
	if !ed25519.Verify(ch.PubKey, msg.Cipher, msg.Signature) {
		return nil, errors.New("signature verification failed")
	}
	pt, err := Decrypt(ch.Key, msg.Cipher)
	if err != nil {
		return nil, err
	}
	return pt, nil
}

// CloseChannel closes the channel and prevents further messages.
func (e *ZeroTrustEngine) CloseChannel(id string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.channels[id]
	if !ok {
		return errors.New("channel not found")
	}
	ch.Open = false
	return nil
}
