package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

// TestIdentityWorkflow exercises identity CLI commands and validates JSON output.
func TestIdentityWorkflow(t *testing.T) {
	out, err := execCommand("identity", "register", "addr1", "Alice", "1990-01-01", "US", "--json")
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

	out, err = execCommand("identity", "verify", "addr1", "passport", "--json")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var verResp map[string]string
	if err := json.Unmarshal([]byte(out), &verResp); err != nil {
		t.Fatalf("unmarshal verify: %v", err)
	}
	if verResp["status"] != "verified" || verResp["method"] != "passport" {
		t.Fatalf("unexpected verify response: %+v", verResp)
	}

	out, err = execCommand("identity", "info", "addr1", "--json")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var infoResp map[string]string
	if err := json.Unmarshal([]byte(out), &infoResp); err != nil {
		t.Fatalf("unmarshal info: %v", err)
	}
	if infoResp["name"] != "Alice" || infoResp["dob"] != "1990-01-01" || infoResp["nationality"] != "US" {
		t.Fatalf("unexpected info response: %+v", infoResp)
	}

	out, err = execCommand("identity", "logs", "addr1", "--json")
	if err != nil {
		t.Fatalf("logs: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var logs []map[string]string
	if err := json.Unmarshal([]byte(out), &logs); err != nil {
		t.Fatalf("unmarshal logs: %v", err)
	}
	if len(logs) != 1 || logs[0]["Method"] != "passport" {
		t.Fatalf("unexpected logs: %+v", logs)
	}

	t.Cleanup(func() { identitySvc = core.NewIdentityService() })
}
