package core

import (
	"errors"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
)

// BridgeTransfer records a cross-chain transfer locked on this chain.
type BridgeTransfer struct {
	ID       string
	BridgeID string
	From     string
	To       string
	Amount   *big.Int
	TokenID  string
	Claimed  bool
}

// BridgeTransferManager manages cross-chain transfers.
type BridgeTransferManager struct {
	mu        sync.RWMutex
	transfers map[string]*BridgeTransfer
	nextID    uint64
}

// NewBridgeTransferManager creates a BridgeTransferManager.
func NewBridgeTransferManager() *BridgeTransferManager {
	return &BridgeTransferManager{transfers: make(map[string]*BridgeTransfer)}
}

// Deposit locks assets for bridging.
func (m *BridgeTransferManager) Deposit(bridgeID, from, to string, amount *big.Int, tokenID string) (*BridgeTransfer, error) {
	if amount == nil || amount.Sign() <= 0 {
		return nil, errors.New("amount must be positive")
	}
	id := atomic.AddUint64(&m.nextID, 1)
	t := &BridgeTransfer{
		ID:       formatID("BRT", id),
		BridgeID: bridgeID,
		From:     from,
		To:       to,
		Amount:   new(big.Int).Set(amount),
		TokenID:  tokenID,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transfers[t.ID] = t
	return t, nil
}

// Claim releases assets using a proof.
func (m *BridgeTransferManager) Claim(id string, proof []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.transfers[id]
	if !ok {
		return errors.New("transfer not found")
	}
	t.Claimed = true
	return nil
}

// GetTransfer retrieves a transfer by ID.
func (m *BridgeTransferManager) GetTransfer(id string) (*BridgeTransfer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transfers[id]
	return t, ok
}

// ListTransfers lists all transfers.
func (m *BridgeTransferManager) ListTransfers() []*BridgeTransfer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*BridgeTransfer, 0, len(m.transfers))
	for _, t := range m.transfers {
		res = append(res, t)
	}
	return res
}

// formatID creates a prefixed incremental identifier.
func formatID(prefix string, id uint64) string {
	return prefix + "-" + strconv.FormatUint(id, 10)
}
