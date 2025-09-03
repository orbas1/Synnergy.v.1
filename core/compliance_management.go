package core

import (
	"errors"
	"sync"

	ilog "synnergy/internal/log"
)

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
func (m *ComplianceManager) Suspend(addr string) {
	m.mu.Lock()
	m.suspended[addr] = true
	m.mu.Unlock()
	ilog.Info("suspend", "address", addr)
}

// Resume lifts a suspension for an address.
func (m *ComplianceManager) Resume(addr string) {
	m.mu.Lock()
	delete(m.suspended, addr)
	m.mu.Unlock()
	ilog.Info("resume", "address", addr)
}

// Whitelist adds an address to the whitelist.
func (m *ComplianceManager) Whitelist(addr string) {
	m.mu.Lock()
	m.whitelist[addr] = true
	m.mu.Unlock()
	ilog.Info("whitelist", "address", addr)
}

// Unwhitelist removes an address from the whitelist.
func (m *ComplianceManager) Unwhitelist(addr string) {
	m.mu.Lock()
	delete(m.whitelist, addr)
	m.mu.Unlock()
	ilog.Info("unwhitelist", "address", addr)
}

// Status returns suspension and whitelist status for an address.
func (m *ComplianceManager) Status(addr string) (suspended, whitelisted bool) {
	m.mu.RLock()
	suspended = m.suspended[addr]
	whitelisted = m.whitelist[addr]
	m.mu.RUnlock()
	ilog.Info("status", "address", addr, "suspended", suspended, "whitelisted", whitelisted)
	return
}

// ReviewTransaction checks if a transaction involves suspended parties.
func (m *ComplianceManager) ReviewTransaction(tx Transaction) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	parties := []string{tx.From, tx.To}
	for _, p := range parties {
		if m.suspended[p] && !m.whitelist[p] {
			ilog.Error("review_tx", "address", p)
			return errors.New("transaction involves suspended address")
		}
	}
	return nil
}
