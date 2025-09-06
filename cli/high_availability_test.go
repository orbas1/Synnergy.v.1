package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestHighAvailabilityWorkflow verifies failover manager CLI commands output JSON.
func TestHighAvailabilityWorkflow(t *testing.T) {
	if _, err := execCommand("highavailability", "init", "p1", "1", "--json"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}

	if _, err := execCommand("highavailability", "add", "b1", "--json"); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}

	if _, err := execCommand("highavailability", "heartbeat", "b1", "--json"); err != nil {
		t.Fatalf("heartbeat: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}

	out, err := execCommand("highavailability", "active", "--json")
	if err != nil {
		t.Fatalf("active: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var resp map[string]string
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["active"] != "p1" {
		t.Fatalf("unexpected active node: %v", resp)
	}

	t.Cleanup(func() { failover = nil })
}
