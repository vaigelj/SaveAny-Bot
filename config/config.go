// Package config handles loading and validation of bot configuration
// from environment variables and config files.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the SaveAny-Bot application.
type Config struct {
	// Telegram bot token from @BotFather
	BotToken string `mapstructure:"bot_token"`

	// Telegram API ID from https://my.telegram.org
	APIID int `mapstructure:"api_id"`

	// Telegram API Hash from https://my.telegram.org
	APIHash string `mapstructure:"api_hash"`

	// Allowed user IDs (empty means all users allowed)
	AllowedUsers []int64 `mapstructure:"allowed_users"`

	// Admin user IDs with elevated permissions
	Admins []int64 `mapstructure:"admins"`

	// Storage backends configuration
	Storage StorageConfig `mapstructure:"storage"`

	// Log level: debug, info, warn, error
	LogLevel string `mapstructure:"log_level"`

	// Temp directory for intermediate file storage
	TempDir string `mapstructure:"temp_dir"`

	// Max concurrent download workers
	Workers int `mapstructure:"workers"`
}

// StorageConfig defines available storage backend configurations.
type StorageConfig struct {
	// Local filesystem storage path
	Local LocalStorage `mapstructure:"local"`

	// Alist storage configuration
	Alist AlistStorage `mapstructure:"alist"`

	// WebDAV storage configuration
	WebDAV WebDAVStorage `mapstructure:"webdav"`
}

// LocalStorage configures local filesystem storage.
type LocalStorage struct {
	Enabled bool   `mapstructure:"enabled"`
	BasePath string `mapstructure:"base_path"`
}

// AlistStorage configures Alist storage backend.
type AlistStorage struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	BasePath string `mapstructure:"base_path"`
}

// WebDAVStorage configures WebDAV storage backend.
type WebDAVStorage struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	BasePath string `mapstructure:"base_path"`
}

var Cfg *Config

// Load reads configuration from the specified file path and environment variables.
// Environment variables take precedence over file values.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("log_level", "info")
	v.SetDefault("temp_dir", os.TempDir())
	v.SetDefault("workers", 3)
	v.SetDefault("storage.local.enabled", true)
	v.SetDefault("storage.local.base_path", "./downloads")

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/saveany-bot/")
	}

	// Allow overriding any config key via environment variables
	v.SetEnvPrefix("SAVEANY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	Cfg = cfg
	return cfg, nil
}

// validate checks that required configuration fields are present and valid.
func (c *Config) validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("bot_token is required")
	}
	if c.APIID == 0 {
		return fmt.Errorf("api_id is required")
	}
	if c.APIHash == "" {
		return fmt.Errorf("api_hash is required")
	}
	if c.Workers < 1 {
		return fmt.Errorf("workers must be at least 1")
	}
	return nil
}
