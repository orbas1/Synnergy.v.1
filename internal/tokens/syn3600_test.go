package tokens

import "testing"

func TestSYN3600Weights(t *testing.T) {
	token := NewSYN3600Token()
	token.SetWeight("alice", 10)
	token.IncreaseWeight("bob", 20)
	if err := token.DecreaseWeight("bob", 5); err != nil {
		t.Fatalf("decrease weight: %v", err)
	}
	if err := token.DecreaseWeight("alice", 20); err != ErrWeightTooLow {
		t.Fatalf("expected ErrWeightTooLow, got %v", err)
	}
	if token.Weight("bob") != 15 {
		t.Fatalf("unexpected bob weight: %d", token.Weight("bob"))
	}
	if token.TotalWeight() != 25 {
		t.Fatalf("unexpected total weight: %d", token.TotalWeight())
	}
	top := token.TopHolders(1)
	if len(top) != 1 || top[0] != "bob" {
		t.Fatalf("unexpected top holder: %v", top)
	}
}
