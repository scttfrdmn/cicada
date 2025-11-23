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

package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/scttfrdmn/cicada/internal/sync"
)

// NewSyncCmd creates the sync command.
func NewSyncCmd() *cobra.Command {
	var (
		source      string
		destination string
		dryRun      bool
		delete      bool
	)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync files between local and S3",
		Long: `Sync files between local filesystem and S3 storage.

Examples:
  # Sync local directory to S3
  cicada sync /data/lab s3://my-bucket/lab-data

  # Sync from S3 to local
  cicada sync s3://my-bucket/lab-data /data/lab

  # Dry run to see what would be synced
  cicada sync --dry-run /data/lab s3://my-bucket/lab-data

  # Sync and delete files not in source
  cicada sync --delete /data/lab s3://my-bucket/lab-data`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			source = args[0]
			destination = args[1]

			ctx := context.Background()

			if verbose {
				fmt.Printf("Syncing from %s to %s\n", source, destination)
				if dryRun {
					fmt.Println("DRY RUN: No changes will be made")
				}
			}

			// Create backends
			srcBackend, srcPath, err := createBackend(ctx, source)
			if err != nil {
				return fmt.Errorf("create source backend: %w", err)
			}
			defer func() { _ = srcBackend.Close() }()

			dstBackend, dstPath, err := createBackend(ctx, destination)
			if err != nil {
				return fmt.Errorf("create destination backend: %w", err)
			}
			defer func() { _ = dstBackend.Close() }()

			// Create sync engine
			engine := sync.NewEngine(srcBackend, dstBackend, sync.SyncOptions{
				DryRun:      dryRun,
				Delete:      delete,
				Concurrency: 4,
				ProgressFunc: func(update sync.ProgressUpdate) {
					if verbose {
						if update.Error != nil {
							fmt.Printf("❌ %s %s: %v\n", update.Operation, update.Path, update.Error)
						} else if update.BytesTotal > 0 {
							fmt.Printf("✓ %s %s (%d bytes)\n", update.Operation, update.Path, update.BytesTotal)
						} else {
							fmt.Printf("✓ %s %s\n", update.Operation, update.Path)
						}
					}
				},
			})

			// Perform sync
			if err := engine.Sync(ctx, srcPath, dstPath); err != nil {
				return fmt.Errorf("sync failed: %w", err)
			}

			if verbose {
				fmt.Println("✓ Sync complete")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be synced without making changes")
	cmd.Flags().BoolVar(&delete, "delete", false, "delete files in destination not present in source")

	return cmd
}

// createBackend creates the appropriate backend based on the path.
func createBackend(ctx context.Context, path string) (sync.Backend, string, error) {
	if strings.HasPrefix(path, "s3://") {
		bucket, key, err := sync.ParseS3URI(path)
		if err != nil {
			return nil, "", err
		}

		backend, err := sync.NewS3Backend(ctx, bucket)
		if err != nil {
			return nil, "", fmt.Errorf("create S3 backend: %w", err)
		}

		return backend, key, nil
	}

	// Local filesystem
	backend, err := sync.NewLocalBackend(path)
	if err != nil {
		return nil, "", fmt.Errorf("create local backend: %w", err)
	}

	return backend, "", nil
}
