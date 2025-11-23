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
	"sync"
)

// SyncOptions configures sync behavior.
type SyncOptions struct {
	// DryRun shows what would be synced without making changes
	DryRun bool

	// Delete removes files in destination not present in source
	Delete bool

	// Concurrency controls parallel transfers
	Concurrency int

	// ProgressFunc is called to report progress
	ProgressFunc func(ProgressUpdate)
}

// ProgressUpdate reports sync progress.
type ProgressUpdate struct {
	Operation string // "upload", "download", "delete", "skip"
	Path      string
	BytesDone int64
	BytesTotal int64
	Error     error
}

// Engine performs sync operations between backends.
type Engine struct {
	source      Backend
	destination Backend
	options     SyncOptions
}

// NewEngine creates a new sync engine.
func NewEngine(source, destination Backend, options SyncOptions) *Engine {
	if options.Concurrency <= 0 {
		options.Concurrency = 4 // Default concurrency
	}

	return &Engine{
		source:      source,
		destination: destination,
		options:     options,
	}
}

// Sync performs the sync operation.
func (e *Engine) Sync(ctx context.Context, sourcePath, destPath string) error {
	// List files from source
	srcFiles, err := e.source.List(ctx, sourcePath)
	if err != nil {
		return fmt.Errorf("list source: %w", err)
	}

	// List files from destination
	dstFiles, err := e.destination.List(ctx, destPath)
	if err != nil {
		return fmt.Errorf("list destination: %w", err)
	}

	// Build destination file map for quick lookup
	dstMap := make(map[string]*FileInfo)
	for i := range dstFiles {
		dstMap[dstFiles[i].Path] = &dstFiles[i]
	}

	// Determine what needs to be synced
	var toSync []FileInfo
	for _, srcFile := range srcFiles {
		if srcFile.IsDir {
			continue // Skip directories
		}

		dstFile, exists := dstMap[srcFile.Path]
		if !exists || needsSync(srcFile, *dstFile) {
			toSync = append(toSync, srcFile)
		}

		// Remove from map to track what's left (for deletion)
		if exists {
			delete(dstMap, srcFile.Path)
		}
	}

	// Files remaining in dstMap are not in source
	var toDelete []string
	if e.options.Delete {
		for path := range dstMap {
			toDelete = append(toDelete, path)
		}
	}

	// Report what will be done
	if e.options.ProgressFunc != nil {
		e.options.ProgressFunc(ProgressUpdate{
			Operation: "summary",
			Path:      fmt.Sprintf("To sync: %d, To delete: %d", len(toSync), len(toDelete)),
		})
	}

	// If dry run, stop here
	if e.options.DryRun {
		return nil
	}

	// Perform sync with concurrency
	if err := e.syncFiles(ctx, toSync); err != nil {
		return err
	}

	// Perform deletes
	if err := e.deleteFiles(ctx, toDelete); err != nil {
		return err
	}

	return nil
}

func (e *Engine) syncFiles(ctx context.Context, files []FileInfo) error {
	sem := make(chan struct{}, e.options.Concurrency)
	var wg sync.WaitGroup
	errCh := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func(f FileInfo) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			if err := e.syncFile(ctx, f); err != nil {
				errCh <- fmt.Errorf("sync %s: %w", f.Path, err)
			}
		}(file)
	}

	wg.Wait()
	close(errCh)

	// Collect errors
	var firstErr error
	for err := range errCh {
		if firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func (e *Engine) syncFile(ctx context.Context, file FileInfo) error {
	// Report progress
	if e.options.ProgressFunc != nil {
		e.options.ProgressFunc(ProgressUpdate{
			Operation:  "upload",
			Path:       file.Path,
			BytesTotal: file.Size,
		})
	}

	// Read from source
	reader, err := e.source.Read(ctx, file.Path)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	defer func() { _ = reader.Close() }()

	// Write to destination
	if err := e.destination.Write(ctx, file.Path, reader, file.Size); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// Report completion
	if e.options.ProgressFunc != nil {
		e.options.ProgressFunc(ProgressUpdate{
			Operation:  "upload",
			Path:       file.Path,
			BytesDone:  file.Size,
			BytesTotal: file.Size,
		})
	}

	return nil
}

func (e *Engine) deleteFiles(ctx context.Context, paths []string) error {
	for _, path := range paths {
		if e.options.ProgressFunc != nil {
			e.options.ProgressFunc(ProgressUpdate{
				Operation: "delete",
				Path:      path,
			})
		}

		if err := e.destination.Delete(ctx, path); err != nil {
			return fmt.Errorf("delete %s: %w", path, err)
		}
	}

	return nil
}

// needsSync determines if a file needs to be synced.
func needsSync(src, dst FileInfo) bool {
	// If ETags are available and match, no sync needed
	if src.ETag != "" && dst.ETag != "" {
		return src.ETag != dst.ETag
	}

	// Otherwise compare size and mod time
	if src.Size != dst.Size {
		return true
	}

	// If source is newer, sync
	return src.ModTime.After(dst.ModTime)
}
