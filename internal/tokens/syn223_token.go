package tokens

import (
	"errors"
	"sync"
	"time"
)

// SYN223Token implements a secure transfer token with whitelist and blacklist controls.
type SYN223Token struct {
	mu        sync.RWMutex
	Name      string
	Symbol    string
	Owner     string
	balances  map[string]uint64
	whitelist map[string]bool
	blacklist map[string]bool
	hooks     map[string]ReceiveHook
	metadata  map[string]map[string]string
	events    []SYN223Event
}

// NewSYN223Token creates a new SYN223 token and assigns the initial supply to the owner.
func NewSYN223Token(name, symbol, owner string, supply uint64) *SYN223Token {
	t := &SYN223Token{
		Name:      name,
		Symbol:    symbol,
		Owner:     owner,
		balances:  map[string]uint64{owner: supply},
		whitelist: make(map[string]bool),
		blacklist: make(map[string]bool),
		hooks:     make(map[string]ReceiveHook),
		metadata:  make(map[string]map[string]string),
	}
	t.whitelist[owner] = true
	return t
}

// AddToWhitelist authorises an address to receive tokens.
func (t *SYN223Token) AddToWhitelist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.whitelist[addr] = true
}

// RemoveFromWhitelist removes an address from the whitelist.
func (t *SYN223Token) RemoveFromWhitelist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.whitelist, addr)
}

// AddToBlacklist blocks an address from participating in transfers.
func (t *SYN223Token) AddToBlacklist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.blacklist[addr] = true
}

// RemoveFromBlacklist lifts a previously applied blacklist restriction.
func (t *SYN223Token) RemoveFromBlacklist(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.blacklist, addr)
}

// SetMetadata stores structured metadata for regulatory reporting.
func (t *SYN223Token) SetMetadata(addr string, data map[string]string) {
	t.mu.Lock()
	cp := make(map[string]string, len(data))
	for k, v := range data {
		cp[k] = v
	}
	t.metadata[addr] = cp
	t.mu.Unlock()
}

// Metadata returns a defensive copy of address metadata.
func (t *SYN223Token) Metadata(addr string) map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	meta := t.metadata[addr]
	cp := make(map[string]string, len(meta))
	for k, v := range meta {
		cp[k] = v
	}
	return cp
}

// ReceiveHook executes contract callbacks similar to ERC223 semantics.
type ReceiveHook func(from string, amount uint64, data []byte) error

// RegisterHook associates a receiver hook with an address.
func (t *SYN223Token) RegisterHook(addr string, hook ReceiveHook) {
	t.mu.Lock()
	if hook == nil {
		delete(t.hooks, addr)
	} else {
		t.hooks[addr] = hook
	}
	t.mu.Unlock()
}

// SYN223Event records transfers for monitoring dashboards.
type SYN223Event struct {
	From      string
	To        string
	Amount    uint64
	Memo      string
	Timestamp time.Time
}

// Transfer performs a safe token transfer verifying whitelist and blacklist rules.
func (t *SYN223Token) Transfer(from, to string, amount uint64) error {
	return t.TransferWithData(from, to, amount, nil, "")
}

// TransferWithData extends Transfer to include payload metadata and hooks.
func (t *SYN223Token) TransferWithData(from, to string, amount uint64, data []byte, memo string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.blacklist[from] || t.blacklist[to] {
		return errors.New("address blacklisted")
	}
	if !t.whitelist[to] {
		return errors.New("recipient not whitelisted")
	}
	bal := t.balances[from]
	if bal < amount {
		return errors.New("insufficient balance")
	}
	t.balances[from] = bal - amount
	t.balances[to] += amount
	if hook := t.hooks[to]; hook != nil {
		if err := hook(from, amount, data); err != nil {
			t.balances[from] = bal
			t.balances[to] -= amount
			return err
		}
	}
	t.events = append(t.events, SYN223Event{From: from, To: to, Amount: amount, Memo: memo, Timestamp: time.Now()})
	return nil
}

// BalanceOf returns the current balance of an address.
func (t *SYN223Token) BalanceOf(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.balances[addr]
}

// Events returns recent transfer events up to the provided limit.
func (t *SYN223Token) Events(limit int) []SYN223Event {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if limit <= 0 || limit >= len(t.events) {
		out := make([]SYN223Event, len(t.events))
		copy(out, t.events)
		return out
	}
	out := make([]SYN223Event, limit)
	copy(out, t.events[len(t.events)-limit:])
	return out
}
