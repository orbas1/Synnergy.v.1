package core

import (
	"errors"
	"sync"
)

// errInsufficientStake is returned when an unstake exceeds the balance.
var errInsufficientStake = errors.New("insufficient stake")

// DAOStaking tracks staked tokens for DAO members per DAO instance.
// It references the DAOManager to enforce membership checks before any
// staking operation is performed.
type DAOStaking struct {
	mu     sync.RWMutex
	stakes map[string]map[string]uint64 // daoID -> addr -> amount
	mgr    *DAOManager
}

// NewDAOStaking creates a new DAOStaking instance bound to a DAOManager.
func NewDAOStaking(mgr *DAOManager) *DAOStaking {
	return &DAOStaking{stakes: make(map[string]map[string]uint64), mgr: mgr}
}

// Stake adds tokens to a member's stake. The caller must be a DAO member.
func (s *DAOStaking) Stake(daoID, addr string, amount uint64) error {
	dao, err := s.mgr.Info(daoID)
	if err != nil {
		return err
	}
	if !dao.IsMember(addr) {
		return errUnauthorized
	}
	s.mu.Lock()
	if s.stakes[daoID] == nil {
		s.stakes[daoID] = make(map[string]uint64)
	}
	s.stakes[daoID][addr] += amount
	s.mu.Unlock()
	return nil
}

// Unstake removes tokens from a member's stake. Only DAO members can unstake.
func (s *DAOStaking) Unstake(daoID, addr string, amount uint64) error {
	dao, err := s.mgr.Info(daoID)
	if err != nil {
		return err
	}
	if !dao.IsMember(addr) {
		return errUnauthorized
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	bal := s.stakes[daoID][addr]
	if bal < amount {
		return errInsufficientStake
	}
	s.stakes[daoID][addr] = bal - amount
	return nil
}

// Balance returns the staked balance for an address within a DAO.
func (s *DAOStaking) Balance(daoID, addr string) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stakes[daoID][addr]
}

// TotalStaked returns the sum of all staked tokens for a DAO.
func (s *DAOStaking) TotalStaked(daoID string) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var total uint64
	for _, amt := range s.stakes[daoID] {
		total += amt
	}
	return total
}
