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
	dispensed    map[string]uint64
	dailyLimit   uint64
	lastReset    time.Time
}

// NewFaucet returns a Faucet with the given balance, dispense amount and cooldown.
func NewFaucet(balance, amount uint64, cooldown time.Duration) *Faucet {
	return &Faucet{
		balance:      balance,
		amount:       amount,
		cooldown:     cooldown,
		lastRequests: make(map[string]time.Time),
		dispensed:    make(map[string]uint64),
		lastReset:    time.Now(),
	}
}

var ErrFaucetDailyLimit = errors.New("faucet daily limit reached")

// SetDailyLimit configures the per-address daily dispense limit. Zero disables the limit.
func (f *Faucet) SetDailyLimit(limit uint64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.dailyLimit = limit
}

func (f *Faucet) resetLocked(now time.Time) {
	if now.Sub(f.lastReset) >= 24*time.Hour {
		f.dispensed = make(map[string]uint64)
		f.lastReset = now
	}
}

// Request sends faucet funds to addr if cooldown has passed.
func (f *Faucet) Request(addr string) (uint64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	f.resetLocked(now)
	if last, ok := f.lastRequests[addr]; ok && now.Sub(last) < f.cooldown {
		return 0, errors.New("cooldown active")
	}
	if f.dailyLimit > 0 && f.dispensed[addr]+f.amount > f.dailyLimit {
		return 0, ErrFaucetDailyLimit
	}
	if f.balance < f.amount {
		return 0, errors.New("insufficient faucet balance")
	}
	f.balance -= f.amount
	f.lastRequests[addr] = now
	f.dispensed[addr] += f.amount
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
