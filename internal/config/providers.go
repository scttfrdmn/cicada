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
	"strings"

	"github.com/joho/godotenv"
)

// ProviderConfig holds configuration for DOI providers
type ProviderConfig struct {
	// Zenodo configuration
	Zenodo ZenodoConfig `mapstructure:"zenodo" yaml:"zenodo"`

	// DataCite configuration
	DataCite DataCiteConfig `mapstructure:"datacite" yaml:"datacite"`

	// Default provider to use
	DefaultProvider string `mapstructure:"default_provider" yaml:"default_provider"`

	// Security settings
	Security SecurityConfig `mapstructure:"security" yaml:"security"`
}

// ZenodoConfig holds Zenodo API credentials
type ZenodoConfig struct {
	// API token
	Token string `mapstructure:"token" yaml:"token"`

	// Environment: "production" or "sandbox"
	Environment string `mapstructure:"environment" yaml:"environment"`
}

// DataCiteConfig holds DataCite API credentials
type DataCiteConfig struct {
	// Repository ID (e.g., "10.5072/FK2" or "CLIENT.MEMBER")
	RepositoryID string `mapstructure:"repository_id" yaml:"repository_id"`

	// Password for repository
	Password string `mapstructure:"password" yaml:"password"`

	// Environment: "production" or "sandbox"
	Environment string `mapstructure:"environment" yaml:"environment"`
}

// SecurityConfig holds security settings
type SecurityConfig struct {
	// Check file permissions (default: true)
	CheckPermissions bool `mapstructure:"check_permissions" yaml:"check_permissions"`

	// Warn about insecure configurations (default: true)
	WarnInsecure bool `mapstructure:"warn_insecure" yaml:"warn_insecure"`
}

// CredentialSource indicates where a credential came from
type CredentialSource string

const (
	SourceCommandLine  CredentialSource = "command-line flag"
	SourceEnvVar       CredentialSource = "environment variable"
	SourceConfigFile   CredentialSource = "config file"
	SourceDotEnv       CredentialSource = ".env file"
	SourceNotFound     CredentialSource = "not found"
)

// Credential represents a credential with its source
type Credential struct {
	Value  string
	Source CredentialSource
}

// ProviderCredentials manages provider credentials with proper precedence
type ProviderCredentials struct {
	// Command-line flags (highest precedence)
	CLIFlags map[string]string

	// Environment variables
	EnvVars map[string]string

	// Config file credentials
	ConfigFile *ProviderConfig

	// .env file credentials
	DotEnv map[string]string

	// Config file path (for security checks)
	ConfigFilePath string

	// .env file path
	DotEnvPath string
}

// NewProviderCredentials creates a new provider credentials manager
func NewProviderCredentials() *ProviderCredentials {
	return &ProviderCredentials{
		CLIFlags: make(map[string]string),
		EnvVars:  make(map[string]string),
		DotEnv:   make(map[string]string),
	}
}

// LoadFromEnvironment loads credentials from environment variables
func (pc *ProviderCredentials) LoadFromEnvironment() {
	// Zenodo credentials
	if token := os.Getenv("CICADA_ZENODO_TOKEN"); token != "" {
		pc.EnvVars["zenodo_token"] = token
	} else if token := os.Getenv("ZENODO_TOKEN"); token != "" {
		pc.EnvVars["zenodo_token"] = token
	}

	if env := os.Getenv("CICADA_ZENODO_SANDBOX"); env != "" {
		pc.EnvVars["zenodo_sandbox"] = env
	}

	// DataCite credentials
	if repoID := os.Getenv("CICADA_DATACITE_REPOSITORY_ID"); repoID != "" {
		pc.EnvVars["datacite_repository_id"] = repoID
	} else if repoID := os.Getenv("DATACITE_REPOSITORY_ID"); repoID != "" {
		pc.EnvVars["datacite_repository_id"] = repoID
	}

	if password := os.Getenv("CICADA_DATACITE_PASSWORD"); password != "" {
		pc.EnvVars["datacite_password"] = password
	} else if password := os.Getenv("DATACITE_PASSWORD"); password != "" {
		pc.EnvVars["datacite_password"] = password
	}

	if env := os.Getenv("CICADA_DATACITE_SANDBOX"); env != "" {
		pc.EnvVars["datacite_sandbox"] = env
	}

	// Default provider
	if provider := os.Getenv("CICADA_DEFAULT_PROVIDER"); provider != "" {
		pc.EnvVars["default_provider"] = provider
	}
}

