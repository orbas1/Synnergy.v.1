package core

import (
	"fmt"
	"sync"
)

// CrossChainTx records a cross-chain asset movement initiated via LockAndMint or BurnAndRelease.
type CrossChainTx struct {
	ID       string
	BridgeID string
	AssetID  string
	Amount   uint64
	To       string
	Type     string // lockmint or burnrelease
	Proof    string
}

// CrossChainTxManager manages cross-chain transaction records.
type CrossChainTxManager struct {
	mu  sync.RWMutex
	seq int
	txs map[string]*CrossChainTx
}

// NewCrossChainTxManager creates a new manager.
func NewCrossChainTxManager() *CrossChainTxManager {
	return &CrossChainTxManager{txs: make(map[string]*CrossChainTx)}
}

// LockAndMint locks native assets and mints wrapped tokens.
func (m *CrossChainTxManager) LockAndMint(bridgeID, assetID string, amount uint64, proof string) (*CrossChainTx, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	id := fmt.Sprintf("cctx-%d", m.seq)
	tx := &CrossChainTx{ID: id, BridgeID: bridgeID, AssetID: assetID, Amount: amount, Type: "lockmint", Proof: proof}
	m.txs[id] = tx
	return tx, nil
}

// BurnAndRelease burns wrapped tokens and releases native assets.
func (m *CrossChainTxManager) BurnAndRelease(bridgeID, to, assetID string, amount uint64) (*CrossChainTx, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	id := fmt.Sprintf("cctx-%d", m.seq)
	tx := &CrossChainTx{ID: id, BridgeID: bridgeID, AssetID: assetID, Amount: amount, To: to, Type: "burnrelease"}
	m.txs[id] = tx
	return tx, nil
}

// GetTx retrieves a transaction by ID.
func (m *CrossChainTxManager) GetTx(id string) (*CrossChainTx, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tx, ok := m.txs[id]
	return tx, ok
}

// ListTxs lists all cross-chain transactions.
func (m *CrossChainTxManager) ListTxs() []*CrossChainTx {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*CrossChainTx, 0, len(m.txs))
	for _, tx := range m.txs {
		out = append(out, tx)
	}
	return out
}
