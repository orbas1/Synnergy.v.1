package cli

import (
	"strings"
	"testing"
)

func TestFullNodeCreateJSON(t *testing.T) {
	out, err := execCommand("fullnode", "create", "--id", "f1", "--mode", "archive", "--json")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "\"status\": \"created\"") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestStakingNodeStake(t *testing.T) {
	out, err := execCommand("staking_node", "stake", "addr1", "10", "--json")
	if err != nil {
		t.Fatalf("stake failed: %v", err)
	}
	if !strings.Contains(out, "\"status\": \"staked\"") {
		t.Fatalf("unexpected output: %s", out)
	}
}
