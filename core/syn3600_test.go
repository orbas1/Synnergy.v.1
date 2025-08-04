package core

import (
	"testing"
	"time"
)

func TestFuturesContract(t *testing.T) {
	exp := time.Unix(100, 0)
	f := NewFuturesContract("BTC", 2, 1000, exp)
	pnl := f.Settle(1100)
	if pnl != 200 {
		t.Fatalf("expected pnl 200, got %d", pnl)
	}
	if !f.Settled {
		t.Fatalf("contract should be settled")
	}
	if !f.IsExpired(time.Unix(200, 0)) {
		t.Fatalf("expected expired")
	}
}
