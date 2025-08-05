package core

import "testing"

func TestPlasmaManagement(t *testing.T) {
	b := NewPlasmaBridge()
	b.Pause()
	if !b.Status() {
		t.Fatalf("expected paused")
	}
	b.Resume()
	if b.Status() {
		t.Fatalf("expected running")
	}
}
