package tokens

import "testing"

func TestSYN3800TokenLifecycle(t *testing.T) {
	tok := NewSYN3800Token(100)
	if err := tok.Mint(0); err != ErrInvalidAmount {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
	if err := tok.Mint(80); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if tok.Supply() != 80 || tok.Remaining() != 20 {
		t.Fatalf("unexpected state supply %d remaining %d", tok.Supply(), tok.Remaining())
	}
	if err := tok.Mint(30); err != ErrCapExceeded {
		t.Fatalf("expected ErrCapExceeded, got %v", err)
	}
	if err := tok.SetCap(60); err != ErrCapExceeded {
		t.Fatalf("expected ErrCapExceeded on lowering cap, got %v", err)
	}
	if err := tok.SetCap(120); err != nil {
		t.Fatalf("set cap: %v", err)
	}
	if err := tok.Burn(90); err != ErrInsufficientSupply {
		t.Fatalf("expected ErrInsufficientSupply, got %v", err)
	}
	if err := tok.Burn(50); err != nil {
		t.Fatalf("burn: %v", err)
	}
	if tok.Supply() != 30 {
		t.Fatalf("unexpected supply %d", tok.Supply())
	}
}
