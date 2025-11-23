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
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	cicadasync "github.com/scttfrdmn/cicada/internal/sync"
)

// Watcher monitors a directory and syncs changes.
type Watcher struct {
	config    Config
	fsWatcher *fsnotify.Watcher
	debouncer *Debouncer
	engine    *cicadasync.Engine
	status    WatchStatus
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// New creates a new watcher.
func New(config Config, engine *cicadasync.Engine) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fsnotify watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := &Watcher{
		config:    config,
		fsWatcher: fsWatcher,
		engine:    engine,
		ctx:       ctx,
		cancel:    cancel,
		status: WatchStatus{
			Source:      config.Source,
			Destination: config.Destination,
			Active:      false,
		},
	}

	// Create debouncer that triggers sync
	w.debouncer = NewDebouncer(config.DebounceDelay, w.triggerSync)

	return w, nil
}

// Start begins watching the configured directory.
func (w *Watcher) Start() error {
	w.mu.Lock()
	if w.status.Active {
		w.mu.Unlock()
		return fmt.Errorf("watcher already active")
	}
	w.status.Active = true
	w.status.StartedAt = time.Now()
	w.mu.Unlock()

	// Add all directories recursively
	if err := w.addRecursive(w.config.Source); err != nil {
		w.status.Active = false
		return fmt.Errorf("add watch directories: %w", err)
	}

	// Perform initial sync if configured
	if w.config.SyncOnStart {
		w.triggerSync()
	}

	// Start event loop
	w.wg.Add(1)
	go w.eventLoop()

	return nil
}

// Stop stops the watcher gracefully.
func (w *Watcher) Stop() error {
	w.mu.Lock()
	if !w.status.Active {
		w.mu.Unlock()
		return fmt.Errorf("watcher not active")
	}
	w.status.Active = false
	w.mu.Unlock()

	// Stop debouncer and cancel context
	w.debouncer.Stop()
	w.cancel()

	// Wait for event loop to finish
	w.wg.Wait()

	// Close fsnotify watcher
	return w.fsWatcher.Close()
}

// Status returns current watcher status.
func (w *Watcher) Status() WatchStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.status
}

// eventLoop processes file system events.
func (w *Watcher) eventLoop() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			w.recordError(err)
		}
	}
}

// handleEvent processes a single file system event.
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Check if file matches exclude patterns
	if w.shouldExclude(event.Name) {
		return
	}

	// Handle different event types
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		// If directory was created, add it to watch list
		if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
			_ = w.addRecursive(event.Name)
		}
		w.debouncer.Trigger()

	case event.Op&fsnotify.Write == fsnotify.Write:
		w.debouncer.Trigger()

	case event.Op&fsnotify.Remove == fsnotify.Remove:
		w.debouncer.Trigger()

	case event.Op&fsnotify.Rename == fsnotify.Rename:
		w.debouncer.Trigger()
	}
}

// triggerSync performs the actual sync operation.
func (w *Watcher) triggerSync() {
	// Wait for min-age if configured
	if w.config.MinAge > 0 {
		time.Sleep(w.config.MinAge)
	}

	// Perform sync
	err := w.engine.Sync(w.ctx, w.config.Source, w.config.Destination)

	w.mu.Lock()
	defer w.mu.Unlock()

	if err != nil {
		w.status.ErrorCount++
		w.status.LastError = err.Error()
	} else {
		w.status.LastSync = time.Now()
		// TODO: Update files/bytes synced from engine feedback
	}
}

// addRecursive adds a directory and all subdirectories to watch.
func (w *Watcher) addRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only watch directories
		if !info.IsDir() {
			return nil
		}

		// Skip excluded directories
		if w.shouldExclude(path) {
			return filepath.SkipDir
		}

		// Add to watcher
		if err := w.fsWatcher.Add(path); err != nil {
			return fmt.Errorf("add watch on %s: %w", path, err)
		}

		return nil
	})
}

// shouldExclude checks if path matches any exclude pattern.
func (w *Watcher) shouldExclude(path string) bool {
	for _, pattern := range w.config.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

// recordError updates error statistics.
func (w *Watcher) recordError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.status.ErrorCount++
	w.status.LastError = err.Error()
}
