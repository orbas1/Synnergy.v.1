package cli

import (
	"sort"
	"strings"
)

// parseMetadata converts CLI key=value inputs into a deterministic map.
// Empty keys are ignored and values default to an empty string when omitted.
func parseMetadata(values []string) map[string]string {
	if len(values) == 0 {
		return nil
	}
	meta := make(map[string]string, len(values))
	for _, kv := range values {
		parts := strings.SplitN(kv, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		val := ""
		if len(parts) == 2 {
			val = strings.TrimSpace(parts[1])
		}
		meta[key] = val
	}
	if len(meta) == 0 {
		return nil
	}
	// normalise key ordering for reproducible serialisation in tests
	ordered := make(map[string]string, len(meta))
	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		ordered[k] = meta[k]
	}
	return ordered
}
