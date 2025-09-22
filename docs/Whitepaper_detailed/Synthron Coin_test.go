package whitepaperdetailed

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSynthronTreasuryUnitAndFunctional(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronCoinTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	diag := treasury.Diagnostics(ctx)
	mintedBefore := diag.Minted

	if _, err := treasury.Issue(ctx, "enterprise/ops", 1_000); err != nil {
		t.Fatalf("issue: %v", err)
	}
	if err := treasury.Burn(ctx, "enterprise/ops", 250); err != nil {
		t.Fatalf("burn: %v", err)
	}

	diag = treasury.Diagnostics(ctx)
	if diag.Minted <= mintedBefore {
		t.Fatalf("minted not updated: before=%d after=%d", mintedBefore, diag.Minted)
	}
	if diag.Burned < 250 {
		t.Fatalf("expected burned >= 250 got %d", diag.Burned)
	}

	if err := treasury.Transfer(ctx, nil, "operations/liquidity", 100, 0); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	diag = treasury.Diagnostics(ctx)
	if diag.Circulating == 0 {
		t.Fatalf("expected positive circulating supply")
	}
	if len(diag.AuditTrail) == 0 {
		t.Fatalf("expected audit trail entries after lifecycle operations")
	}
	trail, err := TreasuryAuditTrail(ctx)
	if err != nil {
		t.Fatalf("audit trail helper: %v", err)
	}
	if len(trail) == 0 {
		t.Fatalf("expected treasury audit helper to return entries")
	}
}

func TestSynthronTreasurySituationalStressAndEvents(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronCoinTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	if _, err := treasury.RegisterConsensusLink(ctx, "pos", "poh"); err != nil {
		t.Fatalf("register consensus: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			addr := fmt.Sprintf("stress-%d", i)
			if _, err := treasury.Issue(ctx, addr, 25); err != nil {
				t.Errorf("issue stress %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	diag := treasury.Diagnostics(ctx)
	if diag.ConsensusNetworks < 2 {
		t.Fatalf("expected at least 2 consensus networks got %d", diag.ConsensusNetworks)
	}
	if diag.AuthorityNodes == 0 {
		t.Fatalf("expected authority nodes to be registered")
	}
	if len(diag.MissingOpcodes) != 0 {
		t.Fatalf("unexpected missing opcodes: %v", diag.MissingOpcodes)
	}
	if len(diag.AuditTrail) == 0 {
		t.Fatalf("expected audit telemetry to be populated")
	}

	events := treasury.SubscribeEvents()
	if _, err := treasury.RegisterAuthority(ctx, "compliance/control", "compliance"); err != nil {
		t.Fatalf("register authority: %v", err)
	}
	select {
	case evt := <-events:
		if evt.Type == "" {
			t.Fatalf("empty event received")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("expected authority event")
	}
}

func TestSynthronTreasuryOperatorDocumentationFlows(t *testing.T) {
	ctx := context.Background()
	const operator = "docs/operator"
	if err := TreasuryAuthorizeOperator(ctx, operator); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	opCtx := TreasuryWithOperator(ctx, operator)
	treasury, err := DefaultSynthronCoinTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}
	if _, err := treasury.Issue(opCtx, "docs/payout", 42); err != nil {
		t.Fatalf("issue via operator: %v", err)
	}
	if err := TreasuryRevokeOperator(ctx, operator); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := treasury.Issue(opCtx, "docs/payout", 1); err == nil {
		t.Fatalf("expected revoked operator to fail")
	}
}
