package core

import (
	"sync"
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
	if _, err := policy.Claim(now); err != ErrPolicyInactive {
		t.Fatalf("expected ErrPolicyInactive, got %v", err)
	}
}

func TestTokenInsurancePolicyConcurrentClaim(t *testing.T) {
	start := time.Unix(0, 0)
	end := start.Add(time.Hour)
	policy := NewTokenInsurancePolicy("p1", "alice", "cov", 1, 100, 0, 100, start, end)
	now := start.Add(30 * time.Minute)
	var wg sync.WaitGroup
	var payouts []uint64
	var mu sync.Mutex
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if p, err := policy.Claim(now); err == nil {
				mu.Lock()
				payouts = append(payouts, p)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if len(payouts) != 1 || payouts[0] != 100 {
		t.Fatalf("expected single payout of 100, got %v", payouts)
	}
}
