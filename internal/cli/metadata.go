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
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/scttfrdmn/cicada/internal/metadata"
)

// NewMetadataCmd creates the metadata command.
func NewMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Extract and manage file metadata",
		Long: `Extract, display, and validate metadata from scientific instrument files.

Cicada automatically detects file formats and uses the appropriate extractor
for microscopy, sequencing, mass spec, and other scientific data formats.`,
	}

	// Add subcommands
	cmd.AddCommand(newMetadataExtractCmd())
	cmd.AddCommand(newMetadataShowCmd())
	cmd.AddCommand(newMetadataValidateCmd())
	cmd.AddCommand(newMetadataListCmd())
	cmd.AddCommand(newMetadataPresetCmd())

	return cmd
}

// newMetadataExtractCmd creates the metadata extract subcommand.
func newMetadataExtractCmd() *cobra.Command {
	var (
		outputFormat string
		outputFile   string
		extractorName string
	)

	cmd := &cobra.Command{
		Use:   "extract <path>",
		Short: "Extract metadata from a file",
		Long: `Extract metadata from a scientific instrument file.

Automatically detects the file format and uses the appropriate extractor.
Outputs metadata in JSON, YAML, or table format.

Examples:
  # Extract metadata and display as JSON
  cicada metadata extract data/image.czi

  # Extract and save to file
  cicada metadata extract data/image.czi --output metadata.json

  # Extract with YAML output
  cicada metadata extract data/image.czi --format yaml

  # Force a specific extractor
  cicada metadata extract data/image.czi --extractor zeiss_czi`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", path)
			}

			// Create registry and register default extractors
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			var result map[string]interface{}
			var err error

			if extractorName != "" {
				// Use specific extractor
				extractor := registry.FindExtractor(filepath.Base(path))
				if extractor == nil || extractor.Name() != extractorName {
					return fmt.Errorf("extractor '%s' not found or cannot handle file", extractorName)
				}
				result, err = extractor.Extract(path)
			} else {
				// Auto-detect extractor
				result, err = registry.Extract(path)
			}

			if err != nil {
				return fmt.Errorf("extraction failed: %w", err)
			}

			// Format output
			var output []byte
			switch strings.ToLower(outputFormat) {
			case "json":
				output, err = json.MarshalIndent(result, "", "  ")
			case "yaml":
				output, err = yaml.Marshal(result)
			case "table":
				output = []byte(formatAsTable(result))
			default:
				return fmt.Errorf("unsupported format: %s (use json, yaml, or table)", outputFormat)
			}

			if err != nil {
				return fmt.Errorf("format output: %w", err)
			}

			// Write output
			if outputFile != "" {
				if err := os.WriteFile(outputFile, output, 0644); err != nil {
					return fmt.Errorf("write output file: %w", err)
				}
				fmt.Printf("Metadata extracted to %s\n", outputFile)
			} else {
				fmt.Println(string(output))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, yaml, table)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringVar(&extractorName, "extractor", "", "Force specific extractor")

	return cmd
}

// newMetadataShowCmd creates the metadata show subcommand.
func newMetadataShowCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "show <path>",
		Short: "Show metadata in human-readable format",
		Long: `Display metadata from a file in a human-readable format.

Similar to 'extract' but optimized for readability with colored output
and formatted tables.

Examples:
  # Show metadata in table format
  cicada metadata show data/image.czi

  # Show as JSON
  cicada metadata show data/image.czi --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", path)
			}

			// Create registry and extract metadata
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			result, err := registry.Extract(path)
			if err != nil {
				return fmt.Errorf("extraction failed: %w", err)
			}

			// Display based on format
			switch strings.ToLower(format) {
			case "table":
				fmt.Println(formatAsTable(result))
			case "json":
				output, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(output))
			case "yaml":
				output, _ := yaml.Marshal(result)
				fmt.Println(string(output))
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")

	return cmd
}

// newMetadataValidateCmd creates the metadata validate subcommand.
func newMetadataValidateCmd() *cobra.Command {
	var presetID string

	cmd := &cobra.Command{
		Use:   "validate <path>",
		Short: "Validate file metadata",
		Long: `Validate that a file can be read and has valid metadata.

