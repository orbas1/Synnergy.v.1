package cli

import (
	"strings"
	"testing"
)

func TestGenesisInitOnce(t *testing.T) {
	out, err := execCommand("genesis", "init")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if !strings.Contains(out, "genesis block") {
		t.Fatalf("unexpected output: %s", out)
	}
	if _, err := execCommand("genesis", "init"); err == nil {
		t.Fatalf("expected error on second init")
	}
}
