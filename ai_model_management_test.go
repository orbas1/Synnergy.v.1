package synnergy

import (
	"sync"
	"testing"
)

// TestModelMarketplaceAddAndGet verifies that listings can be added and retrieved.
func TestModelMarketplaceAddAndGet(t *testing.T) {
	m := NewModelMarketplace()

	id := m.AddListing("hash1", "cid1", "alice", 100)
	if id == "" {
		t.Fatalf("expected non-empty listing ID")
	}

	listing, ok := m.Get(id)
	if !ok {
		t.Fatalf("expected to retrieve listing")
	}
	if listing.ID != id || listing.ModelHash != "hash1" || listing.CID != "cid1" || listing.Seller != "alice" || listing.Price != 100 || !listing.Active {
		t.Fatalf("retrieved listing mismatch: %+v", listing)
	}

	// Fetching non-existent listing should return false
	if _, ok := m.Get("missing"); ok {
		t.Fatalf("expected missing listing to return false")
	}
}

// TestModelMarketplaceList ensures that List returns only active listings.
func TestModelMarketplaceList(t *testing.T) {
	m := NewModelMarketplace()
	id1 := m.AddListing("hash1", "cid1", "alice", 100)
	id2 := m.AddListing("hash2", "cid2", "bob", 200)

	// Mark one listing inactive directly to test filtering
	m.mu.Lock()
	l := m.listings[id2]
	l.Active = false
	m.listings[id2] = l
	m.mu.Unlock()

	listings := m.List()
	if len(listings) != 1 {
		t.Fatalf("expected 1 active listing, got %d", len(listings))
	}
	if listings[0].ID != id1 {
		t.Fatalf("unexpected listing returned: %+v", listings[0])
	}
}

// TestModelMarketplaceUpdate checks price updates and error handling.
func TestModelMarketplaceUpdate(t *testing.T) {
	m := NewModelMarketplace()
	id := m.AddListing("hash1", "cid1", "alice", 100)

	if err := m.Update(id, 150); err != nil {
		t.Fatalf("update returned error: %v", err)
	}
	listing, _ := m.Get(id)
	if listing.Price != 150 {
		t.Fatalf("expected updated price 150, got %d", listing.Price)
	}

	if err := m.Update("missing", 200); err == nil {
		t.Fatalf("expected error for missing listing")
	}
}

// TestModelMarketplaceRemove verifies removal semantics.
func TestModelMarketplaceRemove(t *testing.T) {
	m := NewModelMarketplace()
	id := m.AddListing("hash1", "cid1", "alice", 100)

	if err := m.Remove(id, "bob"); err == nil {
		t.Fatalf("expected error when non-seller attempts removal")
	}

	if err := m.Remove("missing", "alice"); err == nil {
		t.Fatalf("expected error when removing missing listing")
	}

	if err := m.Remove(id, "alice"); err != nil {
		t.Fatalf("failed to remove listing: %v", err)
	}
	if _, ok := m.Get(id); ok {
		t.Fatalf("listing should be removed")
	}
}

// TestModelMarketplaceConcurrency adds listings concurrently to ensure thread safety.
func TestModelMarketplaceConcurrency(t *testing.T) {
	m := NewModelMarketplace()
	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			m.AddListing("hash", "cid", "seller", uint64(i))
		}(i)
	}
	wg.Wait()

	if len(m.List()) != n {
		t.Fatalf("expected %d listings, got %d", n, len(m.List()))
	}
}
