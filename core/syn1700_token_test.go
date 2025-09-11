package core

import (
	"sync"
	"testing"
)

func TestEventTicketsBasic(t *testing.T) {
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

func TestEventTicketSupplyExhausted(t *testing.T) {
	e := NewEvent("Concert", "desc", "loc", 1, 2, 1)
	if _, err := e.IssueTicket("alice", "A", "std", 10); err != nil {
		t.Fatalf("first issue failed: %v", err)
	}
	if _, err := e.IssueTicket("bob", "A", "std", 10); err != ErrTicketSupplyExhausted {
		t.Fatalf("expected ErrTicketSupplyExhausted, got %v", err)
	}
}

func TestEventTicketUnauthorizedTransfer(t *testing.T) {
	e := NewEvent("Concert", "desc", "loc", 1, 2, 1)
	id, _ := e.IssueTicket("alice", "A", "std", 10)
	if err := e.TransferTicket(id, "bob", "carol"); err != ErrTicketNotOwned {
		t.Fatalf("expected ErrTicketNotOwned, got %v", err)
	}
}

func TestEventTicketConcurrentIssue(t *testing.T) {
	e := NewEvent("Concert", "desc", "loc", 1, 2, 10)
	var wg sync.WaitGroup
	successes := make(chan uint64, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(owner string) {
			defer wg.Done()
			if id, err := e.IssueTicket(owner, "A", "std", 1); err == nil {
				successes <- id
			}
		}(string('a' + rune(i)))
	}
	wg.Wait()
	close(successes)
	if got := len(successes); got != 10 {
		t.Fatalf("expected 10 successful issues, got %d", got)
	}
}
