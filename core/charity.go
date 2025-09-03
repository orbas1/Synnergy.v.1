package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// CharityCategory enumerates supported charity categories.
type CharityCategory uint8

const (
	HungerRelief CharityCategory = iota + 1
	ChildrenHelp
	WildlifeHelp
	SeaSupport
	DisasterSupport
	WarSupport
)

func (c CharityCategory) String() string {
	switch c {
	case HungerRelief:
		return "HungerRelief"
	case ChildrenHelp:
		return "ChildrenHelp"
	case WildlifeHelp:
		return "WildlifeHelp"
	case SeaSupport:
		return "SeaSupport"
	case DisasterSupport:
		return "DisasterSupport"
	case WarSupport:
		return "WarSupport"
	default:
		return "Unknown"
	}
}

// CharityRegistration captures a charity's registration information.
type CharityRegistration struct {
	Addr      Address
	Name      string
	Category  CharityCategory
	Cycle     uint64
	VoteCount uint64
}

// electorate defines the minimal behaviour required for voting rights.
type electorate interface {
	IsIDTokenHolder(addr Address) bool
}

// CharityPool coordinates charity registrations, voting and payouts.
type CharityPool struct {
	mu      sync.Mutex
	logger  *logrus.Logger
	led     StateRW
	vote    electorate
	genesis time.Time
}

var (
	CharityPoolAccount     Address = "charity_pool"
	InternalCharityAccount Address = "internal_charity"
)

// NewCharityPool creates a new pool instance.
func NewCharityPool(lg *logrus.Logger, led StateRW, el electorate, genesis time.Time) *CharityPool {
	return &CharityPool{logger: lg, led: led, vote: el, genesis: genesis}
}

// Deposit transfers funds into the charity pool account.
func (cp *CharityPool) Deposit(from Address, amount uint64) error {
	if cp.led == nil {
		return fmt.Errorf("ledger not configured")
	}
	if amount == 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if cp.led.BalanceOf(from) < amount {
		return fmt.Errorf("insufficient balance")
	}
	return cp.led.Transfer(from, CharityPoolAccount, amount)
}

// Register is a stub that records a charity registration.
func (cp *CharityPool) Register(addr Address, name string, cat CharityCategory) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if name == "" {
		return fmt.Errorf("name required")
	}
	key := []byte(fmt.Sprintf("charity:reg:%s", addr.Hex()))
	reg := CharityRegistration{Addr: addr, Name: name, Category: cat}
	cp.led.SetState(key, mustJSON(reg))
	return nil
}

// Vote allows a voter to vote for a charity. The current implementation is a no-op
// aside from basic persistence.
func (cp *CharityPool) Vote(voter, charity Address) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.vote != nil && !cp.vote.IsIDTokenHolder(voter) {
		return fmt.Errorf("ineligible voter")
	}
	key := []byte(fmt.Sprintf("charity:vote:%s:%s", voter.Hex(), charity.Hex()))
	cp.led.SetState(key, []byte(charity.Hex()))
	return nil
}

// Tick performs periodic maintenance such as payouts. For now it is a no-op.
func (cp *CharityPool) Tick(ts time.Time) {}

// Winners returns the list of winning charities for a cycle. This simplified
// implementation returns an empty slice.
func (cp *CharityPool) Winners(cycle uint64) ([]Address, error) {
	return []Address{}, nil
}

// GetRegistration retrieves a charity's registration for a given cycle.
func (cp *CharityPool) GetRegistration(cycle uint64, addr Address) (CharityRegistration, bool, error) {
	key := []byte(fmt.Sprintf("charity:reg:%s", addr.Hex()))
	b, err := cp.led.GetState(key)
	if err != nil || len(b) == 0 {
		return CharityRegistration{}, false, err
	}
	var reg CharityRegistration
	if err := json.Unmarshal(b, &reg); err != nil {
		return CharityRegistration{}, false, err
	}
	return reg, true, nil
}

// mustJSON encodes v to JSON and panics on failure.
func mustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
