package core

import "testing"

func TestConnectionManager(t *testing.T) {
	cm := NewConnectionManager()
	id := cm.OpenConnection("chainA", "chainB")
	if id == 0 {
		t.Fatalf("expected connection id")
	}
	if err := cm.CloseConnection(id); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	c, err := cm.GetConnection(id)
	if err != nil || c.Open {
		t.Fatalf("expected closed connection")
	}
	if len(cm.ListConnections()) != 1 {
		t.Fatalf("unexpected connection list length")
	}
}
