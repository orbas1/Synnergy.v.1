package cli

import (
	"strings"
	"testing"
)

func TestRegulatorCLI(t *testing.T) {
	out, err := execCommand("regulator", "add", "R1", "US", "cap", "100")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "added") {
		t.Fatalf("unexpected add output: %s", out)
	}

	out, err = execCommand("regulator", "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "R1") {
		t.Fatalf("regulation not listed: %s", out)
	}

	out, err = execCommand("regulator", "evaluate", "150")
	if err != nil {
		t.Fatalf("evaluate failed: %v", err)
	}
	if !strings.Contains(out, "R1") {
		t.Fatalf("expected violation: %s", out)
	}

	out, err = execCommand("regulator", "remove", "R1")
	if err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if !strings.Contains(out, "removed") {
		t.Fatalf("unexpected remove output: %s", out)
	}
}
