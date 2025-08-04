package synnergy

import "sync"

// StakingNode tracks token stakes from multiple participants.  It provides
// methods for staking, unstaking and querying balances, enabling the CLI staking
// commands to operate without needing a full ledger implementation.
type StakingNode struct {
	mu    sync.RWMutex
	stake map[string]uint64
	total uint64
}

// NewStakingNode creates an empty staking node instance.
func NewStakingNode() *StakingNode {
	return &StakingNode{stake: make(map[string]uint64)}
}

// Stake locks the specified amount for the given address.
func (s *StakingNode) Stake(addr string, amt uint64) {
	s.mu.Lock()
	s.stake[addr] += amt
	s.total += amt
	s.mu.Unlock()
}

// Unstake releases previously staked tokens for the address.
func (s *StakingNode) Unstake(addr string, amt uint64) {
	s.mu.Lock()
	current := s.stake[addr]
	if amt > current {
		amt = current
	}
	s.stake[addr] = current - amt
	s.total -= amt
	s.mu.Unlock()
}

// Balance returns the staked balance for the provided address.
func (s *StakingNode) Balance(addr string) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stake[addr]
}

// TotalStaked returns the total amount staked across all addresses.
func (s *StakingNode) TotalStaked() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.total
}
