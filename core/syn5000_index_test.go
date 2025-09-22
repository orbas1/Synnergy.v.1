package core

import "testing"

func TestSYN5000IndexRegisterAndSnapshot(t *testing.T) {
	index := NewSYN5000Index()
	token := NewSYN5000Token("Index", "IDX", 0)

	if err := index.Register("IDX", token); err != nil {
		t.Fatalf("register token: %v", err)
	}
	if err := index.Register("IDX", token); err != nil {
		t.Fatalf("register token second time: %v", err)
	}

	if _, ok := index.Snapshot("missing"); ok {
		t.Fatalf("expected snapshot miss for unknown symbol")
	}

	if _, err := token.PlaceBet("alice", 25, 1.6, "dice"); err != nil {
		t.Fatalf("place bet: %v", err)
	}
	if _, err := token.PlaceBet("bob", 40, 1.8, "dice"); err != nil {
		t.Fatalf("place bet: %v", err)
	}

	snap, ok := index.Snapshot("IDX")
	if !ok {
		t.Fatalf("expected snapshot for IDX")
	}
	if snap.Pending != 2 || snap.TotalBets != 2 {
		t.Fatalf("unexpected snapshot %+v", snap)
	}

	bets, ok := index.Bets("IDX", BetFilter{Bettor: "alice"})
	if !ok {
		t.Fatalf("expected bet listing for symbol")
	}
	if len(bets) != 1 || bets[0].Placement.Bettor != "alice" {
		t.Fatalf("unexpected bet listing %+v", bets)
	}

	if err := index.Register("", token); err == nil {
		t.Fatalf("expected error for empty symbol")
	}
	if err := index.Register("NEW", nil); err == nil {
		t.Fatalf("expected error for nil token")
	}
}

func TestSYN5000IndexSymbols(t *testing.T) {
	index := NewSYN5000Index()
	tokens := []string{"ZETA", "ALPHA", "OMEGA"}
	for _, sym := range tokens {
		if err := index.Register(sym, NewSYN5000Token(sym, sym, 0)); err != nil {
			t.Fatalf("register %s: %v", sym, err)
		}
	}
	syms := index.Symbols()
	if len(syms) != len(tokens) {
		t.Fatalf("expected %d symbols got %d", len(tokens), len(syms))
	}
	for i := 1; i < len(syms); i++ {
		if syms[i-1] > syms[i] {
			t.Fatalf("symbols not sorted: %v", syms)
		}
	}
}
