package core

import (
	"errors"
	"math/rand"
	"sort"
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
	index *AuthorityNodeIndex
}

// NewAuthorityNodeRegistry creates a new registry.
func NewAuthorityNodeRegistry() *AuthorityNodeRegistry {
	return &AuthorityNodeRegistry{index: NewAuthorityNodeIndex()}
}

// Register adds a candidate as an authority node.
func (r *AuthorityNodeRegistry) Register(addr, role string) (*AuthorityNode, error) {
	if _, exists := r.index.Get(addr); exists {
		return nil, errors.New("authority node already exists")
	}
	node := &AuthorityNode{Address: addr, Role: role, Votes: make(map[string]bool)}
	r.index.Add(node)
	return node, nil
}

// Vote casts a vote for a candidate authority node.
func (r *AuthorityNodeRegistry) Vote(voterAddr, candidateAddr string) error {
	node, ok := r.index.Get(candidateAddr)
	if !ok {
		return errors.New("candidate not found")
	}
	node.Votes[voterAddr] = true
	return nil
}

// Electorate samples up to size authority nodes weighted by votes.
func (r *AuthorityNodeRegistry) Electorate(size int) []string {
	nodes := r.index.List()
	sort.Slice(nodes, func(i, j int) bool { return len(nodes[i].Votes) > len(nodes[j].Votes) })
	if size > len(nodes) {
		size = len(nodes)
	}
	selected := make([]string, 0, size)
	for i := 0; i < size; i++ {
		selected = append(selected, nodes[i].Address)
	}
	if size > 1 {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(selected), func(i, j int) { selected[i], selected[j] = selected[j], selected[i] })
	}
	return selected
}

// IsAuthorityNode checks if addr is a registered authority node.
func (r *AuthorityNodeRegistry) IsAuthorityNode(addr string) bool {
	_, ok := r.index.Get(addr)
	return ok
}

// Info returns details for an authority node.
func (r *AuthorityNodeRegistry) Info(addr string) (*AuthorityNode, error) {
	n, ok := r.index.Get(addr)
	if !ok {
		return nil, errors.New("authority node not found")
	}
	return n, nil
}

// List returns all authority nodes.
func (r *AuthorityNodeRegistry) List() []*AuthorityNode {
	return r.index.List()
}

// Deregister removes an authority node and its votes.
func (r *AuthorityNodeRegistry) Deregister(addr string) {
	r.index.Remove(addr)
}
