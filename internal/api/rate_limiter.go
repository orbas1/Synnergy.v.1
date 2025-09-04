package api

import "sync"

// RateLimiter limits requests to a fixed count.
type RateLimiter struct {
	mu    sync.Mutex
	count int
	limit int
}

// NewRateLimiter creates a RateLimiter.
func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{limit: limit}
}

// Allow reports if a request is allowed.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.count >= r.limit {
		return false
	}
	r.count++
	return true
}
