package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestRollupManagerCLI verifies pause/resume/status operations.
func TestRollupManagerCLI(t *testing.T) {
	out, err := execCommand("rollupmgr", "pause", "--json")
	if err != nil {
		t.Fatalf("pause: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var resp map[string]string
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal pause: %v", err)
	}
	if resp["status"] != "paused" {
		t.Fatalf("unexpected pause response: %v", resp)
	}

	out, err = execCommand("rollupmgr", "status", "--json")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var st map[string]bool
	if err := json.Unmarshal([]byte(out), &st); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}
	if !st["paused"] {
		t.Fatalf("expected paused true, got %v", st)
	}

	out, err = execCommand("rollupmgr", "resume", "--json")
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal resume: %v", err)
	}
	if resp["status"] != "resumed" {
		t.Fatalf("unexpected resume response: %v", resp)
	}
}
