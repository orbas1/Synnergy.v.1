package tokens

import (
	"math/big"
	"sync"
	"testing"
)

// TestSYN1000Index verifies basic index operations and value calculation.
func TestSYN1000Index(t *testing.T) {
	idx := NewSYN1000Index()
	var mu sync.Mutex
	events := make([]IndexEvent, 0, 3)
	idx.RegisterWatcher(func(evt IndexEvent) {
		mu.Lock()
		events = append(events, evt)
		mu.Unlock()
	})
	id := idx.Create("Stable", "STB", 2)

	if err := idx.AddReserve(id, "USD", big.NewRat(100, 1)); err != nil {
		t.Fatalf("add reserve: %v", err)
	}
	if err := idx.SetReservePrice(id, "USD", big.NewRat(1, 1)); err != nil {
		t.Fatalf("set price: %v", err)
	}
	if err := idx.RemoveReserve(id, "missing"); err != nil {
		t.Fatalf("remove reserve: %v", err)
	}

	v, err := idx.TotalValue(id)
	if err != nil {
		t.Fatalf("total value: %v", err)
	}
	want := big.NewRat(100, 1)
	if v.Cmp(want) != 0 {
		t.Fatalf("want %s got %s", want.String(), v.String())
	}
	ratio, err := idx.Collateralization(id)
	if err != nil {
		t.Fatalf("collateralization: %v", err)
	}
	if ratio.Sign() != 0 {
		t.Fatalf("expected zero ratio with no supply got %s", ratio.String())
	}
	mu.Lock()
	if len(events) < 2 || events[0].Type != IndexEventCreated {
		t.Fatalf("unexpected events %+v", events)
	}
	mu.Unlock()
}
