package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// StateRW defines the minimal ledger operations CharityPool relies on.
type StateRW interface {
	Transfer(from, to Address, amount uint64) error
	SetState(key, value []byte)
	GetState(key []byte) ([]byte, error)
	HasState(key []byte) (bool, error)
	BalanceOf(addr Address) uint64
	PrefixIterator(prefix []byte) Iterator
}

// Iterator abstracts iteration over key/value pairs with a shared prefix.
type Iterator interface {
	Next() bool
	Value() []byte
}

// CharityPool maintains charity registrations, votes and payouts.
type CharityPool struct {
	mu        sync.Mutex
	logger    *logrus.Logger
	led       StateRW
	vote      electorate
	genesis   time.Time
	lastDaily int64
}

// CharityRegistration holds metadata for a charity in a given cycle.
type CharityRegistration struct {
	Addr      Address
	Name      string
	Category  CharityCategory
	Cycle     uint64
	VoteCount uint64
}

// mustJSON encodes v to JSON and panics on failure.
func mustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// voteKey returns the ledger key for a voter's vote in a specific cycle hash.
func voteKey(cycle Hash, voter Address) []byte {
	return []byte(fmt.Sprintf("charity:vote:%x:%s", cycle[:], voter.Hex()))
}
