package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// CrossChainTransaction represents a cross-chain asset movement.
type CrossChainTransaction struct {
	ID       string
	BridgeID string
	Type     string // "lockmint" or "burnrelease"
	AssetID  string
	Amount   uint64
	To       string
	Proof    string
	Created  time.Time
}

// TransactionManager manages cross-chain transfer records.
type TransactionManager struct {
	mu  sync.RWMutex
	txs map[string]CrossChainTransaction
}

// NewTransactionManager creates an empty TransactionManager.
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{txs: make(map[string]CrossChainTransaction)}
}

// LockAndMint records a lock and mint transfer.
func (m *TransactionManager) LockAndMint(bridgeID, assetID string, amount uint64, proof string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d", bridgeID, assetID, time.Now().UnixNano()))))
	m.txs[id] = CrossChainTransaction{
		ID:       id,
		BridgeID: bridgeID,
		Type:     "lockmint",
		AssetID:  assetID,
		Amount:   amount,
		Proof:    proof,
		Created:  time.Now(),
	}
	return id
}

// BurnAndRelease records a burn and release transfer.
func (m *TransactionManager) BurnAndRelease(bridgeID, to, assetID string, amount uint64) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%d|%d", bridgeID, to, assetID, amount, time.Now().UnixNano()))))
	m.txs[id] = CrossChainTransaction{
		ID:       id,
		BridgeID: bridgeID,
		Type:     "burnrelease",
		AssetID:  assetID,
		Amount:   amount,
		To:       to,
		Created:  time.Now(),
	}
	return id
}

// ListTransactions returns all recorded cross-chain transactions.
func (m *TransactionManager) ListTransactions() []CrossChainTransaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]CrossChainTransaction, 0, len(m.txs))
	for _, tx := range m.txs {
		out = append(out, tx)
	}
	return out
}

// GetTransaction retrieves a transaction by ID.
func (m *TransactionManager) GetTransaction(id string) (CrossChainTransaction, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tx, ok := m.txs[id]
	return tx, ok
}
