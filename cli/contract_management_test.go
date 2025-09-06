package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestContractManagerInfoError checks that querying a missing contract prints an error.
func TestContractManagerInfoError(t *testing.T) {
	out, err := execCommand("--json", "contract-mgr", "info", "missing")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	start := strings.LastIndex(out, "{")
	end := strings.LastIndex(out, "}")
	if start != -1 && end != -1 {
		out = out[start : end+1]
	}
	var res map[string]any
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res["error"] == nil {
		t.Fatalf("expected error field, got %v", res)
	}
}
