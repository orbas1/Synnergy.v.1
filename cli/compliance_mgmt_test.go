package cli

import (
	"strings"
	"testing"
)

// TestComplianceMgmtStatusDefault checks default status is not suspended or whitelisted.
func TestComplianceMgmtStatusDefault(t *testing.T) {
	out, err := execCommand("compliance_management", "--json", "status", "addr1")
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	if !strings.Contains(out, "\"suspended\":false") || !strings.Contains(out, "\"whitelisted\":false") {
		t.Fatalf("unexpected output: %s", out)
	}
}
