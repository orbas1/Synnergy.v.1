package core

import "testing"

func TestSYN3500TokenLifecycle(t *testing.T) {
	tok := NewSYN3500Token("Dollar", "USD", "bank", 1.0)
	tok.Mint("alice", 100)
	if bal := tok.BalanceOf("alice"); bal != 100 {
		t.Fatalf("expected 100 got %d", bal)
	}
	if err := tok.Redeem("alice", 40); err != nil {
		t.Fatalf("redeem: %v", err)
	}
	if bal := tok.BalanceOf("alice"); bal != 60 {
		t.Fatalf("expected 60 got %d", bal)
	}
	if err := tok.Redeem("alice", 1000); err == nil {
		t.Fatalf("expected error on excessive redeem")
	}
	tok.SetRate(1.1)
	_, _, rate := tok.Info()
	if rate != 1.1 {
		t.Fatalf("expected rate 1.1 got %f", rate)
	}
}
