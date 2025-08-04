package core

import "sync"

// Swarm groups nodes to enable coordinated operations and simplified message
// broadcasting among related participants.
type Swarm struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

// NewSwarm creates an empty Swarm instance.
func NewSwarm() *Swarm {
	return &Swarm{nodes: make(map[string]*Node)}
}

// Join adds a node to the swarm.
func (s *Swarm) Join(n *Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes[n.ID] = n
}

// Leave removes a node from the swarm by its ID.
func (s *Swarm) Leave(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.nodes, id)
}

// Members returns a snapshot slice of nodes currently in the swarm.
func (s *Swarm) Members() []*Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Node, 0, len(s.nodes))
	for _, n := range s.nodes {
		out = append(out, n)
	}
	return out
}

// Broadcast sends a transaction to all swarm members by invoking AddTransaction
// on each node. Errors are ignored to keep broadcast best-effort.
func (s *Swarm) Broadcast(tx *Transaction) {
	for _, n := range s.Members() {
		_ = n.AddTransaction(tx)
	}
}
