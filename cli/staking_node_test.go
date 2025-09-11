package cli

import (
	"encoding/json"
	"testing"

	"synnergy/core"
)

// TestStakingNodeTotalZero ensures total staked tokens report as zero by default.
func TestStakingNodeTotalZero(t *testing.T) {
	orig := stakingNode
	stakingNode = core.NewStakingNode()
	t.Cleanup(func() { stakingNode = orig })

	out, err := execCommand("--json", "staking_node", "total")
	if err != nil {
		t.Fatalf("total: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var total uint64
	if err := json.Unmarshal([]byte(out), &total); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}
