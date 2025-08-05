package synnergy

import (
	"errors"
	"sync"
)

// DataChannel represents a secure channel with an encryption key.
type DataChannel struct {
	Key      []byte
	Messages [][]byte
	Open     bool
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
	e.channels[id] = &DataChannel{Key: key, Open: true}
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
	e.mu.Lock()
	ch.Messages = append(ch.Messages, cipherText)
	e.mu.Unlock()
	return cipherText, nil
}

// Messages returns encrypted messages for a channel.
func (e *ZeroTrustEngine) Messages(id string) [][]byte {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ch, ok := e.channels[id]
	if !ok {
		return nil
	}
	out := make([][]byte, len(ch.Messages))
	for i, m := range ch.Messages {
		cp := make([]byte, len(m))
		copy(cp, m)
		out[i] = cp
	}
	return out
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
