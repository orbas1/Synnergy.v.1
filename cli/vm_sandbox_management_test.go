package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func resetSandboxMgr() { sandboxMgr = core.NewSandboxManager() }

func TestSandboxStartStatus(t *testing.T) {
	resetSandboxMgr()
	if _, err := execCommand("sandbox", "start", "sb1", "contract", "10", "100"); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	out, err := execCommand("sandbox", "status", "sb1")
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "sb1") {
		t.Fatalf("unexpected output: %s", out)
	}
}
