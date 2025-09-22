package tokens

import (
	"math/big"
	"sync"
	"testing"
)

// TestSYN1000TokenReserveValue verifies reserve accounting and collateralisation calculations.
func TestSYN1000TokenReserveValue(t *testing.T) {
	tok := NewSYN1000Token(1, "Stable", "STB", 2)
	if err := tok.AddReserve("USD", big.NewRat(100, 1)); err != nil {
		t.Fatalf("add reserve: %v", err)
	}
	if err := tok.SetReservePrice("USD", big.NewRat(1, 1)); err != nil {
		t.Fatalf("set price: %v", err)
	}
	if err := tok.AddReserve("EUR", big.NewRat(50, 1)); err != nil {
		t.Fatalf("add reserve: %v", err)
	}
	if err := tok.SetReservePrice("EUR", big.NewRat(2, 1)); err != nil {
		t.Fatalf("set price: %v", err)
	}
	if err := tok.Mint("alice", 100); err != nil {
		t.Fatalf("mint: %v", err)
	}

	got := tok.TotalReserveValue()
	want := big.NewRat(200, 1) // 100*1 + 50*2
	if got.Cmp(want) != 0 {
		t.Fatalf("want %s got %s", want.String(), got.String())
	}
	ratio := tok.CollateralizationRatio()
	expectedRatio := new(big.Rat).Quo(want, big.NewRat(100, 1))
	if ratio.Cmp(expectedRatio) != 0 {
		t.Fatalf("unexpected ratio %s want %s", ratio.String(), expectedRatio.String())
	}
	breakdown := tok.ReserveBreakdown()
	if len(breakdown) != 2 || breakdown["USD"].Amount.Cmp(big.NewRat(100, 1)) != 0 {
		t.Fatalf("unexpected breakdown %+v", breakdown)
	}
}

// TestSYN1000TokenConcurrent ensures concurrent reserve updates are safe.
func TestSYN1000TokenConcurrent(t *testing.T) {
	tok := NewSYN1000Token(1, "Stable", "STB", 2)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := tok.AddReserve("USD", big.NewRat(1, 1)); err != nil {
				t.Errorf("add reserve: %v", err)
			}
		}()
	}
	wg.Wait()
	if err := tok.SetReservePrice("USD", big.NewRat(1, 1)); err != nil {
		t.Fatalf("set price: %v", err)
	}
	got := tok.TotalReserveValue()
	want := big.NewRat(10, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("want %s got %s", want.String(), got.String())
	}
}
