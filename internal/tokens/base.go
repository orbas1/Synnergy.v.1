package tokens

import (
	"errors"
	"math"
	"sort"
	"sync"
	"time"
)

// TokenID uniquely identifies a token instance within the registry.
type TokenID uint64

// Token defines basic behaviours for all tokens.
type Token interface {
	ID() TokenID
	Name() string
	Symbol() string
	Decimals() uint8
	TotalSupply() uint64
	BalanceOf(addr string) uint64
	Transfer(from, to string, amount uint64) error
	TransferFrom(owner, spender, to string, amount uint64) error
	Mint(to string, amount uint64) error
	Burn(from string, amount uint64) error
	Approve(owner, spender string, amount uint64) error
	Allowance(owner, spender string) uint64
}

var (
	// ErrInsufficientBalance is returned when an account lacks funds for the
	// requested operation.
	ErrInsufficientBalance = errors.New("tokens: insufficient balance")
	// ErrAllowanceExceeded indicates the spender attempted to transfer more
	// than the approved allowance.
	ErrAllowanceExceeded = errors.New("tokens: allowance exceeded")
	// ErrInvalidAddress is returned when an empty address is provided to an
	// operation that requires a destination or source account.
	ErrInvalidAddress = errors.New("tokens: address required")
	// ErrAmountZero is returned when an operation attempts to move zero
	// tokens. Enforcing this avoids unnecessary hooks and telemetry noise.
	ErrAmountZero = errors.New("tokens: amount must be greater than zero")
	// ErrOverflow indicates a balance or supply would overflow uint64.
	ErrOverflow = errors.New("tokens: balance overflow detected")
	// ErrSupplyCapExceeded is returned when a mint would exceed the configured cap.
	ErrSupplyCapExceeded = errors.New("tokens: max supply exceeded")
)

// EventType enumerates lifecycle events emitted by the BaseToken.
type EventType string

const (
	EventMint     EventType = "mint"
	EventBurn     EventType = "burn"
	EventTransfer EventType = "transfer"
)

// Event captures the context of a ledger mutation and is delivered to
// registered hooks. Hooks must be fast; they are executed synchronously once
// the internal mutex is released.
type Event struct {
	Type      EventType
	TokenID   TokenID
	From      string
	To        string
	Amount    uint64
	Timestamp time.Time
	Metadata  map[string]string
}

// Hook receives ledger events. Errors should be handled internally as the base
// token will recover from panics to maintain fault tolerance.
type Hook func(Event)

type accountMetadata struct {
	lastUpdated time.Time
}

// BaseToken implements the Token interface providing basic accounting. It is
// safe for concurrent use by multiple goroutines and emits lifecycle events to
// subscribed hooks for CLI and web integrations.
type BaseToken struct {
	mu          sync.RWMutex
	id          TokenID
	name        string
	symbol      string
	decimals    uint8
	balances    map[string]uint64
	supply      uint64
	allowances  map[string]map[string]uint64
	accountInfo map[string]accountMetadata
	hooks       []Hook
	clock       func() time.Time
	maxSupply   uint64
}

// BaseTokenOption customises the behaviour of the base token at construction.
type BaseTokenOption func(*BaseToken)

// WithMaxSupply configures a hard supply cap. Passing zero leaves the token uncapped.
func WithMaxSupply(limit uint64) BaseTokenOption {
	return func(t *BaseToken) {
		t.maxSupply = limit
	}
}

// WithClock overrides the internal clock, primarily used by tests to produce
// deterministic timestamps.
func WithClock(clock func() time.Time) BaseTokenOption {
	return func(t *BaseToken) {
		if clock != nil {
			t.clock = clock
		}
	}
}

// NewBaseToken creates a new base token instance.
func NewBaseToken(id TokenID, name, symbol string, decimals uint8, opts ...BaseTokenOption) *BaseToken {
	token := &BaseToken{
		id:          id,
		name:        name,
		symbol:      symbol,
		decimals:    decimals,
		balances:    make(map[string]uint64),
		allowances:  make(map[string]map[string]uint64),
		accountInfo: make(map[string]accountMetadata),
		clock: func() time.Time {
			return time.Now().UTC()
		},
	}
	for _, opt := range opts {
		opt(token)
	}
	return token
}

