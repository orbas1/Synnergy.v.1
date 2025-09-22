package scripts_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type snapshotStatus struct {
	Network         string `json:"network"`
	Height          int    `json:"height"`
	Hash            string `json:"hash"`
	FinalizedHeight int    `json:"finalized_height"`
	FinalizedHash   string `json:"finalized_hash"`
}

func runSnapshot(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command("bash", append([]string{"../../scripts/chain_state_snapshot.sh"}, args...)...)
	return cmd.CombinedOutput()
}

func writeSnapshotStatus(t *testing.T, dir string, status snapshotStatus) string {
	t.Helper()
	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("failed to marshal status: %v", err)
	}
	path := filepath.Join(dir, "status.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write status file: %v", err)
	}
	return path
}

func TestChainStateSnapshotUsage(t *testing.T) {
	out, err := runSnapshot(t)
	if err == nil {
		t.Fatalf("expected failure without required arguments")
	}
	if !strings.Contains(string(out), "Usage: chain_state_snapshot.sh") {
		t.Fatalf("expected usage output, got: %s", out)
	}
}

func TestChainStateSnapshotCaptureAndRetention(t *testing.T) {
	dir := t.TempDir()
	status := writeSnapshotStatus(t, dir, snapshotStatus{
		Network:         "stage98",
		Height:          1280,
		Hash:            "0x123",
		FinalizedHeight: 1270,
		FinalizedHash:   "0x122",
	})

	metrics := filepath.Join(dir, "metrics.json")
	args := []string{
		"--network", "stage98",
		"--output-dir", dir,
		"--status-file", status,
		"--retain", "1",
		"--metrics-file", metrics,
		"--timestamp", "20240101T000000Z",
	}
	if out, err := runSnapshot(t, args...); err != nil {
		t.Fatalf("snapshot capture failed: %v (%s)", err, out)
	}

	firstSnapshot := filepath.Join(dir, "snapshot-stage98-20240101T000000Z.json")
	if _, err := os.Stat(firstSnapshot); err != nil {
		t.Fatalf("expected snapshot file %s: %v", firstSnapshot, err)
	}

	// Second capture with new timestamp should prune the previous snapshot due to retain=1
	args[len(args)-1] = "20240101T010000Z"
	if out, err := runSnapshot(t, args...); err != nil {
		t.Fatalf("second snapshot capture failed: %v (%s)", err, out)
	}

	files, err := filepath.Glob(filepath.Join(dir, "snapshot-stage98-*.json"))
	if err != nil {
		t.Fatalf("failed to glob snapshots: %v", err)
	}
	if len(files) != 1 || !strings.Contains(files[0], "20240101T010000Z") {
		t.Fatalf("expected only the latest snapshot to remain, got: %v", files)
	}

	data, err := os.ReadFile(metrics)
	if err != nil {
		t.Fatalf("failed to read metrics: %v", err)
	}
	if !strings.Contains(string(data), "\"command\":") && !strings.Contains(string(data), "snapshot_path") {
		// metrics include snapshot_path; ensure data is non-empty and structured
		t.Fatalf("unexpected metrics payload: %s", data)
	}
}
