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

//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cicadasync "github.com/scttfrdmn/cicada/internal/sync"
)

const (
	testProfile = "aws"
	testRegion  = "us-west-2"
	testBucket  = "cicada-integration-test"
)

// TestS3Backend_Integration tests S3 backend operations against real AWS.
func TestS3Backend_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set AWS environment variables for profile and region
	originalProfile := os.Getenv("AWS_PROFILE")
	originalRegion := os.Getenv("AWS_REGION")
	if err := os.Setenv("AWS_PROFILE", testProfile); err != nil {
		t.Fatalf("Failed to set AWS_PROFILE: %v", err)
	}
	if err := os.Setenv("AWS_REGION", testRegion); err != nil {
		t.Fatalf("Failed to set AWS_REGION: %v", err)
	}
	defer func() {
		_ = os.Setenv("AWS_PROFILE", originalProfile)
		_ = os.Setenv("AWS_REGION", originalRegion)
	}()

	ctx := context.Background()

	// Load AWS config with specific profile and region
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(testProfile),
		config.WithRegion(testRegion),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Verify bucket exists
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: stringPtr(testBucket),
	})
	if err != nil {
		t.Fatalf("Test bucket %s does not exist or is not accessible: %v", testBucket, err)
	}

	// Create test prefix with timestamp to avoid conflicts
	testPrefix := fmt.Sprintf("test-run-%d/", time.Now().Unix())

	// Clean up at end of test
	defer func() {
		if err := cleanupS3Prefix(ctx, client, testBucket, testPrefix); err != nil {
			t.Logf("Warning: failed to cleanup S3 prefix: %v", err)
		}
	}()

	// Create S3 backend
	backend, err := cicadasync.NewS3Backend(ctx, testBucket)
	if err != nil {
		t.Fatalf("NewS3Backend() error: %v", err)
	}

	t.Run("Write and Read", func(t *testing.T) {
		testKey := testPrefix + "test-file.txt"
		testContent := []byte("Hello, Cicada Integration Test!")

		// Write file
		if err := backend.Write(ctx, testKey, bytes.NewReader(testContent), int64(len(testContent))); err != nil {
			t.Fatalf("Write() error: %v", err)
		}

		// Read file back
		reader, err := backend.Read(ctx, testKey)
		if err != nil {
			t.Fatalf("Read() error: %v", err)
		}
		defer func() { _ = reader.Close() }()

		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("ReadAll() error: %v", err)
		}

		if string(content) != string(testContent) {
			t.Errorf("Read() content = %s, want %s", string(content), string(testContent))
		}
	})

	t.Run("Stat", func(t *testing.T) {
		testKey := testPrefix + "stat-test.txt"
		testContent := []byte("test content for stat")

		// Write file
		if err := backend.Write(ctx, testKey, bytes.NewReader(testContent), int64(len(testContent))); err != nil {
			t.Fatalf("Write() error: %v", err)
		}

		// Stat file
		info, err := backend.Stat(ctx, testKey)
		if err != nil {
			t.Fatalf("Stat() error: %v", err)
		}

		if info.Size != int64(len(testContent)) {
			t.Errorf("Stat() size = %d, want %d", info.Size, len(testContent))
		}

		if info.ModTime.IsZero() {
			t.Error("Stat() ModTime is zero")
		}

		if info.ETag == "" {
			t.Error("Stat() ETag is empty")
		}
	})

	t.Run("List", func(t *testing.T) {
		// Write multiple test files
		testFiles := map[string][]byte{
			testPrefix + "list-test/file1.txt": []byte("content 1"),
			testPrefix + "list-test/file2.txt": []byte("content 2"),
			testPrefix + "list-test/file3.txt": []byte("content 3"),
		}

		for key, content := range testFiles {
			if err := backend.Write(ctx, key, bytes.NewReader(content), int64(len(content))); err != nil {
				t.Fatalf("Write() error for %s: %v", key, err)
			}
		}

		// List files
		files, err := backend.List(ctx, testPrefix)
		if err != nil {
			t.Fatalf("List() error: %v", err)
		}

		// Check that our test files are in the list
		foundCount := 0
		for _, file := range files {
			if _, exists := testFiles[file.Path]; exists {
				foundCount++
			}
		}

		if foundCount != len(testFiles) {
			t.Errorf("List() found %d test files, want %d", foundCount, len(testFiles))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		testKey := testPrefix + "delete-test.txt"
		testContent := []byte("to be deleted")

		// Write file
		if err := backend.Write(ctx, testKey, bytes.NewReader(testContent), int64(len(testContent))); err != nil {
			t.Fatalf("Write() error: %v", err)
		}

		// Verify it exists
		if _, err := backend.Stat(ctx, testKey); err != nil {
			t.Fatalf("Stat() before delete error: %v", err)
		}

		// Delete file
		if err := backend.Delete(ctx, testKey); err != nil {
			t.Fatalf("Delete() error: %v", err)
		}

		// Verify it's gone
		_, err := backend.Stat(ctx, testKey)
		if err == nil {
			t.Error("Stat() after delete should return error, got nil")
		}
	})
}

