package cli

import (
	"strings"
	"testing"

	"synnergy/internal/tokens"
)

// TestSyn70Workflow covers basic SYN70 asset lifecycle operations.
func TestSyn70Workflow(t *testing.T) {
	tokenRegistry = tokens.NewRegistry()
	syn70 = nil

	out, err := execCommand("syn70", "init", "--name", "Game", "--symbol", "G70", "--decimals", "0")
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	if !strings.Contains(out, "syn70 initialised") {
		t.Fatalf("unexpected init output: %s", out)
	}

	out, err = execCommand("syn70", "register", "a1", "bob", "Sword", "RPG")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if !strings.Contains(out, "asset registered") {
		t.Fatalf("unexpected register output: %s", out)
	}

	out, err = execCommand("syn70", "transfer", "a1", "alice")
	if err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if !strings.Contains(out, "asset transferred") {
		t.Fatalf("unexpected transfer output: %s", out)
	}

	out, err = execCommand("syn70", "setattr", "a1", "damage", "10")
	if err != nil {
		t.Fatalf("setattr: %v", err)
	}
	if !strings.Contains(out, "attribute set") {
		t.Fatalf("unexpected setattr output: %s", out)
	}

	out, err = execCommand("syn70", "achievement", "a1", "boss")
	if err != nil {
		t.Fatalf("achievement: %v", err)
	}
	if !strings.Contains(out, "achievement recorded") {
		t.Fatalf("unexpected achievement output: %s", out)
	}

	out, err = execCommand("syn70", "info", "a1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if !strings.Contains(out, "Owner:alice") || !strings.Contains(out, "damage=10") {
		t.Fatalf("unexpected info output: %s", out)
	}

	out, err = execCommand("syn70", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "a1 alice Sword RPG") {
		t.Fatalf("unexpected list output: %s", out)
	}

	out, err = execCommand("syn70", "balance", "alice")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if strings.TrimSpace(out) != "1" {
		t.Fatalf("unexpected balance output: %s", out)
	}
}
