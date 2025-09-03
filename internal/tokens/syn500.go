package tokens

import (
	"sync"
	"time"
)

// SYN500Token tracks loyalty points with expiration.
type SYN500Token struct {
	mu       sync.RWMutex
	accounts map[string]*loyalty
}

type loyalty struct {
	points uint64
	expiry time.Time
}

// NewSYN500Token initialises the loyalty registry.
func NewSYN500Token() *SYN500Token {
	return &SYN500Token{accounts: make(map[string]*loyalty)}
}

// Mint grants points to an address until the given expiry.
func (t *SYN500Token) Mint(addr string, amount uint64, exp time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.accounts[addr] = &loyalty{points: amount, expiry: exp}
}

// Redeem returns the points for an address if not expired and zeroes the balance.
func (t *SYN500Token) Redeem(addr string, now time.Time) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	acc, ok := t.accounts[addr]
	if !ok || now.After(acc.expiry) {
		return 0
	}
	pts := acc.points
	acc.points = 0
	return pts
}
