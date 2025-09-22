package tokens

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// SYN10Token represents a central bank digital currency pegged to fiat. It
// embeds BaseToken for ledger operations and adds concurrency-safe management of
// issuer metadata and fiat exchange rates.
type SYN10Token struct {
	*BaseToken
	mu           sync.RWMutex
	issuer       string
	exchangeRate float64
	lastUpdate   time.Time
	history      []ExchangeRateSnapshot
}

// ExchangeRateSnapshot tracks historical rate adjustments.
type ExchangeRateSnapshot struct {
	Rate      float64
	Reason    string
	UpdatedAt time.Time
}

// NewSYN10Token initialises a SYN10 token with the given metadata.
func NewSYN10Token(id TokenID, name, symbol, issuer string, rate float64, decimals uint8) *SYN10Token {
	if rate <= 0 {
		rate = 1
	}
	now := time.Now().UTC()
	history := []ExchangeRateSnapshot{{Rate: rate, Reason: "initial", UpdatedAt: now}}
	return &SYN10Token{
		BaseToken:    NewBaseToken(id, name, symbol, decimals),
		issuer:       issuer,
		exchangeRate: rate,
		lastUpdate:   now,
		history:      history,
	}
}

// SetIssuer updates the issuer metadata.
func (t *SYN10Token) SetIssuer(issuer string) error {
	if issuer == "" {
		return errors.New("tokens: issuer required")
	}
	t.mu.Lock()
	t.issuer = issuer
	t.mu.Unlock()
	return nil
}

// Issuer returns the current issuer metadata.
func (t *SYN10Token) Issuer() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.issuer
}

// SetExchangeRate updates the fiat exchange rate for the CBDC.
func (t *SYN10Token) SetExchangeRate(rate float64, reason string) error {
	if rate <= 0 {
		return errors.New("tokens: exchange rate must be positive")
	}
	if reason == "" {
		reason = "update"
	}
	t.mu.Lock()
	t.exchangeRate = rate
	snapshot := ExchangeRateSnapshot{Rate: rate, Reason: reason, UpdatedAt: time.Now().UTC()}
	t.lastUpdate = snapshot.UpdatedAt
	t.history = append(t.history, snapshot)
	if len(t.history) > 64 {
		t.history = t.history[len(t.history)-64:]
	}
	t.mu.Unlock()
	return nil
}

// ExchangeRateHistory returns the most recent history entries up to limit. When
// limit is zero all entries are returned.
func (t *SYN10Token) ExchangeRateHistory(limit int) []ExchangeRateSnapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	hist := make([]ExchangeRateSnapshot, len(t.history))
	copy(hist, t.history)
	if limit <= 0 || limit >= len(hist) {
		return hist
	}
	return hist[len(hist)-limit:]
}

// SYN10Info summarises token configuration.
type SYN10Info struct {
	Name         string
	Symbol       string
	Issuer       string
	ExchangeRate float64
	TotalSupply  uint64
	LastUpdated  time.Time
	History      []ExchangeRateSnapshot
}

// Info returns the current token information.
func (t *SYN10Token) Info() SYN10Info {
	t.mu.RLock()
	defer t.mu.RUnlock()
	history := make([]ExchangeRateSnapshot, len(t.history))
	copy(history, t.history)
	sort.Slice(history, func(i, j int) bool {
		return history[i].UpdatedAt.Before(history[j].UpdatedAt)
	})
	return SYN10Info{
		Name:         t.Name(),
		Symbol:       t.Symbol(),
		Issuer:       t.issuer,
		ExchangeRate: t.exchangeRate,
		TotalSupply:  t.TotalSupply(),
		LastUpdated:  t.lastUpdate,
		History:      history,
	}
}
