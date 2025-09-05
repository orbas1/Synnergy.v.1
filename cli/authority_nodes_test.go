package cli

import (
	"bytes"
	"testing"

	"synnergy/core"
)

// TestAuthorityRegister verifies authority node registration through the CLI.
func TestAuthorityRegister(t *testing.T) {
	authorityReg = core.NewAuthorityNodeRegistry()

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority", "register", "addr1", "validator"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if !authorityReg.IsAuthorityNode("addr1") {
		t.Fatalf("address not registered")
	}
}
