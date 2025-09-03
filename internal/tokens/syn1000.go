package tokens

import (
	"math/big"
	"sync"
)

// ReserveAsset tracks backing assets for a stablecoin using high‑precision
// rational numbers to avoid floating point rounding errors.
type ReserveAsset struct {
	Amount *big.Rat
	Price  *big.Rat
}

// SYN1000Token represents a stablecoin backed by reserve assets.  It embeds
// the thread‑safe BaseToken and guards the reserve ledger with its own mutex so
// callers can modify reserves concurrently.
type SYN1000Token struct {
	*BaseToken
	mu       sync.RWMutex
	reserves map[string]ReserveAsset
}

// NewSYN1000Token creates a new stablecoin token.
func NewSYN1000Token(id TokenID, name, symbol string, decimals uint8) *SYN1000Token {
	return &SYN1000Token{
		BaseToken: NewBaseToken(id, name, symbol, decimals),
		reserves:  make(map[string]ReserveAsset),
	}
}

// AddReserve adds a backing asset to the reserve list. Amounts are expressed as
// rational numbers allowing callers to provide arbitrary precision values.
func (t *SYN1000Token) AddReserve(asset string, amount *big.Rat) {
	t.mu.Lock()
	defer t.mu.Unlock()
	r := t.reserves[asset]
	if r.Amount == nil {
		r.Amount = new(big.Rat)
	}
	r.Amount.Add(r.Amount, amount)
	t.reserves[asset] = r
}

// SetReservePrice updates the unit price of a reserve asset.
func (t *SYN1000Token) SetReservePrice(asset string, price *big.Rat) {
	t.mu.Lock()
	defer t.mu.Unlock()
	r := t.reserves[asset]
	if r.Price == nil {
		r.Price = new(big.Rat)
	}
	r.Price.Set(price)
	t.reserves[asset] = r
}

// TotalReserveValue calculates the total market value of all reserves and
// returns the result as a rational number.
func (t *SYN1000Token) TotalReserveValue() *big.Rat {
	t.mu.RLock()
	defer t.mu.RUnlock()
	total := new(big.Rat)
	for _, r := range t.reserves {
		if r.Amount != nil && r.Price != nil {
			total.Add(total, new(big.Rat).Mul(r.Amount, r.Price))
		}
	}
	return total
}
