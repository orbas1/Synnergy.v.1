package tokens

import (
	"errors"
	"math"
	"math/big"
	"sync"
)

var ErrInvalidRatio = errors.New("tokens: conversion ratio must be greater than zero")

// SYN3200Token converts units using a configurable ratio. It can model
// convertible assets such as wrapped tokens.
type SYN3200Token struct {
	mu    sync.RWMutex
	ratio *big.Rat
}

// NewSYN3200Token creates a new converter with the given ratio. If the ratio is
// invalid the converter defaults to 1:1.
func NewSYN3200Token(ratio float64) *SYN3200Token {
	tok := &SYN3200Token{ratio: big.NewRat(1, 1)}
	_ = tok.SetRatio(ratio)
	return tok
}

// Convert returns the amount multiplied by the current ratio rounded using
// standard half-up rounding.
func (t *SYN3200Token) Convert(amount uint64) uint64 {
	exact := t.ConvertExact(amount)
	num := new(big.Int).Set(exact.Num())
	den := new(big.Int).Set(exact.Denom())
	if den.Sign() == 0 {
		return 0
	}
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.QuoRem(num, den, remainder)
	if remainder.Sign() != 0 {
		doubled := new(big.Int).Lsh(remainder.Abs(remainder), 1)
		if doubled.Cmp(den) >= 0 {
			quotient.Add(quotient, big.NewInt(1))
		}
	}
	return quotient.Uint64()
}

// ConvertExact returns the precise rational representation of the conversion.
func (t *SYN3200Token) ConvertExact(amount uint64) *big.Rat {
	t.mu.RLock()
	defer t.mu.RUnlock()
	product := new(big.Rat).Mul(t.ratio, new(big.Rat).SetUint64(amount))
	return new(big.Rat).Set(product)
}

// Ratio returns a snapshot of the configured ratio.
func (t *SYN3200Token) Ratio() *big.Rat {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return new(big.Rat).Set(t.ratio)
}

// SetRatio updates the conversion ratio using a floating point value.
func (t *SYN3200Token) SetRatio(r float64) error {
	if r <= 0 || math.IsNaN(r) || math.IsInf(r, 0) {
		return ErrInvalidRatio
	}
	rat := new(big.Rat)
	rat.SetFloat64(r)
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ratio = rat
	return nil
}

// SetRatioFraction updates the ratio using a numerator and denominator.
func (t *SYN3200Token) SetRatioFraction(numerator, denominator uint64) error {
	if denominator == 0 {
		return ErrInvalidRatio
	}
	rat := new(big.Rat).SetFrac(new(big.Int).SetUint64(numerator), new(big.Int).SetUint64(denominator))
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ratio = rat
	return nil
}
