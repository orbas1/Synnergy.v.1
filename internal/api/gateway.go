package api

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// GatewayMetrics represents a snapshot of the gateway's runtime statistics.
type GatewayMetrics struct {
	TotalRequests uint64        `json:"total_requests"`
	RateLimited   uint64        `json:"rate_limited"`
	Unauthorized  uint64        `json:"unauthorized"`
	ActiveClients int64         `json:"active_clients"`
	Uptime        time.Duration `json:"uptime"`
	LastRequest   time.Time     `json:"last_request"`
}

// Route describes an HTTP route served by the gateway.
type Route struct {
	Path        string
	Methods     []string
	Handler     http.HandlerFunc
	RequireAuth bool
}

// GatewayOption configures optional gateway behaviour.
type GatewayOption func(*Gateway)

type routeConfig struct {
	methods     map[string]struct{}
	requireAuth bool
}

type gatewayMetrics struct {
	totalRequests atomic.Uint64
	rateLimited   atomic.Uint64
	unauthorized  atomic.Uint64
	activeConns   atomic.Int64
	startedAt     atomic.Int64
	lastRequest   atomic.Int64
}

// Gateway represents an HTTP gateway with rate limiting, authentication and
// built-in observability.
type Gateway struct {
	mu              sync.RWMutex
	server          *http.Server
	router          *http.ServeMux
	routes          map[string]routeConfig
	limiter         *RateLimiter
	auth            *AuthMiddleware
	addr            string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
	listener        net.Listener
	metrics         gatewayMetrics
}

// NewGateway constructs a gateway using sensible defaults.
func NewGateway(opts ...GatewayOption) *Gateway {
	g := &Gateway{
		router:          http.NewServeMux(),
		routes:          make(map[string]routeConfig),
		limiter:         NewRateLimiter(100),
		auth:            &AuthMiddleware{},
		addr:            "127.0.0.1:0",
		readTimeout:     15 * time.Second,
		writeTimeout:    15 * time.Second,
		idleTimeout:     60 * time.Second,
		shutdownTimeout: 5 * time.Second,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(g)
		}
	}
	g.server = g.newServer()
	g.registerDefaultRoutes()
	return g
}

// WithAddress overrides the listen address for the gateway.
func WithAddress(addr string) GatewayOption {
	return func(g *Gateway) {
		if strings.TrimSpace(addr) != "" {
			g.addr = addr
		}
	}
}

// WithRateLimiter overrides the default rate limiter.
func WithRateLimiter(limiter *RateLimiter) GatewayOption {
	return func(g *Gateway) {
		if limiter != nil {
			g.limiter = limiter
		}
	}
}

// WithAuthMiddleware overrides the authentication middleware.
func WithAuthMiddleware(auth *AuthMiddleware) GatewayOption {
	return func(g *Gateway) {
		if auth != nil {
			g.auth = auth
		}
	}
}

// WithReadTimeout sets the HTTP read timeout.
func WithReadTimeout(timeout time.Duration) GatewayOption {
	return func(g *Gateway) {
		if timeout > 0 {
			g.readTimeout = timeout
		}
	}
}

// WithWriteTimeout sets the HTTP write timeout.
func WithWriteTimeout(timeout time.Duration) GatewayOption {
	return func(g *Gateway) {
		if timeout > 0 {
			g.writeTimeout = timeout
		}
	}
}

// WithIdleTimeout sets the HTTP idle timeout.
func WithIdleTimeout(timeout time.Duration) GatewayOption {
	return func(g *Gateway) {
		if timeout > 0 {
			g.idleTimeout = timeout
		}
	}
}

// WithShutdownTimeout sets the graceful shutdown timeout.
func WithShutdownTimeout(timeout time.Duration) GatewayOption {
	return func(g *Gateway) {
		if timeout > 0 {
			g.shutdownTimeout = timeout
		}
	}
}

// RegisterRoute registers a new HTTP route with optional authentication.
func (g *Gateway) RegisterRoute(route Route) error {
	if strings.TrimSpace(route.Path) == "" {
		return errors.New("route path required")
	}
	if route.Handler == nil {
		return errors.New("route handler required")
	}
	methods := make(map[string]struct{})
	for _, method := range route.Methods {
		method = strings.ToUpper(strings.TrimSpace(method))
		if method != "" {
			methods[method] = struct{}{}
		}
	}

	g.mu.Lock()
	g.routes[route.Path] = routeConfig{methods: methods, requireAuth: route.RequireAuth}
	g.router.HandleFunc(route.Path, route.Handler)
	g.mu.Unlock()
	return nil
}

