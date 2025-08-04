package core

import "testing"

func TestDAOProposal(t *testing.T) {
	mgr := NewDAOManager()
	dao := mgr.Create("d", "c")
	pm := NewProposalManager()
	prop := pm.CreateProposal(dao, "c", "desc")
	if prop.DAOID != dao.ID {
		t.Fatalf("dao id mismatch")
	}
	if err := pm.Vote(prop.ID, "v1", 1, true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if prop.YesVotes["v1"] != 1 {
		t.Fatalf("vote not recorded")
	}
}
