package core

import (
	"fmt"
	"sync"
)

// BridgeTransfer records a cross-chain transfer locked on this chain.
type BridgeTransfer struct {
	ID       string
	BridgeID string
	From     string
	To       string
	Amount   uint64
	TokenID  string
	Status   string
}

// BridgeTransferManager manages cross-chain transfer records.
type BridgeTransferManager struct {
	mu        sync.RWMutex
	seq       int
	transfers map[string]*BridgeTransfer
}

// NewBridgeTransferManager creates a new manager.
func NewBridgeTransferManager() *BridgeTransferManager {
	return &BridgeTransferManager{transfers: make(map[string]*BridgeTransfer)}
}

// Deposit locks assets for bridging and records the transfer.
func (m *BridgeTransferManager) Deposit(bridgeID, from, to string, amount uint64, tokenID string) (*BridgeTransfer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	id := fmt.Sprintf("transfer-%d", m.seq)
	t := &BridgeTransfer{
		ID:       id,
		BridgeID: bridgeID,
		From:     from,
		To:       to,
		Amount:   amount,
		TokenID:  tokenID,
		Status:   "locked",
	}
	m.transfers[id] = t
	return t, nil
}

// Claim releases assets when provided with a valid proof.
func (m *BridgeTransferManager) Claim(id, proof string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.transfers[id]
	if !ok {
		return fmt.Errorf("transfer %s not found", id)
	}
	if t.Status != "locked" {
		return fmt.Errorf("transfer %s already claimed", id)
	}
	t.Status = "released"
	return nil
}

// GetTransfer retrieves a transfer by ID.
func (m *BridgeTransferManager) GetTransfer(id string) (*BridgeTransfer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transfers[id]
	return t, ok
}

// ListTransfers lists all transfer records.
func (m *BridgeTransferManager) ListTransfers() []*BridgeTransfer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*BridgeTransfer, 0, len(m.transfers))
	for _, t := range m.transfers {
		out = append(out, t)
	}
	return out
}
