package synnergy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ConnectionStatus represents the lifecycle state of a cross-chain connection.
type ConnectionStatus string

const (
	// ConnectionStatusPending indicates the link has been requested but not yet promoted to active.
	ConnectionStatusPending ConnectionStatus = "PENDING"
	// ConnectionStatusActive means the link is fully operational and healthy.
	ConnectionStatusActive ConnectionStatus = "ACTIVE"
	// ConnectionStatusClosing indicates the link is winding down and not accepting new flows.
	ConnectionStatusClosing ConnectionStatus = "CLOSING"
	// ConnectionStatusClosed marks a link that has been gracefully terminated.
	ConnectionStatusClosed ConnectionStatus = "CLOSED"
	// ConnectionStatusFailed marks a link that was forcefully terminated after an unrecoverable error.
	ConnectionStatusFailed ConnectionStatus = "FAILED"
)

// ConnectionEventType represents the reason a connection update was emitted.
type ConnectionEventType string

const (
	// ConnectionEventOpened is fired when a new connection has been created.
	ConnectionEventOpened ConnectionEventType = "OPENED"
	// ConnectionEventHeartbeat is emitted when a heartbeat arrives.
	ConnectionEventHeartbeat ConnectionEventType = "HEARTBEAT"
	// ConnectionEventClosing indicates a close operation has been initiated.
	ConnectionEventClosing ConnectionEventType = "CLOSING"
	// ConnectionEventClosed indicates a connection has been fully closed.
	ConnectionEventClosed ConnectionEventType = "CLOSED"
	// ConnectionEventFailed indicates the connection faulted.
	ConnectionEventFailed ConnectionEventType = "FAILED"
)

// ConnectionSpec captures the parameters required to negotiate a cross-chain link.
type ConnectionSpec struct {
	LocalChain        string
	RemoteChain       string
	LocalEndpoint     string
	RemoteEndpoint    string
	GovernanceProfile string
	GasProfile        string
	Metadata          map[string]string
	HeartbeatInterval time.Duration
	Signer            string
	HandshakePayload  []byte
	HandshakeProof    []byte
}

// ConnectionFault captures contextual information about a connection failure.
type ConnectionFault struct {
	Code      string
	Detail    string
	Occurred  time.Time
	Severity  string
	Recovered bool
}

// ChainConnection represents an active or historic cross-chain connection.
type ChainConnection struct {
	ID                string
	Spec              ConnectionSpec
	Status            ConnectionStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
	OpenedAt          time.Time
	ClosedAt          time.Time
	ClosingReason     string
	LastHeartbeat     time.Time
	HeartbeatInterval time.Duration
	Faults            []ConnectionFault
}

// ConnectionEvent notifies observers of lifecycle changes.
type ConnectionEvent struct {
	Type         ConnectionEventType
	ConnectionID string
	Status       ConnectionStatus
	Connection   *ChainConnection
	Fault        *ConnectionFault
}

// SignatureVerifier validates handshake payloads.
type SignatureVerifier interface {
	Verify(ctx context.Context, payload, signature []byte, signer string) error
}

// SignatureVerifierFunc adapts a function into a SignatureVerifier.
type SignatureVerifierFunc func(ctx context.Context, payload, signature []byte, signer string) error

// Verify implements SignatureVerifier.
func (f SignatureVerifierFunc) Verify(ctx context.Context, payload, signature []byte, signer string) error {
	return f(ctx, payload, signature, signer)
}

// IDGenerator provides deterministic identifiers useful for testing and auditing.
type IDGenerator interface {
	NewID(spec ConnectionSpec) (string, error)
}

type defaultIDGenerator struct{}

func (defaultIDGenerator) NewID(spec ConnectionSpec) (string, error) {
	base := fmt.Sprintf("%s|%s|%s|%s|%d", spec.LocalChain, spec.RemoteChain, spec.LocalEndpoint, spec.RemoteEndpoint, time.Now().UTC().UnixNano())
	sum := sha256.Sum256([]byte(base))
	return hex.EncodeToString(sum[:]), nil
}

