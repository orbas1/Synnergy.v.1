package cli

import (
	"testing"
)

func TestSyn3500Lifecycle(t *testing.T) {
	syn3500 = nil
	if _, err := execCommand("syn3500", "init", "--name", "Stable", "--symbol", "STB", "--issuer", "bank", "--rate", "1.0"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if _, err := execCommand("syn3500", "setrate", "1.1"); err != nil {
		t.Fatalf("setrate: %v", err)
	}
	if _, err := execCommand("syn3500", "mint", "alice", "100"); err != nil {
		t.Fatalf("mint: %v", err)
	}
	out, err := execCommand("syn3500", "balance", "alice")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if out != "100" {
		t.Fatalf("expected 100, got %s", out)
	}
	if _, err := execCommand("syn3500", "redeem", "alice", "40"); err != nil {
		t.Fatalf("redeem: %v", err)
	}
	out, err = execCommand("syn3500", "balance", "alice")
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if out != "60" {
		t.Fatalf("expected 60, got %s", out)
	}
}

func TestSyn3500Validation(t *testing.T) {
	syn3500 = nil
	if _, err := execCommand("syn3500", "init", "--name", "", "--symbol", "STB", "--issuer", "bank", "--rate", "1.0"); err == nil {
		t.Fatal("expected error for missing name")
	}
	if _, err := execCommand("syn3500", "init", "--name", "Stable", "--symbol", "STB", "--issuer", "bank", "--rate", "-1"); err == nil {
		t.Fatal("expected error for negative rate")
	}
	if _, err := execCommand("syn3500", "init", "--name", "Stable", "--symbol", "STB", "--issuer", "bank", "--rate", "1.0"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if _, err := execCommand("syn3500", "setrate", "0"); err == nil {
		t.Fatal("expected error for non-positive rate")
	}
}
