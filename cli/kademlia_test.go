package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestKademliaStoreGet verifies storing and retrieving values with JSON output.
func TestKademliaStoreGet(t *testing.T) {
	out, err := execCommand("kademlia", "store", "aa", "hello", "--json")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var res map[string]any
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if res["key"] != "aa" || res["stored"] != true {
		t.Fatalf("unexpected response: %v", res)
	}

	out, err = execCommand("kademlia", "get", "aa", "--json")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var val map[string]any
	if err := json.Unmarshal([]byte(out), &val); err != nil {
		t.Fatalf("unmarshal get: %v", err)
	}
	if val["value"] != "hello" {
		t.Fatalf("expected hello, got %v", val)
	}
}

// TestKademliaClosest ensures closest command returns JSON array.
func TestKademliaClosest(t *testing.T) {
	execCommand("kademlia", "store", "a1", "1")
	execCommand("kademlia", "store", "b2", "2")
	out, err := execCommand("kademlia", "closest", "a1", "2", "--json")
	if err != nil {
		t.Fatalf("closest: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var arr []string
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(arr) != 2 {
		t.Fatalf("expected two keys, got %v", arr)
	}
}
