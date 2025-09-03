package cli

import (
	"bytes"
	"context"
	"testing"
)

// TestStorageMarketplaceCLI exercises listing creation via the CLI commands.
func TestStorageMarketplaceCLI(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)

	// Create listing
	rootCmd.SetArgs([]string{"storage_marketplace", "list", "hash", "1", "alice"})
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatalf("expected output")
	}
	buf.Reset()

	// List listings
	rootCmd.SetArgs([]string{"storage_marketplace", "listings"})
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("listings: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatalf("expected listings json")
	}
}
