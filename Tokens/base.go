package tokens

import "fmt"

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
	Mint(to string, amount uint64) error
	Burn(from string, amount uint64) error
}

// BaseToken implements the Token interface providing basic accounting.
type BaseToken struct {
	id       TokenID
	name     string
	symbol   string
	decimals uint8
	balances map[string]uint64
	supply   uint64
}

// NewBaseToken creates a new base token instance.
func NewBaseToken(id TokenID, name, symbol string, decimals uint8) *BaseToken {
	return &BaseToken{
		id:       id,
		name:     name,
		symbol:   symbol,
		decimals: decimals,
		balances: make(map[string]uint64),
	}
}

// ID returns the unique identifier of the token.
func (t *BaseToken) ID() TokenID { return t.id }

// Name returns the human readable token name.
func (t *BaseToken) Name() string { return t.name }

// Symbol returns the token trading symbol.
func (t *BaseToken) Symbol() string { return t.symbol }

// Decimals returns the decimal precision for the token.
func (t *BaseToken) Decimals() uint8 { return t.decimals }

// TotalSupply returns the current token supply.
func (t *BaseToken) TotalSupply() uint64 { return t.supply }

// BalanceOf retrieves the balance for the specified address.
func (t *BaseToken) BalanceOf(addr string) uint64 {
	return t.balances[addr]
}

// Transfer moves tokens between addresses.
func (t *BaseToken) Transfer(from, to string, amount uint64) error {
	if t.balances[from] < amount {
		return fmt.Errorf("insufficient balance")
	}
	t.balances[from] -= amount
	t.balances[to] += amount
	return nil
}

// Mint creates new tokens for the specified address.
func (t *BaseToken) Mint(to string, amount uint64) error {
	t.balances[to] += amount
	t.supply += amount
	return nil
}

// Burn removes tokens from the specified address.
func (t *BaseToken) Burn(from string, amount uint64) error {
	if t.balances[from] < amount {
		return fmt.Errorf("insufficient balance")
	}
	t.balances[from] -= amount
	t.supply -= amount
	return nil
}
