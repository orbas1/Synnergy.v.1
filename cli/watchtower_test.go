package cli

import "testing"

func TestWatchtowerLifecycle(t *testing.T) {
	if err := watchtowerNode.Start(watchtowerCtx); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := watchtowerNode.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
