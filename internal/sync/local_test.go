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

package sync

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalBackend_NewLocalBackend(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test", "nested")

	backend, err := NewLocalBackend(testPath)
	if err != nil {
		t.Fatalf("NewLocalBackend() error = %v", err)
	}
	defer func() { _ = backend.Close() }()

	// Verify directory was created
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("NewLocalBackend() did not create directory")
	}
}

func TestLocalBackend_WriteAndRead(t *testing.T) {
	tmpDir := t.TempDir()
	backend, err := NewLocalBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalBackend() error = %v", err)
	}
	defer func() { _ = backend.Close() }()

	ctx := context.Background()
	testPath := "test/file.txt"
	testContent := "Hello, World!"

	// Write file
	reader := strings.NewReader(testContent)
	err = backend.Write(ctx, testPath, reader, int64(len(testContent)))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Read file
	rc, err := backend.Read(ctx, testPath)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	defer func() { _ = rc.Close() }()

	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Read() content = %q, want %q", string(content), testContent)
	}
}

func TestLocalBackend_List(t *testing.T) {
	tmpDir := t.TempDir()
	backend, err := NewLocalBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalBackend() error = %v", err)
	}
	defer func() { _ = backend.Close() }()

	ctx := context.Background()

	// Create test files
	files := []string{
		"file1.txt",
		"dir1/file2.txt",
		"dir1/file3.txt",
		"dir2/file4.txt",
	}

	for _, file := range files {
		content := "test content"
		err := backend.Write(ctx, file, strings.NewReader(content), int64(len(content)))
		if err != nil {
			t.Fatalf("Write(%s) error = %v", file, err)
		}
	}

	// List all files
	fileList, err := backend.List(ctx, "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// Count non-directory files
	fileCount := 0
	for _, f := range fileList {
		if !f.IsDir {
			fileCount++
		}
	}

	if fileCount != len(files) {
		t.Errorf("List() found %d files, want %d", fileCount, len(files))
	}
}

func TestLocalBackend_Stat(t *testing.T) {
	tmpDir := t.TempDir()
	backend, err := NewLocalBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalBackend() error = %v", err)
	}
	defer func() { _ = backend.Close() }()

	ctx := context.Background()
	testPath := "test.txt"
	testContent := "Hello, World!"

	// Write file
	err = backend.Write(ctx, testPath, strings.NewReader(testContent), int64(len(testContent)))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Stat file
	info, err := backend.Stat(ctx, testPath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	if info.Path != testPath {
		t.Errorf("Stat() path = %q, want %q", info.Path, testPath)
	}

	if info.Size != int64(len(testContent)) {
		t.Errorf("Stat() size = %d, want %d", info.Size, len(testContent))
	}

	if info.ETag == "" {
		t.Errorf("Stat() ETag is empty, want MD5 hash")
	}

	if info.IsDir {
		t.Errorf("Stat() IsDir = true, want false")
	}
}

func TestLocalBackend_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	backend, err := NewLocalBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalBackend() error = %v", err)
	}
	defer func() { _ = backend.Close() }()

	ctx := context.Background()
	testPath := "test.txt"
	testContent := "Hello, World!"

	// Write file
	err = backend.Write(ctx, testPath, strings.NewReader(testContent), int64(len(testContent)))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Delete file
	err = backend.Delete(ctx, testPath)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify file is gone
	_, err = backend.Stat(ctx, testPath)
	if err == nil {
		t.Errorf("Delete() file still exists")
	}
}

func TestLocalBackend_CalculateMD5(t *testing.T) {
	tmpDir := t.TempDir()
	backend, err := NewLocalBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalBackend() error = %v", err)
	}
	defer func() { _ = backend.Close() }()

	ctx := context.Background()
	testPath := "test.txt"
	testContent := "Hello, World!"

	// Write file
	err = backend.Write(ctx, testPath, strings.NewReader(testContent), int64(len(testContent)))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Get MD5 via Stat
	info, err := backend.Stat(ctx, testPath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	// Verify MD5 is valid hex string and not empty
	if len(info.ETag) != 32 {
		t.Errorf("ETag length = %d, want 32 (MD5 hex string)", len(info.ETag))
	}
}
