package core

import "testing"

func TestDAOProposal(t *testing.T) {
	mgr := NewDAOManager()
	mgr.AuthorizeRelayer("c")
	dao, err := mgr.Create("d", "c")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	pm := NewProposalManager()
	prop := pm.CreateProposal(dao, "c", "desc")
	if prop.DAOID != dao.ID {
		t.Fatalf("dao id mismatch")
	}
	if err := pm.Vote(prop.ID, "v1", 4, true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := pm.Vote(prop.ID, "v2", 1, false); err != nil {
		t.Fatalf("vote: %v", err)
	}
	yes, no, err := pm.Results(prop.ID)
	if err != nil || yes != 4 || no != 1 {
		t.Fatalf("unexpected tally: %d %d %v", yes, no, err)
	}
	if err := pm.Execute(prop.ID); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !prop.Executed {
		t.Fatalf("proposal not marked executed")
	}
	if err := pm.Vote(prop.ID, "v3", 1, true); err == nil {
		t.Fatalf("expected vote on executed proposal to fail")
	}
	if err := pm.Execute(prop.ID); err == nil {
		t.Fatalf("expected second execute to fail")
	}
}
