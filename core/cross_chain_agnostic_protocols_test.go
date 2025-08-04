package core

import "testing"

func TestProtocolRegistry(t *testing.T) {
	r := NewProtocolRegistry()
	id := r.Register("ics-20")
	if id == 0 {
		t.Fatalf("expected id > 0")
	}
	if _, ok := r.Get(id); !ok {
		t.Fatalf("protocol not found")
	}
	list := r.List()
	if len(list) != 1 || list[0].ID != id {
		t.Fatalf("unexpected list result")
	}
}
