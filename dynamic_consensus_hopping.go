package synnergy

import (
	"sync"
	"time"
)

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
	TPS            float64 // transactions per second
	LatencySec     float64 // average network latency in seconds
	Validators     int     // active validators in the network
	FinalityLagSec float64 // average time for blocks to finalize
	ForkRate       float64 // proportion of competing forks observed
	QueueDepth     int     // pending transactions waiting for inclusion
}

// ConsensusHopper dynamically switches between consensus modes based on metrics.
type ConsensusHopper struct {
	mu          sync.RWMutex
	mode        ConsensusMode
	last        NetworkMetrics
	history     []NetworkMetrics
	maxHistory  int
	stake       *StakingNode
	failover    *FailoverManager
	historyOnce sync.Once
}

// NewConsensusHopper creates a hopper with the supplied initial mode.
func NewConsensusHopper(initial ConsensusMode) *ConsensusHopper {
	return &ConsensusHopper{mode: initial, maxHistory: 10}
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

// AttachStakingNode connects a staking node to the hopper so staking totals can
// influence consensus mode selection. The attachment is optional; if omitted the
// hopper falls back to purely metric-based heuristics.
func (h *ConsensusHopper) AttachStakingNode(node *StakingNode) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.stake = node
}

// AttachFailoverManager connects a failover manager to the hopper. Failover
// stability is incorporated into consensus selection when present so consensus
// switches respond to recent node promotions and heartbeat gaps.
func (h *ConsensusHopper) AttachFailoverManager(manager *FailoverManager) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failover = manager
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
	h.initializeHistory()
	h.last = m
	h.history = append(h.history, m)
	if len(h.history) > h.maxHistory {
		h.history = h.history[1:]
	}

	agg := h.aggregateHistory()
	stakePerValidator := h.stakePerValidator(m.Validators)
	failoverHealthy, failoverUnstable := h.failoverSignals()

	// Prefer PoS when throughput is high, latency/finality are low, stake
	// participation is healthy and the failover layer reports stability.
	if failoverHealthy && stakePerValidator >= 1500 && agg.tps > 1200 && agg.latency < 1.2 && agg.finality < 3 && agg.fork < 0.04 && agg.queue < 1200 {
		h.mode = ConsensusPoS
		return h.mode
	}

	// Allow PoS under slightly relaxed throughput when validators are
	// plentiful, queues manageable and failover hasn't recently churned.
	if failoverHealthy && !failoverUnstable && stakePerValidator >= 800 && agg.tps > 900 && agg.latency < 1.5 && agg.finality < 3.5 && agg.queue < 1500 {
		h.mode = ConsensusPoS
		return h.mode
	}

	// Few validators or signs of finality/fork instability benefit from a
	// PoH mode that prioritizes fast leader rotation with deterministic slots.
	if m.Validators < 12 || agg.finality >= 4 || agg.fork >= 0.07 || agg.queue >= 2000 || agg.latency >= 2.5 {
		h.mode = ConsensusPoH
		return h.mode
	}

	// Recent failover churn benefits from the battle-tested PoW defaults
	// while the rest of the network rebalances and backups catch up.
	if failoverUnstable {
		h.mode = ConsensusPoW
		return h.mode
	}

	// Default to PoW for steady state when neither PoS nor PoH triggers.
	h.mode = ConsensusPoW
	return h.mode
}

func (h *ConsensusHopper) initializeHistory() {
	h.historyOnce.Do(func() {
		if h.maxHistory <= 0 {
			h.maxHistory = 10
		}
		h.history = make([]NetworkMetrics, 0, h.maxHistory)
	})
}

type aggregatedMetrics struct {
	tps      float64
	latency  float64
	finality float64
	fork     float64
	queue    float64
}

func (h *ConsensusHopper) aggregateHistory() aggregatedMetrics {
	var agg aggregatedMetrics
	if len(h.history) == 0 {
		return agg
	}

	for _, entry := range h.history {
		agg.tps += entry.TPS
		agg.latency += entry.LatencySec
		agg.finality += entry.FinalityLagSec
		agg.fork += entry.ForkRate
		agg.queue += float64(entry.QueueDepth)
	}

	denom := float64(len(h.history))
	agg.tps /= denom
	agg.latency /= denom
	agg.finality /= denom
	agg.fork /= denom
	agg.queue /= denom
	return agg
}

func (h *ConsensusHopper) stakePerValidator(validators int) float64 {
	if h.stake == nil || validators <= 0 {
		return 0
	}
	total := h.stake.TotalStaked()
	if total == 0 {
		return 0
	}
	return float64(total) / float64(validators)
}

func (h *ConsensusHopper) failoverSignals() (healthy bool, unstable bool) {
	if h.failover == nil {
		return true, false
	}
	snapshot := h.failover.Snapshot()
	healthy = snapshot.Healthy
	if snapshot.Failovers == 0 {
		return healthy, false
	}

	if !snapshot.Healthy {
		return false, true
	}

	// Consider the failover layer unstable when a promotion occurred
	// within a short window relative to the configured timeout.
	if snapshot.LastSwitch.IsZero() {
		return healthy, false
	}
	elapsed := time.Since(snapshot.LastSwitch)
	threshold := 3 * snapshot.Timeout
	if threshold <= 0 {
		threshold = 3 * time.Second
	}
	return healthy, elapsed < threshold
}
