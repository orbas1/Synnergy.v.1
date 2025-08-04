package core

import "testing"

func TestLoanPoolProposalLifecycle(t *testing.T) {
	lp := NewLoanPool(1000)
	id, err := lp.SubmitProposal("creator", "rec", "type", 100, "desc")
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if err := lp.VoteProposal("voter1", id); err != nil {
		t.Fatalf("vote: %v", err)
	}
	lp.Tick()
	if err := lp.Disburse(id); err != nil {
		t.Fatalf("disburse: %v", err)
	}
	if _, ok := lp.GetProposal(id); !ok {
		t.Fatalf("proposal not found")
	}
	mgr := NewLoanPoolManager(lp)
	stats := mgr.Stats()
	if stats.DisbursedCount != 1 {
		t.Fatalf("expected disbursed count 1")
	}
	appMgr := NewLoanPoolApply(lp)
	appID := appMgr.Submit("alice", 50, 12, "test")
	if err := appMgr.Vote("bob", appID); err != nil {
		t.Fatalf("app vote: %v", err)
	}
	appMgr.Process()
	if err := appMgr.Disburse(appID); err != nil {
		t.Fatalf("app disburse: %v", err)
	}
}
