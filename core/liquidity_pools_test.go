package core

import "testing"

func TestLiquidityPoolLifecycle(t *testing.T) {
	reg := NewLiquidityPoolRegistry()
	p, err := reg.Create("pool1", "A", "B", 30)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err = p.AddLiquidity("prov", 1000, 1000); err != nil {
		t.Fatalf("add: %v", err)
	}
	if _, err = p.Swap("A", 100, 1); err != nil {
		t.Fatalf("swap: %v", err)
	}
	if _, _, err = p.RemoveLiquidity("prov", 10); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if len(reg.List()) != 1 {
		t.Fatal("expected one pool view")
	}
}

func TestLiquidityPoolValidation(t *testing.T) {
	reg := NewLiquidityPoolRegistry()
	if _, err := reg.Create("", "A", "B", 30); err == nil {
		t.Fatal("expected error for empty id")
	}
	p, err := reg.Create("p1", "A", "B", 30)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := p.AddLiquidity("", 0, 0); err == nil {
		t.Fatal("expected validation error")
	}
	if _, err := p.Swap("X", 10, 1); err == nil {
		t.Fatal("expected unknown token error")
	}
	if _, _, err := p.RemoveLiquidity("prov", 0); err == nil {
		t.Fatal("expected invalid lp token error")
	}
}
