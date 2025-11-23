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
	"time"

	"github.com/spf13/cobra"
	"github.com/scttfrdmn/cicada/internal/watch"
)

var (
	watchManager = watch.NewManager()
)

// NewWatchCmd creates the watch command.
func NewWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch directories for changes and auto-sync",
		Long: `Watch one or more directories for file changes and automatically sync to S3.

Examples:
  # Watch a single directory
  cicada watch /data/microscope s3://my-bucket/microscope-data

  # Watch with configuration file
  cicada watch --config /etc/cicada/watch-config.yaml

  # List active watches
  cicada watch list`,
	}

	cmd.AddCommand(NewWatchAddCmd())
	cmd.AddCommand(NewWatchListCmd())
	cmd.AddCommand(NewWatchRemoveCmd())

	return cmd
}

// NewWatchAddCmd creates the watch add subcommand.
func NewWatchAddCmd() *cobra.Command {
	var (
		debounce     int
		minAge       int
		deleteSource bool
		syncOnStart  bool
	)

	cmd := &cobra.Command{
		Use:   "add <source> <destination>",
		Short: "Add a directory to watch",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			source := args[0]
			destination := args[1]

			if verbose {
				fmt.Printf("Adding watch: %s -> %s\n", source, destination)
			}

			ctx := context.Background()

			// Create backends
			srcBackend, srcPath, err := createBackend(ctx, source)
			if err != nil {
				return fmt.Errorf("create source backend: %w", err)
			}

			dstBackend, dstPath, err := createBackend(ctx, destination)
			if err != nil {
				return fmt.Errorf("create destination backend: %w", err)
			}

			// Create watch config
			config := watch.DefaultConfig()
			config.Source = srcPath
			config.Destination = dstPath
			config.DebounceDelay = time.Duration(debounce) * time.Second
			config.MinAge = time.Duration(minAge) * time.Second
			config.DeleteSource = deleteSource
			config.SyncOnStart = syncOnStart

			// Generate watch ID (simple for now)
			watchID := fmt.Sprintf("%s-%d", source, time.Now().Unix())

			// Add watch
			if err := watchManager.Add(watchID, config, srcBackend, dstBackend); err != nil {
				return fmt.Errorf("add watch: %w", err)
			}

			fmt.Printf("✓ Watch started: %s\n", watchID)
			fmt.Printf("  Source: %s\n", source)
			fmt.Printf("  Destination: %s\n", destination)

			return nil
		},
	}

	cmd.Flags().IntVar(&debounce, "debounce", 5, "debounce delay in seconds")
	cmd.Flags().IntVar(&minAge, "min-age", 10, "minimum file age before sync in seconds")
	cmd.Flags().BoolVar(&deleteSource, "delete-source", false, "delete source files after sync")
	cmd.Flags().BoolVar(&syncOnStart, "sync-on-start", true, "perform initial sync when starting")

	return cmd
}

// NewWatchListCmd creates the watch list subcommand.
func NewWatchListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List active watches",
		RunE: func(cmd *cobra.Command, args []string) error {
			statuses := watchManager.List()

			if len(statuses) == 0 {
				fmt.Println("No active watches")
				return nil
			}

			fmt.Printf("Active watches: %d\n\n", len(statuses))

			for id, status := range statuses {
				fmt.Printf("Watch: %s\n", id)
				fmt.Printf("  Source: %s\n", status.Source)
				fmt.Printf("  Destination: %s\n", status.Destination)
				fmt.Printf("  Active: %v\n", status.Active)
				fmt.Printf("  Started: %s\n", status.StartedAt.Format(time.RFC3339))

				if !status.LastSync.IsZero() {
					fmt.Printf("  Last sync: %s\n", status.LastSync.Format(time.RFC3339))
					fmt.Printf("  Files synced: %d\n", status.FilesSynced)
					fmt.Printf("  Bytes synced: %d\n", status.BytesSynced)
				}

				if status.ErrorCount > 0 {
					fmt.Printf("  Errors: %d\n", status.ErrorCount)
					if status.LastError != "" {
						fmt.Printf("  Last error: %s\n", status.LastError)
					}
				}

				fmt.Println()
			}

			return nil
		},
	}
}

// NewWatchRemoveCmd creates the watch remove subcommand.
func NewWatchRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a watch by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			if verbose {
				fmt.Printf("Removing watch: %s\n", id)
			}

			if err := watchManager.Remove(id); err != nil {
				return fmt.Errorf("remove watch: %w", err)
			}

			fmt.Printf("✓ Watch removed: %s\n", id)
			return nil
		},
	}
}
