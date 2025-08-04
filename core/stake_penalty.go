package core

// StakePenaltyManager adjusts stakes for validators based on behaviour.
type StakePenaltyManager struct{}

// NewStakePenaltyManager creates a new StakePenaltyManager.
func NewStakePenaltyManager() *StakePenaltyManager { return &StakePenaltyManager{} }

// Slash reduces the stake of addr on sn by the given penalty amount.
func (spm *StakePenaltyManager) Slash(sn *StakingNode, addr string, penalty uint64) {
	sn.Unstake(addr, penalty)
}

// Reward increases the stake of addr on sn by the given amount.
func (spm *StakePenaltyManager) Reward(sn *StakingNode, addr string, reward uint64) {
	sn.Stake(addr, reward)
}
