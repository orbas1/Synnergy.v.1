package synnergy

import (
	"errors"
	"sync"
)

// Regulation defines a basic rule for transaction oversight.
type Regulation struct {
	ID           string // Unique identifier of the regulation
	Jurisdiction string // Optional jurisdiction reference
	Description  string // Human readable description
	MaxAmount    uint64 // Maximum allowed amount for a transaction
}

// RegulatoryManager stores and evaluates regulations.
type RegulatoryManager struct {
	mu          sync.RWMutex
	regulations map[string]Regulation
}

// NewRegulatoryManager creates a new RegulatoryManager instance.
func NewRegulatoryManager() *RegulatoryManager {
	return &RegulatoryManager{regulations: make(map[string]Regulation)}
}

// AddRegulation registers a new regulation.
func (m *RegulatoryManager) AddRegulation(reg Regulation) error {
	if reg.ID == "" {
		return errors.New("id required")
	}
	m.mu.Lock()
	m.regulations[reg.ID] = reg
	m.mu.Unlock()
	return nil
}

// RemoveRegulation removes a regulation by ID.
func (m *RegulatoryManager) RemoveRegulation(id string) {
	m.mu.Lock()
	delete(m.regulations, id)
	m.mu.Unlock()
}

// GetRegulation retrieves a regulation by ID.
func (m *RegulatoryManager) GetRegulation(id string) (Regulation, bool) {
	m.mu.RLock()
	reg, ok := m.regulations[id]
	m.mu.RUnlock()
	return reg, ok
}

// ListRegulations returns all registered regulations.
func (m *RegulatoryManager) ListRegulations() []Regulation {
	m.mu.RLock()
	regs := make([]Regulation, 0, len(m.regulations))
	for _, r := range m.regulations {
		regs = append(regs, r)
	}
	m.mu.RUnlock()
	return regs
}

// EvaluateTransaction returns IDs of regulations violated by the transaction.
func (m *RegulatoryManager) EvaluateTransaction(tx Transaction) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var violations []string
	for id, reg := range m.regulations {
		if reg.MaxAmount > 0 && tx.Amount > reg.MaxAmount {
			violations = append(violations, id)
		}
	}
	return violations
}
