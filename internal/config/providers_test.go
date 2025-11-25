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
	"os"
	"path/filepath"
	"testing"
)

func TestProviderCredentialsPrecedence(t *testing.T) {
	pc := NewProviderCredentials()

	// Set up different sources with different values
	pc.SetCLIFlag("zenodo_token", "cli-token")
	pc.EnvVars["zenodo_token"] = "env-token"
	pc.ConfigFile = &ProviderConfig{
		Zenodo: ZenodoConfig{
			Token: "config-token",
		},
	}
	pc.DotEnv["zenodo_token"] = "dotenv-token"

	// Test precedence: CLI > Env > Config > DotEnv
	cred := pc.GetCredential("zenodo_token")
	if cred.Value != "cli-token" {
		t.Errorf("expected CLI token, got %s", cred.Value)
	}
	if cred.Source != SourceCommandLine {
		t.Errorf("expected CLI source, got %s", cred.Source)
	}

	// Remove CLI flag
	delete(pc.CLIFlags, "zenodo_token")
	cred = pc.GetCredential("zenodo_token")
	if cred.Value != "env-token" {
		t.Errorf("expected env token, got %s", cred.Value)
	}
	if cred.Source != SourceEnvVar {
		t.Errorf("expected env source, got %s", cred.Source)
	}

	// Remove env var
	delete(pc.EnvVars, "zenodo_token")
	cred = pc.GetCredential("zenodo_token")
	if cred.Value != "config-token" {
		t.Errorf("expected config token, got %s", cred.Value)
	}
	if cred.Source != SourceConfigFile {
		t.Errorf("expected config source, got %s", cred.Source)
	}

	// Remove config file
	pc.ConfigFile = nil
	cred = pc.GetCredential("zenodo_token")
	if cred.Value != "dotenv-token" {
		t.Errorf("expected dotenv token, got %s", cred.Value)
	}
	if cred.Source != SourceDotEnv {
		t.Errorf("expected dotenv source, got %s", cred.Source)
	}

	// Remove everything
	delete(pc.DotEnv, "zenodo_token")
	cred = pc.GetCredential("zenodo_token")
	if cred.Value != "" {
		t.Errorf("expected empty value, got %s", cred.Value)
	}
	if cred.Source != SourceNotFound {
		t.Errorf("expected not found source, got %s", cred.Source)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Set environment variables
	if err := os.Setenv("CICADA_ZENODO_TOKEN", "test-token-1"); err != nil {
		t.Fatalf("failed to set CICADA_ZENODO_TOKEN: %v", err)
	}
	if err := os.Setenv("ZENODO_TOKEN", "test-token-2"); err != nil {
		t.Fatalf("failed to set ZENODO_TOKEN: %v", err)
	}
	if err := os.Setenv("CICADA_DATACITE_REPOSITORY_ID", "10.5072/TEST"); err != nil {
		t.Fatalf("failed to set CICADA_DATACITE_REPOSITORY_ID: %v", err)
	}
	if err := os.Setenv("DATACITE_PASSWORD", "test-password"); err != nil {
		t.Fatalf("failed to set DATACITE_PASSWORD: %v", err)
	}
	defer func() {
		// Cleanup - errors here are not critical
		_ = os.Unsetenv("CICADA_ZENODO_TOKEN")
		_ = os.Unsetenv("ZENODO_TOKEN")
		_ = os.Unsetenv("CICADA_DATACITE_REPOSITORY_ID")
		_ = os.Unsetenv("DATACITE_PASSWORD")
	}()

	pc := NewProviderCredentials()
	pc.LoadFromEnvironment()

	// Check CICADA_ prefix takes precedence
	if pc.EnvVars["zenodo_token"] != "test-token-1" {
		t.Errorf("expected test-token-1, got %s", pc.EnvVars["zenodo_token"])
	}

	// Check DataCite vars
	if pc.EnvVars["datacite_repository_id"] != "10.5072/TEST" {
		t.Errorf("expected 10.5072/TEST, got %s", pc.EnvVars["datacite_repository_id"])
	}
	if pc.EnvVars["datacite_password"] != "test-password" {
		t.Errorf("expected test-password, got %s", pc.EnvVars["datacite_password"])
	}
}

func TestLoadFromDotEnv(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create .env file
	dotenvPath := filepath.Join(tmpDir, ".env")
	dotenvContent := `CICADA_ZENODO_TOKEN=dotenv-token
CICADA_DATACITE_REPOSITORY_ID=10.5072/DOTENV
DATACITE_PASSWORD=dotenv-password
`
	if err := os.WriteFile(dotenvPath, []byte(dotenvContent), 0600); err != nil {
		t.Fatal(err)
	}

	pc := NewProviderCredentials()
	if err := pc.LoadFromDotEnv(tmpDir); err != nil {
		t.Fatalf("LoadFromDotEnv failed: %v", err)
	}

	// Check values
	if pc.DotEnv["zenodo_token"] != "dotenv-token" {
		t.Errorf("expected dotenv-token, got %s", pc.DotEnv["zenodo_token"])
	}
	if pc.DotEnv["datacite_repository_id"] != "10.5072/DOTENV" {
		t.Errorf("expected 10.5072/DOTENV, got %s", pc.DotEnv["datacite_repository_id"])
	}
}

func TestCheckFileSecurityInsecurePermissions(t *testing.T) {
	// Create temporary file with insecure permissions
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	pc := NewProviderCredentials()
	err := pc.checkFileSecurity(tmpFile)
	if err == nil {
		t.Error("expected error for insecure permissions, got nil")
	}
	if err != nil && !contains(err.Error(), "insecure permissions") {
		t.Errorf("expected 'insecure permissions' error, got: %v", err)
	}
}

func TestCheckFileSecuritySecurePermissions(t *testing.T) {
	// Create temporary file with secure permissions
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(tmpFile, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}

	pc := NewProviderCredentials()
	err := pc.checkFileSecurity(tmpFile)
	if err != nil {
		t.Errorf("expected no error for secure permissions, got: %v", err)
	}
}

func TestValidateZenodoToken(t *testing.T) {
	tests := []struct {
		token   string
		wantErr bool
	}{
		{"", true},                                          // Empty
		{"short", true},                                     // Too short
		{"token with spaces", true},                         // Has spaces
		{"abcdefghijklmnopqrstuvwxyz1234567890123456", false}, // Valid
	}

	for _, tt := range tests {
		err := ValidateZenodoToken(tt.token)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateZenodoToken(%q) error = %v, wantErr %v", tt.token, err, tt.wantErr)
		}
	}
}

