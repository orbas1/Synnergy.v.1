package core

import "testing"

func TestNATManager(t *testing.T) {
	nm := NewNATManager()
	nm.MapPort("node1", 3030)
	if p, ok := nm.GetPort("node1"); !ok || p != 3030 {
		t.Fatalf("unexpected mapping")
	}
	nm.RemoveMapping("node1")
	if _, ok := nm.GetPort("node1"); ok {
		t.Fatalf("remove failed")
	}
}
