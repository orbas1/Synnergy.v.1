package security

import "time"

// RateLimiter is a minimal token bucket implementation.
type RateLimiter struct {
	interval time.Duration
	last     time.Time
}

// NewRateLimiter creates a RateLimiter that allows one event per interval.
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{interval: interval}
}

// Allow reports whether an event is allowed to proceed.
func (r *RateLimiter) Allow() bool {
	now := time.Now()
	if now.Sub(r.last) >= r.interval {
		r.last = now
		return true
	}
	return false
}
