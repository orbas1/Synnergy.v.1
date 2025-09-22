package tokens

import (
	"math/big"
	"testing"
)

func TestSYN3200Conversion(t *testing.T) {
	tok := NewSYN3200Token(1.5)
	if err := tok.SetRatio(2.0); err != nil {
		t.Fatalf("set ratio: %v", err)
	}
	if v := tok.Convert(5); v != 10 {
		t.Fatalf("unexpected conversion result: %d", v)
	}
	exact := tok.ConvertExact(3)
	expected := big.NewRat(6, 1)
	if exact.Cmp(expected) != 0 {
		t.Fatalf("expected exact 6, got %s", exact)
	}
}

func TestSYN3200SetRatioFraction(t *testing.T) {
	tok := NewSYN3200Token(0)
	if err := tok.SetRatioFraction(3, 2); err != nil {
		t.Fatalf("set ratio fraction: %v", err)
	}
	if v := tok.Convert(4); v != 6 {
		t.Fatalf("expected rounded result 6, got %d", v)
	}
	ratio := tok.Ratio()
	if ratio.Cmp(big.NewRat(3, 2)) != 0 {
		t.Fatalf("unexpected ratio %s", ratio)
	}
	if err := tok.SetRatioFraction(1, 0); err != ErrInvalidRatio {
		t.Fatalf("expected ErrInvalidRatio, got %v", err)
	}
}
