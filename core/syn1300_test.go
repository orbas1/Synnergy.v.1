package core

import (
	"testing"
)

func TestSupplyChainRegistryRegisterGet(t *testing.T) {
	reg := NewSupplyChainRegistry()
	asset, err := reg.Register("asset1", "desc", "alice", "loc1")
	if err != nil {
		t.Fatalf("register error: %v", err)
	}
	if asset.ID != "asset1" || asset.Owner != "alice" || asset.Location != "loc1" || asset.Status != "created" {
		t.Fatalf("unexpected asset data: %#v", asset)
	}
	if len(asset.History) != 1 {
		t.Fatalf("expected history length 1 got %d", len(asset.History))
	}
	got, ok := reg.Get("asset1")
	if !ok || got.ID != "asset1" {
		t.Fatalf("failed to get asset")
	}
}

func TestSupplyChainRegistryDuplicate(t *testing.T) {
	reg := NewSupplyChainRegistry()
	if _, err := reg.Register("asset1", "d", "a", "loc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := reg.Register("asset1", "d", "a", "loc"); err == nil {
		t.Fatalf("expected error for duplicate asset")
	}
}

func TestSupplyChainRegistryUpdate(t *testing.T) {
	reg := NewSupplyChainRegistry()
	if _, err := reg.Register("asset1", "d", "a", "loc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := reg.Update("asset1", "loc2", "shipped", "note1"); err != nil {
		t.Fatalf("update error: %v", err)
	}
	asset, _ := reg.Get("asset1")
	if asset.Location != "loc2" || asset.Status != "shipped" {
		t.Fatalf("update failed: %#v", asset)
	}
	if len(asset.History) != 2 {
		t.Fatalf("expected history length 2 got %d", len(asset.History))
	}
	last := asset.History[len(asset.History)-1]
	if last.Status != "shipped" || last.Note != "note1" || last.Location != "loc2" {
		t.Fatalf("unexpected event: %#v", last)
	}
}

func TestSupplyChainRegistryUpdateNonexistent(t *testing.T) {
	reg := NewSupplyChainRegistry()
	if err := reg.Update("missing", "loc", "status", "note"); err == nil {
		t.Fatalf("expected error for missing asset")
	}
}
