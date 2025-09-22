package tokens

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

// SYN1000Index manages multiple SYN1000 stablecoin instances. It protects the
// underlying map with a mutex so concurrent CLI calls remain safe.
type SYN1000Index struct {
	mu       sync.RWMutex
	tokens   map[TokenID]*SYN1000Token
	next     TokenID
	clock    func() time.Time
	watchers []IndexWatcher
}

// IndexWatcher consumes lifecycle events from the index.
type IndexWatcher func(event IndexEvent)

// IndexEventType represents actions emitted by the index.
type IndexEventType string

const (
	IndexEventCreated IndexEventType = "created"
	IndexEventUpdated IndexEventType = "updated"
)

// IndexEvent captures token changes for observability.
type IndexEvent struct {
	Type      IndexEventType
	TokenID   TokenID
	Timestamp time.Time
	Metadata  map[string]string
}

// NewSYN1000Index creates an empty index.
func NewSYN1000Index() *SYN1000Index {
	return &SYN1000Index{
		tokens: make(map[TokenID]*SYN1000Token),
		clock:  func() time.Time { return time.Now().UTC() },
	}
}

// RegisterWatcher attaches an observer to lifecycle events.
func (i *SYN1000Index) RegisterWatcher(w IndexWatcher) {
	if w == nil {
		return
	}
	i.mu.Lock()
	i.watchers = append(i.watchers, w)
	i.mu.Unlock()
}

func (i *SYN1000Index) emit(evt IndexEvent) {
	i.mu.RLock()
	watchers := append([]IndexWatcher(nil), i.watchers...)
	i.mu.RUnlock()
	for _, w := range watchers {
		func(w IndexWatcher) {
			defer func() { _ = recover() }()
			w(evt)
		}(w)
	}
}

// Create instantiates a new SYN1000 token and returns its identifier.
func (i *SYN1000Index) Create(name, symbol string, decimals uint8) TokenID {
	i.mu.Lock()
	i.next++
	id := i.next
	token := NewSYN1000Token(id, name, symbol, decimals)
	i.tokens[id] = token
	ts := i.clock()
	i.mu.Unlock()
	i.emit(IndexEvent{Type: IndexEventCreated, TokenID: id, Timestamp: ts, Metadata: map[string]string{"symbol": symbol}})
	return id
}

// Token retrieves a SYN1000 token by ID.
func (i *SYN1000Index) Token(id TokenID) (*SYN1000Token, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	t, ok := i.tokens[id]
	if !ok {
		return nil, fmt.Errorf("token not found")
	}
	return t, nil
}

// AddReserve adds reserve assets to a specific token.
func (i *SYN1000Index) AddReserve(id TokenID, asset string, amount *big.Rat) error {
	t, err := i.Token(id)
	if err != nil {
		return err
	}
	if err := t.AddReserve(asset, amount); err != nil {
		return err
	}
	i.emit(IndexEvent{Type: IndexEventUpdated, TokenID: id, Timestamp: i.clock(), Metadata: map[string]string{"action": "add_reserve", "asset": asset}})
	return nil
}

// RemoveReserve deletes an asset from the token reserves.
func (i *SYN1000Index) RemoveReserve(id TokenID, asset string) error {
	t, err := i.Token(id)
	if err != nil {
		return err
	}
	t.RemoveReserve(asset)
	i.emit(IndexEvent{Type: IndexEventUpdated, TokenID: id, Timestamp: i.clock(), Metadata: map[string]string{"action": "remove_reserve", "asset": asset}})
	return nil
}

// SetReservePrice updates price information for an asset backing a token.
func (i *SYN1000Index) SetReservePrice(id TokenID, asset string, price *big.Rat) error {
	t, err := i.Token(id)
	if err != nil {
		return err
	}
	if err := t.SetReservePrice(asset, price); err != nil {
		return err
	}
	i.emit(IndexEvent{Type: IndexEventUpdated, TokenID: id, Timestamp: i.clock(), Metadata: map[string]string{"action": "set_price", "asset": asset}})
	return nil
}

// TotalValue returns the current total reserve value for the token.
func (i *SYN1000Index) TotalValue(id TokenID) (*big.Rat, error) {
	t, err := i.Token(id)
	if err != nil {
		return nil, err
	}
	return t.TotalReserveValue(), nil
}

// Collateralization returns the collateralization ratio for the token.
func (i *SYN1000Index) Collateralization(id TokenID) (*big.Rat, error) {
	t, err := i.Token(id)
	if err != nil {
		return nil, err
	}
	return t.CollateralizationRatio(), nil
}
