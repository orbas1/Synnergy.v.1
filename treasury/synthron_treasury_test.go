package treasury

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"synnergy/core"
)

func TestSynthronTreasuryLifecycle(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	before := treasury.Diagnostics(ctx)

	const issued = 500
	const burned = 200
	const transfer = 25

	if _, err := treasury.Issue(ctx, "stage80/lifecycle", issued); err != nil {
		t.Fatalf("issue: %v", err)
	}
	if err := treasury.Burn(ctx, "stage80/lifecycle", burned); err != nil {
		t.Fatalf("burn: %v", err)
	}
	if err := treasury.Transfer(ctx, nil, "stage80/liquidity", transfer, 0); err != nil {
		t.Fatalf("transfer: %v", err)
	}

	diag := treasury.Diagnostics(ctx)
	if diag.Minted < before.Minted+issued {
		t.Fatalf("minted delta too small: before=%d after=%d", before.Minted, diag.Minted)
	}
	if diag.Burned < before.Burned+burned {
		t.Fatalf("burned delta too small: before=%d after=%d", before.Burned, diag.Burned)
	}
	if bal := treasury.ledger.GetBalance("stage80/lifecycle"); bal != issued-burned {
		t.Fatalf("unexpected lifecycle balance: %d", bal)
	}
	if bal := treasury.ledger.GetBalance("stage80/liquidity"); bal < transfer {
		t.Fatalf("expected liquidity balance >= %d got %d", transfer, bal)
	}
	if diag.Circulating == 0 {
		t.Fatalf("expected non-zero circulating supply")
	}
	if diag.Health.VM == "" || diag.Health.Ledger == "" {
		t.Fatalf("expected health information to be populated: %+v", diag.Health)
	}
}

