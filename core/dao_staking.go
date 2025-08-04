package core

import "errors"

// DAOStaking tracks staked tokens for DAO members.
type DAOStaking struct {
	stakes map[string]uint64
}

// NewDAOStaking creates a new DAOStaking instance.
func NewDAOStaking() *DAOStaking {
	return &DAOStaking{stakes: make(map[string]uint64)}
}

// Stake adds tokens to a member's stake.
func (s *DAOStaking) Stake(addr string, amount uint64) {
	s.stakes[addr] += amount
}

// Unstake removes tokens from a member's stake.
func (s *DAOStaking) Unstake(addr string, amount uint64) error {
	bal, ok := s.stakes[addr]
	if !ok || bal < amount {
		return errors.New("insufficient stake")
	}
	s.stakes[addr] = bal - amount
	return nil
}

// Balance returns the staked balance for an address.
func (s *DAOStaking) Balance(addr string) uint64 {
	return s.stakes[addr]
}
