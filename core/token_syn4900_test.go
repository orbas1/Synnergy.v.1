package core

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAgriculturalRegistryConcurrency(t *testing.T) {
	r := NewAgriculturalRegistry()
	harvest := time.Now()
	expiry := harvest.Add(24 * time.Hour)
	if _, err := r.Register("f1", "grain", "alice", "farm", 10, harvest, expiry, "cert"); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			_ = r.Transfer("f1", id)
		}(fmt.Sprintf("owner%d", i))
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(status string) {
			defer wg.Done()
			_ = r.UpdateStatus("f1", status)
		}(fmt.Sprintf("status%d", i))
	}
	wg.Wait()

	asset, ok := r.Get("f1")
	if !ok {
		t.Fatalf("asset not found")
	}
	if asset.Owner == "alice" {
		t.Fatalf("owner should have changed")
	}
	if len(asset.History) == 0 {
		t.Fatalf("expected history entries")
	}
}