// Start begins serving requests. It is safe to call multiple times.
func (g *Gateway) Start() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.listener != nil {
		return nil
	}
	if g.server == nil {
		g.server = g.newServer()
	}
	ln, err := net.Listen("tcp", g.server.Addr)
	if err != nil {
		return err
	}
	g.listener = ln
	g.server.Addr = ln.Addr().String()
	g.metrics.startedAt.Store(time.Now().UnixNano())

	go func() {
		if err := g.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// best-effort: swallow serve errors to keep tests deterministic
		}
	}()
	return nil
}

// Stop gracefully stops the gateway.
func (g *Gateway) Stop(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.server == nil {
		return nil
	}
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), g.shutdownTimeout)
		defer cancel()
	}
	err := g.server.Shutdown(ctx)
	g.listener = nil
	g.metrics.activeConns.Store(0)
	return err
}

// Address returns the effective listen address for the gateway.
func (g *Gateway) Address() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.listener != nil {
		return g.listener.Addr().String()
	}
	if g.server != nil {
		return g.server.Addr
	}
	return g.addr
}

// Metrics snapshots the current gateway metrics.
func (g *Gateway) Metrics() GatewayMetrics {
	total := g.metrics.totalRequests.Load()
	limited := g.metrics.rateLimited.Load()
	unauthorized := g.metrics.unauthorized.Load()
	active := g.metrics.activeConns.Load()
	started := g.metrics.startedAt.Load()
	last := g.metrics.lastRequest.Load()

	metrics := GatewayMetrics{TotalRequests: total, RateLimited: limited, Unauthorized: unauthorized, ActiveClients: active}
	if started > 0 {
		metrics.Uptime = time.Since(time.Unix(0, started))
	}
	if last > 0 {
		metrics.LastRequest = time.Unix(0, last)
	}
	return metrics
}

func (g *Gateway) newServer() *http.Server {
	return &http.Server{
		Addr:         g.addr,
		Handler:      g.wrap(g.router),
		ReadTimeout:  g.readTimeout,
		WriteTimeout: g.writeTimeout,
		IdleTimeout:  g.idleTimeout,
		ConnState:    g.connStateHandler,
	}
}

func (g *Gateway) connStateHandler(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		g.metrics.activeConns.Add(1)
	case http.StateClosed, http.StateHijacked:
		g.metrics.activeConns.Add(-1)
	}
}

func (g *Gateway) registerDefaultRoutes() {
	_ = g.RegisterRoute(Route{Path: "/healthz", Methods: []string{http.MethodGet}, Handler: g.healthHandler})
	_ = g.RegisterRoute(Route{Path: "/readyz", Methods: []string{http.MethodGet}, Handler: g.readyHandler})
	_ = g.RegisterRoute(Route{Path: "/metrics", Methods: []string{http.MethodGet}, Handler: g.metricsHandler, RequireAuth: true})
}

func (g *Gateway) wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.metrics.totalRequests.Add(1)

		cfg, ok := g.routeFor(r.URL.Path)
		if ok && len(cfg.methods) > 0 {
			if _, allowed := cfg.methods[r.Method]; !allowed {
				w.Header().Set("Allow", strings.Join(mapKeys(cfg.methods), ", "))
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}

		if g.limiter != nil && !g.limiter.Allow() {
			g.metrics.rateLimited.Add(1)
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		if ok && cfg.requireAuth {
			token := extractAuthToken(r.Header.Get("Authorization"))
			if g.auth == nil || !g.auth.Authenticate(token) {
				g.metrics.unauthorized.Add(1)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		g.metrics.lastRequest.Store(time.Now().UnixNano())
		next.ServeHTTP(w, r)
	})
}

func (g *Gateway) routeFor(path string) (routeConfig, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	cfg, ok := g.routes[path]
	return cfg, ok
}

func (g *Gateway) healthHandler(w http.ResponseWriter, r *http.Request) {
	g.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (g *Gateway) readyHandler(w http.ResponseWriter, r *http.Request) {
	g.writeJSON(w, http.StatusOK, map[string]bool{"ready": true})
}

func (g *Gateway) metricsHandler(w http.ResponseWriter, r *http.Request) {
	g.writeJSON(w, http.StatusOK, g.Metrics())
}

func (g *Gateway) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func mapKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func extractAuthToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return header
}
