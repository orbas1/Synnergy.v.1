package cli

import (
	"bytes"
	"encoding/json"
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

func TestOrchestratorBootstrapCommand(t *testing.T) {
	defer func() {
		orchestratorJSON = false
		jsonOutput = false
	}()
	out, err := execCommand("--json", "orchestrator", "bootstrap", "--node-id", "stage79-cli", "--authority", "cli-authority=ops")
	if err != nil {
		t.Fatalf("bootstrap command failed: %v", err)
	}
	payload := jsonPayload(out)
	var result struct {
		NodeID         string   `json:"nodeId"`
		AuthorityNodes []string `json:"authorityNodes"`
	}
	if err := json.Unmarshal([]byte(payload), &result); err != nil {
		t.Fatalf("decode bootstrap output: %v\n%s", err, out)
	}
	if result.NodeID != "stage79-cli" {
		t.Fatalf("unexpected node id: %s", result.NodeID)
	}
	if len(result.AuthorityNodes) == 0 {
		t.Fatalf("expected authority nodes in response: %s", out)
	}
}
