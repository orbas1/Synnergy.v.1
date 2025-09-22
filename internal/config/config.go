// Package config provides centralized configuration management.
//
// The module exposes a strongly typed configuration tree that unifies runtime
// settings for the CLI, core consensus, networking stack, virtual machine and
// web control-plane services. Configuration values can be sourced from YAML or
// JSON files, overridden via environment variables and validated against
// enterprise grade constraints that ensure gas metering, governance policies and
// security boundaries remain aligned across subsystems.
package config

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// Config represents the full runtime configuration. It captures operational
// parameters for the CLI, consensus core, node services, gas metering rules and
// telemetry surfaces that the function web consumes for orchestration.
type Config struct {
	Version     string          `mapstructure:"version" validate:"omitempty"`
	Environment string          `mapstructure:"environment" validate:"required,oneof=development staging production"`
	Log         LogConfig       `mapstructure:"log" validate:"required"`
	Server      ServerConfig    `mapstructure:"server" validate:"required"`
	Database    DatabaseConfig  `mapstructure:"database" validate:"required"`
	Consensus   ConsensusConfig `mapstructure:"consensus" validate:"required"`
	VM          VMConfig        `mapstructure:"vm" validate:"required"`
	Wallet      WalletConfig    `mapstructure:"wallet" validate:"required"`
	CLI         CLIConfig       `mapstructure:"cli" validate:"required"`
	Telemetry   TelemetryConfig `mapstructure:"telemetry" validate:"required"`
	Security    SecurityConfig  `mapstructure:"security" validate:"required"`
	Network     NetworkConfig   `mapstructure:"network" validate:"required"`
}

