package core

import "testing"

func TestDAOStaking(t *testing.T) {
	s := NewDAOStaking()
	s.Stake("a", 10)
	if s.Balance("a") != 10 {
		t.Fatalf("expected 10")
	}
	if err := s.Unstake("a", 5); err != nil {
		t.Fatalf("unstake: %v", err)
	}
	if s.Balance("a") != 5 {
		t.Fatalf("expected 5")
	}
}
