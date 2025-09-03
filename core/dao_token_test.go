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
	if err := l.Burn("b", 3); err != nil {
		t.Fatalf("burn: %v", err)
	}
	if l.Balance("b") != 2 {
		t.Fatalf("unexpected burn balance")
	}
	if err := l.Burn("b", 5); err == nil {
		t.Fatalf("expected burn error")
	}
}
