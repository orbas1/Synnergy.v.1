package core

import "math"

const (
	// CoinName represents the name of the native coin.
	CoinName = "Synthron"
	// MaxSupply defines the capped supply of Synthron coins.
	MaxSupply uint64 = 500_000_000
	// GenesisAllocation is the initial amount issued to the creators' wallet.
	GenesisAllocation uint64 = 5_000_000
	// InitialBlockReward is the starting reward for mining a main block.
	InitialBlockReward uint64 = 1252
	// HalvingInterval is the number of blocks between reward halvings.
	HalvingInterval uint64 = 200_000
)

// BlockReward returns the reward for a given block height following the halving
// schedule. Rewards drop by half every HalvingInterval blocks. A reward of zero
// is returned once all halvings have completed.
func BlockReward(height uint64) uint64 {
	halvings := height / HalvingInterval
	reward := InitialBlockReward >> halvings
	return reward
}

// CirculatingSupply approximates the number of Synthron coins in circulation at
// the given block height, taking the halving schedule into account.
func CirculatingSupply(height uint64) uint64 {
	supply := GenesisAllocation
	for i := uint64(0); i < height; i++ {
		supply += BlockReward(i)
		if supply >= MaxSupply {
			return MaxSupply
		}
	}
	return supply
}

// RemainingSupply returns the number of coins left to be minted after the given
// block height.
func RemainingSupply(height uint64) uint64 {
	circ := CirculatingSupply(height)
	if circ >= MaxSupply {
		return 0
	}
	return MaxSupply - circ
}

// InitialPrice calculates the initial price of Synthron according to the
// economic model defined in the specification.
func InitialPrice(C, R, M, V, T, E float64) float64 {
	return (C + R + (M*V)/T) * E
}

// AlphaFactor computes the dynamic alpha scaling factor for staking
// requirements.
func AlphaFactor(volatility, participation, economicStability, normalization float64) float64 {
	return (3*volatility + participation + economicStability) / normalization
}

// MinimumStake calculates the minimum stake required to become a validator using
// network transaction counts, the current mining reward, circulating supply and
// an alpha scaling factor.
func MinimumStake(totalTx, currentReward, circulatingSupply, alpha float64) float64 {
	if currentReward == 0 || circulatingSupply == 0 {
		return 0
	}
	return (totalTx / (currentReward * circulatingSupply)) * alpha
}

// LockupDuration derives the staking lock-up period based on normalized
// transaction volume V and volatility index sigma.
func LockupDuration(base, V, threshold, sigma float64) float64 {
	if threshold == 0 {
		return base
	}
	return base*(V/threshold*10) + (sigma * 20)
}

// PriceToSupplyRatio is a helper exposed for potential economic simulations. It
// returns price divided by circulating supply for basic analytics.
func PriceToSupplyRatio(price float64, height uint64) float64 {
	return price / math.Max(float64(CirculatingSupply(height)), 1)
}