// LogConfig describes structured logging behaviour used by the CLI, node
// services and the governance dashboards.
type LogConfig struct {
	Level         string         `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format        string         `mapstructure:"format" validate:"required,oneof=json text"`
	IncludeCaller bool           `mapstructure:"include_caller"`
	Outputs       []string       `mapstructure:"outputs" validate:"min=1,dive,oneof=stdout stderr file eventstream"`
	Sampling      LogSampling    `mapstructure:"sampling"`
	StaticFields  map[string]any `mapstructure:"static_fields" validate:"omitempty"`
}

// LogSampling controls burst protection for very verbose workloads.
type LogSampling struct {
	Enabled    bool `mapstructure:"enabled"`
	Initial    int  `mapstructure:"initial" validate:"gte=0"`
	Thereafter int  `mapstructure:"thereafter" validate:"gte=0"`
}

// ServerConfig contains HTTP and gRPC server settings exposed through the CLI
// and JavaScript control plane.
type ServerConfig struct {
	Host           string        `mapstructure:"host" validate:"required,hostname|ip"`
	Port           int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	GRPCHost       string        `mapstructure:"grpc_host" validate:"required,hostname|ip"`
	GRPCPort       int           `mapstructure:"grpc_port" validate:"required,min=1,max=65535"`
	PublicEndpoint string        `mapstructure:"public_endpoint" validate:"required,uri"`
	Timeout        time.Duration `mapstructure:"timeout" validate:"required,gt=0"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout" validate:"required,gt=0"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout" validate:"required,gt=0"`
	RateLimit      RateLimit     `mapstructure:"rate_limit" validate:"required"`
	TLS            TLSConfig     `mapstructure:"tls"`
	HealthCheck    bool          `mapstructure:"health_check"`
}

// RateLimit defines throttling behaviour for public RPC and CLI submissions.
type RateLimit struct {
	RequestsPerSecond float64 `mapstructure:"requests_per_second" validate:"gte=0"`
	Burst             int     `mapstructure:"burst" validate:"gte=0"`
}

// TLSConfig encapsulates TLS expectations for all public facing surfaces.
type TLSConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	CertFile        string `mapstructure:"cert_file" validate:"omitempty"`
	KeyFile         string `mapstructure:"key_file" validate:"omitempty"`
	ClientCAFile    string `mapstructure:"client_ca_file" validate:"omitempty"`
	RequireClientCA bool   `mapstructure:"require_client_ca"`
	MinVersion      string `mapstructure:"min_version" validate:"omitempty,oneof=1.2 1.3"`
}

// DatabaseConfig configures relational persistence.
type DatabaseConfig struct {
	URL             string        `mapstructure:"url" validate:"required,url"`
	ReadReplicaURLs []string      `mapstructure:"read_replicas" validate:"omitempty,dive,url"`
	MaxConnections  int           `mapstructure:"max_connections" validate:"gt=0"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" validate:"gte=0"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" validate:"gte=0"`
	MigrationsPath  string        `mapstructure:"migrations_path" validate:"required"`
}

// ConsensusConfig controls validator and committee behaviour to keep gas tables
// aligned with the VM and CLI validators.
type ConsensusConfig struct {
	Engine            string        `mapstructure:"engine" validate:"required,oneof=synnergy-pos synnergy-dpos synnergy-ibft"`
	CommitteeSize     int           `mapstructure:"committee_size" validate:"required,gt=0"`
	BlockTime         time.Duration `mapstructure:"block_time" validate:"required,gt=0"`
	FinalityThreshold float64       `mapstructure:"finality_threshold" validate:"gt=0,lte=1"`
	MessageTTL        time.Duration `mapstructure:"message_ttl" validate:"required,gt=0"`
	MaxDrift          time.Duration `mapstructure:"max_drift" validate:"required,gt=0"`
	SafetyBuffer      int           `mapstructure:"safety_buffer" validate:"gte=0"`
	GasFloor          uint64        `mapstructure:"gas_floor" validate:"gt=0"`
	MaxBlockGas       uint64        `mapstructure:"max_block_gas" validate:"gt=0"`
	AllowManualStop   bool          `mapstructure:"allow_manual_stop"`
}

// VMConfig captures execution policies enforced by the Synnergy VM.
type VMConfig struct {
	ExecutionTimeout  time.Duration `mapstructure:"execution_timeout" validate:"required,gt=0"`
	MaxStackDepth     int           `mapstructure:"max_stack_depth" validate:"gt=0"`
	GasLimit          uint64        `mapstructure:"gas_limit" validate:"gt=0"`
	SchedulePath      string        `mapstructure:"schedule_path" validate:"required"`
	Deterministic     bool          `mapstructure:"deterministic"`
	CacheSize         int           `mapstructure:"cache_size" validate:"gte=0"`
	MeteringPrecision uint64        `mapstructure:"metering_precision" validate:"gt=0"`
	OpcodeSet         string        `mapstructure:"opcode_set" validate:"required"`
	Precompiles       []string      `mapstructure:"precompiles_enabled" validate:"min=1,dive,required"`
}

// WalletConfig controls multisig wallets and CLI key material.
type WalletConfig struct {
	KeystorePath      string `mapstructure:"keystore_path" validate:"required"`
	DefaultAccount    string `mapstructure:"default_account" validate:"required"`
	HSMEnabled        bool   `mapstructure:"hsm_enabled"`
	HSMSlot           string `mapstructure:"hsm_slot" validate:"omitempty"`
	SigningAlgorithm  string `mapstructure:"signing_algorithm" validate:"required,oneof=ed25519 secp256k1 sr25519"`
	ApprovalThreshold int    `mapstructure:"approval_threshold" validate:"gte=1"`
	MultiSigEnabled   bool   `mapstructure:"multi_sig_enabled"`
}

// CLIConfig configures the command line workflow for validators and operators.
type CLIConfig struct {
	DefaultNetwork     string        `mapstructure:"default_network" validate:"required"`
	GasPrice           uint64        `mapstructure:"gas_price" validate:"gt=0"`
	OutputFormat       string        `mapstructure:"output_format" validate:"required,oneof=json text table"`
	AutoComplete       bool          `mapstructure:"auto_complete"`
	DataDir            string        `mapstructure:"data_dir" validate:"required"`
	MaxCommandDuration time.Duration `mapstructure:"max_command_duration" validate:"required,gt=0"`
	BroadcastEndpoint  string        `mapstructure:"broadcast_endpoint" validate:"required,url"`
}

// TelemetryConfig wires metrics, tracing and health probes for the web UX.
type TelemetryConfig struct {
	Metrics MetricsConfig `mapstructure:"metrics" validate:"required"`
	Tracing TracingConfig `mapstructure:"tracing" validate:"required"`
	Health  HealthConfig  `mapstructure:"health" validate:"required"`
}

// MetricsConfig exposes Prometheus friendly settings.
type MetricsConfig struct {
	Enabled            bool          `mapstructure:"enabled"`
	Endpoint           string        `mapstructure:"endpoint" validate:"required,hostname_port"`
	CollectionInterval time.Duration `mapstructure:"collection_interval" validate:"required,gt=0"`
	PushGatewayURL     string        `mapstructure:"push_gateway_url" validate:"omitempty,url"`
}

// TracingConfig exposes OTLP/Zipkin/Jaeger controls used by the function web.
type TracingConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	Provider    string  `mapstructure:"provider" validate:"required,oneof=otlp jaeger zipkin"`
	Endpoint    string  `mapstructure:"endpoint" validate:"required,hostname_port"`
	SampleRatio float64 `mapstructure:"sample_ratio" validate:"gte=0,lte=1"`
}

// HealthConfig powers liveness/readiness endpoints for automation.
type HealthConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port" validate:"min=0,max=65535"`
}

