package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestInitrepWorkflow validates JSON output for initialization replication.
func TestInitrepWorkflow(t *testing.T) {
	out, err := execCommand("initrep", "start", "--json")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var startResp map[string]string
	if err := json.Unmarshal([]byte(out), &startResp); err != nil {
		t.Fatalf("unmarshal start: %v", err)
	}
	if startResp["status"] != "started" {
		t.Fatalf("unexpected start response: %+v", startResp)
	}

	out, err = execCommand("initrep", "stop", "--json")
	if err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var stopResp map[string]string
	if err := json.Unmarshal([]byte(out), &stopResp); err != nil {
		t.Fatalf("unmarshal stop: %v", err)
	}
	if stopResp["status"] != "stopped" {
		t.Fatalf("unexpected stop response: %+v", stopResp)
	}
}
