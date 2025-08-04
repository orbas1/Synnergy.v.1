package core

import "testing"

func TestConnectionRegistry(t *testing.T) {
	reg := NewConnectionRegistry()
	c, err := reg.OpenConnection("chainA", "chainB")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := reg.CloseConnection(c.ID); err != nil {
		t.Fatalf("close: %v", err)
	}
	if c.Active {
		t.Fatalf("expected inactive connection")
	}
	if len(reg.ListConnections()) != 1 {
		t.Fatalf("list: expected 1 connection")
	}
}