// SecurityConfig ensures policy enforcement across wallet, consensus and web.
type SecurityConfig struct {
	EnableAudit         bool          `mapstructure:"enable_audit"`
	AuditTrailRetention time.Duration `mapstructure:"audit_trail_retention" validate:"required,gt=0"`
	PermitCIDRs         []string      `mapstructure:"permit_cidrs" validate:"min=1,dive,cidr"`
	JWTSigningAlgorithm string        `mapstructure:"jwt_signing_algorithm" validate:"required,oneof=ed25519 es256 es256k rs256"`
	JWTIssuer           string        `mapstructure:"jwt_issuer" validate:"required"`
	KeyManagementURI    string        `mapstructure:"key_management_uri" validate:"required,uri"`
	AllowedOrigins      []string      `mapstructure:"allowed_origins" validate:"min=1,dive,uri"`
	MinPasswordLength   int           `mapstructure:"min_password_length" validate:"gte=12"`
	SignatureScheme     string        `mapstructure:"signature_scheme" validate:"required,oneof=ed25519 bls12-381 secp256k1"`
	ConfidentialTx      bool          `mapstructure:"confidential_transactions"`
}

// NetworkConfig tunes peer to peer and RPC transport parameters.
type NetworkConfig struct {
	P2PAddress        string        `mapstructure:"p2p_address" validate:"required,hostname_port"`
	RPCAddress        string        `mapstructure:"rpc_address" validate:"required,hostname_port"`
	MaxPeers          int           `mapstructure:"max_peers" validate:"gt=0"`
	SeedNodes         []string      `mapstructure:"seed_nodes" validate:"min=1,dive,hostname_port"`
	AuthorityNodes    []string      `mapstructure:"authority_nodes" validate:"min=1,dive,hostname_port"`
	AllowPrivatePeers bool          `mapstructure:"allow_private_peers"`
	SyncRetryInterval time.Duration `mapstructure:"sync_retry_interval" validate:"required,gt=0"`
}

var (
	validatorOnce sync.Once
	validate      *validator.Validate
)

// Option customises the behaviour of the Load helper.
type Option func(*loadOptions)

type loadOptions struct {
	envPrefix    string
	defaultPaths []string
}

