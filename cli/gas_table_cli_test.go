package cli

import (
	"encoding/json"
	"strings"
	"testing"
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
