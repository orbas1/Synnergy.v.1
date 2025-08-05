package tokens

// Registry maintains a lookup of token instances by their identifiers.
type Registry struct {
	tokens map[TokenID]Token
	next   TokenID
}

// NewRegistry creates an empty token registry.
func NewRegistry() *Registry {
	return &Registry{tokens: make(map[TokenID]Token)}
}

// NextID returns a unique identifier for a new token.
func (r *Registry) NextID() TokenID {
	r.next++
	return r.next
}

// Register adds the token to the registry using its ID.
func (r *Registry) Register(t Token) {
	r.tokens[t.ID()] = t
}

// Get retrieves a token by ID if present.
func (r *Registry) Get(id TokenID) (Token, bool) {
	t, ok := r.tokens[id]
	return t, ok
}

// GetBySymbol retrieves a token by its symbol.
func (r *Registry) GetBySymbol(symbol string) (Token, bool) {
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
	t, ok := r.tokens[id]
	if !ok {
		return TokenInfo{}, false
	}
	return TokenInfo{
		ID:          t.ID(),
		Name:        t.Name(),
		Symbol:      t.Symbol(),
		Decimals:    t.Decimals(),
		TotalSupply: t.TotalSupply(),
	}, true
}

// InfoBySymbol returns metadata for a token using its symbol.
func (r *Registry) InfoBySymbol(symbol string) (TokenInfo, bool) {
	t, ok := r.GetBySymbol(symbol)
	if !ok {
		return TokenInfo{}, false
	}
	return r.Info(t.ID())
}

// List returns metadata for all registered tokens.
func (r *Registry) List() []TokenInfo {
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
