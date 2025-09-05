package cli

import (
	"strings"
	"testing"
)

// TestConsensusSpecificNodeCreateInfo covers node creation and info output.
func TestConsensusSpecificNodeCreateInfo(t *testing.T) {
	if _, err := execCommand("consensus-node", "create", "pow", "n1", "addr1"); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	out, err := execCommand("consensus-node", "info")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "ID: n1") {
		t.Fatalf("unexpected output: %s", out)
	}
}
