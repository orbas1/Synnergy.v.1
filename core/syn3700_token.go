package core

import (
	"errors"
	"sync"
)

// IndexComponent defines a single asset within an index token.
type IndexComponent struct {
	Token  string
	Weight float64
}

// SYN3700Token aggregates multiple assets into a single index token.
type SYN3700Token struct {
	mu         sync.RWMutex
	Name       string
	Symbol     string
	Components []IndexComponent
}

// NewSYN3700Token creates a new empty index token.
func NewSYN3700Token(name, symbol string) *SYN3700Token {
	return &SYN3700Token{Name: name, Symbol: symbol}
}

// AddComponent adds an asset and its weight to the index.
func (t *SYN3700Token) AddComponent(token string, weight float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Components = append(t.Components, IndexComponent{Token: token, Weight: weight})
}

// RemoveComponent removes an asset from the index by token symbol.
func (t *SYN3700Token) RemoveComponent(token string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, c := range t.Components {
		if c.Token == token {
			t.Components = append(t.Components[:i], t.Components[i+1:]...)
			return nil
		}
	}
	return errors.New("component not found")
}

// ListComponents returns a snapshot of the current index components.
func (t *SYN3700Token) ListComponents() []IndexComponent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]IndexComponent, len(t.Components))
	copy(comps, t.Components)
	return comps
}

// Value computes the weighted index value using the provided price map.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var sum float64
	for _, c := range t.Components {
		price := prices[c.Token]
		sum += price * c.Weight
	}
	return sum
}
