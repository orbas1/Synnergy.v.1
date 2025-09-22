package p2p

import (
	"crypto/rand"
	"encoding/hex"
	"sort"
	"sync"
	"time"

	"synnergy/internal/security"
)

// PeerState reflects the availability of a peer.
type PeerState string

const (
	PeerStateUnknown     PeerState = "unknown"
	PeerStateConnected   PeerState = "connected"
	PeerStateFaulted     PeerState = "faulted"
	PeerStateQuarantined PeerState = "quarantined"
)

// Peer describes the metadata that higher-level services expect when
// interacting with a participant on the network.
type Peer struct {
	ID             string
	Address        string
	PubKey         []byte
	NoiseKey       []byte
	TLSFingerprint []byte
	Capabilities   map[string]bool
	Latency        time.Duration
	LastSeen       time.Time
	Labels         []string
	Region         string
	Metadata       map[string]string
	State          PeerState
	FailureCount   int
}

// PeerEventType enumerates the change stream values.
type PeerEventType string

const (
	PeerEventAdded       PeerEventType = "peer_added"
	PeerEventUpdated     PeerEventType = "peer_updated"
	PeerEventRemoved     PeerEventType = "peer_removed"
	PeerEventQuarantined PeerEventType = "peer_quarantined"
)

// PeerEvent captures changes for consumption by the CLI, consensus engine, and
// the JavaScript dashboard.
type PeerEvent struct {
	Type      PeerEventType
	Peer      Peer
	Timestamp time.Time
	Reason    string
}

// Manager maintains a thread-safe registry of peers with health metrics and
// change feed support.
type Manager struct {
	mu       sync.RWMutex
	peers    map[string]*Peer
	watchers map[string]chan PeerEvent
	ddos     *security.DDoSMitigator
}

// NewManager builds a peer manager. The mitigator may be nil in test
// environments.
func NewManager(ddos *security.DDoSMitigator) *Manager {
	return &Manager{
		peers:    make(map[string]*Peer),
		watchers: make(map[string]chan PeerEvent),
		ddos:     ddos,
	}
}

// AddPeer registers a peer and emits an event. Existing peers are merged with
// the new metadata while preserving transient health fields.
func (m *Manager) AddPeer(peer Peer) Peer {
	if peer.ID == "" {
		peer.ID = generatePeerID()
	}
	peer.State = PeerStateConnected
	peer.LastSeen = time.Now().UTC()
	m.mu.Lock()
	existing := m.peers[peer.ID]
	if existing != nil {
		peer.FailureCount = existing.FailureCount
		if peer.Latency == 0 {
			peer.Latency = existing.Latency
		}
		if peer.Capabilities == nil {
			peer.Capabilities = existing.Capabilities
		}
		if peer.Metadata == nil {
			peer.Metadata = existing.Metadata
		}
	}
	m.peers[peer.ID] = &peer
	m.mu.Unlock()
	m.broadcast(PeerEvent{Type: PeerEventAdded, Peer: peer, Timestamp: peer.LastSeen})
	return peer
}

// UpdatePeer overwrites mutable metadata and notifies subscribers.
func (m *Manager) UpdatePeer(peer Peer, reason string) {
	m.mu.Lock()
	if existing, ok := m.peers[peer.ID]; ok {
		peer.FailureCount = existing.FailureCount
		if peer.State == "" {
			peer.State = existing.State
		}
		m.peers[peer.ID] = &peer
	}
	m.mu.Unlock()
	peer.LastSeen = time.Now().UTC()
	m.broadcast(PeerEvent{Type: PeerEventUpdated, Peer: peer, Timestamp: peer.LastSeen, Reason: reason})
}

// RemovePeer removes a peer from the registry and notifies subscribers.
func (m *Manager) RemovePeer(id string, reason string) {
	m.mu.Lock()
	peer, ok := m.peers[id]
	if ok {
		delete(m.peers, id)
	}
	m.mu.Unlock()
	if ok {
		m.broadcast(PeerEvent{Type: PeerEventRemoved, Peer: *peer, Timestamp: time.Now().UTC(), Reason: reason})
	}
}

// MarkFailure increments the failure counter and transitions the state if the
// peer exceeds the retry budget. When a DDoS mitigator is present the peer may
// be quarantined.
func (m *Manager) MarkFailure(id string, reason string) {
	m.mu.Lock()
	peer, ok := m.peers[id]
	if !ok {
		m.mu.Unlock()
		return
	}
	peer.FailureCount++
	if peer.FailureCount > 3 {
		peer.State = PeerStateFaulted
	}
	if m.ddos != nil {
		m.ddos.Block(peer.Address, time.Now().UTC().Add(time.Minute))
		peer.State = PeerStateQuarantined
	}
	updated := *peer
	m.mu.Unlock()
	m.broadcast(PeerEvent{Type: PeerEventQuarantined, Peer: updated, Timestamp: time.Now().UTC(), Reason: reason})
}

// ListPeers returns a sorted slice of peers.
func (m *Manager) ListPeers() []Peer {
	m.mu.RLock()
	out := make([]Peer, 0, len(m.peers))
	for _, peer := range m.peers {
		out = append(out, *peer)
	}
	m.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// GetPeer retrieves a peer by ID.
func (m *Manager) GetPeer(id string) (Peer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	peer, ok := m.peers[id]
	if !ok {
		return Peer{}, false
	}
	return *peer, true
}

// Subscribe registers for peer events. The returned cancel function should be
// invoked to prevent leaks.
func (m *Manager) Subscribe(buffer int) (<-chan PeerEvent, func()) {
	if buffer <= 0 {
		buffer = 16
	}
	ch := make(chan PeerEvent, buffer)
	id := generatePeerID()
	m.mu.Lock()
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

// Snapshot exposes lightweight data used by the CLI and the JS web dashboard.
func (m *Manager) Snapshot() []Peer {
	peers := m.ListPeers()
	for i := range peers {
		peers[i].Metadata = copyMap(peers[i].Metadata)
		peers[i].Capabilities = copyBoolMap(peers[i].Capabilities)
	}
	return peers
}

func (m *Manager) broadcast(evt PeerEvent) {
	m.mu.RLock()
	watchers := make([]chan PeerEvent, 0, len(m.watchers))
	for _, ch := range m.watchers {
		watchers = append(watchers, ch)
	}
	m.mu.RUnlock()
	for _, ch := range watchers {
		select {
		case ch <- evt:
		default:
		}
	}
}

func generatePeerID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(buf)
}

func copyMap(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func copyBoolMap(in map[string]bool) map[string]bool {
	if in == nil {
		return nil
	}
	out := make(map[string]bool, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
