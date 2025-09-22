package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runCommand(t *testing.T, name string, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func ensureOpenSSL(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("openssl not available: " + err.Error())
	}
}

func generateKeyAndCert(t *testing.T) (string, string) {
	t.Helper()
	ensureOpenSSL(t)

	dir := t.TempDir()
	keyPath := filepath.Join(dir, "node.key")
	certPath := filepath.Join(dir, "node.pem")

	cmd := exec.Command("openssl", "genrsa", "-out", keyPath, "2048")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to generate key: %v (%s)", err, out)
	}

	cmd = exec.Command("openssl", "req", "-x509", "-new", "-key", keyPath, "-out", certPath, "-days", "1", "-subj", "/CN=test.synnergy")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to generate certificate: %v (%s)", err, out)
	}

	return keyPath, certPath
}

func TestCertificateRenewUsage(t *testing.T) {
	out, err := runCommand(t, "bash", "../../scripts/certificate_renew.sh")
	if err == nil {
		t.Fatalf("expected error when no arguments provided")
	}
	if !strings.Contains(string(out), "Usage: certificate_renew.sh") {
		t.Fatalf("expected usage output, got: %s", out)
	}
}

func TestCertificateRenewMissingCertificate(t *testing.T) {
	ensureOpenSSL(t)
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "node.key")

	cmd := exec.Command("openssl", "genrsa", "-out", keyPath, "2048")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to generate key: %v (%s)", err, out)
	}

	missingCert := filepath.Join(dir, "missing.pem")
	out, err := runCommand(t, "bash", "../../scripts/certificate_renew.sh",
		"--cert", missingCert,
		"--key", keyPath,
		"--output-dir", filepath.Join(dir, "out"),
		"--dry-run",
		"--force",
	)
	if err == nil {
		t.Fatalf("expected failure when certificate is missing")
	}
	if !strings.Contains(string(out), "certificate '") {
		t.Fatalf("expected certificate missing message, got: %s", out)
	}
}

func TestCertificateRenewMissingCLI(t *testing.T) {
	keyPath, certPath := generateKeyAndCert(t)
	outDir := t.TempDir()

	out, err := runCommand(t, "bash", "../../scripts/certificate_renew.sh",
		"--cert", certPath,
		"--key", keyPath,
		"--output-dir", outDir,
		"--force",
		"--cli", "/nonexistent/synnergy",
	)
	if err == nil {
		t.Fatalf("expected failure when CLI binary is missing")
	}
	if !strings.Contains(string(out), "CLI binary") {
		t.Fatalf("expected CLI missing message, got: %s", out)
	}
}

func TestCertificateRenewDryRunSuccess(t *testing.T) {
	keyPath, certPath := generateKeyAndCert(t)
	outDir := t.TempDir()
	metricsPath := filepath.Join(outDir, "metrics", "renewal.json")

	out, err := runCommand(t, "bash", "../../scripts/certificate_renew.sh",
		"--cert", certPath,
		"--key", keyPath,
		"--output-dir", outDir,
		"--dry-run",
		"--force",
		"--metrics-file", metricsPath,
	)
	if err != nil {
		t.Fatalf("expected dry-run to succeed, got error: %v (%s)", err, out)
	}

	lines := strings.Split(string(out), "\n")
	var csrPath string
	for _, line := range lines {
		if strings.HasPrefix(line, "CSR_PATH=") {
			csrPath = strings.TrimPrefix(line, "CSR_PATH=")
			break
		}
	}
	if csrPath == "" {
		t.Fatalf("expected CSR_PATH in output, got: %s", out)
	}

	data, err := os.ReadFile(csrPath)
	if err != nil {
		t.Fatalf("expected CSR file, failed to read: %v", err)
	}
	if !strings.Contains(string(data), "BEGIN CERTIFICATE REQUEST") {
		t.Fatalf("unexpected CSR content: %s", data)
	}

	metricsData, err := os.ReadFile(metricsPath)
	if err != nil {
		t.Fatalf("expected metrics file, failed to read: %v", err)
	}
	if !strings.Contains(string(metricsData), "\"dry_run\": 1") {
		t.Fatalf("expected dry_run flag in metrics, got: %s", metricsData)
	}
}
