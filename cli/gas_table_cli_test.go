package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"synnergy/core"
)

// TestGasSnapshotJSON ensures the snapshot command emits valid JSON and reflects
// runtime updates made via the set subcommand.
func TestGasSnapshotJSON(t *testing.T) {
	if _, err := execCommand("gas", "set", "1", "5"); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	out, err := execCommand("gas", "snapshot", "--json")
	if err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	jsonPart := lines[len(lines)-1]
	var m map[string]uint64
	if err := json.Unmarshal([]byte(jsonPart), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(m) == 0 {
		t.Fatalf("expected non-empty snapshot")
	}
}

// TestGasSnapshotWrite verifies the snapshot command can persist the gas table
// to a file in the same deterministic JSON format produced by the core helpers.
func TestGasSnapshotWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	if _, err := execCommand("gas", "snapshot", "--out", path); err != nil {
		t.Fatalf("snapshot write: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	expected, err := core.GasTableSnapshotJSON()
	if err != nil {
		t.Fatalf("snapshot json: %v", err)
	}
	if !bytes.Equal(data, expected) {
		t.Fatalf("snapshot mismatch: %q vs %q", data, expected)
	}
}
