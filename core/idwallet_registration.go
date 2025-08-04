package core

import (
	"errors"
	"sync"
)

// IDRegistry manages on-chain registration of wallets that hold identity tokens.
type IDRegistry struct {
	mu      sync.RWMutex
	wallets map[string]string // address -> metadata/info
}

// NewIDRegistry creates a new IDRegistry instance.
func NewIDRegistry() *IDRegistry {
	return &IDRegistry{wallets: make(map[string]string)}
}

// Register adds a wallet with associated info. Returns error if already registered.
func (r *IDRegistry) Register(addr, info string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.wallets[addr]; exists {
		return errors.New("wallet already registered")
	}
	r.wallets[addr] = info
	return nil
}

// Info returns registration info for an address if present.
func (r *IDRegistry) Info(addr string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.wallets[addr]
	return info, ok
}
