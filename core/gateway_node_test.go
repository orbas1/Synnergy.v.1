package core

import (
	"testing"

	"synnergy/internal/nodes"
)

func TestGatewayNodeEndpoints(t *testing.T) {
	gn := NewGatewayNode(nodes.Address("g1"), GatewayConfig{})
	called := false
	gn.RegisterEndpoint("ping", func(b []byte) error {
		called = true
		return nil
	})
	if err := gn.Handle("ping", nil); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if !called {
		t.Fatalf("endpoint not invoked")
	}
	if len(gn.Endpoints()) != 1 || gn.Endpoints()[0] != "ping" {
		t.Fatalf("unexpected endpoints list")
	}
	gn.RemoveEndpoint("ping")
	if len(gn.Endpoints()) != 0 {
		t.Fatalf("endpoint not removed")
	}
	if err := gn.Handle("ping", nil); err == nil {
		t.Fatalf("expected error for unknown endpoint")
	}
}
