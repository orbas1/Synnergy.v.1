package cli

import (
	"encoding/json"
	"testing"
)

// TestStorageMarketplaceListEmpty ensures listings start empty and return JSON.
func TestStorageMarketplaceListEmpty(t *testing.T) {
	out, err := execCommand("--json", "storage_marketplace", "listings")
	if err != nil {
		t.Fatalf("listings: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var listings []any
	if err := json.Unmarshal([]byte(out), &listings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(listings) != 0 {
		t.Fatalf("expected no listings, got %d", len(listings))
	}
}