func TestSynthronTreasuryConsensusAndEvents(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	baseline := treasury.Diagnostics(ctx)
	events := treasury.SubscribeEvents()

	source := fmt.Sprintf("stage80-src-%d", time.Now().UnixNano())
	target := fmt.Sprintf("stage80-dst-%d", time.Now().UnixNano())
	if _, err := treasury.RegisterConsensusLink(ctx, source, target); err != nil {
		t.Fatalf("register consensus: %v", err)
	}
	addr := fmt.Sprintf("stage80-authority-%d", time.Now().UnixNano())
	if _, err := treasury.RegisterAuthority(ctx, addr, "stage80"); err != nil {
		t.Fatalf("register authority: %v", err)
	}

	select {
	case evt := <-events:
		if evt.Type == "" {
			t.Fatalf("empty event type")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("expected treasury event")
	}

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if _, err := treasury.Issue(ctx, fmt.Sprintf("stage80/parallel-%d", i), 10); err != nil {
				t.Errorf("issue %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	diag := treasury.Diagnostics(ctx)
	if diag.ConsensusNetworks <= baseline.ConsensusNetworks {
		t.Fatalf("expected consensus networks to grow: before=%d after=%d", baseline.ConsensusNetworks, diag.ConsensusNetworks)
	}
	if diag.AuthorityNodes <= baseline.AuthorityNodes {
		t.Fatalf("expected authority nodes to grow: before=%d after=%d", baseline.AuthorityNodes, diag.AuthorityNodes)
	}
	if len(diag.MissingOpcodes) != 0 {
		t.Fatalf("unexpected missing opcodes: %v", diag.MissingOpcodes)
	}
	if len(diag.AuditTrail) == 0 {
		t.Fatalf("expected audit trail to capture events")
	}
}

func TestSynthronTreasuryOperatorGovernance(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	unauthorized := WithOperator(ctx, "stage80/unauthorised")
	if _, err := treasury.Issue(unauthorized, "stage80/blocked", 10); !errors.Is(err, ErrUnauthorizedOperator) {
		t.Fatalf("expected unauthorized error got %v", err)
	}

	managedAddr := fmt.Sprintf("stage80-ops-%d", time.Now().UnixNano())
	if err := treasury.AuthorizeOperator(ctx, managedAddr); err != nil {
		t.Fatalf("authorize: %v", err)
	}

	diag := treasury.Diagnostics(ctx)
	found := false
	for _, op := range diag.Operators {
		if op == managedAddr {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected operators to include %s", managedAddr)
	}

	operatorCtx := WithOperator(ctx, managedAddr)
	if _, err := treasury.Issue(operatorCtx, "stage80/operations", 15); err != nil {
		t.Fatalf("issue by operator: %v", err)
	}
	if err := treasury.RevokeOperator(ctx, managedAddr); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := treasury.Issue(operatorCtx, "stage80/operations", 5); !errors.Is(err, ErrUnauthorizedOperator) {
		t.Fatalf("expected revoked operator to fail, got %v", err)
	}
	if err := treasury.RevokeOperator(ctx, diag.Wallet); !errors.Is(err, ErrProtectedOperator) {
		t.Fatalf("expected protected operator error got %v", err)
	}
}

func TestSynthronTreasurySelfHealingVM(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	if err := treasury.VirtualMachine().Stop(); err != nil {
		t.Fatalf("stop vm: %v", err)
	}
	if !treasury.VirtualMachine().Status() {
		if _, err := treasury.Issue(ctx, "stage80/selfheal", 1); err != nil {
			t.Fatalf("issue to restart vm: %v", err)
		}
	}
	if !treasury.VirtualMachine().Status() {
		t.Fatalf("expected vm to be running after ensure")
	}
}

func TestSynthronTreasurySummary(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	diag := treasury.Diagnostics(ctx)
	summary := SynthronTreasurySummary(diag)
	if summary == "" {
		t.Fatalf("expected non-empty summary")
	}
	if !strings.Contains(summary, "wallet:") || !strings.Contains(summary, "minted:") {
		t.Fatalf("summary missing key fields: %s", summary)
	}
}

func TestSynthronTreasuryAuditTrailSignatures(t *testing.T) {
	ctx := context.Background()
	treasury, err := DefaultSynthronTreasury(ctx)
	if err != nil {
		t.Fatalf("treasury init: %v", err)
	}

	baseline := len(treasury.AuditTrail())
	op := fmt.Sprintf("audit-operator-%d", time.Now().UnixNano())
	if err := treasury.AuthorizeOperator(ctx, op); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	opCtx := WithOperator(ctx, op)
	if _, err := treasury.Issue(opCtx, "audit/account", 5); err != nil {
		t.Fatalf("issue: %v", err)
	}
	diag := treasury.Diagnostics(ctx)
	if len(diag.AuditTrail) == 0 {
		t.Fatalf("expected audit trail entries")
	}

	if len(diag.AuditTrail) <= baseline {
		t.Fatalf("expected new events beyond baseline")
	}
	trail := diag.AuditTrail[baseline:]
	prevDigest := make([]byte, 0)
	if baseline > 0 {
		digestBytes, err := hex.DecodeString(diag.AuditTrail[baseline-1].Digest)
		if err != nil {
			t.Fatalf("decode baseline digest: %v", err)
		}
		prevDigest = digestBytes
	}
	for i, evt := range trail {
		if evt.Signature == "" || evt.Digest == "" {
			t.Fatalf("missing signature metadata on event %#v", evt)
		}
		payload := canonicalEventPayload(evt)
		base := []byte(payload)
		if len(prevDigest) > 0 {
			base = append(base, prevDigest...)
		}
		digest := sha256.Sum256(base)
		if hex.EncodeToString(digest[:]) != evt.Digest {
			t.Fatalf("digest mismatch for event %d", i)
		}
		sig, err := hex.DecodeString(evt.Signature)
		if err != nil {
			t.Fatalf("decode signature: %v", err)
		}
		if !core.VerifyMessage(base, sig, &treasury.wallet.PublicKey) {
			t.Fatalf("signature verification failed for event %d", i)
		}
		prevDigest = digest[:]
	}
}
