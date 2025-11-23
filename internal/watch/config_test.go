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
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.DebounceDelay != 5*time.Second {
		t.Errorf("DebounceDelay = %v, want 5s", config.DebounceDelay)
	}

	if config.MinAge != 10*time.Second {
		t.Errorf("MinAge = %v, want 10s", config.MinAge)
	}

	if config.DeleteSource {
		t.Error("DeleteSource = true, want false")
	}

	if !config.SyncOnStart {
		t.Error("SyncOnStart = false, want true")
	}

	if len(config.ExcludePatterns) == 0 {
		t.Error("ExcludePatterns is empty, want default patterns")
	}

	// Check for common exclude patterns
	expectedPatterns := []string{".git/**", ".DS_Store", "*.tmp", "*.swp"}
	for _, pattern := range expectedPatterns {
		found := false
		for _, p := range config.ExcludePatterns {
			if p == pattern {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ExcludePatterns missing %s", pattern)
		}
	}
}

func TestConfig_CustomValues(t *testing.T) {
	config := Config{
		Source:          "/local/path",
		Destination:     "s3://bucket/prefix",
		DebounceDelay:   1 * time.Second,
		MinAge:          2 * time.Second,
		DeleteSource:    true,
		SyncOnStart:     false,
		ExcludePatterns: []string{"*.log"},
		CronSchedule:    "0 * * * *",
	}

	if config.Source != "/local/path" {
		t.Errorf("Source = %s, want /local/path", config.Source)
	}

	if config.Destination != "s3://bucket/prefix" {
		t.Errorf("Destination = %s, want s3://bucket/prefix", config.Destination)
	}

	if config.DebounceDelay != 1*time.Second {
		t.Errorf("DebounceDelay = %v, want 1s", config.DebounceDelay)
	}

	if config.MinAge != 2*time.Second {
		t.Errorf("MinAge = %v, want 2s", config.MinAge)
	}

	if !config.DeleteSource {
		t.Error("DeleteSource = false, want true")
	}

	if config.SyncOnStart {
		t.Error("SyncOnStart = true, want false")
	}

	if len(config.ExcludePatterns) != 1 || config.ExcludePatterns[0] != "*.log" {
		t.Errorf("ExcludePatterns = %v, want [*.log]", config.ExcludePatterns)
	}

	if config.CronSchedule != "0 * * * *" {
		t.Errorf("CronSchedule = %s, want 0 * * * *", config.CronSchedule)
	}
}

func TestWatchStatus_Defaults(t *testing.T) {
	status := WatchStatus{
		Source:      "/local/path",
		Destination: "s3://bucket/prefix",
		Active:      true,
		StartedAt:   time.Now(),
	}

	if !status.Active {
		t.Error("Active = false, want true")
	}

	if status.StartedAt.IsZero() {
		t.Error("StartedAt is zero, want current time")
	}

	if status.FilesSynced != 0 {
		t.Errorf("FilesSynced = %d, want 0", status.FilesSynced)
	}

	if status.BytesSynced != 0 {
		t.Errorf("BytesSynced = %d, want 0", status.BytesSynced)
	}

	if status.ErrorCount != 0 {
		t.Errorf("ErrorCount = %d, want 0", status.ErrorCount)
	}

	if status.LastError != "" {
		t.Errorf("LastError = %s, want empty", status.LastError)
	}
}
