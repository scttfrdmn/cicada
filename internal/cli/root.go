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

var (
	cfgFile string
	verbose bool
)

// Execute runs the root command.
func Execute(version string) error {
	rootCmd := NewRootCmd(version)
	return rootCmd.Execute()
}

// NewRootCmd creates the root command.
func NewRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cicada",
		Short: "Cicada Data Commons - Dormant data management for academic research",
		Long: `Cicada is a lightweight, cost-effective data commons platform designed for
academic research labs with limited technical expertise and tight budgets.

Like its namesake, Cicada lies dormant most of the time, consuming minimal
resources, but emerges powerfully when needed.`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Cicada Data Commons v" + version)
			fmt.Println("Use 'cicada --help' for available commands")
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cicada/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(NewSyncCmd())
	rootCmd.AddCommand(NewWatchCmd())
	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.AddCommand(NewMetadataCmd())
	rootCmd.AddCommand(NewDOICmd())
	rootCmd.AddCommand(NewVersionCmd(version))

	return rootCmd
}
