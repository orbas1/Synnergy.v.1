package cli

import (
	"strings"
	"testing"
)

// TestContractOpcodesGas ensures opcode list includes gas costs.
func TestContractOpcodesGas(t *testing.T) {
	out, err := execCommand("contractopcodes")
	if err != nil {
		t.Fatalf("contractopcodes failed: %v", err)
	}
	if !strings.Contains(out, "gas") {
		t.Fatalf("expected gas annotation, got %q", out)
	}
}
