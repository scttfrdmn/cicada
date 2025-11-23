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
	"sync"
	"testing"
	"time"
)

func TestNeedsSync(t *testing.T) {
	baseTime := time.Now()

	tests := []struct {
		name     string
		src      FileInfo
		dst      FileInfo
		expected bool
	}{
		{
			name: "same etag - no sync",
			src: FileInfo{
				Path:  "test.txt",
				Size:  100,
				ETag:  "abc123",
			},
			dst: FileInfo{
				Path:  "test.txt",
				Size:  100,
				ETag:  "abc123",
			},
			expected: false,
		},
		{
			name: "different etag - needs sync",
			src: FileInfo{
				Path:  "test.txt",
				Size:  100,
				ETag:  "abc123",
			},
			dst: FileInfo{
				Path:  "test.txt",
				Size:  100,
				ETag:  "xyz789",
			},
			expected: true,
		},
		{
			name: "no etag, different size - needs sync",
			src: FileInfo{
				Path:    "test.txt",
				Size:    100,
				ModTime: baseTime,
			},
			dst: FileInfo{
				Path:    "test.txt",
				Size:    50,
				ModTime: baseTime,
			},
			expected: true,
		},
		{
			name: "no etag, same size, src newer - needs sync",
			src: FileInfo{
				Path:    "test.txt",
				Size:    100,
				ModTime: baseTime.Add(1 * time.Hour),
			},
			dst: FileInfo{
				Path:    "test.txt",
				Size:    100,
				ModTime: baseTime,
			},
			expected: true,
		},
		{
			name: "no etag, same size, dst newer - no sync",
			src: FileInfo{
				Path:    "test.txt",
				Size:    100,
				ModTime: baseTime,
			},
			dst: FileInfo{
				Path:    "test.txt",
				Size:    100,
				ModTime: baseTime.Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "no etag, same size and time - no sync",
			src: FileInfo{
				Path:    "test.txt",
				Size:    100,
				ModTime: baseTime,
			},
			dst: FileInfo{
				Path:    "test.txt",
				Size:    100,
				ModTime: baseTime,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := needsSync(tt.src, tt.dst)
			if result != tt.expected {
				t.Errorf("needsSync() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewEngine(t *testing.T) {
	tests := []struct {
		name               string
		options            SyncOptions
		expectedConcurrency int
	}{
		{
			name: "default concurrency",
			options: SyncOptions{
				Concurrency: 0,
			},
			expectedConcurrency: 4,
		},
		{
			name: "custom concurrency",
			options: SyncOptions{
				Concurrency: 10,
			},
			expectedConcurrency: 10,
		},
		{
			name: "negative concurrency defaults to 4",
			options: SyncOptions{
				Concurrency: -1,
			},
			expectedConcurrency: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine(nil, nil, tt.options)
			if engine.options.Concurrency != tt.expectedConcurrency {
				t.Errorf("NewEngine() concurrency = %v, want %v",
					engine.options.Concurrency, tt.expectedConcurrency)
			}
		})
	}
}

func TestEngine_Sync(t *testing.T) {
	baseTime := time.Now()

	t.Run("sync new files", func(t *testing.T) {
		src := newMockBackend()
		dst := newMockBackend()

		// Add files to source
		src.addFile("file1.txt", "content1", "etag1", baseTime)
		src.addFile("file2.txt", "content2", "etag2", baseTime)

		engine := NewEngine(src, dst, SyncOptions{
			Concurrency: 2,
		})

		ctx := context.Background()
		err := engine.Sync(ctx, "", "")
		if err != nil {
			t.Fatalf("Sync() error = %v", err)
		}

		// Verify files were synced
		dstFiles, _ := dst.List(ctx, "")
		if len(dstFiles) != 2 {
			t.Errorf("Sync() synced %d files, want 2", len(dstFiles))
		}
	})

	t.Run("skip unchanged files", func(t *testing.T) {
		src := newMockBackend()
		dst := newMockBackend()

		// Add same file to both with same etag
		src.addFile("file.txt", "content", "same-etag", baseTime)
		dst.addFile("file.txt", "content", "same-etag", baseTime)

		syncCount := 0
		engine := NewEngine(src, dst, SyncOptions{
			ProgressFunc: func(update ProgressUpdate) {
				if update.Operation == "upload" {
					syncCount++
				}
			},
		})

		ctx := context.Background()
		err := engine.Sync(ctx, "", "")
		if err != nil {
			t.Fatalf("Sync() error = %v", err)
		}

		if syncCount != 0 {
			t.Errorf("Sync() synced %d files, want 0 (should skip unchanged)", syncCount)
		}
	})

	t.Run("sync modified files", func(t *testing.T) {
		src := newMockBackend()
		dst := newMockBackend()

		// Add file with different etag
		src.addFile("file.txt", "new content", "etag-new", baseTime.Add(1*time.Hour))
		dst.addFile("file.txt", "old content", "etag-old", baseTime)

		var syncCount int
		var mu sync.Mutex
		engine := NewEngine(src, dst, SyncOptions{
			ProgressFunc: func(update ProgressUpdate) {
				if update.Operation == "upload" && update.BytesDone == update.BytesTotal && update.BytesTotal > 0 {
					mu.Lock()
					syncCount++
					mu.Unlock()
				}
			},
		})

		ctx := context.Background()
		err := engine.Sync(ctx, "", "")
		if err != nil {
			t.Fatalf("Sync() error = %v", err)
		}

		if syncCount != 1 {
			t.Errorf("Sync() synced %d files, want 1", syncCount)
		}
	})

	t.Run("delete extra files", func(t *testing.T) {
		src := newMockBackend()
		dst := newMockBackend()

		// Add file only to destination
		src.addFile("file1.txt", "content1", "etag1", baseTime)
		dst.addFile("file1.txt", "content1", "etag1", baseTime)
		dst.addFile("file2.txt", "content2", "etag2", baseTime)

		deleteCount := 0
		engine := NewEngine(src, dst, SyncOptions{
			Delete: true,
			ProgressFunc: func(update ProgressUpdate) {
				if update.Operation == "delete" {
					deleteCount++
				}
			},
		})

		ctx := context.Background()
		err := engine.Sync(ctx, "", "")
		if err != nil {
			t.Fatalf("Sync() error = %v", err)
		}

		if deleteCount != 1 {
			t.Errorf("Sync() deleted %d files, want 1", deleteCount)
		}

		// Verify file was deleted
		dstFiles, _ := dst.List(ctx, "")
		if len(dstFiles) != 1 {
			t.Errorf("After delete, destination has %d files, want 1", len(dstFiles))
		}
	})

	t.Run("dry run mode", func(t *testing.T) {
		src := newMockBackend()
		dst := newMockBackend()

		src.addFile("file.txt", "content", "etag1", baseTime)

		engine := NewEngine(src, dst, SyncOptions{
			DryRun: true,
		})

		ctx := context.Background()
		err := engine.Sync(ctx, "", "")
		if err != nil {
			t.Fatalf("Sync() error = %v", err)
		}

		// Verify nothing was actually synced
		dstFiles, _ := dst.List(ctx, "")
		if len(dstFiles) != 0 {
			t.Errorf("DryRun synced %d files, want 0", len(dstFiles))
		}
	})

	t.Run("progress reporting", func(t *testing.T) {
		src := newMockBackend()
		dst := newMockBackend()

		src.addFile("file1.txt", "content1", "etag1", baseTime)
		src.addFile("file2.txt", "content2", "etag2", baseTime)

		var updates []ProgressUpdate
		var mu sync.Mutex
		engine := NewEngine(src, dst, SyncOptions{
			ProgressFunc: func(update ProgressUpdate) {
				mu.Lock()
				updates = append(updates, update)
				mu.Unlock()
			},
		})

		ctx := context.Background()
		err := engine.Sync(ctx, "", "")
		if err != nil {
			t.Fatalf("Sync() error = %v", err)
		}

		mu.Lock()
		updateCount := len(updates)
		mu.Unlock()

		if updateCount == 0 {
			t.Error("Sync() did not report any progress")
		}
	})
}
