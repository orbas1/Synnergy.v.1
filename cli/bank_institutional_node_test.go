package cli

import (
	"bytes"
	"testing"

	"synnergy/core"
)

// TestBankInstitutionalRegister ensures institutions can enrol via the CLI.
func TestBankInstitutionalRegister(t *testing.T) {
	ledger = core.NewLedger()
	bankInstNode = core.NewBankInstitutionalNode("bank1", "addr1", ledger)

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"bankinst", "register", "BankA"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if !bankInstNode.IsRegistered("BankA") {
		t.Fatalf("bank not registered")
	}
}
