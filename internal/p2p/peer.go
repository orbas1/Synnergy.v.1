package p2p

import "sync"

// Peer represents a network participant.
type Peer struct {
	ID      string
	Address string
	PubKey  []byte
}

// Manager maintains a thread-safe registry of peers.
type Manager struct {
	mu    sync.RWMutex
	peers map[string]Peer
}

// NewManager creates an empty peer manager.
func NewManager() *Manager {
	return &Manager{peers: make(map[string]Peer)}
}

// AddPeer adds or updates a peer in the registry.
func (m *Manager) AddPeer(p Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peers[p.ID] = p
}

// RemovePeer removes a peer by its ID.
func (m *Manager) RemovePeer(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.peers, id)
}

// GetPeer retrieves a peer by ID.
func (m *Manager) GetPeer(id string) (Peer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.peers[id]
	return p, ok
}

// ListPeers returns all known peers.
func (m *Manager) ListPeers() []Peer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Peer, 0, len(m.peers))
	for _, p := range m.peers {
		out = append(out, p)
	}
	return out
}
