package cli

import "testing"

// TestBiometricVerifyNoEnrollment ensures the verify command returns false when the user is not enrolled.
func TestBiometricVerifyNoEnrollment(t *testing.T) {
	out, err := execCommand("biometric", "verify", "user1", "data", "00")
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if out != "false" {
		t.Fatalf("expected false, got %q", out)
	}
}
