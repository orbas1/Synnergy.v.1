package core

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

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

	// RegNode performs regulatory checks on transactions during
	// consensus validation. When nil, regulatory checks are bypassed.
	RegNode *RegulatoryNode

	// PoWTimeout bounds the amount of time spent searching for a valid
	// nonce before aborting the mining attempt. A zero value falls back
	// to powDefaultTimeout.
	PoWTimeout time.Duration
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
		PoWTimeout:   10 * time.Second,
	}
}

// SetRegulatoryNode attaches a regulatory node so that sub-block validation
// includes regulatory transaction checks.
func (sc *SynnergyConsensus) SetRegulatoryNode(rn *RegulatoryNode) {
	sc.RegNode = rn
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
	if err := sb.Validate(); err != nil {
		return false
	}
	if sc.RegNode == nil {
		// No regulatory node configured; bypass compliance checks.
		return true
	}
	for _, tx := range sb.Transactions {
		if err := sc.RegNode.ApproveTransaction(*tx); err != nil {
			ilog.Info("regulatory_reject", "tx", tx.ID, "err", err)
			return false
		}
	}
	return true
}

var (
	errMiningCancelled = errors.New("mining cancelled")
	errMiningTimeout   = errors.New("mining timed out")
)

const powCheckInterval = 1 << 12
const powDefaultTimeout = 10 * time.Second

var maxPoWTarget = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

// MineBlock performs a SHA-256 proof-of-work using the provided difficulty.
// Difficulty is defined as the number of leading zeroes required in the block
// hash. The search honours context cancellation and a configurable timeout so
// callers can bound resource usage. A nil block or zero difficulty returns an
// error.
func (sc *SynnergyConsensus) MineBlock(ctx context.Context, b *Block, difficulty uint8) error {
	if b == nil {
		return fmt.Errorf("mine block: %w", ErrNilBlock)
	}
	if difficulty == 0 {
		return fmt.Errorf("mine block: difficulty must be > 0")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	target := strings.Repeat("0", int(difficulty))
	timeout := sc.PoWTimeout
	if timeout <= 0 {
		timeout = powDefaultTimeout
	}
	deadline := time.Now().Add(timeout)

	var nonce uint64
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("mine block: %w: %w", errMiningCancelled, ctx.Err())
		default:
		}

		hash := b.HeaderHash(nonce)
		if strings.HasPrefix(hash, target) {
			b.Nonce = nonce
			b.Hash = hash
			ilog.Info("mine_block", "nonce", nonce)
			return nil
		}

		nonce++
		if nonce == 0 { // overflowed
			return fmt.Errorf("mine block: nonce space exhausted")
		}

		if nonce%powCheckInterval == 0 && time.Now().After(deadline) {
			return fmt.Errorf("mine block: %w after %d iterations", errMiningTimeout, nonce)
		}
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

// ChooseChain selects the highest scoring chain from the provided candidates.
// Candidates that are discontinuous or contain nil blocks are ignored. Chain
// quality is assessed using height, number of finalized blocks, and cumulative
// proof-of-work to provide deterministic tie-breaking between competing forks.
func (sc *SynnergyConsensus) ChooseChain(chains [][]*Block) []*Block {
	var best []*Block
	bestScore := chainScore{work: big.NewInt(0)}
	for _, chain := range chains {
		if !isContinuousChain(chain) {
			continue
		}
		score := computeChainScore(chain)
		if best == nil || score.betterThan(bestScore) {
			best = chain
			bestScore = score
		}
	}
	ilog.Info("fork_choice", "length", bestScore.length, "finalized", bestScore.finalized, "work", bestScore.work.String())
	return best
}

type chainScore struct {
	length    int
	finalized int
	work      *big.Int
}

func (s chainScore) betterThan(other chainScore) bool {
	if s.length != other.length {
		return s.length > other.length
	}
	if s.finalized != other.finalized {
		return s.finalized > other.finalized
	}
	if other.work == nil {
		return true
	}
	return s.work.Cmp(other.work) > 0
}

func computeChainScore(chain []*Block) chainScore {
	score := chainScore{length: len(chain), work: big.NewInt(0)}
	for _, b := range chain {
		if b == nil {
			continue
		}
		if b.Finalized {
			score.finalized++
		}
		if hash := strings.TrimSpace(b.Hash); hash != "" {
			if val, ok := new(big.Int).SetString(hash, 16); ok {
				work := new(big.Int).Sub(maxPoWTarget, val)
				if work.Sign() > 0 {
					score.work.Add(score.work, work)
				}
			}
		}
	}
	return score
}

func isContinuousChain(chain []*Block) bool {
	if len(chain) == 0 {
		return false
	}
	for i := range chain {
		if chain[i] == nil {
			return false
		}
		if i == 0 {
			continue
		}
		prev := chain[i-1]
		if prev == nil || chain[i].PrevHash != prev.Hash {
			return false
		}
	}
	return true
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
