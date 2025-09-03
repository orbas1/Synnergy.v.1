package tokens

import "sync"

// SYN10Token represents a central bank digital currency pegged to fiat. It
// embeds BaseToken for ledger operations and adds concurrency-safe management of
// issuer metadata and fiat exchange rates.
type SYN10Token struct {
	*BaseToken
	mu           sync.RWMutex
	issuer       string
	exchangeRate float64
}

// NewSYN10Token initialises a SYN10 token with the given metadata.
func NewSYN10Token(id TokenID, name, symbol, issuer string, rate float64, decimals uint8) *SYN10Token {
	return &SYN10Token{
		BaseToken:    NewBaseToken(id, name, symbol, decimals),
		issuer:       issuer,
		exchangeRate: rate,
	}
}

// SetExchangeRate updates the fiat exchange rate for the CBDC.
func (t *SYN10Token) SetExchangeRate(rate float64) {
	t.mu.Lock()
	t.exchangeRate = rate
	t.mu.Unlock()
}

// SYN10Info summarises token configuration.
type SYN10Info struct {
	Name         string
	Symbol       string
	Issuer       string
	ExchangeRate float64
	TotalSupply  uint64
}

// Info returns the current token information.
func (t *SYN10Token) Info() SYN10Info {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return SYN10Info{
		Name:         t.Name(),
		Symbol:       t.Symbol(),
		Issuer:       t.issuer,
		ExchangeRate: t.exchangeRate,
		TotalSupply:  t.TotalSupply(),
	}
}