Checks:
  - File is readable
  - Format is recognized
  - Metadata can be extracted
  - Required fields are present

With --preset flag:
  - Validates against instrument-specific requirements
  - Checks required and optional fields
  - Provides quality score (0-100)

Examples:
  # Validate a single file
  cicada metadata validate data/image.czi

  # Validate multiple files
  cicada metadata validate data/*.czi

  # Validate against a preset
  cicada metadata validate data/image.czi --preset zeiss-lsm-880

  # Validate Illumina FASTQ
  cicada metadata validate data/sample_R1.fastq.gz --preset illumina-novaseq`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			extractorRegistry := metadata.NewExtractorRegistry()
			extractorRegistry.RegisterDefaults()

			var presetRegistry *metadata.PresetRegistry
			var preset *metadata.InstrumentPreset
			if presetID != "" {
				presetRegistry = metadata.NewPresetRegistry()
				presetRegistry.RegisterDefaults()
				var err error
				preset, err = presetRegistry.GetPreset(presetID)
				if err != nil {
					return fmt.Errorf("preset not found: %s", presetID)
				}
			}

			var hasErrors bool

			for _, path := range args {
				// Check if file exists
				if _, err := os.Stat(path); os.IsNotExist(err) {
					fmt.Printf("❌ %s: file not found\n", path)
					hasErrors = true
					continue
				}

				// Try to extract metadata
				result, err := extractorRegistry.Extract(path)
				if err != nil {
					fmt.Printf("❌ %s: %v\n", path, err)
					hasErrors = true
					continue
				}

				// Validate against preset if specified
				if preset != nil {
					validation := preset.Validate(result)
					if !validation.IsValid {
						fmt.Printf("❌ %s: validation failed\n", path)
						for _, err := range validation.Errors {
							fmt.Printf("     Error: %s\n", err)
						}
						for _, warn := range validation.Warnings {
							fmt.Printf("     Warning: %s\n", warn)
						}
						fmt.Printf("     Quality Score: %.1f/100\n", validation.QualityScore())
						hasErrors = true
					} else {
						score := validation.QualityScore()
						fmt.Printf("✓ %s: valid (%s)\n", path, result["format"])
						if len(validation.Warnings) > 0 {
							fmt.Printf("     %d warnings\n", len(validation.Warnings))
							for _, warn := range validation.Warnings {
								fmt.Printf("       - %s\n", warn)
							}
						}
						fmt.Printf("     Quality Score: %.1f/100\n", score)
					}
				} else {
					// Basic validation (no preset)
					requiredFields := []string{"format", "file_name"}
					var missingFields []string
					for _, field := range requiredFields {
						if _, ok := result[field]; !ok {
							missingFields = append(missingFields, field)
						}
					}

					if len(missingFields) > 0 {
						fmt.Printf("⚠️  %s: missing fields: %v\n", path, missingFields)
						hasErrors = true
					} else {
						fmt.Printf("✓ %s: valid (%s)\n", path, result["format"])
					}
				}
			}

			if hasErrors {
				return fmt.Errorf("validation failed for one or more files")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&presetID, "preset", "p", "", "Validate against instrument preset")

	return cmd
}

// newMetadataListCmd creates the metadata list subcommand.
func newMetadataListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available metadata extractors",
		Long: `List all registered metadata extractors and supported formats.

Shows which file formats are supported and which extractor handles each format.

Example:
  cicada metadata list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			extractors := registry.ListExtractors()

			fmt.Println("Available Metadata Extractors:")
			fmt.Println()

			for _, ext := range extractors {
				fmt.Printf("  %s\n", ext.Name)
				fmt.Printf("    Formats: %s\n", strings.Join(ext.Formats, ", "))
				fmt.Println()
			}

			fmt.Printf("Total: %d extractors\n", len(extractors))

			return nil
		},
	}

	return cmd
}

// formatAsTable formats metadata as a human-readable table.
func formatAsTable(data map[string]interface{}) string {
	var sb strings.Builder

	// Find max key length for alignment
	maxKeyLen := 0
	for key := range data {
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}
	}

	// Format each field
	for key, value := range data {
		// Skip complex types (arrays, nested objects)
		switch v := value.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("%-*s : %s\n", maxKeyLen, key, v))
		case int, int64, float64, bool:
			sb.WriteString(fmt.Sprintf("%-*s : %v\n", maxKeyLen, key, v))
		case []interface{}:
			sb.WriteString(fmt.Sprintf("%-*s : [%d items]\n", maxKeyLen, key, len(v)))
		default:
			sb.WriteString(fmt.Sprintf("%-*s : %T\n", maxKeyLen, key, v))
		}
	}

	return sb.String()
}

// newMetadataPresetCmd creates the metadata preset command.
func newMetadataPresetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preset",
		Short: "Manage instrument presets",
		Long: `Manage instrument presets for metadata validation.

Presets define expected metadata fields for specific instruments and enable
validation of extracted metadata against instrument specifications.

Examples:
  # List all available presets
  cicada metadata preset list

  # Show details of a specific preset
  cicada metadata preset show zeiss-lsm-880

  # Validate metadata against a preset
  cicada metadata validate data/image.czi --preset zeiss-lsm-880`,
	}

	// Add subcommands
	cmd.AddCommand(newMetadataPresetListCmd())
	cmd.AddCommand(newMetadataPresetShowCmd())

	return cmd
}

// newMetadataPresetListCmd creates the preset list subcommand.
func newMetadataPresetListCmd() *cobra.Command {
	var (
		outputFormat string
		manufacturer string
		instrumentType string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available instrument presets",
		Long: `List all available instrument presets.

Presets can be filtered by manufacturer or instrument type.

Examples:
  # List all presets
  cicada metadata preset list

  # List Zeiss presets
  cicada metadata preset list --manufacturer Zeiss

  # List microscopy presets
  cicada metadata preset list --type microscopy

  # List in YAML format
  cicada metadata preset list --format yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create registry and register defaults
			registry := metadata.NewPresetRegistry()
			registry.RegisterDefaults()

			// Find presets
			var presets []*metadata.InstrumentPreset
			if manufacturer != "" || instrumentType != "" {
				presets = registry.FindPresets(manufacturer, instrumentType)
			} else {
				presets = registry.ListPresets()
			}

			// Format output
			switch strings.ToLower(outputFormat) {
			case "json":
				output, err := json.MarshalIndent(presets, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal JSON: %w", err)
				}
				fmt.Println(string(output))

			case "yaml":
				output, err := yaml.Marshal(presets)
				if err != nil {
					return fmt.Errorf("marshal YAML: %w", err)
				}
				fmt.Println(string(output))

			case "table", "":
				fmt.Println("Available Instrument Presets:")
				fmt.Println()

				for _, preset := range presets {
					fmt.Printf("  %s\n", preset.Name)
					fmt.Printf("    ID: %s\n", preset.ID)
					fmt.Printf("    Manufacturer: %s\n", preset.Manufacturer)
					fmt.Printf("    Type: %s\n", preset.InstrumentType)
					if len(preset.Models) > 0 {
						fmt.Printf("    Models: %s\n", strings.Join(preset.Models, ", "))
					}
					if len(preset.FileFormats) > 0 {
						fmt.Printf("    Formats: %s\n", strings.Join(preset.FileFormats, ", "))
					}
					fmt.Println()
				}

				fmt.Printf("Total: %d presets\n", len(presets))

			default:
				return fmt.Errorf("unsupported format: %s (use json, yaml, or table)", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format (json, yaml, table)")
	cmd.Flags().StringVarP(&manufacturer, "manufacturer", "m", "", "Filter by manufacturer")
	cmd.Flags().StringVarP(&instrumentType, "type", "t", "", "Filter by instrument type")

	return cmd
}

// newMetadataPresetShowCmd creates the preset show subcommand.
func newMetadataPresetShowCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "show <preset-id>",
		Short: "Show details of an instrument preset",
		Long: `Show detailed information about an instrument preset.

Displays required fields, optional fields, validation rules, and examples.

Examples:
  # Show Zeiss LSM 880 preset
  cicada metadata preset show zeiss-lsm-880

  # Show in JSON format
  cicada metadata preset show illumina-novaseq --format json

  # Show generic microscopy preset
  cicada metadata preset show generic-microscopy`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			presetID := args[0]

			// Create registry and register defaults
			registry := metadata.NewPresetRegistry()
			registry.RegisterDefaults()

			// Get preset
			preset, err := registry.GetPreset(presetID)
			if err != nil {
				return fmt.Errorf("preset not found: %s", presetID)
			}

			// Format output
			switch strings.ToLower(outputFormat) {
			case "json":
				output, err := json.MarshalIndent(preset, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal JSON: %w", err)
				}
				fmt.Println(string(output))

			case "yaml":
				output, err := yaml.Marshal(preset)
				if err != nil {
					return fmt.Errorf("marshal YAML: %w", err)
				}
				fmt.Println(string(output))

			case "table", "":
				fmt.Printf("Preset: %s\n", preset.Name)
				fmt.Println()
				fmt.Printf("ID:           %s\n", preset.ID)
				fmt.Printf("Manufacturer: %s\n", preset.Manufacturer)
				fmt.Printf("Type:         %s\n", preset.InstrumentType)
				if len(preset.Models) > 0 {
					fmt.Printf("Models:       %s\n", strings.Join(preset.Models, ", "))
				}
				if len(preset.FileFormats) > 0 {
					fmt.Printf("Formats:      %s\n", strings.Join(preset.FileFormats, ", "))
				}
				if preset.Description != "" {
					fmt.Printf("Description:  %s\n", preset.Description)
				}
				fmt.Println()

				if len(preset.RequiredFields) > 0 {
					fmt.Println("Required Fields:")
					for _, field := range preset.RequiredFields {
						fmt.Printf("  - %s (%s)\n", field.Name, field.Type)
						if field.Description != "" {
							fmt.Printf("      %s\n", field.Description)
						}
						if len(field.Enum) > 0 {
							fmt.Printf("      Allowed values: %s\n", strings.Join(field.Enum, ", "))
						}
						if field.MinValue != nil || field.MaxValue != nil {
							if field.MinValue != nil && field.MaxValue != nil {
								fmt.Printf("      Range: %.2f - %.2f\n", *field.MinValue, *field.MaxValue)
							} else if field.MinValue != nil {
								fmt.Printf("      Minimum: %.2f\n", *field.MinValue)
							} else {
								fmt.Printf("      Maximum: %.2f\n", *field.MaxValue)
							}
						}
						if field.Example != nil {
							fmt.Printf("      Example: %v\n", field.Example)
						}
					}
					fmt.Println()
				}

				if len(preset.OptionalFields) > 0 {
					fmt.Println("Optional Fields:")
					for _, field := range preset.OptionalFields {
						fmt.Printf("  - %s (%s)\n", field.Name, field.Type)
						if field.Description != "" {
							fmt.Printf("      %s\n", field.Description)
						}
					}
					fmt.Println()
				}

				if len(preset.References) > 0 {
					fmt.Println("References:")
					for _, ref := range preset.References {
						fmt.Printf("  - %s\n", ref)
					}
					fmt.Println()
				}

			default:
				return fmt.Errorf("unsupported format: %s (use json, yaml, or table)", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format (json, yaml, table)")

	return cmd
}
