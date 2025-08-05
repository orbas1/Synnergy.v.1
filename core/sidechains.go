package core

import (
	"fmt"
	"sync"
)

// Sidechain represents a registered side-chain with metadata and state.
type Sidechain struct {
	ID         string
	Metadata   string
	Header     string
	Validators []string
	Paused     bool
	Deposits   map[string]uint64 // escrow balances per address
}

// SidechainRegistry tracks registered side-chains.
type SidechainRegistry struct {
	mu     sync.RWMutex
	chains map[string]*Sidechain
}

// NewSidechainRegistry creates an empty registry.
func NewSidechainRegistry() *SidechainRegistry {
	return &SidechainRegistry{chains: make(map[string]*Sidechain)}
}

// Register registers a new side-chain.
func (r *SidechainRegistry) Register(id, meta string, validators []string) (*Sidechain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.chains[id]; exists {
		return nil, fmt.Errorf("sidechain %s exists", id)
	}
	sc := &Sidechain{ID: id, Metadata: meta, Validators: validators, Deposits: make(map[string]uint64)}
	r.chains[id] = sc
	return sc, nil
}

// SubmitHeader records the latest header for a side-chain.
func (r *SidechainRegistry) SubmitHeader(id, header string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sc, ok := r.chains[id]
	if !ok {
		return fmt.Errorf("sidechain %s not found", id)
	}
	sc.Header = header
	return nil
}

// GetHeader retrieves the latest header for a side-chain.
func (r *SidechainRegistry) GetHeader(id string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sc, ok := r.chains[id]
	if !ok {
		return "", false
	}
	return sc.Header, true
}

// Meta returns metadata for a side-chain.
func (r *SidechainRegistry) Meta(id string) (Sidechain, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sc, ok := r.chains[id]
	if !ok {
		return Sidechain{}, false
	}
	return *sc, true
}

// List returns all registered side-chains.
func (r *SidechainRegistry) List() []Sidechain {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Sidechain, 0, len(r.chains))
	for _, sc := range r.chains {
		out = append(out, *sc)
	}
	return out
}

// Pause suspends a side-chain.
func (r *SidechainRegistry) Pause(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sc, ok := r.chains[id]
	if !ok {
		return fmt.Errorf("sidechain %s not found", id)
	}
	sc.Paused = true
	return nil
}

// Resume resumes a paused side-chain.
func (r *SidechainRegistry) Resume(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sc, ok := r.chains[id]
	if !ok {
		return fmt.Errorf("sidechain %s not found", id)
	}
	sc.Paused = false
	return nil
}

// UpdateValidators updates the validator set for a side-chain.
func (r *SidechainRegistry) UpdateValidators(id string, validators []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sc, ok := r.chains[id]
	if !ok {
		return fmt.Errorf("sidechain %s not found", id)
	}
	sc.Validators = validators
	return nil
}

// Remove deletes a side-chain and its data.
func (r *SidechainRegistry) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.chains[id]; !ok {
		return fmt.Errorf("sidechain %s not found", id)
	}
	delete(r.chains, id)
	return nil
}
