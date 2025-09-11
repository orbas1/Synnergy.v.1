package core

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestInvestmentRedeem(t *testing.T) {
	reg := NewInvestmentRegistry()
	rec, err := reg.Issue("inv1", "alice", 1000, 0.10, time.Now().Add(365*24*time.Hour))
	if err != nil {
		t.Fatalf("issue failed: %v", err)
	}
	maturity := rec.LastAccrued.Add(365 * 24 * time.Hour)
	payout, err := reg.Redeem("inv1", "alice", maturity)
	if err != nil {
		t.Fatalf("redeem failed: %v", err)
	}
	if payout != 1100 {
		t.Fatalf("expected payout 1100 got %d", payout)
	}
	if _, ok := reg.Get("inv1"); ok {
		t.Fatalf("record should be removed after redemption")
	}
}

func TestInvestmentDuplicateIssue(t *testing.T) {
	reg := NewInvestmentRegistry()
	if _, err := reg.Issue("inv1", "alice", 1000, 0.10, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("first issue failed: %v", err)
	}
	if _, err := reg.Issue("inv1", "alice", 1000, 0.10, time.Now().Add(time.Hour)); err != ErrInvestmentExists {
		t.Fatalf("expected ErrInvestmentExists got %v", err)
	}
}

func TestInvestmentConcurrentIssue(t *testing.T) {
	reg := NewInvestmentRegistry()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("inv%d", i)
			if _, err := reg.Issue(id, "owner", 100, 0.05, time.Now().Add(time.Hour)); err != nil {
				t.Errorf("issue %s failed: %v", id, err)
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < 50; i++ {
		if _, ok := reg.Get(fmt.Sprintf("inv%d", i)); !ok {
			t.Fatalf("missing record inv%d", i)
		}
	}
}

func TestInvestmentUnauthorizedRedeem(t *testing.T) {
	reg := NewInvestmentRegistry()
	_, err := reg.Issue("inv1", "alice", 1000, 0.10, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("issue failed: %v", err)
	}
	if _, err := reg.Redeem("inv1", "bob", time.Now().Add(2*time.Hour)); err != ErrUnauthorizedRedeemer {
		t.Fatalf("expected ErrUnauthorizedRedeemer got %v", err)
	}
}
