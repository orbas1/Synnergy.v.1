package tokens

import (
	"errors"
	"sync"
)

var ErrAccountFrozen = errors.New("tokens: account is frozen")

// SYN3500Token represents a currency or stablecoin token.
type SYN3500Token struct {
	mu         sync.RWMutex
	Name       string
	Symbol     string
	Issuer     string
	Rate       float64
	Balances   map[string]uint64
	allowances map[string]map[string]uint64
	frozen     map[string]bool
	total      uint64
}

// NewSYN3500Token creates a new currency token instance.
func NewSYN3500Token(name, symbol, issuer string, rate float64) *SYN3500Token {
	return &SYN3500Token{
		Name:       name,
		Symbol:     symbol,
		Issuer:     issuer,
		Rate:       rate,
		Balances:   make(map[string]uint64),
		allowances: make(map[string]map[string]uint64),
		frozen:     make(map[string]bool),
	}
}

// SetRate updates the fiat exchange rate of the token.
func (t *SYN3500Token) SetRate(rate float64) error {
	if rate <= 0 {
		return ErrInvalidForexRate
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Rate = rate
	return nil
}

// Info returns token symbol, issuer and current rate.
func (t *SYN3500Token) Info() (string, string, float64) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Symbol, t.Issuer, t.Rate
}

// Mint creates new tokens for the specified address.
func (t *SYN3500Token) Mint(to string, amt uint64) error {
	if amt == 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.frozen[to] {
		return ErrAccountFrozen
	}
	t.Balances[to] += amt
	t.total += amt
	return nil
}

// Redeem removes tokens from circulation for the given address.
func (t *SYN3500Token) Redeem(from string, amt uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.frozen[from] {
		return ErrAccountFrozen
	}
	bal := t.Balances[from]
	if bal < amt {
		return ErrInsufficientBalance
	}
	t.Balances[from] = bal - amt
	t.total -= amt
	return nil
}

// Transfer moves tokens between accounts respecting freeze status.
func (t *SYN3500Token) Transfer(from, to string, amt uint64) error {
	if amt == 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.frozen[from] || t.frozen[to] {
		return ErrAccountFrozen
	}
	bal := t.Balances[from]
	if bal < amt {
		return ErrInsufficientBalance
	}
	t.Balances[from] = bal - amt
	t.Balances[to] += amt
	return nil
}

// Approve allows the spender to transfer up to amt from the owner's account.
func (t *SYN3500Token) Approve(owner, spender string, amt uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.frozen[owner] {
		return ErrAccountFrozen
	}
	if _, ok := t.allowances[owner]; !ok {
		t.allowances[owner] = make(map[string]uint64)
	}
	t.allowances[owner][spender] = amt
	return nil
}

// Allowance returns the approved spending limit for a spender.
func (t *SYN3500Token) Allowance(owner, spender string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.allowances[owner][spender]
}

// TransferFrom spends the allowance and moves funds to the recipient.
func (t *SYN3500Token) TransferFrom(owner, spender, to string, amt uint64) error {
	if amt == 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.frozen[owner] || t.frozen[to] {
		return ErrAccountFrozen
	}
	allowed := t.allowances[owner][spender]
	if allowed < amt {
		return ErrInsufficientBalance
	}
	bal := t.Balances[owner]
	if bal < amt {
		return ErrInsufficientBalance
	}
	t.allowances[owner][spender] = allowed - amt
	t.Balances[owner] = bal - amt
	t.Balances[to] += amt
	return nil
}

// BalanceOf returns the balance of the specified address.
func (t *SYN3500Token) BalanceOf(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Balances[addr]
}

// TotalSupply returns the total supply.
func (t *SYN3500Token) TotalSupply() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.total
}

// FreezeAccount prevents transfers to or from the address.
func (t *SYN3500Token) FreezeAccount(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.frozen[addr] = true
}

// UnfreezeAccount removes the freeze restriction.
func (t *SYN3500Token) UnfreezeAccount(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.frozen, addr)
}

// IsFrozen reports whether an account is currently frozen.
func (t *SYN3500Token) IsFrozen(addr string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.frozen[addr]
}
