package tokens

import "sync"

// Registry maintains a lookup of token instances by their identifiers. It is
// safe for concurrent use by multiple goroutines.
type Registry struct {
	mu     sync.RWMutex
	tokens map[TokenID]Token
	next   TokenID
}

// NewRegistry creates an empty token registry.
func NewRegistry() *Registry {
	return &Registry{tokens: make(map[TokenID]Token)}
}

// NextID returns a unique identifier for a new token.
func (r *Registry) NextID() TokenID {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	return r.next
}

// Register adds the token to the registry using its ID.
func (r *Registry) Register(t Token) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[t.ID()] = t
}

// Get retrieves a token by ID if present.
func (r *Registry) Get(id TokenID) (Token, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[id]
	return t, ok
}

// GetBySymbol retrieves a token by its symbol.
func (r *Registry) GetBySymbol(symbol string) (Token, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tokens {
		if t.Symbol() == symbol {
			return t, true
		}
	}
	return nil, false
}

// TokenInfo summarises metadata for a registered token.
type TokenInfo struct {
	ID          TokenID
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply uint64
}

// Info returns metadata for a token by ID.
func (r *Registry) Info(id TokenID) (TokenInfo, bool) {
	r.mu.RLock()
	t, ok := r.tokens[id]
	if !ok {
		r.mu.RUnlock()
		return TokenInfo{}, false
	}
	info := TokenInfo{
		ID:          t.ID(),
		Name:        t.Name(),
		Symbol:      t.Symbol(),
		Decimals:    t.Decimals(),
		TotalSupply: t.TotalSupply(),
	}
	r.mu.RUnlock()
	return info, true
}

// InfoBySymbol returns metadata for a token using its symbol.
func (r *Registry) InfoBySymbol(symbol string) (TokenInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tokens {
		if t.Symbol() == symbol {
			return TokenInfo{
				ID:          t.ID(),
				Name:        t.Name(),
				Symbol:      t.Symbol(),
				Decimals:    t.Decimals(),
				TotalSupply: t.TotalSupply(),
			}, true
		}
	}
	return TokenInfo{}, false
}

// List returns metadata for all registered tokens.
func (r *Registry) List() []TokenInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	infos := make([]TokenInfo, 0, len(r.tokens))
	for _, t := range r.tokens {
		infos = append(infos, TokenInfo{
			ID:          t.ID(),
			Name:        t.Name(),
			Symbol:      t.Symbol(),
			Decimals:    t.Decimals(),
			TotalSupply: t.TotalSupply(),
		})
	}
	return infos
}
