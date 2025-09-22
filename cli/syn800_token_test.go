package cli

import (
	"encoding/json"
	"testing"

	"synnergy/core"
)

// TestSyn800TokenWorkflow covers registering and updating assets.
func TestSyn800TokenWorkflow(t *testing.T) {
	assetRegistry = core.NewAssetRegistry()

	out, err := execCommand("--json", "syn800_token", "register", "A1", "desc", "100", "loc", "type", "cert")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	payload := map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse register: %v", err)
	}
	if payload["status"] != "asset registered" {
		t.Fatalf("unexpected register output: %s", out)
	}

	out, err = execCommand("--json", "syn800_token", "custodian", "A1", "guardian")
	if err != nil {
		t.Fatalf("custodian: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse custodian: %v", err)
	}
	if payload["status"] != "custodian assigned" {
		t.Fatalf("unexpected custodian output: %s", out)
	}

	out, err = execCommand("--json", "syn800_token", "update", "A1", "150")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse update: %v", err)
	}
	if payload["value"].(float64) != 150 {
		t.Fatalf("unexpected update output: %s", out)
	}

	out, err = execCommand("--json", "syn800_token", "info", "A1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse info: %v", err)
	}
	if payload["Valuation"].(float64) != 150 {
		t.Fatalf("unexpected info output: %s", out)
	}

	out, err = execCommand("--json", "syn800_token", "snapshot")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse snapshot: %v", err)
	}
	assets := payload["assets"].([]any)
	if len(assets) != 1 {
		t.Fatalf("unexpected snapshot entries: %v", assets)
	}

	out, err = execCommand("--json", "syn800_token", "history", "--limit", "1")
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse history: %v", err)
	}
	entries := payload["entries"].([]any)
	if len(entries) != 1 {
		t.Fatalf("unexpected history entries: %v", entries)
	}
}
