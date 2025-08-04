package core

import (
	"errors"
	"time"
)

// License represents a usage license for an IP token.
type License struct {
	ID       string
	Type     string
	Licensee string
	Royalty  uint64
	Created  time.Time
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
}

// IPRegistry manages SYN700 tokens.
type IPRegistry struct {
	assets map[string]*IPTokens
}

// NewIPRegistry creates a new IP registry.
func NewIPRegistry() *IPRegistry {
	return &IPRegistry{assets: make(map[string]*IPTokens)}
}

// Register adds a new IP token.
func (r *IPRegistry) Register(tokenID, title, desc, creator, owner string) (*IPTokens, error) {
	if _, exists := r.assets[tokenID]; exists {
		return nil, errors.New("token already exists")
	}
	asset := &IPTokens{TokenID: tokenID, Title: title, Description: desc, Creator: creator, Owner: owner, Licenses: make(map[string]*License)}
	r.assets[tokenID] = asset
	return asset, nil
}

// CreateLicense associates a license with an IP token.
func (r *IPRegistry) CreateLicense(tokenID, licID, licType, licensee string, royalty uint64) error {
	asset, ok := r.assets[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	if _, exists := asset.Licenses[licID]; exists {
		return errors.New("license exists")
	}
	asset.Licenses[licID] = &License{ID: licID, Type: licType, Licensee: licensee, Royalty: royalty, Created: time.Now()}
	return nil
}

// RecordRoyalty records a royalty payment for a license.
func (r *IPRegistry) RecordRoyalty(tokenID, licID, licensee string, amount uint64) error {
	asset, ok := r.assets[tokenID]
	if !ok {
		return errors.New("token not found")
	}
	if _, exists := asset.Licenses[licID]; !exists {
		return errors.New("license not found")
	}
	asset.Royalties = append(asset.Royalties, RoyaltyPayment{LicenseID: licID, Licensee: licensee, Amount: amount, Timestamp: time.Now()})
	return nil
}

// Get returns an IP token.
func (r *IPRegistry) Get(tokenID string) (*IPTokens, bool) {
	asset, ok := r.assets[tokenID]
	return asset, ok
}
