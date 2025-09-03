package tokens

import "sync"

// SYN2700Token distributes dividends to registered holders based on their share
// of the total supply. It is concurrency-safe for use by multiple goroutines.
type SYN2700Token struct {
	mu      sync.RWMutex
	holders map[string]uint64
	total   uint64
}

// NewSYN2700Token initialises an empty dividend token.
func NewSYN2700Token() *SYN2700Token {
	return &SYN2700Token{holders: make(map[string]uint64)}
}

// AddHolder registers balance for an address. The total supply is increased
// accordingly.
func (t *SYN2700Token) AddHolder(addr string, amount uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.holders[addr] += amount
	t.total += amount
}

// Distribute splits the dividend across all holders proportionally and returns
// the distribution map. The token state is not modified so callers can apply
// transfers separately.
func (t *SYN2700Token) Distribute(dividend uint64) map[string]uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[string]uint64)
	if t.total == 0 {
		return out
	}
	for addr, bal := range t.holders {
		out[addr] = dividend * bal / t.total
	}
	return out
}
