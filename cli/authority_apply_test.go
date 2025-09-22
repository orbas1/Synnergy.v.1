package cli

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"synnergy/core"
)

// TestAuthorityApplySubmit exercises the authority application workflow via the CLI
// to ensure candidates can be submitted and subsequently listed.
func TestAuthorityApplySubmit(t *testing.T) {
	// reset global state used by the CLI commands
	ledger = core.NewLedger()
	authorityValidators = core.NewValidatorManager(core.MinStake)
	authorityReg = core.NewAuthorityNodeRegistry(ledger, authorityValidators, 1)
	authorityRegistry = authorityReg
	applyManager = core.NewAuthorityApplicationManager(authorityRegistry, time.Hour)

	// submit an application
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority_apply", "submit", "node1", "validator", "desc"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("submit failed: %v", err)
	}

	// ensure application stored
	apps := applyManager.List()
	if len(apps) != 1 || apps[0].Candidate != "node1" {
		t.Fatalf("unexpected applications: %+v", apps)
	}

	// cast a signed vote approving the application
	pub, priv, _ := ed25519.GenerateKey(nil)
	voter := hex.EncodeToString(pub)
	msg := fmt.Sprintf("%s:%t", apps[0].ID, true)
	sig := ed25519.Sign(priv, []byte(msg))
	ledger.Mint(voter, 10)
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority_apply", "vote", voter, apps[0].ID, "true", "--pub", hex.EncodeToString(pub), "--sig", hex.EncodeToString(sig)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("vote failed: %v", err)
	}

	// finalize via CLI
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority_apply", "finalize", apps[0].ID})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("finalize failed: %v", err)
	}

	if !authorityRegistry.IsAuthorityNode("node1") {
		t.Fatalf("candidate not registered after finalize")
	}
}
