package core

import "testing"

func TestDAOTokenLedger(t *testing.T) {
	l := NewDAOTokenLedger()
	l.Mint("a", 10)
	if l.Balance("a") != 10 {
		t.Fatalf("expected 10")
	}
	if err := l.Transfer("a", "b", 5); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if l.Balance("a") != 5 || l.Balance("b") != 5 {
		t.Fatalf("unexpected balances")
	}
}