// ID returns the unique identifier of the token.
func (t *BaseToken) ID() TokenID {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.id
}

// Name returns the human readable token name.
func (t *BaseToken) Name() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.name
}

// Symbol returns the token trading symbol.
func (t *BaseToken) Symbol() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.symbol
}

// Decimals returns the decimal precision for the token.
func (t *BaseToken) Decimals() uint8 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.decimals
}

// TotalSupply returns the current token supply.
func (t *BaseToken) TotalSupply() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.supply
}

// MaxSupply returns the configured supply cap, if any.
func (t *BaseToken) MaxSupply() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.maxSupply
}

// SetMaxSupply updates the supply cap. Setting the cap below the current supply
// has no effect until the circulating supply is reduced.
func (t *BaseToken) SetMaxSupply(limit uint64) {
	t.mu.Lock()
	t.maxSupply = limit
	t.mu.Unlock()
}

// BalanceOf retrieves the balance for the specified address.
func (t *BaseToken) BalanceOf(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.balances[addr]
}

// Snapshot returns metadata for a specific address including the last updated timestamp.
func (t *BaseToken) Snapshot(addr string) AccountSnapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	meta := t.accountInfo[addr]
	return AccountSnapshot{Address: addr, Balance: t.balances[addr], LastUpdated: meta.lastUpdated}
}

// Accounts returns deterministic snapshots for all tracked accounts. The list is
// sorted by address to ensure CLI and web dashboards remain stable across runs.
func (t *BaseToken) Accounts() []AccountSnapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	accounts := make([]AccountSnapshot, 0, len(t.balances))
	for addr, bal := range t.balances {
		accounts = append(accounts, AccountSnapshot{Address: addr, Balance: bal, LastUpdated: t.accountInfo[addr].lastUpdated})
	}
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Address < accounts[j].Address
	})
	return accounts
}

// AccountSnapshot summarises ledger information for an account.
type AccountSnapshot struct {
	Address     string
	Balance     uint64
	LastUpdated time.Time
}

// RegisterHook installs an event hook. Hooks are invoked synchronously after
// each ledger mutation but outside the internal mutex ensuring they cannot
// block account operations.
func (t *BaseToken) RegisterHook(h Hook) {
	if h == nil {
		return
	}
	t.mu.Lock()
	t.hooks = append(t.hooks, h)
	t.mu.Unlock()
}

// Transfer moves tokens between addresses.
func (t *BaseToken) Transfer(from, to string, amount uint64) error {
	if err := validateTransferInputs(from, to, amount); err != nil {
		return err
	}
	ts := t.clock()
	t.mu.Lock()
	if t.balances[from] < amount {
		t.mu.Unlock()
		return ErrInsufficientBalance
	}
	if math.MaxUint64-t.balances[to] < amount {
		t.mu.Unlock()
		return ErrOverflow
	}
	t.balances[from] -= amount
	t.balances[to] += amount
	t.touchLocked(from, ts)
	t.touchLocked(to, ts)
	hooks := append([]Hook(nil), t.hooks...)
	id := t.id
	t.mu.Unlock()

	t.emit(Event{Type: EventTransfer, TokenID: id, From: from, To: to, Amount: amount, Timestamp: ts}, hooks)
	return nil
}

