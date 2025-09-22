package core

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	synn "synnergy"
	"synnergy/internal/telemetry"
)

// StorageListing represents a single storage offer on the marketplace.
type StorageListing struct {
	ID       string
	DataHash string
	Price    uint64
	Owner    string
}

// StorageDeal represents an agreement between a buyer and a listing owner.
type StorageDeal struct {
	ID        string
	ListingID string
	Buyer     string
}

// StorageMarketplace exposes a minimal in‑memory marketplace for storage deals.
// It is concurrency safe and relies on gas priced opcodes to limit abuse. The
// marketplace accepts an optional virtual machine so contracts can be executed
// alongside marketplace operations in future extensions.
type StorageMarketplace struct {
	vm       VirtualMachine
	mu       sync.RWMutex
	listings map[string]StorageListing
	deals    map[string]StorageDeal
	counter  int
}

// NewStorageMarketplace creates a new marketplace. If a VM is supplied it will
// be started to ensure handlers are ready for execution.
func NewStorageMarketplace(vm VirtualMachine) *StorageMarketplace {
	if vm != nil && !vm.Status() {
		_ = vm.Start()
	}
	return &StorageMarketplace{
		vm:       vm,
		listings: make(map[string]StorageListing),
		deals:    make(map[string]StorageDeal),
	}
}

// CreateListing registers a new storage offer identified by a content hash. The
// caller must provide sufficient gas or the operation is rejected.
func (m *StorageMarketplace) CreateListing(ctx context.Context, hash string, price uint64, owner string, gasLimit uint64) (string, error) {
	ctx, span := telemetry.Tracer("core.storage_market").Start(ctx, "StorageMarketplace.CreateListing")
	defer span.End()

	required := synn.GasCost("CreateListing")
	if gasLimit < required {
		return "", fmt.Errorf("%w: need %d", ErrInsufficientGas, required)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter++
	id := fmt.Sprintf("L%d", m.counter)
	m.listings[id] = StorageListing{ID: id, DataHash: hash, Price: price, Owner: owner}
	return id, nil
}

// ListListings returns all current storage offers.
func (m *StorageMarketplace) ListListings(ctx context.Context) ([]StorageListing, error) {
	ctx, span := telemetry.Tracer("core.storage_market").Start(ctx, "StorageMarketplace.ListListings")
	defer span.End()
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]StorageListing, 0, len(m.listings))
	for _, l := range m.listings {
		out = append(out, l)
	}
	return out, nil
}

// OpenDeal creates a deal for a given listing and buyer.
func (m *StorageMarketplace) OpenDeal(ctx context.Context, listingID, buyer string, gasLimit uint64) (string, error) {
	ctx, span := telemetry.Tracer("core.storage_market").Start(ctx, "StorageMarketplace.OpenDeal")
	defer span.End()
	required := synn.GasCost("OpenDeal")
	if gasLimit < required {
		return "", fmt.Errorf("%w: need %d", ErrInsufficientGas, required)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.listings[listingID]; !ok {
		return "", fmt.Errorf("listing %s not found", listingID)
	}
	m.counter++
	id := fmt.Sprintf("D%d", m.counter)
	m.deals[id] = StorageDeal{ID: id, ListingID: listingID, Buyer: buyer}
	return id, nil
}

// CloseDeal removes an open deal.
func (m *StorageMarketplace) CloseDeal(ctx context.Context, dealID string) error {
	ctx, span := telemetry.Tracer("core.storage_market").Start(ctx, "StorageMarketplace.CloseDeal")
	defer span.End()
	required := synn.GasCost("CloseDeal")
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.deals[dealID]; !ok {
		return fmt.Errorf("deal %s not found", dealID)
	}
	delete(m.deals, dealID)
	_ = required // gas accounted at call site via RegisterGasCost
	return nil
}

// GetListing returns the listing with the given ID.
func (m *StorageMarketplace) GetListing(ctx context.Context, id string) (StorageListing, bool) {
	ctx, span := telemetry.Tracer("core.storage_market").Start(ctx, "StorageMarketplace.GetListing")
	defer span.End()
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.listings[id]
	return l, ok
}

// GetDeal returns the deal with the given ID.
func (m *StorageMarketplace) GetDeal(ctx context.Context, id string) (StorageDeal, bool) {
	ctx, span := telemetry.Tracer("core.storage_market").Start(ctx, "StorageMarketplace.GetDeal")
	defer span.End()
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.deals[id]
	return d, ok
}

// ListDeals returns all open deals.
func (m *StorageMarketplace) ListDeals(ctx context.Context) ([]StorageDeal, error) {
	ctx, span := telemetry.Tracer("core.storage_market").Start(ctx, "StorageMarketplace.ListDeals")
	defer span.End()
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]StorageDeal, 0, len(m.deals))
	for _, d := range m.deals {
		out = append(out, d)
	}
	return out, nil
}

// MarshalListings encodes listings to JSON – convenient for CLI/GUI responses.
func MarshalListings(listings []StorageListing) ([]byte, error) {
	return json.Marshal(listings)
}

// MarshalDeals encodes deals to JSON for client consumption.
func MarshalDeals(deals []StorageDeal) ([]byte, error) {
	return json.Marshal(deals)
}
