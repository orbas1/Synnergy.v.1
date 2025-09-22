//go:build experimental

package cli

import (
	"strings"
	"testing"
)

func TestExperimentalNodeLifecycle(t *testing.T) {
	expNode = nil
	t.Cleanup(func() { expNode = nil })

	out, err := executeCLICommand(t, "experimental", "status")
	if err != nil {
		t.Fatalf("status without node: %v", err)
	}
	if strings.TrimSpace(out) != "no node" {
		t.Fatalf("expected 'no node', got %q", out)
	}

	out, err = executeCLICommand(t, "experimental", "create", "exp1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if strings.TrimSpace(out) != "created" {
		t.Fatalf("unexpected create output: %q", out)
	}
	if expNode == nil {
		t.Fatalf("expected expNode to be created")
	}

	if _, err := executeCLICommand(t, "experimental", "start"); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !expNode.IsRunning() {
		t.Fatalf("node should be running")
	}

	if _, err := executeCLICommand(t, "experimental", "dial", "peer1"); err != nil {
		t.Fatalf("dial: %v", err)
	}

	if _, err := executeCLICommand(t, "experimental", "stop"); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if expNode.IsRunning() {
		t.Fatalf("node should be stopped")
	}
}
