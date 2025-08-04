package core

import "testing"

func TestSYN5000Token(t *testing.T) {
	token := NewSYN5000Token("Gamble", "GMB", 0)
	id := token.PlaceBet("alice", 100, 2.0, "dice")
	payout, err := token.ResolveBet(id, true)
	if err != nil || payout != 200 {
		t.Fatalf("expected payout 200 got %d err %v", payout, err)
	}
	if b, ok := token.GetBet(id); !ok || !b.Resolved || !b.Won {
		t.Fatalf("bet not resolved correctly")
	}
}
