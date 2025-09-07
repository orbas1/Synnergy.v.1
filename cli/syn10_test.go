package cli

import (
	"encoding/json"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn10InfoJSON initialises the token and retrieves info via JSON output.
func TestSyn10InfoJSON(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	syn10 = nil

	if _, err := execCommand("--json", "syn10", "init", "--name", "SYN10", "--symbol", "S10", "--issuer", "Gov"); err != nil {
		t.Fatalf("init: %v", err)
	}
	out, err := execCommand("--json", "syn10", "info")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var info tokens.SYN10Info
	if err := json.Unmarshal([]byte(out), &info); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if info.Symbol != "S10" {
		t.Fatalf("expected symbol S10, got %s", info.Symbol)
	}
}
