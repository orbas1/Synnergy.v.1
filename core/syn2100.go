package core

import (
	"errors"
	"time"
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
func (t *TradeFinanceToken) RegisterDocument(docID, issuer, recipient string, amount uint64, issue, due time.Time, desc string) {
	t.Documents[docID] = &FinancialDocument{
		DocID:       docID,
		Issuer:      issuer,
		Recipient:   recipient,
		Amount:      amount,
		IssueDate:   issue,
		DueDate:     due,
		Description: desc,
	}
}

// FinanceDocument marks a document as financed by a financier.
func (t *TradeFinanceToken) FinanceDocument(docID, financier string) error {
	d, ok := t.Documents[docID]
	if !ok {
		return errors.New("document not found")
	}
	if d.Financed {
		return errors.New("document already financed")
	}
	d.Financed = true
	d.Financier = financier
	return nil
}

// GetDocument fetches a document by ID.
func (t *TradeFinanceToken) GetDocument(docID string) (*FinancialDocument, bool) {
	d, ok := t.Documents[docID]
	return d, ok
}

// ListDocuments returns all registered documents.
func (t *TradeFinanceToken) ListDocuments() []*FinancialDocument {
	res := make([]*FinancialDocument, 0, len(t.Documents))
	for _, d := range t.Documents {
		res = append(res, d)
	}
	return res
}

// AddLiquidity adds funds to the liquidity pool from an address.
func (t *TradeFinanceToken) AddLiquidity(addr string, amt uint64) {
	t.Liquidity[addr] += amt
}

// RemoveLiquidity removes funds from the liquidity pool.
func (t *TradeFinanceToken) RemoveLiquidity(addr string, amt uint64) error {
	if t.Liquidity[addr] < amt {
		return errors.New("insufficient liquidity")
	}
	t.Liquidity[addr] -= amt
	return nil
}
