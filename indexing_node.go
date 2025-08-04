package synnergy

import "sync"

// IndexingNode provides fast query capabilities by indexing ledger data into
// an in-memory key/value store.
type IndexingNode struct {
	mu    sync.RWMutex
	index map[string][]byte
}

// NewIndexingNode constructs an empty IndexingNode.
func NewIndexingNode() *IndexingNode {
	return &IndexingNode{index: make(map[string][]byte)}
}

// Index inserts or updates the value for a given key.
func (n *IndexingNode) Index(key string, value []byte) {
	n.mu.Lock()
	dup := append([]byte(nil), value...)
	n.index[key] = dup
	n.mu.Unlock()
}

// Query retrieves a copy of the value associated with the key.
func (n *IndexingNode) Query(key string) ([]byte, bool) {
	n.mu.RLock()
	val, ok := n.index[key]
	n.mu.RUnlock()
	if !ok {
		return nil, false
	}
	out := make([]byte, len(val))
	copy(out, val)
	return out, true
}

// Remove deletes the key from the index if present.
func (n *IndexingNode) Remove(key string) {
	n.mu.Lock()
	delete(n.index, key)
	n.mu.Unlock()
}

// Keys returns a snapshot of all indexed keys.
func (n *IndexingNode) Keys() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	keys := make([]string, 0, len(n.index))
	for k := range n.index {
		keys = append(keys, k)
	}
	return keys
}
