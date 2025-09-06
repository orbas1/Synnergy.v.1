package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

// TestIDWalletWorkflow verifies idwallet CLI commands emit JSON responses.
func TestIDWalletWorkflow(t *testing.T) {
	out, err := execCommand("idwallet", "register", "addr1", "info1", "--json")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var regResp map[string]string
	if err := json.Unmarshal([]byte(out), &regResp); err != nil {
		t.Fatalf("unmarshal register: %v", err)
	}
	if regResp["status"] != "registered" || regResp["address"] != "addr1" {
		t.Fatalf("unexpected register response: %+v", regResp)
	}

	out, err = execCommand("idwallet", "check", "addr1", "--json")
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var chkResp map[string]string
	if err := json.Unmarshal([]byte(out), &chkResp); err != nil {
		t.Fatalf("unmarshal check: %v", err)
	}
	if chkResp["info"] != "info1" {
		t.Fatalf("unexpected check response: %+v", chkResp)
	}

	t.Cleanup(func() { idRegistry = core.NewIDRegistry() })
}
