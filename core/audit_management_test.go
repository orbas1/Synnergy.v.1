package core

import "testing"

func TestAuditManager_LogAndList(t *testing.T) {
	m := NewAuditManager()
	meta := map[string]string{"foo": "bar"}
	if err := m.Log("addr1", "event1", meta); err != nil {
		t.Fatalf("log failed: %v", err)
	}
	events := m.List("addr1")
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Event != "event1" {
		t.Fatalf("unexpected event: %s", events[0].Event)
	}
	if events[0].Metadata["foo"] != "bar" {
		t.Fatalf("metadata not stored")
	}
}
