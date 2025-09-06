package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestLoanmgrPauseResume validates pause and resume commands with JSON output.
func TestLoanmgrPauseResume(t *testing.T) {
	if _, err := execCommand("loanmgr", "pause", "--json"); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	out, err := execCommand("loanmgr", "resume", "--json")
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var res string
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if res != "resumed" {
		t.Fatalf("expected resumed, got %s", res)
	}
}
