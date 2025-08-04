package core

import (
	"errors"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
)

// CrossChainTx records a cross-chain asset movement initiated via LockAndMint or BurnAndRelease.
type CrossChainTx struct {
	ID       string
	BridgeID string
	AssetID  string
	To       string
	Amount   *big.Int
	Type     string
}

// CrossChainTxManager manages cross-chain transactions.
type CrossChainTxManager struct {
	mu     sync.RWMutex
	txs    map[string]*CrossChainTx
	nextID uint64
}

// NewCrossChainTxManager creates a CrossChainTxManager.
func NewCrossChainTxManager() *CrossChainTxManager {
	return &CrossChainTxManager{txs: make(map[string]*CrossChainTx)}
}

// LockAndMint locks native assets and mints wrapped tokens.
func (m *CrossChainTxManager) LockAndMint(bridgeID, assetID string, amount *big.Int, proof string) (*CrossChainTx, error) {
	if amount == nil || amount.Sign() <= 0 {
		return nil, errors.New("amount must be positive")
	}
	id := atomic.AddUint64(&m.nextID, 1)
	tx := &CrossChainTx{
		ID:       formatID("CCT", id),
		BridgeID: bridgeID,
		AssetID:  assetID,
		Amount:   new(big.Int).Set(amount),
		Type:     "lockmint",
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.txs[tx.ID] = tx
	return tx, nil
}

// BurnAndRelease burns wrapped tokens and releases native assets.
func (m *CrossChainTxManager) BurnAndRelease(bridgeID, to, assetID string, amount *big.Int) (*CrossChainTx, error) {
	if amount == nil || amount.Sign() <= 0 {
		return nil, errors.New("amount must be positive")
	}
	id := atomic.AddUint64(&m.nextID, 1)
	tx := &CrossChainTx{
		ID:       formatID("CCT", id),
		BridgeID: bridgeID,
		AssetID:  assetID,
		To:       to,
		Amount:   new(big.Int).Set(amount),
		Type:     "burnrelease",
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.txs[tx.ID] = tx
	return tx, nil
}

// GetTx retrieves a cross-chain transaction by ID.
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
	res := make([]*CrossChainTx, 0, len(m.txs))
	for _, tx := range m.txs {
		res = append(res, tx)
	}
	return res
}

// formatID creates a prefixed incremental identifier.
func formatID(prefix string, id uint64) string {
	return prefix + "-" + strconv.FormatUint(id, 10)
}
