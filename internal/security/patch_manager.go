package security

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// PatchMetadata captures everything required to validate and audit a software
// update applied to the network.
type PatchMetadata struct {
	ID          string
	Version     string
	Description string
	Digest      []byte
	Signature   []byte
	SubmittedAt time.Time
	AppliedAt   time.Time
	ApprovedBy  string
	Metadata    map[string]string
}

// PatchValidator validates metadata before it is accepted. Validators typically
// verify digital signatures or cross-check governance approvals.
type PatchValidator func(PatchMetadata) error

// PatchManager stores patch metadata and exposes helpers for the CLI and web UI.
type PatchManager struct {
	mu        sync.RWMutex
	patches   map[string]PatchMetadata
	order     []string
	validator PatchValidator
	authority ed25519.PublicKey
}

// NewPatchManager creates a manager. When an authority key is supplied the
// manager verifies patch signatures with the provided public key.
func NewPatchManager(authority ...ed25519.PublicKey) *PatchManager {
	var key ed25519.PublicKey
	if len(authority) > 0 {
		key = authority[0]
	}
	return &PatchManager{
		patches:   make(map[string]PatchMetadata),
		validator: func(PatchMetadata) error { return nil },
		authority: append(ed25519.PublicKey(nil), key...),
	}
}

// SetValidator overrides the default validator. It can be used to plug in
// on-chain governance checks.
func (p *PatchManager) SetValidator(v PatchValidator) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if v == nil {
		v = func(PatchMetadata) error { return nil }
	}
	p.validator = v
}

// Apply records a patch ID using default metadata for backwards compatibility
// with existing tooling.
func (p *PatchManager) Apply(id string) {
	if id == "" {
		return
	}
	_ = p.ApplyMetadata(PatchMetadata{ID: id, Version: "legacy", SubmittedAt: time.Now().UTC()})
}

// ApplyMetadata records a patch after running all validators.
func (p *PatchManager) ApplyMetadata(meta PatchMetadata) error {
	if meta.ID == "" {
		return errors.New("security: patch id required")
	}
	if meta.SubmittedAt.IsZero() {
		meta.SubmittedAt = time.Now().UTC()
	}
	if meta.AppliedAt.IsZero() {
		meta.AppliedAt = time.Now().UTC()
	}
	if meta.Metadata == nil {
		meta.Metadata = make(map[string]string)
	}
	if err := p.verifySignature(meta); err != nil {
		return err
	}
	if err := p.validator(meta); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.patches[meta.ID]; exists {
		return fmt.Errorf("security: patch %s already recorded", meta.ID)
	}
	p.patches[meta.ID] = meta
	p.order = append(p.order, meta.ID)
	return nil
}

// Applied returns the identifiers of all applied patches ordered by application
// time. The method maintains backwards compatibility with legacy CLI commands.
func (p *PatchManager) Applied() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]string, len(p.order))
	copy(out, p.order)
	return out
}

// Metadata returns a snapshot of all known patches including detailed metadata.
func (p *PatchManager) Metadata() []PatchMetadata {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]PatchMetadata, 0, len(p.order))
	for _, id := range p.order {
		out = append(out, p.patches[id])
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].AppliedAt.Equal(out[j].AppliedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].AppliedAt.Before(out[j].AppliedAt)
	})
	return out
}

// verifySignature validates meta.Signature when an authority key is registered.
func (p *PatchManager) verifySignature(meta PatchMetadata) error {
	if len(p.authority) == 0 || len(meta.Signature) == 0 || len(meta.Digest) == 0 {
		return nil
	}
	payload := patchDigestPayload(meta)
	if !ed25519.Verify(p.authority, payload, meta.Signature) {
		return errors.New("security: patch signature invalid")
	}
	return nil
}

func patchDigestPayload(meta PatchMetadata) []byte {
	buf := []byte(meta.ID)
	buf = append(buf, meta.Version...)
	buf = append(buf, meta.Digest...)
	keys := make([]string, 0, len(meta.Metadata))
	for k := range meta.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := meta.Metadata[k]
		buf = append(buf, k...)
		buf = append(buf, v...)
	}
	return buf
}
