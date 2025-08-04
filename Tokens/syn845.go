package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// DebtMetadata stores comprehensive data about a debt instrument.
type DebtRecord struct {
	DebtID    string
	Borrower  string
	Principal uint64
	Rate      float64
	Penalty   float64
	Due       time.Time
	Paid      uint64
}

// DebtToken represents a SYN845 debt token.
type DebtToken struct {
	ID     string
	Name   string
	Symbol string
	Owner  string
	Supply uint64
	Debts  map[string]*DebtRecord
}

// DebtRegistry manages debt tokens.
type DebtRegistry struct {
	mu      sync.RWMutex
	tokens  map[string]*DebtToken
	counter uint64
}

// NewDebtRegistry creates an empty debt token registry.
func NewDebtRegistry() *DebtRegistry {
	return &DebtRegistry{tokens: make(map[string]*DebtToken)}
}

// CreateToken initializes a new debt token.
func (r *DebtRegistry) CreateToken(name, symbol, owner string, supply uint64) (string, *DebtToken) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("DEBT-%d", r.counter)
	t := &DebtToken{ID: id, Name: name, Symbol: symbol, Owner: owner, Supply: supply, Debts: make(map[string]*DebtRecord)}
	r.tokens[id] = t
	return id, t
}

// IssueDebt issues a debt instrument under a token.
func (r *DebtRegistry) IssueDebt(tokenID, debtID, borrower string, principal uint64, rate, penalty float64, due time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	if _, exists := t.Debts[debtID]; exists {
		return errors.New("debt already exists")
	}
	t.Debts[debtID] = &DebtRecord{DebtID: debtID, Borrower: borrower, Principal: principal, Rate: rate, Penalty: penalty, Due: due}
	return nil
}

// RecordPayment records a payment toward a debt.
func (r *DebtRegistry) RecordPayment(tokenID, debtID string, amount uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	d, ok := t.Debts[debtID]
	if !ok {
		return errors.New("debt not found")
	}
	d.Paid += amount
	return nil
}

// GetDebt returns information about a specific debt.
func (r *DebtRegistry) GetDebt(tokenID, debtID string) (*DebtRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[tokenID]
	if !ok {
		return nil, errors.New("token not found")
	}
	d, ok := t.Debts[debtID]
	if !ok {
		return nil, errors.New("debt not found")
	}
	cp := *d
	return &cp, nil
}
