package core

import (
        "context"
        "crypto/sha256"
        "fmt"
        "math"
        "math/big"
        "sort"
        "sync"
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

var maxPoWTarget = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

const (
        expectedBlockIntervalSeconds = 12.0
        stalePenaltyWindow           = 5 * time.Minute
        maxStalePenalty              = 0.45
        diversityBoostCap            = 0.05
)

// SynnergyConsensus encapsulates the consensus algorithms and their dynamic
// weighting.
type SynnergyConsensus struct {
	mu sync.RWMutex

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

// SetRegulatoryNode attaches a regulatory node so that sub-block validation
// includes regulatory transaction checks.
func (sc *SynnergyConsensus) SetRegulatoryNode(rn *RegulatoryNode) {
	sc.mu.Lock()
	sc.RegNode = rn
	sc.mu.Unlock()
}

// Threshold computes the switching threshold based on network demand (D) and
// stake concentration (S).
func (sc *SynnergyConsensus) Threshold(D, S float64) float64 {
	sc.mu.RLock()
	alpha := sc.Alpha
	beta := sc.Beta
	sc.mu.RUnlock()
	return alpha*D + beta*S
}

// AdjustWeights modifies the internal consensus weightings based on current
// network demand (D) and stake concentration (S).  Adjustments are smoothed and
// normalized so that the weights always sum to one and never drop below 7.5%
// unless network conditions force them lower.
func (sc *SynnergyConsensus) AdjustWeights(D, S float64) {
	sc.mu.Lock()
	if sc.Dmax == 0 {
		sc.Dmax = 1
	}
	if sc.Smax == 0 {
		sc.Smax = 1
	}
	weights := sc.Weights
	weights.PoW = clamp(weights.PoW+sc.Gamma*((D/sc.Dmax)-weights.PoW), 0.075, 1)
	weights.PoS = clamp(weights.PoS+sc.Gamma*((S/sc.Smax)-weights.PoS), 0.075, 1)
	loadInv := 1 - (D / sc.Dmax)
	weights.PoH = clamp(weights.PoH+sc.Gamma*(loadInv-weights.PoH), 0.075, 1)

	if !sc.PoWAvailable || !sc.PoWRewards {
		weights.PoW = 0
	}
	if !sc.PoSAvailable {
		weights.PoS = 0
	}
	if !sc.PoHAvailable {
		weights.PoH = 0
	}

	weights = normalizeWeights(weights)
	sc.Weights = weights
	sc.mu.Unlock()
	ilog.Info("adjust_weights", "pow", weights.PoW, "pos", weights.PoS, "poh", weights.PoH)
}

// Tload computes the network-load component of the transition threshold.
func (sc *SynnergyConsensus) Tload(D float64) float64 {
	sc.mu.RLock()
	dmax := sc.Dmax
	alpha := sc.Alpha
	sc.mu.RUnlock()
	if dmax == 0 {
		return 0
	}
	return alpha * (D / dmax)
}

// Tsecurity computes the security-threat component of the transition threshold.
func (sc *SynnergyConsensus) Tsecurity(threat float64) float64 {
	sc.mu.RLock()
	smax := sc.Smax
	beta := sc.Beta
	sc.mu.RUnlock()
	if smax == 0 {
		return 0
	}
	return beta * (threat / smax)
}

// Tstake computes the stake-concentration component of the transition threshold.
func (sc *SynnergyConsensus) Tstake(S float64) float64 {
	sc.mu.RLock()
	smax := sc.Smax
	gamma := sc.Gamma
	sc.mu.RUnlock()
	if smax == 0 {
		return 0
	}
	return gamma * (S / smax)
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
	sc.mu.Lock()
	sc.PoWAvailable = pow
	sc.PoSAvailable = pos
	sc.PoHAvailable = poh
	sc.mu.Unlock()
	ilog.Info("set_availability", "pow", pow, "pos", pos, "poh", poh)
}

// SetPoWRewards indicates whether mining rewards remain for PoW miners.
func (sc *SynnergyConsensus) SetPoWRewards(enabled bool) {
	sc.mu.Lock()
	sc.PoWRewards = enabled
	sc.mu.Unlock()
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
	type candidate struct {
		addr  string
		stake uint64
	}
	candidates := make([]candidate, 0, len(stakes))
	var total uint64
	for addr, stake := range stakes {
		if stake == 0 {
			continue
		}
		candidates = append(candidates, candidate{addr: addr, stake: stake})
		total += stake
	}
	if total == 0 {
		ilog.Info("select_validator", "result", "")
		return ""
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].addr < candidates[j].addr
	})
	h := sha256.Sum256([]byte(seed))
	entropy := new(big.Int).SetBytes(h[:])
	pick := new(big.Int).Mod(entropy, new(big.Int).SetUint64(total)).Uint64()
	var cumulative uint64
	for _, cand := range candidates {
		cumulative += cand.stake
		if pick < cumulative {
			ilog.Info("select_validator", "result", cand.addr, "stake", cand.stake, "total", total)
			return cand.addr
		}
	}
	chosen := candidates[len(candidates)-1].addr
	ilog.Info("select_validator", "result", chosen, "stake", candidates[len(candidates)-1].stake, "total", total)
	return chosen
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
	regNode := sc.getRegNode()
	if regNode == nil {
		// No regulatory node configured; bypass compliance checks.
		return true
	}
	for _, tx := range sb.Transactions {
		if err := regNode.ApproveTransaction(*tx); err != nil {
			ilog.Info("regulatory_reject", "tx", tx.ID, "err", err)
			return false
		}
	}
	return true
}

