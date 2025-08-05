package core

import (
	"math"
	"math/big"
	"time"
)

// FuturesContract defines the essential metadata of a futures contract.
type FuturesContract struct {
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
	return !now.Before(f.Expiration)
}

// Settle marks the contract settled and returns PnL for the long side.
func (f *FuturesContract) Settle(marketPrice uint64) int64 {
	f.Settled = true
	m := new(big.Int).SetUint64(marketPrice)
	p := new(big.Int).SetUint64(f.Price)
	q := new(big.Int).SetUint64(f.Quantity)
	diff := new(big.Int).Sub(m, p)
	pnl := new(big.Int).Mul(diff, q)
	if pnl.IsInt64() {
		return pnl.Int64()
	}
	if pnl.Sign() >= 0 {
		return math.MaxInt64
	}
	return math.MinInt64
}
