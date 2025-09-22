package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func TestPlasmaOperationsCommands(t *testing.T) {
	plasmaBridge = core.NewPlasmaBridge()
	t.Cleanup(func() { plasmaBridge = core.NewPlasmaBridge() })

	if _, err := executeCLICommand(t, "plasma-ops", "deposit", "alice", "token", "10"); err != nil {
		t.Fatalf("deposit: %v", err)
	}
	out, err := executeCLICommand(t, "plasma-ops", "exit", "alice", "token", "5")
	if err != nil || strings.TrimSpace(out) != "1" {
		t.Fatalf("unexpected exit output: %q err=%v", out, err)
	}
	if _, err := executeCLICommand(t, "plasma-ops", "finalize", "1"); err != nil {
		t.Fatalf("finalize: %v", err)
	}
	out, err = executeCLICommand(t, "plasma-ops", "get", "1")
	if err != nil || !strings.Contains(out, "Nonce:1") {
		t.Fatalf("unexpected get output: %q err=%v", out, err)
	}
	out, err = executeCLICommand(t, "plasma-ops", "list")
	if err != nil || !strings.Contains(out, "Finalized:true") {
		t.Fatalf("unexpected list output: %q err=%v", out, err)
	}
}
