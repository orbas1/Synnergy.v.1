package tokens

import (
	"fmt"
	"sync"
)

// SYN70Asset represents a single game asset tracked by the SYN70 token standard.
type SYN70Asset struct {
	ID           string
	Owner        string
	Name         string
	Game         string
	Attributes   map[string]string
	Achievements []string
}

// SYN70Token manages in-game assets.
type SYN70Token struct {
	*BaseToken
	mu     sync.RWMutex
	assets map[string]*SYN70Asset
}

// NewSYN70Token constructs a SYN70 token instance.
func NewSYN70Token(id TokenID, name, symbol string, decimals uint8) *SYN70Token {
	return &SYN70Token{
		BaseToken: NewBaseToken(id, name, symbol, decimals),
		assets:    make(map[string]*SYN70Asset),
	}
}

// RegisterAsset creates a new asset and assigns it to the owner.
func (t *SYN70Token) RegisterAsset(id, owner, name, game string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.assets[id]; exists {
		return fmt.Errorf("asset already exists")
	}
	t.assets[id] = &SYN70Asset{ID: id, Owner: owner, Name: name, Game: game, Attributes: map[string]string{}}
	return t.BaseToken.Mint(owner, 1)
}

// TransferAsset moves an asset to a new owner.
func (t *SYN70Token) TransferAsset(id, newOwner string) error {
	t.mu.Lock()
	asset, ok := t.assets[id]
	if !ok {
		t.mu.Unlock()
		return fmt.Errorf("asset not found")
	}
	owner := asset.Owner
	t.mu.Unlock()
	if err := t.BaseToken.Transfer(owner, newOwner, 1); err != nil {
		return err
	}
	t.mu.Lock()
	asset.Owner = newOwner
	t.mu.Unlock()
	return nil
}

// SetAttribute sets a custom attribute on an asset.
func (t *SYN70Token) SetAttribute(id, key, value string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	asset, ok := t.assets[id]
	if !ok {
		return fmt.Errorf("asset not found")
	}
	asset.Attributes[key] = value
	return nil
}

// AddAchievement records an achievement for an asset.
func (t *SYN70Token) AddAchievement(id, name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	asset, ok := t.assets[id]
	if !ok {
		return fmt.Errorf("asset not found")
	}
	asset.Achievements = append(asset.Achievements, name)
	return nil
}

// AssetInfo returns asset information if present.
func (t *SYN70Token) AssetInfo(id string) (SYN70Asset, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	asset, ok := t.assets[id]
	if !ok {
		return SYN70Asset{}, fmt.Errorf("asset not found")
	}
	cp := *asset
	cp.Attributes = make(map[string]string)
	for k, v := range asset.Attributes {
		cp.Attributes[k] = v
	}
	cp.Achievements = append([]string(nil), asset.Achievements...)
	return cp, nil
}

// ListAssets returns all registered assets.
func (t *SYN70Token) ListAssets() []SYN70Asset {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]SYN70Asset, 0, len(t.assets))
	for _, a := range t.assets {
		cp := *a
		cp.Attributes = make(map[string]string)
		for k, v := range a.Attributes {
			cp.Attributes[k] = v
		}
		cp.Achievements = append([]string(nil), a.Achievements...)
		out = append(out, cp)
	}
	return out
}
