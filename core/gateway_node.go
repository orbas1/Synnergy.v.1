package core

import (
	"fmt"

	"synnergy/nodes"
)

// GatewayConfig bundles dependencies required for a GatewayNode.
type GatewayConfig struct {
	Adapter nodes.NodeInterface
}

// GatewayNode exposes a minimal interface for bridging external data sources.
type GatewayNode struct {
	*BaseNode
	cfg      GatewayConfig
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
func (g *GatewayNode) RegisterEndpoint(name string, fn func([]byte) error) {
	g.handlers[name] = fn
}

// Handle invokes a registered endpoint handler.
func (g *GatewayNode) Handle(name string, data []byte) error {
	if h, ok := g.handlers[name]; ok {
		return h(data)
	}
	return fmt.Errorf("unknown endpoint: %s", name)
}