// LoadFromConfigFile loads credentials from config file
func (pc *ProviderCredentials) LoadFromConfigFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Not an error, file just doesn't exist
	}

	pc.ConfigFilePath = path

	// Check file security
	if err := pc.checkFileSecurity(path); err != nil {
		return err
	}

	// Load the full config (includes provider config)
	_, err := Load(path)
	if err != nil {
		return fmt.Errorf("load config file: %w", err)
	}

	// Extract provider config if it exists
	// Note: This will be implemented when we extend the Config struct
	// For now, we'll load it separately
	// TODO: Integrate with main Config struct

	return nil
}

// LoadFromDotEnv loads credentials from .env file
func (pc *ProviderCredentials) LoadFromDotEnv(workDir string) error {
	dotenvPath := filepath.Join(workDir, ".env")

	// Check if file exists
	if _, err := os.Stat(dotenvPath); os.IsNotExist(err) {
		return nil // Not an error, file just doesn't exist
	}

	pc.DotEnvPath = dotenvPath

	// Check file security
	if err := pc.checkFileSecurity(dotenvPath); err != nil {
		return err
	}

	// Check if in git repository
	if err := pc.checkGitRepository(dotenvPath); err != nil {
		// This is a warning, not an error
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	// Load .env file
	envMap, err := godotenv.Read(dotenvPath)
	if err != nil {
		return fmt.Errorf("read .env file: %w", err)
	}

	// Store in DotEnv map
	for key, value := range envMap {
		// Normalize key to lowercase
		normalizedKey := strings.ToLower(strings.TrimPrefix(key, "CICADA_"))
		pc.DotEnv[normalizedKey] = value
	}

	return nil
}

// GetCredential retrieves a credential with proper precedence
// key should be like "zenodo_token", "datacite_repository_id", etc.
func (pc *ProviderCredentials) GetCredential(key string) Credential {
	// Normalize key
	key = strings.ToLower(key)

	// 1. Check command-line flags (highest precedence)
	if val, ok := pc.CLIFlags[key]; ok && val != "" {
		return Credential{Value: val, Source: SourceCommandLine}
	}

	// 2. Check environment variables
	if val, ok := pc.EnvVars[key]; ok && val != "" {
		return Credential{Value: val, Source: SourceEnvVar}
	}

	// 3. Check config file
	if pc.ConfigFile != nil {
		val := pc.getFromConfigFile(key)
		if val != "" {
			return Credential{Value: val, Source: SourceConfigFile}
		}
	}

	// 4. Check .env file (lowest precedence)
	if val, ok := pc.DotEnv[key]; ok && val != "" {
		return Credential{Value: val, Source: SourceDotEnv}
	}

	// Not found
	return Credential{Value: "", Source: SourceNotFound}
}

// getFromConfigFile retrieves value from loaded config file
func (pc *ProviderCredentials) getFromConfigFile(key string) string {
	if pc.ConfigFile == nil {
		return ""
	}

	switch key {
	case "zenodo_token":
		return pc.ConfigFile.Zenodo.Token
	case "zenodo_environment", "zenodo_sandbox":
		return pc.ConfigFile.Zenodo.Environment
	case "datacite_repository_id":
		return pc.ConfigFile.DataCite.RepositoryID
	case "datacite_password":
		return pc.ConfigFile.DataCite.Password
	case "datacite_environment", "datacite_sandbox":
		return pc.ConfigFile.DataCite.Environment
	case "default_provider":
		return pc.ConfigFile.DefaultProvider
	default:
		return ""
	}
}

// SetCLIFlag sets a command-line flag value
func (pc *ProviderCredentials) SetCLIFlag(key, value string) {
	pc.CLIFlags[strings.ToLower(key)] = value
}

// checkFileSecurity checks file permissions and warns if insecure
func (pc *ProviderCredentials) checkFileSecurity(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	mode := fileInfo.Mode()
	perm := mode.Perm()

	// Check if group or others have read permission (Unix only)
	// Permissions: 0600 means owner read/write only
	if perm&0077 != 0 {
		return fmt.Errorf(`insecure permissions on %s: %04o

  Your config file contains sensitive credentials but has insecure permissions.
  Other users on this system can read your credentials.

  Fix with: chmod 600 %s

  Current:  %s (%04o)
  Required: -rw------- (0600)`,
			path, perm, path, mode, perm)
	}

	return nil
}

// checkGitRepository checks if file is in a git repository and not in .gitignore
func (pc *ProviderCredentials) checkGitRepository(path string) error {
	// Walk up directory tree looking for .git
	dir := filepath.Dir(path)
	gitDir := ""

	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			gitDir = dir
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, no .git found
			return nil
		}
		dir = parent
	}

	// Found .git directory, check if file is in .gitignore
	gitignorePath := filepath.Join(gitDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// No .gitignore, warn user
		return fmt.Errorf(`config file is in a git repository without .gitignore

  File: %s
  Repo: %s

  Your config file contains credentials but is in a git repository.
  Ensure it's in .gitignore to prevent accidental commit.

  Add to .gitignore:
    echo "%s" >> %s

  Or move config to: ~/.config/cicada/config.yaml`,
			path, gitDir, filepath.Base(path), gitignorePath)
	}

	// TODO: Check if file is actually ignored (more complex)
	// For now, just warn if .git exists and file is named .env
	if filepath.Base(path) == ".env" {
		return fmt.Errorf(`ensure .env file is in .gitignore

  File: %s
  Repo: %s

  Add to .gitignore:
    echo ".env" >> %s`,
			path, gitDir, gitignorePath)
	}

	return nil
}

