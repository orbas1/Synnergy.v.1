package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ErrSYN800AssetExists indicates an asset ID already exists.
var ErrSYN800AssetExists = errors.New("syn800: asset exists")

// ErrSYN800AssetNotFound indicates an asset cannot be located.
var ErrSYN800AssetNotFound = errors.New("syn800: asset not found")

// AssetMetadata describes a real world asset backing a SYN800 token.
type AssetMetadata struct {
	ID            string
	Description   string
	Valuation     uint64
	Location      string
	AssetType     string
	Certification string
	Updated       time.Time
	Custodian     string
}

// AssetRegistry stores SYN800 assets.
type AssetRegistry struct {
	mu      sync.RWMutex
	assets  map[string]*AssetMetadata
	history []AssetMetadata
}

// NewAssetRegistry creates an empty registry.
func NewAssetRegistry() *AssetRegistry {
	return &AssetRegistry{assets: make(map[string]*AssetMetadata)}
}

// Register adds a new asset and returns it.
func (r *AssetRegistry) Register(id, desc string, valuation uint64, loc, typ, cert string) (*AssetMetadata, error) {
	if id == "" || valuation == 0 {
		return nil, fmt.Errorf("syn800: invalid asset parameters")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.assets[id]; exists {
		return nil, ErrSYN800AssetExists
	}
	now := time.Now().UTC()
	a := &AssetMetadata{ID: id, Description: desc, Valuation: valuation, Location: loc, AssetType: typ, Certification: cert, Updated: now}
	r.assets[id] = a
	r.history = append(r.history, *a)
	return a, nil
}

// AssignCustodian updates the custodian responsible for the asset.
func (r *AssetRegistry) AssignCustodian(id, custodian string) error {
	if custodian == "" {
		return fmt.Errorf("syn800: custodian required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	asset, ok := r.assets[id]
	if !ok {
		return ErrSYN800AssetNotFound
	}
	asset.Custodian = custodian
	asset.Updated = time.Now().UTC()
	r.history = append(r.history, *asset)
	return nil
}

// UpdateValuation updates the valuation of an asset.
func (r *AssetRegistry) UpdateValuation(id string, valuation uint64) error {
	if valuation == 0 {
		return fmt.Errorf("syn800: valuation must be positive")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.assets[id]
	if !ok {
		return ErrSYN800AssetNotFound
	}
	a.Valuation = valuation
	a.Updated = time.Now().UTC()
	r.history = append(r.history, *a)
	return nil
}

// Get returns asset metadata by id.
func (r *AssetRegistry) Get(id string) (*AssetMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.assets[id]
	if !ok {
		return nil, false
	}
	clone := *a
	return &clone, true
}

// Snapshot returns all assets sorted by valuation.
func (r *AssetRegistry) Snapshot() []AssetMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]AssetMetadata, 0, len(r.assets))
	for _, a := range r.assets {
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Valuation > out[j].Valuation })
	return out
}

// History returns the chronological history of updates limited to the most
// recent n entries when n > 0.
func (r *AssetRegistry) History(n int) []AssetMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if n <= 0 || n >= len(r.history) {
		out := make([]AssetMetadata, len(r.history))
		copy(out, r.history)
		return out
	}
	start := len(r.history) - n
	out := make([]AssetMetadata, len(r.history[start:]))
	copy(out, r.history[start:])
	return out
}
