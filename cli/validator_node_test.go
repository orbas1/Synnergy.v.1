package cli

import (
	"strings"
	"testing"
)

func resetValidatorNode() { validatorNode = nil }

func TestValidatorNodeCreateQuorum(t *testing.T) {
	resetValidatorNode()
	out, err := execCommand("validatornode", "create", "--id", "n1", "--addr", "addr", "--minstake", "1", "--quorum", "1")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "created") {
		t.Fatalf("unexpected output: %s", out)
	}
	out, err = execCommand("validatornode", "quorum")
	if err != nil {
		t.Fatalf("quorum failed: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "false") {
		t.Fatalf("unexpected quorum output: %s", out)
	}
}
