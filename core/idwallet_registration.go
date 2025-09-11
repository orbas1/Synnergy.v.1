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

var (
	// ErrWalletExists signals a wallet is already registered.
	ErrWalletExists = errors.New("wallet already registered")
	// ErrEmptyAddress is returned when addr is empty.
	ErrEmptyAddress = errors.New("address required")
)

// NewIDRegistry creates a new IDRegistry instance.
func NewIDRegistry() *IDRegistry {
	return &IDRegistry{wallets: make(map[string]string)}
}

// Register adds a wallet with associated info. Returns error if already registered.
func (r *IDRegistry) Register(addr, info string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if addr == "" {
		return ErrEmptyAddress
	}
	if _, exists := r.wallets[addr]; exists {
		return ErrWalletExists
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

// Unregister removes a wallet from the registry. It returns true if the wallet
// was present.
func (r *IDRegistry) Unregister(addr string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.wallets[addr]; ok {
		delete(r.wallets, addr)
		return true
	}
	return false
}
