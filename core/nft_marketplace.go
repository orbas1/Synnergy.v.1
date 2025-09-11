package core

import (
	"context"
	"fmt"
	"sync"

	synn "synnergy"
	"synnergy/internal/telemetry"
)

// NFT represents a non-fungible token minted through the marketplace.
type NFT struct {
	ID       string
	Owner    string
	Metadata string
	Price    uint64
}

// NFTMarketplace provides minimal mint and trade capabilities for NFTs.
type NFTMarketplace struct {
	mu   sync.RWMutex
	nfts map[string]*NFT
}

// NewNFTMarketplace creates an empty marketplace.
func NewNFTMarketplace() *NFTMarketplace {
	return &NFTMarketplace{nfts: make(map[string]*NFT)}
}

// Mint creates a new NFT owned by the given address. A gas limit must be
// supplied and is validated against the registered cost.
func (m *NFTMarketplace) Mint(ctx context.Context, id, owner, metadata string, price, gasLimit uint64) (*NFT, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "NFTMarketplace.Mint")
	defer span.End()

	required := synn.GasCost("MintNFT")
	if gasLimit < required {
		return nil, fmt.Errorf("%w: need %d", ErrInsufficientGas, required)
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.nfts[id]; exists {
		return nil, fmt.Errorf("nft exists")
	}
	nft := &NFT{ID: id, Owner: owner, Metadata: metadata, Price: price}
	m.nfts[id] = nft
	return nft, nil
}

// List returns the NFT with the given identifier.
func (m *NFTMarketplace) List(id string) (*NFT, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	nft, ok := m.nfts[id]
	if !ok {
		return nil, fmt.Errorf("nft not found")
	}
	return nft, nil
}

// Buy transfers ownership of the NFT to the new owner. The buyer must supply a
// gas limit that meets the registered cost.
func (m *NFTMarketplace) Buy(ctx context.Context, id, newOwner string, gasLimit uint64) error {
	ctx, span := telemetry.Tracer().Start(ctx, "NFTMarketplace.Buy")
	defer span.End()

	required := synn.GasCost("BuyNFT")
	if gasLimit < required {
		return fmt.Errorf("%w: need %d", ErrInsufficientGas, required)
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	nft, ok := m.nfts[id]
	if !ok {
		return fmt.Errorf("nft not found")
	}
	nft.Owner = newOwner
	return nil
}

// UpdatePrice adjusts the listed price of an existing NFT.
func (m *NFTMarketplace) UpdatePrice(id string, price uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	nft, ok := m.nfts[id]
	if !ok {
		return fmt.Errorf("nft not found")
	}
	nft.Price = price
	return nil
}

// ListAll returns a snapshot of all NFTs.
func (m *NFTMarketplace) ListAll() []*NFT {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*NFT, 0, len(m.nfts))
	for _, n := range m.nfts {
		out = append(out, n)
	}
	return out
}
