package api

// Gateway represents a minimal API gateway.
type Gateway struct{}

// NewGateway creates a new Gateway.
func NewGateway() *Gateway { return &Gateway{} }

// Start begins serving requests (placeholder).
func (g *Gateway) Start() error { return nil }
