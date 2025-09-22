//go:build prod

package config

// DefaultConfigPath points to the production configuration file.
// Authority nodes compile with `-tags prod` to enforce the hardened operating
// posture and telemetry requirements encoded in configs/prod.yaml.
const DefaultConfigPath = "configs/prod.yaml"
