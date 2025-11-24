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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/scttfrdmn/cicada/internal/config"
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
		providerName   string
		dryRun         bool
		draft          bool
		publisher      string
		license        string
		enrichmentFile string
		presetID       string
		sandbox        bool
	)

	cmd := &cobra.Command{
		Use:   "mint <file>",
		Short: "Mint a DOI for a dataset",
		Long: `Mint a new DOI for a scientific dataset.

Validates metadata, interacts with configured provider (DataCite or Zenodo),
and returns the minted DOI.

NOTE: This requires provider configuration. Set credentials via environment
variables, config file (~/.config/cicada/config.yaml), or .env file.

Environment Variables:
  Zenodo:
    CICADA_ZENODO_TOKEN or ZENODO_TOKEN
    CICADA_ZENODO_SANDBOX=true (optional)

  DataCite:
    CICADA_DATACITE_REPOSITORY_ID or DATACITE_REPOSITORY_ID
    CICADA_DATACITE_PASSWORD or DATACITE_PASSWORD
    CICADA_DATACITE_SANDBOX=true (optional)

Examples:
  # Mint DOI with Zenodo (free)
  export CICADA_ZENODO_TOKEN="your-token"
  cicada doi mint data/experiment.fastq --provider zenodo

  # Mint with DataCite
  export CICADA_DATACITE_REPOSITORY_ID="10.5072/FK2"
  export CICADA_DATACITE_PASSWORD="your-password"
  cicada doi mint data/sample.fastq --provider datacite

  # Dry run (validate without minting)
  cicada doi mint data/sample.fastq --provider zenodo --dry-run

  # With enrichment metadata
  cicada doi mint data/image.czi --provider zenodo --enrich metadata.yaml

  # Use sandbox for testing
  cicada doi mint data/test.fastq --provider zenodo --sandbox`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDOIMint(args[0], providerName, publisher, license, enrichmentFile, presetID, sandbox, dryRun, draft)
		},
	}

	cmd.Flags().StringVar(&providerName, "provider", "zenodo", "DOI provider (datacite, zenodo)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate without actually minting")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as draft (not published)")
	cmd.Flags().StringVar(&publisher, "publisher", "", "Publisher name")
	cmd.Flags().StringVar(&license, "license", "CC-BY-4.0", "License (default: CC-BY-4.0)")
	cmd.Flags().StringVar(&enrichmentFile, "enrich", "", "Enrichment metadata file (JSON or YAML)")
	cmd.Flags().StringVar(&presetID, "preset", "", "Instrument preset ID")
	cmd.Flags().BoolVar(&sandbox, "sandbox", false, "Use sandbox/test environment")

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

