package tokens

import (
	"errors"
	"sort"
	"sync"
)

var (
	ErrComponentNotFound = errors.New("tokens: component not found")
	ErrMissingPrice      = errors.New("tokens: price missing for component")
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
func (t *SYN3700Token) AddComponent(token string, weight float64) error {
	if weight < 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Components = append(t.Components, IndexComponent{Token: token, Weight: weight})
	return nil
}

// UpdateComponent updates the weight for an existing component.
func (t *SYN3700Token) UpdateComponent(token string, weight float64) error {
	if weight < 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, c := range t.Components {
		if c.Token == token {
			t.Components[i].Weight = weight
			return nil
		}
	}
	return ErrComponentNotFound
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
	return ErrComponentNotFound
}

// ListComponents returns a snapshot of the current index components.
func (t *SYN3700Token) ListComponents() []IndexComponent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	comps := make([]IndexComponent, len(t.Components))
	copy(comps, t.Components)
	sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
	return comps
}

// NormalizeWeights scales the weights so that they sum to 1.
func (t *SYN3700Token) NormalizeWeights() {
	t.mu.Lock()
	defer t.mu.Unlock()
	var sum float64
	for _, c := range t.Components {
		sum += c.Weight
	}
	if sum == 0 {
		return
	}
	for i := range t.Components {
		t.Components[i].Weight = t.Components[i].Weight / sum
	}
}

// ComponentWeight returns the weight for a token.
func (t *SYN3700Token) ComponentWeight(token string) (float64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, c := range t.Components {
		if c.Token == token {
			return c.Weight, nil
		}
	}
	return 0, ErrComponentNotFound
}

// Value computes the weighted index value using the provided price map.
func (t *SYN3700Token) Value(prices map[string]float64) float64 {
	value, _, err := t.ValueDetailed(prices)
	if err != nil {
		return 0
	}
	return value
}

// ValueDetailed returns the index value plus per-component contributions.
func (t *SYN3700Token) ValueDetailed(prices map[string]float64) (float64, map[string]float64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	contributions := make(map[string]float64, len(t.Components))
	var sum float64
	for _, c := range t.Components {
		price, ok := prices[c.Token]
		if !ok {
			return 0, nil, ErrMissingPrice
		}
		contrib := price * c.Weight
		contributions[c.Token] = contrib
		sum += contrib
	}
	return sum, contributions, nil
}
