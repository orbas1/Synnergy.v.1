package nodes

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// NodeInterface defines the lifecycle behaviour for specialised nodes.  The
// interface intentionally mirrors the minimal expectations of CLI tooling and
// monitoring code so it can be shared without importing the runtime packages.
type NodeInterface interface {
	ID() string
	Start() error
	Stop() error
}

// NodeMetrics tracks the lifecycle state for a registered node.
type NodeMetrics struct {
	Running     bool
	LastStarted time.Time
	LastStopped time.Time
	LastError   string
}

// Registry manages specialised nodes for tests and CLI tooling, providing
// concurrency-safe lifecycle orchestration.
type Registry struct {
	mu      sync.RWMutex
	nodes   map[string]NodeInterface
	metrics map[string]NodeMetrics
}

// NewRegistry constructs an empty registry.
func NewRegistry() *Registry {
	return &Registry{
		nodes:   make(map[string]NodeInterface),
		metrics: make(map[string]NodeMetrics),
	}
}

// Register adds a node to the registry. It returns an error when a duplicate ID
// is provided.
func (r *Registry) Register(node NodeInterface) error {
	if node == nil {
		return errors.New("node cannot be nil")
	}
	id := node.ID()
	if id == "" {
		return errors.New("node id cannot be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.nodes[id]; exists {
		return fmt.Errorf("node with id %s already registered", id)
	}
	r.nodes[id] = node
	r.metrics[id] = NodeMetrics{}
	return nil
}

// Unregister removes a node from the registry. It is safe to call multiple
// times.
func (r *Registry) Unregister(id string) {
	r.mu.Lock()
	delete(r.nodes, id)
	delete(r.metrics, id)
	r.mu.Unlock()
}

// Start activates a node by identifier and records lifecycle metrics.
func (r *Registry) Start(id string) error {
	node, metrics, err := r.get(id)
	if err != nil {
		return err
	}
	if err := node.Start(); err != nil {
		r.setMetrics(id, metrics.withError(err))
		return err
	}
	metrics.Running = true
	metrics.LastStarted = time.Now().UTC()
	metrics.LastError = ""
	r.setMetrics(id, metrics)
	return nil
}

// Stop halts a node by identifier and records lifecycle metrics.
func (r *Registry) Stop(id string) error {
	node, metrics, err := r.get(id)
	if err != nil {
		return err
	}
	if err := node.Stop(); err != nil {
		r.setMetrics(id, metrics.withError(err))
		return err
	}
	metrics.Running = false
	metrics.LastStopped = time.Now().UTC()
	r.setMetrics(id, metrics)
	return nil
}

// StartAll activates all registered nodes, returning the first error
// encountered.
func (r *Registry) StartAll() error {
	for _, id := range r.sortedIDs() {
		if err := r.Start(id); err != nil {
			return err
		}
	}
	return nil
}

// StopAll halts all registered nodes, returning the first error encountered.
func (r *Registry) StopAll() error {
	ids := r.sortedIDs()
	for i := len(ids) - 1; i >= 0; i-- {
		if err := r.Stop(ids[i]); err != nil {
			return err
		}
	}
	return nil
}

// Metrics returns lifecycle metrics for a node.
func (r *Registry) Metrics(id string) (NodeMetrics, bool) {
	r.mu.RLock()
	metrics, ok := r.metrics[id]
	r.mu.RUnlock()
	return metrics, ok
}

func (r *Registry) get(id string) (NodeInterface, NodeMetrics, error) {
	r.mu.RLock()
	node, ok := r.nodes[id]
	metrics := r.metrics[id]
	r.mu.RUnlock()
	if !ok {
		return nil, NodeMetrics{}, fmt.Errorf("node %s not registered", id)
	}
	return node, metrics, nil
}

func (r *Registry) setMetrics(id string, metrics NodeMetrics) {
	r.mu.Lock()
	r.metrics[id] = metrics
	r.mu.Unlock()
}

func (r *Registry) sortedIDs() []string {
	r.mu.RLock()
	ids := make([]string, 0, len(r.nodes))
	for id := range r.nodes {
		ids = append(ids, id)
	}
	r.mu.RUnlock()
	sort.Strings(ids)
	return ids
}

func (m NodeMetrics) withError(err error) NodeMetrics {
	m.LastError = err.Error()
	return m
}
