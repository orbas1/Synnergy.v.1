package synnergy

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

// language registry ensures thread-safe tracking of supported contract languages.
var (
	langMu         sync.RWMutex
	supportedLangs = map[string]struct{}{
		"wasm":       {},
		"golang":     {},
		"javascript": {},
		"solidity":   {},
		"rust":       {},
		"python":     {},
		"yul":        {},
	}
)

// SupportedContractLanguages returns a sorted slice of registered languages.
func SupportedContractLanguages() []string {
	langMu.RLock()
	defer langMu.RUnlock()
	out := make([]string, 0, len(supportedLangs))
	for l := range supportedLangs {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}

// AddSupportedLanguage registers a new language; it's no-op if already present.
func AddSupportedLanguage(lang string) error {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return errors.New("language cannot be empty")
	}
	langMu.Lock()
	supportedLangs[lang] = struct{}{}
	langMu.Unlock()
	return nil
}

// RemoveSupportedLanguage deletes a language from the registry.
func RemoveSupportedLanguage(lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))
	langMu.Lock()
	_, ok := supportedLangs[lang]
	if ok {
		delete(supportedLangs, lang)
	}
	langMu.Unlock()
	return ok
}

// IsLanguageSupported reports whether the supplied language is recognised as
// compatible with the VM. The check is case-insensitive.
func IsLanguageSupported(lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))
	langMu.RLock()
	_, ok := supportedLangs[lang]
	langMu.RUnlock()
	return ok
}
