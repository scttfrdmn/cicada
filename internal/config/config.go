// Copyright 2025 Scott Friedman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the Cicada configuration.
type Config struct {
	// Version of config format
	Version string `mapstructure:"version" yaml:"version"`

	// AWS configuration
	AWS AWSConfig `mapstructure:"aws" yaml:"aws"`

	// Default sync options
	Sync SyncConfig `mapstructure:"sync" yaml:"sync"`

	// Watch configurations
	Watches []WatchConfig `mapstructure:"watches" yaml:"watches"`

	// Global settings
	Settings SettingsConfig `mapstructure:"settings" yaml:"settings"`
}

// AWSConfig holds AWS-specific configuration.
type AWSConfig struct {
	// Profile to use from ~/.aws/credentials
	Profile string `mapstructure:"profile" yaml:"profile"`

	// Region override (optional)
	Region string `mapstructure:"region" yaml:"region"`

	// Endpoint override for testing (optional)
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint"`
}

// SyncConfig holds default sync options.
type SyncConfig struct {
	// Default concurrency level
	Concurrency int `mapstructure:"concurrency" yaml:"concurrency"`

	// Default delete behavior
	Delete bool `mapstructure:"delete" yaml:"delete"`

	// Exclude patterns
	Exclude []string `mapstructure:"exclude" yaml:"exclude"`
}

// WatchConfig holds a watch configuration.
type WatchConfig struct {
	// Unique ID for this watch
	ID string `mapstructure:"id" yaml:"id"`

	// Source path
	Source string `mapstructure:"source" yaml:"source"`

	// Destination path
	Destination string `mapstructure:"destination" yaml:"destination"`

	// Debounce delay in seconds
	DebounceSeconds int `mapstructure:"debounce_seconds" yaml:"debounce_seconds"`

	// Min age in seconds
	MinAgeSeconds int `mapstructure:"min_age_seconds" yaml:"min_age_seconds"`

	// Delete source after sync
	DeleteSource bool `mapstructure:"delete_source" yaml:"delete_source"`

	// Sync on start
	SyncOnStart bool `mapstructure:"sync_on_start" yaml:"sync_on_start"`

	// Exclude patterns
	Exclude []string `mapstructure:"exclude" yaml:"exclude"`

	// Enabled flag
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
}

// SettingsConfig holds global settings.
type SettingsConfig struct {
	// Verbose logging
	Verbose bool `mapstructure:"verbose" yaml:"verbose"`

	// Log file path
	LogFile string `mapstructure:"log_file" yaml:"log_file"`

	// Check for updates on startup
	CheckUpdates bool `mapstructure:"check_updates" yaml:"check_updates"`
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Version: "1",
		AWS: AWSConfig{
			Profile: "default",
		},
		Sync: SyncConfig{
			Concurrency: 4,
			Delete:      false,
			Exclude:     []string{".git/**", ".DS_Store", "*.tmp", "*.swp"},
		},
		Watches: []WatchConfig{},
		Settings: SettingsConfig{
			Verbose:      false,
			CheckUpdates: true,
		},
	}
}

// ConfigPath returns the default config file path.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".cicada")
	return filepath.Join(configDir, "config.yaml"), nil
}

// ConfigDir returns the config directory path.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	return filepath.Join(home, ".cicada"), nil
}

// Load reads configuration from file.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	config := DefaultConfig()
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return config, nil
}

// LoadOrDefault loads config from default path or returns default config.
func LoadOrDefault() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	return Load(path)
}

// Save writes configuration to file.
func Save(config *Config, path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Create config directory if needed
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	// Marshal config to viper
	if err := v.MergeConfigMap(configToMap(config)); err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// Write to file
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// configToMap converts Config to map for viper.
func configToMap(c *Config) map[string]interface{} {
	return map[string]interface{}{
		"version": c.Version,
		"aws": map[string]interface{}{
			"profile":  c.AWS.Profile,
			"region":   c.AWS.Region,
			"endpoint": c.AWS.Endpoint,
		},
		"sync": map[string]interface{}{
			"concurrency": c.Sync.Concurrency,
			"delete":      c.Sync.Delete,
			"exclude":     c.Sync.Exclude,
		},
		"watches": watchesToMaps(c.Watches),
		"settings": map[string]interface{}{
			"verbose":       c.Settings.Verbose,
			"log_file":      c.Settings.LogFile,
			"check_updates": c.Settings.CheckUpdates,
		},
	}
}

// watchesToMaps converts watch configs to maps.
func watchesToMaps(watches []WatchConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, len(watches))
	for i, w := range watches {
		result[i] = map[string]interface{}{
			"id":               w.ID,
			"source":           w.Source,
			"destination":      w.Destination,
			"debounce_seconds": w.DebounceSeconds,
			"min_age_seconds":  w.MinAgeSeconds,
			"delete_source":    w.DeleteSource,
			"sync_on_start":    w.SyncOnStart,
			"exclude":          w.Exclude,
			"enabled":          w.Enabled,
		}
	}
	return result
}
