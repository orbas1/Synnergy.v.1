package scripts_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type statusPayload struct {
	Network         string `json:"network"`
	Height          int    `json:"height"`
	Hash            string `json:"hash"`
	FinalizedHeight int    `json:"finalized_height"`
	FinalizedHash   string `json:"finalized_hash"`
}

func runScript(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command("bash", append([]string{"../../scripts/chain_rollback_prevention.sh"}, args...)...)
	return cmd.CombinedOutput()
}

func writeStatus(t *testing.T, dir string, payload statusPayload) string {
	t.Helper()
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	path := filepath.Join(dir, "status.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write status file: %v", err)
	}
	return path
}

func writeCheckpoint(t *testing.T, dir string, payload statusPayload) string {
	t.Helper()
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal checkpoint: %v", err)
	}
	path := filepath.Join(dir, "checkpoint.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write checkpoint file: %v", err)
	}
	return path
}

func TestChainRollbackPreventionUsage(t *testing.T) {
	out, err := runScript(t)
	if err == nil {
		t.Fatalf("expected error when no arguments provided")
	}
	if !strings.Contains(string(out), "Usage: chain_rollback_prevention.sh") {
		t.Fatalf("expected usage output, got: %s", out)
	}
}

func TestChainRollbackInitializationCreatesCheckpoint(t *testing.T) {
	dir := t.TempDir()
	status := writeStatus(t, dir, statusPayload{
		Network:         "stage98",
		Height:          1200,
		Hash:            "0xaaa",
		FinalizedHeight: 1180,
		FinalizedHash:   "0xbbb",
	})

	checkpoint := filepath.Join(dir, "checkpoint.json")
	out, err := runScript(t,
		"--network", "stage98",
		"--checkpoint-file", checkpoint,
		"--status-file", status,
		"--mitigation-action", "none",
		"--dry-run",
	)
	if err != nil {
		t.Fatalf("expected initialization to succeed, got error: %v (%s)", err, out)
	}

	if _, err := os.Stat(checkpoint); err != nil {
		t.Fatalf("expected checkpoint file to exist: %v", err)
	}
	data, err := os.ReadFile(checkpoint)
	if err != nil {
		t.Fatalf("failed reading checkpoint: %v", err)
	}
	if !strings.Contains(string(data), "\"finalized_height\": 1180") {
		t.Fatalf("expected finalized height in checkpoint, got: %s", data)
	}
}

func TestChainRollbackPreventionDetectsRollback(t *testing.T) {
	dir := t.TempDir()
	status := writeStatus(t, dir, statusPayload{
		Network:         "stage98",
		Height:          1110,
		Hash:            "0xccc",
		FinalizedHeight: 1000,
		FinalizedHash:   "0xddd",
	})

	checkpoint := writeCheckpoint(t, dir, statusPayload{
		Network:         "stage98",
		Height:          1200,
		Hash:            "0xffff",
		FinalizedHeight: 1100,
		FinalizedHash:   "0xeeee",
	})

	out, err := runScript(t,
		"--network", "stage98",
		"--checkpoint-file", checkpoint,
		"--status-file", status,
		"--mitigation-action", "none",
		"--dry-run",
	)
	if err == nil {
		t.Fatalf("expected rollback detection to fail")
	}
	if !strings.Contains(string(out), "rollback detected") {
		t.Fatalf("expected rollback warning, got: %s", out)
	}
}

func TestChainRollbackPreventionMetricsAndUpdate(t *testing.T) {
	dir := t.TempDir()
	status := writeStatus(t, dir, statusPayload{
		Network:         "stage98",
		Height:          2100,
		Hash:            "0xabc123",
		FinalizedHeight: 2000,
		FinalizedHash:   "0xabc777",
	})

	checkpoint := writeCheckpoint(t, dir, statusPayload{
		Network:         "stage98",
		Height:          1900,
		Hash:            "0xold",
		FinalizedHeight: 1850,
		FinalizedHash:   "0xoldhash",
	})

	metrics := filepath.Join(dir, "metrics.json")
	out, err := runScript(t,
		"--network", "stage98",
		"--checkpoint-file", checkpoint,
		"--status-file", status,
		"--metrics-file", metrics,
		"--mitigation-action", "none",
	)
	if err != nil {
		t.Fatalf("expected successful evaluation, got error: %v (%s)", err, out)
	}

	data, err := os.ReadFile(metrics)
	if err != nil {
		t.Fatalf("failed reading metrics: %v", err)
	}
	if !strings.Contains(string(data), "\"rollback_detected\": 0") {
		t.Fatalf("expected rollback_detected flag to be 0, got: %s", data)
	}
	if !strings.Contains(string(data), "\"updated_checkpoint\": 1") {
		t.Fatalf("expected checkpoint update flag, got: %s", data)
	}

	newCheckpoint, err := os.ReadFile(checkpoint)
	if err != nil {
		t.Fatalf("failed reading updated checkpoint: %v", err)
	}
	if !strings.Contains(string(newCheckpoint), "\"finalized_height\": 2000") {
		t.Fatalf("expected checkpoint to reflect new finalized height, got: %s", newCheckpoint)
	}
}
