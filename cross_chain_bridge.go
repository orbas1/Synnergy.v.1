package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// BridgeTransferStatus enumerates the lifecycle of a bridge transfer.
type BridgeTransferStatus string

const (
	TransferPending BridgeTransferStatus = "pending"
	TransferClaimed BridgeTransferStatus = "claimed"
	TransferFailed  BridgeTransferStatus = "failed"
	TransferExpired BridgeTransferStatus = "expired"
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
	Status    BridgeTransferStatus
	Metadata  map[string]string
	ClaimedAt time.Time
	CreatedAt time.Time
	ExpiresAt time.Time
}

// BridgeTransferOption customises transfer creation.
type BridgeTransferOption func(*BridgeTransfer)

// WithTransferExpiry sets a deadline for claiming the transfer.
func WithTransferExpiry(expiry time.Time) BridgeTransferOption {
	return func(t *BridgeTransfer) { t.ExpiresAt = expiry }
}

// WithTransferMetadata annotates the transfer with metadata.
func WithTransferMetadata(key, value string) BridgeTransferOption {
	return func(t *BridgeTransfer) {
		if t.Metadata == nil {
			t.Metadata = make(map[string]string)
		}
		t.Metadata[key] = value
	}
}

// BridgeTransferEventType enumerates events emitted by the transfer manager.
type BridgeTransferEventType string

const (
	TransferEventDeposited BridgeTransferEventType = "deposited"
	TransferEventClaimed   BridgeTransferEventType = "claimed"
	TransferEventFailed    BridgeTransferEventType = "failed"
	TransferEventExpired   BridgeTransferEventType = "expired"
)

// BridgeTransferEvent captures a transfer lifecycle transition.
type BridgeTransferEvent struct {
	Type      BridgeTransferEventType
	Transfer  BridgeTransfer
	Timestamp time.Time
	Err       error
}

// BridgeTransferListener receives transfer events.
type BridgeTransferListener func(BridgeTransferEvent)

// BridgeTransferMetrics aggregates counters for observability.
type BridgeTransferMetrics struct {
	Total   uint64
	Claimed uint64
	Failed  uint64
	Expired uint64
}

type bridgeTransferMetrics struct {
	total   uint64
	claimed uint64
	failed  uint64
	expired uint64
}

// BridgeTransferManager tracks cross-chain deposits and claims.
type BridgeTransferManager struct {
	mu        sync.RWMutex
	transfers map[string]*BridgeTransfer
	listeners []BridgeTransferListener
	metrics   bridgeTransferMetrics
}

// NewBridgeTransferManager creates an empty BridgeTransferManager.
func NewBridgeTransferManager() *BridgeTransferManager {
	return &BridgeTransferManager{transfers: make(map[string]*BridgeTransfer)}
}

// RegisterListener attaches a listener for transfer events.
func (m *BridgeTransferManager) RegisterListener(l BridgeTransferListener) {
	if l == nil {
		return
	}
	m.mu.Lock()
	m.listeners = append(m.listeners, l)
	m.mu.Unlock()
}

func (m *BridgeTransferManager) snapshotListeners() []BridgeTransferListener {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.listeners) == 0 {
		return nil
	}
	out := make([]BridgeTransferListener, len(m.listeners))
	copy(out, m.listeners)
	return out
}

func (m *BridgeTransferManager) emit(event BridgeTransferEvent) {
	for _, l := range m.snapshotListeners() {
		func() {
			defer func() { _ = recover() }()
			l(event)
		}()
	}
}

// Deposit locks assets for bridging and returns a transfer ID.
func (m *BridgeTransferManager) Deposit(bridgeID, from, to string, amount uint64, tokenID string, opts ...BridgeTransferOption) string {
	now := time.Now()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%d|%d", bridgeID, from, to, amount, now.UnixNano()))))
	transfer := &BridgeTransfer{
		ID:        id,
		BridgeID:  bridgeID,
		From:      from,
		To:        to,
		Amount:    amount,
		TokenID:   tokenID,
		Status:    TransferPending,
		Metadata:  make(map[string]string),
		CreatedAt: now,
	}
	for _, opt := range opts {
		opt(transfer)
	}

	m.mu.Lock()
	m.transfers[id] = transfer
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.total, 1)
	m.emit(BridgeTransferEvent{Type: TransferEventDeposited, Transfer: *transfer, Timestamp: now})
	return id
}

// Claim releases bridged assets using a provided proof.
func (m *BridgeTransferManager) Claim(id string, proof []byte) error {
	now := time.Now()
	m.mu.Lock()
	t, ok := m.transfers[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("transfer not found")
	}
	if t.Status != TransferPending {
		m.mu.Unlock()
		return fmt.Errorf("transfer not claimable")
	}
	t.Proof = proof
	t.Status = TransferClaimed
	t.ClaimedAt = now
	snapshot := *t
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.claimed, 1)
	m.emit(BridgeTransferEvent{Type: TransferEventClaimed, Transfer: snapshot, Timestamp: now})
	return nil
}

// Fail marks a transfer as failed.
func (m *BridgeTransferManager) Fail(id string, err error) error {
	now := time.Now()
	m.mu.Lock()
	t, ok := m.transfers[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("transfer not found")
	}
	t.Status = TransferFailed
	snapshot := *t
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.failed, 1)
	m.emit(BridgeTransferEvent{Type: TransferEventFailed, Transfer: snapshot, Timestamp: now, Err: err})
	return nil
}

// SweepExpired marks and returns transfers that expired before now.
func (m *BridgeTransferManager) SweepExpired(now time.Time) []BridgeTransfer {
	var expired []BridgeTransfer
	m.mu.Lock()
	for id, t := range m.transfers {
		if t.Status != TransferPending {
			continue
		}
		if t.ExpiresAt.IsZero() || now.Before(t.ExpiresAt) {
			continue
		}
		t.Status = TransferExpired
		expired = append(expired, *t)
		delete(m.transfers, id)
	}
	m.mu.Unlock()

	if len(expired) > 0 {
		atomic.AddUint64(&m.metrics.expired, uint64(len(expired)))
		for _, tr := range expired {
			m.emit(BridgeTransferEvent{Type: TransferEventExpired, Transfer: tr, Timestamp: now})
		}
	}
	return expired
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

// Metrics returns a snapshot of transfer counters.
func (m *BridgeTransferManager) Metrics() BridgeTransferMetrics {
	return BridgeTransferMetrics{
		Total:   atomic.LoadUint64(&m.metrics.total),
		Claimed: atomic.LoadUint64(&m.metrics.claimed),
		Failed:  atomic.LoadUint64(&m.metrics.failed),
		Expired: atomic.LoadUint64(&m.metrics.expired),
	}
}
