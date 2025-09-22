package synnergy

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CrossChainProtocol defines a protocol standard understood across chains.
type CrossChainProtocol struct {
	ID        string
	Name      string
	Version   string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Metadata  map[string]string
}

// ProtocolOption customises protocol registration.
type ProtocolOption func(*CrossChainProtocol)

// WithProtocolVersion assigns a semantic version to the protocol.
func WithProtocolVersion(version string) ProtocolOption {
	return func(p *CrossChainProtocol) {
		p.Version = version
	}
}

// WithProtocolMetadata annotates the protocol with additional metadata.
func WithProtocolMetadata(key, value string) ProtocolOption {
	return func(p *CrossChainProtocol) {
		if p.Metadata == nil {
			p.Metadata = make(map[string]string)
		}
		p.Metadata[key] = value
	}
}

// WithProtocolInactive registers the protocol in an inactive state.
func WithProtocolInactive() ProtocolOption {
	return func(p *CrossChainProtocol) { p.Active = false }
}

// ProtocolEventType enumerates lifecycle events emitted by the registry.
type ProtocolEventType string

const (
	ProtocolEventRegistered  ProtocolEventType = "registered"
	ProtocolEventUpdated     ProtocolEventType = "updated"
	ProtocolEventActivated   ProtocolEventType = "activated"
	ProtocolEventDeactivated ProtocolEventType = "deactivated"
)

// ProtocolEvent captures a change to a protocol definition.
type ProtocolEvent struct {
	Type      ProtocolEventType
	Protocol  CrossChainProtocol
	Timestamp time.Time
}

// ProtocolListener receives protocol events.
type ProtocolListener func(ProtocolEvent)

// ProtocolMetrics exposes aggregate counters for observability.
type ProtocolMetrics struct {
	Total   uint64
	Active  uint64
	Updates uint64
}

type protocolMetrics struct {
	total   uint64
	active  uint64
	updates uint64
}

// ProtocolRegistry stores registered protocol definitions.
type ProtocolRegistry struct {
	mu        sync.RWMutex
	protocols map[string]CrossChainProtocol
	listeners []ProtocolListener
	metrics   protocolMetrics
}

// NewProtocolRegistry creates an empty ProtocolRegistry.
func NewProtocolRegistry() *ProtocolRegistry {
	return &ProtocolRegistry{protocols: make(map[string]CrossChainProtocol)}
}

// RegisterListener attaches a listener for protocol events.
func (r *ProtocolRegistry) RegisterListener(l ProtocolListener) {
	if l == nil {
		return
	}
	r.mu.Lock()
	r.listeners = append(r.listeners, l)
	r.mu.Unlock()
}

func (r *ProtocolRegistry) snapshotListeners() []ProtocolListener {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.listeners) == 0 {
		return nil
	}
	out := make([]ProtocolListener, len(r.listeners))
	copy(out, r.listeners)
	return out
}

func (r *ProtocolRegistry) emit(event ProtocolEvent) {
	for _, l := range r.snapshotListeners() {
		func() {
			defer func() { _ = recover() }()
			l(event)
		}()
	}
}

// RegisterProtocol registers a new protocol and returns its identifier.
func (r *ProtocolRegistry) RegisterProtocol(name string, opts ...ProtocolOption) string {
	if name == "" {
		return ""
	}
	now := time.Now()
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s|%d", name, now.UnixNano()))))
	protocol := CrossChainProtocol{
		ID:        id,
		Name:      name,
		Version:   "v1",
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  make(map[string]string),
	}
	for _, opt := range opts {
		opt(&protocol)
	}

	r.mu.Lock()
	r.protocols[id] = protocol
	r.mu.Unlock()

	atomic.AddUint64(&r.metrics.total, 1)
	if protocol.Active {
		atomic.AddUint64(&r.metrics.active, 1)
	}
	r.emit(ProtocolEvent{Type: ProtocolEventRegistered, Protocol: protocol, Timestamp: now})
	return id
}

// UpdateProtocol applies options to an existing protocol definition.
func (r *ProtocolRegistry) UpdateProtocol(id string, opts ...ProtocolOption) (CrossChainProtocol, bool) {
	now := time.Now()
	r.mu.Lock()
	protocol, ok := r.protocols[id]
	if !ok {
		r.mu.Unlock()
		return CrossChainProtocol{}, false
	}
	for _, opt := range opts {
		opt(&protocol)
	}
	protocol.UpdatedAt = now
	r.protocols[id] = protocol
	r.mu.Unlock()

	atomic.AddUint64(&r.metrics.updates, 1)
	r.emit(ProtocolEvent{Type: ProtocolEventUpdated, Protocol: protocol, Timestamp: now})
	return protocol, true
}

// DeactivateProtocol marks a protocol as inactive.
func (r *ProtocolRegistry) DeactivateProtocol(id string) bool {
	return r.setProtocolActive(id, false)
}

// ActivateProtocol marks a protocol as active.
func (r *ProtocolRegistry) ActivateProtocol(id string) bool {
	return r.setProtocolActive(id, true)
}

func (r *ProtocolRegistry) setProtocolActive(id string, active bool) bool {
	now := time.Now()
	r.mu.Lock()
	protocol, ok := r.protocols[id]
	if !ok {
		r.mu.Unlock()
		return false
	}
	if protocol.Active == active {
		r.mu.Unlock()
		return true
	}
	protocol.Active = active
	protocol.UpdatedAt = now
	r.protocols[id] = protocol
	r.mu.Unlock()

	if active {
		atomic.AddUint64(&r.metrics.active, 1)
		r.emit(ProtocolEvent{Type: ProtocolEventActivated, Protocol: protocol, Timestamp: now})
	} else {
		atomic.AddUint64(&r.metrics.active, ^uint64(0))
		r.emit(ProtocolEvent{Type: ProtocolEventDeactivated, Protocol: protocol, Timestamp: now})
	}
	return true
}

// ListProtocols returns all registered protocols.
func (r *ProtocolRegistry) ListProtocols() []CrossChainProtocol {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]CrossChainProtocol, 0, len(r.protocols))
	for _, p := range r.protocols {
		out = append(out, p)
	}
	return out
}

// GetProtocol retrieves a protocol definition by ID.
func (r *ProtocolRegistry) GetProtocol(id string) (CrossChainProtocol, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.protocols[id]
	return p, ok
}

// Metrics returns a snapshot of registry counters.
func (r *ProtocolRegistry) Metrics() ProtocolMetrics {
	return ProtocolMetrics{
		Total:   atomic.LoadUint64(&r.metrics.total),
		Active:  atomic.LoadUint64(&r.metrics.active),
		Updates: atomic.LoadUint64(&r.metrics.updates),
	}
}
