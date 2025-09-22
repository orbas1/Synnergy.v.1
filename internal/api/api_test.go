package api

import (
	"context"
	"net/http"
	"testing"
)

func TestClaimsFromContext(t *testing.T) {
	ctx := context.Background()
	claims := &TokenClaims{Subject: "user"}
	ctx = context.WithValue(ctx, claimsContextKey{}, claims)
	got, ok := ClaimsFromContext(ctx)
	if !ok || got.Subject != "user" {
		t.Fatalf("expected claims in context")
	}
}

func TestGatewayRoutes(t *testing.T) {
	g := NewGateway()
	g.HandleFunc("/one", "", func(http.ResponseWriter, *http.Request) {})
	g.HandleFunc("/two", "", func(http.ResponseWriter, *http.Request) {})
	routes := g.Routes()
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %v", routes)
	}
	if routes[0] != "/one" || routes[1] != "/two" {
		t.Fatalf("unexpected route order: %v", routes)
	}
}
