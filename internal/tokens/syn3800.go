package tokens

import (
	"errors"
	"sync"
)

// SYN3800Token enforces a capped supply with mint and burn operations.
type SYN3800Token struct {
	mu     sync.Mutex
	cap    uint64
	supply uint64
}

// NewSYN3800Token creates a new capped token.
func NewSYN3800Token(cap uint64) *SYN3800Token {
	return &SYN3800Token{cap: cap}
}

// Mint increases the supply if the cap allows it.
func (t *SYN3800Token) Mint(amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.supply+amount > t.cap {
		return errors.New("cap exceeded")
	}
	t.supply += amount
	return nil
}

// Burn reduces the supply.
func (t *SYN3800Token) Burn(amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.supply < amount {
		return errors.New("insufficient supply")
	}
	t.supply -= amount
	return nil
}

// Supply reports the current circulating supply.
func (t *SYN3800Token) Supply() uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.supply
}
