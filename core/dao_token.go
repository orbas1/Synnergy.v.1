package core

import "errors"

// DAOTokenLedger tracks DAO membership token balances.
type DAOTokenLedger struct {
	balances map[string]uint64
}

// NewDAOTokenLedger returns an initialised ledger.
func NewDAOTokenLedger() *DAOTokenLedger {
	return &DAOTokenLedger{balances: make(map[string]uint64)}
}

// Mint creates tokens for an address.
func (l *DAOTokenLedger) Mint(addr string, amount uint64) {
	l.balances[addr] += amount
}

// Transfer moves tokens between addresses.
func (l *DAOTokenLedger) Transfer(from, to string, amount uint64) error {
	if l.balances[from] < amount {
		return errors.New("insufficient balance")
	}
	l.balances[from] -= amount
	l.balances[to] += amount
	return nil
}

// Balance returns the token balance for an address.
func (l *DAOTokenLedger) Balance(addr string) uint64 {
	return l.balances[addr]
}
