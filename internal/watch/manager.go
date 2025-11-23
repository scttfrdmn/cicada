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

package watch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/scttfrdmn/cicada/internal/config"
	cicadasync "github.com/scttfrdmn/cicada/internal/sync"
)

// Manager manages multiple watchers.
type Manager struct {
	watchers map[string]*Watcher
	mu       sync.RWMutex
}

// NewManager creates a new watch manager.
func NewManager() *Manager {
	return &Manager{
		watchers: make(map[string]*Watcher),
	}
}

// Add creates and starts a new watcher.
func (m *Manager) Add(id string, config Config, srcBackend, dstBackend cicadasync.Backend) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.watchers[id]; exists {
		return fmt.Errorf("watch %s already exists", id)
	}

	// Create sync engine for this watch
	engine := cicadasync.NewEngine(srcBackend, dstBackend, cicadasync.SyncOptions{
		Concurrency: 4,
		ProgressFunc: func(update cicadasync.ProgressUpdate) {
			// TODO: Log progress
		},
	})

	// Create watcher
	watcher, err := New(config, engine)
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}

	// Start watching
	if err := watcher.Start(); err != nil {
		return fmt.Errorf("start watcher: %w", err)
	}

	m.watchers[id] = watcher
	return nil
}

// Remove stops and removes a watcher.
func (m *Manager) Remove(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	watcher, exists := m.watchers[id]
	if !exists {
		return fmt.Errorf("watch %s not found", id)
	}

	if err := watcher.Stop(); err != nil {
		return fmt.Errorf("stop watcher: %w", err)
	}

	delete(m.watchers, id)
	return nil
}

// Get retrieves a watcher by ID.
func (m *Manager) Get(id string) (*Watcher, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	watcher, exists := m.watchers[id]
	return watcher, exists
}

// List returns all active watchers.
func (m *Manager) List() map[string]WatchStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make(map[string]WatchStatus)
	for id, watcher := range m.watchers {
		statuses[id] = watcher.Status()
	}
	return statuses
}

// StopAll stops all watchers.
func (m *Manager) StopAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error
	for id, watcher := range m.watchers {
		if err := watcher.Stop(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("stop watcher %s: %w", id, err)
		}
	}

	m.watchers = make(map[string]*Watcher)
	return firstErr
}

// SaveConfig persists all watches to the configuration file.
func (m *Manager) SaveConfig() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Load current config
	cfg, err := config.LoadOrDefault()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Convert watches to config format
	cfg.Watches = make([]config.WatchConfig, 0, len(m.watchers))
	for id, watcher := range m.watchers {
		status := watcher.Status()
		watchConfig := config.WatchConfig{
			ID:              id,
			Source:          status.Source,
			Destination:     status.Destination,
			DebounceSeconds: int(watcher.config.DebounceDelay.Seconds()),
			MinAgeSeconds:   int(watcher.config.MinAge.Seconds()),
			DeleteSource:    watcher.config.DeleteSource,
			SyncOnStart:     watcher.config.SyncOnStart,
			Exclude:         watcher.config.ExcludePatterns,
			Enabled:         status.Active,
		}
		cfg.Watches = append(cfg.Watches, watchConfig)
	}

	// Save config
	path, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("get config path: %w", err)
	}

	if err := config.Save(cfg, path); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// AddWatch creates, starts, and persists a new watcher.
func (m *Manager) AddWatch(id string, cfg Config, srcBackend, dstBackend cicadasync.Backend) error {
	if err := m.Add(id, cfg, srcBackend, dstBackend); err != nil {
		return err
	}

	if err := m.SaveConfig(); err != nil {
		// Rollback the add operation
		_ = m.Remove(id)
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// RemoveWatch stops, removes, and unpersists a watcher.
func (m *Manager) RemoveWatch(id string) error {
	if err := m.Remove(id); err != nil {
		return err
	}

	if err := m.SaveConfig(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// LoadFromConfig loads and starts all enabled watches from the configuration file.
func (m *Manager) LoadFromConfig(createBackend func(context.Context, string) (cicadasync.Backend, string, error)) error {
	// Load config
	cfg, err := config.LoadOrDefault()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx := context.Background()

	// Start each enabled watch
	for _, watchConfig := range cfg.Watches {
		if !watchConfig.Enabled {
			continue
		}

		// Create backends
		srcBackend, srcPath, err := createBackend(ctx, watchConfig.Source)
		if err != nil {
			return fmt.Errorf("create source backend for %s: %w", watchConfig.ID, err)
		}

		dstBackend, dstPath, err := createBackend(ctx, watchConfig.Destination)
		if err != nil {
			return fmt.Errorf("create destination backend for %s: %w", watchConfig.ID, err)
		}

		// Convert config
		config := Config{
			Source:          srcPath,
			Destination:     dstPath,
			DebounceDelay:   time.Duration(watchConfig.DebounceSeconds) * time.Second,
			MinAge:          time.Duration(watchConfig.MinAgeSeconds) * time.Second,
			DeleteSource:    watchConfig.DeleteSource,
			SyncOnStart:     watchConfig.SyncOnStart,
			ExcludePatterns: watchConfig.Exclude,
		}

		// Add watch (without persisting again)
		if err := m.Add(watchConfig.ID, config, srcBackend, dstBackend); err != nil {
			return fmt.Errorf("add watch %s: %w", watchConfig.ID, err)
		}
	}

	return nil
}
