package core

import "sync"

// WebRTCRPC simulates an RPC mechanism over WebRTC style peer connections. It
// maintains in-memory channels for peers to exchange messages.
type WebRTCRPC struct {
	mu    sync.RWMutex
	peers map[string]chan []byte
}

// NewWebRTCRPC creates a new WebRTCRPC instance.
func NewWebRTCRPC() *WebRTCRPC {
	return &WebRTCRPC{peers: make(map[string]chan []byte)}
}

// Connect registers a peer and returns a receive-only channel for incoming
// messages.
func (r *WebRTCRPC) Connect(id string) <-chan []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan []byte, 1)
	r.peers[id] = ch
	return ch
}

// Send delivers a message to the specified peer if it exists.
func (r *WebRTCRPC) Send(id string, msg []byte) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, ok := r.peers[id]
	if !ok {
		return false
	}
	select {
	case ch <- append([]byte(nil), msg...):
		return true
	default:
		return false
	}
}

// Broadcast sends a message to all connected peers. Returns count of peers
// successfully receiving the message.
func (r *WebRTCRPC) Broadcast(msg []byte) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, ch := range r.peers {
		select {
		case ch <- append([]byte(nil), msg...):
			count++
		default:
		}
	}
	return count
}

// Disconnect removes a peer from the RPC network.
func (r *WebRTCRPC) Disconnect(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.peers[id]; ok {
		close(ch)
		delete(r.peers, id)
	}
}

// Peers returns the list of connected peer IDs.
func (r *WebRTCRPC) Peers() []string {
	r.mu.RLock()
	ids := make([]string, 0, len(r.peers))
	for id := range r.peers {
		ids = append(ids, id)
	}
	r.mu.RUnlock()
	return ids
}
