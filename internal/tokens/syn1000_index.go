package tokens

import (
	"fmt"
	"math/big"
	"sync"
)

// SYN1000Index manages multiple SYN1000 stablecoin instances. It protects the
// underlying map with a mutex so concurrent CLI calls remain safe.
type SYN1000Index struct {
	mu     sync.RWMutex
	tokens map[TokenID]*SYN1000Token
	next   TokenID
}

// NewSYN1000Index creates an empty index.
func NewSYN1000Index() *SYN1000Index {
	return &SYN1000Index{tokens: make(map[TokenID]*SYN1000Token)}
}

// Create instantiates a new SYN1000 token and returns its identifier.
func (i *SYN1000Index) Create(name, symbol string, decimals uint8) TokenID {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.next++
	id := i.next
	i.tokens[id] = NewSYN1000Token(id, name, symbol, decimals)
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
	t.AddReserve(asset, amount)
	return nil
}

// SetReservePrice updates price information for an asset backing a token.
func (i *SYN1000Index) SetReservePrice(id TokenID, asset string, price *big.Rat) error {
	t, err := i.Token(id)
	if err != nil {
		return err
	}
	t.SetReservePrice(asset, price)
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
