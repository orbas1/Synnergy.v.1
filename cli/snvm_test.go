package cli

import (
	"encoding/json"
	"testing"

	"synnergy/core"
)

// TestSnvmExecJSON runs a simple program and checks JSON output.
func TestSnvmExecJSON(t *testing.T) {
	vm = core.NewSNVM()
	out, err := execCommand("--json", "snvm", "exec", "add", "2", "3")
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var resp struct {
		Result int64 `json:"result"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Result != 5 {
		t.Fatalf("expected 5, got %d", resp.Result)
	}
}
