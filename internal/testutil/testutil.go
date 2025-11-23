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

// Package testutil provides testing utilities for Cicada.
package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory for testing.
// The directory is automatically cleaned up after the test.
func TempDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "cicada-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	return dir
}

// WriteFile writes content to a file in the given directory.
func WriteFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create parent dir: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	return path
}
