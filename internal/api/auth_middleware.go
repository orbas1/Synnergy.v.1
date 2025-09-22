package api

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"synnergy/internal/auth"
)

// TokenClaims represents validated authentication claims extracted from a token.
type TokenClaims struct {
	Subject     string
	IssuedAt    time.Time
	NotBefore   time.Time
	ExpiresAt   time.Time
	Permissions []auth.Permission
	Metadata    map[string]any
	TokenID     string
}

// HasPermission reports whether the claims contain the required permission.
func (c TokenClaims) HasPermission(perm auth.Permission) bool {
	if perm == "" {
		return true
	}
	for _, p := range c.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// AuthMiddlewareOption configures the middleware instance.
type AuthMiddlewareOption func(*AuthMiddleware)

// WithPolicyEnforcer wires an RBAC policy enforcer into the middleware.
func WithPolicyEnforcer(enforcer *auth.PolicyEnforcer) AuthMiddlewareOption {
	return func(a *AuthMiddleware) {
		a.enforcer = enforcer
	}
}

// WithRateLimiter sets a limiter used to protect against brute-force attacks.
func WithRateLimiter(limiter *RateLimiter) AuthMiddlewareOption {
	return func(a *AuthMiddleware) {
		a.rateLimiter = limiter
	}
}

// WithClock overrides the clock used for token validation.
func WithClock(clock func() time.Time) AuthMiddlewareOption {
	return func(a *AuthMiddleware) {
		if clock != nil {
			a.now = clock
		}
	}
}

// WithLeeway configures clock skew tolerance for token validation.
func WithLeeway(leeway time.Duration) AuthMiddlewareOption {
	return func(a *AuthMiddleware) {
		a.leeway = leeway
	}
}

var (
	// ErrTokenRequired is returned when no token is provided.
	ErrTokenRequired = errors.New("auth: token required")
	// ErrInvalidToken indicates the token structure or signature is invalid.
	ErrInvalidToken = errors.New("auth: invalid token")
	// ErrExpiredToken indicates the token is expired.
	ErrExpiredToken = errors.New("auth: token expired")
)

// ErrRateLimited indicates the request should be retried after a delay.
type ErrRateLimited struct {
	RetryAfter time.Duration
}

// Error implements the error interface.
func (e ErrRateLimited) Error() string {
	if e.RetryAfter <= 0 {
		return "auth: rate limit exceeded"
	}
	return fmt.Sprintf("auth: rate limit exceeded, retry after %s", e.RetryAfter)
}

// AuthMiddleware performs token validation and permission enforcement.
type AuthMiddleware struct {
	secret      []byte
	enforcer    *auth.PolicyEnforcer
	rateLimiter *RateLimiter
	now         func() time.Time
	leeway      time.Duration
}

// NewAuthMiddleware initialises middleware secured with the provided secret.
func NewAuthMiddleware(secret []byte, opts ...AuthMiddlewareOption) *AuthMiddleware {
	if len(secret) == 0 {
		secret = make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			panic(fmt.Errorf("auth: unable to seed secret: %w", err))
		}
	}
	m := &AuthMiddleware{
		secret: secret,
		now:    time.Now,
		leeway: 15 * time.Second,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(m)
		}
	}
	return m
}

// Authenticate validates a token returning true if it is structurally valid.
func (a *AuthMiddleware) Authenticate(token string) bool {
	_, err := a.AuthenticateContext(context.Background(), token, "", nil)
	return err == nil
}

// AuthenticateContext validates the token and enforces the required permission.
func (a *AuthMiddleware) AuthenticateContext(ctx context.Context, token string, required auth.Permission, metadata map[string]any) (*TokenClaims, error) {
	if token == "" {
		return nil, ErrTokenRequired
	}
	claims, err := a.parseToken(token)
	if err != nil {
		return nil, err
	}
	now := a.now()
	if !claims.NotBefore.IsZero() && now.Add(a.leeway).Before(claims.NotBefore) {
		return nil, ErrInvalidToken
	}
	if !claims.ExpiresAt.IsZero() && now.After(claims.ExpiresAt.Add(a.leeway)) {
		return nil, ErrExpiredToken
	}
	if a.rateLimiter != nil {
		if retry, ok := a.rateLimiter.AllowFor(claims.Subject); !ok {
			return nil, ErrRateLimited{RetryAfter: retry}
		}
	}
	if required != "" {
		if !claims.HasPermission(required) {
			return nil, auth.ErrUnauthorized
		}
		if a.enforcer != nil {
			meta := map[string]any{"token_id": claims.TokenID, "issued_at": claims.IssuedAt, "expires_at": claims.ExpiresAt}
			for k, v := range metadata {
				meta[k] = v
			}
			if err := a.enforcer.Authorize(claims.Subject, required, meta); err != nil {
				return nil, err
			}
		}
	}
	return claims, nil
}

// SignToken creates an HMAC-signed token for the provided claims. Intended for testing and CLI tooling.
func (a *AuthMiddleware) SignToken(claims TokenClaims) (string, error) {
	if len(a.secret) == 0 {
		return "", ErrInvalidToken
	}
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := map[string]any{
		"sub":   claims.Subject,
		"iat":   claims.IssuedAt.Unix(),
		"nbf":   claims.NotBefore.Unix(),
		"exp":   claims.ExpiresAt.Unix(),
		"perms": permissionsToStrings(claims.Permissions),
		"meta":  claims.Metadata,
		"jti":   claims.TokenID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := header + "." + payloadEncoded
	mac := hmac.New(sha256.New, a.secret)
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signingInput + "." + signature, nil
}

func (a *AuthMiddleware) parseToken(token string) (*TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}
	mac := hmac.New(sha256.New, a.secret)
	mac.Write([]byte(parts[0]))
	mac.Write([]byte("."))
	mac.Write([]byte(parts[1]))
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return nil, ErrInvalidToken
	}
	var payload struct {
		Subject     string         `json:"sub"`
		IssuedAt    int64          `json:"iat"`
		NotBefore   int64          `json:"nbf"`
		ExpiresAt   int64          `json:"exp"`
		Permissions []string       `json:"perms"`
		Metadata    map[string]any `json:"meta"`
		TokenID     string         `json:"jti"`
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, ErrInvalidToken
	}
	claims := &TokenClaims{
		Subject:     payload.Subject,
		IssuedAt:    unixToTime(payload.IssuedAt),
		NotBefore:   unixToTime(payload.NotBefore),
		ExpiresAt:   unixToTime(payload.ExpiresAt),
		Permissions: stringsToPermissions(payload.Permissions),
		Metadata:    payload.Metadata,
		TokenID:     payload.TokenID,
	}
	if claims.Metadata == nil {
		claims.Metadata = map[string]any{}
	}
	return claims, nil
}

func unixToTime(sec int64) time.Time {
	if sec == 0 {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}

func permissionsToStrings(perms []auth.Permission) []string {
	out := make([]string, len(perms))
	for i, p := range perms {
		out[i] = string(p)
	}
	return out
}

func stringsToPermissions(perms []string) []auth.Permission {
	out := make([]auth.Permission, len(perms))
	for i, p := range perms {
		out[i] = auth.Permission(p)
	}
	return out
}
