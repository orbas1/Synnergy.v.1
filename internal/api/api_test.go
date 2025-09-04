package api

import "testing"

func TestGateway(t *testing.T) {
	g := NewGateway()
	if err := g.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
}

func TestAuthMiddleware(t *testing.T) {
	a := &AuthMiddleware{}
	if !a.Authenticate("token") {
		t.Fatalf("expected token to be valid")
	}
}

func TestRateLimiter(t *testing.T) {
	r := NewRateLimiter(1)
	if !r.Allow() || r.Allow() {
		t.Fatalf("rate limiting failed")
	}
}
