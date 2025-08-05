package tokens

import "fmt"

// SYN1000Index manages multiple SYN1000 stablecoin instances.
type SYN1000Index struct {
	tokens map[TokenID]*SYN1000Token
	next   TokenID
}

// NewSYN1000Index creates an empty index.
func NewSYN1000Index() *SYN1000Index {
	return &SYN1000Index{tokens: make(map[TokenID]*SYN1000Token)}
}

// Create instantiates a new SYN1000 token and returns its identifier.
func (i *SYN1000Index) Create(name, symbol string, decimals uint8) TokenID {
	i.next++
	id := i.next
	i.tokens[id] = NewSYN1000Token(id, name, symbol, decimals)
	return id
}

// Token retrieves a SYN1000 token by ID.
func (i *SYN1000Index) Token(id TokenID) (*SYN1000Token, error) {
	t, ok := i.tokens[id]
	if !ok {
		return nil, fmt.Errorf("token not found")
	}
	return t, nil
}

// AddReserve adds reserve assets to a specific token.
func (i *SYN1000Index) AddReserve(id TokenID, asset string, amount float64) error {
	t, err := i.Token(id)
	if err != nil {
		return err
	}
	t.AddReserve(asset, amount)
	return nil
}

// SetReservePrice updates price information for an asset backing a token.
func (i *SYN1000Index) SetReservePrice(id TokenID, asset string, price float64) error {
	t, err := i.Token(id)
	if err != nil {
		return err
	}
	t.SetReservePrice(asset, price)
	return nil
}

// TotalValue returns the current total reserve value for the token.
func (i *SYN1000Index) TotalValue(id TokenID) (float64, error) {
	t, err := i.Token(id)
	if err != nil {
		return 0, err
	}
	return t.TotalReserveValue(), nil
}