// TestLocalToS3Sync_Integration tests syncing from local filesystem to S3.
func TestLocalToS3Sync_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set AWS environment variables
	originalProfile := os.Getenv("AWS_PROFILE")
	originalRegion := os.Getenv("AWS_REGION")
	if err := os.Setenv("AWS_PROFILE", testProfile); err != nil {
		t.Fatalf("Failed to set AWS_PROFILE: %v", err)
	}
	if err := os.Setenv("AWS_REGION", testRegion); err != nil {
		t.Fatalf("Failed to set AWS_REGION: %v", err)
	}
	defer func() {
		_ = os.Setenv("AWS_PROFILE", originalProfile)
		_ = os.Setenv("AWS_REGION", originalRegion)
	}()

	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(testProfile),
		config.WithRegion(testRegion),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create test prefix
	testPrefix := fmt.Sprintf("sync-test-%d/", time.Now().Unix())

	// Clean up at end
	defer func() {
		if err := cleanupS3Prefix(ctx, client, testBucket, testPrefix); err != nil {
			t.Logf("Warning: failed to cleanup S3 prefix: %v", err)
		}
	}()

	// Create temporary local directory with test files
	tmpDir := t.TempDir()

	testFiles := map[string]string{
		"file1.txt":           "Hello from file 1",
		"file2.txt":           "Hello from file 2",
		"subdir/file3.txt":    "Hello from subdirectory",
		"subdir/nested/file4.txt": "Hello from nested directory",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("MkdirAll() error: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}
	}

	// Create backends
	srcBackend, err := cicadasync.NewLocalBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalBackend() error: %v", err)
	}

	dstBackend, err := cicadasync.NewS3Backend(ctx, testBucket)
	if err != nil {
		t.Fatalf("NewS3Backend() error: %v", err)
	}

	// Create sync engine
	engine := cicadasync.NewEngine(srcBackend, dstBackend, cicadasync.SyncOptions{
		Concurrency: 4,
		DryRun:      false,
		Delete:      false,
	})

	// Perform sync
	if err := engine.Sync(ctx, "", testPrefix); err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	// Verify files exist in S3
	files, err := dstBackend.List(ctx, testPrefix)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(files) != len(testFiles) {
		t.Errorf("S3 has %d files, want %d", len(files), len(testFiles))
	}

	// Verify content of each file
	for path, expectedContent := range testFiles {
		s3Path := testPrefix + path
		reader, err := dstBackend.Read(ctx, s3Path)
		if err != nil {
			t.Errorf("Read() error for %s: %v", path, err)
			continue
		}

		content, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			t.Errorf("ReadAll() error for %s: %v", path, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("Content of %s = %s, want %s", path, string(content), expectedContent)
		}
	}
}

// TestS3ToLocalSync_Integration tests syncing from S3 to local filesystem.
func TestS3ToLocalSync_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set AWS environment variables
	originalProfile := os.Getenv("AWS_PROFILE")
	originalRegion := os.Getenv("AWS_REGION")
	if err := os.Setenv("AWS_PROFILE", testProfile); err != nil {
		t.Fatalf("Failed to set AWS_PROFILE: %v", err)
	}
	if err := os.Setenv("AWS_REGION", testRegion); err != nil {
		t.Fatalf("Failed to set AWS_REGION: %v", err)
	}
	defer func() {
		_ = os.Setenv("AWS_PROFILE", originalProfile)
		_ = os.Setenv("AWS_REGION", originalRegion)
	}()

	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(testProfile),
		config.WithRegion(testRegion),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create test prefix
	testPrefix := fmt.Sprintf("download-test-%d/", time.Now().Unix())

	// Clean up at end
	defer func() {
		if err := cleanupS3Prefix(ctx, client, testBucket, testPrefix); err != nil {
			t.Logf("Warning: failed to cleanup S3 prefix: %v", err)
		}
	}()

	// Create S3 backend and populate with test files
	srcBackend, err := cicadasync.NewS3Backend(ctx, testBucket)
	if err != nil {
		t.Fatalf("NewS3Backend() error: %v", err)
	}

	testFiles := map[string]string{
		"download1.txt":        "Content from S3 file 1",
		"download2.txt":        "Content from S3 file 2",
		"subdir/download3.txt": "Content from S3 subdirectory",
	}

	for path, content := range testFiles {
		s3Path := testPrefix + path
		if err := srcBackend.Write(ctx, s3Path, bytes.NewReader([]byte(content)), int64(len(content))); err != nil {
			t.Fatalf("Write() error for %s: %v", path, err)
		}
	}

	// Create temporary local directory
	tmpDir := t.TempDir()
	dstBackend, err := cicadasync.NewLocalBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalBackend() error: %v", err)
	}

	// Create sync engine
	engine := cicadasync.NewEngine(srcBackend, dstBackend, cicadasync.SyncOptions{
		Concurrency: 4,
		DryRun:      false,
		Delete:      false,
	})

	// Perform sync
	if err := engine.Sync(ctx, testPrefix, ""); err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	// Verify files exist locally
	for path, expectedContent := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("ReadFile() error for %s: %v", path, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("Content of %s = %s, want %s", path, string(content), expectedContent)
		}
	}
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func cleanupS3Prefix(ctx context.Context, client *s3.Client, bucket, prefix string) error {
	// List all objects with the prefix
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: stringPtr(bucket),
		Prefix: stringPtr(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("list objects: %w", err)
		}

		// Delete objects in batches
		for _, obj := range page.Contents {
			_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: stringPtr(bucket),
				Key:    obj.Key,
			})
			if err != nil {
				return fmt.Errorf("delete object %s: %w", *obj.Key, err)
			}
		}
	}

	return nil
}
