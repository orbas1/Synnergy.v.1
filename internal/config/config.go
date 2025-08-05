// Package config provides centralized configuration management.
package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config represents the application configuration. Fields can be populated from
// configuration files (YAML or JSON) and environment variables. Validation is
// performed using struct tags.
type Config struct {
	Environment string         `mapstructure:"environment" validate:"required,oneof=development staging production"`
	LogLevel    string         `mapstructure:"log_level" validate:"required,oneof=debug info warn error"`
	Server      ServerConfig   `mapstructure:"server" validate:"required"`
	Database    DatabaseConfig `mapstructure:"database" validate:"required"`
}

// ServerConfig holds configuration for the API server.
type ServerConfig struct {
	Host string `mapstructure:"host" validate:"required"`
	Port int    `mapstructure:"port" validate:"required,min=1,max=65535"`
}

// DatabaseConfig holds configuration for database connectivity.
type DatabaseConfig struct {
	URL string `mapstructure:"url" validate:"required,url"`
}

var validate = validator.New()

// Load reads the configuration from the provided file path. The file may be
// YAML or JSON. Environment variables with the prefix "SYN" override file
// values. Nested keys are represented with underscores. For example,
// "SYN_SERVER_PORT" overrides server.port.
func Load(path string) (*Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	v.SetEnvPrefix("syn")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set some sensible defaults.
	v.SetDefault("environment", "development")
	v.SetDefault("log_level", "info")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}
