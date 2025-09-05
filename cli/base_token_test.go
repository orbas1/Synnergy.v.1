package cli

import (
	"bytes"
	"testing"

	"synnergy/internal/tokens"
)

// TestBaseTokenMint ensures base token initialisation and minting work via the CLI.
func TestBaseTokenMint(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	baseToken = nil

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"basetoken", "init", "--name", "Base", "--symbol", "B", "--decimals", "18"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	rootCmd.SetOut(new(bytes.Buffer))
	rootCmd.SetArgs([]string{"basetoken", "mint", "alice", "100"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("mint failed: %v", err)
	}

	if baseToken.BalanceOf("alice") != 100 {
		t.Fatalf("expected balance 100, got %d", baseToken.BalanceOf("alice"))
	}
}
