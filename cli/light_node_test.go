package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestLightNodeHeader verifies header addition and retrieval via JSON output.
func TestLightNodeHeader(t *testing.T) {
	if _, err := execCommand("light", "add-header", "h1", "1", "p0", "--json"); err != nil {
		t.Fatalf("add-header: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	out, err := execCommand("light", "latest", "--json")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var header map[string]any
	if err := json.Unmarshal([]byte(out), &header); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if header["Hash"] != "h1" {
		t.Fatalf("unexpected header: %v", header)
	}
}
