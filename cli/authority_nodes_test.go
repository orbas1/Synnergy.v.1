package cli

import (
	"bytes"
	"encoding/hex"
	"testing"

	"crypto/ed25519"
	"synnergy/core"
)

// TestAuthorityRegister verifies authority node registration through the CLI.
func TestAuthorityRegister(t *testing.T) {
	ledger = core.NewLedger()
	authorityValidators = core.NewValidatorManager(core.MinStake)
	authorityReg = core.NewAuthorityNodeRegistry(ledger, authorityValidators, 1)

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority", "register", "addr1", "validator"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if !authorityReg.IsAuthorityNode("addr1") {
		t.Fatalf("address not registered")
	}
}

// TestAuthorityVote verifies signature-based voting through the CLI.
func TestAuthorityVote(t *testing.T) {
	ledger = core.NewLedger()
	authorityValidators = core.NewValidatorManager(core.MinStake)
	authorityReg = core.NewAuthorityNodeRegistry(ledger, authorityValidators, 1)
	// register candidate
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority", "register", "addr1", "validator"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	pub, priv, _ := ed25519.GenerateKey(nil)
	voter := hex.EncodeToString(pub)
	sig := ed25519.Sign(priv, []byte("addr1"))
	sigHex := hex.EncodeToString(sig)
	ledger.Mint(voter, 10)
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority", "vote", voter, "addr1", "--pub", voter, "--sig", sigHex})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("vote failed: %v", err)
	}

	info, _ := authorityReg.Info("addr1")
	if len(info.Votes) != 1 {
		t.Fatalf("expected 1 vote, got %d", len(info.Votes))
	}
}
