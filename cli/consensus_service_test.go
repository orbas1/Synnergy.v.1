package cli

import (
	"strings"
	"testing"
)

// TestConsensusServiceStartStop exercises start, info and stop commands.
func TestConsensusServiceStartStop(t *testing.T) {
	if _, err := execCommand("consensus-service", "start", "1"); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	out, err := execCommand("consensus-service", "info")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "running: true") {
		t.Fatalf("unexpected info output: %s", out)
	}
	if _, err := execCommand("consensus-service", "stop"); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}
