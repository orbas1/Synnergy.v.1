package core

// StateIterator provides a minimal interface for iterating over state entries
// with a common prefix.
type StateIterator interface {
	Next() bool
	Value() []byte
}

// StateRW defines the read/write operations on chain state required by the
// CharityPool. Implementations are provided elsewhere in the codebase.
type StateRW interface {
	Transfer(from, to Address, amount uint64) error
	SetState(key []byte, value []byte)
	GetState(key []byte) ([]byte, error)
	HasState(key []byte) (bool, error)
	PrefixIterator(prefix []byte) StateIterator
	BalanceOf(addr Address) uint64
}
