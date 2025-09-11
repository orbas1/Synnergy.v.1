package core

import (
	"fmt"
	"sync"
	"testing"

	"synnergy/internal/nodes"
)

func TestGatewayNodeEndpoints(t *testing.T) {
	gn := NewGatewayNode(nodes.Address("g1"), GatewayConfig{})
	if err := gn.RegisterEndpoint("ping", func(b []byte) error { return nil }); err == nil {
		t.Fatalf("expected error when node not running")
	}

	if err := gn.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	called := false
	if err := gn.RegisterEndpoint("ping", func(b []byte) error { called = true; return nil }); err != nil {
		t.Fatalf("register: %v", err)
	}
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

// TestGatewayNodeConcurrency registers and invokes endpoints concurrently.
func TestGatewayNodeConcurrency(t *testing.T) {
	gn := NewGatewayNode(nodes.Address("g2"), GatewayConfig{})
	_ = gn.Start()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("ep%d", i)
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			_ = gn.RegisterEndpoint(n, func([]byte) error { return nil })
			_ = gn.Handle(n, nil)
		}(name)
	}
	wg.Wait()
	if got := len(gn.Endpoints()); got != 100 {
		t.Fatalf("expected 100 endpoints, got %d", got)
	}
}
