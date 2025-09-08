package cli

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"synnergy/core"
)

// TestBankInstitutionalRegister ensures institutions can enrol via the CLI.
func TestBankInstitutionalRegister(t *testing.T) {
	ledger = core.NewLedger()
	pub, priv, _ := ed25519.GenerateKey(nil)
	addr := hex.EncodeToString(pub)
	bankInstNode = core.NewBankInstitutionalNode("bank1", addr, ledger)

	sig := ed25519.Sign(priv, []byte("register:BankA"))
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"bankinst", "register", "BankA", "--pub", addr, "--sig", hex.EncodeToString(sig)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if !bankInstNode.IsRegistered("BankA") {
		t.Fatalf("bank not registered")
	}

	sigRem := ed25519.Sign(priv, []byte("remove:BankA"))
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"bankinst", "remove", "BankA", "--pub", addr, "--sig", hex.EncodeToString(sigRem)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if bankInstNode.IsRegistered("BankA") {
		t.Fatalf("bank not removed")
	}
}
