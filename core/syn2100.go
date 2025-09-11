package core

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrDocumentExists signals a duplicate document registration.
	ErrDocumentExists = errors.New("document exists")
	// ErrDocumentNotFound indicates a missing document lookup.
	ErrDocumentNotFound = errors.New("document not found")
	// ErrDocumentAlreadyFinanced is returned when financing an already financed document.
	ErrDocumentAlreadyFinanced = errors.New("document already financed")
	// ErrInsufficientLiquidity is returned when removing more liquidity than available.
	ErrInsufficientLiquidity = errors.New("insufficient liquidity")
)

// FinancialDocument captures metadata for a trade finance instrument.
type FinancialDocument struct {
	DocID       string
	Issuer      string
	Recipient   string
	Amount      uint64
	IssueDate   time.Time
	DueDate     time.Time
	Description string
	Financed    bool
	Financier   string
}

// TradeFinanceToken manages financial documents and liquidity pools for SYN2100.
type TradeFinanceToken struct {
	mu        sync.RWMutex
	Documents map[string]*FinancialDocument
	Liquidity map[string]uint64
}

// NewTradeFinanceToken creates a new registry instance.
func NewTradeFinanceToken() *TradeFinanceToken {
	return &TradeFinanceToken{
		Documents: make(map[string]*FinancialDocument),
		Liquidity: make(map[string]uint64),
	}
}

// RegisterDocument registers a financing document.
func (t *TradeFinanceToken) RegisterDocument(docID, issuer, recipient string, amount uint64, issue, due time.Time, desc string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.Documents[docID]; exists {
		return ErrDocumentExists
	}
	t.Documents[docID] = &FinancialDocument{
		DocID:       docID,
		Issuer:      issuer,
		Recipient:   recipient,
		Amount:      amount,
		IssueDate:   issue,
		DueDate:     due,
		Description: desc,
	}
	return nil
}

// FinanceDocument marks a document as financed by a financier.
func (t *TradeFinanceToken) FinanceDocument(docID, financier string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	d, ok := t.Documents[docID]
	if !ok {
		return ErrDocumentNotFound
	}
	if d.Financed {
		return ErrDocumentAlreadyFinanced
	}
	d.Financed = true
	d.Financier = financier
	return nil
}

// GetDocument fetches a document by ID.
func (t *TradeFinanceToken) GetDocument(docID string) (*FinancialDocument, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	d, ok := t.Documents[docID]
	if !ok {
		return nil, false
	}
	cp := *d
	return &cp, true
}

// ListDocuments returns all registered documents.
func (t *TradeFinanceToken) ListDocuments() []*FinancialDocument {
	t.mu.RLock()
	defer t.mu.RUnlock()
	res := make([]*FinancialDocument, 0, len(t.Documents))
	for _, d := range t.Documents {
		cp := *d
		res = append(res, &cp)
	}
	return res
}

// AddLiquidity adds funds to the liquidity pool from an address.
func (t *TradeFinanceToken) AddLiquidity(addr string, amt uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Liquidity[addr] += amt
}

// RemoveLiquidity removes funds from the liquidity pool.
func (t *TradeFinanceToken) RemoveLiquidity(addr string, amt uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Liquidity[addr] < amt {
		return ErrInsufficientLiquidity
	}
	t.Liquidity[addr] -= amt
	return nil
}
