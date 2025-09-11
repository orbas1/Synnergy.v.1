package core

import (
	"context"
	"crypto/sha256"
	"math/big"
	"strings"

	ilog "synnergy/internal/log"
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
	if sc.Dmax == 0 {
		sc.Dmax = 1
	}
	if sc.Smax == 0 {
		sc.Smax = 1
	}
	sc.Weights.PoW = clamp(sc.Weights.PoW+sc.Gamma*((D/sc.Dmax)-sc.Weights.PoW), 0.075, 1)
	sc.Weights.PoS = clamp(sc.Weights.PoS+sc.Gamma*((S/sc.Smax)-sc.Weights.PoS), 0.075, 1)
	loadInv := 1 - (D / sc.Dmax)
	sc.Weights.PoH = clamp(sc.Weights.PoH+sc.Gamma*(loadInv-sc.Weights.PoH), 0.075, 1)

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
	ilog.Info("adjust_weights", "pow", sc.Weights.PoW, "pos", sc.Weights.PoS, "poh", sc.Weights.PoH)
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
	ilog.Info("set_availability", "pow", pow, "pos", pos, "poh", poh)
}

// SetPoWRewards indicates whether mining rewards remain for PoW miners.
func (sc *SynnergyConsensus) SetPoWRewards(enabled bool) {
	sc.PoWRewards = enabled
	ilog.Info("set_pow_rewards", "enabled", enabled)
}

// SelectValidator deterministically selects a validator using a VRF-style
// hash of the provided seed and validator address. The validator with the
// smallest weighted hash value wins, ensuring all nodes arrive at the same
// result.
func (sc *SynnergyConsensus) SelectValidator(seed string, stakes map[string]uint64) string {
	if len(stakes) == 0 {
		ilog.Info("select_validator", "result", "")
		return ""
	}
	var total, max uint64
	for _, s := range stakes {
		total += s
		if s > max {
			max = s
		}
	}
	if len(stakes) > 1 && max*2 > total {
		ilog.Info("select_validator", "result", "")
		return ""
	}
	var bestAddr string
	var bestScore *big.Int
	for addr, stake := range stakes {
		h := sha256.Sum256([]byte(seed + addr))
		score := new(big.Int).SetBytes(h[:])
		if stake > 0 {
			score.Div(score, new(big.Int).SetUint64(stake))
		}
		if bestAddr == "" || score.Cmp(bestScore) < 0 {
			bestAddr = addr
			bestScore = score
		}
	}
	ilog.Info("select_validator", "result", bestAddr)
	return bestAddr
}

// ValidateSubBlock performs simple PoS and PoH validation on a sub-block.
// It verifies that the sub-block is non-nil, contains transactions and has a
// valid signature from its declared validator.
func (sc *SynnergyConsensus) ValidateSubBlock(sb *SubBlock) bool {
	if sb == nil {
		return false
	}
	return sb.Validate() == nil
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
			ilog.Info("mine_block", "nonce", nonce)
			return
		}
		nonce++

	}
}

// FinalizeBlock applies a simple BFT-style vote on the block. If at least two
// thirds of votes are affirmative the block is marked finalized and validators
// contributing sub-blocks receive a stake reward via the provided manager.
func (sc *SynnergyConsensus) FinalizeBlock(b *Block, votes map[string]bool, vm *ValidatorManager, reward uint64) bool {
	yes := 0
	for _, v := range votes {
		if v {
			yes++
		}
	}
	if yes*3 >= len(votes)*2 {
		b.Finalized = true
		if vm != nil && reward > 0 {
			for _, sb := range b.SubBlocks {
				vm.Reward(context.Background(), sb.Validator, reward)
			}
		}
		ilog.Info("finalize_block", "hash", b.Hash)
		return true
	}
	return false
}

// ChooseChain selects the longest chain from the candidates. This placeholder
// fork-choice rule enables nodes to converge on a canonical history.
func (sc *SynnergyConsensus) ChooseChain(chains [][]*Block) []*Block {
	var best []*Block
	max := 0
	for _, c := range chains {
		if len(c) > max {
			best = c
			max = len(c)
		}
	}
	ilog.Info("fork_choice", "length", max)
	return best
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

// ValidateBlock checks that a block and its sub-blocks are well-formed.
// It delegates to the block's internal validation routine.
func (sc *SynnergyConsensus) ValidateBlock(b *Block) bool {
	if b == nil {
		return false
	}
	return b.Validate() == nil
}
