package cli

import (
	"bytes"
	"testing"

	"synnergy/core"
)

// TestAuthorityNodeIndexAdd ensures nodes can be added and retrieved via the CLI.
func TestAuthorityNodeIndexAdd(t *testing.T) {
	authorityIndex = core.NewAuthorityNodeIndex()

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority_index", "add", "addr1", "validator"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	if _, ok := authorityIndex.Get("addr1"); !ok {
		t.Fatalf("node not indexed")
	}
}
