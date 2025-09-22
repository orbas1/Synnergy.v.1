package api

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"synnergy/internal/auth"
)

func TestAuthenticateContextSuccess(t *testing.T) {
	secret := bytes.Repeat([]byte("s"), 32)
	rbac := auth.NewRBAC()
	rbac.AddRole("admin")
	if err := rbac.AddPermissionToRole("admin", auth.Permission("write")); err != nil {
		t.Fatalf("add permission: %v", err)
	}
	if err := rbac.AssignRole("user1", "admin"); err != nil {
		t.Fatalf("assign role: %v", err)
	}
	logger := auth.NewMemoryAuditLogger()
	enforcer := auth.NewPolicyEnforcer(rbac, logger)
	limiter := NewRateLimiter(5)
	middleware := NewAuthMiddleware(secret, WithPolicyEnforcer(enforcer), WithRateLimiter(limiter), WithLeeway(time.Second))

	claims := TokenClaims{
		Subject:     "user1",
		IssuedAt:    time.Now().Add(-time.Minute),
		NotBefore:   time.Now().Add(-time.Minute),
		ExpiresAt:   time.Now().Add(5 * time.Minute),
		Permissions: []auth.Permission{"write"},
		TokenID:     "tok-1",
	}
	token, err := middleware.SignToken(claims)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	got, err := middleware.AuthenticateContext(context.Background(), token, auth.Permission("write"), map[string]any{"resource": "doc"})
	if err != nil {
		t.Fatalf("authenticate: %v", err)
	}
	if got.Subject != "user1" {
		t.Fatalf("unexpected subject: %v", got.Subject)
	}
	if len(logger.Events()) == 0 {
		t.Fatalf("expected audit events to be captured")
	}
}

func TestAuthenticateContextRateLimit(t *testing.T) {
	secret := bytes.Repeat([]byte("r"), 32)
	middleware := NewAuthMiddleware(secret, WithRateLimiter(NewRateLimiter(1)))
	claims := TokenClaims{
		Subject:     "user2",
		IssuedAt:    time.Now(),
		ExpiresAt:   time.Now().Add(time.Minute),
		Permissions: []auth.Permission{"read"},
	}
	token, err := middleware.SignToken(claims)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	if _, err := middleware.AuthenticateContext(context.Background(), token, "", nil); err != nil {
		t.Fatalf("first call expected success: %v", err)
	}
	if _, err := middleware.AuthenticateContext(context.Background(), token, "", nil); err == nil {
		t.Fatalf("expected rate limit error on second call")
	}
}

func TestAuthenticateContextInvalidToken(t *testing.T) {
	secret := bytes.Repeat([]byte("x"), 32)
	middleware := NewAuthMiddleware(secret)
	if middleware.Authenticate("bad.token") {
		t.Fatalf("expected invalid token to fail")
	}
	claims := TokenClaims{Subject: "user", IssuedAt: time.Now(), ExpiresAt: time.Now().Add(time.Minute)}
	token, err := middleware.SignToken(claims)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	// tamper with token
	tampered := token + "extra"
	if _, err := middleware.AuthenticateContext(context.Background(), tampered, "", nil); err == nil {
		t.Fatalf("expected tampered token to fail")
	}
}

func TestAuthenticateContextPermissionDenied(t *testing.T) {
	secret := bytes.Repeat([]byte("p"), 32)
	rbac := auth.NewRBAC()
	logger := auth.NewMemoryAuditLogger()
	middleware := NewAuthMiddleware(secret, WithPolicyEnforcer(auth.NewPolicyEnforcer(rbac, logger)))
	claims := TokenClaims{
		Subject:     "user3",
		IssuedAt:    time.Now(),
		ExpiresAt:   time.Now().Add(time.Minute),
		Permissions: []auth.Permission{"read"},
	}
	token, err := middleware.SignToken(claims)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	if _, err := middleware.AuthenticateContext(context.Background(), token, auth.Permission("admin"), nil); !errors.Is(err, auth.ErrUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}
