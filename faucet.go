package synnergy

import (
	"errors"
	"sync"
	"time"
)

// Faucet dispenses small amounts of test funds to users while enforcing a
// cooldown period between requests.
type Faucet struct {
	mu       sync.Mutex
	balance  uint64
	amount   uint64
	cooldown time.Duration
	lastReq  map[string]time.Time
}

// NewFaucet creates a faucet with an initial balance, dispense amount and
// cooldown period.
func NewFaucet(balance, amount uint64, cooldown time.Duration) *Faucet {
	return &Faucet{balance: balance, amount: amount, cooldown: cooldown, lastReq: make(map[string]time.Time)}
}

// Request dispenses funds to the given address if enough balance remains and the
// cooldown period has elapsed.  The amount granted is returned.
func (f *Faucet) Request(addr string, now time.Time) (uint64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.balance < f.amount {
		return 0, errors.New("faucet empty")
	}
	if t, ok := f.lastReq[addr]; ok && now.Sub(t) < f.cooldown {
		return 0, errors.New("cooldown period not met")
	}
	f.balance -= f.amount
	f.lastReq[addr] = now
	return f.amount, nil
}

// Balance returns the remaining faucet balance.
func (f *Faucet) Balance() uint64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.balance
}

// Configure updates the dispense amount and cooldown duration.
func (f *Faucet) Configure(amount uint64, cooldown time.Duration) {
	f.mu.Lock()
	f.amount = amount
	f.cooldown = cooldown
	f.mu.Unlock()
}
