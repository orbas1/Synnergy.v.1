package core

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrIPExists indicates the IP asset already exists in the registry.
var ErrIPExists = errors.New("syn700: token already exists")

// ErrLicenseExists indicates the license identifier already exists.
var ErrLicenseExists = errors.New("syn700: license exists")

// ErrLicenseNotFound indicates the license identifier cannot be located.
var ErrLicenseNotFound = errors.New("syn700: license not found")

// License represents a usage license for an IP token.
type License struct {
	ID       string
	Type     string
	Licensee string
	Royalty  uint64
	Created  time.Time
	Expires  time.Time
}

// RoyaltyPayment records a royalty payment for a license.
type RoyaltyPayment struct {
	LicenseID string
	Licensee  string
	Amount    uint64
	Timestamp time.Time
}

// IPTokens represents an intellectual property asset.
type IPTokens struct {
	TokenID     string
	Title       string
	Description string
	Creator     string
	Owner       string
	Licenses    map[string]*License
	Royalties   []RoyaltyPayment
	Updated     time.Time
}

// IPRegistry manages SYN700 tokens with concurrent access safety.
type IPRegistry struct {
	mu     sync.RWMutex
	assets map[string]*IPTokens
}

// NewIPRegistry creates a new IP registry.
func NewIPRegistry() *IPRegistry {
	return &IPRegistry{assets: make(map[string]*IPTokens)}
}

// Register adds a new IP token.
func (r *IPRegistry) Register(tokenID, title, desc, creator, owner string) (*IPTokens, error) {
	if tokenID == "" || title == "" || creator == "" || owner == "" {
		return nil, fmt.Errorf("syn700: missing required fields")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.assets[tokenID]; exists {
		return nil, ErrIPExists
	}
	now := time.Now().UTC()
	asset := &IPTokens{TokenID: tokenID, Title: title, Description: desc, Creator: creator, Owner: owner, Licenses: make(map[string]*License), Updated: now}
	r.assets[tokenID] = asset
	return asset, nil
}

// Transfer changes ownership of an IP token.
func (r *IPRegistry) Transfer(tokenID, newOwner string) error {
	if newOwner == "" {
		return fmt.Errorf("syn700: new owner required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	asset, ok := r.assets[tokenID]
	if !ok {
		return errors.New("syn700: token not found")
	}
	asset.Owner = newOwner
	asset.Updated = time.Now().UTC()
	return nil
}

// CreateLicense associates a license with an IP token.
func (r *IPRegistry) CreateLicense(tokenID, licID, licType, licensee string, royalty uint64, expires time.Time) error {
	if licID == "" || licType == "" || licensee == "" || royalty == 0 {
		return fmt.Errorf("syn700: invalid license parameters")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	asset, ok := r.assets[tokenID]
	if !ok {
		return errors.New("syn700: token not found")
	}
	if _, exists := asset.Licenses[licID]; exists {
		return ErrLicenseExists
	}
	asset.Licenses[licID] = &License{ID: licID, Type: licType, Licensee: licensee, Royalty: royalty, Created: time.Now().UTC(), Expires: expires}
	asset.Updated = time.Now().UTC()
	return nil
}

// RecordRoyalty records a royalty payment for a license.
func (r *IPRegistry) RecordRoyalty(tokenID, licID, licensee string, amount uint64) error {
	if amount == 0 {
		return fmt.Errorf("syn700: royalty amount must be positive")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	asset, ok := r.assets[tokenID]
	if !ok {
		return errors.New("syn700: token not found")
	}
	license, exists := asset.Licenses[licID]
	if !exists {
		return ErrLicenseNotFound
	}
	if license.Expires.Before(time.Now().UTC()) {
		return fmt.Errorf("syn700: license expired")
	}
	asset.Royalties = append(asset.Royalties, RoyaltyPayment{LicenseID: licID, Licensee: licensee, Amount: amount, Timestamp: time.Now().UTC()})
	asset.Updated = time.Now().UTC()
	return nil
}

// RoyaltySummary calculates the total royalties for a license.
func (r *IPRegistry) RoyaltySummary(tokenID, licID string) (uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	asset, ok := r.assets[tokenID]
	if !ok {
		return 0, errors.New("syn700: token not found")
	}
	var total uint64
	for _, rpy := range asset.Royalties {
		if rpy.LicenseID == licID {
			total += rpy.Amount
		}
	}
	if _, exists := asset.Licenses[licID]; !exists {
		return 0, ErrLicenseNotFound
	}
	return total, nil
}

// Get returns an IP token copy for read-only operations.
func (r *IPRegistry) Get(tokenID string) (*IPTokens, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	asset, ok := r.assets[tokenID]
	if !ok {
		return nil, false
	}
	clone := *asset
	clone.Licenses = make(map[string]*License, len(asset.Licenses))
	for k, v := range asset.Licenses {
		lic := *v
		clone.Licenses[k] = &lic
	}
	clone.Royalties = append([]RoyaltyPayment(nil), asset.Royalties...)
	return &clone, true
}

// List returns all registered token identifiers.
func (r *IPRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.assets))
	for id := range r.assets {
		ids = append(ids, id)
	}
	return ids
}
