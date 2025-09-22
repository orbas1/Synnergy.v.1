package core

import (
	"sync"
	"testing"
	"time"
)

func TestIPRegistryLifecycle(t *testing.T) {
	registry := NewIPRegistry()
	asset, err := registry.Register("IP1", "Title", "Desc", "Creator", "Owner")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if asset.Owner != "Owner" {
		t.Fatalf("unexpected owner: %s", asset.Owner)
	}
	if err := registry.Transfer("IP1", "NewOwner"); err != nil {
		t.Fatalf("transfer failed: %v", err)
	}
	copy, ok := registry.Get("IP1")
	if !ok {
		t.Fatal("expected asset")
	}
	if copy.Owner != "NewOwner" {
		t.Fatalf("expected new owner, got %s", copy.Owner)
	}
	expires := time.Now().Add(time.Hour)
	if err := registry.CreateLicense("IP1", "LIC1", "exclusive", "Bob", 10, expires); err != nil {
		t.Fatalf("license create failed: %v", err)
	}
	if err := registry.RecordRoyalty("IP1", "LIC1", "Bob", 5); err != nil {
		t.Fatalf("record royalty failed: %v", err)
	}
	total, err := registry.RoyaltySummary("IP1", "LIC1")
	if err != nil {
		t.Fatalf("summary failed: %v", err)
	}
	if total != 5 {
		t.Fatalf("unexpected royalty total %d", total)
	}
	ids := registry.List()
	if len(ids) != 1 || ids[0] != "IP1" {
		t.Fatalf("unexpected list: %#v", ids)
	}
}

func TestIPRegistryValidation(t *testing.T) {
	registry := NewIPRegistry()
	if _, err := registry.Register("", "Title", "Desc", "Creator", "Owner"); err == nil {
		t.Fatal("expected validation error for missing id")
	}
	if _, err := registry.Register("IP1", "Title", "Desc", "Creator", "Owner"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if _, err := registry.Register("IP1", "Title", "Desc", "Creator", "Owner"); err != ErrIPExists {
		t.Fatalf("expected duplicate error, got %v", err)
	}
	if err := registry.CreateLicense("IP1", "LIC1", "type", "", 10, time.Now().Add(time.Hour)); err == nil {
		t.Fatal("expected invalid license params")
	}
	if err := registry.CreateLicense("missing", "LIC1", "type", "Bob", 10, time.Now().Add(time.Hour)); err == nil {
		t.Fatal("expected missing token error")
	}
	if err := registry.CreateLicense("IP1", "LIC1", "type", "Bob", 10, time.Now().Add(-time.Hour)); err != nil {
		t.Fatalf("create license failed: %v", err)
	}
	if err := registry.CreateLicense("IP1", "LIC1", "type", "Bob", 10, time.Now().Add(time.Hour)); err != ErrLicenseExists {
		t.Fatalf("expected duplicate license error, got %v", err)
	}
	if err := registry.RecordRoyalty("IP1", "missing", "Bob", 5); err != ErrLicenseNotFound {
		t.Fatalf("expected missing license error, got %v", err)
	}
	if err := registry.RecordRoyalty("IP1", "LIC1", "Bob", 0); err == nil {
		t.Fatal("expected invalid royalty amount")
	}
	if err := registry.RecordRoyalty("IP1", "LIC1", "Bob", 5); err == nil {
		t.Fatal("expected expired license error")
	}
}

func TestIPRegistryConcurrentRoyalties(t *testing.T) {
	registry := NewIPRegistry()
	if _, err := registry.Register("IP1", "Title", "Desc", "Creator", "Owner"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := registry.CreateLicense("IP1", "LIC1", "type", "Bob", 1, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("license create failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := registry.RecordRoyalty("IP1", "LIC1", "Bob", 1); err != nil {
				t.Errorf("record royalty: %v", err)
			}
		}()
	}
	wg.Wait()

	total, err := registry.RoyaltySummary("IP1", "LIC1")
	if err != nil {
		t.Fatalf("summary failed: %v", err)
	}
	if total != 20 {
		t.Fatalf("expected total 20, got %d", total)
	}
}
