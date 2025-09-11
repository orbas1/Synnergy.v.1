package core

import "sync"

// PeerManager maintains a simple peer table for discovery and connection
// management.
type PeerManager struct {
	mu      sync.RWMutex
	peers   map[string]string // id -> address
	adverts map[string][]string
}

// NewPeerManager creates an empty peer manager.
func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers:   make(map[string]string),
		adverts: make(map[string][]string),
	}
}

// AddPeer records a peer and its address.
func (pm *PeerManager) AddPeer(id, addr string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.peers[id] = addr
}

// RemovePeer deletes a peer from the table.
func (pm *PeerManager) RemovePeer(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.peers, id)
}

// GetPeer returns the address for the peer if known.
func (pm *PeerManager) GetPeer(id string) (string, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	addr, ok := pm.peers[id]
	return addr, ok
}

// ListPeers returns a slice of all peer identifiers.
func (pm *PeerManager) ListPeers() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	ids := make([]string, 0, len(pm.peers))
	for id := range pm.peers {
		ids = append(ids, id)
	}
	return ids
}

// Count returns the number of known peers.
func (pm *PeerManager) Count() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.peers)
}

// Connect records a peer by its network address and returns the derived
// identifier. This is a helper for CLI commands where the address doubles as the
// identifier.
func (pm *PeerManager) Connect(addr string) string {
	pm.AddPeer(addr, addr)
	return addr
}

// Advertise notes that the given peer is advertising on a topic. It does not
// perform any network operations but allows discovery via Discover.
func (pm *PeerManager) Advertise(id, topic string) {
	pm.mu.Lock()
	pm.adverts[topic] = append(pm.adverts[topic], id)
	pm.mu.Unlock()
}

// Discover returns peers that have advertised under the specified topic.
func (pm *PeerManager) Discover(topic string) []string {
	pm.mu.RLock()
	ids := append([]string(nil), pm.adverts[topic]...)
	pm.mu.RUnlock()
	return ids
}
