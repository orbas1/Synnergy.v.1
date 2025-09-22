//go:build dev

package config

// DefaultConfigPath points to the development configuration file.
// Developers and integration tests build with `-tags dev` to ensure feature
// flags align with the richer defaults captured in configs/dev.yaml.
const DefaultConfigPath = "configs/dev.yaml"
