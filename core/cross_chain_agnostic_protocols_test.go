package core

import "testing"

func TestProtocolRegistry(t *testing.T) {
	r := NewProtocolRegistry()

	if _, err := r.Register("ics-20", "relayer1"); err == nil {
		t.Fatalf("expected unauthorized relayer error")
	}

	r.AuthorizeRelayer("relayer1")
	id, err := r.Register("ics-20", "relayer1")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if _, ok := r.Get(id); !ok {
		t.Fatalf("protocol not found")
	}

	list := r.List()
	if len(list) != 1 || list[0].ID != id {
		t.Fatalf("unexpected list result")
	}

	if err := r.Remove(id, "bad"); err == nil {
		t.Fatalf("expected unauthorized removal error")
	}
	if err := r.Remove(id, "relayer1"); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if _, ok := r.Get(id); ok {
		t.Fatalf("expected protocol to be removed")
	}
}
