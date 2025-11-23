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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalBackend implements Backend for local filesystem.
type LocalBackend struct {
	root string
}

// NewLocalBackend creates a new local filesystem backend.
func NewLocalBackend(root string) (*LocalBackend, error) {
	// Ensure root exists
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("create root directory: %w", err)
	}

	return &LocalBackend{root: root}, nil
}

// List returns all files under the given prefix.
func (b *LocalBackend) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	var files []FileInfo

	fullPath := filepath.Join(b.root, prefix)

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(b.root, path)
		if err != nil {
			return err
		}

		// Calculate ETag (MD5 hash) for files
		var etag string
		if !info.IsDir() && info.Size() > 0 {
			etag, _ = b.calculateMD5(path)
		}

		files = append(files, FileInfo{
			Path:    relPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			ETag:    etag,
			IsDir:   info.IsDir(),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk directory: %w", err)
	}

	return files, nil
}

// Read opens a file for reading.
func (b *LocalBackend) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(b.root, path)
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return f, nil
}

// Write writes a file.
func (b *LocalBackend) Write(ctx context.Context, path string, r io.Reader, size int64) error {
	fullPath := filepath.Join(b.root, path)

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}

	// Create file
	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() { _ = f.Close() }()

	// Copy data
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	return nil
}

// Delete deletes a file.
func (b *LocalBackend) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(b.root, path)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("remove file: %w", err)
	}
	return nil
}

// Stat gets file metadata.
func (b *LocalBackend) Stat(ctx context.Context, path string) (*FileInfo, error) {
	fullPath := filepath.Join(b.root, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	var etag string
	if !info.IsDir() && info.Size() > 0 {
		etag, _ = b.calculateMD5(fullPath)
	}

	return &FileInfo{
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		ETag:    etag,
		IsDir:   info.IsDir(),
	}, nil
}

// Close closes the backend.
func (b *LocalBackend) Close() error {
	return nil // Nothing to close for local backend
}

// calculateMD5 computes MD5 hash of a file.
func (b *LocalBackend) calculateMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
