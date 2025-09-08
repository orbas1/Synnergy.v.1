package synnergy

import (
	"errors"
	"sync"
)

// ContentNetworkNode mirrors basic information about a content node in the network.
// It tracks which content items are hosted by this node so others can discover
// available resources.
type ContentNetworkNode struct {
	ID      string
	Address string

	mu       sync.RWMutex
	contents map[string]ContentMeta
}

// NewContentNetworkNode creates a registry entry for a content node.
func NewContentNetworkNode(id, addr string) *ContentNetworkNode {
	return &ContentNetworkNode{
		ID:       id,
		Address:  addr,
		contents: make(map[string]ContentMeta),
	}
}

// Register records the availability of a content item on this node.
// An error is returned if the metadata is missing an ID or already exists.
func (n *ContentNetworkNode) Register(meta ContentMeta) error {
	if meta.ID == "" {
		return errors.New("missing content id")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, exists := n.contents[meta.ID]; exists {
		return errors.New("content already registered")
	}
	n.contents[meta.ID] = meta
	return nil
}

// Unregister removes a content item from this node's registry.
// An error is returned if the item is not present.
func (n *ContentNetworkNode) Unregister(id string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.contents[id]; !ok {
		return errors.New("content not found")
	}
	delete(n.contents, id)
	return nil
}

// Content returns metadata for a hosted content item.
func (n *ContentNetworkNode) Content(id string) (ContentMeta, bool) {
	n.mu.RLock()
	meta, ok := n.contents[id]
	n.mu.RUnlock()
	return meta, ok
}

// List returns all content metadata stored on this node.
func (n *ContentNetworkNode) List() []ContentMeta {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]ContentMeta, 0, len(n.contents))
	for _, m := range n.contents {
		out = append(out, m)
	}
	return out
}
