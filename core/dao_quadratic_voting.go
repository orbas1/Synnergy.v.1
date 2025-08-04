package core

import "math"

// QuadraticWeight returns the quadratic voting weight for a token amount.
func QuadraticWeight(tokens uint64) uint64 {
	return uint64(math.Sqrt(float64(tokens)))
}

// CastQuadraticVote records a quadratic vote on a proposal.
func (pm *ProposalManager) CastQuadraticVote(id, voter string, tokens uint64, support bool) error {
	weight := QuadraticWeight(tokens)
	return pm.Vote(id, voter, weight, support)
}
