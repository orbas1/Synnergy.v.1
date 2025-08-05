package core

import "testing"

// helper to create a pool with initial liquidity
func newTestLiquidityPool(t *testing.T, id, tokenA, tokenB string, fee uint16, amtA, amtB uint64) *LiquidityPool {
    t.Helper()
    p := NewLiquidityPool(id, tokenA, tokenB, fee)
    if _, err := p.AddLiquidity("provider", amtA, amtB); err != nil {
        t.Fatalf("add liquidity: %v", err)
    }
    return p
}

func TestNewLiquidityPoolViewSnapshot(t *testing.T) {
    p := newTestLiquidityPool(t, "pool1", "A", "B", 30, 500, 1000)
    view := NewLiquidityPoolView(p)

    if view.ID != p.ID || view.TokenA != p.TokenA || view.TokenB != p.TokenB ||
        view.ReserveA != p.ReserveA || view.ReserveB != p.ReserveB || view.FeeBps != p.FeeBps {
        t.Fatalf("view does not match pool: %#v vs %#v", view, p)
    }

    // mutate pool and ensure the view remains unchanged (snapshot)
    p.ReserveA = 9999
    if view.ReserveA == p.ReserveA {
        t.Fatalf("expected view to be snapshot, got ReserveA %d", view.ReserveA)
    }
}

func TestPoolInfo(t *testing.T) {
    reg := NewLiquidityPoolRegistry()
    p := newTestLiquidityPool(t, "pool1", "A", "B", 30, 100, 200)
    reg.pools[p.ID] = p

    view, ok := reg.PoolInfo("pool1")
    if !ok {
        t.Fatal("expected pool info to be found")
    }
    if view.ID != "pool1" || view.ReserveA != 100 || view.ReserveB != 200 {
        t.Fatalf("unexpected view: %#v", view)
    }

    if _, ok := reg.PoolInfo("missing"); ok {
        t.Fatal("expected missing pool to return false")
    }
}

func TestPoolViews(t *testing.T) {
    reg := NewLiquidityPoolRegistry()
    p1 := newTestLiquidityPool(t, "p1", "A", "B", 30, 1000, 1000)
    p2 := newTestLiquidityPool(t, "p2", "C", "D", 25, 500, 400)
    reg.pools[p1.ID] = p1
    reg.pools[p2.ID] = p2

    views := reg.PoolViews()
    if len(views) != 2 {
        t.Fatalf("expected 2 views, got %d", len(views))
    }
    m := make(map[string]LiquidityPoolView)
    for _, v := range views {
        m[v.ID] = v
    }

    v1, ok := m["p1"]
    if !ok || v1.ReserveA != 1000 || v1.ReserveB != 1000 {
        t.Fatalf("unexpected view for p1: %#v", v1)
    }
    v2, ok := m["p2"]
    if !ok || v2.ReserveA != 500 || v2.ReserveB != 400 || v2.FeeBps != 25 {
        t.Fatalf("unexpected view for p2: %#v", v2)
    }
}

