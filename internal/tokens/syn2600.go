package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// InvestorTokenMeta defines metadata for SYN2600 investor tokens.
type InvestorTokenMeta struct {
	ID       string
	Asset    string
	Owner    string
	Shares   uint64
	IssuedAt time.Time
	Expiry   time.Time
	Active   bool
	Returns  []ReturnRecord
}

// ReturnRecord logs a distribution paid to the investor.
type ReturnRecord struct {
	Amount uint64
	Time   time.Time
}

// InvestorRegistry tracks issued investor tokens.
type InvestorRegistry struct {
	mu      sync.RWMutex
	tokens  map[string]*InvestorTokenMeta
	counter uint64
}

// NewInvestorRegistry creates an empty registry.
func NewInvestorRegistry() *InvestorRegistry {
	return &InvestorRegistry{tokens: make(map[string]*InvestorTokenMeta)}
}

// Issue creates a new investor token for a given asset.
func (r *InvestorRegistry) Issue(asset, owner string, shares uint64, expiry time.Time) *InvestorTokenMeta {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("INV-%d", r.counter)
	tok := &InvestorTokenMeta{
		ID:       id,
		Asset:    asset,
		Owner:    owner,
		Shares:   shares,
		IssuedAt: time.Now(),
		Expiry:   expiry,
		Active:   true,
	}
	r.tokens[id] = tok
	return tok
}

// Transfer moves ownership of a token to a new owner.
func (r *InvestorRegistry) Transfer(tokenID, newOwner string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	tok.Owner = newOwner
	return nil
}

// RecordReturn records a return payment to an investor token.
func (r *InvestorRegistry) RecordReturn(tokenID string, amount uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	tok.Returns = append(tok.Returns, ReturnRecord{Amount: amount, Time: time.Now()})
	return nil
}

// Get retrieves a token by identifier.
func (r *InvestorRegistry) Get(tokenID string) (*InvestorTokenMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return nil, false
	}
	cp := *tok
	cp.Returns = append([]ReturnRecord(nil), tok.Returns...)
	return &cp, true
}

// List returns all issued investor tokens.
func (r *InvestorRegistry) List() []*InvestorTokenMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*InvestorTokenMeta, 0, len(r.tokens))
	for _, tok := range r.tokens {
		cp := *tok
		cp.Returns = append([]ReturnRecord(nil), tok.Returns...)
		res = append(res, &cp)
	}
	return res
}
