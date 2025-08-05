package tokens

// ReserveAsset tracks backing assets for a stablecoin.
type ReserveAsset struct {
	Amount float64
	Price  float64
}

// SYN1000Token represents a stablecoin backed by reserve assets.
type SYN1000Token struct {
	*BaseToken
	reserves map[string]ReserveAsset
}

// NewSYN1000Token creates a new stablecoin token.
func NewSYN1000Token(id TokenID, name, symbol string, decimals uint8) *SYN1000Token {
	return &SYN1000Token{
		BaseToken: NewBaseToken(id, name, symbol, decimals),
		reserves:  make(map[string]ReserveAsset),
	}
}

// AddReserve adds a backing asset to the reserve list.
func (t *SYN1000Token) AddReserve(asset string, amount float64) {
	r := t.reserves[asset]
	r.Amount += amount
	t.reserves[asset] = r
}

// SetReservePrice updates the unit price of a reserve asset.
func (t *SYN1000Token) SetReservePrice(asset string, price float64) {
	r := t.reserves[asset]
	r.Price = price
	t.reserves[asset] = r
}

// TotalReserveValue calculates the total market value of all reserves.
func (t *SYN1000Token) TotalReserveValue() float64 {
	var v float64
	for _, r := range t.reserves {
		v += r.Amount * r.Price
	}
	return v
}
