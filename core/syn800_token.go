package core

import (
	"errors"
	"time"
)

// AssetMetadata describes a real world asset backing a SYN800 token.
type AssetMetadata struct {
	ID            string
	Description   string
	Valuation     uint64
	Location      string
	AssetType     string
	Certification string
	Updated       time.Time
}

// AssetRegistry stores SYN800 assets.
type AssetRegistry struct {
	assets map[string]*AssetMetadata
}

// NewAssetRegistry creates an empty registry.
func NewAssetRegistry() *AssetRegistry {
	return &AssetRegistry{assets: make(map[string]*AssetMetadata)}
}

// Register adds a new asset and returns it.
func (r *AssetRegistry) Register(id, desc string, valuation uint64, loc, typ, cert string) (*AssetMetadata, error) {
	if _, exists := r.assets[id]; exists {
		return nil, errors.New("asset exists")
	}
	a := &AssetMetadata{ID: id, Description: desc, Valuation: valuation, Location: loc, AssetType: typ, Certification: cert, Updated: time.Now()}
	r.assets[id] = a
	return a, nil
}

// UpdateValuation updates the valuation of an asset.
func (r *AssetRegistry) UpdateValuation(id string, valuation uint64) error {
	a, ok := r.assets[id]
	if !ok {
		return errors.New("asset not found")
	}
	a.Valuation = valuation
	a.Updated = time.Now()
	return nil
}

// Get returns asset metadata by id.
func (r *AssetRegistry) Get(id string) (*AssetMetadata, bool) {
	a, ok := r.assets[id]
	return a, ok
}
