//go:build !dev && !test && !prod

package config

// DefaultConfigPath specifies the configuration file used when no build tags are provided.
// The development profile is selected so local CLI, VM and web workflows pick up
// sane defaults without additional flags.
const DefaultConfigPath = "configs/dev.yaml"
