package security

import "testing"

func validateSecurityAssessment(t *testing.T, filename string) {
	t.Helper()
	res, err := AnalyzeAssessment(filename)
	if err != nil {
		t.Fatalf("%s: %v", filename, err)
	}
	if res.Vulnerabilities == 0 || res.Mitigations == 0 {
		t.Fatalf("%s: expected non-zero vulnerabilities and mitigations, got %d and %d", filename, res.Vulnerabilities, res.Mitigations)
	}
}
