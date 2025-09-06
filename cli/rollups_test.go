package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestRollupsCLI covers batch submission and listing.
func TestRollupsCLI(t *testing.T) {
	out, err := execCommand("rollups", "submit", "tx1", "--json")
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var resp map[string]string
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal submit: %v", err)
	}
	id := resp["id"]
	if id == "" {
		t.Fatalf("missing batch id")
	}

	out, err = execCommand("rollups", "list", "--json")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, id) {
		t.Fatalf("expected batch id in list: %s", out)
	}

	out, err = execCommand("rollups", "txs", id, "--json")
	if err != nil {
		t.Fatalf("txs: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var txs []string
	if err := json.Unmarshal([]byte(out), &txs); err != nil {
		t.Fatalf("unmarshal txs: %v", err)
	}
	if len(txs) != 1 || txs[0] != "tx1" {
		t.Fatalf("unexpected txs: %v", txs)
	}
}
