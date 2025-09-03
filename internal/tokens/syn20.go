package tokens

import (
	"fmt"
	"sync"
)

// SYN20Token extends BaseToken with pause and freeze capabilities.
type SYN20Token struct {
	*BaseToken
	mu     sync.RWMutex
	paused bool
	frozen map[string]bool
}

// NewSYN20Token creates a new SYN20 token.
func NewSYN20Token(id TokenID, name, symbol string, decimals uint8) *SYN20Token {
	return &SYN20Token{
		BaseToken: NewBaseToken(id, name, symbol, decimals),
		frozen:    make(map[string]bool),
	}
}

// Pause halts all transfer, mint and burn operations.
func (t *SYN20Token) Pause() {
	t.mu.Lock()
	t.paused = true
	t.mu.Unlock()
}

// Unpause resumes operations.
func (t *SYN20Token) Unpause() {
	t.mu.Lock()
	t.paused = false
	t.mu.Unlock()
}

// Freeze prevents an address from participating in transfers.
func (t *SYN20Token) Freeze(addr string) {
	t.mu.Lock()
	t.frozen[addr] = true
	t.mu.Unlock()
}

// Unfreeze lifts restrictions on an address.
func (t *SYN20Token) Unfreeze(addr string) {
	t.mu.Lock()
	delete(t.frozen, addr)
	t.mu.Unlock()
}

// Transfer overrides BaseToken.Transfer adding pause/freeze checks.
func (t *SYN20Token) Transfer(from, to string, amount uint64) error {
	t.mu.RLock()
	paused := t.paused
	fromFrozen := t.frozen[from]
	toFrozen := t.frozen[to]
	t.mu.RUnlock()
	if paused {
		return fmt.Errorf("token transfers are paused")
	}
	if fromFrozen || toFrozen {
		return fmt.Errorf("address frozen")
	}
	return t.BaseToken.Transfer(from, to, amount)
}

// Mint creates tokens if operations are not paused or frozen.
func (t *SYN20Token) Mint(to string, amount uint64) error {
	t.mu.RLock()
	paused := t.paused
	frozen := t.frozen[to]
	t.mu.RUnlock()
	if paused || frozen {
		return fmt.Errorf("minting restricted")
	}
	return t.BaseToken.Mint(to, amount)
}

// Burn destroys tokens if operations are allowed.
func (t *SYN20Token) Burn(from string, amount uint64) error {
	t.mu.RLock()
	paused := t.paused
	frozen := t.frozen[from]
	t.mu.RUnlock()
	if paused || frozen {
		return fmt.Errorf("burning restricted")
	}
	return t.BaseToken.Burn(from, amount)
}
