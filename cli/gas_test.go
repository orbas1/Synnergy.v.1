package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGasListCommand(t *testing.T) {
	out, err := execCommand("gas", "list")
	if err != nil {
		t.Fatalf("gas list: %v", err)
	}
	if !strings.Contains(out, "OPCODE") {
		t.Fatalf("expected table header in output: %s", out)
	}

	jsonOut, err := execCommand("--json", "gas", "list")
	if err != nil {
		t.Fatalf("gas list json: %v", err)
	}
	var entries []struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	payload := strings.TrimSpace(jsonOut)
	if payload == "" {
		t.Fatalf("empty json output")
	}
	lines := strings.Split(payload, "\n")
	for len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "gas cost:") {
		lines = lines[1:]
	}
	raw := strings.Join(lines, "\n")
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("expected at least one gas entry")
	}
	found := false
	for _, entry := range entries {
		if entry.Name == "EnterpriseBootstrap" {
			found = true
			if entry.Category == "" {
				t.Fatalf("expected category for EnterpriseBootstrap")
			}
			break
		}
	}
	if !found {
		t.Fatalf("expected EnterpriseBootstrap in gas catalogue")
	}
}
