package core

import "testing"

func TestQuadraticWeight(t *testing.T) {
	if w := QuadraticWeight(9); w != 3 {
		t.Fatalf("expected 3 got %d", w)
	}
}