// runDOIMint implements the DOI minting logic
func runDOIMint(filePath, providerName, publisher, license, enrichmentFile, presetID string, sandbox, dryRun, draft bool) error {
	ctx := context.Background()

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	fmt.Printf("Minting DOI for: %s\n", filepath.Base(filePath))
	fmt.Printf("Provider: %s\n", providerName)
	if sandbox {
		fmt.Printf("Environment: SANDBOX (test)\n")
	} else {
		fmt.Printf("Environment: PRODUCTION\n")
	}
	if dryRun {
		fmt.Printf("Mode: DRY RUN (no actual minting)\n")
	}
	fmt.Println()

	// Step 1: Load credentials
	fmt.Println("→ Loading credentials...")
	credentials := config.NewProviderCredentials()
	
	// Load from all sources
	credentials.LoadFromEnvironment()
	
	// Try to load from config file
	configPath, err := config.ConfigPath()
	if err == nil {
		_ = credentials.LoadFromConfigFile(configPath)
	}
	
	// Try to load from .env in current directory
	workDir, err := os.Getwd()
	if err == nil {
		_ = credentials.LoadFromDotEnv(workDir)
	}

	// Step 2: Create provider
	fmt.Printf("→ Initializing %s provider...\n", providerName)
	var provider doi.Provider
	
	switch strings.ToLower(providerName) {
	case "zenodo":
		tokenCred := credentials.GetCredential("zenodo_token")
		if tokenCred.Source == config.SourceNotFound {
			return fmt.Errorf("Zenodo token not found. Set CICADA_ZENODO_TOKEN environment variable or configure in ~/.config/cicada/config.yaml")
		}
		
		// Validate token
		if err := config.ValidateZenodoToken(tokenCred.Value); err != nil {
			return fmt.Errorf("invalid Zenodo token: %w", err)
		}
		
		fmt.Printf("  Using token from %s: %s\n", tokenCred.Source, config.RedactToken(tokenCred.Value))
		
		zenodoConfig := &doi.ZenodoConfig{
			Token:   tokenCred.Value,
			Sandbox: sandbox,
		}
		
		provider, err = doi.NewZenodoProvider(zenodoConfig)
		if err != nil {
			return fmt.Errorf("create Zenodo provider: %w", err)
		}
		
	case "datacite":
		repoIDCred := credentials.GetCredential("datacite_repository_id")
		passwordCred := credentials.GetCredential("datacite_password")
		
		if repoIDCred.Source == config.SourceNotFound {
			return fmt.Errorf("DataCite repository ID not found. Set CICADA_DATACITE_REPOSITORY_ID environment variable or configure in ~/.config/cicada/config.yaml")
		}
		if passwordCred.Source == config.SourceNotFound {
			return fmt.Errorf("DataCite password not found. Set CICADA_DATACITE_PASSWORD environment variable or configure in ~/.config/cicada/config.yaml")
		}
		
		// Validate credentials
		if err := config.ValidateDataCiteRepositoryID(repoIDCred.Value); err != nil {
			return fmt.Errorf("invalid DataCite repository ID: %w", err)
		}
		if err := config.ValidateDataCitePassword(passwordCred.Value); err != nil {
			return fmt.Errorf("invalid DataCite password: %w", err)
		}
		
		fmt.Printf("  Using repository ID from %s: %s\n", repoIDCred.Source, repoIDCred.Value)
		fmt.Printf("  Using password from %s: %s\n", passwordCred.Source, config.RedactToken(passwordCred.Value))
		
		dataciteConfig := &doi.DataCiteConfig{
			RepositoryID: repoIDCred.Value,
			Password:     passwordCred.Value,
			Sandbox:      sandbox,
		}
		
		provider, err = doi.NewDataCiteProvider(dataciteConfig)
		if err != nil {
			return fmt.Errorf("create DataCite provider: %w", err)
		}
		
	default:
		return fmt.Errorf("unknown provider: %s (supported: zenodo, datacite)", providerName)
	}

	// Step 3: Extract metadata
	fmt.Println("\n→ Extracting metadata...")
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()
	
	extractedMeta, err := registry.Extract(filePath)
	if err != nil {
		return fmt.Errorf("extract metadata: %w", err)
	}

	fmt.Printf("  Extracted %d metadata fields\n", len(extractedMeta))

	// Step 4: Load enrichment if provided
	var enrichment map[string]interface{}
	if enrichmentFile != "" {
		fmt.Printf("\n→ Loading enrichment from %s...\n", enrichmentFile)
		enrichData, err := os.ReadFile(enrichmentFile)
		if err != nil {
			return fmt.Errorf("read enrichment file: %w", err)
		}
		
		if err := yaml.Unmarshal(enrichData, &enrichment); err != nil {
			if err := json.Unmarshal(enrichData, &enrichment); err != nil {
				return fmt.Errorf("parse enrichment file: %w", err)
			}
		}
	}

	// Step 5: Prepare dataset
	fmt.Println("\n→ Preparing DOI metadata...")
	
	providerRegistry := doi.NewProviderRegistry()
	providerRegistry.Register(provider)
	_ = providerRegistry.SetActive(provider.Name())
	
	workflowConfig := &doi.WorkflowConfig{
		Publisher:          publisher,
		License:            license,
		MinQualityScore:    60.0,
		RequireRealAuthors: true,
		RequireDescription: true,
	}
	
	workflow := doi.NewDOIWorkflow(workflowConfig, providerRegistry)
	
	prepReq := &doi.PrepareRequest{
		FilePath:   filePath,
		Metadata:   extractedMeta,
		Enrichment: enrichment,
		PresetID:   presetID,
	}
	
	result, err := workflow.Prepare(prepReq)
	if err != nil {
		return fmt.Errorf("prepare metadata: %w", err)
	}

	// Step 6: Validate
	fmt.Printf("\n→ Validating metadata...\n")
	fmt.Printf("  Quality Score: %.1f/100 (%s)\n", result.Validation.Score, doi.GetQualityLevel(result.Validation.Score))
	
	if len(result.Validation.Errors) > 0 {
		fmt.Printf("  ✗ %d errors found:\n", len(result.Validation.Errors))
		for _, err := range result.Validation.Errors {
			fmt.Printf("    • %s\n", err)
		}
		return fmt.Errorf("validation failed")
	}
	
	if len(result.Validation.Warnings) > 0 {
		fmt.Printf("  ⚠ %d warnings:\n", len(result.Validation.Warnings))
		for _, warning := range result.Validation.Warnings {
			fmt.Printf("    • %s\n", warning)
		}
	}
	
	if result.Validation.IsReady {
		fmt.Printf("  ✓ Ready for DOI minting\n")
	}

	// Show dataset info
	fmt.Printf("\n→ Dataset Information:\n")
	fmt.Printf("  Title: %s\n", result.Dataset.Title)
	fmt.Printf("  Authors: %d\n", len(result.Dataset.Authors))
	for i, author := range result.Dataset.Authors {
		fmt.Printf("    %d. %s", i+1, author.Name)
		if author.ORCID != "" {
			fmt.Printf(" (ORCID: %s)", author.ORCID)
		}
		if author.Affiliation != "" {
			fmt.Printf(" - %s", author.Affiliation)
		}
		fmt.Println()
	}
	fmt.Printf("  Publisher: %s\n", result.Dataset.Publisher)
	fmt.Printf("  License: %s\n", result.Dataset.License)

	// Step 7: Dry run check
	if dryRun {
		fmt.Println("\n✓ DRY RUN COMPLETE")
		fmt.Println("  Metadata is valid and ready for minting.")
		fmt.Println("  Remove --dry-run flag to actually mint the DOI.")
		return nil
	}

	// Step 8: Mint DOI
	fmt.Println("\n→ Minting DOI...")
	fmt.Println("  This may take a few moments...")
	
	startTime := time.Now()
	mintedDOI, err := provider.Mint(ctx, result.Dataset)
	if err != nil {
		return fmt.Errorf("mint DOI: %w", err)
	}
	duration := time.Since(startTime)

	// Step 9: Success!
	fmt.Printf("\n✓ DOI MINTED SUCCESSFULLY in %.1fs\n\n", duration.Seconds())
	fmt.Printf("DOI: %s\n", mintedDOI.DOI)
	if mintedDOI.URL != "" {
		fmt.Printf("URL: %s\n", mintedDOI.URL)
	}
	fmt.Printf("State: %s\n", mintedDOI.State)
	fmt.Printf("Created: %s\n", mintedDOI.CreatedAt.Format(time.RFC3339))
	
	fmt.Println("\n🎉 Your data now has a permanent identifier!")
	fmt.Println("   You can cite this DOI in publications.")
	
	if sandbox {
		fmt.Println("\n⚠ NOTE: This is a SANDBOX DOI (test only)")
		fmt.Println("   Remove --sandbox flag to mint a production DOI.")
	}

	return nil
}
