package cli

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWalletNewCLI(t *testing.T) {
	path := "wallet_cli.json"
	out, err := execCommand("wallet", "new", "--out", path, "--password", "pass", "--json")
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	defer os.Remove(path)
	var resp struct {
		Address string `json:"address"`
		Path    string `json:"path"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("json: %v", err)
	}
	if len(resp.Address) != 40 {
		t.Fatalf("bad address: %s", resp.Address)
	}
	if resp.Path != path {
		t.Fatalf("expected path %s", path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("wallet file missing")
	}
}
