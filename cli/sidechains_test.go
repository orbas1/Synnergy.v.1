package cli

import (
    "encoding/json"
    "strings"
    "testing"
)

// TestSidechainListEmpty verifies listing side-chains returns an empty set by default.
func TestSidechainListEmpty(t *testing.T) {
    out, err := execCommand("--json", "sidechain", "list")
    if err != nil {
        t.Fatalf("list: %v", err)
    }
    if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
        t.Fatalf("reset json: %v", err)
    }
    if idx := strings.Index(out, "\n"); idx != -1 {
        out = out[idx+1:]
    }
    var chains []any
    if err := json.Unmarshal([]byte(out), &chains); err != nil {
        t.Fatalf("unmarshal: %v", err)
    }
    if len(chains) != 0 {
        t.Fatalf("expected empty list, got %d", len(chains))
    }
}

