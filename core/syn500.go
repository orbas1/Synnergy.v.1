package core

import "errors"

// ServiceTier defines access tiers for SYN500 utility tokens.
type ServiceTier struct {
	Tier int
	Max  uint64
	Used uint64
}

// SYN500Token defines a simple utility token with usage tracking.
type SYN500Token struct {
	Name     string
	Symbol   string
	Owner    string
	Decimals uint8
	Supply   uint64
	Grants   map[string]*ServiceTier
}

// NewSYN500Token creates a new utility token.
func NewSYN500Token(name, symbol, owner string, decimals uint8, supply uint64) *SYN500Token {
	return &SYN500Token{Name: name, Symbol: symbol, Owner: owner, Decimals: decimals, Supply: supply, Grants: make(map[string]*ServiceTier)}
}

// Grant assigns a service tier to an address.
func (t *SYN500Token) Grant(addr string, tier int, max uint64) {
	t.Grants[addr] = &ServiceTier{Tier: tier, Max: max}
}

// Use records a token usage for an address.
func (t *SYN500Token) Use(addr string) error {
	g, ok := t.Grants[addr]
	if !ok {
		return errors.New("no tier granted")
	}
	if g.Used >= g.Max {
		return errors.New("usage limit reached")
	}
	g.Used++
	return nil
}
