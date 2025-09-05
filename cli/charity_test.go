package cli

import (
	"strings"
	"testing"
)

// TestCharityRegistrationJSON verifies registration info is emitted as JSON when requested.
func TestCharityRegistrationJSON(t *testing.T) {
	if _, err := execCommand("charity_pool", "register", "addr1", "1", "CharityOne"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	out, err := execCommand("charity_pool", "--json", "registration", "addr1")
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}
	if !strings.Contains(out, "CharityOne") {
		t.Fatalf("unexpected output: %s", out)
	}
}

// TestCharityBalancesJSON ensures balances command respects the --json flag.
func TestCharityBalancesJSON(t *testing.T) {
	if _, err := execCommand("charity_mgmt", "donate", "addr2", "10"); err != nil {
		t.Fatalf("donate failed: %v", err)
	}
	out, err := execCommand("charity_mgmt", "--json", "balances")
	if err != nil {
		t.Fatalf("balances failed: %v", err)
	}
	if !strings.Contains(out, "\"pool\":") {
		t.Fatalf("expected json output, got %s", out)
	}
}
