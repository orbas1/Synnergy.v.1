package cli

import "testing"

// TestBSNAuthNoEnrollment verifies authentication fails when no template exists.
func TestBSNAuthNoEnrollment(t *testing.T) {
	out, err := execCommand("bsn", "auth", "addr1", "data", "00")
	if err != nil {
		t.Fatalf("auth failed: %v", err)
	}
	if out != "false" {
		t.Fatalf("expected false, got %q", out)
	}
}
