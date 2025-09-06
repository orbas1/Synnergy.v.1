package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestShardingCLI verifies leader mapping commands.
func TestShardingCLI(t *testing.T) {
	if _, err := execCommand("sharding", "leader", "set", "1", "addr1", "--json"); err != nil {
		t.Fatalf("set leader: %v", err)
	}

	out, err := execCommand("sharding", "map", "--json")
	if err != nil {
		t.Fatalf("map: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}
	if m["1"] != "addr1" {
		t.Fatalf("unexpected leader map: %v", m)
	}
}
