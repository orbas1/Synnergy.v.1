package core

// StateIterator provides sequential access over key/value pairs stored in the
// ledger. It is deliberately minimal for test compilation purposes.
type StateIterator interface {
	Next() bool
	Value() []byte
}

// StateRW abstracts the ledger access that CharityPool depends on. Only the
// methods required by charity.go are included here.
type StateRW interface {
	Transfer(from, to Address, amount uint64) error
	SetState(key, value []byte)
	GetState(key []byte) ([]byte, error)
	HasState(key []byte) (bool, error)
	PrefixIterator(prefix []byte) StateIterator
	BalanceOf(addr Address) uint64
}
