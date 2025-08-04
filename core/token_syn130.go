package core

import (
	"errors"
	"time"
)

// Syn130SaleRecord tracks sale price history for tangible assets.
type Syn130SaleRecord struct {
	Buyer     string
	Price     uint64
	Timestamp time.Time
}

// LeaseInfo captures lease details for an asset.
type LeaseInfo struct {
	Lessee  string
	Payment uint64
	Start   time.Time
	End     time.Time
	Active  bool
}

// TangibleAsset represents a tokenised real-world asset.
type TangibleAsset struct {
	ID        string
	Owner     string
	Metadata  string
	Valuation uint64
	Sales     []Syn130SaleRecord
	Lease     *LeaseInfo
}

// TangibleAssetRegistry manages tangible assets.
type TangibleAssetRegistry struct {
	assets map[string]*TangibleAsset
}

// NewTangibleAssetRegistry creates a new registry.
func NewTangibleAssetRegistry() *TangibleAssetRegistry {
	return &TangibleAssetRegistry{assets: make(map[string]*TangibleAsset)}
}

// Register adds a new tangible asset.
func (r *TangibleAssetRegistry) Register(id, owner, meta string, value uint64) (*TangibleAsset, error) {
	if _, exists := r.assets[id]; exists {
		return nil, errors.New("asset exists")
	}
	asset := &TangibleAsset{ID: id, Owner: owner, Metadata: meta, Valuation: value}
	r.assets[id] = asset
	return asset, nil
}

// UpdateValuation updates the valuation of an asset.
func (r *TangibleAssetRegistry) UpdateValuation(id string, val uint64) error {
	asset, ok := r.assets[id]
	if !ok {
		return errors.New("asset not found")
	}
	asset.Valuation = val
	return nil
}

// RecordSale records a sale for the asset.
func (r *TangibleAssetRegistry) RecordSale(id, buyer string, price uint64) error {
	asset, ok := r.assets[id]
	if !ok {
		return errors.New("asset not found")
	}
	asset.Owner = buyer
	asset.Sales = append(asset.Sales, Syn130SaleRecord{Buyer: buyer, Price: price, Timestamp: time.Now()})
	return nil
}

// StartLease begins a lease for the asset.
func (r *TangibleAssetRegistry) StartLease(id, lessee string, payment uint64, start, end time.Time) error {
	asset, ok := r.assets[id]
	if !ok {
		return errors.New("asset not found")
	}
	asset.Lease = &LeaseInfo{Lessee: lessee, Payment: payment, Start: start, End: end, Active: true}
	return nil
}

// EndLease terminates an active lease.
func (r *TangibleAssetRegistry) EndLease(id string) error {
	asset, ok := r.assets[id]
	if !ok {
		return errors.New("asset not found")
	}
	if asset.Lease != nil {
		asset.Lease.Active = false
	}
	return nil
}

// Get returns a tangible asset by id.
func (r *TangibleAssetRegistry) Get(id string) (*TangibleAsset, bool) {
	asset, ok := r.assets[id]
	return asset, ok
}
