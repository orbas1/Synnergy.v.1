package core

import "testing"

func TestEnsureGenesisDexSeedsPool(t *testing.T) {
	reg := NewLiquidityPoolRegistry()
	if err := EnsureGenesisDex(reg); err != nil {
		t.Fatalf("ensure genesis dex: %v", err)
	}
	pool, ok := reg.Get("SYNN-BTC")
	if !ok {
		t.Fatalf("expected SYNN-BTC pool")
	}
	if pool.ReserveA != 100 || pool.ReserveB != 1 {
		t.Fatalf("unexpected reserves: %d:%d", pool.ReserveA, pool.ReserveB)
	}
	if lp := pool.LPBalances["genesis-liquidity"]; lp == 0 {
		t.Fatalf("expected genesis LP tokens, got %d", lp)
	}
}

func TestEnsureGenesisDexDoesNotOverwrite(t *testing.T) {
	reg := NewLiquidityPoolRegistry()
	pool, err := reg.Create("SYNN-BTC", "SYNN", "BTC", 25)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if _, err := pool.AddLiquidity("alice", 500, 5); err != nil {
		t.Fatalf("add liquidity: %v", err)
	}
	if err := EnsureGenesisDex(reg); err != nil {
		t.Fatalf("ensure genesis dex: %v", err)
	}
	if pool.ReserveA != 500 || pool.ReserveB != 5 {
		t.Fatalf("reserves modified: %d:%d", pool.ReserveA, pool.ReserveB)
	}
	if lp := pool.LPBalances["genesis-liquidity"]; lp != 0 {
		t.Fatalf("unexpected genesis liquidity injection")
	}
}
