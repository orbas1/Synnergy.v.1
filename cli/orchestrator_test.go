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
	if !strings.Contains(out, "sealed=") {
		t.Fatalf("expected sealed status in output: %s", out)
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
	if !strings.Contains(buf.String(), "gasLastSyncedAt") {
		t.Fatalf("expected gas sync timestamp in json: %s", buf.String())
	}
}

func TestOrchestratorBootstrapCommand(t *testing.T) {
	cmd := RootCmd()
	defer cmd.SetArgs(nil)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"orchestrator", "bootstrap"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing bootstrap: %v", err)
	}
	if !strings.Contains(buf.String(), "Bootstrap complete") {
		t.Fatalf("expected bootstrap summary, got: %s", buf.String())
	}
}
