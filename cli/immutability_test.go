package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

// TestImmutabilityWorkflow ensures immutability CLI emits JSON.
func TestImmutabilityWorkflow(t *testing.T) {
	origLedger := ledger
	ledger = core.NewLedger()
	t.Cleanup(func() { ledger = origLedger; enforcer = nil })

	out, err := execCommand("immutability", "init", "--json")
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.LastIndex(out, "{"); idx != -1 {
		out = out[idx:]
	}
	var initResp map[string]string
	if err := json.Unmarshal([]byte(out), &initResp); err != nil {
		t.Fatalf("unmarshal init: %v", err)
	}
	if initResp["genesis"] == "" {
		t.Fatalf("expected genesis hash, got %+v", initResp)
	}

	out, err = execCommand("immutability", "check", "--json")
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.LastIndex(out, "{"); idx != -1 {
		out = out[idx:]
	}
	var chkResp map[string]string
	if err := json.Unmarshal([]byte(out), &chkResp); err != nil {
		t.Fatalf("unmarshal check: %v", err)
	}
	if chkResp["status"] != "ok" {
		t.Fatalf("unexpected check response: %+v", chkResp)
	}
}
