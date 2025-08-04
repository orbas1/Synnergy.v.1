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
