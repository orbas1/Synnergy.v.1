package cli

import (
	"encoding/json"
	"testing"
)

func TestLiquidityViewsList(t *testing.T) {
	if _, err := execCommand("liquidity_pools", "create", "AAA", "BBB"); err != nil {
		t.Fatalf("create: %v", err)
	}
	out, err := execCommand("liquidity_views", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var views []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &views); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(views) != 1 {
		t.Fatalf("expected one view, got %d", len(views))
	}
}
