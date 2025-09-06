package cli

import (
	"strings"
	"testing"
)

// TestGenesisShow verifies the show command outputs deterministic wallet addresses.
func TestGenesisShow(t *testing.T) {
	out, err := execCommand("genesis", "show")
	if err != nil {
		t.Fatalf("show failed: %v", err)
	}
	if !strings.Contains(out, genesisWallets.Genesis) {
		t.Fatalf("genesis address missing from output: %s", out)
	}
	if !strings.Contains(out, genesisWallets.AuthorityNodes) {
		t.Fatalf("authority nodes address missing from output: %s", out)
	}
}

// TestGenesisAllocate ensures allocations are displayed for the default wallets.
func TestGenesisAllocate(t *testing.T) {
	out, err := execCommand("genesis", "allocate", "1000")
	if err != nil {
		t.Fatalf("allocate failed: %v", err)
	}
	if !strings.Contains(out, genesisWallets.InternalDevelopment) {
		t.Fatalf("allocation missing internal development address: %s", out)
	}
}
