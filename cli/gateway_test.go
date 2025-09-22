package cli

import (
	"strings"
	"testing"

	"synnergy/core"
	nodes "synnergy/internal/nodes"
)

func TestGatewayCommands(t *testing.T) {
	gateway = core.NewGatewayNode(nodes.Address("gw-test"), core.GatewayConfig{})
	if err := gateway.Start(); err != nil {
		t.Fatalf("start gateway: %v", err)
	}
	t.Cleanup(func() {
		_ = gateway.Stop()
		gateway = core.NewGatewayNode(nodes.Address("gw1"), core.GatewayConfig{})
		_ = gateway.Start()
	})

	if _, err := executeCLICommand(t, "gateway", "register", "foo"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if endpoints := gateway.Endpoints(); len(endpoints) != 1 || endpoints[0] != "foo" {
		t.Fatalf("unexpected endpoints: %v", endpoints)
	}

	out, err := executeCLICommand(t, "gateway", "call", "foo", "payload")
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if strings.TrimSpace(out) != "foo received: payload" {
		t.Fatalf("unexpected call output: %q", out)
	}

	out, err = executeCLICommand(t, "gateway", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if strings.TrimSpace(out) != "foo" {
		t.Fatalf("unexpected list output: %q", out)
	}

	if _, err := executeCLICommand(t, "gateway", "remove", "foo"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if len(gateway.Endpoints()) != 0 {
		t.Fatalf("expected no endpoints after removal")
	}
}
