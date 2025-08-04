package core

// ConsensusSpecificNode represents a node locked to a single consensus mode.
type ConsensusSpecificNode struct {
	*Node
	Mode ConsensusMode
}

// NewConsensusSpecificNode creates a node with the supplied consensus mode.
func NewConsensusSpecificNode(mode ConsensusMode, id, addr string, ledger *Ledger) *ConsensusSpecificNode {
	n := NewNode(id, addr, ledger)
	csn := &ConsensusSpecificNode{Node: n, Mode: mode}
	csn.configure()
	return csn
}

// configure adjusts the underlying consensus engine to only allow the
// specified mode.
func (n *ConsensusSpecificNode) configure() {
	switch n.Mode {
	case ModePoW:
		n.Consensus.SetAvailability(true, false, false)
		n.Consensus.Weights = ConsensusWeights{PoW: 1}
	case ModePoS:
		n.Consensus.SetAvailability(false, true, false)
		n.Consensus.Weights = ConsensusWeights{PoS: 1}
	case ModePoH:
		n.Consensus.SetAvailability(false, false, true)
		n.Consensus.Weights = ConsensusWeights{PoH: 1}
	}
}
