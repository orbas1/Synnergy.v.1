package api

import "testing"

func TestGatewayStart(t *testing.T) {
	g := NewGateway()
	if err := g.Start(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
