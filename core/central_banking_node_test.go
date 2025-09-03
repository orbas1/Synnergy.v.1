package core

import "testing"

func TestCentralBankingMint(t *testing.T) {
	l := NewLedger()
	cb := NewCentralBankingNode("id", "addr", l, "neutral")

	if err := cb.Mint("alice", 10); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if bal := l.GetBalance("alice"); bal != 10 {
		t.Fatalf("unexpected balance %d", bal)
	}
	// attempt to mint beyond remaining supply
	if err := cb.Mint("alice", RemainingSupply(0)+1); err == nil {
		t.Fatalf("expected supply error")
	}
}
