package cli

import (
	"strings"
	"testing"

	on "synnergy/internal/nodes/optimization_nodes"
)

func TestOptimizationFeeOrdering(t *testing.T) {
	feeOpt = on.FeeOptimizer{}

	out, err := executeCLICommand(t, "optimize", "fee", "tx1:10:2", "tx2:9:1")
	if err != nil {
		t.Fatalf("fee optimize: %v", err)
	}
	if !strings.Contains(out, "Hash:tx2") || !strings.Contains(out, "Hash:tx1") {
		t.Fatalf("unexpected optimize output: %q", out)
	}
	if strings.Index(out, "Hash:tx2") > strings.Index(out, "Hash:tx1") {
		t.Fatalf("expected tx2 before tx1, got %q", out)
	}
}
