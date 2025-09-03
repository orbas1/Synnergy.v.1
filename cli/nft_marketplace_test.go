package cli

import (
	"bytes"
	"context"
	"testing"
)

// TestNFTMarketplaceCLI exercises basic minting and buying through the CLI.
func TestNFTMarketplaceCLI(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"nft", "mint", "id1", "alice", "meta", "100"})
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("mint: %v", err)
	}
	buf.Reset()

	rootCmd.SetArgs([]string{"nft", "list", "id1"})
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("list: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatalf("expected output")
	}
	buf.Reset()

	rootCmd.SetArgs([]string{"nft", "buy", "id1", "bob"})
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("buy: %v", err)
	}
}
