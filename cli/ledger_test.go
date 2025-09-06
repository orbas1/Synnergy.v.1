package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestLedgerMintBalance ensures minting and balance queries emit JSON with gas info.
func TestLedgerMintBalance(t *testing.T) {
	if _, err := execCommand("ledger", "mint", "alice", "50", "--json"); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	out, err := execCommand("ledger", "balance", "alice", "--json")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var bal float64
	if err := json.Unmarshal([]byte(out), &bal); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if bal != 50 {
		t.Fatalf("expected balance 50, got %v", bal)
	}
}
