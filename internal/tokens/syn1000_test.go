package tokens

import (
	"math/big"
	"sync"
	"testing"
)

// TestSYN1000TokenReserveValue verifies that reserve additions and price updates
// are reflected in the calculated total value using high precision arithmetic.
func TestSYN1000TokenReserveValue(t *testing.T) {
	tok := NewSYN1000Token(1, "Stable", "STB", 2)
	tok.AddReserve("USD", big.NewRat(100, 1))
	tok.SetReservePrice("USD", big.NewRat(1, 1))
	tok.AddReserve("EUR", big.NewRat(50, 1))
	tok.SetReservePrice("EUR", big.NewRat(2, 1))

	got := tok.TotalReserveValue()
	want := big.NewRat(200, 1) // 100*1 + 50*2
	if got.Cmp(want) != 0 {
		t.Fatalf("want %s got %s", want.String(), got.String())
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
			tok.AddReserve("USD", big.NewRat(1, 1))
		}()
	}
	wg.Wait()
	tok.SetReservePrice("USD", big.NewRat(1, 1))
	got := tok.TotalReserveValue()
	want := big.NewRat(10, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("want %s got %s", want.String(), got.String())
	}
}
