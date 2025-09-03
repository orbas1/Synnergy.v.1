package tokens

import (
	"sync"
	"time"
)

// SYN3900Token manages vesting grants that unlock after a release time.
type SYN3900Token struct {
	mu     sync.RWMutex
	grants map[string]*vestingGrant
}

type vestingGrant struct {
	amount   uint64
	release  time.Time
	released bool
}

// NewSYN3900Token creates an empty vesting registry.
func NewSYN3900Token() *SYN3900Token {
	return &SYN3900Token{grants: make(map[string]*vestingGrant)}
}

// Grant assigns a vesting schedule to an address.
func (t *SYN3900Token) Grant(addr string, amount uint64, release time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.grants[addr] = &vestingGrant{amount: amount, release: release}
}

// Release returns the vested amount if the release time has passed. Subsequent
// calls return zero to prevent double spending.
func (t *SYN3900Token) Release(addr string, now time.Time) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	g, ok := t.grants[addr]
	if !ok || g.released || now.Before(g.release) {
		return 0
	}
	g.released = true
	return g.amount
}
