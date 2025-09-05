package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestConsensusNetworkJSON(t *testing.T) {
	out, err := execCommand("cross-consensus", "register", "src", "dst", "--json")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	out, err = execCommand("cross-consensus", "list", "--json")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "src") {
		t.Fatalf("network not listed: %s", out)
	}
}
