package cli

import (
	"strings"
	"testing"
)

func TestRegNodeCLI(t *testing.T) {
	// ensure regulation exists for rejection
	if out, err := execCommand("regulator", "add", "R1", "US", "cap", "100"); err != nil {
		t.Fatalf("setup add failed: %v", err)
	} else if !strings.Contains(out, "added") {
		t.Fatalf("setup add unexpected: %s", out)
	}

	out, err := execCommand("regnode", "approve", "alice", "150")
	if err != nil {
		t.Fatalf("approve failed: %v", err)
	}
	if !strings.Contains(out, "rejected") || !strings.Contains(out, "gas cost") {
		t.Fatalf("unexpected approve output: %s", out)
	}

	out, err = execCommand("regnode", "flag", "alice", "suspicious")
	if err != nil {
		t.Fatalf("flag failed: %v", err)
	}
	if !strings.Contains(out, "flagged") {
		t.Fatalf("unexpected flag output: %s", out)
	}

	out, err = execCommand("regnode", "logs", "alice")
	if err != nil {
		t.Fatalf("logs failed: %v", err)
	}
	if !strings.Contains(out, "R1") || !strings.Contains(out, "suspicious") {
		t.Fatalf("missing logs: %s", out)
	}
}
