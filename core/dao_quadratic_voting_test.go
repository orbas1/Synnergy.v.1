package core

import "testing"

func TestQuadraticWeight(t *testing.T) {
	if w := QuadraticWeight(9); w != 3 {
		t.Fatalf("expected 3 got %d", w)
	}
}

func TestCastQuadraticVoteZeroTokens(t *testing.T) {
	pm := NewProposalManager()
	dao := NewDAOManager().Create("d", "c")
	p := pm.CreateProposal(dao, "c", "desc")
	if err := pm.CastQuadraticVote(p.ID, "v", 0, true); err == nil {
		t.Fatalf("expected error for zero tokens")
	}
}
