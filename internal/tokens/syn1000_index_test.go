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

	tok, err := idx.Token(id)
	if err != nil {
		t.Fatalf("token lookup: %v", err)
	}
	tok.Mint("treasury", 90)

	v, err := idx.TotalValue(id)
	if err != nil {
		t.Fatalf("total value: %v", err)
	}
	want := big.NewRat(100, 1)
	if v.Cmp(want) != 0 {
		t.Fatalf("want %s got %s", want.String(), v.String())
	}

	cov, err := idx.CoverageRatio(id)
	if err != nil {
		t.Fatalf("coverage: %v", err)
	}
	if cov.Cmp(big.NewRat(100, 90)) != 0 {
		t.Fatalf("unexpected coverage ratio: %s", cov)
	}

	ok, err := idx.StressTest(id, 80)
	if err != nil {
		t.Fatalf("stress test: %v", err)
	}
	if !ok {
		t.Fatal("expected stress test to succeed")
	}

	snapshot, err := idx.ReserveSnapshot(id)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(snapshot) != 1 || snapshot[0].Asset != "USD" {
		t.Fatalf("unexpected snapshot: %+v", snapshot)
	}
}
