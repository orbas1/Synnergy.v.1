package core

import "math"

// SubBlock groups transactions validated via POS and POH.
type SubBlock struct {
	Transactions []*Transaction
	Validator    string
}

// Block is composed of multiple SubBlocks and finalized via POW.
type Block struct {
	SubBlocks []SubBlock
	Nonce     uint64
}

// ConsensusWeights holds the relative weights assigned to each consensus
// mechanism.  Values are represented as percentages that sum to 1.0.
type ConsensusWeights struct {
	PoW float64
	PoS float64
	PoH float64
}

// SynnergyConsensus encapsulates the consensus algorithms and their dynamic
// weighting.
type SynnergyConsensus struct {
	Weights ConsensusWeights
	Alpha   float64
	Beta    float64
	Gamma   float64
	Dmax    float64
	Smax    float64
}

// NewSynnergyConsensus returns a new consensus engine with default parameters
// derived from the Synnergy specification.
func NewSynnergyConsensus() *SynnergyConsensus {
	return &SynnergyConsensus{
		Weights: ConsensusWeights{PoW: 0.40, PoS: 0.30, PoH: 0.30},
		Alpha:   0.5,
		Beta:    0.5,
		Gamma:   0.1,
		Dmax:    1,
		Smax:    1,
	}
}

// Threshold computes the switching threshold based on network demand (D) and
// stake concentration (S).
func (sc *SynnergyConsensus) Threshold(D, S float64) float64 {
	return sc.Alpha*D + sc.Beta*S
}

// AdjustWeights modifies the internal consensus weightings based on current
// network demand (D) and stake concentration (S).  Adjustments are smoothed and
// normalized so that the weights always sum to one and never drop below 7.5%
// unless network conditions force them lower.
func (sc *SynnergyConsensus) AdjustWeights(D, S float64) {
	adj := sc.Gamma * ((D / sc.Dmax) + (S / sc.Smax))
	sc.Weights.PoW = clamp(sc.Weights.PoW+adj, 0.075, 1)
	sc.Weights.PoS = clamp(sc.Weights.PoS+adj, 0.075, 1)
	sc.Weights.PoH = clamp(sc.Weights.PoH+adj, 0.075, 1)
	total := sc.Weights.PoW + sc.Weights.PoS + sc.Weights.PoH
	sc.Weights.PoW /= total
	sc.Weights.PoS /= total
	sc.Weights.PoH /= total
}

// ValidateSubBlock performs POS and POH validation on a sub-block.  For the
// prototype this simply returns true.
func (sc *SynnergyConsensus) ValidateSubBlock(sb *SubBlock) bool {
	return true
}

// MineBlock performs POW to finalize a block.  This prototype implements a
// trivial nonce incrementer rather than a full hashing routine.
func (sc *SynnergyConsensus) MineBlock(b *Block, difficulty uint64) {
	target := uint64(math.MaxUint64) / difficulty
	for b.Nonce < target {
		b.Nonce++
	}
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
