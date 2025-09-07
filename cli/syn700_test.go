package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

// TestSyn700Workflow covers registering IP assets and royalties.
func TestSyn700Workflow(t *testing.T) {
	ipRegistry = core.NewIPRegistry()

	out, err := execCommand("syn700", "register", "IP1", "Title", "Desc", "Creator", "Owner")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if !strings.Contains(out, "registered") {
		t.Fatalf("unexpected register output: %s", out)
	}

	out, err = execCommand("syn700", "license", "IP1", "LIC1", "exclusive", "Bob", "10")
	if err != nil {
		t.Fatalf("license: %v", err)
	}
	if !strings.Contains(out, "license created") {
		t.Fatalf("unexpected license output: %s", out)
	}

	out, err = execCommand("syn700", "royalty", "IP1", "LIC1", "Bob", "5")
	if err != nil {
		t.Fatalf("royalty: %v", err)
	}
	if !strings.Contains(out, "royalty recorded") {
		t.Fatalf("unexpected royalty output: %s", out)
	}

	out, err = execCommand("syn700", "info", "IP1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if !strings.Contains(out, "\"Title\": \"Title\"") {
		t.Fatalf("unexpected info output: %s", out)
	}
}
