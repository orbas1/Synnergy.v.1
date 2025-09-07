package cli

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestWalletNewCLI(t *testing.T) {
	path := "wallet_cli.json"
	out, err := execCommand("wallet", "new", "--out", path, "--password", "pass", "--json")
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	defer os.Remove(path)
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
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
