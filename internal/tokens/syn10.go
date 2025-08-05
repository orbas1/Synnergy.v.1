package tokens

// SYN10Token represents a central bank digital currency pegged to fiat.
type SYN10Token struct {
	*BaseToken
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
	t.exchangeRate = rate
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
	return SYN10Info{
		Name:         t.Name(),
		Symbol:       t.Symbol(),
		Issuer:       t.issuer,
		ExchangeRate: t.exchangeRate,
		TotalSupply:  t.TotalSupply(),
	}
}
