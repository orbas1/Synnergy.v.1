package core

import (
	"errors"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// AuthorityNode represents a node eligible for governance actions.
type AuthorityNode struct {
	Address string
	Role    string
	Votes   map[string]bool // voter address -> approved
}

// AuthorityNodeRegistry manages authority nodes and voting.
type AuthorityNodeRegistry struct {
	mu    sync.RWMutex
	index *AuthorityNodeIndex
}

// NewAuthorityNodeRegistry creates a new registry.
func NewAuthorityNodeRegistry() *AuthorityNodeRegistry {
	return &AuthorityNodeRegistry{index: NewAuthorityNodeIndex()}
}

// Register adds a candidate as an authority node.
func (r *AuthorityNodeRegistry) Register(addr, role string) (*AuthorityNode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.index.Get(addr); exists {
		return nil, errors.New("authority node already exists")
	}
	node := &AuthorityNode{Address: addr, Role: role, Votes: make(map[string]bool)}
	r.index.Add(node)
	return node, nil
}

// Vote casts a vote for a candidate authority node.
func (r *AuthorityNodeRegistry) Vote(voterAddr, candidateAddr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	node, ok := r.index.Get(candidateAddr)
	if !ok {
		return errors.New("candidate not found")
	}
	node.Votes[voterAddr] = true
	return nil
}

// Electorate samples up to size authority nodes weighted by votes.
func (r *AuthorityNodeRegistry) Electorate(size int) []string {
	r.mu.RLock()
	nodes := r.index.List()
	r.mu.RUnlock()
	sort.Slice(nodes, func(i, j int) bool { return len(nodes[i].Votes) > len(nodes[j].Votes) })
	if size > len(nodes) {
		size = len(nodes)
	}
	selected := make([]string, 0, size)
	for i := 0; i < size; i++ {
		selected = append(selected, nodes[i].Address)
	}
	if size > 1 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(selected), func(i, j int) { selected[i], selected[j] = selected[j], selected[i] })
	}
	return selected
}

// IsAuthorityNode checks if addr is a registered authority node.
func (r *AuthorityNodeRegistry) IsAuthorityNode(addr string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.index.Get(addr)
	return ok
}

// Info returns details for an authority node.
func (r *AuthorityNodeRegistry) Info(addr string) (*AuthorityNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.index.Get(addr)
	if !ok {
		return nil, errors.New("authority node not found")
	}
	return n, nil
}

// List returns all authority nodes.
func (r *AuthorityNodeRegistry) List() []*AuthorityNode {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.index.List()
}

// Deregister removes an authority node and its votes.
func (r *AuthorityNodeRegistry) Deregister(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.index.Remove(addr)
}
