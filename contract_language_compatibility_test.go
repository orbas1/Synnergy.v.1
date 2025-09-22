package synnergy

import (
	"sort"
	"testing"
	"time"
)

type recordingSink struct {
	saved   []LanguageMetadata
	removed []string
	saveErr error
}

func (r *recordingSink) SaveLanguage(meta LanguageMetadata) error {
	r.saved = append(r.saved, meta)
	return r.saveErr
}

func (r *recordingSink) RemoveLanguage(name string) error {
	r.removed = append(r.removed, name)
	return nil
}

// resetLanguages restores the default set for test isolation.
func resetLanguages() {
	langMu.Lock()
	langSink = nil
	langMu.Unlock()
	initialiseDefaultLanguages()
}

func TestLanguageRegistry(t *testing.T) {
	resetLanguages()

	defaults := []string{"golang", "javascript", "python", "rust", "solidity", "wasm", "yul"}
	for _, l := range defaults {
		if !IsLanguageSupported(l) {
			t.Fatalf("%s should be supported", l)
		}
		meta, ok := GetLanguageMetadata(l)
		if !ok {
			t.Fatalf("expected metadata for %s", l)
		}
		if !meta.Core {
			t.Fatalf("expected %s to be core", l)
		}
		if meta.AddedAt.IsZero() {
			t.Fatalf("expected AddedAt to be set for %s", l)
		}
	}

	sink := &recordingSink{}
	SetLanguageMetadataSink(sink)
	defer SetLanguageMetadataSink(nil)

	moveMeta := LanguageMetadata{
		Name:     "Move",
		Version:  "1.0",
		Features: []string{"resource", "modules"},
	}
	if err := AddLanguageMetadata(moveMeta); err != nil {
		t.Fatalf("add: %v", err)
	}
	if len(sink.saved) != 1 || sink.saved[0].Name != "move" {
		t.Fatalf("expected sink to record save, got %#v", sink.saved)
	}
	if !IsLanguageSupported("MOVE") {
		t.Fatalf("move should be supported after addition")
	}
	meta, ok := GetLanguageMetadata("move")
	if !ok {
		t.Fatalf("expected metadata for move")
	}
	if meta.Version != "1.0" {
		t.Fatalf("unexpected version: %s", meta.Version)
	}
	if meta.AddedAt.After(time.Now().UTC()) {
		t.Fatalf("expected AddedAt to be in the past")
	}
	if !sort.StringsAreSorted(meta.Features) {
		t.Fatalf("expected features to be sorted, got %v", meta.Features)
	}

	langs := SupportedContractLanguages()
	if !sort.StringsAreSorted(langs) {
		t.Fatalf("languages not sorted: %v", langs)
	}
	if RemoveSupportedLanguage("wasm") {
		t.Fatalf("core language should not be removable")
	}
	if !RemoveSupportedLanguage("move") {
		t.Fatalf("expected removal to succeed")
	}
	if len(sink.removed) != 1 || sink.removed[0] != "move" {
		t.Fatalf("expected sink removal for move, got %#v", sink.removed)
	}
	if IsLanguageSupported("move") {
		t.Fatalf("move should not be supported after removal")
	}
}
