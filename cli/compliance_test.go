package cli

import (
	"strings"
	"testing"
)

// TestComplianceRiskDefault ensures the risk command returns zero for unknown address.
func TestComplianceRiskDefault(t *testing.T) {
	out, err := execCommand("compliance", "--json", "risk", "addr1")
	if err != nil {
		t.Fatalf("risk failed: %v", err)
	}
	if !strings.Contains(out, "\"risk\":0") {
		t.Fatalf("unexpected output: %s", out)
	}
}
