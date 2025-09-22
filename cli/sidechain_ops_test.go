package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func TestSidechainOpsCommands(t *testing.T) {
	sideReg = core.NewSidechainRegistry()
	sideOps = core.NewSidechainOps(sideReg)
	t.Cleanup(func() {
		sideReg = core.NewSidechainRegistry()
		sideOps = core.NewSidechainOps(sideReg)
	})

	if _, err := sideReg.Register("chain1", "meta", nil); err != nil {
		t.Fatalf("register chain: %v", err)
	}

	if _, err := executeCLICommand(t, "sidechainops", "deposit", "chain1", "alice", "40"); err != nil {
		t.Fatalf("deposit: %v", err)
	}
	out, err := executeCLICommand(t, "sidechainops", "balance", "chain1", "alice")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if strings.TrimSpace(out) != "40" {
		t.Fatalf("expected balance 40, got %q", out)
	}

	if _, err := executeCLICommand(t, "sidechainops", "withdraw", "chain1", "alice", "10", "proof"); err != nil {
		t.Fatalf("withdraw: %v", err)
	}
	out, err = executeCLICommand(t, "sidechainops", "balance", "chain1", "alice")
	if err != nil {
		t.Fatalf("balance after withdraw: %v", err)
	}
	if strings.TrimSpace(out) != "30" {
		t.Fatalf("expected balance 30, got %q", out)
	}
}
