package core

import "sync"

// PeerManager maintains a simple peer table for discovery and connection
// management.
type PeerManager struct {
	mu    sync.RWMutex
	peers map[string]string // id -> address
}

// NewPeerManager creates an empty peer manager.
func NewPeerManager() *PeerManager {
	return &PeerManager{peers: make(map[string]string)}
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
