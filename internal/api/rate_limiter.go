package api

import (
	"math"
	"sync"
	"time"
)

// RateLimiter limits requests across identities using a token bucket algorithm.
type RateLimiter struct {
	mu      sync.Mutex
	limit   float64
	burst   float64
	window  time.Duration
	entries map[string]*rateBucket
	now     func() time.Time
}

type rateBucket struct {
	tokens float64
	last   time.Time
}

// RateLimiterOption configures the limiter.
type RateLimiterOption func(*RateLimiter)

// WithRateWindow configures the time window for refilling tokens.
func WithRateWindow(window time.Duration) RateLimiterOption {
	return func(r *RateLimiter) {
		if window > 0 {
			r.window = window
		}
	}
}

// WithRateBurst configures the burst size for the limiter.
func WithRateBurst(burst int) RateLimiterOption {
	return func(r *RateLimiter) {
		if burst > 0 {
			r.burst = float64(burst)
		}
	}
}

// WithRateClock overrides the clock used by the limiter.
func WithRateClock(clock func() time.Time) RateLimiterOption {
	return func(r *RateLimiter) {
		if clock != nil {
			r.now = clock
		}
	}
}

// NewRateLimiter creates a RateLimiter with the provided limit per window.
func NewRateLimiter(limit int, opts ...RateLimiterOption) *RateLimiter {
	if limit <= 0 {
		limit = 1
	}
	rl := &RateLimiter{
		limit:   float64(limit),
		burst:   float64(limit),
		window:  time.Second,
		entries: make(map[string]*rateBucket),
		now:     time.Now,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(rl)
		}
	}
	if rl.burst < rl.limit {
		rl.burst = rl.limit
	}
	return rl
}

// Allow reports whether a request is allowed for the global bucket.
func (r *RateLimiter) Allow() bool {
	_, ok := r.AllowFor("__global__")
	return ok
}

// AllowFor reports whether a request for the given key can proceed and returns the retry-after duration when limited.
func (r *RateLimiter) AllowFor(key string) (time.Duration, bool) {
	return r.AllowN(key, 1)
}

// AllowN consumes the specified weight from the key's bucket.
func (r *RateLimiter) AllowN(key string, weight float64) (time.Duration, bool) {
	if weight <= 0 {
		return 0, true
	}
	now := r.now()
	r.mu.Lock()
	defer r.mu.Unlock()
	bucket := r.entries[key]
	if bucket == nil {
		bucket = &rateBucket{tokens: r.burst, last: now}
		r.entries[key] = bucket
	} else {
		r.refill(bucket, now)
	}
	if bucket.tokens < weight {
		needed := weight - bucket.tokens
		retry := r.retryAfter(needed)
		return retry, false
	}
	bucket.tokens -= weight
	return 0, true
}

// Reset clears the limiter state for a given key.
func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	delete(r.entries, key)
	r.mu.Unlock()
}

// Active returns the number of tracked identities.
func (r *RateLimiter) Active() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

func (r *RateLimiter) refill(bucket *rateBucket, now time.Time) {
	if r.window <= 0 {
		bucket.tokens = r.burst
		bucket.last = now
		return
	}
	elapsed := now.Sub(bucket.last)
	if elapsed <= 0 {
		return
	}
	ratePerSecond := r.limit / r.window.Seconds()
	bucket.tokens = math.Min(r.burst, bucket.tokens+elapsed.Seconds()*ratePerSecond)
	bucket.last = now
}

func (r *RateLimiter) retryAfter(needed float64) time.Duration {
	if r.window <= 0 {
		return 0
	}
	ratePerSecond := r.limit / r.window.Seconds()
	if ratePerSecond <= 0 {
		return time.Second
	}
	seconds := math.Ceil(needed / ratePerSecond)
	if seconds < 1 {
		seconds = 1
	}
	return time.Duration(seconds) * time.Second
}