// ValidateZenodoToken validates a Zenodo token format
func ValidateZenodoToken(token string) error {
	if token == "" {
		return fmt.Errorf("Zenodo token is empty")
	}
	if len(token) < 20 {
		return fmt.Errorf("Zenodo token too short (expected 40+ characters, got %d)", len(token))
	}
	if strings.Contains(token, " ") {
		return fmt.Errorf("Zenodo token contains whitespace")
	}
	return nil
}

// ValidateDataCiteRepositoryID validates a DataCite repository ID format
func ValidateDataCiteRepositoryID(repoID string) error {
	if repoID == "" {
		return fmt.Errorf("DataCite repository ID is empty")
	}
	// Repository ID should be like "10.5072/FK2" or "CLIENT.MEMBER"
	if !strings.Contains(repoID, ".") && !strings.Contains(repoID, "/") {
		return fmt.Errorf("DataCite repository ID format invalid (expected '10.XXXX/YYY' or 'CLIENT.MEMBER')")
	}
	return nil
}

// ValidateDataCitePassword validates a DataCite password
func ValidateDataCitePassword(password string) error {
	if password == "" {
		return fmt.Errorf("DataCite password is empty")
	}
	if len(password) < 8 {
		return fmt.Errorf("DataCite password too short (minimum 8 characters)")
	}
	return nil
}

// RedactToken redacts a token for display (shows last 5 chars only)
func RedactToken(token string) string {
	if token == "" {
		return "(empty)"
	}
	if len(token) <= 5 {
		return "*****"
	}
	return "****" + token[len(token)-5:]
}
