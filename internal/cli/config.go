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

package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scttfrdmn/cicada/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewConfigCmd creates the config command.
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Cicada configuration",
		Long:  "Manage Cicada configuration file and settings",
	}

	cmd.AddCommand(
		NewConfigInitCmd(),
		NewConfigSetCmd(),
		NewConfigGetCmd(),
		NewConfigListCmd(),
	)

	return cmd
}

// NewConfigInitCmd creates the config init command.
func NewConfigInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		Long:  "Create a new configuration file with default values at ~/.cicada/config.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.ConfigPath()
			if err != nil {
				return fmt.Errorf("get config path: %w", err)
			}

			// Check if config already exists
			if _, err := os.Stat(path); err == nil && !force {
				return fmt.Errorf("config file already exists at %s (use --force to overwrite)", path)
			}

			// Create default config
			cfg := config.DefaultConfig()

			// Save config
			if err := config.Save(cfg, path); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			fmt.Printf("✓ Configuration initialized at %s\n", path)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing config file")

	return cmd
}

// NewConfigSetCmd creates the config set command.
func NewConfigSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value using dot notation.

Examples:
  cicada config set aws.profile myprofile
  cicada config set aws.region us-west-2
  cicada config set sync.concurrency 8
  cicada config set sync.delete true
  cicada config set settings.verbose true`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			// Load or create config
			cfg, err := config.LoadOrDefault()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			// Set the value
			if err := setConfigValue(cfg, key, value); err != nil {
				return err
			}

			// Save config
			path, err := config.ConfigPath()
			if err != nil {
				return fmt.Errorf("get config path: %w", err)
			}

			if err := config.Save(cfg, path); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			fmt.Printf("✓ Set %s = %s\n", key, value)
			return nil
		},
	}

	return cmd
}

// NewConfigGetCmd creates the config get command.
func NewConfigGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Long: `Get a configuration value using dot notation.

Examples:
  cicada config get aws.profile
  cicada config get aws.region
  cicada config get sync.concurrency
  cicada config get settings.verbose`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			// Load config
			cfg, err := config.LoadOrDefault()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			// Get the value
			value, err := getConfigValue(cfg, key)
			if err != nil {
				return err
			}

			fmt.Println(value)
			return nil
		},
	}

	return cmd
}

// NewConfigListCmd creates the config list command.
func NewConfigListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration",
		Long:  "Display the current configuration in YAML format",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			cfg, err := config.LoadOrDefault()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			// Marshal to YAML
			data, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("marshal config: %w", err)
			}

			fmt.Print(string(data))
			return nil
		},
	}

	return cmd
}

// setConfigValue sets a configuration value using dot notation.
func setConfigValue(cfg *config.Config, key, value string) error {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format, expected section.field (e.g., aws.profile)")
	}

	section, field := parts[0], parts[1]

	switch section {
	case "aws":
		return setAWSValue(&cfg.AWS, field, value)
	case "sync":
		return setSyncValue(&cfg.Sync, field, value)
	case "settings":
		return setSettingsValue(&cfg.Settings, field, value)
	default:
		return fmt.Errorf("unknown section: %s (valid: aws, sync, settings)", section)
	}
}

// setAWSValue sets an AWS configuration value.
func setAWSValue(aws *config.AWSConfig, field, value string) error {
	switch field {
	case "profile":
		aws.Profile = value
	case "region":
		aws.Region = value
	case "endpoint":
		aws.Endpoint = value
	default:
		return fmt.Errorf("unknown AWS field: %s (valid: profile, region, endpoint)", field)
	}
	return nil
}

// setSyncValue sets a sync configuration value.
func setSyncValue(sync *config.SyncConfig, field, value string) error {
	switch field {
	case "concurrency":
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("concurrency must be an integer: %w", err)
		}
		sync.Concurrency = i
	case "delete":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("delete must be a boolean: %w", err)
		}
		sync.Delete = b
	default:
		return fmt.Errorf("unknown sync field: %s (valid: concurrency, delete)", field)
	}
	return nil
}

// setSettingsValue sets a settings configuration value.
func setSettingsValue(settings *config.SettingsConfig, field, value string) error {
	switch field {
	case "verbose":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("verbose must be a boolean: %w", err)
		}
		settings.Verbose = b
	case "log_file":
		settings.LogFile = value
	case "check_updates":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("check_updates must be a boolean: %w", err)
		}
		settings.CheckUpdates = b
	default:
		return fmt.Errorf("unknown settings field: %s (valid: verbose, log_file, check_updates)", field)
	}
	return nil
}

// getConfigValue gets a configuration value using dot notation.
func getConfigValue(cfg *config.Config, key string) (string, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid key format, expected section.field (e.g., aws.profile)")
	}

	section, field := parts[0], parts[1]

	switch section {
	case "aws":
		return getAWSValue(&cfg.AWS, field)
	case "sync":
		return getSyncValue(&cfg.Sync, field)
	case "settings":
		return getSettingsValue(&cfg.Settings, field)
	default:
		return "", fmt.Errorf("unknown section: %s (valid: aws, sync, settings)", section)
	}
}

// getAWSValue gets an AWS configuration value.
func getAWSValue(aws *config.AWSConfig, field string) (string, error) {
	switch field {
	case "profile":
		return aws.Profile, nil
	case "region":
		return aws.Region, nil
	case "endpoint":
		return aws.Endpoint, nil
	default:
		return "", fmt.Errorf("unknown AWS field: %s (valid: profile, region, endpoint)", field)
	}
}

// getSyncValue gets a sync configuration value.
func getSyncValue(sync *config.SyncConfig, field string) (string, error) {
	switch field {
	case "concurrency":
		return strconv.Itoa(sync.Concurrency), nil
	case "delete":
		return strconv.FormatBool(sync.Delete), nil
	case "exclude":
		return strings.Join(sync.Exclude, ", "), nil
	default:
		return "", fmt.Errorf("unknown sync field: %s (valid: concurrency, delete, exclude)", field)
	}
}

// getSettingsValue gets a settings configuration value.
func getSettingsValue(settings *config.SettingsConfig, field string) (string, error) {
	switch field {
	case "verbose":
		return strconv.FormatBool(settings.Verbose), nil
	case "log_file":
		return settings.LogFile, nil
	case "check_updates":
		return strconv.FormatBool(settings.CheckUpdates), nil
	default:
		return "", fmt.Errorf("unknown settings field: %s (valid: verbose, log_file, check_updates)", field)
	}
}
