package cli

import (
	"strings"
	"testing"
)

// TestMobileMiningCLI covers basic mobile mining operations.
func TestMobileMiningCLI(t *testing.T) {
	out, err := execCommand("mobile_mining", "start", "--json")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if !strings.Contains(out, "started") {
		t.Fatalf("unexpected start output: %s", out)
	}

	out, err = execCommand("mobile_mining", "status", "--json")
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	if !strings.Contains(out, "mining") {
		t.Fatalf("unexpected status output: %s", out)
	}

	out, err = execCommand("mobile_mining", "set-power", "50", "--json")
	if err != nil {
		t.Fatalf("set-power failed: %v", err)
	}
	if !strings.Contains(out, "50") {
		t.Fatalf("unexpected set-power output: %s", out)
	}

	out, err = execCommand("mobile_mining", "power", "--json")
	if err != nil {
		t.Fatalf("power failed: %v", err)
	}
	if !strings.Contains(out, "50") {
		t.Fatalf("unexpected power output: %s", out)
	}

	out, err = execCommand("mobile_mining", "stop", "--json")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if !strings.Contains(out, "stopped") {
		t.Fatalf("unexpected stop output: %s", out)
	}
}