// WithEnvPrefix customises the environment variable prefix used by Load.
func WithEnvPrefix(prefix string) Option {
	return func(lo *loadOptions) {
		lo.envPrefix = prefix
	}
}

// WithDefaultPaths overrides the search order for configuration files when a
// path is not explicitly provided.
func WithDefaultPaths(paths ...string) Option {
	return func(lo *loadOptions) {
		lo.defaultPaths = append([]string{}, paths...)
	}
}

// Load reads configuration from disk, environment variables or defaults. The
// function guarantees deterministic validation so that CLI and web orchestrators
// observe an identical runtime view.
func Load(path string, opts ...Option) (*Config, error) {
	options := loadOptions{
		envPrefix:    "SYN",
		defaultPaths: []string{DefaultConfigPath},
	}
	for _, opt := range opts {
		opt(&options)
	}

	v := viper.New()
	applyDefaults(v)
	registerAliases(v)

	if path != "" {
		v.SetConfigFile(path)
		if err := readConfigFile(v); err != nil {
			return nil, err
		}
	} else {
		if err := loadFromDefaultPaths(v, options.defaultPaths); err != nil {
			return nil, err
		}
	}

	v.SetEnvPrefix(strings.ToLower(options.envPrefix))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg, viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc())); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.applyBackwardsCompatibility(v); err != nil {
		return nil, err
	}

	if err := ensureValidator().Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate performs cross-field validations that cannot be expressed with
// struct tags alone.
func (c *Config) Validate() error {
	if c.Server.TLS.Enabled {
		if c.Server.TLS.CertFile == "" || c.Server.TLS.KeyFile == "" {
			return errors.New("server.tls requires cert_file and key_file when enabled")
		}
	}

	if c.Consensus.MaxBlockGas < c.Consensus.GasFloor {
		return fmt.Errorf("consensus.max_block_gas must be >= consensus.gas_floor")
	}

	if c.VM.GasLimit < c.Consensus.GasFloor {
		return fmt.Errorf("vm.gas_limit must be >= consensus.gas_floor")
	}

	if err := validateCIDRs(c.Security.PermitCIDRs); err != nil {
		return fmt.Errorf("security.permit_cidrs: %w", err)
	}

	return nil
}

// Clone returns a deep copy of the configuration which can safely be mutated by
// callers without affecting the shared runtime configuration.
func (c *Config) Clone() Config {
	clone := *c
	clone.Log.Outputs = append([]string(nil), c.Log.Outputs...)
	clone.Log.StaticFields = copyMap(c.Log.StaticFields)
	clone.Database.ReadReplicaURLs = append([]string(nil), c.Database.ReadReplicaURLs...)
	clone.VM.Precompiles = append([]string(nil), c.VM.Precompiles...)
	clone.Security.PermitCIDRs = append([]string(nil), c.Security.PermitCIDRs...)
	clone.Security.AllowedOrigins = append([]string(nil), c.Security.AllowedOrigins...)
	clone.Network.SeedNodes = append([]string(nil), c.Network.SeedNodes...)
	clone.Network.AuthorityNodes = append([]string(nil), c.Network.AuthorityNodes...)
	return clone
}

// EffectiveHTTPAddress returns host:port string for HTTP server bindings.
func (c *Config) EffectiveHTTPAddress() string {
	return net.JoinHostPort(c.Server.Host, strconv.Itoa(c.Server.Port))
}

// EffectiveGRPCAddress returns host:port for gRPC services.
func (c *Config) EffectiveGRPCAddress() string {
	return net.JoinHostPort(c.Server.GRPCHost, strconv.Itoa(c.Server.GRPCPort))
}

