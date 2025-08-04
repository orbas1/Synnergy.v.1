package core

import (
	"errors"
	"sync"
	"time"
)

// Faucet dispenses test tokens with rate limiting.
type Faucet struct {
	mu           sync.Mutex
	balance      uint64
	amount       uint64
	cooldown     time.Duration
	lastRequests map[string]time.Time
}

// NewFaucet returns a Faucet with the given balance, dispense amount and cooldown.
func NewFaucet(balance, amount uint64, cooldown time.Duration) *Faucet {
	return &Faucet{
		balance:      balance,
		amount:       amount,
		cooldown:     cooldown,
		lastRequests: make(map[string]time.Time),
	}
}

// Request sends faucet funds to addr if cooldown has passed.
func (f *Faucet) Request(addr string) (uint64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	if last, ok := f.lastRequests[addr]; ok && now.Sub(last) < f.cooldown {
		return 0, errors.New("cooldown active")
	}
	if f.balance < f.amount {
		return 0, errors.New("insufficient faucet balance")
	}
	f.balance -= f.amount
	f.lastRequests[addr] = now
	return f.amount, nil
}

// Balance returns the remaining faucet balance.
func (f *Faucet) Balance() uint64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.balance
}

// UpdateConfig sets the dispense amount and cooldown period.
func (f *Faucet) UpdateConfig(amount uint64, cooldown time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.amount = amount
	f.cooldown = cooldown
}
