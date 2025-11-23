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
	"time"
)

// Config holds configuration for a watch operation.
type Config struct {
	// Source is the local directory to watch
	Source string

	// Destination is the S3 URI to sync to
	Destination string

	// DebounceDelay is how long to wait after last change before syncing
	DebounceDelay time.Duration

	// MinAge is minimum file age before syncing (prevents partial files)
	MinAge time.Duration

	// DeleteSource removes source files after successful sync
	DeleteSource bool

	// SyncOnStart performs initial sync when watch starts
	SyncOnStart bool

	// Exclude patterns for files to ignore
	ExcludePatterns []string

	// CronSchedule for periodic syncs (optional)
	CronSchedule string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		DebounceDelay:   5 * time.Second,
		MinAge:          10 * time.Second,
		DeleteSource:    false,
		SyncOnStart:     true,
		ExcludePatterns: []string{".git/**", ".DS_Store", "*.tmp", "*.swp"},
	}
}

// WatchStatus represents the current state of a watcher.
type WatchStatus struct {
	Source       string    `json:"source"`
	Destination  string    `json:"destination"`
	Active       bool      `json:"active"`
	StartedAt    time.Time `json:"started_at"`
	LastSync     time.Time `json:"last_sync"`
	FilesSynced  int64     `json:"files_synced"`
	BytesSynced  int64     `json:"bytes_synced"`
	ErrorCount   int       `json:"error_count"`
	LastError    string    `json:"last_error,omitempty"`
}
