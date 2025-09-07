package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

// TestSyn800TokenWorkflow covers registering and updating assets.
func TestSyn800TokenWorkflow(t *testing.T) {
	assetRegistry = core.NewAssetRegistry()

	out, err := execCommand("syn800_token", "register", "A1", "desc", "100", "loc", "type", "cert")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if !strings.Contains(out, "asset registered") {
		t.Fatalf("unexpected register output: %s", out)
	}

	out, err = execCommand("syn800_token", "update", "A1", "150")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if !strings.Contains(out, "valuation updated") {
		t.Fatalf("unexpected update output: %s", out)
	}

	out, err = execCommand("syn800_token", "info", "A1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if !strings.Contains(out, "\"Valuation\": 150") {
		t.Fatalf("unexpected info output: %s", out)
	}
}
