package security

import (
	"math"
	"sync"
	"time"

	"synnergy/internal/telemetry"
)

// RateLimiterOption configures the enterprise-grade rate limiter.
type RateLimiterOption func(*RateLimiter)

// WithBurst configures the number of tokens that may accumulate and be spent in
// a burst. When unset the limiter allows a single token burst.
func WithBurst(burst int) RateLimiterOption {
	return func(r *RateLimiter) {
		if burst > 0 {
			r.burst = float64(burst)
		}
	}
}

// WithClock overrides the time source. It is primarily used by tests to drive
// deterministic scenarios.
func WithClock(clock func() time.Time) RateLimiterOption {
	return func(r *RateLimiter) {
		if clock != nil {
			r.now = clock
		}
	}
}

// WithComponent tags emitted telemetry samples with a logical component name
// (for example "rpc" or "cli").
func WithComponent(component string) RateLimiterOption {
	return func(r *RateLimiter) {
		if component != "" {
			r.component = component
		}
	}
}

// WithSink overrides the telemetry sink used to publish rate limit decisions.
func WithSink(sink telemetry.RateLimitRecorder) RateLimiterOption {
	return func(r *RateLimiter) {
		if sink != nil {
			r.sink = sink
		}
	}
}

// RateLimitState exposes the internal state of a particular identity bucket so
// that CLI and monitoring integrations can render the remaining budget.
type RateLimitState struct {
	Remaining  float64
	LastRefill time.Time
}

type bucket struct {
	tokens float64
	last   time.Time
}

// RateLimiter implements a concurrent-safe, per-identity token bucket with
// burst support, telemetry emission and dynamic reconfiguration hooks.
type RateLimiter struct {
	mu        sync.Mutex
	interval  time.Duration
	burst     float64
	now       func() time.Time
	buckets   map[string]*bucket
	component string
	sink      telemetry.RateLimitRecorder
}

// NewRateLimiter creates a new limiter. The interval represents the cost of a
// single token; every interval a new token is released into each bucket up to
// the configured burst capacity.
func NewRateLimiter(interval time.Duration, opts ...RateLimiterOption) *RateLimiter {
	if interval <= 0 {
		interval = time.Second
	}
	rl := &RateLimiter{
		interval:  interval,
		burst:     1,
		now:       time.Now,
		buckets:   make(map[string]*bucket),
		component: "security",
		sink:      telemetry.GlobalRateLimitSink(),
	}
	for _, opt := range opts {
		opt(rl)
	}
	return rl
}

// Allow records an anonymous request using the default identity. The method is
// preserved for backwards compatibility with legacy CLI integrations.
func (r *RateLimiter) Allow() bool {
	allowed, _ := r.AllowN("default", 1)
	return allowed
}

// AllowN attempts to consume cost tokens for the specified identity. The
// returned duration indicates how long the caller should wait before retrying
// when the request is denied.
func (r *RateLimiter) AllowN(identity string, cost float64) (bool, time.Duration) {
	if identity == "" {
		identity = "default"
	}
	if cost <= 0 {
		cost = 1
	}
	now := r.now()

	r.mu.Lock()
	b := r.ensureBucket(identity)
	elapsed := now.Sub(b.last)
	if elapsed < 0 {
		elapsed = 0
	}
	if b.last.IsZero() {
		b.tokens = r.burst
	} else if elapsed > 0 {
		refill := elapsed.Seconds() / r.interval.Seconds()
		if refill > 0 {
			b.tokens = math.Min(r.burst, b.tokens+refill)
		}
	}
	b.last = now

	var allowed bool
	var retry time.Duration
	if b.tokens >= cost {
		b.tokens -= cost
		allowed = true
	} else {
		shortfall := cost - b.tokens
		retry = time.Duration(math.Ceil(shortfall * float64(r.interval)))
		b.tokens = 0
	}
	remaining := b.tokens
	r.mu.Unlock()

	r.record(identity, allowed, remaining, retry)
	return allowed, retry
}

// Snapshot exposes the limiter state for all tracked identities.
func (r *RateLimiter) Snapshot() map[string]RateLimitState {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]RateLimitState, len(r.buckets))
	for id, b := range r.buckets {
		out[id] = RateLimitState{
			Remaining:  b.tokens,
			LastRefill: b.last,
		}
	}
	return out
}

// Reconfigure dynamically updates the limiter interval and burst capacity.
func (r *RateLimiter) Reconfigure(interval time.Duration, burst int) {
	if interval <= 0 && burst <= 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if interval > 0 {
		r.interval = interval
	}
	if burst > 0 {
		r.burst = float64(burst)
		for _, b := range r.buckets {
			if b.tokens > r.burst {
				b.tokens = r.burst
			}
		}
	}
}

func (r *RateLimiter) ensureBucket(identity string) *bucket {
	b, ok := r.buckets[identity]
	if !ok {
		b = &bucket{tokens: r.burst}
		r.buckets[identity] = b
	}
	return b
}

func (r *RateLimiter) record(identity string, allowed bool, remaining float64, retry time.Duration) {
	if r.sink == nil {
		return
	}
	r.sink.RecordRateLimit(r.component, identity, allowed, remaining, retry)
}
