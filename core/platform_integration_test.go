package core

import (
	"context"
	"sync"
	"testing"
)

func TestPlatformIntegrationDiagnostics(t *testing.T) {
	integration, err := NewPlatformIntegration()
	if err != nil {
		t.Fatalf("integration: %v", err)
	}
	defer func() {
		if err := integration.Close(); err != nil {
			t.Fatalf("close: %v", err)
		}
	}()

	status := integration.Diagnostics(context.Background())
	if !status.VM.Running {
		t.Fatalf("expected vm running")
	}
	if status.Wallet.Address == "" {
		t.Fatalf("wallet address missing")
	}
	if status.Authority.Registered == 0 {
		t.Fatalf("expected authority registration")
	}
	if status.Consensus.Networks == 0 {
		t.Fatalf("expected consensus networks registered")
	}
	if status.Node.BlockHeight == 0 {
		t.Fatalf("expected mined diagnostic block")
	}
	if diag, ok := status.Diagnostics["overall"]; !ok || diag == "" {
		t.Fatalf("missing overall diagnostic summary")
	}
	if len(status.Issues) != 0 {
		t.Fatalf("unexpected issues reported: %v", status.Issues)
	}
	checks := map[string]DiagnosticCheck{
		"security":         status.Enterprise.Security,
		"scalability":      status.Enterprise.Scalability,
		"privacy":          status.Enterprise.Privacy,
		"governance":       status.Enterprise.Governance,
		"interoperability": status.Enterprise.Interoperability,
		"compliance":       status.Enterprise.Compliance,
	}
	for name, check := range checks {
		if !check.Healthy {
			t.Fatalf("expected %s check healthy: %+v", name, check)
		}
		if check.Detail == "" {
			t.Fatalf("missing detail for %s check", name)
		}
		if check.Latency == "" {
			t.Fatalf("missing latency for %s check", name)
		}
	}
}

func TestPlatformIntegrationDiagnosticsTimeout(t *testing.T) {
	integration, err := NewPlatformIntegration()
	if err != nil {
		t.Fatalf("integration: %v", err)
	}
	defer integration.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	status := integration.Diagnostics(ctx)
	if len(status.Issues) == 0 {
		t.Fatalf("expected issues when context cancelled")
	}
	if status.Diagnostics["overall"] != "integration cancelled" {
		t.Fatalf("expected cancelled status, got %s", status.Diagnostics["overall"])
	}
}

func TestPlatformIntegrationDiagnosticsConcurrent(t *testing.T) {
	integration, err := NewPlatformIntegration()
	if err != nil {
		t.Fatalf("integration: %v", err)
	}
	defer integration.Close()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			status := integration.Diagnostics(context.Background())
			if len(status.Issues) != 0 && status.Diagnostics["overall"] == "integration healthy" {
				t.Errorf("unexpected healthy status with issues: %+v", status.Issues)
			}
		}()
	}
	wg.Wait()
}
