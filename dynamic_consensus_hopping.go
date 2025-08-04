package synnergy

import "sync"

// ConsensusMode represents an available consensus algorithm.
type ConsensusMode string

const (
	// ConsensusPoW represents a proof-of-work style consensus.
	ConsensusPoW ConsensusMode = "pow"
	// ConsensusPoS represents a proof-of-stake consensus.
	ConsensusPoS ConsensusMode = "pos"
	// ConsensusPoH represents a proof-of-history consensus.
	ConsensusPoH ConsensusMode = "poh"
)

// NetworkMetrics captures network observations used to select a consensus mode.
type NetworkMetrics struct {
	TPS        float64 // transactions per second
	LatencySec float64 // average network latency in seconds
	Validators int     // active validators in the network
}

// ConsensusHopper dynamically switches between consensus modes based on metrics.
type ConsensusHopper struct {
	mu   sync.RWMutex
	mode ConsensusMode
	last NetworkMetrics
}

// NewConsensusHopper creates a hopper with the supplied initial mode.
func NewConsensusHopper(initial ConsensusMode) *ConsensusHopper {
	return &ConsensusHopper{mode: initial}
}

// Mode returns the currently selected consensus mode.
func (h *ConsensusHopper) Mode() ConsensusMode {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.mode
}

// SetMode forcibly sets the consensus mode, bypassing evaluation heuristics.
func (h *ConsensusHopper) SetMode(m ConsensusMode) {
	h.mu.Lock()
	h.mode = m
	h.mu.Unlock()
}

// LastMetrics returns the metrics used during the most recent evaluation.
func (h *ConsensusHopper) LastMetrics() NetworkMetrics {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.last
}

// Evaluate selects a consensus mode based on network metrics and returns it.
func (h *ConsensusHopper) Evaluate(m NetworkMetrics) ConsensusMode {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.last = m
	switch {
	case m.TPS > 1000 && m.LatencySec < 1:
		h.mode = ConsensusPoS
	case m.Validators < 10:
		h.mode = ConsensusPoH
	default:
		h.mode = ConsensusPoW
	}
	return h.mode
}
