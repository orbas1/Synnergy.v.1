package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLiquidityViewsList(t *testing.T) {
	if _, err := execCommand("liquidity_pools", "create", "AAA", "BBB", "--json"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	out, err := execCommand("liquidity_views", "list", "--json")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var views []map[string]any
	if err := json.Unmarshal([]byte(out), &views); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(views) != 1 {
		t.Fatalf("expected one view, got %d", len(views))
	}
}
