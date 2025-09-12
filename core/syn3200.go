package core

import (
	"errors"
	"sync"
	"time"
)

// Bill represents a bill record managed by the SYN3200 standard.
type Bill struct {
	ID       string
	Issuer   string
	Payer    string
	Amount   uint64
	DueDate  time.Time
	Metadata string
	Payments []BillPayment
}

// BillPayment captures a single payment toward a bill.
type BillPayment struct {
	Payer     string
	Amount    uint64
	Timestamp time.Time
}

// BillRegistry manages bills and payments.
type BillRegistry struct {
	mu    sync.RWMutex
	bills map[string]*Bill
}

// NewBillRegistry creates a new registry.
func NewBillRegistry() *BillRegistry {
	return &BillRegistry{bills: make(map[string]*Bill)}
}

// Create adds a new bill to the registry.
func (r *BillRegistry) Create(id, issuer, payer string, amt uint64, due time.Time, meta string) (*Bill, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.bills[id]; exists {
		return nil, errors.New("bill already exists")
	}
	b := &Bill{ID: id, Issuer: issuer, Payer: payer, Amount: amt, DueDate: due, Metadata: meta}
	r.bills[id] = b
	return b, nil
}

// Pay records a payment toward the bill.
func (r *BillRegistry) Pay(id, payer string, amt uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bills[id]
	if !ok {
		return errors.New("bill not found")
	}
	if amt == 0 {
		return errors.New("amount must be positive")
	}
	b.Payments = append(b.Payments, BillPayment{Payer: payer, Amount: amt, Timestamp: time.Now()})
	if amt >= b.Amount {
		b.Amount = 0
	} else {
		b.Amount -= amt
	}
	return nil
}

// Adjust changes the amount due on a bill.
func (r *BillRegistry) Adjust(id string, amt uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bills[id]
	if !ok {
		return errors.New("bill not found")
	}
	b.Amount = amt
	return nil
}

// Get returns bill information.
func (r *BillRegistry) Get(id string) (*Bill, bool) {
	r.mu.RLock()
	b, ok := r.bills[id]
	r.mu.RUnlock()
	return b, ok
}