func TestValidateDataCiteRepositoryID(t *testing.T) {
	tests := []struct {
		repoID  string
		wantErr bool
	}{
		{"", true},                 // Empty
		{"invalid", true},          // No dot or slash
		{"10.5072/TEST", false},    // Valid prefix format
		{"CLIENT.MEMBER", false},   // Valid client format
	}

	for _, tt := range tests {
		err := ValidateDataCiteRepositoryID(tt.repoID)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateDataCiteRepositoryID(%q) error = %v, wantErr %v", tt.repoID, err, tt.wantErr)
		}
	}
}

func TestValidateDataCitePassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
	}{
		{"", true},          // Empty
		{"short", true},     // Too short
		{"password123", false}, // Valid
	}

	for _, tt := range tests {
		err := ValidateDataCitePassword(tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateDataCitePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
		}
	}
}

func TestRedactToken(t *testing.T) {
	tests := []struct {
		token string
		want  string
	}{
		{"", "(empty)"},
		{"abc", "*****"},
		{"abcdef", "****bcdef"},  // Shows last 5 chars
		{"abcdefghijklmnop", "****lmnop"},
	}

	for _, tt := range tests {
		got := RedactToken(tt.token)
		if got != tt.want {
			t.Errorf("RedactToken(%q) = %q, want %q", tt.token, got, tt.want)
		}
	}
}

func TestGetCredentialAllProviders(t *testing.T) {
	pc := NewProviderCredentials()
	pc.ConfigFile = &ProviderConfig{
		Zenodo: ZenodoConfig{
			Token:       "zenodo-token",
			Environment: "production",
		},
		DataCite: DataCiteConfig{
			RepositoryID: "10.5072/TEST",
			Password:     "datacite-password",
			Environment:  "sandbox",
		},
		DefaultProvider: "zenodo",
	}

	tests := []struct {
		key       string
		wantValue string
		wantSource CredentialSource
	}{
		{"zenodo_token", "zenodo-token", SourceConfigFile},
		{"zenodo_environment", "production", SourceConfigFile},
		{"datacite_repository_id", "10.5072/TEST", SourceConfigFile},
		{"datacite_password", "datacite-password", SourceConfigFile},
		{"datacite_environment", "sandbox", SourceConfigFile},
		{"default_provider", "zenodo", SourceConfigFile},
		{"nonexistent", "", SourceNotFound},
	}

	for _, tt := range tests {
		cred := pc.GetCredential(tt.key)
		if cred.Value != tt.wantValue {
			t.Errorf("GetCredential(%q) value = %q, want %q", tt.key, cred.Value, tt.wantValue)
		}
		if cred.Source != tt.wantSource {
			t.Errorf("GetCredential(%q) source = %q, want %q", tt.key, cred.Source, tt.wantSource)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
