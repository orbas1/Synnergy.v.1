package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func TestSystemHealthLoggingCommands(t *testing.T) {
	sysLogger = core.NewSystemHealthLogger()
	sysLogs = nil
	t.Cleanup(func() {
		sysLogger = core.NewSystemHealthLogger()
		sysLogs = nil
	})

	out, err := executeCLICommand(t, "system_health", "snapshot")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if !strings.Contains(out, "\"CPUUsage\"") {
		t.Fatalf("expected CPUUsage in output, got %q", out)
	}

	if _, err := executeCLICommand(t, "system_health", "log", "warn", "check"); err != nil {
		t.Fatalf("log: %v", err)
	}
	if len(sysLogs) != 1 || sysLogs[0] != "warn: check" {
		t.Fatalf("unexpected sysLogs: %v", sysLogs)
	}
}
