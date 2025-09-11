package core

import (
	"fmt"
	"sync"

	"synnergy/internal/nodes"
)

// GatewayConfig bundles dependencies required for a GatewayNode.
type GatewayConfig struct {
	Adapter nodes.NodeInterface
}

// GatewayNode exposes a minimal interface for bridging external data sources.
type GatewayNode struct {
	*BaseNode
	cfg GatewayConfig

	mu       sync.RWMutex
	handlers map[string]func([]byte) error
}

// NewGatewayNode constructs a gateway node using the provided config.
func NewGatewayNode(id nodes.Address, cfg GatewayConfig) *GatewayNode {
	return &GatewayNode{
		BaseNode: NewBaseNode(id),
		cfg:      cfg,
		handlers: make(map[string]func([]byte) error),
	}
}

// RegisterEndpoint registers a handler for an external endpoint name.
func (g *GatewayNode) RegisterEndpoint(name string, fn func([]byte) error) error {
	if !g.IsRunning() {
		return fmt.Errorf("gateway not running")
	}
	g.mu.Lock()
	g.handlers[name] = fn
	g.mu.Unlock()
	return nil
}

// Handle invokes a registered endpoint handler.
func (g *GatewayNode) Handle(name string, data []byte) error {
	if !g.IsRunning() {
		return fmt.Errorf("gateway not running")
	}
	g.mu.RLock()
	h, ok := g.handlers[name]
	g.mu.RUnlock()
	if !ok {
		return fmt.Errorf("unknown endpoint: %s", name)
	}
	return h(data)
}

// RemoveEndpoint removes a previously registered endpoint handler.
func (g *GatewayNode) RemoveEndpoint(name string) {
	g.mu.Lock()
	delete(g.handlers, name)
	g.mu.Unlock()
}

// Endpoints returns a list of registered endpoint names.
func (g *GatewayNode) Endpoints() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]string, 0, len(g.handlers))
	for k := range g.handlers {
		out = append(out, k)
	}
	return out
}
