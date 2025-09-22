package core

import (
	"context"
	"sync"
)

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

// Peers returns the identifiers of nodes in the swarm.
func (s *Swarm) Peers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]string, 0, len(s.nodes))
	for id := range s.nodes {
		ids = append(ids, id)
	}
	return ids
}

// StartConsensus triggers block production on all swarm members and returns the
// mined blocks. Nodes without pending transactions simply skip mining.
func (s *Swarm) StartConsensus() []*Block {
	members := s.Members()
	blocks := make([]*Block, 0, len(members))
	for _, n := range members {
		if b, err := n.MineBlock(context.Background()); err == nil && b != nil {
			blocks = append(blocks, b)
		}
	}
	return blocks
}
