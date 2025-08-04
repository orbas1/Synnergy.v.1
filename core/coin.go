package core

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
// schedule. Rewards drop by half every HalvingInterval blocks.
func BlockReward(height uint64) uint64 {
	halvings := height / HalvingInterval
	reward := InitialBlockReward >> halvings
	return reward
}

