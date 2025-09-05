package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCustodialNodeJSON(t *testing.T) {
	out, err := execCommand("custodial", "custody", "user1", "5", "--json")
	if err != nil {
		t.Fatalf("custody failed: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	out, err = execCommand("custodial", "holdings", "user1", "--json")
	if err != nil {
		t.Fatalf("holdings failed: %v", err)
	}
	if !strings.Contains(out, "user1") {
		t.Fatalf("expected user in holdings: %s", out)
	}
}
