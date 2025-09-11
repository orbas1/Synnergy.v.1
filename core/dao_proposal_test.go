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
	// Non-member cannot create proposal
	if _, err := pm.CreateProposal(dao, "x", "desc"); err == nil {
		t.Fatalf("expected non-member create failure")
	}
	prop, err := pm.CreateProposal(dao, "c", "desc")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if prop.DAOID != dao.ID {
		t.Fatalf("dao id mismatch")
	}
	// Non-member cannot vote
	if err := pm.Vote(dao, prop.ID, "v1", 4, true); err == nil {
		t.Fatalf("expected non-member vote failure")
	}
	// Add members for voting
	if err := dao.AddMember("v1", RoleMember); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if err := dao.AddMember("v2", RoleMember); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if err := pm.Vote(dao, prop.ID, "v1", 4, true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := pm.Vote(dao, prop.ID, "v2", 1, false); err != nil {
		t.Fatalf("vote: %v", err)
	}
	yes, no, err := pm.Results(prop.ID)
	if err != nil || yes != 4 || no != 1 {
		t.Fatalf("unexpected tally: %d %d %v", yes, no, err)
	}
	// Non-admin cannot execute
	if err := pm.Execute(dao, prop.ID, "v1"); err == nil {
		t.Fatalf("expected non-admin execute failure")
	}
	if err := pm.Execute(dao, prop.ID, "c"); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !prop.Executed {
		t.Fatalf("proposal not marked executed")
	}
	if err := pm.Vote(dao, prop.ID, "v1", 1, true); err == nil {
		t.Fatalf("expected vote on executed proposal to fail")
	}
	if err := pm.Execute(dao, prop.ID, "c"); err == nil {
		t.Fatalf("expected second execute to fail")
	}
}
