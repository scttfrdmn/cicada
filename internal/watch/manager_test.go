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
	"testing"
	"time"
)

func TestManager_NewManager(t *testing.T) {
	manager := NewManager()

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	if manager.watchers == nil {
		t.Error("watchers map not initialized")
	}

	if len(manager.watchers) != 0 {
		t.Errorf("watchers map has %d entries, want 0", len(manager.watchers))
	}
}

func TestManager_List_Empty(t *testing.T) {
	manager := NewManager()

	statuses := manager.List()

	if statuses == nil {
		t.Fatal("List() returned nil")
	}

	if len(statuses) != 0 {
		t.Errorf("List() returned %d statuses, want 0", len(statuses))
	}
}

func TestManager_Get_NotFound(t *testing.T) {
	manager := NewManager()

	watcher, exists := manager.Get("nonexistent")

	if exists {
		t.Error("Get() returned exists=true for nonexistent watch")
	}

	if watcher != nil {
		t.Error("Get() returned non-nil watcher for nonexistent watch")
	}
}

func TestManager_Remove_NotFound(t *testing.T) {
	manager := NewManager()

	err := manager.Remove("nonexistent")

	if err == nil {
		t.Error("Remove() returned nil error for nonexistent watch")
	}
}

func TestManager_StopAll_Empty(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	err := manager.StopAll(ctx)

	if err != nil {
		t.Errorf("StopAll() returned error: %v", err)
	}

	statuses := manager.List()
	if len(statuses) != 0 {
		t.Errorf("After StopAll(), %d watches remain, want 0", len(statuses))
	}
}

func TestConfig_ZeroValues(t *testing.T) {
	config := Config{}

	if config.Source != "" {
		t.Errorf("Source = %s, want empty", config.Source)
	}

	if config.Destination != "" {
		t.Errorf("Destination = %s, want empty", config.Destination)
	}

	if config.DebounceDelay != 0 {
		t.Errorf("DebounceDelay = %v, want 0", config.DebounceDelay)
	}

	if config.MinAge != 0 {
		t.Errorf("MinAge = %v, want 0", config.MinAge)
	}
}

func TestWatchStatus_ErrorTracking(t *testing.T) {
	status := WatchStatus{
		Source:      "/test",
		Destination: "s3://test",
		Active:      true,
		StartedAt:   time.Now(),
		ErrorCount:  5,
		LastError:   "test error",
	}

	if status.ErrorCount != 5 {
		t.Errorf("ErrorCount = %d, want 5", status.ErrorCount)
	}

	if status.LastError != "test error" {
		t.Errorf("LastError = %s, want 'test error'", status.LastError)
	}
}

func TestWatchStatus_SyncTracking(t *testing.T) {
	now := time.Now()
	status := WatchStatus{
		Source:      "/test",
		Destination: "s3://test",
		Active:      true,
		StartedAt:   now,
		LastSync:    now.Add(1 * time.Minute),
		FilesSynced: 100,
		BytesSynced: 1024 * 1024,
	}

	if status.FilesSynced != 100 {
		t.Errorf("FilesSynced = %d, want 100", status.FilesSynced)
	}

	if status.BytesSynced != 1024*1024 {
		t.Errorf("BytesSynced = %d, want %d", status.BytesSynced, 1024*1024)
	}

	if status.LastSync.Before(status.StartedAt) {
		t.Error("LastSync is before StartedAt")
	}
}
