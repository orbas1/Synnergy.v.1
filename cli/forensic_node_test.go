package cli

import (
	"strings"
	"testing"
)

func TestForensicRecordTx(t *testing.T) {
	if _, err := execCommand("forensic", "record-tx", "--hash", "h1", "--from", "a", "--to", "b", "--value", "1"); err != nil {
		t.Fatalf("record-tx failed: %v", err)
	}
	out, err := execCommand("forensic", "txs", "--json")
	if err != nil {
		t.Fatalf("txs failed: %v", err)
	}
	if !strings.Contains(out, "h1") {
		t.Fatalf("expected transaction in output: %s", out)
	}
}
