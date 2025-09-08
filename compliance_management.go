package synnergy

import (
	"errors"
	"sync"
)

// Transaction is a minimal representation used for compliance checks.
type Transaction struct {
	From   string // Sender address
	To     string // Receiver address
	Amount uint64 // Amount to transfer
}

// ComplianceManager manages address suspensions and whitelists.
type ComplianceManager struct {
	mu        sync.RWMutex
	suspended map[string]bool
	whitelist map[string]bool
}

// NewComplianceManager creates a new ComplianceManager.
func NewComplianceManager() *ComplianceManager {
	return &ComplianceManager{
		suspended: make(map[string]bool),
		whitelist: make(map[string]bool),
	}
}

// Suspend marks an address as suspended from transfers.
// It returns an error if the address is empty or already suspended.
func (m *ComplianceManager) Suspend(addr string) error {
	if addr == "" {
		return errors.New("address required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.suspended[addr] {
		return errors.New("address already suspended")
	}
	m.suspended[addr] = true
	return nil
}

// Resume lifts a suspension for an address.
// It returns an error if the address is not currently suspended.
func (m *ComplianceManager) Resume(addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.suspended[addr] {
		return errors.New("address not suspended")
	}
	delete(m.suspended, addr)
	return nil
}

// Whitelist adds an address to the whitelist.
// It returns an error if the address is empty or already whitelisted.
func (m *ComplianceManager) Whitelist(addr string) error {
	if addr == "" {
		return errors.New("address required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.whitelist[addr] {
		return errors.New("address already whitelisted")
	}
	m.whitelist[addr] = true
	return nil
}

// Unwhitelist removes an address from the whitelist.
// It returns an error if the address is not whitelisted.
func (m *ComplianceManager) Unwhitelist(addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.whitelist[addr] {
		return errors.New("address not whitelisted")
	}
	delete(m.whitelist, addr)
	return nil
}

// Status returns suspension and whitelist status for an address.
func (m *ComplianceManager) Status(addr string) (suspended, whitelisted bool) {
	m.mu.RLock()
	suspended = m.suspended[addr]
	whitelisted = m.whitelist[addr]
	m.mu.RUnlock()
	return
}

// ReviewTransaction checks if a transaction involves suspended parties.
func (m *ComplianceManager) ReviewTransaction(tx Transaction) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	parties := []string{tx.From, tx.To}
	for _, p := range parties {
		if m.suspended[p] && !m.whitelist[p] {
			return errors.New("transaction involves suspended address")
		}
	}
	return nil
}
