package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
)

// errInsufficientStake is returned when an unstake exceeds the balance.
var errInsufficientStake = errors.New("insufficient stake")

// DAOStaking tracks staked tokens for DAO members per DAO instance.  It wires
// ledger transfers into the staking flow so that DAO treasuries reflect the
// locked capital backing governance decisions.
type DAOStaking struct {
	mu     sync.RWMutex
	stakes map[string]map[string]uint64 // daoID -> addr -> amount
	mgr    *DAOManager
	ledger *Ledger
}

// NewDAOStaking creates a new DAOStaking instance bound to a DAOManager and
// optional ledger. When ledger is provided stake and unstake operations will
// transfer the underlying funds to and from a deterministic DAO treasury
// address.
func NewDAOStaking(mgr *DAOManager, ledger *Ledger) *DAOStaking {
	return &DAOStaking{stakes: make(map[string]map[string]uint64), mgr: mgr, ledger: ledger}
}

// Stake adds tokens to a member's stake. The caller must be a DAO member. When
// backed by a ledger the staked amount is moved into the DAO treasury.
func (s *DAOStaking) Stake(daoID, addr string, amount uint64) error {
	dao, err := s.mgr.Info(daoID)
	if err != nil {
		return err
	}
	if !dao.IsMember(addr) {
		return errUnauthorized
	}
	if s.ledger != nil && amount > 0 {
		if err := s.ledger.Transfer(addr, s.treasuryAddress(daoID), amount, 0); err != nil {
			return err
		}
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
// Funds are returned from the DAO treasury when a ledger is configured.
func (s *DAOStaking) Unstake(daoID, addr string, amount uint64) error {
	dao, err := s.mgr.Info(daoID)
	if err != nil {
		return err
	}
	if !dao.IsMember(addr) {
		return errUnauthorized
	}
	s.mu.Lock()
	bal := s.stakes[daoID][addr]
	if bal < amount {
		s.mu.Unlock()
		return errInsufficientStake
	}
	s.stakes[daoID][addr] = bal - amount
	remaining := s.stakes[daoID][addr]
	if remaining == 0 {
		delete(s.stakes[daoID], addr)
	}
	s.mu.Unlock()
	if s.ledger != nil && amount > 0 {
		if err := s.ledger.Transfer(s.treasuryAddress(daoID), addr, amount, 0); err != nil {
			s.mu.Lock()
			s.stakes[daoID][addr] += amount
			s.mu.Unlock()
			return err
		}
	}
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

func (s *DAOStaking) treasuryAddress(daoID string) string {
	sum := sha256.Sum256([]byte("dao_treasury:" + daoID))
	return hex.EncodeToString(sum[:])
}
