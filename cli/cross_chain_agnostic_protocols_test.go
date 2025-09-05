package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

func TestCrossChainAgnosticProtocolsCommands(t *testing.T) {
	protocolRegistry = core.NewProtocolRegistry()
	if cmd, _, err := rootCmd.Find([]string{"cross_chain_agnostic_protocols", "list"}); err == nil {
		_ = cmd.Flags().Set("json", "false")
	}
	if cmd, _, err := rootCmd.Find([]string{"cross_chain_agnostic_protocols", "get"}); err == nil {
		_ = cmd.Flags().Set("json", "false")
	}

	out, err := execCLI(t, "cross_chain_agnostic_protocols", "register", "protoA")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if !strings.Contains(out, "1") {
		t.Fatalf("expected protocol ID, got %q", out)
	}

	out, err = execCLI(t, "cross_chain_agnostic_protocols", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "1: protoA") {
		t.Fatalf("unexpected list output %q", out)
	}

	out, err = execCLI(t, "cross_chain_agnostic_protocols", "get", "1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !strings.Contains(out, "1: protoA") {
		t.Fatalf("unexpected get output %q", out)
	}

	out, err = execCLI(t, "cross_chain_agnostic_protocols", "get", "--json", "1")
	if err != nil {
		t.Fatalf("get json: %v", err)
	}
	var p core.ProtocolDefinition
	if err := json.Unmarshal([]byte(out), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Name != "protoA" {
		t.Fatalf("expected protoA, got %s", p.Name)
	}
	if cmd, _, err := rootCmd.Find([]string{"cross_chain_agnostic_protocols", "get"}); err == nil {
		_ = cmd.Flags().Set("json", "false")
	}
}
