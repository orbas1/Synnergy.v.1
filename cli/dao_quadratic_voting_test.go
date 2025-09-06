package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"synnergy/core"
)

// TestDAOQuadraticWeightJSON verifies JSON output for weight calculation.
func TestDAOQuadraticWeightJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dao-qv", "weight", "9", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	var resp map[string]uint64
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	expected := core.QuadraticWeight(9)
	if resp["weight"] != expected {
		t.Fatalf("unexpected weight: %v", resp)
	}
}
