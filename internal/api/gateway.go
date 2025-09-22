package api

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"synnergy/internal/auth"
)

// ErrGatewayRunning indicates the gateway server is already active.
var ErrGatewayRunning = errors.New("gateway: already running")

// routeConfig stores handler metadata.
type routeConfig struct {
	handler    http.Handler
	permission auth.Permission
}

// Gateway represents an API gateway secured with authentication and rate limiting.
type Gateway struct {
	mu          sync.RWMutex
	routes      map[string]routeConfig
	auth        *AuthMiddleware
	rateLimiter *RateLimiter
	server      *http.Server
	listener    net.Listener
	addr        string
	started     atomic.Bool
	logger      func(string)
}

// GatewayOption configures a Gateway instance.
type GatewayOption func(*Gateway)

// WithGatewayAuth configures the authentication middleware.
func WithGatewayAuth(auth *AuthMiddleware) GatewayOption {
	return func(g *Gateway) {
		g.auth = auth
	}
}

// WithGatewayRateLimiter configures the rate limiter for incoming requests.
func WithGatewayRateLimiter(limiter *RateLimiter) GatewayOption {
	return func(g *Gateway) {
		g.rateLimiter = limiter
	}
}

// WithGatewayAddress sets the listening address.
func WithGatewayAddress(addr string) GatewayOption {
	return func(g *Gateway) {
		g.addr = addr
	}
}

// WithGatewayLogger sets a logger function.
func WithGatewayLogger(logger func(string)) GatewayOption {
	return func(g *Gateway) {
		g.logger = logger
	}
}

// NewGateway creates a new Gateway.
func NewGateway(opts ...GatewayOption) *Gateway {
	g := &Gateway{
		routes: make(map[string]routeConfig),
		addr:   "127.0.0.1:0",
	}
	for _, opt := range opts {
		if opt != nil {
			opt(g)
		}
	}
	if g.logger == nil {
		g.logger = func(string) {}
	}
	return g
}

// Handle registers a route guarded by the provided permission.
func (g *Gateway) Handle(pattern string, perm auth.Permission, handler http.Handler) {
	if handler == nil {
		handler = http.NotFoundHandler()
	}
	g.mu.Lock()
	g.routes[pattern] = routeConfig{handler: handler, permission: perm}
	g.mu.Unlock()
}

// HandleFunc registers a route using a function handler.
func (g *Gateway) HandleFunc(pattern string, perm auth.Permission, handler func(http.ResponseWriter, *http.Request)) {
	g.Handle(pattern, perm, http.HandlerFunc(handler))
}

// Start boots the gateway HTTP server.
func (g *Gateway) Start(ctx context.Context) error {
	if g.started.Load() {
		return ErrGatewayRunning
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started.Load() {
		return ErrGatewayRunning
	}
	mux := http.NewServeMux()
	for pattern, cfg := range g.routes {
		mux.Handle(pattern, g.wrapHandler(cfg))
	}
	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	ln, err := net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}
	g.server = srv
	g.listener = ln
	g.started.Store(true)

	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			g.logger("gateway serve error: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			g.logger("gateway shutdown error: " + err.Error())
		}
	}()
	return nil
}

// Shutdown stops the gateway server.
func (g *Gateway) Shutdown(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.started.Load() {
		return nil
	}
	g.started.Store(false)
	if g.server != nil {
		return g.server.Shutdown(ctx)
	}
	return nil
}

// Addr returns the listening address once started.
func (g *Gateway) Addr() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.listener != nil {
		return g.listener.Addr().String()
	}
	return g.addr
}

func (g *Gateway) wrapHandler(cfg routeConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g.rateLimiter != nil {
			if retry, ok := g.rateLimiter.AllowFor(clientKey(r)); !ok {
				if retry > 0 {
					w.Header().Set("Retry-After", fmtDuration(retry))
				}
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		}
		if g.auth != nil {
			token := extractToken(r.Header.Get("Authorization"))
			claims, err := g.auth.AuthenticateContext(r.Context(), token, cfg.permission, map[string]any{"path": r.URL.Path})
			if err != nil {
				status := statusFromAuthError(err)
				http.Error(w, err.Error(), status)
				return
			}
			ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
			r = r.WithContext(ctx)
		}
		cfg.handler.ServeHTTP(w, r)
	})
}

// ClaimsFromContext extracts token claims from the request context.
func ClaimsFromContext(ctx context.Context) (*TokenClaims, bool) {
	val := ctx.Value(claimsContextKey{})
	if val == nil {
		return nil, false
	}
	claims, ok := val.(*TokenClaims)
	return claims, ok
}

type claimsContextKey struct{}

func extractToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return header
}

func clientKey(r *http.Request) string {
	if r.RemoteAddr == "" {
		return r.Host
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func statusFromAuthError(err error) int {
	switch {
	case errors.Is(err, ErrTokenRequired), errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, auth.ErrUnauthorized):
		return http.StatusForbidden
	case errors.As(err, new(ErrRateLimited)):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrExpiredToken):
		return http.StatusUnauthorized
	default:
		return http.StatusUnauthorized
	}
}

func fmtDuration(d time.Duration) string {
	secs := int(math.Ceil(d.Seconds()))
	if secs < 1 {
		secs = 1
	}
	return fmt.Sprintf("%d", secs)
}

// Routes returns registered route patterns for diagnostics.
func (g *Gateway) Routes() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	keys := make([]string, 0, len(g.routes))
	for pattern := range g.routes {
		keys = append(keys, pattern)
	}
	sort.Strings(keys)
	return keys
}
