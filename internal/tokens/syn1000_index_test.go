package tokens

import (
	"math/big"
	"testing"
)

// TestSYN1000Index verifies basic index operations and value calculation.
func TestSYN1000Index(t *testing.T) {
	idx := NewSYN1000Index()
	id := idx.Create("Stable", "STB", 2)

	if err := idx.AddReserve(id, "USD", big.NewRat(100, 1)); err != nil {
		t.Fatalf("add reserve: %v", err)
	}
	if err := idx.SetReservePrice(id, "USD", big.NewRat(1, 1)); err != nil {
		t.Fatalf("set price: %v", err)
	}

	v, err := idx.TotalValue(id)
	if err != nil {
		t.Fatalf("total value: %v", err)
	}
	want := big.NewRat(100, 1)
	if v.Cmp(want) != 0 {
		t.Fatalf("want %s got %s", want.String(), v.String())
	}
}
