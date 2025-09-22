package tokens

import (
	"sort"
	"sync"
	"time"
)

// Registry maintains a lookup of token instances by their identifiers. It is
// safe for concurrent use by multiple goroutines and emits events when tokens
// are registered or removed.
type Registry struct {
	mu        sync.RWMutex
	tokens    map[TokenID]Token
	next      TokenID
	observers []RegistryObserver
	clock     func() time.Time
}

// RegistryOption customises registry behaviour.
type RegistryOption func(*Registry)

// WithRegistryClock overrides the registry clock primarily for deterministic testing.
func WithRegistryClock(clock func() time.Time) RegistryOption {
	return func(r *Registry) {
		if clock != nil {
			r.clock = clock
		}
	}
}

// NewRegistry creates an empty token registry.
func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{
		tokens: make(map[TokenID]Token),
		clock: func() time.Time {
			return time.Now().UTC()
		},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// RegistryEventType enumerates lifecycle events.
type RegistryEventType string

const (
	RegistryEventRegistered RegistryEventType = "registered"
	RegistryEventRemoved    RegistryEventType = "removed"
)

// RegistryEvent represents a registry change delivered to observers.
type RegistryEvent struct {
	Type      RegistryEventType
	Info      TokenInfo
	Timestamp time.Time
}

// RegistryObserver consumes registry events.
type RegistryObserver func(RegistryEvent)

// RegisterObserver attaches an observer. Observers are invoked synchronously
// after the registry lock is released to avoid blocking writers.
func (r *Registry) RegisterObserver(obs RegistryObserver) {
	if obs == nil {
		return
	}
	r.mu.Lock()
	r.observers = append(r.observers, obs)
	r.mu.Unlock()
}

// NextID returns a unique identifier for a new token.
func (r *Registry) NextID() TokenID {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	return r.next
}

// Register adds the token to the registry using its ID.
func (r *Registry) Register(t Token) {
	if t == nil {
		return
	}
	r.mu.Lock()
	r.tokens[t.ID()] = t
	observers := append([]RegistryObserver(nil), r.observers...)
	info := r.infoLocked(t)
	ts := r.clock()
	r.mu.Unlock()
	r.emit(RegistryEvent{Type: RegistryEventRegistered, Info: info, Timestamp: ts}, observers)
}

// Remove deletes a token from the registry returning true when successful.
func (r *Registry) Remove(id TokenID) bool {
	r.mu.Lock()
	t, ok := r.tokens[id]
	if !ok {
		r.mu.Unlock()
		return false
	}
	delete(r.tokens, id)
	observers := append([]RegistryObserver(nil), r.observers...)
	info := r.infoLocked(t)
	ts := r.clock()
	r.mu.Unlock()
	r.emit(RegistryEvent{Type: RegistryEventRemoved, Info: info, Timestamp: ts}, observers)
	return true
}

// Get retrieves a token by ID if present.
func (r *Registry) Get(id TokenID) (Token, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[id]
	return t, ok
}

// GetBySymbol retrieves a token by its symbol.
func (r *Registry) GetBySymbol(symbol string) (Token, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tokens {
		if t.Symbol() == symbol {
			return t, true
		}
	}
	return nil, false
}

// TokenInfo summarises metadata for a registered token.
type TokenInfo struct {
	ID          TokenID
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply uint64
}

// Info returns metadata for a token by ID.
func (r *Registry) Info(id TokenID) (TokenInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[id]
	if !ok {
		return TokenInfo{}, false
	}
	return r.infoLocked(t), true
}

func (r *Registry) infoLocked(t Token) TokenInfo {
	return TokenInfo{
		ID:          t.ID(),
		Name:        t.Name(),
		Symbol:      t.Symbol(),
		Decimals:    t.Decimals(),
		TotalSupply: t.TotalSupply(),
	}
}

// List returns metadata for all registered tokens sorted by identifier.
func (r *Registry) List() []TokenInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	infos := make([]TokenInfo, 0, len(r.tokens))
	for _, t := range r.tokens {
		infos = append(infos, r.infoLocked(t))
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ID < infos[j].ID
	})
	return infos
}

// Snapshot returns the same data as List and is provided for API symmetry.
func (r *Registry) Snapshot() []TokenInfo {
	return r.List()
}

func (r *Registry) emit(evt RegistryEvent, observers []RegistryObserver) {
	for _, obs := range observers {
		func(o RegistryObserver) {
			defer func() {
				_ = recover()
			}()
			o(evt)
		}(obs)
	}
}
