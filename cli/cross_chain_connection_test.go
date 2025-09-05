package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestConnectionManagerJSON(t *testing.T) {
	out, err := execCommand("cross_chain_connection", "open", "c1", "c2", "--json")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	id := resp["id"]
	out, err = execCommand("cross_chain_connection", "get", fmt.Sprintf("%v", id), "--json")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if !strings.Contains(out, "c1") {
		t.Fatalf("missing connection info: %s", out)
	}
}
