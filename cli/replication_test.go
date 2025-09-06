package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestReplicationCLI exercises the replication command set with JSON output.
func TestReplicationCLI(t *testing.T) {
	out, err := execCommand("replication", "start", "--json")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var resp map[string]string
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal start: %v", err)
	}
	if resp["status"] != "started" {
		t.Fatalf("unexpected start response: %v", resp)
	}

	out, err = execCommand("replication", "status", "--json")
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
	if !st["running"] {
		t.Fatalf("expected running true, got %v", st)
	}

	out, err = execCommand("replication", "replicate", "h1", "--json")
	if err != nil {
		t.Fatalf("replicate: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var rep map[string]string
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatalf("unmarshal replicate: %v", err)
	}
	if rep["status"] != "replicated" {
		t.Fatalf("unexpected replicate response: %v", rep)
	}

	if _, err = execCommand("replication", "stop", "--json"); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