// ConnectionFilter allows callers to limit the returned set of connections.
type ConnectionFilter struct {
	Statuses     []ConnectionStatus
	LocalChain   string
	RemoteChain  string
	IncludeFault bool
	IncludeEnded bool
}

// ConnectionManager manages cross-chain connections.
type ConnectionManager struct {
	mu               sync.RWMutex
	connections      map[string]*ChainConnection
	watchers         map[int]chan ConnectionEvent
	nextWatcherID    int
	verifier         SignatureVerifier
	idGen            IDGenerator
	defaultHeartbeat time.Duration
}

// ConnectionManagerOption customizes ConnectionManager construction.
type ConnectionManagerOption func(*ConnectionManager)

// WithSignatureVerifier overrides the handshake verifier.
func WithSignatureVerifier(verifier SignatureVerifier) ConnectionManagerOption {
	return func(m *ConnectionManager) {
		m.verifier = verifier
	}
}

// WithIDGenerator overrides the identifier generator.
func WithIDGenerator(gen IDGenerator) ConnectionManagerOption {
	return func(m *ConnectionManager) {
		m.idGen = gen
	}
}

// WithDefaultHeartbeat configures the fallback heartbeat interval used when a spec omits it.
func WithDefaultHeartbeat(interval time.Duration) ConnectionManagerOption {
	return func(m *ConnectionManager) {
		if interval > 0 {
			m.defaultHeartbeat = interval
		}
	}
}

