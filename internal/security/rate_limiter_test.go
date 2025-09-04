package security

import (
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	r := NewRateLimiter(time.Millisecond)
	if !r.Allow() {
		t.Fatal("expected first event to be allowed")
	}
	if r.Allow() {
		t.Fatal("expected second event to be blocked")
	}
	time.Sleep(time.Millisecond)
	if !r.Allow() {
		t.Fatal("expected event after interval to be allowed")
	}
}
