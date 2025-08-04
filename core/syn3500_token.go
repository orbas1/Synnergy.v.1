package core

import (
	"errors"
	"sync"
)

// SYN3500Token represents a currency or stablecoin token.
type SYN3500Token struct {
	mu       sync.RWMutex
	Name     string
	Symbol   string
	Issuer   string
	Rate     float64
	Balances map[string]uint64
}

// NewSYN3500Token creates a new currency token instance.
func NewSYN3500Token(name, symbol, issuer string, rate float64) *SYN3500Token {
	return &SYN3500Token{
		Name:     name,
		Symbol:   symbol,
		Issuer:   issuer,
		Rate:     rate,
		Balances: make(map[string]uint64),
	}
}

// SetRate updates the fiat exchange rate of the token.
func (t *SYN3500Token) SetRate(rate float64) {
	t.mu.Lock()
	t.Rate = rate
	t.mu.Unlock()
}

// Info returns token symbol, issuer and current rate.
func (t *SYN3500Token) Info() (string, string, float64) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Symbol, t.Issuer, t.Rate
}

// Mint creates new tokens for the specified address.
func (t *SYN3500Token) Mint(to string, amt uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Balances[to] += amt
}

// Redeem removes tokens from circulation for the given address.
func (t *SYN3500Token) Redeem(from string, amt uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	bal := t.Balances[from]
	if bal < amt {
		return errors.New("insufficient balance")
	}
	t.Balances[from] = bal - amt
	return nil
}

// BalanceOf returns the balance of the specified address.
func (t *SYN3500Token) BalanceOf(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Balances[addr]
}
