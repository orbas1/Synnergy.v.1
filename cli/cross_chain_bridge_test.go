package cli

import (
	"encoding/json"
	"testing"
)

func TestCrossChainBridgeDepositJSON(t *testing.T) {
	out, err := execCommand("cross_chain_bridge", "deposit", "b1", "alice", "bob", "1", "--json")
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["id"] == "" {
		t.Fatalf("expected id in response: %v", resp)
	}
}
