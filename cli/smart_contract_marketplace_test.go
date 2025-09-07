package cli

import (
	"encoding/json"
	"os"
	"testing"

	"synnergy/core"
)

// TestSmartContractMarketplaceDeployJSON deploys and trades a contract verifying JSON output.
func TestSmartContractMarketplaceDeployJSON(t *testing.T) {
	marketplace = core.NewSmartContractMarketplace(core.NewSimpleVM())

	tmp, err := os.CreateTemp(t.TempDir(), "scm-*.wasm")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	if _, err := tmp.Write([]byte{0x00}); err != nil {
		t.Fatalf("write wasm: %v", err)
	}
	tmp.Close()

	out, err := execCommand("--json", "marketplace", "deploy", tmp.Name(), "alice")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	var resp struct {
		Address string `json:"address"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal deploy: %v", err)
	}
	if resp.Address == "" {
		t.Fatalf("expected address")
	}

	out, err = execCommand("--json", "marketplace", "trade", resp.Address, "bob")
	if err != nil {
		t.Fatalf("trade: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	var tradeResp map[string]string
	if err := json.Unmarshal([]byte(out), &tradeResp); err != nil {
		t.Fatalf("unmarshal trade: %v", err)
	}
	if tradeResp["status"] != "traded" {
		t.Fatalf("expected status traded, got %s", tradeResp["status"])
	}
}
