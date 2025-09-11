package core

import (
	"context"
	"fmt"
	"testing"

	synn "synnergy"
)

func TestNFTMarketplace(t *testing.T) {
	m := NewNFTMarketplace()
	if _, err := m.Mint(context.Background(), "id1", "alice", "meta", 100, synn.GasCost("MintNFT")); err != nil {
		t.Fatalf("mint: %v", err)
	}
	nft, err := m.List("id1")
	if err != nil || nft.Owner != "alice" {
		t.Fatalf("list after mint: %v", err)
	}
	if err := m.Buy(context.Background(), "id1", "bob", synn.GasCost("BuyNFT")); err != nil {
		t.Fatalf("buy: %v", err)
	}
	nft, err = m.List("id1")
	if err != nil {
		t.Fatalf("list after buy: %v", err)
	}
	if nft.Owner != "bob" {
		t.Fatalf("expected bob, got %s", nft.Owner)
	}
	if len(m.ListAll()) != 1 {
		t.Fatalf("expected 1 nft")
	}
}

func TestNFTMarketplaceDuplicate(t *testing.T) {
	m := NewNFTMarketplace()
	if _, err := m.Mint(context.Background(), "id1", "alice", "m", 1, synn.GasCost("MintNFT")); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if _, err := m.Mint(context.Background(), "id1", "alice", "m", 1, synn.GasCost("MintNFT")); err == nil {
		t.Fatalf("expected duplicate error")
	}
}

func BenchmarkNFTMarketplaceMint(b *testing.B) {
	m := NewNFTMarketplace()
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("id%d", i)
		_, _ = m.Mint(context.Background(), id, "alice", "meta", 1, synn.GasCost("MintNFT"))
	}
}

func TestNFTMarketplaceUpdatePrice(t *testing.T) {
	m := NewNFTMarketplace()
	if _, err := m.Mint(context.Background(), "id1", "alice", "meta", 100, synn.GasCost("MintNFT")); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if err := m.UpdatePrice("id1", 200); err != nil {
		t.Fatalf("update: %v", err)
	}
	nft, err := m.List("id1")
	if err != nil || nft.Price != 200 {
		t.Fatalf("price not updated: %+v %v", nft, err)
	}
}
