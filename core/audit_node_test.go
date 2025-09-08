package core

import "testing"

type mockBootstrap struct{ started bool }

func (m *mockBootstrap) Start() error {
	m.started = true
	return nil
}

func TestAuditNode_StartLogAndVerify(t *testing.T) {
	mgr := NewAuditManager()
	bs := &mockBootstrap{}
	node := NewAuditNode(bs, mgr)
	if err := node.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if !bs.started {
		t.Fatalf("bootstrap not started")
	}
	if err := node.LogEvent("addr", "evt", nil); err != nil {
		t.Fatalf("log event failed: %v", err)
	}
	if err := node.LogEvent("", "evt", nil); err == nil {
		t.Fatalf("expected validation error")
	}
	events := node.ListEvents("addr")
	if len(events) != 1 || events[0].Event != "evt" {
		t.Fatalf("event not recorded")
	}
	ok, err := node.VerifyEvent("addr", 0)
	if err != nil || !ok {
		t.Fatalf("verify failed: %v", err)
	}
	if _, err := node.VerifyEvent("addr", 1); err == nil {
		t.Fatalf("expected index error")
	}
	// manager missing
	badNode := NewAuditNode(bs, nil)
	if _, err := badNode.VerifyEvent("addr", 0); err == nil {
		t.Fatalf("expected manager error")
	}
}
