package core

import "testing"

func TestQuadraticWeight(t *testing.T) {
	if w := QuadraticWeight(9); w != 3 {
		t.Fatalf("expected 3 got %d", w)
	}
}

func TestCastQuadraticVoteZeroTokens(t *testing.T) {
	pm := NewProposalManager()
	mgr := NewDAOManager()
	mgr.AuthorizeRelayer("c")
	dao, err := mgr.Create("d", "c")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	p, err := pm.CreateProposal(dao, "c", "desc")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if err := pm.CastQuadraticVote(dao, p.ID, "v", 0, true); err == nil {
		t.Fatalf("expected error for zero tokens")
	}
}

func TestCastQuadraticVoteRequiresMembership(t *testing.T) {
	pm := NewProposalManager()
	mgr := NewDAOManager()
	mgr.AuthorizeRelayer("c")
	dao, err := mgr.Create("d", "c")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	p, err := pm.CreateProposal(dao, "c", "desc")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if err := pm.CastQuadraticVote(dao, p.ID, "v", 4, true); err != errNotMember {
		t.Fatalf("expected errNotMember got %v", err)
	}
}

func TestCastQuadraticVoteSuccess(t *testing.T) {
	pm := NewProposalManager()
	mgr := NewDAOManager()
	mgr.AuthorizeRelayer("c")
	dao, err := mgr.Create("d", "c")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := dao.AddMember("v", RoleMember); err != nil {
		t.Fatalf("add member: %v", err)
	}
	p, err := pm.CreateProposal(dao, "c", "desc")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if err := pm.CastQuadraticVote(dao, p.ID, "v", 9, true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	yes, no, err := pm.Results(p.ID)
	if err != nil {
		t.Fatalf("results: %v", err)
	}
	expected := QuadraticWeight(9)
	if yes != expected || no != 0 {
		t.Fatalf("unexpected results: %d %d", yes, no)
	}
}
