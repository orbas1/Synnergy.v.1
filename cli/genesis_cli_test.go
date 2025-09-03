package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func TestGenesisInitBlock(t *testing.T) {
	ledger = core.NewLedger()
	currentNode = core.NewNode("node1", "localhost", ledger)
	out, err := execCommand("genesis", "init-block")
	if err != nil {
		t.Fatalf("init-block failed: %v", err)
	}
	if !strings.Contains(out, "hash:") {
		t.Fatalf("unexpected output: %s", out)
	}
	// reset state for other tests
	ledger = core.NewLedger()
	currentNode = core.NewNode("node1", "localhost", ledger)
}
