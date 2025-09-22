package tokens

import (
	"errors"
	"sync"
)

var (
	ErrCapExceeded        = errors.New("tokens: cap exceeded")
	ErrInsufficientSupply = errors.New("tokens: insufficient supply")
)

// SYN3800Token enforces a capped supply with mint and burn operations.
type SYN3800Token struct {
	mu     sync.RWMutex
	cap    uint64
	supply uint64
}

// NewSYN3800Token creates a new capped token.
func NewSYN3800Token(cap uint64) *SYN3800Token {
	return &SYN3800Token{cap: cap}
}

// Mint increases the supply if the cap allows it.
func (t *SYN3800Token) Mint(amount uint64) error {
	if amount == 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.supply+amount > t.cap {
		return ErrCapExceeded
	}
	t.supply += amount
	return nil
}

// Burn reduces the supply.
func (t *SYN3800Token) Burn(amount uint64) error {
	if amount == 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.supply < amount {
		return ErrInsufficientSupply
	}
	t.supply -= amount
	return nil
}

// Supply reports the current circulating supply.
func (t *SYN3800Token) Supply() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.supply
}

// Cap returns the maximum supply.
func (t *SYN3800Token) Cap() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.cap
}

// Remaining returns the capacity left before hitting the cap.
func (t *SYN3800Token) Remaining() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.cap <= t.supply {
		return 0
	}
	return t.cap - t.supply
}

// SetCap updates the maximum supply if the new value is not below the current supply.
func (t *SYN3800Token) SetCap(newCap uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if newCap < t.supply {
		return ErrCapExceeded
	}
	t.cap = newCap
	return nil
}
