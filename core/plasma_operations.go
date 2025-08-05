package core

import "fmt"

// Deposit records a deposit into the Plasma bridge. For this simplified
// implementation deposits are acknowledged without additional state.
func (b *PlasmaBridge) Deposit(owner, token string, amount uint64) error {
	if b.Status() {
		return fmt.Errorf("plasma bridge paused")
	}
	// Deposits would normally update on-chain state; here we just accept them.
	return nil
}

// StartExit initiates an exit and returns its nonce.
func (b *PlasmaBridge) StartExit(owner, token string, amount uint64) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.paused {
		return 0, fmt.Errorf("plasma bridge paused")
	}
	b.seq++
	nonce := b.seq
	b.exits[nonce] = &PlasmaExit{Nonce: nonce, Owner: owner, Token: token, Amount: amount}
	return nonce, nil
}

// FinalizeExit finalizes a pending exit.
func (b *PlasmaBridge) FinalizeExit(nonce uint64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	exit, ok := b.exits[nonce]
	if !ok {
		return fmt.Errorf("exit %d not found", nonce)
	}
	exit.Finalized = true
	return nil
}

// GetExit retrieves details of an exit.
func (b *PlasmaBridge) GetExit(nonce uint64) (*PlasmaExit, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	exit, ok := b.exits[nonce]
	if !ok {
		return nil, fmt.Errorf("exit %d not found", nonce)
	}
	return exit, nil
}

// ListExits lists exits initiated by an owner. If owner is empty all exits are returned.
func (b *PlasmaBridge) ListExits(owner string) []*PlasmaExit {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := []*PlasmaExit{}
	for _, ex := range b.exits {
		if owner == "" || ex.Owner == owner {
			out = append(out, ex)
		}
	}
	return out
}
