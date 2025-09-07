package cli

import (
	"testing"

	"synnergy/core"
)

// TestWalletNew verifies core wallet creation.
func TestWalletNew(t *testing.T) {
	w, err := core.NewWallet()
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if len(w.Address) != 40 {
		t.Fatalf("bad address: %s", w.Address)
	}
}