// MineBlock performs a simple SHA-256 proof-of-work using the provided
// difficulty, defined as the number of leading zeroes required in the block
// hash.
func (sc *SynnergyConsensus) MineBlock(b *Block, difficulty uint8) {
	const maxDifficultyBits = 64 // clamp work so the stub miner cannot hang indefinitely
	bits := int(difficulty) * 4
	if bits > maxDifficultyBits {
		bits = maxDifficultyBits
	}
	shift := 256 - bits
	if shift < 1 {
		shift = 1
	}
	target := new(big.Int).Lsh(big.NewInt(1), uint(shift))
	target.Sub(target, big.NewInt(1))
	start := time.Now()
	var nonce uint64
	for {
		hash := b.HeaderHash(nonce)
		hashInt, ok := new(big.Int).SetString(hash, 16)
		if ok && hashInt.Cmp(target) <= 0 {
			b.Nonce = nonce
			b.Hash = hash
			ilog.Info("mine_block", "nonce", nonce, "difficulty_bits", bits)
			return
		}
		nonce++
		if nonce == 0 || time.Since(start) > time.Second {
			// Refresh the timestamp to alter the work target and avoid infinite loops.
			b.Timestamp = time.Now().Unix()
			start = time.Now()
		}
	}
}

// FinalizeBlock applies a simple BFT-style vote on the block. If at least two
// thirds of votes are affirmative the block is marked finalized and validators
// contributing sub-blocks receive a stake reward via the provided manager.
func (sc *SynnergyConsensus) FinalizeBlock(b *Block, votes map[string]bool, vm *ValidatorManager, reward uint64) bool {
	if b == nil || vm == nil {
		return false
	}
	eligible := vm.Eligible()
	if len(eligible) == 0 {
		return false
	}
	var totalStake uint64
	for _, stake := range eligible {
		totalStake += stake
	}
	if totalStake == 0 {
		return false
	}
	var participatingStake uint64
	var affirmativeStake uint64
	for addr, approve := range votes {
		stake, ok := eligible[addr]
		if !ok || stake == 0 {
			continue
		}
		participatingStake += stake
		if approve {
			affirmativeStake += stake
		}
	}
	if participatingStake*3 < totalStake*2 {
		ilog.Info("finalize_block", "hash", b.Hash, "result", "insufficient_participation", "participating", participatingStake, "total", totalStake)
		return false
	}
	if affirmativeStake*3 < totalStake*2 {
		ilog.Info("finalize_block", "hash", b.Hash, "result", "insufficient_affirmative", "yes", affirmativeStake, "total", totalStake)
		return false
	}
	b.Finalized = true
	if reward > 0 {
		for _, sb := range b.SubBlocks {
			if sb == nil {
				continue
			}
			vm.Reward(context.Background(), sb.Validator, reward)
		}
	}
	ilog.Info("finalize_block", "hash", b.Hash, "yes_stake", affirmativeStake, "total_stake", totalStake)
	return true
}

