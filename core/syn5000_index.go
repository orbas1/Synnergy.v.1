package core

// GamblingToken exposes methods of the SYN5000 token.
type GamblingToken interface {
	PlaceBet(bettor string, amount uint64, odds float64, game string) uint64
	ResolveBet(betID uint64, win bool) (uint64, error)
	GetBet(betID uint64) (*BetRecord, bool)
}
