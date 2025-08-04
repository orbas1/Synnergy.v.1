package core

// LiquidityPoolView provides a serialisable snapshot of a pool's state.
type LiquidityPoolView struct {
	ID       string `json:"id"`
	TokenA   string `json:"token_a"`
	TokenB   string `json:"token_b"`
	ReserveA uint64 `json:"reserve_a"`
	ReserveB uint64 `json:"reserve_b"`
	FeeBps   uint16 `json:"fee_bps"`
}

// NewLiquidityPoolView constructs a view from a pool.
func NewLiquidityPoolView(p *LiquidityPool) LiquidityPoolView {
	return LiquidityPoolView{
		ID:       p.ID,
		TokenA:   p.TokenA,
		TokenB:   p.TokenB,
		ReserveA: p.ReserveA,
		ReserveB: p.ReserveB,
		FeeBps:   p.FeeBps,
	}
}

// PoolInfo returns a view for a pool by ID if it exists.
func (r *LiquidityPoolRegistry) PoolInfo(id string) (LiquidityPoolView, bool) {
	p, ok := r.Get(id)
	if !ok {
		return LiquidityPoolView{}, false
	}
	return NewLiquidityPoolView(p), true
}

// PoolViews lists all pools in view form.
func (r *LiquidityPoolRegistry) PoolViews() []LiquidityPoolView {
	pools := r.List()
	views := make([]LiquidityPoolView, 0, len(pools))
	for _, p := range pools {
		views = append(views, NewLiquidityPoolView(p))
	}
	return views
}
