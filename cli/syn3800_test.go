package cli

import (
	"encoding/json"
	"testing"

	"synnergy/core"
)

func TestSyn3800Lifecycle(t *testing.T) {
	grantRegistry = core.NewGrantRegistry()
	out, err := execCommand("syn3800", "create", "alice", "research", "100")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if out != "1" {
		t.Fatalf("expected id 1, got %s", out)
	}
	if _, err := execCommand("syn3800", "release", "1", "40", "phase1"); err != nil {
		t.Fatalf("release: %v", err)
	}
	out, err = execCommand("syn3800", "get", "1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	var g struct{ Released uint64 }
	if err := json.Unmarshal([]byte(out), &g); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if g.Released != 40 {
		t.Fatalf("expected released 40, got %d", g.Released)
	}
	out, err = execCommand("syn3800", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var gs []struct{}
	if err := json.Unmarshal([]byte(out), &gs); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if len(gs) != 1 {
		t.Fatalf("expected 1 grant, got %d", len(gs))
	}
}

func TestSyn3800Validation(t *testing.T) {
	grantRegistry = core.NewGrantRegistry()
	if _, err := execCommand("syn3800", "create", "", "name", "10"); err == nil {
		t.Fatal("expected error for missing beneficiary")
	}
	if _, err := execCommand("syn3800", "create", "bob", "grant", "0"); err == nil {
		t.Fatal("expected error for amount")
	}
	if _, err := execCommand("syn3800", "release", "1", "10"); err == nil {
		t.Fatal("expected error for unknown id")
	}
}
