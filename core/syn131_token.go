package core

import "errors"

// SYN131Token represents an intangible asset token.
type SYN131Token struct {
	ID        string
	Name      string
	Symbol    string
	Owner     string
	Valuation uint64
}

// SYN131Registry stores issued SYN131 tokens.
type SYN131Registry struct {
	tokens map[string]*SYN131Token
}

// NewSYN131Registry creates a new registry.
func NewSYN131Registry() *SYN131Registry {
	return &SYN131Registry{tokens: make(map[string]*SYN131Token)}
}

// Create issues a new SYN131 token.
func (r *SYN131Registry) Create(id, name, symbol, owner string, valuation uint64) (*SYN131Token, error) {
	if _, exists := r.tokens[id]; exists {
		return nil, errors.New("token already exists")
	}
	t := &SYN131Token{ID: id, Name: name, Symbol: symbol, Owner: owner, Valuation: valuation}
	r.tokens[id] = t
	return t, nil
}

// UpdateValuation sets a new valuation for the token.
func (r *SYN131Registry) UpdateValuation(id string, val uint64) error {
	tok, ok := r.tokens[id]
	if !ok {
		return errors.New("token not found")
	}
	tok.Valuation = val
	return nil
}

// Get fetches a token by id.
func (r *SYN131Registry) Get(id string) (*SYN131Token, bool) {
	tok, ok := r.tokens[id]
	return tok, ok
}
