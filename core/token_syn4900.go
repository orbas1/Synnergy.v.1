package core

import (
	"errors"
	"sync"
	"time"
)

// AgriEvent captures a notable event for an agricultural asset.
type AgriEvent struct {
	Timestamp   time.Time
	Description string
}

// AgriculturalAsset holds detailed metadata for SYN4900 tokens.
type AgriculturalAsset struct {
	ID            string
	AssetType     string
	Quantity      uint64
	Owner         string
	Origin        string
	HarvestDate   time.Time
	ExpiryDate    time.Time
	Status        string
	Certification string
	History       []AgriEvent
}

// AgriculturalRegistry manages agricultural assets.
type AgriculturalRegistry struct {
	mu     sync.RWMutex
	assets map[string]*AgriculturalAsset
}

// NewAgriculturalRegistry creates a new registry.
func NewAgriculturalRegistry() *AgriculturalRegistry {
	return &AgriculturalRegistry{assets: make(map[string]*AgriculturalAsset)}
}

// Register adds a new agricultural asset to the registry.
func (r *AgriculturalRegistry) Register(id, assetType, owner, origin string, qty uint64, harvest, expiry time.Time, cert string) (*AgriculturalAsset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.assets[id]; exists {
		return nil, errors.New("asset exists")
	}
	a := &AgriculturalAsset{ID: id, AssetType: assetType, Quantity: qty, Owner: owner, Origin: origin, HarvestDate: harvest, ExpiryDate: expiry, Status: "fresh", Certification: cert}
	r.assets[id] = a
	return a, nil
}

// Transfer moves ownership of an asset.
func (r *AgriculturalRegistry) Transfer(id, newOwner string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.assets[id]
	if !ok {
		return errors.New("asset not found")
	}
	a.Owner = newOwner
	a.History = append(a.History, AgriEvent{Timestamp: time.Now(), Description: "transfer"})
	return nil
}

// UpdateStatus updates the current status of an asset.
func (r *AgriculturalRegistry) UpdateStatus(id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.assets[id]
	if !ok {
		return errors.New("asset not found")
	}
	a.Status = status
	a.History = append(a.History, AgriEvent{Timestamp: time.Now(), Description: status})
	return nil
}

// Get returns an agricultural asset by id.
func (r *AgriculturalRegistry) Get(id string) (*AgriculturalAsset, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.assets[id]
	if !ok {
		return nil, false
	}
	cp := *a
	return &cp, true
}
