package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStage73StoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stage73_state.json")

	store := NewStage73Store(path)
	if err := store.Load(); err != nil {
		t.Fatalf("load fresh store: %v", err)
	}

	idx := NewSYN3700Token("Stage73 Index", "S73")
	idx.AddController("controller-1")
	if err := idx.AddComponent("controller-1", "AAA", 0.6, 0.1); err != nil {
		t.Fatalf("add component: %v", err)
	}
	store.SetIndex(idx)

	grants := store.Grants()
	grantID := grants.CreateGrant("beneficiary", "education", 100)
	if err := grants.AddAuthorizer(grantID, "authorizer-1"); err != nil {
		t.Fatalf("add authorizer: %v", err)
	}
	if err := grants.DisburseWithActor(grantID, 40, "first", "authorizer-1"); err != nil {
		t.Fatalf("disburse: %v", err)
	}

	benefits := store.Benefits()
	benefitID := benefits.RegisterBenefit("recipient", "health", 75)
	if err := benefits.AddApprover(benefitID, "approver-1"); err != nil {
		t.Fatalf("benefit approver: %v", err)
	}
	if err := benefits.Claim(benefitID, "recipient"); err != nil {
		t.Fatalf("benefit claim: %v", err)
	}

	charity := store.Charity()
	charity.Donate("HELP", "donor", 25, "support")

	legal := store.Legal()
	token := NewLegalToken("case-1", "Contract", "LGL", "nda", "hash", "owner", time.Now().Add(24*time.Hour), 10, []string{"a", "b"})
	legal.Add(token)
	if err := token.Sign("a", "sig-a"); err != nil {
		t.Fatalf("sign token: %v", err)
	}

	util := NewSYN500Token("Utility", "UTL", "issuer", 4, 1000)
	util.Grant("consumer", 1, 10, time.Minute)
	if err := util.Use("consumer"); err != nil {
		t.Fatalf("use utility: %v", err)
	}
	store.SetUtility(util)

	store.MarkDirty()
	if err := store.Save(); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected snapshot data")
	}

	reloaded := NewStage73Store(path)
	if err := reloaded.Load(); err != nil {
		t.Fatalf("reload snapshot: %v", err)
	}

	snapshot := reloaded.Snapshot()
	if snapshot.Index == nil || snapshot.Index.Name != "Stage73 Index" {
		t.Fatalf("missing index snapshot: %+v", snapshot.Index)
	}
	if snapshot.Grants.Summary.Total != 1 || snapshot.Grants.Summary.Completed != 0 {
		t.Fatalf("unexpected grant summary: %+v", snapshot.Grants.Summary)
	}
	if len(snapshot.Benefits.Records) != 1 || snapshot.Benefits.Summary.Claimed != 1 {
		t.Fatalf("unexpected benefit snapshot: %+v", snapshot.Benefits.Summary)
	}
	if len(snapshot.Charity) != 1 || snapshot.Charity[0].Raised != 25 {
		t.Fatalf("unexpected charity snapshot: %+v", snapshot.Charity)
	}
	if len(snapshot.Legal) != 1 || snapshot.Legal[0].Signatures["a"] != "sig-a" {
		t.Fatalf("unexpected legal snapshot: %+v", snapshot.Legal)
	}
	if snapshot.Utility == nil || snapshot.Utility.Telemetry.Grants != 1 {
		t.Fatalf("unexpected utility snapshot: %+v", snapshot.Utility)
	}

	digest, err := reloaded.Digest()
	if err != nil {
		t.Fatalf("digest error: %v", err)
	}
	if digest == "" {
		t.Fatalf("expected digest value")
	}
}

func TestStage73StoreNoPath(t *testing.T) {
	store := NewStage73Store("")
	if err := store.Load(); err != nil {
		t.Fatalf("load in-memory: %v", err)
	}
	store.MarkDirty()
	if err := store.Save(); err != nil {
		t.Fatalf("save in-memory: %v", err)
	}
}

func TestStage73StoreConcurrentSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stage73_state.json")

	store := NewStage73Store(path)
	if err := store.Load(); err != nil {
		t.Fatalf("load fresh store: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			registry := store.Grants()
			registry.CreateGrant(fmt.Sprintf("beneficiary-%d", i), fmt.Sprintf("program-%d", i), uint64(100+i))
			store.MarkDirty()
			if err := store.Save(); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("save failure: %v", err)
		}
	}

	if got := len(store.Grants().ListGrants()); got != 8 {
		t.Fatalf("expected 8 grants in memory, got %d", got)
	}
	store.MarkDirty()
	if err := store.Save(); err != nil {
		t.Fatalf("final save: %v", err)
	}

	reloaded := NewStage73Store(path)
	if err := reloaded.Load(); err != nil {
		t.Fatalf("reload snapshot: %v", err)
	}
	snap := reloaded.Snapshot()
	if snap.Grants.Summary.Total != 8 {
		t.Fatalf("expected 8 grants, got %d", snap.Grants.Summary.Total)
	}
	if digest, err := reloaded.Digest(); err != nil || digest == "" {
		t.Fatalf("digest mismatch: %v %q", err, digest)
	}
}

func TestStage73StoreCorruptSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")
	if err := os.WriteFile(path, []byte("{\"index\":"), 0o600); err != nil {
		t.Fatalf("write corrupt snapshot: %v", err)
	}
	store := NewStage73Store(path)
	err := store.Load()
	if err == nil {
		t.Fatal("expected load error for corrupt snapshot")
	}
	if !strings.Contains(err.Error(), "decode snapshot") {
		t.Fatalf("unexpected error: %v", err)
	}
}
