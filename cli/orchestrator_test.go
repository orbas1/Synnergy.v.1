package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestOrchestratorStatusCommand(t *testing.T) {
	defer func() { orchestratorJSON = false }()
	cmd := RootCmd()
	defer cmd.SetArgs(nil)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"orchestrator", "status"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing status: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Orchestrator status") {
		t.Fatalf("expected status output, got: %s", out)
	}
}

func TestOrchestratorStatusJSON(t *testing.T) {
	defer func() { orchestratorJSON = false }()
	cmd := RootCmd()
	defer cmd.SetArgs(nil)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"orchestrator", "status", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing status json: %v", err)
	}
	if !strings.Contains(buf.String(), "\"vmMode\"") {
		t.Fatalf("expected json payload, got: %s", buf.String())
	}
}
