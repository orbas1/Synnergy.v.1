package synnergy

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

// language registry ensures thread-safe tracking of supported contract languages.
type languageEntry struct {
	metadata LanguageMetadata
}

// LanguageMetadata captures metadata about a supported language.
type LanguageMetadata struct {
	Name     string
	Version  string
	Features []string
	AddedAt  time.Time
	Core     bool
}

// LanguageMetadataSink persists language metadata changes.
type LanguageMetadataSink interface {
	SaveLanguage(meta LanguageMetadata) error
	RemoveLanguage(name string) error
}

var (
	langMu                 sync.RWMutex
	supportedLangs         = map[string]languageEntry{}
	langSink               LanguageMetadataSink
	errEmptyLanguage       = errors.New("language cannot be empty")
	ErrCoreLanguageRemoval = errors.New("cannot remove core language")
)

func init() {
	initialiseDefaultLanguages()
}

func initialiseDefaultLanguages() {
	defaults := []LanguageMetadata{
		{Name: "wasm", Core: true},
		{Name: "golang", Core: true},
		{Name: "javascript", Core: true},
		{Name: "solidity", Core: true},
		{Name: "rust", Core: true},
		{Name: "python", Core: true},
		{Name: "yul", Core: true},
	}
	langMu.Lock()
	defer langMu.Unlock()
	supportedLangs = make(map[string]languageEntry, len(defaults))
	for _, meta := range defaults {
		if meta.AddedAt.IsZero() {
			meta.AddedAt = time.Now().UTC()
		}
		key := strings.ToLower(meta.Name)
		supportedLangs[key] = languageEntry{metadata: normaliseMetadata(meta)}
	}
}

// SetLanguageMetadataSink configures an optional persistence sink.
func SetLanguageMetadataSink(sink LanguageMetadataSink) {
	langMu.Lock()
	langSink = sink
	langMu.Unlock()
}

// SupportedContractLanguages returns a sorted slice of registered languages.
func SupportedContractLanguages() []string {
	langMu.RLock()
	defer langMu.RUnlock()
	out := make([]string, 0, len(supportedLangs))
	for _, entry := range supportedLangs {
		out = append(out, entry.metadata.Name)
	}
	sort.Strings(out)
	return out
}

// AddSupportedLanguage registers a new language; it's no-op if already present.
func AddSupportedLanguage(lang string) error {
	return AddLanguageMetadata(LanguageMetadata{Name: lang})
}

// AddLanguageMetadata registers metadata for a language. Existing entries are
// replaced with the supplied metadata.
func AddLanguageMetadata(meta LanguageMetadata) error {
	name := strings.ToLower(strings.TrimSpace(meta.Name))
	if name == "" {
		return errEmptyLanguage
	}
	normalised := normaliseMetadata(LanguageMetadata{
		Name:     name,
		Version:  meta.Version,
		Features: append([]string(nil), meta.Features...),
		AddedAt:  meta.AddedAt,
		Core:     meta.Core,
	})

	langMu.Lock()
	prev, existed := supportedLangs[name]
	supportedLangs[name] = languageEntry{metadata: normalised}
	sink := langSink
	langMu.Unlock()

	if sink != nil {
		if err := sink.SaveLanguage(copyMetadata(normalised)); err != nil {
			langMu.Lock()
			if existed {
				supportedLangs[name] = prev
			} else {
				delete(supportedLangs, name)
			}
			langMu.Unlock()
			return err
		}
	}
	return nil
}

// RemoveSupportedLanguage deletes a language from the registry.
func RemoveSupportedLanguage(lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return false
	}
	langMu.Lock()
	entry, ok := supportedLangs[lang]
	sink := langSink
	if !ok || entry.metadata.Core {
		langMu.Unlock()
		return false
	}
	delete(supportedLangs, lang)
	langMu.Unlock()
	if sink != nil {
		if err := sink.RemoveLanguage(lang); err != nil {
			langMu.Lock()
			supportedLangs[lang] = entry
			langMu.Unlock()
			return false
		}
	}
	return true
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

// SupportedLanguageMetadata returns a sorted copy of all language metadata.
func SupportedLanguageMetadata() []LanguageMetadata {
	langMu.RLock()
	out := make([]LanguageMetadata, 0, len(supportedLangs))
	for _, entry := range supportedLangs {
		out = append(out, copyMetadata(entry.metadata))
	}
	langMu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// GetLanguageMetadata returns metadata for the requested language.
func GetLanguageMetadata(name string) (LanguageMetadata, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	langMu.RLock()
	entry, ok := supportedLangs[key]
	langMu.RUnlock()
	if !ok {
		return LanguageMetadata{}, false
	}
	return copyMetadata(entry.metadata), true
}

func normaliseMetadata(meta LanguageMetadata) LanguageMetadata {
	meta.Name = strings.ToLower(strings.TrimSpace(meta.Name))
	if meta.AddedAt.IsZero() {
		meta.AddedAt = time.Now().UTC()
	}
	if meta.Features == nil {
		meta.Features = []string{}
	} else {
		features := make([]string, len(meta.Features))
		copy(features, meta.Features)
		sort.Strings(features)
		meta.Features = features
	}
	return meta
}

func copyMetadata(meta LanguageMetadata) LanguageMetadata {
	features := make([]string, len(meta.Features))
	copy(features, meta.Features)
	return LanguageMetadata{
		Name:     meta.Name,
		Version:  meta.Version,
		Features: features,
		AddedAt:  meta.AddedAt,
		Core:     meta.Core,
	}
}
