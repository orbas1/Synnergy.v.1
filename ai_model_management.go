package synnergy

import (
	"errors"
	"fmt"
	"sync"
)

// ModelListing represents a listing of an AI model.
type ModelListing struct {
	ID        string
	ModelHash string
	CID       string
	Seller    string
	Price     uint64
	Active    bool
}

// ModelMarketplace manages AI model listings.
type ModelMarketplace struct {
	mu       sync.RWMutex
	listings map[string]ModelListing
	nextID   uint64
}

// NewModelMarketplace constructs an empty marketplace.
func NewModelMarketplace() *ModelMarketplace {
	return &ModelMarketplace{listings: make(map[string]ModelListing)}
}

// AddListing registers a new model listing and returns its ID.
func (m *ModelMarketplace) AddListing(hash, cid, seller string, price uint64) string {
	m.mu.Lock()
	m.nextID++
	id := fmt.Sprintf("listing-%d", m.nextID)
	m.listings[id] = ModelListing{
		ID:        id,
		ModelHash: hash,
		CID:       cid,
		Seller:    seller,
		Price:     price,
		Active:    true,
	}
	m.mu.Unlock()
	return id
}

// Get retrieves a listing by ID.
func (m *ModelMarketplace) Get(id string) (ModelListing, bool) {
	m.mu.RLock()
	listing, ok := m.listings[id]
	m.mu.RUnlock()
	return listing, ok
}

// List returns all active listings.
func (m *ModelMarketplace) List() []ModelListing {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]ModelListing, 0, len(m.listings))
	for _, l := range m.listings {
		if l.Active {
			out = append(out, l)
		}
	}
	return out
}

// Update modifies the price of an existing listing.
func (m *ModelMarketplace) Update(id string, price uint64) error {
	m.mu.Lock()
	listing, ok := m.listings[id]
	if !ok {
		m.mu.Unlock()
		return errors.New("listing not found")
	}
	listing.Price = price
	m.listings[id] = listing
	m.mu.Unlock()
	return nil
}

// Remove deletes a listing owned by the seller.
func (m *ModelMarketplace) Remove(id, seller string) error {
	m.mu.Lock()
	listing, ok := m.listings[id]
	if !ok {
		m.mu.Unlock()
		return errors.New("listing not found")
	}
	if listing.Seller != seller {
		m.mu.Unlock()
		return errors.New("not the seller")
	}
	delete(m.listings, id)
	m.mu.Unlock()
	return nil
}
