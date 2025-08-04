package core

import (
	"errors"
	"fmt"
	"sync"
)

// AssetBridge defines a token bridge between two chains with a relayer whitelist.
type AssetBridge struct {
	ID       int
	Source   string
	Target   string
	Relayers map[string]struct{}
}

// BridgeTransferRecord records a cross-chain transfer locked on this chain.
type BridgeTransferRecord struct {
	ID       int
	BridgeID int
	From     string
	To       string
	Amount   uint64
	TokenID  string
	Claimed  bool
}

// BridgeManager manages bridges and transfer records.
type BridgeManager struct {
	mu             sync.RWMutex
	bridges        map[int]*AssetBridge
	transfers      map[int]*BridgeTransferRecord
	nextBridgeID   int
	nextTransferID int
	ledger         *Ledger
}

// NewBridgeManager creates a new manager using the provided ledger for balance operations.
func NewBridgeManager(l *Ledger) *BridgeManager {
	return &BridgeManager{
		bridges:   make(map[int]*AssetBridge),
		transfers: make(map[int]*BridgeTransferRecord),
		ledger:    l,
	}
}

// RegisterBridge creates a new bridge definition and returns its ID.
func (m *BridgeManager) RegisterBridge(source, target, relayer string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextBridgeID++
	id := m.nextBridgeID
	relayers := make(map[string]struct{})
	if relayer != "" {
		relayers[relayer] = struct{}{}
	}
	m.bridges[id] = &AssetBridge{ID: id, Source: source, Target: target, Relayers: relayers}
	return id
}

// ListBridges returns all registered bridges.
func (m *BridgeManager) ListBridges() []*AssetBridge {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*AssetBridge, 0, len(m.bridges))
	for _, b := range m.bridges {
		out = append(out, b)
	}
	return out
}

// GetBridge retrieves a bridge by ID.
func (m *BridgeManager) GetBridge(id int) (*AssetBridge, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bridges[id]
	if !ok {
		return nil, errors.New("bridge not found")
	}
	return b, nil
}

// AuthorizeRelayer adds an address to the bridge's relayer whitelist.
func (m *BridgeManager) AuthorizeRelayer(id int, addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bridges[id]
	if !ok {
		return errors.New("bridge not found")
	}
	b.Relayers[addr] = struct{}{}
	return nil
}

// RevokeRelayer removes an address from the bridge's whitelist.
func (m *BridgeManager) RevokeRelayer(id int, addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bridges[id]
	if !ok {
		return errors.New("bridge not found")
	}
	delete(b.Relayers, addr)
	return nil
}

// Deposit locks assets on the source chain creating a transfer record.
func (m *BridgeManager) Deposit(bridgeID int, from, to string, amount uint64, tokenID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bridges[bridgeID]; !ok {
		return 0, errors.New("bridge not found")
	}
	if m.ledger == nil {
		return 0, errors.New("ledger not configured")
	}
	tx := &Transaction{From: from, To: "bridge_escrow", Amount: amount}
	if err := m.ledger.ApplyTransaction(tx); err != nil {
		return 0, err
	}
	m.nextTransferID++
	id := m.nextTransferID
	m.transfers[id] = &BridgeTransferRecord{ID: id, BridgeID: bridgeID, From: from, To: to, Amount: amount, TokenID: tokenID}
	return id, nil
}

// Claim releases locked assets to the recipient using a proof placeholder.
func (m *BridgeManager) Claim(transferID int, proof string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.transfers[transferID]
	if !ok {
		return errors.New("transfer not found")
	}
	if t.Claimed {
		return errors.New("transfer already claimed")
	}
	if m.ledger == nil {
		return errors.New("ledger not configured")
	}
	m.ledger.Credit(t.To, t.Amount)
	t.Claimed = true
	return nil
}

// GetTransfer returns a transfer record by ID.
func (m *BridgeManager) GetTransfer(id int) (*BridgeTransferRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transfers[id]
	if !ok {
		return nil, errors.New("transfer not found")
	}
	return t, nil
}

// ListTransfers returns all transfer records.
func (m *BridgeManager) ListTransfers() []*BridgeTransferRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*BridgeTransferRecord, 0, len(m.transfers))
	for _, t := range m.transfers {
		out = append(out, t)
	}
	return out
}

// BridgeTransfer represents a transfer record used by BridgeTransferManager with string IDs.
type BridgeTransfer struct {
	ID       string
	BridgeID string
	From     string
	To       string
	Amount   uint64
	TokenID  string
	Status   string
}

// BridgeTransferManager manages cross-chain transfer records with string identifiers.
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
