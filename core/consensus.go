package core

import (
	"math/rand"
	"strings"
	"time"
)

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

	// Availability flags allow weightings to be dropped to zero when no
	// validators are present for a given mechanism or when rewards are
	// exhausted.
	PoWAvailable bool
	PoSAvailable bool
	PoHAvailable bool
	PoWRewards   bool
}

// NewSynnergyConsensus returns a new consensus engine with default parameters
// derived from the Synnergy specification.
func NewSynnergyConsensus() *SynnergyConsensus {
	return &SynnergyConsensus{
		Weights:      ConsensusWeights{PoW: 0.40, PoS: 0.30, PoH: 0.30},
		Alpha:        0.5,
		Beta:         0.5,
		Gamma:        0.1,
		Dmax:         1,
		Smax:         1,
		PoWAvailable: true,
		PoSAvailable: true,
		PoHAvailable: true,
		PoWRewards:   true,
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

	if !sc.PoWAvailable || !sc.PoWRewards {
		sc.Weights.PoW = 0
	}
	if !sc.PoSAvailable {
		sc.Weights.PoS = 0
	}
	if !sc.PoHAvailable {
		sc.Weights.PoH = 0
	}

	total := sc.Weights.PoW + sc.Weights.PoS + sc.Weights.PoH
	if total == 0 {
		return
	}
	sc.Weights.PoW /= total
	sc.Weights.PoS /= total
	sc.Weights.PoH /= total
}

// Tload computes the network-load component of the transition threshold.
func (sc *SynnergyConsensus) Tload(D float64) float64 {
	if sc.Dmax == 0 {
		return 0
	}
	return sc.Alpha * (D / sc.Dmax)
}

// Tsecurity computes the security-threat component of the transition threshold.
func (sc *SynnergyConsensus) Tsecurity(threat float64) float64 {
	if sc.Smax == 0 {
		return 0
	}
	return sc.Beta * (threat / sc.Smax)
}

// Tstake computes the stake-concentration component of the transition threshold.
func (sc *SynnergyConsensus) Tstake(S float64) float64 {
	if sc.Smax == 0 {
		return 0
	}
	return sc.Gamma * (S / sc.Smax)
}

// TransitionThreshold combines load, security and stake factors to determine
// whether a shift in weighting should occur.
func (sc *SynnergyConsensus) TransitionThreshold(D, threat, S float64) float64 {
	return sc.Tload(D) + sc.Tsecurity(threat) + sc.Tstake(S)
}

// DifficultyAdjust recalculates mining difficulty based on the time required to
// mine the previous window of blocks.
func (sc *SynnergyConsensus) DifficultyAdjust(oldDifficulty, actualTime, expectedTime float64) float64 {
	if expectedTime == 0 {
		return oldDifficulty
	}
	return oldDifficulty * (actualTime / expectedTime)
}

// SetAvailability toggles the availability flags for consensus methods.
func (sc *SynnergyConsensus) SetAvailability(pow, pos, poh bool) {
	sc.PoWAvailable = pow
	sc.PoSAvailable = pos
	sc.PoHAvailable = poh
}

// SetPoWRewards indicates whether mining rewards remain for PoW miners.
func (sc *SynnergyConsensus) SetPoWRewards(enabled bool) {
	sc.PoWRewards = enabled
}

// SelectValidator returns a validator address using a weighted random selection
// based on the provided stake map.
func (sc *SynnergyConsensus) SelectValidator(stakes map[string]uint64) string {
	if len(stakes) == 0 {
		return ""
	}
	var total uint64
	for _, s := range stakes {
		total += s
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano())).Uint64() % total
	var cumulative uint64
	for addr, s := range stakes {
		cumulative += s
		if r < cumulative {
			return addr
		}
	}
	return ""
}

// ValidateSubBlock performs POS and POH validation on a sub-block.  For the
// prototype this simply returns true.
func (sc *SynnergyConsensus) ValidateSubBlock(sb *SubBlock) bool {
	return true
}

// MineBlock performs a simple SHA-256 proof-of-work using the provided
// difficulty, defined as the number of leading zeroes required in the block
// hash.
func (sc *SynnergyConsensus) MineBlock(b *Block, difficulty uint8) {
	target := strings.Repeat("0", int(difficulty))
	var nonce uint64
	for {
		hash := b.HeaderHash(nonce)
		if strings.HasPrefix(hash, target) {
			b.Nonce = nonce
			b.Hash = hash
			return
		}
		nonce++
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
