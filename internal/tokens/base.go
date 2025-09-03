package tokens

import (
	"errors"
	"sync"
)

// TokenID uniquely identifies a token instance within the registry.
type TokenID uint64

// Token defines basic behaviours for all tokens.
type Token interface {
	ID() TokenID
	Name() string
	Symbol() string
	Decimals() uint8
	TotalSupply() uint64
	BalanceOf(addr string) uint64
	Transfer(from, to string, amount uint64) error
	TransferFrom(owner, spender, to string, amount uint64) error
	Mint(to string, amount uint64) error
	Burn(from string, amount uint64) error
	Approve(owner, spender string, amount uint64) error
	Allowance(owner, spender string) uint64
}

var (
	// ErrInsufficientBalance is returned when an account lacks funds for the
	// requested operation.
	ErrInsufficientBalance = errors.New("insufficient balance")
	// ErrAllowanceExceeded indicates the spender attempted to transfer more
	// than the approved allowance.
	ErrAllowanceExceeded = errors.New("allowance exceeded")
)

// BaseToken implements the Token interface providing basic accounting. It is
// safe for concurrent use by multiple goroutines.
type BaseToken struct {
	mu         sync.RWMutex
	id         TokenID
	name       string
	symbol     string
	decimals   uint8
	balances   map[string]uint64
	supply     uint64
	allowances map[string]map[string]uint64
}

// NewBaseToken creates a new base token instance.
func NewBaseToken(id TokenID, name, symbol string, decimals uint8) *BaseToken {
	return &BaseToken{
		id:         id,
		name:       name,
		symbol:     symbol,
		decimals:   decimals,
		balances:   make(map[string]uint64),
		allowances: make(map[string]map[string]uint64),
	}
}

// ID returns the unique identifier of the token.
func (t *BaseToken) ID() TokenID {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.id
}

// Name returns the human readable token name.
func (t *BaseToken) Name() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.name
}

// Symbol returns the token trading symbol.
func (t *BaseToken) Symbol() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.symbol
}

// Decimals returns the decimal precision for the token.
func (t *BaseToken) Decimals() uint8 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.decimals
}

// TotalSupply returns the current token supply.
func (t *BaseToken) TotalSupply() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.supply
}

// BalanceOf retrieves the balance for the specified address.
func (t *BaseToken) BalanceOf(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.balances[addr]
}

// Transfer moves tokens between addresses.
func (t *BaseToken) Transfer(from, to string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.balances[from] < amount {
		return ErrInsufficientBalance
	}
	t.balances[from] -= amount
	t.balances[to] += amount
	return nil
}

// TransferFrom moves tokens on behalf of an owner using an approved allowance.
func (t *BaseToken) TransferFrom(owner, spender, to string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.allowances[owner][spender] < amount {
		return ErrAllowanceExceeded
	}
	if t.balances[owner] < amount {
		return ErrInsufficientBalance
	}
	t.allowances[owner][spender] -= amount
	t.balances[owner] -= amount
	t.balances[to] += amount
	return nil
}

// Mint creates new tokens for the specified address.
func (t *BaseToken) Mint(to string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.balances[to] += amount
	t.supply += amount
	return nil
}

// Burn removes tokens from the specified address.
func (t *BaseToken) Burn(from string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.balances[from] < amount {
		return ErrInsufficientBalance
	}
	t.balances[from] -= amount
	t.supply -= amount
	return nil
}

// Approve sets the allowance for a spender on behalf of an owner.
func (t *BaseToken) Approve(owner, spender string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.allowances[owner] == nil {
		t.allowances[owner] = make(map[string]uint64)
	}
	t.allowances[owner][spender] = amount
	return nil
}

// Allowance returns the remaining approved amount a spender can use from an owner.
func (t *BaseToken) Allowance(owner, spender string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.allowances[owner][spender]
}