// ChooseChain selects the longest chain from the candidates. This placeholder
// fork-choice rule enables nodes to converge on a canonical history.
func (sc *SynnergyConsensus) ChooseChain(chains [][]*Block) []*Block {
        if len(chains) == 0 {
                ilog.Info("fork_choice", "score", 0.0, "reason", "no_candidates")
                return nil
        }

        weights := normalizeWeights(sc.WeightsSnapshot())
        var (
                chosen   []*Block
                bestEval chainEvaluation
                haveBest bool
        )

        for _, candidate := range chains {
                eval, err := sc.evaluateChain(candidate, weights)
                if err != nil {
                        ilog.Info("fork_choice_skip", "reason", err.Error())
                        continue
                }
                if !haveBest || eval.betterThan(bestEval) {
                        chosen = candidate
                        bestEval = eval
                        haveBest = true
                }
        }

        if !haveBest {
                ilog.Info("fork_choice", "score", 0.0, "reason", "no_valid_chain")
                return nil
        }

        ilog.Info(
                "fork_choice",
                "score", bestEval.score,
                "pow_component", bestEval.powComponent,
                "pos_component", bestEval.posComponent,
                "poh_component", bestEval.pohComponent,
                "finalized", bestEval.finalized,
                "length", bestEval.length,
                "unique_validators", bestEval.uniqueValidators,
                "tip_time", bestEval.tipTime.Unix(),
        )

        return chosen
}

type chainEvaluation struct {
        score            float64
        powComponent     float64
        posComponent     float64
        pohComponent     float64
        finalized        int
        length           int
        tipTime          time.Time
        tipHash          string
        uniqueValidators int
        participation    float64
}

func (e chainEvaluation) betterThan(other chainEvaluation) bool {
        if e.score != other.score {
                return e.score > other.score
        }
        if e.finalized != other.finalized {
                return e.finalized > other.finalized
        }
        if e.length != other.length {
                return e.length > other.length
        }
        switch {
        case other.tipTime.IsZero() && !e.tipTime.IsZero():
                return true
        case e.tipTime.IsZero() && !other.tipTime.IsZero():
                return false
        case !e.tipTime.Equal(other.tipTime):
                return e.tipTime.After(other.tipTime)
        }
        return e.tipHash > other.tipHash
}

func (sc *SynnergyConsensus) evaluateChain(chain []*Block, weights ConsensusWeights) (chainEvaluation, error) {
        if len(chain) == 0 {
                return chainEvaluation{}, fmt.Errorf("empty chain")
        }

        eval := chainEvaluation{length: len(chain)}
        validators := make(map[string]struct{})

        var (
                powAccum      float64
                totalSubBlock int
                prevHash      string
                prevTimestamp int64
        )

        for i, blk := range chain {
                if blk == nil {
                        return chainEvaluation{}, fmt.Errorf("nil block at index %d", i)
                }
                if err := blk.Validate(); err != nil {
                        return chainEvaluation{}, fmt.Errorf("block %d invalid: %w", i, err)
                }
                if i > 0 {
                        if blk.PrevHash != prevHash {
                                return chainEvaluation{}, fmt.Errorf("block %d prev hash mismatch", i)
                        }
                        if blk.Timestamp < prevTimestamp {
                                return chainEvaluation{}, fmt.Errorf("block %d timestamp regressed", i)
                        }
                }
                if blk.Hash != "" {
                        powAccum += powQuality(blk.Hash)
                } else if blk.Nonce > 0 {
                        powAccum += 1 / float64(blk.Nonce)
                }
                if blk.Finalized {
                        eval.finalized++
                }
                for _, sb := range blk.SubBlocks {
                        if sb == nil {
                                continue
                        }
                        totalSubBlock++
                        validators[sb.Validator] = struct{}{}
                }
                prevHash = blk.Hash
                prevTimestamp = blk.Timestamp
        }

        eval.uniqueValidators = len(validators)
        if last := chain[len(chain)-1]; last != nil {
                eval.tipHash = last.Hash
                if last.Timestamp > 0 {
                        eval.tipTime = time.Unix(last.Timestamp, 0)
                }
        }

        eval.powComponent = clamp(powAccum/float64(len(chain)), 0, 1)
        if totalSubBlock > 0 {
                eval.posComponent = clamp(float64(len(validators))/float64(totalSubBlock), 0, 1)
                eval.participation = eval.posComponent
        }
        eval.pohComponent = pohContinuityScore(chain)
        if totalSubBlock > 0 {
                avg := float64(totalSubBlock) / float64(len(chain))
                eval.pohComponent = 0.7*eval.pohComponent + 0.3*math.Min(1, avg/8)
        }

        score := weights.PoW*eval.powComponent + weights.PoS*eval.posComponent + weights.PoH*eval.pohComponent
        if eval.length > 0 {
                score += 0.2 * float64(eval.finalized) / float64(eval.length)
        }

        if !eval.tipTime.IsZero() {
                if age := time.Since(eval.tipTime); age > 0 {
                        penalty := (age.Seconds() / stalePenaltyWindow.Seconds()) * 0.05
                        if penalty > maxStalePenalty {
                                penalty = maxStalePenalty
                        }
                        score -= penalty
                }
        }

        diversity := math.Min(float64(len(validators))/25.0, diversityBoostCap)
        score += diversity

        if math.IsNaN(score) || math.IsInf(score, 0) {
                return chainEvaluation{}, fmt.Errorf("invalid score computed")
        }

        eval.score = score
        return eval, nil
}

