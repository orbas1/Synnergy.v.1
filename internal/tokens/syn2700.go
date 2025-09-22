package tokens

import (
	"errors"
	"math/big"
	"sort"
	"sync"
	"time"
)

// ErrHolderNotFound is returned when the requested holder does not exist.
var ErrHolderNotFound = errors.New("tokens: holder not found")

// ErrInvalidAmount indicates that a provided amount is zero.
var ErrInvalidAmount = errors.New("tokens: amount must be greater than zero")

// ErrNoHolders is returned when an operation requires at least one registered holder.
var ErrNoHolders = errors.New("tokens: no holders registered")

// DividendResult captures the outcome of a dividend distribution including any
// unallocated remainder caused by integer rounding.
type DividendResult struct {
	Amounts     map[string]uint64
	Distributed uint64
	Unallocated uint64
	TotalSupply uint64
}

// DividendRecord stores a historical dividend distribution with the timestamp
// of when it was calculated.
type DividendRecord struct {
	DividendResult
	Timestamp time.Time
}

// SYN2700Token distributes dividends to registered holders based on their share
// of the total supply. It is concurrency-safe for use by multiple goroutines and
// keeps an auditable ledger of dividend calculations.
type SYN2700Token struct {
	mu      sync.RWMutex
	holders map[string]uint64
	total   uint64
	history []DividendRecord
}

// NewSYN2700Token initialises an empty dividend token.
func NewSYN2700Token() *SYN2700Token {
	return &SYN2700Token{holders: make(map[string]uint64)}
}

// AddHolder registers balance for an address. The total supply is increased
// accordingly. Zero amounts are rejected so upstream accounting mistakes do not
// silently pass through.
func (t *SYN2700Token) AddHolder(addr string, amount uint64) error {
	if amount == 0 {
		return ErrInvalidAmount
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.holders[addr] += amount
	t.total += amount
	return nil
}

// UpdateHolder sets the balance for a holder, adjusting the total supply to
// reflect the new value.
func (t *SYN2700Token) UpdateHolder(addr string, amount uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	current, ok := t.holders[addr]
	if !ok {
		return ErrHolderNotFound
	}
	t.holders[addr] = amount
	if amount > current {
		t.total += amount - current
	} else {
		t.total -= current - amount
	}
	return nil
}

// RemoveHolder removes a holder from the register.
func (t *SYN2700Token) RemoveHolder(addr string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	bal, ok := t.holders[addr]
	if !ok {
		return ErrHolderNotFound
	}
	delete(t.holders, addr)
	t.total -= bal
	return nil
}

// HolderBalance returns the balance of an address and whether it exists.
func (t *SYN2700Token) HolderBalance(addr string) (uint64, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	bal, ok := t.holders[addr]
	return bal, ok
}

// TotalSupply returns the aggregate balance of all holders.
func (t *SYN2700Token) TotalSupply() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.total
}

// Snapshot returns a copy of the holders map for inspection or reporting.
func (t *SYN2700Token) Snapshot() map[string]uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	cp := make(map[string]uint64, len(t.holders))
	for addr, bal := range t.holders {
		cp[addr] = bal
	}
	return cp
}

// Distribute splits the dividend across all holders proportionally and returns
// only the distribution map for backwards compatibility.
func (t *SYN2700Token) Distribute(dividend uint64) map[string]uint64 {
	res, _ := t.DistributeDetailed(dividend)
	return res.Amounts
}

// DistributeDetailed performs a pro-rata distribution using the largest
// remainder method to fairly allocate any remainder. The resulting allocation is
// recorded in the internal history for auditability.
func (t *SYN2700Token) DistributeDetailed(dividend uint64) (DividendResult, error) {
	holders := t.Snapshot()
	total := uint64(0)
	for _, bal := range holders {
		total += bal
	}
	if total == 0 {
		return DividendResult{Amounts: map[string]uint64{}, Unallocated: dividend}, ErrNoHolders
	}

	allocations := make(map[string]uint64, len(holders))
	type remainder struct {
		addr      string
		remainder uint64
	}
	remainders := make([]remainder, 0, len(holders))

	divBig := new(big.Int).SetUint64(dividend)
	totalBig := new(big.Int).SetUint64(total)

	var distributed uint64
	for addr, bal := range holders {
		balBig := new(big.Int).SetUint64(bal)
		product := new(big.Int).Mul(divBig, balBig)
		quotient := new(big.Int)
		rem := new(big.Int)
		quotient.QuoRem(product, totalBig, rem)
		amount := quotient.Uint64()
		allocations[addr] = amount
		distributed += amount
		remainders = append(remainders, remainder{addr: addr, remainder: rem.Uint64()})
	}

	leftover := dividend - distributed
	sort.SliceStable(remainders, func(i, j int) bool {
		if remainders[i].remainder == remainders[j].remainder {
			return remainders[i].addr < remainders[j].addr
		}
		return remainders[i].remainder > remainders[j].remainder
	})

	for i := 0; i < len(remainders) && leftover > 0; i++ {
		if remainders[i].remainder == 0 {
			break
		}
		allocations[remainders[i].addr]++
		leftover--
		distributed++
	}

	result := DividendResult{
		Amounts:     allocations,
		Distributed: distributed,
		Unallocated: leftover,
		TotalSupply: total,
	}

	record := DividendRecord{DividendResult: cloneDividendResult(result), Timestamp: time.Now()}
	t.mu.Lock()
	t.history = append(t.history, record)
	t.mu.Unlock()

	return result, nil
}

// History returns the most recent dividend records up to the specified limit.
func (t *SYN2700Token) History(limit int) []DividendRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if limit <= 0 || limit > len(t.history) {
		limit = len(t.history)
	}
	res := make([]DividendRecord, limit)
	for i := 0; i < limit; i++ {
		res[i] = cloneDividendRecord(t.history[len(t.history)-1-i])
	}
	return res
}

// Reset removes all holders and history.
func (t *SYN2700Token) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.holders = make(map[string]uint64)
	t.total = 0
	t.history = nil
}

func cloneDividendResult(res DividendResult) DividendResult {
	cp := make(map[string]uint64, len(res.Amounts))
	for k, v := range res.Amounts {
		cp[k] = v
	}
	res.Amounts = cp
	return res
}

func cloneDividendRecord(rec DividendRecord) DividendRecord {
	cloned := DividendRecord{
		DividendResult: cloneDividendResult(rec.DividendResult),
		Timestamp:      rec.Timestamp,
	}
	return cloned
}
