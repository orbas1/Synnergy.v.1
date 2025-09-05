package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestProtocolRegistryJSON(t *testing.T) {
	out, err := execCommand("cross_chain_agnostic_protocols", "register", "proto1", "--json")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if _, ok := resp["id"]; !ok {
		t.Fatalf("missing id: %v", resp)
	}
	out, err = execCommand("cross_chain_agnostic_protocols", "list", "--json")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "proto1") {
		t.Fatalf("expected protocol in list: %s", out)
	}
}
