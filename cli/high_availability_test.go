package cli

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"synnergy/core"
)

// TestHighAvailabilityWorkflow verifies failover manager CLI commands emit structured JSON and honour signed heartbeats.
func TestHighAvailabilityWorkflow(t *testing.T) {
	primaryWallet, err := core.NewWallet()
	if err != nil {
		t.Fatalf("wallet init: %v", err)
	}
	primaryPub := base64.StdEncoding.EncodeToString(primaryWallet.PublicKeyBytes())

	if _, err := execCommand("highavailability", "init", "primary-node", "1", "--primary-region", "us-east", "--primary-role", "orchestrator", "--primary-pubkey", primaryPub, "--json"); err != nil {
		t.Fatalf("init: %v", err)
	}

	backupWallet, err := core.NewWallet()
	if err != nil {
		t.Fatalf("backup wallet: %v", err)
	}
	backupPub := base64.StdEncoding.EncodeToString(backupWallet.PublicKeyBytes())

	if _, err := execCommand("highavailability", "add", "backup-1", "--region", "eu-west", "--role", "validator", "--pubkey", backupPub, "--json"); err != nil {
		t.Fatalf("add: %v", err)
	}

	payload := "backup-1-heartbeat"
	sig, err := backupWallet.SignMessage([]byte(payload))
	if err != nil {
		t.Fatalf("sign heartbeat: %v", err)
	}
	sigB64 := base64.StdEncoding.EncodeToString(sig)

	if _, err := execCommand("highavailability", "heartbeat", "backup-1", "--payload", payload, "--signature", sigB64, "--latency", "5ms", "--json"); err != nil {
		t.Fatalf("heartbeat: %v", err)
	}

	activeOut, err := execCommand("highavailability", "active", "--json")
	if err != nil {
		t.Fatalf("active: %v", err)
	}
	var activeResp map[string]string
	if err := json.Unmarshal([]byte(jsonPayload(activeOut)), &activeResp); err != nil {
		t.Fatalf("active json: %v", err)
	}
	if activeResp["active"] != "primary-node" {
		t.Fatalf("unexpected active node: %v", activeResp)
	}

        reportOut, err := execCommand("highavailability", "report", "--json")
        if err != nil {
                t.Fatalf("report: %v", err)
        }
        var report map[string]any
        if err := json.Unmarshal([]byte(jsonPayload(reportOut)), &report); err != nil {
                t.Fatalf("report json: %v", err)
        }
	if report["activeNode"].(string) != "primary-node" {
		t.Fatalf("unexpected report active node: %v", report["activeNode"])
	}
	if report["signature"].(string) == "" {
		t.Fatalf("report signature missing: %v", report)
	}
	backups, ok := report["backups"].([]any)
	if !ok || len(backups) == 0 {
		t.Fatalf("expected backups in report: %v", report)
	}
}
