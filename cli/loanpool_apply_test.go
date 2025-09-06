package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestLoanpoolApplySubmitList checks application submission and listing via JSON output.
func TestLoanpoolApplySubmitList(t *testing.T) {
	out, err := execCommand("loanpool_apply", "submit", "bob", "100", "12", "car", "--json")
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var res map[string]any
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("unmarshal submit: %v", err)
	}
	out, err = execCommand("loanpool_apply", "list", "--json")
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
		t.Fatalf("expected applications, got %v", views)
	}
	_ = res // ensure submit succeeded
}
