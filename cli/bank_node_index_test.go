package cli

import (
	"bytes"
	"testing"

	"synnergy/core"
)

// TestBankNodeIndexAdd ensures records can be added via the CLI.
func TestBankNodeIndexAdd(t *testing.T) {
	bankIndex = core.NewBankNodeIndex()

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"bank_index", "add", "id1", core.BankInstitutionalNodeType})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	if _, ok := bankIndex.Get("id1"); !ok {
		t.Fatalf("record not indexed")
	}
}
