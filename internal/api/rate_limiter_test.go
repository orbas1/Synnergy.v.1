package api

import (
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	limiter := NewRateLimiter(2, WithRateWindow(time.Second))
	if !limiter.Allow() {
		t.Fatalf("expected first request to pass")
	}
	if !limiter.Allow() {
		t.Fatalf("expected second request to pass within burst")
	}
	if limiter.Allow() {
		t.Fatalf("expected third request to be limited")
	}
}

func TestRateLimiterAllowForKey(t *testing.T) {
	limiter := NewRateLimiter(1, WithRateWindow(time.Second))
	if _, ok := limiter.AllowFor("client1"); !ok {
		t.Fatalf("expected first request to pass")
	}
	if retry, ok := limiter.AllowFor("client1"); ok || retry <= 0 {
		t.Fatalf("expected rate limit with retry, got %v %v", ok, retry)
	}
	limiter.Reset("client1")
	if _, ok := limiter.AllowFor("client1"); !ok {
		t.Fatalf("expected reset to allow request")
	}
}

func TestRateLimiterCustomClock(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	limiter := NewRateLimiter(1, WithRateWindow(time.Second), WithRateClock(clock))
	if _, ok := limiter.AllowFor("key"); !ok {
		t.Fatalf("expected initial request")
	}
	if _, ok := limiter.AllowFor("key"); ok {
		t.Fatalf("expected rate limit before advance")
	}
	now = now.Add(2 * time.Second)
	if _, ok := limiter.AllowFor("key"); !ok {
		t.Fatalf("expected request after clock advance")
	}
}
