package core

import "errors"

// GenesisDexSeed defines the liquidity bootstrap configuration for a DEX pool.
type GenesisDexSeed struct {
	PoolID   string
	TokenA   string
	TokenB   string
	FeeBps   uint16
	ReserveA uint64
	ReserveB uint64
	Provider string
}

var defaultGenesisDexSeed = GenesisDexSeed{
	PoolID:   "SYNN-BTC",
	TokenA:   "SYNN",
	TokenB:   "BTC",
	FeeBps:   30,
	ReserveA: 100,
	ReserveB: 1,
	Provider: "genesis-liquidity",
}

// EnsureGenesisDex initialises the canonical SYNN/BTC pool with a
// 100 SYNN : 1 BTC price ratio at genesis. It is safe to call multiple times;
// if liquidity already exists it will not be modified.
func EnsureGenesisDex(reg *LiquidityPoolRegistry) error {
	return SeedGenesisDexLiquidity(reg, defaultGenesisDexSeed)
}

// SeedGenesisDexLiquidity applies a liquidity seed to the registry, creating the
// pool if it does not exist and adding liquidity only when the pool is empty.
func SeedGenesisDexLiquidity(reg *LiquidityPoolRegistry, seed GenesisDexSeed) error {
	if reg == nil {
		return errors.New("liquidity registry required")
	}
	if seed.PoolID == "" || seed.TokenA == "" || seed.TokenB == "" || seed.Provider == "" {
		return errors.New("invalid seed configuration")
	}
	if seed.ReserveA == 0 || seed.ReserveB == 0 {
		return errors.New("seed reserves must be greater than zero")
	}
	pool, ok := reg.Get(seed.PoolID)
	if !ok {
		var err error
		pool, err = reg.Create(seed.PoolID, seed.TokenA, seed.TokenB, seed.FeeBps)
		if err != nil {
			return err
		}
	}
	if pool.totalLP != 0 || pool.ReserveA != 0 || pool.ReserveB != 0 {
		return nil
	}
	_, err := pool.AddLiquidity(seed.Provider, seed.ReserveA, seed.ReserveB)
	return err
}
