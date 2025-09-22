package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Bridge represents a configured link between two chains.
type Bridge struct {
	ID          string
	SourceChain string
	TargetChain string
	Relayers    map[string]struct{}
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Metadata    map[string]string
}

// BridgeOption customises bridge creation.
type BridgeOption func(*Bridge)

// WithBridgeMetadata annotates a bridge with the provided key/value pair.
func WithBridgeMetadata(key, value string) BridgeOption {
	return func(b *Bridge) {
		if b.Metadata == nil {
			b.Metadata = make(map[string]string)
		}
		b.Metadata[key] = value
	}
}

// WithBridgeInactive creates the bridge in an inactive state.
func WithBridgeInactive() BridgeOption {
	return func(b *Bridge) {
		b.Active = false
	}
}

// BridgeEventType enumerates lifecycle notifications emitted by the manager.
type BridgeEventType string

const (
	BridgeEventRegistered        BridgeEventType = "registered"
	BridgeEventActivated         BridgeEventType = "activated"
	BridgeEventDeactivated       BridgeEventType = "deactivated"
	BridgeEventRelayerAuthorized BridgeEventType = "relayer_authorized"
	BridgeEventRelayerRevoked    BridgeEventType = "relayer_revoked"
	BridgeEventMetadataUpdated   BridgeEventType = "metadata_updated"
)

// BridgeEvent describes a single lifecycle transition.
type BridgeEvent struct {
	Type      BridgeEventType
	Bridge    Bridge
	Relayer   string
	Timestamp time.Time
}

// BridgeListener receives bridge lifecycle events.
type BridgeListener func(BridgeEvent)

// BridgeMetrics exposes aggregate counters for observability.
type BridgeMetrics struct {
	Total           uint64
	Active          uint64
	AuthorizedRelay uint64
	RevokedRelay    uint64
}

type bridgeMetrics struct {
	total           uint64
	active          uint64
	authorizedRelay uint64
	revokedRelay    uint64
}

// CrossChainManager manages bridge configurations and authorized relayers.
type CrossChainManager struct {
	mu        sync.RWMutex
	bridges   map[string]*Bridge
	relayers  map[string]struct{}
	listeners []BridgeListener
	metrics   bridgeMetrics
}

// NewCrossChainManager creates an empty CrossChainManager instance.
func NewCrossChainManager() *CrossChainManager {
	return &CrossChainManager{
		bridges:  make(map[string]*Bridge),
		relayers: make(map[string]struct{}),
	}
}

// RegisterListener attaches a listener for bridge events.
func (m *CrossChainManager) RegisterListener(l BridgeListener) {
	if l == nil {
		return
	}
	m.mu.Lock()
	m.listeners = append(m.listeners, l)
	m.mu.Unlock()
}

func (m *CrossChainManager) snapshotListeners() []BridgeListener {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.listeners) == 0 {
		return nil
	}
	out := make([]BridgeListener, len(m.listeners))
	copy(out, m.listeners)
	return out
}

func (m *CrossChainManager) emit(event BridgeEvent) {
	for _, l := range m.snapshotListeners() {
		func() {
			defer func() { _ = recover() }()
			l(event)
		}()
	}
}

// RegisterBridge registers a new bridge configuration and returns its ID.
func (m *CrossChainManager) RegisterBridge(sourceChain, targetChain, relayerAddr string, opts ...BridgeOption) string {
	if sourceChain == "" || targetChain == "" {
		return ""
	}
	now := time.Now()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d", sourceChain, targetChain, now.UnixNano()))))
	bridge := &Bridge{
		ID:          id,
		SourceChain: sourceChain,
		TargetChain: targetChain,
		Relayers:    make(map[string]struct{}),
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    make(map[string]string),
	}
	for _, opt := range opts {
		opt(bridge)
	}
	if relayerAddr != "" {
		bridge.Relayers[relayerAddr] = struct{}{}
	}

	m.mu.Lock()
	m.bridges[id] = bridge
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.total, 1)
	if bridge.Active {
		atomic.AddUint64(&m.metrics.active, 1)
	}
	m.emit(BridgeEvent{Type: BridgeEventRegistered, Bridge: *bridge, Relayer: relayerAddr, Timestamp: now})
	return id
}

// ListBridges returns all registered bridges.
func (m *CrossChainManager) ListBridges() []*Bridge {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Bridge, 0, len(m.bridges))
	for _, b := range m.bridges {
		out = append(out, b)
	}
	return out
}

