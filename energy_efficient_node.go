package synnergy

import (
	"sync"
	"time"
)

// SustainabilityCertificate summarises the efficiency metrics for a validator.
type SustainabilityCertificate struct {
	Validator  string
	Efficiency float64 // transactions per kWh
	Offsets    float64 // carbon offset credits
	IssuedAt   time.Time
}

// EnergyEfficientNode records energy usage and computes certificates.
type EnergyEfficientNode struct {
	id            string
	tracker       *EnergyEfficiencyTracker
	mu            sync.RWMutex
	offsetCredits float64
	cert          SustainabilityCertificate
}

// NewEnergyEfficientNode constructs a node bound to an efficiency tracker.
func NewEnergyEfficientNode(id string, tracker *EnergyEfficiencyTracker) *EnergyEfficientNode {
	return &EnergyEfficientNode{id: id, tracker: tracker}
}

// ID returns the node identifier.
func (n *EnergyEfficientNode) ID() string { return n.id }

// RecordUsage records processed transactions and energy consumption for the node.
func (n *EnergyEfficientNode) RecordUsage(txProcessed int, energyKWh float64) {
	n.tracker.Record(n.id, txProcessed, energyKWh)
}

// AddOffset credits carbon offsets to the node.
func (n *EnergyEfficientNode) AddOffset(credits float64) {
	n.mu.Lock()
	n.offsetCredits += credits
	n.mu.Unlock()
}

// OffsetCredits returns the current accumulated carbon offset credits for the node.
func (n *EnergyEfficientNode) OffsetCredits() float64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.offsetCredits
}

// Certify recomputes and stores a sustainability certificate based on latest metrics.
func (n *EnergyEfficientNode) Certify() SustainabilityCertificate {
	eff, _ := n.tracker.Efficiency(n.id)
	n.mu.Lock()
	n.cert = SustainabilityCertificate{
		Validator:  n.id,
		Efficiency: eff,
		Offsets:    n.offsetCredits,
		IssuedAt:   time.Now().UTC(),
	}
	cert := n.cert
	n.mu.Unlock()
	return cert
}

// Certificate returns the last issued certificate.
func (n *EnergyEfficientNode) Certificate() SustainabilityCertificate {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.cert
}

// ShouldThrottle reports whether efficiency falls below the given threshold.
func (n *EnergyEfficientNode) ShouldThrottle(threshold float64) bool {
	eff, ok := n.tracker.Efficiency(n.id)
	if !ok {
		return true
	}
	return eff < threshold
}
