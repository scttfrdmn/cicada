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
	cmd := &cobra.Command{
		Use:   "validate <path>",
		Short: "Validate file metadata",
		Long: `Validate that a file can be read and has valid metadata.

Checks:
  - File is readable
  - Format is recognized
  - Metadata can be extracted
  - Required fields are present

Examples:
  # Validate a single file
  cicada metadata validate data/image.czi

  # Validate multiple files
  cicada metadata validate data/*.czi`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			var hasErrors bool

			for _, path := range args {
				// Check if file exists
				if _, err := os.Stat(path); os.IsNotExist(err) {
					fmt.Printf("❌ %s: file not found\n", path)
					hasErrors = true
					continue
				}

				// Try to extract metadata
				result, err := registry.Extract(path)
				if err != nil {
					fmt.Printf("❌ %s: %v\n", path, err)
					hasErrors = true
					continue
				}

				// Check for required fields
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

			if hasErrors {
				return fmt.Errorf("validation failed for one or more files")
			}

			return nil
		},
	}

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
