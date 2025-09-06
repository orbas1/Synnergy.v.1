package cli

import (
	"bytes"
	"testing"

	"synnergy/core"
)

// TestNATCLIMapAndUnmap verifies mapping and unmapping through the CLI commands.
func TestNATCLIMapAndUnmap(t *testing.T) {
	natMgr = core.NewNATManager()
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetErr(new(bytes.Buffer))

	rootCmd.SetArgs([]string{"nat", "map", "1234"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("map failed: %v", err)
	}
	if port, ok := natMgr.GetPort("self"); !ok || port != 1234 {
		t.Fatalf("expected port 1234, got %d", port)
	}

	rootCmd.SetArgs([]string{"nat", "unmap"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unmap failed: %v", err)
	}
	if _, ok := natMgr.GetPort("self"); ok {
		t.Fatalf("expected mapping removed")
	}
}

// TestNATCLIInvalidPort ensures invalid ports return errors.
func TestNATCLIInvalidPort(t *testing.T) {
	natMgr = core.NewNATManager()
	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetErr(new(bytes.Buffer))

	rootCmd.SetArgs([]string{"nat", "map", "bad"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error for non-numeric port")
	}

	rootCmd.SetArgs([]string{"nat", "map", "70000"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error for out-of-range port")
	}
}
