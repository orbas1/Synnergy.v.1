package core

import (
	"errors"
	"math"
)

// QuadraticWeight returns the quadratic voting weight for a token amount.
func QuadraticWeight(tokens uint64) uint64 {
	return uint64(math.Sqrt(float64(tokens)))
}

// CastQuadraticVote records a quadratic vote on a proposal. A vote with zero
// tokens is rejected to avoid no-op entries.
func (pm *ProposalManager) CastQuadraticVote(dao *DAO, id, voter string, tokens uint64, support bool) error {
	if !dao.IsMember(voter) {
		return errNotMember
	}
	if tokens == 0 {
		return errors.New("tokens must be > 0")
	}
	weight := QuadraticWeight(tokens)
	return pm.Vote(dao, id, voter, weight, support)
}
