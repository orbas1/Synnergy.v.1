package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestXContractJSON(t *testing.T) {
	if _, err := execCommand("xcontract", "register", "local", "remote", "raddr", "--json"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	out, err := execCommand("xcontract", "list", "--json")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "local") {
		t.Fatalf("mapping missing: %s", out)
	}
	out, err = execCommand("xcontract", "remove", "local", "--json")
	if err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if _, ok := resp["gas"]; !ok {
		t.Fatalf("expected gas field: %v", resp)
	}
}
