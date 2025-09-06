package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"synnergy/core"
)

// TestHistoricalWorkflow ensures historical node CLI commands emit JSON output.
func TestHistoricalWorkflow(t *testing.T) {
	out, err := execCommand("historical", "archive", "1", "h1", "--json")
	if err != nil {
		t.Fatalf("archive: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var archResp struct {
		Status string `json:"status"`
		Height uint64 `json:"height"`
		Hash   string `json:"hash"`
	}
	if err := json.Unmarshal([]byte(out), &archResp); err != nil {
		t.Fatalf("unmarshal archive: %v", err)
	}
	if archResp.Status != "archived" || archResp.Height != 1 || archResp.Hash != "h1" {
		t.Fatalf("unexpected archive response: %+v", archResp)
	}

	out, err = execCommand("historical", "height", "1", "--json")
	if err != nil {
		t.Fatalf("height: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var blockResp struct {
		Height    uint64 `json:"height"`
		Hash      string `json:"hash"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal([]byte(out), &blockResp); err != nil {
		t.Fatalf("unmarshal height: %v", err)
	}
	if blockResp.Height != 1 || blockResp.Hash != "h1" {
		t.Fatalf("unexpected height response: %+v", blockResp)
	}

	out, err = execCommand("historical", "hash", "h1", "--json")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var hashResp struct {
		Height    uint64 `json:"height"`
		Hash      string `json:"hash"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal([]byte(out), &hashResp); err != nil {
		t.Fatalf("unmarshal hash: %v", err)
	}
	if hashResp.Height != 1 || hashResp.Hash != "h1" {
		t.Fatalf("unexpected hash response: %+v", hashResp)
	}

	out, err = execCommand("historical", "total", "--json")
	if err != nil {
		t.Fatalf("total: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var totalResp struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal([]byte(out), &totalResp); err != nil {
		t.Fatalf("unmarshal total: %v", err)
	}
	if totalResp.Total != 1 {
		t.Fatalf("unexpected total response: %+v", totalResp)
	}

	t.Cleanup(func() { historicalNode = core.NewHistoricalNode() })
}
