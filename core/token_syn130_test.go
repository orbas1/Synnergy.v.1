package core

import (
	"sync"
	"testing"
	"time"
)

func TestTangibleAssetRegistryConcurrency(t *testing.T) {
	r := NewTangibleAssetRegistry()
	if _, err := r.Register("a1", "alice", "house", 100); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.UpdateValuation("a1", 200)
		}()
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.RecordSale("a1", "bob", 150)
		}()
	}
	start := time.Now()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = r.StartLease("a1", "lessee", 10, start, start.Add(time.Hour))
	}()
	wg.Wait()

	asset, ok := r.Get("a1")
	if !ok {
		t.Fatalf("asset not found")
	}
	if asset.Valuation != 200 {
		t.Fatalf("expected valuation 200, got %d", asset.Valuation)
	}
	if asset.Owner != "bob" {
		t.Fatalf("expected owner bob, got %s", asset.Owner)
	}
	if asset.Lease == nil || asset.Lease.Lessee != "lessee" {
		t.Fatalf("lease not set correctly")
	}
}
