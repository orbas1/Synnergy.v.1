package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestGovernmentNew ensures a government authority node can be created.
func TestGovernmentNew(t *testing.T) {
	out, err := execCommand("government", "new", "addr1", "role1", "dept1", "--json")
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var resp map[string]string
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["address"] != "addr1" || resp["role"] != "role1" || resp["department"] != "dept1" {
		t.Fatalf("unexpected response: %v", resp)
	}
}
