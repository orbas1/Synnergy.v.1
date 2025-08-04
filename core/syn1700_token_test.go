package core

import "testing"

func TestEventTickets(t *testing.T) {
	e := NewEvent("Concert", "Live show", "NYC", 1000, 2000, 2)
	id, err := e.IssueTicket("alice", "VIP", "standard", 100)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if !e.VerifyTicket(id, "alice") {
		t.Fatalf("expected alice to own ticket")
	}
	if err := e.TransferTicket(id, "alice", "bob"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if !e.VerifyTicket(id, "bob") {
		t.Fatalf("expected bob to own ticket")
	}
}
