package tokens

import (
	"math/big"
	"sync"
	"time"
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
	mu                 sync.RWMutex
	reserves           map[string]ReserveAsset
	liquidityThreshold *big.Rat
	attestations       []ReserveAttestation
}

// NewSYN1000Token creates a new stablecoin token.
func NewSYN1000Token(id TokenID, name, symbol string, decimals uint8) *SYN1000Token {
	return &SYN1000Token{
		BaseToken:          NewBaseToken(id, name, symbol, decimals),
		reserves:           make(map[string]ReserveAsset),
		liquidityThreshold: big.NewRat(1, 1),
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

// ReserveReport exposes a consistent snapshot of reserve state for CLI, VM and
// web dashboards.
type ReserveReport struct {
	Asset  string
	Amount *big.Rat
	Price  *big.Rat
	Value  *big.Rat
}

// ReserveAttestation captures signed reports from custodians.
type ReserveAttestation struct {
	Source    string
	Reported  time.Time
	Value     *big.Rat
	Statement string
}

// SnapshotReserves returns a deep copy of reserves for analytics pipelines.
func (t *SYN1000Token) SnapshotReserves() []ReserveReport {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]ReserveReport, 0, len(t.reserves))
	for asset, res := range t.reserves {
		report := ReserveReport{Asset: asset}
		if res.Amount != nil {
			report.Amount = new(big.Rat).Set(res.Amount)
		}
		if res.Price != nil {
			report.Price = new(big.Rat).Set(res.Price)
		}
		if res.Amount != nil && res.Price != nil {
			report.Value = new(big.Rat).Mul(res.Amount, res.Price)
		}
		out = append(out, report)
	}
	return out
}

// CoverageRatio reports the ratio between reserves and circulating supply.
func (t *SYN1000Token) CoverageRatio() *big.Rat {
	supply := t.TotalSupply()
	totalValue := t.TotalReserveValue()
	if supply == 0 {
		return big.NewRat(1, 1)
	}
	return totalValue.Quo(totalValue, new(big.Rat).SetUint64(supply))
}

// SetLiquidityThreshold configures the minimum acceptable coverage ratio.
func (t *SYN1000Token) SetLiquidityThreshold(ratio *big.Rat) {
	t.mu.Lock()
	t.liquidityThreshold = new(big.Rat).Set(ratio)
	t.mu.Unlock()
}

// IsCollateralised indicates whether reserves meet the configured threshold.
func (t *SYN1000Token) IsCollateralised() bool {
	t.mu.RLock()
	threshold := new(big.Rat).Set(t.liquidityThreshold)
	t.mu.RUnlock()
	ratio := t.CoverageRatio()
	return ratio.Cmp(threshold) >= 0
}

// StressTestRedemption validates whether reserves can satisfy a redemption
// request of the given nominal amount.
func (t *SYN1000Token) StressTestRedemption(amount uint64) bool {
	reserves := t.TotalReserveValue()
	need := new(big.Rat).SetUint64(amount)
	return reserves.Cmp(need) >= 0
}

// RecordAttestation appends an externally signed reserve statement.
func (t *SYN1000Token) RecordAttestation(source, statement string, value *big.Rat) {
	t.mu.Lock()
	defer t.mu.Unlock()
	att := ReserveAttestation{Source: source, Statement: statement, Reported: time.Now()}
	if value != nil {
		att.Value = new(big.Rat).Set(value)
	}
	t.attestations = append(t.attestations, att)
}

// LatestAttestation returns the newest attestation if available.
func (t *SYN1000Token) LatestAttestation() (ReserveAttestation, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if len(t.attestations) == 0 {
		return ReserveAttestation{}, false
	}
	att := t.attestations[len(t.attestations)-1]
	if att.Value != nil {
		att.Value = new(big.Rat).Set(att.Value)
	}
	return att, true
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
