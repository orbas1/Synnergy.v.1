package cli

import (
	"encoding/json"
	"testing"
	"time"

	"synnergy/core"
)

// TestSyn700Workflow covers registering IP assets and royalties.
func TestSyn700Workflow(t *testing.T) {
	ipRegistry = core.NewIPRegistry()

	out, err := execCommand("--json", "syn700", "register", "IP1", "Title", "Desc", "Creator", "Owner")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	payload := map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse register: %v", err)
	}
	if payload["status"] != "registered" {
		t.Fatalf("unexpected register output: %s", out)
	}

	expiry := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	out, err = execCommand("--json", "syn700", "license", "IP1", "LIC1", "exclusive", "Bob", "10", "--expiry", expiry)
	if err != nil {
		t.Fatalf("license: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse license: %v", err)
	}
	if payload["status"] != "license created" {
		t.Fatalf("unexpected license output: %s", out)
	}

	out, err = execCommand("--json", "syn700", "royalty", "IP1", "LIC1", "Bob", "5")
	if err != nil {
		t.Fatalf("royalty: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse royalty: %v", err)
	}
	if payload["status"] != "royalty recorded" {
		t.Fatalf("unexpected royalty output: %s", out)
	}

	out, err = execCommand("--json", "syn700", "royalties", "IP1", "LIC1")
	if err != nil {
		t.Fatalf("royalty summary: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse summary: %v", err)
	}
	if payload["total"].(float64) != 5 {
		t.Fatalf("unexpected royalty total: %s", out)
	}

	out, err = execCommand("--json", "syn700", "info", "IP1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse info: %v", err)
	}
	if payload["Title"] != "Title" {
		t.Fatalf("unexpected info output: %s", out)
	}

	out, err = execCommand("--json", "syn700", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	payload = map[string]any{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse list: %v", err)
	}
	ids := payload["ids"].([]any)
	if len(ids) != 1 || ids[0] != "IP1" {
		t.Fatalf("unexpected ids: %v", ids)
	}
}
