package tokens

import "sync"

// SYN3600Token stores governance weights for addresses.
type SYN3600Token struct {
	mu      sync.RWMutex
	weights map[string]uint64
}

// NewSYN3600Token creates an empty governance weight ledger.
func NewSYN3600Token() *SYN3600Token {
	return &SYN3600Token{weights: make(map[string]uint64)}
}

// SetWeight assigns voting weight to an address.
func (t *SYN3600Token) SetWeight(addr string, w uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.weights[addr] = w
}

// Weight returns the voting weight of an address.
func (t *SYN3600Token) Weight(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.weights[addr]
}
