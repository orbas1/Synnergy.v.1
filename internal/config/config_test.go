package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeTempFile creates a temporary file with provided contents and returns its path.
func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	file := filepath.Join(dir, name)
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return file
}

func TestLoadYAML(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeTempFile(t, dir, "config.yaml", `
environment: production
log:
  level: debug
  outputs: ["stdout", "eventstream"]
  format: json
server:
  host: "127.0.0.1"
  port: 9000
  grpc_port: 9100
  grpc_host: "127.0.0.1"
  public_endpoint: "https://api.synnergy.example"
  timeout: 30s
  read_timeout: 15s
  write_timeout: 15s
  rate_limit:
    requests_per_second: 500
    burst: 800
database:
  url: "postgres://admin:pass@db:5432/synnergy?sslmode=disable"
  migrations_path: "schema"
consensus:
  engine: synnergy-ibft
  committee_size: 21
  block_time: 1s
  finality_threshold: 0.75
  gas_floor: 120000
  max_block_gas: 40000000
vm:
  gas_limit: 42000000
  schedule_path: "configs/enterprise_gas.yaml"
  execution_timeout: 7s
wallet:
  keystore_path: "/secure"
  default_account: "validator-1"
  signing_algorithm: ed25519
  approval_threshold: 2
cli:
  default_network: "synnergy-mainnet"
  gas_price: 2100000
  broadcast_endpoint: "https://cli.synnergy.example/api"
telemetry:
  tracing:
    enabled: true
    endpoint: "otel.collector:4317"
security:
  allowed_origins: ["https://console.synnergy.example"]
network:
  seed_nodes: ["seed-a.synnergy.example:7070", "seed-b.synnergy.example:7070"]
  authority_nodes: ["auth-a.synnergy.example:7070"]
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load YAML config: %v", err)
	}

	if cfg.Server.Port != 9000 {
		t.Fatalf("expected server port 9000 got %d", cfg.Server.Port)
	}
	if cfg.Consensus.CommitteeSize != 21 {
		t.Fatalf("expected committee size 21 got %d", cfg.Consensus.CommitteeSize)
	}
	if cfg.VM.SchedulePath != "configs/enterprise_gas.yaml" {
		t.Fatalf("expected custom gas schedule path, got %s", cfg.VM.SchedulePath)
	}
	if cfg.Security.AllowedOrigins[0] != "https://console.synnergy.example" {
		t.Fatalf("expected allowed origin propagated")
	}
}

func TestEnvOverrideAndJSON(t *testing.T) {
	dir := t.TempDir()
	raw := map[string]any{
		"environment": "staging",
		"log_level":   "info",
		"server": map[string]any{
			"host": "0.0.0.0",
			"port": 8080,
		},
	}
	content, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	path := writeTempFile(t, dir, "config.json", string(content))

	t.Setenv("SYN_CLI_GAS_PRICE", "9000000")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load JSON config: %v", err)
	}
	if cfg.CLI.GasPrice != 9_000_000 {
		t.Fatalf("expected env override to 9000000 got %d", cfg.CLI.GasPrice)
	}
}

func TestValidation(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.yaml", `server:
  tls:
    enabled: true
`)

	if _, err := Load(path); err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}

func TestWithDefaultPathsFallback(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	configPath := writeTempFile(t, dir, "custom.yaml", `environment: development`)

	cfg, err := Load("", WithDefaultPaths(configPath))
	if err != nil {
		t.Fatalf("load with default paths: %v", err)
	}
	if cfg.Environment != "development" {
		t.Fatalf("expected development environment fallback")
	}
}

func TestFingerprintDeterministic(t *testing.T) {
	t.Parallel()
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	fp1 := cfg.Fingerprint()
	fp2 := cfg.Clone().Fingerprint()
	if fp1 != fp2 {
		t.Fatalf("expected deterministic fingerprint, got %s and %s", fp1, fp2)
	}
	if len(fp1) == 0 {
		t.Fatalf("fingerprint should not be empty")
	}
}

func TestCloneIsolated(t *testing.T) {
	t.Parallel()
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	clone := cfg.Clone()
	clone.Log.Outputs[0] = "stderr"
	clone.Network.SeedNodes[0] = "overridden:7070"
	if cfg.Log.Outputs[0] == "stderr" {
		t.Fatalf("mutation leaked to original log outputs")
	}
	if cfg.Network.SeedNodes[0] == "overridden:7070" {
		t.Fatalf("mutation leaked to original seed nodes")
	}
}

func TestValidateCIDRFails(t *testing.T) {
	t.Parallel()
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	cfg.Security.PermitCIDRs = []string{"invalid"}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected cidr validation failure")
	}
}

func TestConsensusGasFloorValidation(t *testing.T) {
	t.Parallel()
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	cfg.Consensus.MaxBlockGas = cfg.Consensus.GasFloor - 1
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error when gas floor exceeds max block gas")
	}
}

func TestTimeoutsRemainPositive(t *testing.T) {
	t.Parallel()
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	cfg.Server.Timeout = -1
	if err := ensureValidator().Struct(cfg); err == nil {
		t.Fatalf("expected validator to reject negative timeout")
	}
}

func TestLoadHonoursBackwardsCompatibleKeys(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeTempFile(t, dir, "legacy.yaml", `
log_level: warn
environment: production
server:
  host: "0.0.0.0"
database:
  url: "postgres://legacy@localhost:5432/synnergy"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load legacy config: %v", err)
	}
	if cfg.Log.Level != "warn" {
		t.Fatalf("expected legacy log level to map to structured log config, got %s", cfg.Log.Level)
	}
}

func TestWithEnvPrefix(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "config.yaml", `environment: development`)
	t.Setenv("CUSTOM_ENVIRONMENT", "production")
	cfg, err := Load(path, WithEnvPrefix("CUSTOM"))
	if err != nil {
		t.Fatalf("load with custom prefix: %v", err)
	}
	if cfg.Environment != "production" {
		t.Fatalf("expected environment override from custom prefix")
	}
}

func TestLoadRejectsMissingTLSMaterial(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeTempFile(t, dir, "config.yaml", `
server:
  tls:
    enabled: true
`)
	if _, err := Load(path); err == nil {
		t.Fatalf("expected tls validation error")
	}
}

func TestApplyBackwardsCompatibilityDefaultsOutputs(t *testing.T) {
	t.Parallel()
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	if len(cfg.Log.Outputs) == 0 {
		t.Fatalf("log outputs should default to stdout")
	}
}

func TestLoadFromDefaultPathsIgnoresMissing(t *testing.T) {
	t.Parallel()
	cfg, err := Load("", WithDefaultPaths(filepath.Join(t.TempDir(), "missing.yaml")))
	if err != nil {
		t.Fatalf("load should succeed even if default file missing: %v", err)
	}
	if cfg.Server.Timeout <= 0 {
		t.Fatalf("expected defaults to apply")
	}
}

func TestMaxCommandDurationDecodes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeTempFile(t, dir, "config.yaml", `cli:
  max_command_duration: 45s
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load with custom duration: %v", err)
	}
	if cfg.CLI.MaxCommandDuration != 45*time.Second {
		t.Fatalf("expected CLI duration override, got %s", cfg.CLI.MaxCommandDuration)
	}
}
