package cli

import "testing"

// TestPoolViews ensures the registry returns created pool views.
func TestPoolViews(t *testing.T) {
	poolRegistry.Create("X-Y", "X", "Y", 30)
	views := poolRegistry.PoolViews()
	if len(views) == 0 {
		t.Fatalf("expected views, got %d", len(views))
	}
}
