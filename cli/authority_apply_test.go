package cli

import (
	"bytes"
	"testing"
	"time"

	"synnergy/core"
)

// TestAuthorityApplySubmit exercises the authority application workflow via the CLI
// to ensure candidates can be submitted and subsequently listed.
func TestAuthorityApplySubmit(t *testing.T) {
	// reset global state used by the CLI commands
	authorityRegistry = core.NewAuthorityNodeRegistry()
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

	// list via CLI to ensure command executes without error
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"authority_apply", "list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list failed: %v", err)
	}
}
