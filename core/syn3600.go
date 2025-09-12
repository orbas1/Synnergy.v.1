package core

import (
	"errors"
	"math"
	"math/big"
	"sync"
	"time"
)

// FuturesContract defines the essential metadata of a futures contract.
type FuturesContract struct {
	mu         sync.RWMutex
	Underlying string
	Quantity   uint64
	Price      uint64 // entry price per unit
	Expiration time.Time
	Settled    bool
}

// NewFuturesContract creates a new futures contract.
func NewFuturesContract(underlying string, quantity, price uint64, expiration time.Time) *FuturesContract {
	return &FuturesContract{Underlying: underlying, Quantity: quantity, Price: price, Expiration: expiration}
}

// IsExpired returns true if the contract has reached expiration.
func (f *FuturesContract) IsExpired(now time.Time) bool {
	f.mu.RLock()
	exp := f.Expiration
	f.mu.RUnlock()
	return !now.Before(exp)
}

// Settle marks the contract settled and returns PnL for the long side.
// Subsequent settlement attempts return an error to guard against double
// settlement in concurrent environments.
func (f *FuturesContract) Settle(marketPrice uint64) (int64, error) {
	f.mu.Lock()
	if f.Settled {
		f.mu.Unlock()
		return 0, errors.New("contract already settled")
	}
	f.Settled = true
	f.mu.Unlock()
	m := new(big.Int).SetUint64(marketPrice)
	p := new(big.Int).SetUint64(f.Price)
	q := new(big.Int).SetUint64(f.Quantity)
	diff := new(big.Int).Sub(m, p)
	pnl := new(big.Int).Mul(diff, q)
	if pnl.IsInt64() {
		return pnl.Int64(), nil
	}
	if pnl.Sign() >= 0 {
		return math.MaxInt64, nil
	}
	return math.MinInt64, nil
}
