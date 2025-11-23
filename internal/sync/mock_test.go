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
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// mockBackend is a simple in-memory backend for testing.
type mockBackend struct {
	mu    sync.RWMutex
	files map[string]mockFile
}

type mockFile struct {
	content string
	modTime time.Time
	etag    string
}

func newMockBackend() *mockBackend {
	return &mockBackend{
		files: make(map[string]mockFile),
	}
}

func (m *mockBackend) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var files []FileInfo
	for path, file := range m.files {
		if prefix == "" || strings.HasPrefix(path, prefix) {
			files = append(files, FileInfo{
				Path:    path,
				Size:    int64(len(file.content)),
				ModTime: file.modTime,
				ETag:    file.etag,
				IsDir:   false,
			})
		}
	}
	return files, nil
}

func (m *mockBackend) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	file, exists := m.files[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return io.NopCloser(strings.NewReader(file.content)), nil
}

func (m *mockBackend) Write(ctx context.Context, path string, r io.Reader, size int64) error {
	content, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.files[path] = mockFile{
		content: string(content),
		modTime: time.Now(),
		etag:    fmt.Sprintf("mock-etag-%s", path),
	}
	return nil
}

func (m *mockBackend) Delete(ctx context.Context, path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.files[path]; !exists {
		return fmt.Errorf("file not found: %s", path)
	}
	delete(m.files, path)
	return nil
}

func (m *mockBackend) Stat(ctx context.Context, path string) (*FileInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	file, exists := m.files[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return &FileInfo{
		Path:    path,
		Size:    int64(len(file.content)),
		ModTime: file.modTime,
		ETag:    file.etag,
		IsDir:   false,
	}, nil
}

func (m *mockBackend) Close() error {
	return nil
}

// addFile is a helper for tests to add files to the mock backend.
func (m *mockBackend) addFile(path, content, etag string, modTime time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.files[path] = mockFile{
		content: content,
		modTime: modTime,
		etag:    etag,
	}
}
