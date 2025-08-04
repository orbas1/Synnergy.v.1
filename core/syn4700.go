package core

import (
	"errors"
	"sync"
	"time"
)

// LegalTokenStatus enumerates the status values for a legal token.
type LegalTokenStatus string

const (
	LegalTokenStatusPending   LegalTokenStatus = "pending"
	LegalTokenStatusActive    LegalTokenStatus = "active"
	LegalTokenStatusCompleted LegalTokenStatus = "completed"
	LegalTokenStatusDisputed  LegalTokenStatus = "disputed"
)

// Dispute records dispute actions and optional results.
type Dispute struct {
	Action    string
	Result    string
	Timestamp time.Time
}

// LegalToken represents a SYN4700 legal token tied to a legal document.
type LegalToken struct {
	mu           sync.RWMutex
	ID           string
	Name         string
	Symbol       string
	DocumentType string
	DocumentHash string
	Expiry       time.Time
	Owner        string
	Supply       uint64
	Parties      []string
	Signatures   map[string]string
	Status       LegalTokenStatus
	Disputes     []Dispute
}

// NewLegalToken creates a new legal token instance.
func NewLegalToken(id, name, symbol, docType, hash, owner string, expiry time.Time, supply uint64, parties []string) *LegalToken {
	cp := make([]string, len(parties))
	copy(cp, parties)
	return &LegalToken{
		ID:           id,
		Name:         name,
		Symbol:       symbol,
		DocumentType: docType,
		DocumentHash: hash,
		Expiry:       expiry,
		Owner:        owner,
		Supply:       supply,
		Parties:      cp,
		Signatures:   make(map[string]string),
		Status:       LegalTokenStatusPending,
	}
}

// Sign records a party signature on the legal token.
func (t *LegalToken) Sign(party, sig string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.partyExists(party) {
		return errors.New("unknown party")
	}
	t.Signatures[party] = sig
	return nil
}

// RevokeSignature removes a party's signature.
func (t *LegalToken) RevokeSignature(party string) {
	t.mu.Lock()
	delete(t.Signatures, party)
	t.mu.Unlock()
}

// UpdateStatus sets the current status of the legal token.
func (t *LegalToken) UpdateStatus(status LegalTokenStatus) {
	t.mu.Lock()
	t.Status = status
	t.mu.Unlock()
}

// Dispute records a dispute action and optional result. Status becomes disputed.
func (t *LegalToken) Dispute(action, result string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Status = LegalTokenStatusDisputed
	t.Disputes = append(t.Disputes, Dispute{Action: action, Result: result, Timestamp: time.Now()})
}

func (t *LegalToken) partyExists(party string) bool {
	for _, p := range t.Parties {
		if p == party {
			return true
		}
	}
	return false
}

// LegalTokenRegistry manages legal tokens by ID.
type LegalTokenRegistry struct {
	mu     sync.RWMutex
	tokens map[string]*LegalToken
}

// NewLegalTokenRegistry creates an empty registry.
func NewLegalTokenRegistry() *LegalTokenRegistry {
	return &LegalTokenRegistry{tokens: make(map[string]*LegalToken)}
}

// Add inserts or replaces a legal token in the registry.
func (r *LegalTokenRegistry) Add(t *LegalToken) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[t.ID] = t
}

// Get retrieves a legal token by ID.
func (r *LegalTokenRegistry) Get(id string) (*LegalToken, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[id]
	return t, ok
}

// Remove deletes a legal token from the registry.
func (r *LegalTokenRegistry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, id)
}

// List returns all legal tokens in the registry.
func (r *LegalTokenRegistry) List() []*LegalToken {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*LegalToken, 0, len(r.tokens))
	for _, t := range r.tokens {
		list = append(list, t)
	}
	return list
}
