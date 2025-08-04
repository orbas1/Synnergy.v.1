package core

import (
	"errors"
	"sync"
)

// CrossChainTxType identifies the type of cross-chain transaction.
type CrossChainTxType string

const (
	// TxTypeLockMint locks native assets and mints wrapped tokens.
	TxTypeLockMint CrossChainTxType = "lockmint"
	// TxTypeBurnRelease burns wrapped tokens and releases native assets.
	TxTypeBurnRelease CrossChainTxType = "burnrelease"
)

// CrossChainTransfer records a cross-chain asset transfer.
type CrossChainTransfer struct {
	ID        int
	BridgeID  int
	From      string
	To        string
	AssetID   string
	Amount    uint64
	Type      CrossChainTxType
	Completed bool
}

// CrossChainTxManager manages cross-chain transfers.
type CrossChainTxManager struct {
	mu     sync.RWMutex
	txs    map[int]*CrossChainTransfer
	nextID int
	ledger *Ledger
}

// NewCrossChainTxManager creates a new manager bound to a ledger.
func NewCrossChainTxManager(l *Ledger) *CrossChainTxManager {
	return &CrossChainTxManager{txs: make(map[int]*CrossChainTransfer), ledger: l}
}

// LockMint locks native assets from the sender and credits wrapped tokens to the recipient.
func (m *CrossChainTxManager) LockMint(bridgeID int, from, to, assetID string, amount uint64, proof string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ledger == nil {
		return 0, errors.New("ledger not configured")
	}
	tx := &Transaction{From: from, To: "lockmint_escrow", Amount: amount}
	if err := m.ledger.ApplyTransaction(tx); err != nil {
		return 0, err
	}
	m.ledger.Credit(to, amount)
	m.nextID++
	id := m.nextID
	m.txs[id] = &CrossChainTransfer{ID: id, BridgeID: bridgeID, From: from, To: to, AssetID: assetID, Amount: amount, Type: TxTypeLockMint}
	return id, nil
}

// BurnRelease burns wrapped tokens from the sender and releases native assets to the recipient.
func (m *CrossChainTxManager) BurnRelease(bridgeID int, from, to, assetID string, amount uint64) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ledger == nil {
		return 0, errors.New("ledger not configured")
	}
	tx := &Transaction{From: from, To: "burn_vault", Amount: amount}
	if err := m.ledger.ApplyTransaction(tx); err != nil {
		return 0, err
	}
	m.ledger.Credit(to, amount)
	m.nextID++
	id := m.nextID
	m.txs[id] = &CrossChainTransfer{ID: id, BridgeID: bridgeID, From: from, To: to, AssetID: assetID, Amount: amount, Type: TxTypeBurnRelease, Completed: true}
	return id, nil
}

// ListTransfers returns all cross-chain transfer records.
func (m *CrossChainTxManager) ListTransfers() []*CrossChainTransfer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*CrossChainTransfer, 0, len(m.txs))
	for _, t := range m.txs {
		out = append(out, t)
	}
	return out
}

// GetTransfer retrieves a transfer by ID.
func (m *CrossChainTxManager) GetTransfer(id int) (*CrossChainTransfer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.txs[id]
	if !ok {
		return nil, errors.New("transfer not found")
	}
	return t, nil

}
