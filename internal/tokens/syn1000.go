package tokens

import (
	"errors"
	"math/big"
	"sync"
	"time"
)

// ReserveAsset tracks backing assets for a stablecoin using high-precision
// rational numbers to avoid floating point rounding errors.
type ReserveAsset struct {
	Amount    *big.Rat
	Price     *big.Rat
	UpdatedAt time.Time
}

// SYN1000Token represents a stablecoin backed by reserve assets. It embeds the
// thread-safe BaseToken and guards the reserve ledger with its own mutex so
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
func (t *SYN1000Token) AddReserve(asset string, amount *big.Rat) error {
	if asset == "" {
		return errors.New("tokens: asset identifier required")
	}
	if amount == nil || amount.Sign() <= 0 {
		return errors.New("tokens: reserve amount must be positive")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	entry := t.reserves[asset]
	if entry.Amount == nil {
		entry.Amount = new(big.Rat)
	}
	entry.Amount.Add(entry.Amount, new(big.Rat).Set(amount))
	entry.UpdatedAt = time.Now().UTC()
	t.reserves[asset] = entry
	return nil
}

// RemoveReserve removes a backing asset entirely.
func (t *SYN1000Token) RemoveReserve(asset string) {
	if asset == "" {
		return
	}
	t.mu.Lock()
	delete(t.reserves, asset)
	t.mu.Unlock()
}

// SetReservePrice updates the unit price of a reserve asset.
func (t *SYN1000Token) SetReservePrice(asset string, price *big.Rat) error {
	if asset == "" {
		return errors.New("tokens: asset identifier required")
	}
	if price == nil || price.Sign() <= 0 {
		return errors.New("tokens: reserve price must be positive")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	entry := t.reserves[asset]
	entry.Price = new(big.Rat).Set(price)
	entry.UpdatedAt = time.Now().UTC()
	t.reserves[asset] = entry
	return nil
}

// ReserveBreakdown returns a snapshot of all reserves with defensive copies.
func (t *SYN1000Token) ReserveBreakdown() map[string]ReserveAsset {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[string]ReserveAsset, len(t.reserves))
	for asset, res := range t.reserves {
		copy := ReserveAsset{UpdatedAt: res.UpdatedAt}
		if res.Amount != nil {
			copy.Amount = new(big.Rat).Set(res.Amount)
		}
		if res.Price != nil {
			copy.Price = new(big.Rat).Set(res.Price)
		}
		out[asset] = copy
	}
	return out
}

// TotalReserveValue calculates the total market value of all reserves and
// returns the result as a rational number.
func (t *SYN1000Token) TotalReserveValue() *big.Rat {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.totalValueLockedLocked()
}

// CollateralizationRatio reports the ratio between reserve value and circulating supply.
func (t *SYN1000Token) CollateralizationRatio() *big.Rat {
	t.mu.RLock()
	defer t.mu.RUnlock()
	supply := t.TotalSupply()
	if supply == 0 {
		return new(big.Rat)
	}
	value := t.totalValueLockedLocked()
	return new(big.Rat).Quo(value, new(big.Rat).SetUint64(supply))
}

func (t *SYN1000Token) totalValueLockedLocked() *big.Rat {
	total := new(big.Rat)
	for _, r := range t.reserves {
		if r.Amount != nil && r.Price != nil {
			total.Add(total, new(big.Rat).Mul(r.Amount, r.Price))
		}
	}
	return total
}
