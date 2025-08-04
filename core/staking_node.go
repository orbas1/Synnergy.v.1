package core

import "sync"

// StakingNode manages token staking for governance or validation purposes.
type StakingNode struct {
	mu     sync.Mutex
	stakes map[string]uint64
}

// NewStakingNode initializes and returns a StakingNode instance.
func NewStakingNode() *StakingNode {
	return &StakingNode{stakes: make(map[string]uint64)}
}

// Stake locks tokens from addr, increasing their staked balance.
func (s *StakingNode) Stake(addr string, amt uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stakes[addr] += amt
}

// Unstake releases tokens for addr, reducing their staked balance. If the
// amount exceeds the current stake it removes the stake entirely.
func (s *StakingNode) Unstake(addr string, amt uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bal := s.stakes[addr]
	if amt >= bal {
		delete(s.stakes, addr)
		return
	}
	s.stakes[addr] = bal - amt
}

// Balance returns the current staked balance for addr.
func (s *StakingNode) Balance(addr string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stakes[addr]
}

// TotalStaked returns the total tokens staked across all addresses.
func (s *StakingNode) TotalStaked() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	var total uint64
	for _, amt := range s.stakes {
		total += amt
	}
	return total
}
