package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestLiquidityPoolsCreateList verifies creation and listing of pools with JSON output.
func TestLiquidityPoolsCreateList(t *testing.T) {
	out, err := execCommand("liquidity_pools", "create", "AAA", "BBB", "--json")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var res map[string]any
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	if res["status"] != "created" {
		t.Fatalf("unexpected create response: %v", res)
	}
	out, err = execCommand("liquidity_pools", "list", "--json")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var views []any
	if err := json.Unmarshal([]byte(out), &views); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if len(views) == 0 {
		t.Fatalf("expected pools, got %v", views)
	}
}
