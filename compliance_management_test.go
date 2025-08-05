package synnergy

import "testing"

func TestComplianceManager(t *testing.T) {
	m := NewComplianceManager()
	m.Suspend("addr1")
	if s, w := m.Status("addr1"); !s || w {
		t.Fatalf("unexpected status")
	}
	tx := Transaction{From: "addr1", To: "addr2", Amount: 1}
	if err := m.ReviewTransaction(tx); err == nil {
		t.Fatalf("expected review error")
	}
	m.Whitelist("addr1")
	if err := m.ReviewTransaction(tx); err != nil {
		t.Fatalf("review failed after whitelist: %v", err)
	}
	m.Resume("addr1")
	if s, _ := m.Status("addr1"); s {
		t.Fatalf("suspension not cleared")
	}
}
