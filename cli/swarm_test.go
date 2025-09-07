package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestSwarmPeersEmpty verifies listing peers yields an empty set by default.
func TestSwarmPeersEmpty(t *testing.T) {
	out, err := execCommand("--json", "swarm", "peers")
	if err != nil {
		t.Fatalf("peers: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var peers []string
	if err := json.Unmarshal([]byte(out), &peers); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(peers) != 0 {
		t.Fatalf("expected empty list, got %d", len(peers))
	}
}
