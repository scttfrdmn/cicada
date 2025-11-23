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
	"time"
)

// FileInfo represents metadata about a file.
type FileInfo struct {
	Path         string
	Size         int64
	ModTime      time.Time
	ETag         string // Checksum/hash
	IsDir        bool
	StorageClass string // For S3: STANDARD, GLACIER, etc.
}

// Backend represents a storage backend (local filesystem or S3).
type Backend interface {
	// List returns all files with the given prefix.
	List(ctx context.Context, prefix string) ([]FileInfo, error)

	// Read opens a file for reading.
	Read(ctx context.Context, path string) (io.ReadCloser, error)

	// Write writes a file.
	Write(ctx context.Context, path string, r io.Reader, size int64) error

	// Delete deletes a file.
	Delete(ctx context.Context, path string) error

	// Stat gets file metadata.
	Stat(ctx context.Context, path string) (*FileInfo, error)

	// Close closes the backend and releases resources.
	Close() error
}
