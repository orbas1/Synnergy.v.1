package cli

import (
	"encoding/json"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn1000ValueZeroJSON initialises the token and ensures the reserve value reports zero via JSON.
func TestSyn1000ValueZeroJSON(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	syn1000 = nil

	if _, err := execCommand("--json", "syn1000", "init", "--name", "SYN1000", "--symbol", "S1000"); err != nil {
		t.Fatalf("init: %v", err)
	}
	out, err := execCommand("--json", "syn1000", "value")
	if err != nil {
		t.Fatalf("value: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Value != "0.00" {
		t.Fatalf("expected value 0.00, got %s", resp.Value)
	}
}
