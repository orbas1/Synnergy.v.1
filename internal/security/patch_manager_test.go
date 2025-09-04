package security

import "testing"

func TestPatchManager(t *testing.T) {
	p := NewPatchManager()
	p.Apply("p1")
	applied := p.Applied()
	if len(applied) != 1 || applied[0] != "p1" {
		t.Fatalf("unexpected applied patches: %v", applied)
	}
}
