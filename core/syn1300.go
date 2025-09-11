package core

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrAssetExists is returned when attempting to register an asset twice.
	ErrAssetExists = errors.New("asset already exists")
	// ErrAssetNotFound is returned when an unknown asset id is referenced.
	ErrAssetNotFound = errors.New("asset not found")
)

// SupplyChainEvent records a movement or status change for a supply chain asset.
type SupplyChainEvent struct {
	Timestamp time.Time
	Location  string
	Status    string
	Note      string
}

// SupplyChainAsset holds metadata and history for an asset tracked through the supply chain.
type SupplyChainAsset struct {
	ID          string
	Description string
	Owner       string
	Location    string
	Status      string
	History     []SupplyChainEvent
}

// SupplyChainRegistry manages supply chain assets and their events.
type SupplyChainRegistry struct {
	mu     sync.RWMutex
	assets map[string]*SupplyChainAsset
}

// NewSupplyChainRegistry creates a new empty registry.
func NewSupplyChainRegistry() *SupplyChainRegistry {
	return &SupplyChainRegistry{assets: make(map[string]*SupplyChainAsset)}
}

// Register inserts a new asset into the registry.
func (r *SupplyChainRegistry) Register(id, desc, owner, location string) (*SupplyChainAsset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.assets[id]; exists {
		return nil, ErrAssetExists
	}
	asset := &SupplyChainAsset{ID: id, Description: desc, Owner: owner, Location: location, Status: "created"}
	asset.History = append(asset.History, SupplyChainEvent{Timestamp: time.Now(), Location: location, Status: "created"})
	r.assets[id] = asset
	return asset, nil
}

// Update records a new event for the given asset.
func (r *SupplyChainRegistry) Update(id, location, status, note string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	asset, ok := r.assets[id]
	if !ok {
		return ErrAssetNotFound
	}
	asset.Location = location
	asset.Status = status
	asset.History = append(asset.History, SupplyChainEvent{Timestamp: time.Now(), Location: location, Status: status, Note: note})
	return nil
}

// Get returns the asset with the given id.
func (r *SupplyChainRegistry) Get(id string) (*SupplyChainAsset, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	asset, ok := r.assets[id]
	return asset, ok
}
