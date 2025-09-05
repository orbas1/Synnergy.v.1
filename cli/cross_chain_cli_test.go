package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCrossChainCLI(t *testing.T) {
	out, err := execCommand("cross_chain", "register", "chainA", "chainB", "relayer1")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if out == "" {
		t.Fatalf("expected bridge id")
	}
	out, err = execCommand("cross_chain", "list", "--json")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "chainA") {
		t.Fatalf("missing bridge in list: %s", out)
	}
}

func TestPlasmaMgmtCLI(t *testing.T) {
	if _, err := execCommand("plasma", "plasma-mgmt", "pause"); err != nil {
		t.Fatalf("pause failed: %v", err)
	}
	out, err := execCommand("plasma", "plasma-mgmt", "status")
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	if !strings.Contains(out, "true") {
		t.Fatalf("expected paused status, got %s", out)
	}
	if _, err := execCommand("plasma", "plasma-mgmt", "resume"); err != nil {
		t.Fatalf("resume failed: %v", err)
	}
}

func TestBridgeDepositCLI(t *testing.T) {
	out, err := execCommand("cross_chain_bridge", "deposit", "bridge1", "alice", "bob", "5", "--json")
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["id"] == "" {
		t.Fatalf("expected id in response: %v", resp)
	}
	if _, ok := resp["gas"]; !ok {
		t.Fatalf("expected gas field: %v", resp)
	}
}
