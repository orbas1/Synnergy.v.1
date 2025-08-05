package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// BridgeTransfer represents a token movement through a bridge.
type BridgeTransfer struct {
	ID        string
	BridgeID  string
	From      string
	To        string
	Amount    uint64
	TokenID   string
	Proof     []byte
	Claimed   bool
	CreatedAt time.Time
	ClaimedAt time.Time
}

// BridgeTransferManager tracks cross-chain deposits and claims.
type BridgeTransferManager struct {
	mu        sync.RWMutex
	transfers map[string]*BridgeTransfer
}

// NewBridgeTransferManager creates an empty BridgeTransferManager.
func NewBridgeTransferManager() *BridgeTransferManager {
	return &BridgeTransferManager{transfers: make(map[string]*BridgeTransfer)}
}

// Deposit locks assets for bridging and returns a transfer ID.
func (m *BridgeTransferManager) Deposit(bridgeID, from, to string, amount uint64, tokenID string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%d|%d", bridgeID, from, to, amount, time.Now().UnixNano()))))
	m.transfers[id] = &BridgeTransfer{
		ID:        id,
		BridgeID:  bridgeID,
		From:      from,
		To:        to,
		Amount:    amount,
		TokenID:   tokenID,
		CreatedAt: time.Now(),
	}
	return id
}

// Claim releases bridged assets using a provided proof.
func (m *BridgeTransferManager) Claim(id string, proof []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.transfers[id]
	if !ok {
		return fmt.Errorf("transfer not found")
	}
	if t.Claimed {
		return fmt.Errorf("transfer already claimed")
	}
	t.Proof = proof
	t.Claimed = true
	t.ClaimedAt = time.Now()
	return nil
}

// GetTransfer retrieves a transfer record by ID.
func (m *BridgeTransferManager) GetTransfer(id string) (*BridgeTransfer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transfers[id]
	return t, ok
}

// ListTransfers returns all recorded transfers.
func (m *BridgeTransferManager) ListTransfers() []*BridgeTransfer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*BridgeTransfer, 0, len(m.transfers))
	for _, t := range m.transfers {
		out = append(out, t)
	}
	return out
}
