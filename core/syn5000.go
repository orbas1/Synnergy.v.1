package core

import "errors"

// BetRecord stores betting activity for SYN5000 tokens.
type BetRecord struct {
	ID       uint64
	Bettor   string
	Amount   uint64
	Odds     float64
	Game     string
	Resolved bool
	Won      bool
}

// SYN5000Token implements the GamblingToken interface.
type SYN5000Token struct {
	Name     string
	Symbol   string
	Decimals uint8

	nextBetID uint64
	bets      map[uint64]*BetRecord
}

// NewSYN5000Token creates a new gambling token instance.
func NewSYN5000Token(name, symbol string, decimals uint8) *SYN5000Token {
	return &SYN5000Token{Name: name, Symbol: symbol, Decimals: decimals, bets: make(map[uint64]*BetRecord)}
}

// PlaceBet records a new bet and returns its ID.
func (t *SYN5000Token) PlaceBet(bettor string, amount uint64, odds float64, game string) uint64 {
	t.nextBetID++
	id := t.nextBetID
	t.bets[id] = &BetRecord{ID: id, Bettor: bettor, Amount: amount, Odds: odds, Game: game}
	return id
}

// ResolveBet resolves a bet and returns the payout if won.
func (t *SYN5000Token) ResolveBet(betID uint64, win bool) (uint64, error) {
	b, ok := t.bets[betID]
	if !ok {
		return 0, errors.New("bet not found")
	}
	if b.Resolved {
		return 0, errors.New("bet already resolved")
	}
	b.Resolved = true
	b.Won = win
	if win {
		return uint64(float64(b.Amount) * b.Odds), nil
	}
	return 0, nil
}

// GetBet returns a bet record by ID.
func (t *SYN5000Token) GetBet(betID uint64) (*BetRecord, bool) {
	b, ok := t.bets[betID]
	return b, ok
}

var _ GamblingToken = (*SYN5000Token)(nil)
