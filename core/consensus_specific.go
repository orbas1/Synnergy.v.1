package core

// ConsensusMode identifies the active consensus mechanism for a node.
type ConsensusMode string

const (
	// ModePoW selects proof-of-work validation.
	ModePoW ConsensusMode = "pow"
	// ModePoS selects proof-of-stake validation.
	ModePoS ConsensusMode = "pos"
	// ModePoH selects proof-of-history validation.
	ModePoH ConsensusMode = "poh"
)

// ConsensusSwitcher evaluates consensus weights and picks the dominant mode.
type ConsensusSwitcher struct {
	mode ConsensusMode
}

// NewConsensusSwitcher creates a switcher with the provided starting mode.
func NewConsensusSwitcher(mode ConsensusMode) *ConsensusSwitcher {
	return &ConsensusSwitcher{mode: mode}
}

// Evaluate updates the current mode based on the highest weight in the
// consensus engine and returns the selected mode.
func (cs *ConsensusSwitcher) Evaluate(sc *SynnergyConsensus) ConsensusMode {
	if sc == nil {
		return cs.mode
	}
	snapshot := sc.WeightsSnapshot()
	weights := map[ConsensusMode]float64{
		ModePoW: snapshot.PoW,
		ModePoS: snapshot.PoS,
		ModePoH: snapshot.PoH,
	}
	var maxMode ConsensusMode
	maxWeight := -1.0
	for m, w := range weights {
		if w > maxWeight {
			maxWeight = w
			maxMode = m
		}
	}
	cs.mode = maxMode
	return cs.mode
}

// Mode returns the last evaluated consensus mode.
func (cs *ConsensusSwitcher) Mode() ConsensusMode {
	return cs.mode
}
