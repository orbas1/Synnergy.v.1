package cli

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	cmd := RootCmd()
	if f := cmd.PersistentFlags().Lookup("config"); f == nil {
		t.Fatal("config flag not registered")
	}
	if f := cmd.PersistentFlags().Lookup("log-level"); f == nil {
		t.Fatal("log-level flag not registered")
	}
	cmd.SetArgs([]string{"--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
