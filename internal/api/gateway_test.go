package api

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"synnergy/internal/auth"
)

func TestGatewayStartAndRequests(t *testing.T) {
	secret := []byte("abcdefghijklmnopqrstuvwxyz123456")
	rbac := auth.NewRBAC()
	rbac.AddRole("writer")
	if err := rbac.AddPermissionToRole("writer", auth.Permission("write")); err != nil {
		t.Fatalf("add permission: %v", err)
	}
	if err := rbac.AssignRole("user", "writer"); err != nil {
		t.Fatalf("assign role: %v", err)
	}
	middleware := NewAuthMiddleware(secret, WithPolicyEnforcer(auth.NewPolicyEnforcer(rbac, auth.NewMemoryAuditLogger())))
	gateway := NewGateway(WithGatewayAuth(middleware), WithGatewayRateLimiter(NewRateLimiter(1)))
	gateway.HandleFunc("/secured", auth.Permission("write"), func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := gateway.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	addr := gateway.Addr()
	time.Sleep(50 * time.Millisecond)

	// valid request
	token, err := middleware.SignToken(TokenClaims{
		Subject:     "user",
		IssuedAt:    time.Now().Add(-time.Minute),
		ExpiresAt:   time.Now().Add(time.Minute),
		Permissions: []auth.Permission{"write"},
	})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+"/secured", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// second request should hit rate limit
	req2, _ := http.NewRequest(http.MethodGet, "http://"+addr+"/secured", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("rate limit request: %v", err)
	}
	if resp2.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", resp2.StatusCode)
	}
	resp2.Body.Close()

	time.Sleep(time.Second)

	// unauthorized request without token
	req3, _ := http.NewRequest(http.MethodGet, "http://"+addr+"/secured", nil)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("unauthorized request: %v", err)
	}
	if resp3.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp3.StatusCode)
	}
	resp3.Body.Close()

	cancel()
	time.Sleep(50 * time.Millisecond)
}
