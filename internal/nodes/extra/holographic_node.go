package nodes

import (
	"errors"
	"sync"
	"time"

	"synnergy"
)

// HolographicNode provides holographic data distribution and redundancy.
type HolographicNode struct {
	id         string
	mu         sync.RWMutex
	store      map[string]synnergy.HolographicFrame
	accessTime map[string]time.Time
	cfg        holographicConfig
}

type holographicConfig struct {
	maxFrames int
	retention time.Duration
}

// HolographicOption configures optional behaviour for holographic nodes.
type HolographicOption interface {
	applyHolographicOption(*holographicConfig)
}

type holographicOptionFunc func(*holographicConfig)

func (f holographicOptionFunc) applyHolographicOption(cfg *holographicConfig) { f(cfg) }

// WithFrameLimit restricts the number of frames held in memory.
func WithFrameLimit(limit int) HolographicOption {
	return holographicOptionFunc(func(cfg *holographicConfig) {
		if limit > 0 {
			cfg.maxFrames = limit
		}
	})
}

// WithFrameRetention configures a retention window for frames.
func WithFrameRetention(window time.Duration) HolographicOption {
	return holographicOptionFunc(func(cfg *holographicConfig) {
		if window > 0 {
			cfg.retention = window
		}
	})
}

// NewHolographicNode creates a new HolographicNode with the given identifier.
func NewHolographicNode(id string, opts ...HolographicOption) *HolographicNode {
	cfg := holographicConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt.applyHolographicOption(&cfg)
		}
	}
	return &HolographicNode{
		id:         id,
		store:      make(map[string]synnergy.HolographicFrame),
		accessTime: make(map[string]time.Time),
		cfg:        cfg,
	}
}

// ID returns the node identifier.
func (n *HolographicNode) ID() string { return n.id }

// Start implements the NodeInterface; holographic nodes currently have no
// background processes so Start validates configuration.
func (n *HolographicNode) Start() error {
	if n.cfg.maxFrames < 0 {
		return errors.New("max frames cannot be negative")
	}
	return nil
}

// Stop implements the NodeInterface; holographic nodes currently have no
// background processes so Stop is a no-op.
func (n *HolographicNode) Stop() error { return nil }

// Store saves a holographic frame in the node's internal storage.
func (n *HolographicNode) Store(frame synnergy.HolographicFrame) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.store[frame.ID] = frame
	n.accessTime[frame.ID] = time.Now().UTC()
	n.enforcePolicies()
}

// Retrieve fetches a holographic frame by ID. The returned boolean indicates
// whether the frame was found.
func (n *HolographicNode) Retrieve(id string) (synnergy.HolographicFrame, bool) {
	n.mu.RLock()
	frame, ok := n.store[id]
	n.mu.RUnlock()
	if ok {
		n.mu.Lock()
		n.accessTime[id] = time.Now().UTC()
		n.mu.Unlock()
	}
	return frame, ok
}

// Frames returns the identifiers for all stored frames.
func (n *HolographicNode) Frames() []string {
	n.mu.RLock()
	ids := make([]string, 0, len(n.store))
	for id := range n.store {
		ids = append(ids, id)
	}
	n.mu.RUnlock()
	return ids
}

func (n *HolographicNode) enforcePolicies() {
	if n.cfg.retention > 0 {
		cutoff := time.Now().UTC().Add(-n.cfg.retention)
		for id, ts := range n.accessTime {
			if ts.Before(cutoff) {
				delete(n.store, id)
				delete(n.accessTime, id)
			}
		}
	}
	if n.cfg.maxFrames > 0 && len(n.store) > n.cfg.maxFrames {
		// Evict least recently accessed frame.
		var lruID string
		var lruTime time.Time
		for id, ts := range n.accessTime {
			if lruID == "" || ts.Before(lruTime) {
				lruID = id
				lruTime = ts
			}
		}
		delete(n.store, lruID)
		delete(n.accessTime, lruID)
	}
}
