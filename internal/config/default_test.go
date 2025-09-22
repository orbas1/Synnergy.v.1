//go:build test

package config

// DefaultConfigPath points to the staging/test configuration file.
// Integration test binaries rely on this tag to exercise quorum safe defaults
// without mutating production grade configuration.
const DefaultConfigPath = "configs/test.yaml"
