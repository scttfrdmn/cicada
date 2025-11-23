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

// NewConfigCmd creates the config command.
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Cicada configuration",
		Long:  `View and modify Cicada configuration settings.`,
	}

	cmd.AddCommand(NewConfigShowCmd())
	cmd.AddCommand(NewConfigInitCmd())

	return cmd
}

// NewConfigShowCmd creates the config show subcommand.
func NewConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement config show logic
			fmt.Println("Configuration:")
			fmt.Println("  Config file: ~/.cicada/config.yaml")
			fmt.Println("  (not yet configured)")
			return nil
		},
	}
}

// NewConfigInitCmd creates the config init subcommand.
func NewConfigInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize Cicada configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement config init logic
			return fmt.Errorf("config init not yet implemented")
		},
	}
}
