package tokens

import (
	"testing"
	"time"
)

func TestSYN300TokenSupplyAndTransfer(t *testing.T) {
	token := NewSYN300Token(map[string]uint64{"alice": 100})
	if err := token.Mint("bob", 0); err != ErrInvalidAmount {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
	if err := token.Mint("bob", 50); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if token.TotalSupply() != 150 {
		t.Fatalf("unexpected supply: %d", token.TotalSupply())
	}
	if err := token.Transfer("alice", "bob", 40); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if token.BalanceOf("alice") != 60 {
		t.Fatalf("unexpected alice balance: %d", token.BalanceOf("alice"))
	}
	if err := token.Burn("bob", 10); err != nil {
		t.Fatalf("burn: %v", err)
	}
	if token.TotalSupply() != 140 {
		t.Fatalf("unexpected supply after burn: %d", token.TotalSupply())
	}
}

func TestSYN300GovernanceLifecycle(t *testing.T) {
	token := NewSYN300Token(map[string]uint64{"alice": 100})
	if _, err := token.CreateProposalWithDeadline("carol", "unauthorised", time.Now().Add(time.Hour)); err != ErrNoVotingPower {
		t.Fatalf("expected ErrNoVotingPower, got %v", err)
	}
	if _, err := token.CreateProposalWithDeadline("alice", "expired", time.Now().Add(-time.Hour)); err != ErrProposalExpired {
		t.Fatalf("expected ErrProposalExpired, got %v", err)
	}

	id, err := token.CreateProposal("alice", "upgrade network")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if err := token.Vote(id, "alice", true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := token.Execute(id, 50); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if err := token.Vote(id, "alice", true); err != ErrProposalExecuted {
		t.Fatalf("expected ErrProposalExecuted, got %v", err)
	}

	status, err := token.ProposalStatus(id)
	if err != nil {
		t.Fatalf("proposal status: %v", err)
	}
	if !status.Executed || status.ExecutedAt.IsZero() {
		t.Fatalf("expected executed proposal")
	}

	proposals := token.ListProposals()
	if len(proposals) != 1 || proposals[0].ID != id {
		t.Fatalf("unexpected proposals list: %+v", proposals)
	}
}
