package cli

import (
	"encoding/json"
	"strings"
	"testing"

	nodes "synnergy/internal/nodes"
)

// TestHolographicNodeWorkflow verifies holographic node CLI commands emit JSON output.
func TestHolographicNodeWorkflow(t *testing.T) {
	out, err := execCommand("holographic", "store", "id1", "data", "2", "--json")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var storeResp struct {
		Status string `json:"status"`
		ID     string `json:"id"`
		Shards int    `json:"shards"`
	}
	if err := json.Unmarshal([]byte(out), &storeResp); err != nil {
		t.Fatalf("unmarshal store: %v", err)
	}
	if storeResp.Status != "stored" || storeResp.ID != "id1" || storeResp.Shards != 2 {
		t.Fatalf("unexpected store response: %+v", storeResp)
	}

	out, err = execCommand("holographic", "retrieve", "id1", "--json")
	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var retResp struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &retResp); err != nil {
		t.Fatalf("unmarshal retrieve: %v", err)
	}
	if retResp.ID != "id1" || retResp.Data != "data" {
		t.Fatalf("unexpected retrieve response: %+v", retResp)
	}

	if _, err := execCommand("holographic", "dial", "peer1", "--json"); err != nil {
		t.Fatalf("dial: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}

	out, err = execCommand("holographic", "peers", "--json")
	if err != nil {
		t.Fatalf("peers: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var peersResp struct {
		Peers []string `json:"peers"`
	}
	if err := json.Unmarshal([]byte(out), &peersResp); err != nil {
		t.Fatalf("unmarshal peers: %v", err)
	}
	if len(peersResp.Peers) != 1 || peersResp.Peers[0] != string(nodes.Address("peer1")) {
		t.Fatalf("unexpected peers response: %+v", peersResp)
	}

	t.Cleanup(func() { holoNode = nodes.NewHolographicNode(nodes.Address("holo-1")) })
}
