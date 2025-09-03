package core

import (
	"context"
	"testing"

	synn "synnergy"
)

// TestStorageMarketplaceLifecycle ensures listings and deals can be created and closed
// with proper gas accounting.
func TestStorageMarketplaceLifecycle(t *testing.T) {
	vm := NewSimpleVM()
	if err := vm.Start(); err != nil {
		t.Fatalf("vm start: %v", err)
	}
	m := NewStorageMarketplace(vm)

	listingGas := synn.GasCost("CreateListing")
	id, err := m.CreateListing(context.Background(), "hash", 10, "alice", listingGas)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if id == "" {
		t.Fatalf("expected listing id")
	}
	if _, ok := m.GetListing(context.Background(), id); !ok {
		t.Fatalf("listing not stored")
	}

	dealGas := synn.GasCost("OpenDeal")
	dealID, err := m.OpenDeal(context.Background(), id, "bob", dealGas)
	if err != nil {
		t.Fatalf("open deal: %v", err)
	}
	if dealID == "" {
		t.Fatalf("expected deal id")
	}
	if _, ok := m.GetDeal(context.Background(), dealID); !ok {
		t.Fatalf("deal not stored")
	}

	if err := m.CloseDeal(context.Background(), dealID); err != nil {
		t.Fatalf("close deal: %v", err)
	}
	if _, ok := m.GetDeal(context.Background(), dealID); ok {
		t.Fatalf("deal still exists")
	}
}
