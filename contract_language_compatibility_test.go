package synnergy

import (
	"sort"
	"testing"
)

// resetLanguages restores the default set for test isolation.
func resetLanguages() {
	langMu.Lock()
	supportedLangs = map[string]struct{}{
		"wasm":       {},
		"golang":     {},
		"javascript": {},
		"solidity":   {},
		"rust":       {},
		"python":     {},
		"yul":        {},
	}
	langMu.Unlock()
}

func TestLanguageRegistry(t *testing.T) {
	resetLanguages()

	// Default languages should be supported
	defaults := []string{"wasm", "golang", "javascript", "solidity", "rust", "python", "yul"}
	for _, l := range defaults {
		if !IsLanguageSupported(l) {
			t.Fatalf("%s should be supported", l)
		}
	}

	// Adding and removing languages is case insensitive
	if err := AddSupportedLanguage("Move"); err != nil {
		t.Fatalf("add: %v", err)
	}
	if !IsLanguageSupported("MOVE") {
		t.Fatalf("move should be supported after addition")
	}
	if !RemoveSupportedLanguage("move") {
		t.Fatalf("expected removal to succeed")
	}
	if IsLanguageSupported("move") {
		t.Fatalf("move should not be supported after removal")
	}

	// SupportedContractLanguages should return a sorted list
	langs := SupportedContractLanguages()
	if !sort.StringsAreSorted(langs) {
		t.Fatalf("languages not sorted: %v", langs)
	}
}
