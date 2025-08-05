//go:build !dev && !test && !prod

package config

// DefaultConfigPath specifies the configuration file used when no build tags are provided.
// It points to the development configuration as a sensible fallback.
const DefaultConfigPath = "configs/dev.yaml"
