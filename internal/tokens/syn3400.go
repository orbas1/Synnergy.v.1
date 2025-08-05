package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ForexMetadata defines pair specific information for SYN3400 tokens.
type ForexPair struct {
	PairID    string
	Base      string
	Quote     string
	Rate      float64
	UpdatedAt time.Time
}

// ForexRegistry tracks registered Forex pairs.
type ForexRegistry struct {
	mu      sync.RWMutex
	pairs   map[string]*ForexPair
	counter uint64
}

// NewForexRegistry creates an empty registry.
func NewForexRegistry() *ForexRegistry {
	return &ForexRegistry{pairs: make(map[string]*ForexPair)}
}

// Register adds a new forex pair to the registry.
func (r *ForexRegistry) Register(base, quote string, rate float64) *ForexPair {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("FX-%d", r.counter)
	p := &ForexPair{PairID: id, Base: base, Quote: quote, Rate: rate, UpdatedAt: time.Now()}
	r.pairs[id] = p
	return p
}

// UpdateRate updates the exchange rate for a pair.
func (r *ForexRegistry) UpdateRate(pairID string, rate float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.pairs[pairID]
	if !ok {
		return errors.New("pair not found")
	}
	p.Rate = rate
	p.UpdatedAt = time.Now()
	return nil
}

// Get retrieves a forex pair by ID.
func (r *ForexRegistry) Get(pairID string) (*ForexPair, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pairs[pairID]
	if !ok {
		return nil, false
	}
	cp := *p
	return &cp, true
}

// List returns all registered forex pairs.
func (r *ForexRegistry) List() []*ForexPair {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*ForexPair, 0, len(r.pairs))
	for _, p := range r.pairs {
		cp := *p
		res = append(res, &cp)
	}
	return res
}
