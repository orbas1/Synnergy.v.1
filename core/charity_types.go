package core

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Address represents a simple hexadecimal account identifier.
type Address string

// Hex returns the hexadecimal string form of the address.
func (a Address) Hex() string { return string(a) }

// Bytes returns the binary representation of the address.
func (a Address) Bytes() []byte {
	s := strings.TrimPrefix(string(a), "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return []byte{}
	}
	return b
}

// Short returns a shortened representation useful for logs.
func (a Address) Short() string {
	s := a.Hex()
	if len(s) <= 10 {
		return s
	}
	return s[:6] + "..." + s[len(s)-4:]
}

// StringToAddress converts a hex string into an Address type.
func StringToAddress(s string) (Address, error) {
	if s == "" {
		return "", errors.New("empty address")
	}
	if !strings.HasPrefix(s, "0x") {
		return "", errors.New("invalid address format")
	}
	if _, err := hex.DecodeString(s[2:]); err != nil {
		return "", err
	}
	return Address(strings.ToLower(s)), nil
}

// Hash is a 32-byte identifier used for cycle hashing.
type Hash [32]byte

// StateIterator iterates over key/value pairs with a prefix.
type StateIterator interface {
	Next() bool
	Value() []byte
}

// StateRW abstracts state reads and writes required by CharityPool.
type StateRW interface {
	Transfer(from, to Address, amount uint64) error
	BalanceOf(addr Address) uint64
	SetState(key, value []byte)
	GetState(key []byte) ([]byte, error)
	HasState(key []byte) (bool, error)
	PrefixIterator(prefix []byte) StateIterator
}

// CharityPool holds charity mechanics state and dependencies.
type CharityPool struct {
	logger    *logrus.Logger
	led       StateRW
	vote      electorate
	genesis   time.Time
	lastDaily int64
	mu        sync.Mutex
}

// CharityRegistration represents a single charity entry.
type CharityRegistration struct {
	Addr      Address
	Name      string
	Category  CharityCategory
	Cycle     uint64
	VoteCount int
}

// mustJSON marshals v to JSON and panics on error.
func mustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// voteKey computes the ledger key for a voter in a given cycle.
func voteKey(cycle Hash, voter Address) []byte {
	return []byte(fmt.Sprintf("charity:vote:%x:%s", cycle[:], voter.Hex()))
}
