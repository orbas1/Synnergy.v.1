package core

import (
	"testing"
)

func TestAssetRegistryLifecycle(t *testing.T) {
	registry := NewAssetRegistry()
	asset, err := registry.Register("A1", "desc", 100, "loc", "type", "cert")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if asset.Valuation != 100 {
		t.Fatalf("unexpected valuation: %d", asset.Valuation)
	}
	if err := registry.AssignCustodian("A1", "custodian"); err != nil {
		t.Fatalf("assign custodian: %v", err)
	}
	if err := registry.UpdateValuation("A1", 150); err != nil {
		t.Fatalf("update valuation: %v", err)
	}
	snapshot := registry.Snapshot()
	if len(snapshot) != 1 || snapshot[0].Valuation != 150 {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
	history := registry.History(2)
	if len(history) != 2 {
		t.Fatalf("unexpected history length: %d", len(history))
	}
	clone, ok := registry.Get("A1")
	if !ok || clone.Valuation != 150 {
		t.Fatalf("unexpected clone: %#v", clone)
	}
}

func TestAssetRegistryValidation(t *testing.T) {
	registry := NewAssetRegistry()
	if _, err := registry.Register("", "desc", 100, "loc", "type", "cert"); err == nil {
		t.Fatal("expected invalid parameters error")
	}
	if _, err := registry.Register("A1", "desc", 100, "loc", "type", "cert"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if _, err := registry.Register("A1", "desc", 100, "loc", "type", "cert"); err != ErrSYN800AssetExists {
		t.Fatalf("expected duplicate error, got %v", err)
	}
	if err := registry.AssignCustodian("A1", ""); err == nil {
		t.Fatal("expected custodian validation error")
	}
	if err := registry.AssignCustodian("missing", "custodian"); err != ErrSYN800AssetNotFound {
		t.Fatalf("expected missing asset error, got %v", err)
	}
	if err := registry.UpdateValuation("missing", 100); err != ErrSYN800AssetNotFound {
		t.Fatalf("expected missing asset error, got %v", err)
	}
	if err := registry.UpdateValuation("A1", 0); err == nil {
		t.Fatal("expected invalid valuation error")
	}
}
