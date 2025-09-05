package cli

import (
	"strings"
	"testing"
)

// TestBSNAuthNoEnrollment verifies authentication fails when no template exists.
func TestBSNAuthNoEnrollment(t *testing.T) {
	out, err := execCommand("bsn", "--json", "auth", "addr1", "data", "00")
	if err != nil {
		t.Fatalf("auth failed: %v", err)
	}
	if !strings.Contains(out, "\"authenticated\":false") {
		t.Fatalf("unexpected output: %s", out)
	}
}
