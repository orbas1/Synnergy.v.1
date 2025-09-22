package tokens

import (
	"testing"
)

func TestForexRegistryLifecycle(t *testing.T) {
	reg := NewForexRegistry()
	pair, err := reg.Register("USD", "EUR", 0.9)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := reg.Register("", "EUR", 1); err == nil {
		t.Fatalf("expected validation error for empty base")
	}
	fetched, err := reg.Get(pair.PairID)
	if err != nil {
		t.Fatalf("get pair: %v", err)
	}
	if fetched.Rate != 0.9 {
		t.Fatalf("unexpected rate: %f", fetched.Rate)
	}
	if err := reg.UpdateRate(pair.PairID, 0); err != ErrInvalidForexRate {
		t.Fatalf("expected invalid rate error, got %v", err)
	}
	if err := reg.UpdateRate(pair.PairID, 0.95); err != nil {
		t.Fatalf("update rate: %v", err)
	}
	symbolPair, err := reg.GetBySymbol("USD", "EUR")
	if err != nil {
		t.Fatalf("get by symbol: %v", err)
	}
	if symbolPair.Rate != 0.95 {
		t.Fatalf("unexpected symbol rate: %f", symbolPair.Rate)
	}
	list := reg.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
	if err := reg.Remove(pair.PairID); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, err := reg.Get(pair.PairID); err != ErrForexPairNotFound {
		t.Fatalf("expected pair not found, got %v", err)
	}
}
