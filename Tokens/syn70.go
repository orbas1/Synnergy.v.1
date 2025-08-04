package tokens

import "fmt"

// SYN70Asset represents a single game asset tracked by the SYN70 token standard.
type SYN70Asset struct {
	ID       string
	Owner    string
	Metadata string
}

// SYN70Token manages in-game assets.
type SYN70Token struct {
	*BaseToken
	assets map[string]SYN70Asset
}

// NewSYN70Token constructs a SYN70 token instance.
func NewSYN70Token(id TokenID, name, symbol string, decimals uint8) *SYN70Token {
	return &SYN70Token{
		BaseToken: NewBaseToken(id, name, symbol, decimals),
		assets:    make(map[string]SYN70Asset),
	}
}

// MintAsset creates a new asset and assigns it to the owner.
func (t *SYN70Token) MintAsset(owner, assetID, metadata string) error {
	if _, exists := t.assets[assetID]; exists {
		return fmt.Errorf("asset already exists")
	}
	t.assets[assetID] = SYN70Asset{ID: assetID, Owner: owner, Metadata: metadata}
	return t.BaseToken.Mint(owner, 1)
}

// TransferAsset moves an asset from one owner to another.
func (t *SYN70Token) TransferAsset(assetID, from, to string) error {
	asset, ok := t.assets[assetID]
	if !ok {
		return fmt.Errorf("asset not found")
	}
	if asset.Owner != from {
		return fmt.Errorf("not asset owner")
	}
	if err := t.BaseToken.Transfer(from, to, 1); err != nil {
		return err
	}
	asset.Owner = to
	t.assets[assetID] = asset
	return nil
}

// GetAsset returns asset information if present.
func (t *SYN70Token) GetAsset(assetID string) (SYN70Asset, bool) {
	a, ok := t.assets[assetID]
	return a, ok
}
