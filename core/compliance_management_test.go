package core

import "testing"

func TestComplianceManager(t *testing.T) {
	m := NewComplianceManager()
	if err := m.Suspend("addr1"); err != nil {
		t.Fatalf("suspend failed: %v", err)
	}
	if s, w := m.Status("addr1"); !s || w {
		t.Fatalf("unexpected status")
	}
	tx := Transaction{From: "addr1", To: "addr2", Amount: 1}
	if err := m.ReviewTransaction(tx); err == nil {
		t.Fatalf("expected review error")
	}
	if err := m.Whitelist("addr1"); err != nil {
		t.Fatalf("whitelist failed: %v", err)
	}
	if err := m.ReviewTransaction(tx); err != nil {
		t.Fatalf("review failed after whitelist: %v", err)
	}
	if err := m.Resume("addr1"); err != nil {
		t.Fatalf("resume failed: %v", err)
	}
	if s, _ := m.Status("addr1"); s {
		t.Fatalf("suspension not cleared")
	}

	if err := m.Unwhitelist("addr1"); err != nil {
		t.Fatalf("unwhitelist failed: %v", err)
	}

	if err := m.Suspend(""); err == nil {
		t.Fatalf("expected error for empty suspend")
	}
}
