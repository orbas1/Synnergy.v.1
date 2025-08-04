package tokens

import "time"

// SYN12Metadata defines treasury bill properties for CBDC instruments.
type SYN12Metadata struct {
	BillID    string
	Issuer    string
	IssueDate time.Time
	Maturity  time.Time
	Discount  float64
	FaceValue uint64
}

// SYN12Token represents a tokenised treasury bill.
type SYN12Token struct {
	*BaseToken
	Metadata SYN12Metadata
}

// NewSYN12Token creates a new SYN12 token with the provided metadata.
func NewSYN12Token(id TokenID, name, symbol string, meta SYN12Metadata, decimals uint8) *SYN12Token {
	return &SYN12Token{
		BaseToken: NewBaseToken(id, name, symbol, decimals),
		Metadata:  meta,
	}
}
