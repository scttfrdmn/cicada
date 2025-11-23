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

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Version != "1" {
		t.Errorf("Version = %s, want 1", cfg.Version)
	}

	if cfg.AWS.Profile != "default" {
		t.Errorf("AWS.Profile = %s, want default", cfg.AWS.Profile)
	}

	if cfg.Sync.Concurrency != 4 {
		t.Errorf("Sync.Concurrency = %d, want 4", cfg.Sync.Concurrency)
	}

	if cfg.Sync.Delete {
		t.Error("Sync.Delete = true, want false")
	}

	if len(cfg.Sync.Exclude) == 0 {
		t.Error("Sync.Exclude is empty, want default patterns")
	}

	if len(cfg.Watches) != 0 {
		t.Errorf("Watches has %d entries, want 0", len(cfg.Watches))
	}

	if cfg.Settings.Verbose {
		t.Error("Settings.Verbose = true, want false")
	}

	if !cfg.Settings.CheckUpdates {
		t.Error("Settings.CheckUpdates = false, want true")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create test config
	cfg := DefaultConfig()
	cfg.AWS.Profile = "test-profile"
	cfg.AWS.Region = "us-west-2"
	cfg.Sync.Concurrency = 8
	cfg.Sync.Delete = true
	cfg.Settings.Verbose = true

	// Add a watch
	cfg.Watches = append(cfg.Watches, WatchConfig{
		ID:              "test-watch",
		Source:          "/local/path",
		Destination:     "s3://bucket/prefix",
		DebounceSeconds: 5,
		MinAgeSeconds:   10,
		DeleteSource:    false,
		SyncOnStart:     true,
		Exclude:         []string{"*.tmp"},
		Enabled:         true,
	})

	// Save config
	if err := Save(cfg, configPath); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load config
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Verify loaded config matches
	if loaded.Version != cfg.Version {
		t.Errorf("Version = %s, want %s", loaded.Version, cfg.Version)
	}

	if loaded.AWS.Profile != cfg.AWS.Profile {
		t.Errorf("AWS.Profile = %s, want %s", loaded.AWS.Profile, cfg.AWS.Profile)
	}

	if loaded.AWS.Region != cfg.AWS.Region {
		t.Errorf("AWS.Region = %s, want %s", loaded.AWS.Region, cfg.AWS.Region)
	}

	if loaded.Sync.Concurrency != cfg.Sync.Concurrency {
		t.Errorf("Sync.Concurrency = %d, want %d", loaded.Sync.Concurrency, cfg.Sync.Concurrency)
	}

	if loaded.Sync.Delete != cfg.Sync.Delete {
		t.Errorf("Sync.Delete = %v, want %v", loaded.Sync.Delete, cfg.Sync.Delete)
	}

	if loaded.Settings.Verbose != cfg.Settings.Verbose {
		t.Errorf("Settings.Verbose = %v, want %v", loaded.Settings.Verbose, cfg.Settings.Verbose)
	}

	if len(loaded.Watches) != 1 {
		t.Fatalf("Watches has %d entries, want 1", len(loaded.Watches))
	}

	watch := loaded.Watches[0]
	if watch.ID != "test-watch" {
		t.Errorf("Watch.ID = %s, want test-watch", watch.ID)
	}

	if watch.Source != "/local/path" {
		t.Errorf("Watch.Source = %s, want /local/path", watch.Source)
	}

	if watch.Destination != "s3://bucket/prefix" {
		t.Errorf("Watch.Destination = %s, want s3://bucket/prefix", watch.Destination)
	}
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Load() on nonexistent file returned nil error")
	}
}

func TestLoadOrDefault_Nonexistent(t *testing.T) {
	// Temporarily override config path
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("Setenv() error: %v", err)
	}
	defer func() { _ = os.Setenv("HOME", originalHome) }()

	// LoadOrDefault should return default config when file doesn't exist
	cfg, err := LoadOrDefault()
	if err != nil {
		t.Fatalf("LoadOrDefault() error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadOrDefault() returned nil config")
	}

	// Verify it's the default config
	if cfg.Version != "1" {
		t.Errorf("Version = %s, want 1", cfg.Version)
	}

	if cfg.AWS.Profile != "default" {
		t.Errorf("AWS.Profile = %s, want default", cfg.AWS.Profile)
	}
}

func TestLoadOrDefault_Existing(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Set HOME to tmpDir
	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("Setenv() error: %v", err)
	}
	defer func() { _ = os.Setenv("HOME", originalHome) }()

	// Create config directory
	configDir := filepath.Join(tmpDir, ".cicada")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}

	// Create config file with custom values
	configPath := filepath.Join(configDir, "config.yaml")
	cfg := DefaultConfig()
	cfg.AWS.Profile = "custom-profile"
	cfg.AWS.Region = "eu-west-1"

	if err := Save(cfg, configPath); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// LoadOrDefault should load the existing config
	loaded, err := LoadOrDefault()
	if err != nil {
		t.Fatalf("LoadOrDefault() error: %v", err)
	}

	if loaded.AWS.Profile != "custom-profile" {
		t.Errorf("AWS.Profile = %s, want custom-profile", loaded.AWS.Profile)
	}

	if loaded.AWS.Region != "eu-west-1" {
		t.Errorf("AWS.Region = %s, want eu-west-1", loaded.AWS.Region)
	}
}

func TestConfigPath(t *testing.T) {
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error: %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() returned empty string")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("ConfigPath() = %s is not absolute", path)
	}

	// Should end with .cicada/config.yaml
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("ConfigPath() basename = %s, want config.yaml", filepath.Base(path))
	}

	if filepath.Base(filepath.Dir(path)) != ".cicada" {
		t.Errorf("ConfigPath() parent dir = %s, want .cicada", filepath.Base(filepath.Dir(path)))
	}
}

func TestConfigDir(t *testing.T) {
	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() error: %v", err)
	}

	if dir == "" {
		t.Error("ConfigDir() returned empty string")
	}

	if !filepath.IsAbs(dir) {
		t.Errorf("ConfigDir() = %s is not absolute", dir)
	}

	// Should end with .cicada
	if filepath.Base(dir) != ".cicada" {
		t.Errorf("ConfigDir() basename = %s, want .cicada", filepath.Base(dir))
	}
}

func TestWatchConfig(t *testing.T) {
	cfg := WatchConfig{
		ID:              "test",
		Source:          "/src",
		Destination:     "s3://dst",
		DebounceSeconds: 5,
		MinAgeSeconds:   10,
		DeleteSource:    true,
		SyncOnStart:     false,
		Exclude:         []string{"*.log", "*.tmp"},
		Enabled:         true,
	}

	if cfg.ID != "test" {
		t.Errorf("ID = %s, want test", cfg.ID)
	}

	if cfg.DebounceSeconds != 5 {
		t.Errorf("DebounceSeconds = %d, want 5", cfg.DebounceSeconds)
	}

	if !cfg.DeleteSource {
		t.Error("DeleteSource = false, want true")
	}

	if cfg.SyncOnStart {
		t.Error("SyncOnStart = true, want false")
	}

	if len(cfg.Exclude) != 2 {
		t.Errorf("Exclude has %d entries, want 2", len(cfg.Exclude))
	}
}
