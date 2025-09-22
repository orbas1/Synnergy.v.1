package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func TestPlasmaManagementCommands(t *testing.T) {
	plasmaBridge = core.NewPlasmaBridge()
	plasmaJSON = false
	t.Cleanup(func() {
		plasmaBridge = core.NewPlasmaBridge()
		plasmaJSON = false
	})

	out, err := executeCLICommand(t, "plasma-mgmt", "pause")
	if err != nil {
		t.Fatalf("pause: %v", err)
	}
	if !strings.Contains(out, "gas:") {
		t.Fatalf("expected gas output, got %q", out)
	}
	if !plasmaBridge.IsPaused() {
		t.Fatalf("bridge should be paused")
	}

	out, err = executeCLICommand(t, "plasma-mgmt", "status")
	if err != nil {
		t.Fatalf("status after pause: %v", err)
	}
	if strings.TrimSpace(out) != "true" {
		t.Fatalf("expected paused status, got %q", out)
	}

	if _, err = executeCLICommand(t, "plasma-mgmt", "resume"); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if plasmaBridge.IsPaused() {
		t.Fatalf("bridge should be resumed")
	}

	out, err = executeCLICommand(t, "plasma-mgmt", "status")
	if err != nil {
		t.Fatalf("status after resume: %v", err)
	}
	if strings.TrimSpace(out) != "false" {
		t.Fatalf("expected resumed status, got %q", out)
	}
}
