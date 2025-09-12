package core

import (
	"errors"
	"sync"
)

// Component represents a token and its weight within the index.
type Component struct {
	Token  string
	Weight float64
}

// SYN3700Token models a weighted index token.
type SYN3700Token struct {
	mu         sync.RWMutex
	Name       string
	Symbol     string
	Components map[string]float64
}

// NewSYN3700Token creates an empty index token.
func NewSYN3700Token(name, symbol string) *SYN3700Token {
	return &SYN3700Token{
		Name:       name,
		Symbol:     symbol,
		Components: make(map[string]float64),
	}
}

// AddComponent adds or updates a component token with a given weight.
func (t *SYN3700Token) AddComponent(token string, weight float64) {
	t.mu.Lock()
	t.Components[token] = weight
	t.mu.Unlock()
}

// RemoveComponent deletes a component from the index.
func (t *SYN3700Token) RemoveComponent(token string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.Components[token]; !ok {
		return errors.New("component not found")
	}
	delete(t.Components, token)
	return nil
}

// ListComponents returns a snapshot of all components.
func (t *SYN3700Token) ListComponents() []Component {
	t.mu.RLock()
	comps := make([]Component, 0, len(t.Components))
	for tok, w := range t.Components {
		comps = append(comps, Component{Token: tok, Weight: w})
	}
	t.mu.RUnlock()
	return comps
}

// Value computes the index value given a map of token prices.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var total float64
	for tok, w := range t.Components {
		total += w * prices[tok]
	}
	return total
}