// TransferFrom moves tokens on behalf of an owner using an approved allowance.
func (t *BaseToken) TransferFrom(owner, spender, to string, amount uint64) error {
	if err := validateTransferInputs(owner, to, amount); err != nil {
		return err
	}
	if spender == "" {
		return ErrInvalidAddress
	}
	ts := t.clock()
	t.mu.Lock()
	if t.allowances[owner][spender] < amount {
		t.mu.Unlock()
		return ErrAllowanceExceeded
	}
	if t.balances[owner] < amount {
		t.mu.Unlock()
		return ErrInsufficientBalance
	}
	if math.MaxUint64-t.balances[to] < amount {
		t.mu.Unlock()
		return ErrOverflow
	}
	t.allowances[owner][spender] -= amount
	t.balances[owner] -= amount
	t.balances[to] += amount
	t.touchLocked(owner, ts)
	t.touchLocked(to, ts)
	hooks := append([]Hook(nil), t.hooks...)
	id := t.id
	t.mu.Unlock()

	t.emit(Event{Type: EventTransfer, TokenID: id, From: owner, To: to, Amount: amount, Timestamp: ts, Metadata: map[string]string{"spender": spender}}, hooks)
	return nil
}

// Mint creates new tokens for the specified address.
func (t *BaseToken) Mint(to string, amount uint64) error {
	if err := validateAddress(to); err != nil {
		return err
	}
	if amount == 0 {
		return ErrAmountZero
	}
	ts := t.clock()
	t.mu.Lock()
	if t.maxSupply > 0 && t.supply+amount > t.maxSupply {
		t.mu.Unlock()
		return ErrSupplyCapExceeded
	}
	if math.MaxUint64-t.balances[to] < amount || math.MaxUint64-t.supply < amount {
		t.mu.Unlock()
		return ErrOverflow
	}
	t.balances[to] += amount
	t.supply += amount
	t.touchLocked(to, ts)
	hooks := append([]Hook(nil), t.hooks...)
	id := t.id
	t.mu.Unlock()

	t.emit(Event{Type: EventMint, TokenID: id, To: to, Amount: amount, Timestamp: ts}, hooks)
	return nil
}

// Burn removes tokens from the specified address.
func (t *BaseToken) Burn(from string, amount uint64) error {
	if err := validateAddress(from); err != nil {
		return err
	}
	if amount == 0 {
		return ErrAmountZero
	}
	ts := t.clock()
	t.mu.Lock()
	if t.balances[from] < amount {
		t.mu.Unlock()
		return ErrInsufficientBalance
	}
	t.balances[from] -= amount
	t.supply -= amount
	t.touchLocked(from, ts)
	hooks := append([]Hook(nil), t.hooks...)
	id := t.id
	t.mu.Unlock()

	t.emit(Event{Type: EventBurn, TokenID: id, From: from, Amount: amount, Timestamp: ts}, hooks)
	return nil
}

// Approve sets the allowance for a spender on behalf of an owner.
func (t *BaseToken) Approve(owner, spender string, amount uint64) error {
	if err := validateAddress(owner); err != nil {
		return err
	}
	if err := validateAddress(spender); err != nil {
		return err
	}
	t.mu.Lock()
	if t.allowances[owner] == nil {
		t.allowances[owner] = make(map[string]uint64)
	}
	t.allowances[owner][spender] = amount
	t.mu.Unlock()
	return nil
}

// Allowance returns the remaining approved amount a spender can use from an owner.
func (t *BaseToken) Allowance(owner, spender string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.allowances[owner][spender]
}

func validateTransferInputs(from, to string, amount uint64) error {
	if err := validateAddress(from); err != nil {
		return err
	}
	if err := validateAddress(to); err != nil {
		return err
	}
	if amount == 0 {
		return ErrAmountZero
	}
	return nil
}

func validateAddress(addr string) error {
	if addr == "" {
		return ErrInvalidAddress
	}
	return nil
}

func (t *BaseToken) touchLocked(addr string, ts time.Time) {
	if addr == "" {
		return
	}
	if _, exists := t.accountInfo[addr]; !exists {
		t.accountInfo[addr] = accountMetadata{}
	}
	info := t.accountInfo[addr]
	info.lastUpdated = ts
	t.accountInfo[addr] = info
}

func (t *BaseToken) emit(evt Event, hooks []Hook) {
	for _, h := range hooks {
		func(h Hook) {
			defer func() {
				_ = recover()
			}()
			h(evt)
		}(h)
	}
}
