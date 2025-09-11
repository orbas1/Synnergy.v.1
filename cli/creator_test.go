package cli

import (
	"path/filepath"
	"testing"

	"synnergy/core"
)

func TestCreatorDisableDistribution(t *testing.T) {
	dir := t.TempDir()
	walletPath := filepath.Join(dir, "creator.wallet")
	if _, err := execCommand("wallet", "new", "--out", walletPath, "--password", "pass"); err != nil {
		t.Fatalf("wallet new failed: %v", err)
	}
	w, err := core.LoadWallet(walletPath, "pass")
	if err != nil {
		t.Fatalf("load wallet failed: %v", err)
	}
	genesisWallets.CreatorWallet = w.Address
	if _, err := execCommand("creator", "disable-distribution", "--wallet", walletPath, "--password", "pass"); err != nil {
		t.Fatalf("disable-distribution failed: %v", err)
	}
	if core.IsCreatorDistributionEnabled() {
		t.Fatalf("creator distribution should be disabled")
	}
	core.SetCreatorDistribution(true)
}
