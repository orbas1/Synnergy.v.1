package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPatchManagerApplyMetadata(t *testing.T) {
	t.Helper()
	mgr := NewPatchManager()
	now := time.Now().UTC()
	meta := PatchMetadata{
		ID:          "p1",
		Version:     "1.2.3",
		Description: "validator hotfix",
		Digest:      []byte("digest"),
		Metadata:    map[string]string{"component": "validator"},
		SubmittedAt: now,
		AppliedAt:   now.Add(time.Second),
		ApprovedBy:  "alice",
	}
	if err := mgr.ApplyMetadata(meta); err != nil {
		t.Fatalf("apply metadata: %v", err)
	}
	if err := mgr.ApplyMetadata(meta); err == nil || !strings.Contains(err.Error(), "already recorded") {
		t.Fatalf("expected duplicate patch error, got %v", err)
	}
	got := mgr.Metadata()
	if len(got) != 1 || got[0].ID != "p1" || got[0].Metadata["component"] != "validator" {
		t.Fatalf("unexpected metadata snapshot: %+v", got)
	}
	applied := mgr.Applied()
	if len(applied) != 1 || applied[0] != "p1" {
		t.Fatalf("unexpected applied order: %v", applied)
	}
}

func TestPatchManagerValidatorAndSignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	mgr := NewPatchManager(pub)
	digest := sha256.Sum256([]byte("payload"))
	meta := PatchMetadata{
		ID:          "p2",
		Version:     "1.0.1",
		Description: "consensus tuning",
		Digest:      digest[:],
		Metadata:    map[string]string{"rollout": "staged"},
	}
	meta.Signature = ed25519.Sign(priv, patchDigestPayload(meta))
	validateCount := 0
	mgr.SetValidator(func(m PatchMetadata) error {
		validateCount++
		if m.Metadata["rollout"] != "staged" {
			return errors.New("unexpected metadata")
		}
		return nil
	})
	if err := mgr.ApplyMetadata(meta); err != nil {
		t.Fatalf("apply metadata: %v", err)
	}
	if validateCount != 1 {
		t.Fatalf("expected validator invoked once got %d", validateCount)
	}
	// Modify signature and ensure verification fails.
	meta.ID = "p3"
	meta.Signature = append([]byte(nil), meta.Signature...)
	meta.Signature[0] ^= 0xff
	if err := mgr.ApplyMetadata(meta); err == nil {
		t.Fatalf("expected signature validation error")
	}
}

func TestPatchManagerConcurrentAccess(t *testing.T) {
	mgr := NewPatchManager()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := "p" + string(rune('a'+i))
			digest := sha256.Sum256([]byte(id))
			err := mgr.ApplyMetadata(PatchMetadata{ID: id, Digest: digest[:], Metadata: map[string]string{"idx": id}})
			if err != nil {
				t.Errorf("apply %s: %v", id, err)
			}
		}(i)
	}
	wg.Wait()
	metas := mgr.Metadata()
	if len(metas) != 10 {
		t.Fatalf("expected 10 patches got %d", len(metas))
	}
	// Ensure deterministic ordering by AppliedAt then ID.
	last := metas[0].AppliedAt
	for _, m := range metas {
		if m.AppliedAt.Before(last) {
			t.Fatalf("metadata not sorted by AppliedAt: %v < %v", m.AppliedAt, last)
		}
		last = m.AppliedAt
	}
}
