package core

import (
	"errors"
	"sync"
)

// DAOTokenLedger tracks DAO membership token balances per DAO and
// enforces membership and admin permissions for minting, transferring
// and burning tokens.
type DAOTokenLedger struct {
	mu       sync.RWMutex
	balances map[string]map[string]uint64 // daoID -> addr -> balance
	daoMgr   *DAOManager
}

// NewDAOTokenLedger returns an initialised ledger bound to a DAO manager.
func NewDAOTokenLedger(mgr *DAOManager) *DAOTokenLedger {
	return &DAOTokenLedger{balances: make(map[string]map[string]uint64), daoMgr: mgr}
}

// Mint creates tokens for a DAO member. Only a DAO admin can mint tokens.
func (l *DAOTokenLedger) Mint(daoID, admin, addr string, amount uint64) error {
	dao, err := l.daoMgr.Info(daoID)
	if err != nil {
		return err
	}
	if !dao.IsAdmin(admin) {
		return errUnauthorized
	}
	if !dao.IsMember(addr) {
		return errMemberMissing
	}
	l.mu.Lock()
	if l.balances[daoID] == nil {
		l.balances[daoID] = make(map[string]uint64)
	}
	l.balances[daoID][addr] += amount
	l.mu.Unlock()
	return nil
}

// Transfer moves tokens between DAO members.
func (l *DAOTokenLedger) Transfer(daoID, from, to string, amount uint64) error {
	dao, err := l.daoMgr.Info(daoID)
	if err != nil {
		return err
	}
	if !dao.IsMember(from) || !dao.IsMember(to) {
		return errMemberMissing
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.balances[daoID] == nil || l.balances[daoID][from] < amount {
		return errors.New("insufficient balance")
	}
	l.balances[daoID][from] -= amount
	l.balances[daoID][to] += amount
	return nil
}

// Balance returns the token balance for a DAO member. Non-members always return 0.
func (l *DAOTokenLedger) Balance(daoID, addr string) uint64 {
	dao, err := l.daoMgr.Info(daoID)
	if err != nil || !dao.IsMember(addr) {
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.balances[daoID][addr]
}

// Burn removes tokens from a DAO member. Only a DAO admin can burn tokens.
func (l *DAOTokenLedger) Burn(daoID, admin, addr string, amount uint64) error {
	dao, err := l.daoMgr.Info(daoID)
	if err != nil {
		return err
	}
	if !dao.IsAdmin(admin) {
		return errUnauthorized
	}
	if !dao.IsMember(addr) {
		return errMemberMissing
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	bal := l.balances[daoID][addr]
	if bal < amount {
		return errors.New("insufficient balance")
	}
	l.balances[daoID][addr] = bal - amount
	return nil
}
