package core

import "testing"

// TestSYN300ProposalLifecycle exercises the happy path for creating,
// voting on and executing a proposal.
func TestSYN300ProposalLifecycle(t *testing.T) {
	token := NewSYN300Token(map[string]uint64{"alice": 100, "bob": 50})
	if err := token.Delegate("bob", "alice"); err != nil {
		t.Fatalf("delegate: %v", err)
	}
	id, err := token.CreateProposal("alice", "upgrade network")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if err := token.Vote(id, "alice", true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := token.Execute(id, 100); err != nil {
		t.Fatalf("execute: %v", err)
	}
	p, err := token.ProposalStatus(id)
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !p.Executed {
		t.Fatalf("expected proposal executed")
	}
}

// TestSYN300Validations checks common validation failures.
func TestSYN300Validations(t *testing.T) {
	token := NewSYN300Token(map[string]uint64{"alice": 100})
	if err := token.Delegate("alice", "alice"); err == nil {
		t.Fatalf("expected self delegation error")
	}
	if _, err := token.CreateProposal("carol", "test"); err == nil {
		t.Fatalf("expected error for creator with no power")
	}
	if _, err := token.CreateProposal("alice", ""); err == nil {
		t.Fatalf("expected error for empty description")
	}
	id, err := token.CreateProposal("alice", "governance")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := token.Vote(id, "bob", true); err == nil {
		t.Fatalf("expected error for voter with no power")
	}
	if err := token.Vote(id, "alice", true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := token.Vote(id, "alice", true); err == nil {
		t.Fatalf("expected double vote error")
	}
}