func (c *Config) applyBackwardsCompatibility(v *viper.Viper) error {
	if c.Log.Level == "" {
		if lvl := v.GetString("log_level"); lvl != "" {
			c.Log.Level = lvl
		} else {
			c.Log.Level = "info"
		}
	}
	if c.Log.Format == "" {
		if format := v.GetString("log_format"); format != "" {
			c.Log.Format = format
		} else {
			c.Log.Format = "json"
		}
	}
	if len(c.Log.Outputs) == 0 {
		c.Log.Outputs = []string{"stdout"}
	}
	return nil
}

func copyMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func ensureValidator() *validator.Validate {
	validatorOnce.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
		_ = validate.RegisterValidation("cidr", func(fl validator.FieldLevel) bool {
			_, _, err := net.ParseCIDR(fl.Field().String())
			return err == nil
		})
	})
	return validate
}

func readConfigFile(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if errors.As(err, &nf) {
			return fmt.Errorf("config file %s not found", v.ConfigFileUsed())
		}
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}

func loadFromDefaultPaths(v *viper.Viper, paths []string) error {
	var lastErr error
	for _, candidate := range paths {
		if candidate == "" {
			continue
		}
		abs := candidate
		if !filepath.IsAbs(candidate) {
			abs = filepath.Clean(candidate)
		}
		if _, err := os.Stat(abs); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				lastErr = err
				continue
			}
			return fmt.Errorf("stat config %s: %w", abs, err)
		}
		v.SetConfigFile(abs)
		if err := v.ReadInConfig(); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	if lastErr != nil {
		// No config file was found, but defaults still apply.
		return nil
	}
	return nil
}

func registerAliases(v *viper.Viper) {
	v.RegisterAlias("log.level", "log_level")
	v.RegisterAlias("log.format", "log_format")
}

