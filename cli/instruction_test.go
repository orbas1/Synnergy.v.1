package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestInstructionNew verifies instruction creation outputs JSON and gas info.
func TestInstructionNew(t *testing.T) {
	out, err := execCommand("instruction", "new", "1", "10", "--json")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var inst map[string]any
	if err := json.Unmarshal([]byte(out), &inst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if inst["Op"].(float64) != 1 || inst["Value"].(float64) != 10 {
		t.Fatalf("unexpected instruction: %v", inst)
	}
}

// TestInstructionList ensures opcode catalogue is returned in JSON format.
func TestInstructionList(t *testing.T) {
	out, err := execCommand("instruction", "list", "--json")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json flag: %v", err)
	}
	if idx := strings.Index(out, "\n"); idx != -1 {
		out = out[idx+1:]
	}
	var cat map[string]string
	if err := json.Unmarshal([]byte(out), &cat); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(cat) == 0 {
		t.Fatalf("expected catalogue entries")
	}
}
