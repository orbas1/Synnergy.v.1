package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

// execCLI executes the root command with the given args and returns output.
func execCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	out := strings.TrimSpace(buf.String())
	rootCmd.SetOut(nil)
	rootCmd.SetErr(nil)
	rootCmd.SetArgs(nil)
	return out, err
}

func TestCrossChainBridgeCommands(t *testing.T) {
	bridgeRegistry = core.NewBridgeRegistry()
	if cmd, _, err := rootCmd.Find([]string{"cross_chain", "list"}); err == nil {
		_ = cmd.Flags().Set("json", "false")
	}
	if cmd, _, err := rootCmd.Find([]string{"cross_chain", "get"}); err == nil {
		_ = cmd.Flags().Set("json", "false")
	}

	out, err := execCLI(t, "cross_chain", "register", "chainA", "chainB", "relayer1")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if !strings.Contains(out, "bridge-1") {
		t.Fatalf("expected bridge ID, got %q", out)
	}

	out, err = execCLI(t, "cross_chain", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "bridge-1: chainA -> chainB") {
		t.Fatalf("unexpected list output %q", out)
	}

	out, err = execCLI(t, "cross_chain", "get", "bridge-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !strings.Contains(out, "relayers=1") {
		t.Fatalf("expected relayers=1, got %q", out)
	}

	if _, err = execCLI(t, "cross_chain", "authorize", "bridge-1", "relayer2"); err != nil {
		t.Fatalf("authorize: %v", err)
	}

	out, err = execCLI(t, "cross_chain", "get", "--json", "bridge-1")
	if err != nil {
		t.Fatalf("get json: %v", err)
	}
	var b core.Bridge
	if err := json.Unmarshal([]byte(out), &b); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(b.Relayers) != 2 {
		t.Fatalf("expected 2 relayers, got %d", len(b.Relayers))
	}

	if _, err = execCLI(t, "cross_chain", "revoke", "bridge-1", "relayer2"); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	// reset json flag before calling without --json
	if cmd, _, errFind := rootCmd.Find([]string{"cross_chain", "get"}); errFind == nil {
		_ = cmd.Flags().Set("json", "false")
	}
	out, err = execCLI(t, "cross_chain", "get", "bridge-1")
	if err != nil {
		t.Fatalf("final get: %v", err)
	}
	if !strings.Contains(out, "relayers=1") {
		t.Fatalf("expected relayers=1, got %q", out)
	}
}
