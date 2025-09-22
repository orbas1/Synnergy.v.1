package cli

import (
	"encoding/json"
	"testing"
)

// TestSyn500Workflow verifies creating a token, granting usage and recording audits.
func TestSyn500Workflow(t *testing.T) {
	syn500Token = nil

	out, err := execCommand("--json", "syn500", "create", "--name", "Utility", "--symbol", "UTL", "--owner", "owner", "--dec", "2", "--supply", "1000")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	payload := map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse create output: %v", err)
	}
	if payload["status"] != "token created" {
		t.Fatalf("unexpected create output: %s", out)
	}

	out, err = execCommand("--json", "syn500", "grant", "alice", "--tier", "1", "--max", "10")
	if err != nil {
		t.Fatalf("grant: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse grant output: %v", err)
	}
	if payload["status"] != "granted" || payload["address"] != "alice" {
		t.Fatalf("unexpected grant output: %s", out)
	}

	out, err = execCommand("--json", "syn500", "use", "alice", "--amount", "3", "--note", "api")
	if err != nil {
		t.Fatalf("use: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse use output: %v", err)
	}
	if payload["status"] != "usage recorded" || payload["remaining"].(float64) != 7 {
		t.Fatalf("unexpected use output: %s", out)
	}

	out, err = execCommand("--json", "syn500", "snapshot")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse snapshot: %v", err)
	}
	grants := payload["grants"].(map[string]any)
	alice := grants["alice"].(map[string]any)
	if alice["Used"].(float64) != 3 {
		t.Fatalf("unexpected snapshot usage: %v", alice)
	}

	out, err = execCommand("--json", "syn500", "audit", "--limit", "1")
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse audit: %v", err)
	}
	entries := payload["entries"].([]any)
	if len(entries) != 1 {
		t.Fatalf("unexpected audit entries: %v", entries)
	}
	entry := entries[0].(map[string]any)
	if entry["Amount"].(float64) != 3 {
		t.Fatalf("unexpected audit amount: %v", entry)
	}

	out, err = execCommand("--json", "syn500", "revoke", "alice")
	if err != nil {
		t.Fatalf("revoke: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse revoke: %v", err)
	}
	if payload["status"] != "revoked" {
		t.Fatalf("unexpected revoke output: %s", out)
	}
}
