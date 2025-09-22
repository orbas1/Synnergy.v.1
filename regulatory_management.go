package synnergy

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// Regulation defines rules for transaction oversight.  Additional optional
// fields provide granular control for enterprise deployments while remaining
// backward compatible with earlier modules that only populate MaxAmount.
type Regulation struct {
	ID               string
	Jurisdiction     string
	Description      string
	MaxAmount        uint64
	MinAmount        uint64
	AllowedEntities  []string
	DeniedEntities   []string
	RequireWhitelist bool
	Expiry           time.Time
}

// Violation captures a specific regulation breach alongside a human readable
// description that can be surfaced in CLI and UI audit logs.
type Violation struct {
	RegulationID string
	Reason       string
}

// EvaluationResult provides structured feedback from the compliance check.
type EvaluationResult struct {
	Violations []Violation
	CheckedAt  time.Time
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
		return errors.New("regulation id required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.regulations[reg.ID]; exists {
		return fmt.Errorf("regulation %s already exists", reg.ID)
	}
	m.regulations[reg.ID] = normaliseRegulation(reg)
	return nil
}

// UpdateRegulation replaces an existing regulation. Returns error if not found.
func (m *RegulatoryManager) UpdateRegulation(reg Regulation) error {
	if reg.ID == "" {
		return errors.New("regulation id required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.regulations[reg.ID]; !ok {
		return errors.New("regulation not found")
	}
	m.regulations[reg.ID] = normaliseRegulation(reg)
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

// ListRegulations returns all registered regulations ordered by identifier for
// deterministic CLI output.
func (m *RegulatoryManager) ListRegulations() []Regulation {
	m.mu.RLock()
	regs := make([]Regulation, 0, len(m.regulations))
	for _, r := range m.regulations {
		regs = append(regs, r)
	}
	m.mu.RUnlock()
	sort.Slice(regs, func(i, j int) bool { return regs[i].ID < regs[j].ID })
	return regs
}

// EvaluateTransaction returns IDs of regulations violated by the transaction.
func (m *RegulatoryManager) EvaluateTransaction(tx Transaction) []string {
	res := m.EvaluateTransactionDetailed(tx)
	ids := make([]string, 0, len(res.Violations))
	for _, v := range res.Violations {
		ids = append(ids, v.RegulationID)
	}
	return ids
}

// EvaluateTransactionDetailed returns detailed violation data for auditing.
func (m *RegulatoryManager) EvaluateTransactionDetailed(tx Transaction) EvaluationResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := time.Now().UTC()
	result := EvaluationResult{CheckedAt: now}

	for id, reg := range m.regulations {
		if reg.Expiry.IsZero() == false && now.After(reg.Expiry) {
			continue
		}
		if reg.MaxAmount > 0 && tx.Amount > reg.MaxAmount {
			result.Violations = append(result.Violations, Violation{RegulationID: id, Reason: fmt.Sprintf("amount %d exceeds max %d", tx.Amount, reg.MaxAmount)})
			continue
		}
		if reg.MinAmount > 0 && tx.Amount < reg.MinAmount {
			result.Violations = append(result.Violations, Violation{RegulationID: id, Reason: fmt.Sprintf("amount %d below minimum %d", tx.Amount, reg.MinAmount)})
			continue
		}
		if entityDenied(reg.DeniedEntities, tx.From) || entityDenied(reg.DeniedEntities, tx.To) {
			result.Violations = append(result.Violations, Violation{RegulationID: id, Reason: "counterparty denied"})
			continue
		}
		if reg.RequireWhitelist && !entityAllowed(reg.AllowedEntities, tx.From, tx.To) {
			result.Violations = append(result.Violations, Violation{RegulationID: id, Reason: "parties not whitelisted"})
			continue
		}
	}

	sort.Slice(result.Violations, func(i, j int) bool { return result.Violations[i].RegulationID < result.Violations[j].RegulationID })
	return result
}

// ValidateTransaction returns an error describing violated regulations, or nil
// if the transaction complies with all registered rules.
func (m *RegulatoryManager) ValidateTransaction(tx Transaction) error {
	res := m.EvaluateTransactionDetailed(tx)
	if len(res.Violations) == 0 {
		return nil
	}
	var parts []string
	for _, v := range res.Violations {
		parts = append(parts, fmt.Sprintf("%s: %s", v.RegulationID, v.Reason))
	}
	return errors.New("violations: " + strings.Join(parts, "; "))
}

func normaliseRegulation(reg Regulation) Regulation {
	reg.AllowedEntities = normaliseEntities(reg.AllowedEntities)
	reg.DeniedEntities = normaliseEntities(reg.DeniedEntities)
	return reg
}

func normaliseEntities(entities []string) []string {
	out := make([]string, 0, len(entities))
	seen := make(map[string]struct{}, len(entities))
	for _, e := range entities {
		e = strings.TrimSpace(strings.ToLower(e))
		if e == "" {
			continue
		}
		if _, ok := seen[e]; ok {
			continue
		}
		seen[e] = struct{}{}
		out = append(out, e)
	}
	sort.Strings(out)
	return out
}

func entityDenied(denied []string, entities ...string) bool {
	for _, entity := range entities {
		norm := strings.ToLower(entity)
		for _, d := range denied {
			if norm == d {
				return true
			}
		}
	}
	return false
}

func entityAllowed(allowed []string, entities ...string) bool {
	if len(allowed) == 0 {
		return false
	}
	for _, entity := range entities {
		norm := strings.ToLower(entity)
		for _, a := range allowed {
			if norm == a {
				return true
			}
		}
	}
	return false
}
