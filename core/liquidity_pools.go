package core

import (
	"errors"
	"fmt"
	"sync"
)

// LiquidityPool represents a constant-product AMM pool.
type LiquidityPool struct {
	mu         sync.Mutex
	ID         string
	TokenA     string
	TokenB     string
	ReserveA   uint64
	ReserveB   uint64
	FeeBps     uint16
	LPBalances map[string]uint64
	totalLP    uint64
}

// NewLiquidityPool creates a new liquidity pool.
func NewLiquidityPool(id, tokenA, tokenB string, feeBps uint16) *LiquidityPool {
	return &LiquidityPool{
		ID:         id,
		TokenA:     tokenA,
		TokenB:     tokenB,
		FeeBps:     feeBps,
		LPBalances: make(map[string]uint64),
	}
}

// AddLiquidity adds tokens to the pool and mints LP tokens.
func (p *LiquidityPool) AddLiquidity(provider string, amtA, amtB uint64) (uint64, error) {
	if provider == "" || amtA == 0 || amtB == 0 {
		return 0, errors.New("invalid liquidity parameters")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.totalLP == 0 {
		p.ReserveA += amtA
		p.ReserveB += amtB
		lp := sqrt(amtA * amtB)
		p.totalLP = lp
		p.LPBalances[provider] += lp
		return lp, nil
	}
	if p.ReserveA == 0 || p.ReserveB == 0 {
		return 0, errors.New("invalid pool state")
	}
	requiredB := amtA * p.ReserveB / p.ReserveA
	if amtB < requiredB {
		return 0, fmt.Errorf("insufficient amount of token B: need %d", requiredB)
	}
	lp := amtA * p.totalLP / p.ReserveA
	p.ReserveA += amtA
	p.ReserveB += requiredB
	p.totalLP += lp
	p.LPBalances[provider] += lp
	return lp, nil
}

// RemoveLiquidity burns LP tokens and returns underlying assets.
func (p *LiquidityPool) RemoveLiquidity(provider string, lpTokens uint64) (uint64, uint64, error) {
	if provider == "" || lpTokens == 0 {
		return 0, 0, errors.New("invalid liquidity parameters")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	bal := p.LPBalances[provider]
	if bal < lpTokens || p.totalLP == 0 {
		return 0, 0, errors.New("insufficient LP balance")
	}
	amtA := lpTokens * p.ReserveA / p.totalLP
	amtB := lpTokens * p.ReserveB / p.totalLP
	p.ReserveA -= amtA
	p.ReserveB -= amtB
	p.totalLP -= lpTokens
	p.LPBalances[provider] -= lpTokens
	return amtA, amtB, nil
}

// Swap executes a token swap within the pool.
func (p *LiquidityPool) Swap(tokenIn string, amtIn, minOut uint64) (uint64, error) {
	if amtIn == 0 || minOut == 0 {
		return 0, errors.New("amounts must be > 0")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	var reserveIn, reserveOut *uint64
	switch tokenIn {
	case p.TokenA:
		reserveIn = &p.ReserveA
		reserveOut = &p.ReserveB
	case p.TokenB:
		reserveIn = &p.ReserveB
		reserveOut = &p.ReserveA
	default:
		return 0, fmt.Errorf("token %s not in pool", tokenIn)
	}
	fee := amtIn * uint64(p.FeeBps) / 10000
	amtInWithFee := amtIn - fee
	newIn := *reserveIn + amtInWithFee
	k := uint64(*reserveIn * *reserveOut)
	newOut := k / newIn
	out := *reserveOut - newOut
	if out < minOut {
		return 0, errors.New("insufficient output amount")
	}
	*reserveIn = newIn
	*reserveOut = newOut
	return out, nil
}

// sqrt returns floor square root of n using simple integer method.
func sqrt(n uint64) uint64 {
	var x, y uint64 = n, (n + 1) / 2
	for y < x {
		x = y
		y = (n/y + y) / 2
	}
	return x
}

// LiquidityPoolRegistry tracks all pools.
type LiquidityPoolRegistry struct {
	mu    sync.RWMutex
	pools map[string]*LiquidityPool
}

// NewLiquidityPoolRegistry creates an empty registry.
func NewLiquidityPoolRegistry() *LiquidityPoolRegistry {
	return &LiquidityPoolRegistry{pools: make(map[string]*LiquidityPool)}
}

// Create registers a new liquidity pool.
func (r *LiquidityPoolRegistry) Create(id, tokenA, tokenB string, feeBps uint16) (*LiquidityPool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if id == "" || tokenA == "" || tokenB == "" || tokenA == tokenB {
		return nil, errors.New("invalid pool parameters")
	}
	if _, exists := r.pools[id]; exists {
		return nil, fmt.Errorf("pool %s exists", id)
	}
	p := NewLiquidityPool(id, tokenA, tokenB, feeBps)
	r.pools[id] = p
	return p, nil
}

// Get returns a pool by ID.
func (r *LiquidityPoolRegistry) Get(id string) (*LiquidityPool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pools[id]
	return p, ok
}

// List returns all pools.
func (r *LiquidityPoolRegistry) List() []*LiquidityPool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*LiquidityPool, 0, len(r.pools))
	for _, p := range r.pools {
		res = append(res, p)
	}
	return res
}
