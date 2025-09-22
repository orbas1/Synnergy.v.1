package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestCoinInfoJSON ensures info subcommand emits structured JSON.
func TestCoinInfoJSON(t *testing.T) {
	out, err := execCommand("coin", "--json", "info")
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}
	if !strings.Contains(out, "\"name\"") {
		t.Fatalf("expected JSON output, got %s", out)
	}
}

// TestCoinRewardValidation ensures invalid heights are rejected.
func TestCoinRewardValidation(t *testing.T) {
	if _, err := execCommand("coin", "reward", "abc"); err == nil {
		t.Fatalf("expected error for invalid height")
	}
}

func TestCoinTelemetryJSON(t *testing.T) {
	out, err := execCommand("coin", "telemetry", "--json")
	if err != nil {
		t.Fatalf("telemetry failed: %v", err)
	}
	payload := jsonPayload(out)
	var diag map[string]any
	if err := json.Unmarshal([]byte(payload), &diag); err != nil {
		t.Fatalf("decode telemetry: %v", err)
	}
	if _, ok := diag["operators"]; !ok {
		t.Fatalf("expected operators field in telemetry: %s", payload)
	}
}

func TestCoinTelemetryOperatorLifecycle(t *testing.T) {
	operator := "cli-ops-stage80"
	out, err := execCommand("coin", "telemetry", "--json", "--authorize-operator", operator)
	if err != nil {
		t.Fatalf("authorize failed: %v", err)
	}
	if !strings.Contains(jsonPayload(out), operator) {
		t.Fatalf("expected telemetry to list operator")
	}
	if _, err := execCommand("coin", "telemetry", "--json", "--operator", operator, "--transfer", "stage80-cli:1"); err != nil {
		t.Fatalf("transfer with operator failed: %v", err)
	}
	if _, err := execCommand("coin", "telemetry", "--json", "--revoke-operator", operator); err != nil {
		t.Fatalf("revoke failed: %v", err)
	}
	if _, err := execCommand("coin", "telemetry", "--json", "--operator", operator, "--transfer", "stage80-cli:1"); err == nil {
		t.Fatalf("expected revoked operator to fail")
	}
}
