package api

import "testing"

func TestRateLimiterAllow(t *testing.T) {
	r := NewRateLimiter(1)
	if !r.Allow() {
		t.Fatal("expected first request to be allowed")
	}
	if r.Allow() {
		t.Fatal("expected second request to be denied")
	}
}
