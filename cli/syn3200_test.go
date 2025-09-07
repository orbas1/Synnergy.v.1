package cli

import (
	"strings"
	"testing"
	"time"

	"synnergy/core"
)

func TestSyn3200Lifecycle(t *testing.T) {
	bills = core.NewBillRegistry()
	due := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	if _, err := execCommand("syn3200", "create", "--id", "b1", "--issuer", "iss", "--payer", "pay", "--amount", "100", "--due", due, "--meta", "m"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := execCommand("syn3200", "pay", "b1", "pay", "40"); err != nil {
		t.Fatalf("pay: %v", err)
	}
	out, err := execCommand("syn3200", "get", "b1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !strings.Contains(out, "Amount:60") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestSyn3200InvalidDue(t *testing.T) {
	bills = core.NewBillRegistry()
	if _, err := execCommand("syn3200", "create", "--id", "b1", "--issuer", "iss", "--payer", "pay", "--amount", "100", "--due", "bad", "--meta", "m"); err == nil {
		t.Fatal("expected error")
	}
}
