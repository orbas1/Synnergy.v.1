package cli

import (
	"encoding/json"
	"testing"
)

// TestStateBalanceZero verifies querying balance returns zero when unset.
func TestStateBalanceZero(t *testing.T) {
	out, err := execCommand("--json", "state", "balance", "addr1")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var resp struct {
		Balance uint64 `json:"balance"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Balance != 0 {
		t.Fatalf("expected balance 0, got %d", resp.Balance)
	}
}
