package core

import (
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
