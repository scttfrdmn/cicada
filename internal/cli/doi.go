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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/scttfrdmn/cicada/internal/doi"
	"github.com/scttfrdmn/cicada/internal/metadata"
)

// NewDOICmd creates the DOI command
func NewDOICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doi",
		Short: "Manage Digital Object Identifiers (DOIs)",
		Long: `Manage DOI assignment for scientific datasets.

Prepare metadata, validate readiness, and mint DOIs through DataCite or Zenodo.
Ensures metadata meets DataCite Schema v4.5 requirements before minting.`,
	}

	// Add subcommands
	cmd.AddCommand(newDOIPrepareCmd())
	cmd.AddCommand(newDOIValidateCmd())
	cmd.AddCommand(newDOIMintCmd())
	cmd.AddCommand(newDOIConfigCmd())
	cmd.AddCommand(newDOIProviderCmd())

	return cmd
}

// newDOIPrepareCmd creates the doi prepare subcommand
func newDOIPrepareCmd() *cobra.Command {
	var (
		outputFormat   string
		outputFile     string
		presetID       string
		publisher      string
		license        string
		enrichmentFile string
	)

	cmd := &cobra.Command{
		Use:   "prepare <file>",
		Short: "Prepare metadata for DOI minting",
		Long: `Extract and map metadata to DataCite schema, validate readiness.

Analyzes file metadata, maps to DataCite required fields, and provides
quality assessment and recommendations for improvement.

Examples:
  # Prepare metadata from a file
  cicada doi prepare data/experiment.czi

  # Prepare with custom publisher and license
  cicada doi prepare data/sample.fastq --publisher "My Lab" --license CC-BY-4.0

  # Prepare with instrument preset
  cicada doi prepare data/image.czi --preset zeiss-lsm-880

  # Prepare with enrichment from file
  cicada doi prepare data/sample.fastq --enrich metadata.yaml

  # Save prepared metadata
  cicada doi prepare data/image.czi --output prepared.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", path)
			}

			// Extract metadata
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			extractedMeta, err := registry.Extract(path)
			if err != nil {
				return fmt.Errorf("failed to extract metadata: %w", err)
			}

			// Load enrichment if provided
			var enrichment map[string]interface{}
			if enrichmentFile != "" {
				enrichData, err := os.ReadFile(enrichmentFile)
				if err != nil {
					return fmt.Errorf("failed to read enrichment file: %w", err)
				}

				if err := yaml.Unmarshal(enrichData, &enrichment); err != nil {
					if err := json.Unmarshal(enrichData, &enrichment); err != nil {
						return fmt.Errorf("failed to parse enrichment file: %w", err)
					}
				}
			}

			// Create workflow
			config := &doi.WorkflowConfig{
				Publisher:          publisher,
				License:            license,
				MinQualityScore:    60.0,
				RequireRealAuthors: true,
				RequireDescription: true,
			}

			// For now, use a disabled provider registry
			providerRegistry := doi.NewProviderRegistry()
			providerRegistry.Register(doi.NewDisabledProvider())
			_ = providerRegistry.SetActive("disabled")

			workflow := doi.NewDOIWorkflow(config, providerRegistry)

			// Prepare
			prepReq := &doi.PrepareRequest{
				FilePath:   path,
				Metadata:   extractedMeta,
				Enrichment: enrichment,
				PresetID:   presetID,
			}

			result, err := workflow.Prepare(prepReq)
			if err != nil {
				return fmt.Errorf("failed to prepare metadata: %w", err)
			}

			// Display results
			if outputFormat == "json" || outputFile != "" {
				output := map[string]interface{}{
					"dataset":    result.Dataset,
					"validation": result.Validation,
					"warnings":   result.Warnings,
				}

				var data []byte
				if outputFormat == "yaml" {
					data, err = yaml.Marshal(output)
				} else {
					data, err = json.MarshalIndent(output, "", "  ")
				}
				if err != nil {
					return err
				}

				if outputFile != "" {
					return os.WriteFile(outputFile, data, 0644)
				}
				fmt.Println(string(data))
			} else {
				// Table format
				fmt.Printf("DOI Preparation Results\n")
				fmt.Printf("=======================\n\n")

				fmt.Printf("File: %s\n\n", filepath.Base(path))

				// Dataset info
				fmt.Printf("Dataset Information:\n")
				fmt.Printf("  Title: %s\n", result.Dataset.Title)
				fmt.Printf("  Authors: %d\n", len(result.Dataset.Authors))
				for i, author := range result.Dataset.Authors {
					fmt.Printf("    %d. %s", i+1, author.Name)
					if author.ORCID != "" {
						fmt.Printf(" (ORCID: %s)", author.ORCID)
					}
					fmt.Println()
				}
				fmt.Printf("  Publisher: %s\n", result.Dataset.Publisher)
				fmt.Printf("  Resource Type: %s\n", result.Dataset.ResourceType)
				fmt.Printf("  Keywords: %d\n", len(result.Dataset.Keywords))
				fmt.Println()

				// Validation
				fmt.Printf("Validation:\n")
				if result.Validation.IsReady {
					fmt.Printf("  ✓ Ready for DOI minting\n")
				} else {
					fmt.Printf("  ✗ Not ready for DOI minting\n")
				}
				fmt.Printf("  Quality Score: %.1f/100 (%s)\n",
					result.Validation.Score, doi.GetQualityLevel(result.Validation.Score))
				fmt.Printf("  Errors: %d\n", len(result.Validation.Errors))
				fmt.Printf("  Warnings: %d\n", len(result.Validation.Warnings))
				fmt.Println()

				// Errors
				if len(result.Validation.Errors) > 0 {
					fmt.Printf("Errors:\n")
					for _, err := range result.Validation.Errors {
						fmt.Printf("  • %s\n", err)
					}
					fmt.Println()
				}

				// Warnings
				if len(result.Validation.Warnings) > 0 {
					fmt.Printf("Warnings:\n")
					for _, warning := range result.Validation.Warnings {
						fmt.Printf("  • %s\n", warning)
					}
					fmt.Println()
				}

				// Recommendations
				recommendations := workflow.GetRecommendations(result.Validation)
				if len(recommendations) > 0 {
					fmt.Printf("Recommendations:\n")
					for _, rec := range recommendations {
						fmt.Printf("  %s\n", rec)
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")
	cmd.Flags().StringVar(&presetID, "preset", "", "Instrument preset ID")
	cmd.Flags().StringVar(&publisher, "publisher", "", "Publisher name")
	cmd.Flags().StringVar(&license, "license", "CC-BY-4.0", "License (default: CC-BY-4.0)")
	cmd.Flags().StringVar(&enrichmentFile, "enrich", "", "Enrichment metadata file (JSON or YAML)")

	return cmd
}

// newDOIValidateCmd creates the doi validate subcommand
func newDOIValidateCmd() *cobra.Command {
	var (
		outputFormat string
		minScore     float64
	)

	cmd := &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate metadata for DOI readiness",
		Long: `Validate that file metadata meets DOI requirements.

Checks DataCite Schema v4.5 required fields, calculates quality score,
and provides recommendations for improvement.

Examples:
  # Validate a file
  cicada doi validate data/experiment.czi

  # Validate with minimum quality score
  cicada doi validate data/sample.fastq --min-score 80

  # Validate and output JSON
  cicada doi validate data/image.czi --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", path)
			}

			// Extract metadata
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			extractedMeta, err := registry.Extract(path)
			if err != nil {
				return fmt.Errorf("failed to extract metadata: %w", err)
			}

			// Create workflow with validation settings
			config := &doi.WorkflowConfig{
				MinQualityScore:    minScore,
				RequireRealAuthors: true,
				RequireDescription: true,
			}

			providerRegistry := doi.NewProviderRegistry()
			providerRegistry.Register(doi.NewDisabledProvider())
			_ = providerRegistry.SetActive("disabled")

			workflow := doi.NewDOIWorkflow(config, providerRegistry)

			// Validate
			result, err := workflow.ValidateMetadata(extractedMeta, path)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// Display results
			if outputFormat == "json" {
				data, err := json.MarshalIndent(result.Validation, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else if outputFormat == "yaml" {
				data, err := yaml.Marshal(result.Validation)
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				// Table format
				fmt.Printf("DOI Validation Results\n")
				fmt.Printf("======================\n\n")

				fmt.Printf("File: %s\n\n", filepath.Base(path))

				if result.Validation.IsReady {
					fmt.Printf("✓ READY for DOI minting\n\n")
				} else {
					fmt.Printf("✗ NOT READY for DOI minting\n\n")
				}

				fmt.Printf("Quality Score: %.1f/100 (%s)\n\n",
					result.Validation.Score, doi.GetQualityLevel(result.Validation.Score))

				// Present fields
				if len(result.Validation.Present) > 0 {
					fmt.Printf("Present Fields (%d):\n", len(result.Validation.Present))
					for _, field := range result.Validation.Present {
						fmt.Printf("  ✓ %s\n", field)
					}
					fmt.Println()
				}

				// Missing fields
				if len(result.Validation.Missing) > 0 {
					fmt.Printf("Missing Fields (%d):\n", len(result.Validation.Missing))
					for _, field := range result.Validation.Missing {
						fmt.Printf("  ✗ %s\n", field)
					}
					fmt.Println()
				}

				// Errors
				if len(result.Validation.Errors) > 0 {
					fmt.Printf("Errors (%d):\n", len(result.Validation.Errors))
					for _, err := range result.Validation.Errors {
						fmt.Printf("  • %s\n", err)
					}
					fmt.Println()
				}

				// Warnings
				if len(result.Validation.Warnings) > 0 {
					fmt.Printf("Warnings (%d):\n", len(result.Validation.Warnings))
					for _, warning := range result.Validation.Warnings {
						fmt.Printf("  • %s\n", warning)
					}
					fmt.Println()
				}

				// Recommendations
				recommendations := workflow.GetRecommendations(result.Validation)
				if len(recommendations) > 0 {
					fmt.Printf("Recommendations:\n")
					for _, rec := range recommendations {
						fmt.Printf("  %s\n", rec)
					}
				}
			}

			// Exit with error code if not ready
			if !result.Validation.IsReady {
				return fmt.Errorf("validation failed: %d errors", len(result.Validation.Errors))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().Float64Var(&minScore, "min-score", 60.0, "Minimum quality score threshold")

	return cmd
}

// newDOIMintCmd creates the doi mint subcommand
func newDOIMintCmd() *cobra.Command {
	var (
		provider string
		dryRun   bool
		draft    bool
	)

	cmd := &cobra.Command{
		Use:   "mint <file>",
		Short: "Mint a DOI for a dataset",
		Long: `Mint a new DOI for a scientific dataset.

Validates metadata, interacts with configured provider (DataCite or Zenodo),
and returns the minted DOI.

NOTE: This requires provider configuration. See 'cicada doi config' for setup.

Examples:
  # Mint DOI with default provider
  cicada doi mint data/experiment.czi

  # Dry run (validate without minting)
  cicada doi mint data/sample.fastq --dry-run

  # Mint as draft (not published)
  cicada doi mint data/image.czi --draft

  # Mint with specific provider
  cicada doi mint data/sample.fastq --provider zenodo`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("DOI minting not yet implemented - provider configuration required\n" +
				"See 'cicada doi config' for setup instructions")
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "", "DOI provider (datacite, zenodo)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate without actually minting")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as draft (not published)")

	return cmd
}

// newDOIConfigCmd creates the doi config subcommand
func newDOIConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure DOI providers",
		Long: `Configure DataCite, Zenodo, or other DOI providers.

Providers require credentials and configuration. See provider-specific
documentation for setup instructions.

DataCite Requirements:
  - Repository ID
  - Password
  - DOI Prefix (e.g., 10.12345)

Zenodo Requirements:
  - Access Token
  - Optional: Community ID

Examples:
  # Show current configuration
  cicada doi config show

  # Set default publisher
  cicada doi config set publisher "My Research Lab"

  # Set default license
  cicada doi config set license CC-BY-4.0

  # Configure DataCite
  cicada doi config provider datacite --repo-id INST.LAB --prefix 10.12345

  # Configure Zenodo
  cicada doi config provider zenodo --token YOUR_TOKEN`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("DOI provider configuration")
			fmt.Println("==========================")
			fmt.Println()
			fmt.Println("Configuration file: $HOME/.cicada/doi.yaml")
			fmt.Println()
			fmt.Println("To configure providers, create a YAML file with the following structure:")
			fmt.Println()
			fmt.Println("provider: datacite  # or zenodo")
			fmt.Println("organization: \"Your Lab Name\"")
			fmt.Println("publisher: \"Your Lab Name\"")
			fmt.Println("license: CC-BY-4.0")
			fmt.Println()
			fmt.Println("datacite:")
			fmt.Println("  repository_id: INST.LAB")
			fmt.Println("  password: your-password")
			fmt.Println("  prefix: \"10.12345\"")
			fmt.Println("  test_mode: true")
			fmt.Println()
			fmt.Println("zenodo:")
			fmt.Println("  access_token: your-token")
			fmt.Println("  sandbox: true")
			fmt.Println()
			fmt.Println("See documentation for full configuration options.")
		},
	}

	return cmd
}

// newDOIProviderCmd creates the doi provider subcommand
func newDOIProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Manage DOI providers",
		Long: `List and manage configured DOI providers.

Shows available providers, their status, and cost information.

Examples:
  # List providers
  cicada doi provider list

  # Show provider details
  cicada doi provider show datacite

  # Set active provider
  cicada doi provider set datacite`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available providers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available DOI Providers")
			fmt.Println("=======================")
			fmt.Println()
			fmt.Println("• DataCite")
			fmt.Println("  Status: Not configured")
			fmt.Println("  Cost: ~$1-5 USD per DOI")
			fmt.Println()
			fmt.Println("• Zenodo")
			fmt.Println("  Status: Not configured")
			fmt.Println("  Cost: Free")
			fmt.Println()
			fmt.Println("• Disabled")
			fmt.Println("  Status: Active (default)")
			fmt.Println("  Cost: N/A")
			fmt.Println()
			fmt.Println("Configure providers with 'cicada doi config'")
		},
	})

	return cmd
}