// GetBridge retrieves a bridge by its identifier.
func (m *CrossChainManager) GetBridge(id string) (*Bridge, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bridges[id]
	return b, ok
}

// ActivateBridge marks the specified bridge as active.
func (m *CrossChainManager) ActivateBridge(id string) error {
	return m.setBridgeActive(id, true)
}

// DeactivateBridge marks the specified bridge as inactive.
func (m *CrossChainManager) DeactivateBridge(id string) error {
	return m.setBridgeActive(id, false)
}

func (m *CrossChainManager) setBridgeActive(id string, active bool) error {
	now := time.Now()
	m.mu.Lock()
	bridge, ok := m.bridges[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("bridge %s not found", id)
	}
	if bridge.Active == active {
		m.mu.Unlock()
		return nil
	}
	bridge.Active = active
	bridge.UpdatedAt = now
	snapshot := *bridge
	m.mu.Unlock()

	if active {
		atomic.AddUint64(&m.metrics.active, 1)
		m.emit(BridgeEvent{Type: BridgeEventActivated, Bridge: snapshot, Timestamp: now})
	} else {
		atomic.AddUint64(&m.metrics.active, ^uint64(0)) // subtract 1 using wrap-around
		m.emit(BridgeEvent{Type: BridgeEventDeactivated, Bridge: snapshot, Timestamp: now})
	}
	return nil
}

// AuthorizeRelayer whitelists a relayer address for all bridges.
func (m *CrossChainManager) AuthorizeRelayer(addr string) {
	if addr == "" {
		return
	}
	m.mu.Lock()
	if _, exists := m.relayers[addr]; exists {
		m.mu.Unlock()
		return
	}
	m.relayers[addr] = struct{}{}
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.authorizedRelay, 1)
	m.emit(BridgeEvent{Type: BridgeEventRelayerAuthorized, Relayer: addr, Timestamp: time.Now()})
}

// AuthorizeBridgeRelayer whitelists a relayer for a specific bridge.
func (m *CrossChainManager) AuthorizeBridgeRelayer(id, addr string) error {
	if addr == "" {
		return fmt.Errorf("relayer address required")
	}
	now := time.Now()
	m.mu.Lock()
	bridge, ok := m.bridges[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("bridge %s not found", id)
	}
	if bridge.Relayers == nil {
		bridge.Relayers = make(map[string]struct{})
	}
	bridge.Relayers[addr] = struct{}{}
	bridge.UpdatedAt = now
	snapshot := *bridge
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.authorizedRelay, 1)
	m.emit(BridgeEvent{Type: BridgeEventRelayerAuthorized, Bridge: snapshot, Relayer: addr, Timestamp: now})
	return nil
}

// RevokeRelayer removes a relayer from the whitelist.
func (m *CrossChainManager) RevokeRelayer(addr string) {
	if addr == "" {
		return
	}
	m.mu.Lock()
	delete(m.relayers, addr)
	for _, bridge := range m.bridges {
		delete(bridge.Relayers, addr)
	}
	m.mu.Unlock()

	atomic.AddUint64(&m.metrics.revokedRelay, 1)
	m.emit(BridgeEvent{Type: BridgeEventRelayerRevoked, Relayer: addr, Timestamp: time.Now()})
}

// UpdateBridgeMetadata sets a metadata entry on the bridge.
func (m *CrossChainManager) UpdateBridgeMetadata(id, key, value string) error {
	now := time.Now()
	m.mu.Lock()
	bridge, ok := m.bridges[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("bridge %s not found", id)
	}
	if bridge.Metadata == nil {
		bridge.Metadata = make(map[string]string)
	}
	bridge.Metadata[key] = value
	bridge.UpdatedAt = now
	snapshot := *bridge
	m.mu.Unlock()

	m.emit(BridgeEvent{Type: BridgeEventMetadataUpdated, Bridge: snapshot, Timestamp: now})
	return nil
}

// IsRelayerAuthorized returns true if the relayer is currently whitelisted.
func (m *CrossChainManager) IsRelayerAuthorized(addr string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.relayers[addr]
	return ok
}

// Metrics returns a snapshot of bridge counters.
func (m *CrossChainManager) Metrics() BridgeMetrics {
	return BridgeMetrics{
		Total:           atomic.LoadUint64(&m.metrics.total),
		Active:          atomic.LoadUint64(&m.metrics.active),
		AuthorizedRelay: atomic.LoadUint64(&m.metrics.authorizedRelay),
		RevokedRelay:    atomic.LoadUint64(&m.metrics.revokedRelay),
	}
}
