package core

import (
	"testing"

	"synnergy/internal/tokens"
)

func TestCentralBankingMintCBDC(t *testing.T) {
	l := NewLedger()
	tok := tokens.NewSYN10Token(1, "CBDC", "cSYN", "central", 1, 2)
	cb := NewCentralBankingNode("id", "addr", l, "neutral", tok)

	if err := cb.MintCBDC("alice", 10); err != nil {
		t.Fatalf("mintCBDC: %v", err)
	}
	if bal := tok.BalanceOf("alice"); bal != 10 {
		t.Fatalf("unexpected token balance %d", bal)
	}
	if bal := l.GetBalance("alice"); bal != 0 {
		t.Fatalf("ledger should remain unchanged, got %d", bal)
	}
	if err := cb.MintCBDC("alice", 0); err == nil {
		t.Fatalf("expected error on zero mint")
	}
}
