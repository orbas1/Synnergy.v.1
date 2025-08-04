package synnergy

import "sync"

// DataSet represents a piece of content offered through the network and the
// nodes that host it.
type DataSet struct {
	Meta  ContentMeta
	Nodes map[string]struct{}
}

// DataDistribution keeps track of which nodes serve specific datasets.
type DataDistribution struct {
	mu   sync.RWMutex
	sets map[string]*DataSet
}

// NewDataDistribution creates an empty DataDistribution registry.
func NewDataDistribution() *DataDistribution {
	return &DataDistribution{sets: make(map[string]*DataSet)}
}

// Offer registers that a node hosts the given content dataset.
func (d *DataDistribution) Offer(nodeID string, meta ContentMeta) {
	d.mu.Lock()
	ds, ok := d.sets[meta.ID]
	if !ok {
		ds = &DataSet{Meta: meta, Nodes: make(map[string]struct{})}
		d.sets[meta.ID] = ds
	}
	ds.Nodes[nodeID] = struct{}{}
	d.mu.Unlock()
}

// Locations returns the ids of nodes currently hosting the specified dataset.
func (d *DataDistribution) Locations(contentID string) []string {
	d.mu.RLock()
	ds, ok := d.sets[contentID]
	d.mu.RUnlock()
	if !ok {
		return nil
	}
	nodes := make([]string, 0, len(ds.Nodes))
	for n := range ds.Nodes {
		nodes = append(nodes, n)
	}
	return nodes
}
