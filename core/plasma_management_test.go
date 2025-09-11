package core

import "testing"

func TestPlasmaManagement(t *testing.T) {
	b := NewPlasmaBridge()
	b.Pause()
	if !b.IsPaused() {
		t.Fatalf("expected paused")
	}
	b.Resume()
	if b.IsPaused() {
		t.Fatalf("expected running")
	}
}
