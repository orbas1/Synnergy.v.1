package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

// TestSyn5000Workflow verifies creating a token, placing a bet and resolving it.
func TestSyn5000Workflow(t *testing.T) {
	syn5000Tokens = map[string]*core.SYN5000Token{}

	out, err := execCommand("syn5000", "create", "--name", "Gamble", "--symbol", "GMB", "--dec", "0")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !strings.Contains(out, "token created GMB") {
		t.Fatalf("unexpected create output: %s", out)
	}

	out, err = execCommand("syn5000", "bet", "alice", "--id", "GMB", "--amt", "10", "--odds", "2", "--game", "poker")
	if err != nil {
		t.Fatalf("bet: %v", err)
	}
	if !strings.Contains(out, "bet placed 1") {
		t.Fatalf("unexpected bet output: %s", out)
	}

	out, err = execCommand("syn5000", "resolve", "--id", "GMB", "--bet", "1", "--win")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !strings.Contains(out, "payout 20") {
		t.Fatalf("unexpected resolve output: %s", out)
	}
}
