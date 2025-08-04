package core

import (
	"testing"
	"time"
)

func TestTokenInsurancePolicy(t *testing.T) {
	start := time.Unix(0, 0)
	end := start.Add(24 * time.Hour)
	policy := NewTokenInsurancePolicy("p1", "alice", "coverage", 10, 1000, 0, 1000, start, end)
	now := start.Add(12 * time.Hour)
	if !policy.IsActive(now) {
		t.Fatalf("policy should be active")
	}
	payout, err := policy.Claim(now)
	if err != nil || payout != 1000 {
		t.Fatalf("claim failed: %v %d", err, payout)
	}
	if policy.IsActive(now) {
		t.Fatalf("policy should be inactive after claim")
	}
}
