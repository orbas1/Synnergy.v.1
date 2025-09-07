package cli

import (
	"encoding/json"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn1000IndexValueZeroJSON creates a token and verifies value reporting via JSON.
func TestSyn1000IndexValueZeroJSON(t *testing.T) {
	syn1000Index = tokens.NewSYN1000Index()

	if _, err := execCommand("--json", "syn1000index", "create", "Token", "TK"); err != nil {
		t.Fatalf("create: %v", err)
	}
	out, err := execCommand("--json", "syn1000index", "value", "1")
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