func applyDefaults(v *viper.Viper) {
	v.SetDefault("version", "latest")
	v.SetDefault("environment", "development")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.include_caller", true)
	v.SetDefault("log.outputs", []string{"stdout"})
	v.SetDefault("log.sampling.enabled", true)
	v.SetDefault("log.sampling.initial", 100)
	v.SetDefault("log.sampling.thereafter", 100)

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.grpc_host", "0.0.0.0")
	v.SetDefault("server.grpc_port", 9090)
	v.SetDefault("server.public_endpoint", "http://localhost:8080")
	v.SetDefault("server.timeout", 15*time.Second)
	v.SetDefault("server.read_timeout", 10*time.Second)
	v.SetDefault("server.write_timeout", 10*time.Second)
	v.SetDefault("server.rate_limit.requests_per_second", 250.0)
	v.SetDefault("server.rate_limit.burst", 500)
	v.SetDefault("server.health_check", true)

	v.SetDefault("database.url", "postgres://synnergy:synnergy@localhost:5432/synnergy?sslmode=disable")
	v.SetDefault("database.max_connections", 32)
	v.SetDefault("database.conn_max_idle_time", 5*time.Minute)
	v.SetDefault("database.conn_max_lifetime", time.Hour)
	v.SetDefault("database.migrations_path", "migrations")

	v.SetDefault("consensus.engine", "synnergy-pos")
	v.SetDefault("consensus.committee_size", 11)
	v.SetDefault("consensus.block_time", 2*time.Second)
	v.SetDefault("consensus.finality_threshold", 0.67)
	v.SetDefault("consensus.message_ttl", 30*time.Second)
	v.SetDefault("consensus.max_drift", 3*time.Second)
	v.SetDefault("consensus.safety_buffer", 2)
	v.SetDefault("consensus.gas_floor", uint64(50_000))
	v.SetDefault("consensus.max_block_gas", uint64(25_000_000))

	v.SetDefault("vm.execution_timeout", 5*time.Second)
	v.SetDefault("vm.max_stack_depth", 2048)
	v.SetDefault("vm.gas_limit", uint64(100_000_000))
	v.SetDefault("vm.schedule_path", "configs/gas_table.yaml")
	v.SetDefault("vm.deterministic", true)
	v.SetDefault("vm.cache_size", 1024)
	v.SetDefault("vm.metering_precision", uint64(4))
	v.SetDefault("vm.opcode_set", "synnergy/v1")
	v.SetDefault("vm.precompiles_enabled", []string{"bls12-381", "poseidon"})

	v.SetDefault("wallet.keystore_path", "data/wallets")
	v.SetDefault("wallet.default_account", "authority-0")
	v.SetDefault("wallet.hsm_enabled", false)
	v.SetDefault("wallet.signing_algorithm", "ed25519")
	v.SetDefault("wallet.approval_threshold", 2)
	v.SetDefault("wallet.multi_sig_enabled", true)

	v.SetDefault("cli.default_network", "synnergy-devnet")
	v.SetDefault("cli.gas_price", uint64(1_000_000))
	v.SetDefault("cli.output_format", "table")
	v.SetDefault("cli.auto_complete", true)
	v.SetDefault("cli.data_dir", "~/.synnergy")
	v.SetDefault("cli.max_command_duration", 2*time.Minute)
	v.SetDefault("cli.broadcast_endpoint", "http://localhost:8080")

	v.SetDefault("telemetry.metrics.enabled", true)
	v.SetDefault("telemetry.metrics.endpoint", "0.0.0.0:9100")
	v.SetDefault("telemetry.metrics.collection_interval", 15*time.Second)
	v.SetDefault("telemetry.tracing.enabled", false)
	v.SetDefault("telemetry.tracing.provider", "otlp")
	v.SetDefault("telemetry.tracing.endpoint", "localhost:4317")
	v.SetDefault("telemetry.tracing.sample_ratio", 0.1)
	v.SetDefault("telemetry.health.enabled", true)
	v.SetDefault("telemetry.health.port", 9091)

	v.SetDefault("security.enable_audit", true)
	v.SetDefault("security.audit_trail_retention", 30*24*time.Hour)
	v.SetDefault("security.permit_cidrs", []string{"10.0.0.0/8", "192.168.0.0/16"})
	v.SetDefault("security.jwt_signing_algorithm", "ed25519")
	v.SetDefault("security.jwt_issuer", "synnergy-network")
	v.SetDefault("security.key_management_uri", "kms://synnergy/local")
	v.SetDefault("security.allowed_origins", []string{"https://console.synnergy.io"})
	v.SetDefault("security.min_password_length", 16)
	v.SetDefault("security.signature_scheme", "ed25519")
	v.SetDefault("security.confidential_transactions", true)

	v.SetDefault("network.p2p_address", "0.0.0.0:7070")
	v.SetDefault("network.rpc_address", "0.0.0.0:8545")
	v.SetDefault("network.max_peers", 2000)
	v.SetDefault("network.seed_nodes", []string{"seed-1.synnergy.io:7070", "seed-2.synnergy.io:7070"})
	v.SetDefault("network.authority_nodes", []string{"authority-1.synnergy.io:7070"})
	v.SetDefault("network.allow_private_peers", false)
	v.SetDefault("network.sync_retry_interval", 5*time.Second)
}

func validateCIDRs(cidrs []string) error {
	sort.Strings(cidrs)
	for _, cidr := range cidrs {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return fmt.Errorf("invalid cidr %s: %w", cidr, err)
		}
	}
	return nil
}

// Fingerprint returns a deterministic hash of the configuration that can be
// shared between CLI and authority nodes to ensure consensus over operational
// parameters.
func (c Config) Fingerprint() string {
	clone := c.Clone()
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%d|%d|%d|%t|%t", clone.Version, clone.Environment, clone.Log.Level, clone.Consensus.Engine, clone.VM.OpcodeSet, clone.Consensus.GasFloor, clone.Consensus.MaxBlockGas, clone.Server.Port, clone.Security.ConfidentialTx, clone.Wallet.HSMEnabled)
	for _, output := range clone.Log.Outputs {
		data += "|log:" + output
	}
	for _, seed := range clone.Network.SeedNodes {
		data += "|seed:" + seed
	}
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}
