package cli

import (
	"strings"
	"testing"
)

func TestFirewallAllowCheck(t *testing.T) {
	if _, err := execCommand("firewall", "allow", "1.2.3.4"); err != nil {
		t.Fatalf("allow failed: %v", err)
	}
	out, err := execCommand("firewall", "check", "1.2.3.4")
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if !strings.Contains(out, "true") {
		t.Fatalf("expected true, got %s", out)
	}
}
