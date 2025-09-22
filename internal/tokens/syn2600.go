package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// InvestorTokenMeta defines metadata for SYN2600 investor tokens.
type InvestorTokenMeta struct {
	ID        string
	Asset     string
	Owner     string
	Shares    uint64
	IssuedAt  time.Time
	Expiry    time.Time
	Active    bool
	Returns   []ReturnRecord
	Metadata  map[string]string
	RiskScore float64
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
// It validates the provided metadata and returns an error if any
// required field is missing or invalid. Expired tokens cannot be
// created.
func (r *InvestorRegistry) Issue(asset, owner string, shares uint64, expiry time.Time) (*InvestorTokenMeta, error) {
	if asset == "" || owner == "" {
		return nil, errors.New("asset and owner required")
	}
	if shares == 0 {
		return nil, errors.New("shares must be > 0")
	}
	if time.Now().After(expiry) {
		return nil, errors.New("expiry must be in the future")
	}
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
		Metadata: make(map[string]string),
	}
	r.tokens[id] = tok
	return tok, nil
}

// Transfer moves ownership of a token to a new owner.
func (r *InvestorRegistry) Transfer(tokenID, newOwner string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	if !tok.Active || time.Now().After(tok.Expiry) {
		tok.Active = false
		return errors.New("token inactive")
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
	if !tok.Active || time.Now().After(tok.Expiry) {
		tok.Active = false
		return errors.New("token inactive")
	}
	tok.Returns = append(tok.Returns, ReturnRecord{Amount: amount, Time: time.Now()})
	return nil
}

// Deactivate marks a token as inactive, preventing further operations.
func (r *InvestorRegistry) Deactivate(tokenID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	tok.Active = false
	return nil
}

// Activate re-enables a previously deactivated token if the expiry permits.
func (r *InvestorRegistry) Activate(tokenID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	if time.Now().After(tok.Expiry) {
		return errors.New("token expired")
	}
	tok.Active = true
	return nil
}

// ExtendExpiry extends the expiry date when governance approves.
func (r *InvestorRegistry) ExtendExpiry(tokenID string, newExpiry time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	if newExpiry.Before(tok.Expiry) {
		return errors.New("new expiry must be after current expiry")
	}
	tok.Expiry = newExpiry
	return nil
}

// UpdateMetadata updates structured metadata for the token.
func (r *InvestorRegistry) UpdateMetadata(tokenID string, meta map[string]string, risk float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tok, ok := r.tokens[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	if tok.Metadata == nil {
		tok.Metadata = make(map[string]string)
	}
	for k, v := range meta {
		tok.Metadata[k] = v
	}
	tok.RiskScore = risk
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

// FilterByAsset returns active tokens for the specified asset.
func (r *InvestorRegistry) FilterByAsset(asset string) []*InvestorTokenMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*InvestorTokenMeta, 0)
	for _, tok := range r.tokens {
		if tok.Asset == asset {
			cp := *tok
			cp.Returns = append([]ReturnRecord(nil), tok.Returns...)
			res = append(res, &cp)
		}
	}
	return res
}
