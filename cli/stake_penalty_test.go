package cli

import (
	"encoding/json"
	"testing"

	"synnergy/core"
)

// TestStakePenaltySlashJSON applies a slash and confirms JSON output and balance.
func TestStakePenaltySlashJSON(t *testing.T) {
	stakingNode = core.NewStakingNode()
	if _, err := execCommand("--json", "staking_node", "stake", "alice", "10"); err != nil {
		t.Fatalf("stake: %v", err)
	}
	out, err := execCommand("--json", "stake_penalty", "slash", "alice", "5")
	if err != nil {
		t.Fatalf("slash: %v", err)
	}
	balOut, err := execCommand("--json", "staking_node", "balance", "alice")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var resp map[string]string
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal slash: %v", err)
	}
	if resp["status"] != "slashed" {
		t.Fatalf("expected slashed, got %s", resp["status"])
	}
	var bal uint64
	if err := json.Unmarshal([]byte(balOut), &bal); err != nil {
		t.Fatalf("unmarshal balance: %v", err)
	}
	if bal != 5 {
		t.Fatalf("expected balance 5, got %d", bal)
	}
}