func powQuality(hash string) float64 {
        if hash == "" {
                return 0
        }
        value := new(big.Int)
        if _, ok := value.SetString(hash, 16); !ok {
                return 0
        }
        if value.Sign() <= 0 {
                return 1
        }
        remaining := new(big.Int).Sub(maxPoWTarget, value)
        if remaining.Sign() <= 0 {
                return 0
        }
        ratio := new(big.Float).Quo(new(big.Float).SetInt(remaining), new(big.Float).SetInt(maxPoWTarget))
        f, _ := ratio.Float64()
        if f < 0 {
                return 0
        }
        if f > 1 {
                return 1
        }
        return f
}

func pohContinuityScore(chain []*Block) float64 {
        if len(chain) <= 1 {
                return 1
        }
        var (
                total float64
                prev  = chain[0].Timestamp
        )
        for i := 1; i < len(chain); i++ {
                ts := chain[i].Timestamp
                if ts <= prev {
                        return 0
                }
                delta := float64(ts - prev)
                if delta <= 0 {
                        return 0
                }
                ratio := delta / expectedBlockIntervalSeconds
                if ratio > 1 {
                        ratio = 1 / ratio
                }
                if ratio > 1 {
                        ratio = 1
                }
                total += ratio
                prev = ts
        }
        return clamp(total/float64(len(chain)-1), 0, 1)
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

// WeightsSnapshot returns the current consensus weight distribution.
func (sc *SynnergyConsensus) WeightsSnapshot() ConsensusWeights {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.Weights
}

// SetWeights overrides the consensus weights, normalising the provided values.
func (sc *SynnergyConsensus) SetWeights(w ConsensusWeights) {
	sc.mu.Lock()
	sc.Weights = normalizeWeights(w)
	sc.mu.Unlock()
}

// ValidateBlock checks that a block and its sub-blocks are well-formed.
// It delegates to the block's internal validation routine.
func (sc *SynnergyConsensus) ValidateBlock(b *Block) bool {
	if b == nil {
		return false
	}
	return b.Validate() == nil
}

func (sc *SynnergyConsensus) getRegNode() *RegulatoryNode {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.RegNode
}

func normalizeWeights(w ConsensusWeights) ConsensusWeights {
	if w.PoW < 0 {
		w.PoW = 0
	}
	if w.PoS < 0 {
		w.PoS = 0
	}
	if w.PoH < 0 {
		w.PoH = 0
	}
	total := w.PoW + w.PoS + w.PoH
	if total == 0 {
		return ConsensusWeights{}
	}
	return ConsensusWeights{
		PoW: w.PoW / total,
		PoS: w.PoS / total,
		PoH: w.PoH / total,
	}
}
