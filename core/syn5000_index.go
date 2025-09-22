package core

import (
	"errors"
	"sort"
	"sync"
)

// GamblingToken exposes methods of the SYN5000 token used by the CLI and web
// dashboards.
type GamblingToken interface {
	PlaceBet(bettor string, amount uint64, odds float64, game string) (uint64, error)
	ResolveBet(betID uint64, win bool) (uint64, error)
	CancelBet(betID uint64, note string) error
	GetBet(betID uint64) (*BetRecord, bool)
	ListBets(filter BetFilter) []BetRecord
	Snapshot() SYN5000Snapshot
}

// ErrIndexSymbolRequired is returned when a caller attempts to interact with
// the index without providing a symbol.
var ErrIndexSymbolRequired = errors.New("syn5000 index: symbol required")

// ErrIndexTokenNil is returned when attempting to register a nil token.
var ErrIndexTokenNil = errors.New("syn5000 index: token required")

// SYN5000Index maintains a thread-safe registry of gambling tokens, exposing
// helper functions for CLI rendering and API integrations.
type SYN5000Index struct {
	mu     sync.RWMutex
	tokens map[string]GamblingToken
}

// NewSYN5000Index constructs an empty index.
func NewSYN5000Index() *SYN5000Index {
	return &SYN5000Index{tokens: make(map[string]GamblingToken)}
}

// Register adds or refreshes a token reference under the provided symbol.
func (i *SYN5000Index) Register(symbol string, token GamblingToken) error {
	if symbol == "" {
		return ErrIndexSymbolRequired
	}
	if token == nil {
		return ErrIndexTokenNil
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.tokens[symbol] = token
	return nil
}

// Snapshot returns the current snapshot for the registered symbol.
func (i *SYN5000Index) Snapshot(symbol string) (SYN5000Snapshot, bool) {
	if symbol == "" {
		return SYN5000Snapshot{}, false
	}
	i.mu.RLock()
	token, ok := i.tokens[symbol]
	i.mu.RUnlock()
	if !ok {
		return SYN5000Snapshot{}, false
	}
	return token.Snapshot(), true
}

// Bets returns a filtered list of bets for the requested symbol.
func (i *SYN5000Index) Bets(symbol string, filter BetFilter) ([]BetRecord, bool) {
	if symbol == "" {
		return nil, false
	}
	i.mu.RLock()
	token, ok := i.tokens[symbol]
	i.mu.RUnlock()
	if !ok {
		return nil, false
	}
	bets := token.ListBets(filter)
	return bets, true
}

// Symbols returns all registered token symbols in deterministic order.
func (i *SYN5000Index) Symbols() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	symbols := make([]string, 0, len(i.tokens))
	for symbol := range i.tokens {
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)
	return symbols
}
