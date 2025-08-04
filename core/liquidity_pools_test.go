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
	if len(reg.PoolViews()) != 1 {
		t.Fatal("expected one pool view")
	}
}
