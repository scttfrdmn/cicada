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
	"fmt"

	"github.com/spf13/cobra"
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

			if verbose {
				fmt.Printf("Syncing from %s to %s\n", source, destination)
				if dryRun {
					fmt.Println("DRY RUN: No changes will be made")
				}
			}

			// TODO: Implement sync logic
			return fmt.Errorf("sync not yet implemented")
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be synced without making changes")
	cmd.Flags().BoolVar(&delete, "delete", false, "delete files in destination not present in source")

	return cmd
}
