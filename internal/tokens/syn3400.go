package tokens

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

var (
	ErrForexPairNotFound = errors.New("tokens: forex pair not found")
	ErrInvalidForexRate  = errors.New("tokens: forex rate must be positive")
)

// ForexPair defines pair specific information for SYN3400 tokens.
type ForexPair struct {
	PairID    string
	Base      string
	Quote     string
	Rate      float64
	UpdatedAt time.Time
}

// ForexRegistry tracks registered Forex pairs.
type ForexRegistry struct {
	mu       sync.RWMutex
	pairs    map[string]*ForexPair
	bySymbol map[string]string
	counter  uint64
}

// NewForexRegistry creates an empty registry.
func NewForexRegistry() *ForexRegistry {
	return &ForexRegistry{pairs: make(map[string]*ForexPair), bySymbol: make(map[string]string)}
}

func key(base, quote string) string {
	return fmt.Sprintf("%s-%s", base, quote)
}

// Register adds a new forex pair to the registry.
func (r *ForexRegistry) Register(base, quote string, rate float64) (*ForexPair, error) {
	if base == "" || quote == "" {
		return nil, errors.New("base and quote required")
	}
	if rate <= 0 {
		return nil, ErrInvalidForexRate
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if existingID, ok := r.bySymbol[key(base, quote)]; ok {
		pair := r.pairs[existingID]
		pair.Rate = rate
		pair.UpdatedAt = time.Now()
		return cloneForexPair(pair), nil
	}
	r.counter++
	id := fmt.Sprintf("FX-%d", r.counter)
	p := &ForexPair{PairID: id, Base: base, Quote: quote, Rate: rate, UpdatedAt: time.Now()}
	r.pairs[id] = p
	r.bySymbol[key(base, quote)] = id
	return cloneForexPair(p), nil
}

// UpdateRate updates the exchange rate for a pair.
func (r *ForexRegistry) UpdateRate(pairID string, rate float64) error {
	if rate <= 0 {
		return ErrInvalidForexRate
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.pairs[pairID]
	if !ok {
		return ErrForexPairNotFound
	}
	p.Rate = rate
	p.UpdatedAt = time.Now()
	return nil
}

// Get retrieves a forex pair by ID.
func (r *ForexRegistry) Get(pairID string) (*ForexPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pairs[pairID]
	if !ok {
		return nil, ErrForexPairNotFound
	}
	return cloneForexPair(p), nil
}

// GetBySymbol retrieves a pair by base and quote symbols.
func (r *ForexRegistry) GetBySymbol(base, quote string) (*ForexPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.bySymbol[key(base, quote)]
	if !ok {
		return nil, ErrForexPairNotFound
	}
	return cloneForexPair(r.pairs[id]), nil
}

// List returns all registered forex pairs sorted by identifier.
func (r *ForexRegistry) List() []*ForexPair {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*ForexPair, 0, len(r.pairs))
	for _, p := range r.pairs {
		res = append(res, cloneForexPair(p))
	}
	sort.Slice(res, func(i, j int) bool { return res[i].PairID < res[j].PairID })
	return res
}

// Remove deletes a pair from the registry.
func (r *ForexRegistry) Remove(pairID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.pairs[pairID]
	if !ok {
		return ErrForexPairNotFound
	}
	delete(r.pairs, pairID)
	delete(r.bySymbol, key(p.Base, p.Quote))
	return nil
}

func cloneForexPair(p *ForexPair) *ForexPair {
	cp := *p
	return &cp
}
