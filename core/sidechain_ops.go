package core

import "fmt"

// SidechainOps provides deposit and withdrawal helpers for side-chains.
type SidechainOps struct {
	registry *SidechainRegistry
}

// NewSidechainOps creates operations bound to a registry.
func NewSidechainOps(reg *SidechainRegistry) *SidechainOps {
	return &SidechainOps{registry: reg}
}

// Deposit credits escrow balance for an address on a side-chain.
func (o *SidechainOps) Deposit(chainID, from string, amount uint64) error {
	o.registry.mu.Lock()
	defer o.registry.mu.Unlock()
	sc, ok := o.registry.chains[chainID]
	if !ok {
		return fmt.Errorf("sidechain %s not found", chainID)
	}
	sc.Deposits[from] += amount
	return nil
}

// Withdraw debits escrow balance if sufficient and a dummy proof is provided.
func (o *SidechainOps) Withdraw(chainID, from string, amount uint64, proof string) error {
	if proof == "" {
		return fmt.Errorf("missing proof")
	}
	o.registry.mu.Lock()
	defer o.registry.mu.Unlock()
	sc, ok := o.registry.chains[chainID]
	if !ok {
		return fmt.Errorf("sidechain %s not found", chainID)
	}
	bal := sc.Deposits[from]
	if bal < amount {
		return fmt.Errorf("insufficient balance")
	}
	sc.Deposits[from] = bal - amount
	return nil
}

// EscrowBalance returns the escrowed balance for an address on a side-chain.
func (o *SidechainOps) EscrowBalance(chainID, addr string) (uint64, error) {
	o.registry.mu.RLock()
	defer o.registry.mu.RUnlock()
	sc, ok := o.registry.chains[chainID]
	if !ok {
		return 0, fmt.Errorf("sidechain %s not found", chainID)
	}
	return sc.Deposits[addr], nil
}

// ListDeposits returns a copy of all deposits for a side-chain.
func (o *SidechainOps) ListDeposits(chainID string) (map[string]uint64, error) {
	o.registry.mu.RLock()
	defer o.registry.mu.RUnlock()
	sc, ok := o.registry.chains[chainID]
	if !ok {
		return nil, fmt.Errorf("sidechain %s not found", chainID)
	}
	out := make(map[string]uint64, len(sc.Deposits))
	for a, v := range sc.Deposits {
		out[a] = v
	}
	return out, nil
}