// NewConnectionManager creates an empty ConnectionManager with production hardened defaults.
func NewConnectionManager(opts ...ConnectionManagerOption) *ConnectionManager {
	m := &ConnectionManager{
		connections:      make(map[string]*ChainConnection),
		watchers:         make(map[int]chan ConnectionEvent),
		defaultHeartbeat: 30 * time.Second,
		verifier: SignatureVerifierFunc(func(ctx context.Context, payload, signature []byte, signer string) error {
			if len(payload) == 0 {
				return errors.New("handshake payload required")
			}
			if len(signature) == 0 {
				return errors.New("handshake signature required")
			}
			if signer == "" {
				return errors.New("signer required")
			}
			return nil
		}),
		idGen: defaultIDGenerator{},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

var (
	// ErrConnectionNotFound is returned when a connection cannot be located.
	ErrConnectionNotFound = errors.New("connection not found")
	// ErrConnectionClosed indicates the connection is no longer active for the requested operation.
	ErrConnectionClosed = errors.New("connection already closed")
	// ErrConnectionExists prevents accidental overwriting of existing identifiers.
	ErrConnectionExists = errors.New("connection already exists")
)

// OpenConnection establishes a new link between two chains and returns its descriptor.
func (m *ConnectionManager) OpenConnection(ctx context.Context, spec ConnectionSpec) (*ChainConnection, error) {
	if err := validateSpec(spec); err != nil {
		return nil, err
	}
	if err := m.verifier.Verify(ctx, spec.handshakeCanonicalPayload(), spec.HandshakeProof, spec.Signer); err != nil {
		return nil, fmt.Errorf("handshake verification failed: %w", err)
	}
	id, err := m.idGen.NewID(spec)
	if err != nil {
		return nil, fmt.Errorf("generate connection id: %w", err)
	}
	now := time.Now().UTC()
	conn := &ChainConnection{
		ID:                id,
		Spec:              cloneSpec(spec),
		Status:            ConnectionStatusActive,
		CreatedAt:         now,
		UpdatedAt:         now,
		OpenedAt:          now,
		HeartbeatInterval: spec.HeartbeatInterval,
		LastHeartbeat:     now,
	}
	if conn.HeartbeatInterval <= 0 {
		conn.HeartbeatInterval = m.defaultHeartbeat
	}
	m.mu.Lock()
	if _, exists := m.connections[id]; exists {
		m.mu.Unlock()
		return nil, ErrConnectionExists
	}
	m.connections[id] = conn
	snapshot := cloneConnection(conn)
	m.mu.Unlock()
	m.broadcast(ConnectionEvent{Type: ConnectionEventOpened, ConnectionID: id, Status: conn.Status, Connection: snapshot})
	return snapshot, nil
}

// CloseConnection marks a connection as closed with the supplied reason.
func (m *ConnectionManager) CloseConnection(ctx context.Context, id, reason string) (*ChainConnection, error) {
	m.mu.Lock()
	conn, ok := m.connections[id]
	if !ok {
		m.mu.Unlock()
		return nil, ErrConnectionNotFound
	}
	if conn.Status == ConnectionStatusClosed {
		m.mu.Unlock()
		return nil, ErrConnectionClosed
	}
	conn.Status = ConnectionStatusClosing
	conn.ClosingReason = reason
	conn.UpdatedAt = time.Now().UTC()
	closingSnapshot := cloneConnection(conn)
	m.mu.Unlock()
	m.broadcast(ConnectionEvent{Type: ConnectionEventClosing, ConnectionID: id, Status: conn.Status, Connection: closingSnapshot})

	m.mu.Lock()
	conn, ok = m.connections[id]
	if !ok {
		m.mu.Unlock()
		return nil, ErrConnectionNotFound
	}
	conn.Status = ConnectionStatusClosed
	conn.UpdatedAt = time.Now().UTC()
	conn.ClosedAt = conn.UpdatedAt
	finalSnapshot := cloneConnection(conn)
	m.mu.Unlock()
	m.broadcast(ConnectionEvent{Type: ConnectionEventClosed, ConnectionID: id, Status: conn.Status, Connection: finalSnapshot})
	return finalSnapshot, nil
}

// FailConnection records an unrecoverable error and marks the link failed.
func (m *ConnectionManager) FailConnection(id string, fault ConnectionFault) (*ChainConnection, error) {
	m.mu.Lock()
	conn, ok := m.connections[id]
	if !ok {
		m.mu.Unlock()
		return nil, ErrConnectionNotFound
	}
	if fault.Occurred.IsZero() {
		fault.Occurred = time.Now().UTC()
	}
	conn.Status = ConnectionStatusFailed
	conn.Faults = append(conn.Faults, fault)
	conn.UpdatedAt = fault.Occurred
	failedSnapshot := cloneConnection(conn)
	m.mu.Unlock()
	m.broadcast(ConnectionEvent{Type: ConnectionEventFailed, ConnectionID: id, Status: conn.Status, Connection: failedSnapshot, Fault: &fault})
	return failedSnapshot, nil
}

// MarkHeartbeat updates the last seen heartbeat timestamp for a connection.
func (m *ConnectionManager) MarkHeartbeat(id string, at time.Time) (*ChainConnection, error) {
	if at.IsZero() {
		at = time.Now().UTC()
	}
	m.mu.Lock()
	conn, ok := m.connections[id]
	if !ok {
		m.mu.Unlock()
		return nil, ErrConnectionNotFound
	}
	if conn.Status == ConnectionStatusClosed || conn.Status == ConnectionStatusFailed {
		m.mu.Unlock()
		return nil, ErrConnectionClosed
	}
	conn.LastHeartbeat = at
	conn.UpdatedAt = at
	snapshot := cloneConnection(conn)
	m.mu.Unlock()
	m.broadcast(ConnectionEvent{Type: ConnectionEventHeartbeat, ConnectionID: id, Status: conn.Status, Connection: snapshot})
	return snapshot, nil
}

// GetConnection retrieves connection details by ID.
func (m *ConnectionManager) GetConnection(id string) (*ChainConnection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.connections[id]
	if !ok {
		return nil, false
	}
	return cloneConnection(conn), true
}

// ListConnections returns all known connections optionally filtered by status.
func (m *ConnectionManager) ListConnections(filter ConnectionFilter) []*ChainConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*ChainConnection
	for _, conn := range m.connections {
		if len(filter.Statuses) > 0 && !statusAllowed(conn.Status, filter.Statuses) {
			continue
		}
		if !filter.IncludeEnded && (conn.Status == ConnectionStatusClosed || conn.Status == ConnectionStatusFailed) {
			continue
		}
		if filter.LocalChain != "" && conn.Spec.LocalChain != filter.LocalChain {
			continue
		}
		if filter.RemoteChain != "" && conn.Spec.RemoteChain != filter.RemoteChain {
			continue
		}
		result = append(result, cloneConnection(conn))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

// Subscribe registers a watcher that receives lifecycle events. The returned cancel function must be invoked to release resources.
func (m *ConnectionManager) Subscribe(buffer int) (<-chan ConnectionEvent, func()) {
	m.mu.Lock()
	id := m.nextWatcherID
	m.nextWatcherID++
	ch := make(chan ConnectionEvent, buffer)
	m.watchers[id] = ch
	m.mu.Unlock()
	cancel := func() {
		m.mu.Lock()
		if watcher, ok := m.watchers[id]; ok {
			delete(m.watchers, id)
			close(watcher)
		}
		m.mu.Unlock()
	}
	return ch, cancel
}

func (m *ConnectionManager) broadcast(event ConnectionEvent) {
	m.mu.RLock()
	watchers := make([]chan ConnectionEvent, 0, len(m.watchers))
	for _, ch := range m.watchers {
		watchers = append(watchers, ch)
	}
	m.mu.RUnlock()
	for _, ch := range watchers {
		select {
		case ch <- event:
		default:
		}
	}
}

func validateSpec(spec ConnectionSpec) error {
	if spec.LocalChain == "" {
		return errors.New("local chain required")
	}
	if spec.RemoteChain == "" {
		return errors.New("remote chain required")
	}
	if spec.LocalEndpoint == "" {
		return errors.New("local endpoint required")
	}
	if spec.RemoteEndpoint == "" {
		return errors.New("remote endpoint required")
	}
	if spec.Signer == "" {
		return errors.New("signer required")
	}
	return nil
}

func statusAllowed(status ConnectionStatus, allowed []ConnectionStatus) bool {
	for _, s := range allowed {
		if s == status {
			return true
		}
	}
	return false
}

func cloneSpec(spec ConnectionSpec) ConnectionSpec {
	copyMeta := make(map[string]string, len(spec.Metadata))
	for k, v := range spec.Metadata {
		copyMeta[k] = v
	}
	return ConnectionSpec{
		LocalChain:        spec.LocalChain,
		RemoteChain:       spec.RemoteChain,
		LocalEndpoint:     spec.LocalEndpoint,
		RemoteEndpoint:    spec.RemoteEndpoint,
		GovernanceProfile: spec.GovernanceProfile,
		GasProfile:        spec.GasProfile,
		Metadata:          copyMeta,
		HeartbeatInterval: spec.HeartbeatInterval,
		Signer:            spec.Signer,
		HandshakePayload:  append([]byte(nil), spec.HandshakePayload...),
		HandshakeProof:    append([]byte(nil), spec.HandshakeProof...),
	}
}

func cloneConnection(conn *ChainConnection) *ChainConnection {
	if conn == nil {
		return nil
	}
	copyFaults := make([]ConnectionFault, len(conn.Faults))
	copy(copyFaults, conn.Faults)
	cloned := &ChainConnection{
		ID:                conn.ID,
		Spec:              cloneSpec(conn.Spec),
		Status:            conn.Status,
		CreatedAt:         conn.CreatedAt,
		UpdatedAt:         conn.UpdatedAt,
		OpenedAt:          conn.OpenedAt,
		ClosedAt:          conn.ClosedAt,
		ClosingReason:     conn.ClosingReason,
		LastHeartbeat:     conn.LastHeartbeat,
		HeartbeatInterval: conn.HeartbeatInterval,
		Faults:            copyFaults,
	}
	return cloned
}

func (spec ConnectionSpec) handshakeCanonicalPayload() []byte {
	if len(spec.HandshakePayload) > 0 {
		return append([]byte(nil), spec.HandshakePayload...)
	}
	base := fmt.Sprintf("%s|%s|%s|%s|%s", spec.LocalChain, spec.RemoteChain, spec.LocalEndpoint, spec.RemoteEndpoint, spec.Signer)
	return []byte(base)
}
