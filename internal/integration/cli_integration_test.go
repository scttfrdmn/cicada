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

package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scttfrdmn/cicada/internal/cli"
)

// TestCLI_MetadataExtract tests the metadata extract command
func TestCLI_MetadataExtract(t *testing.T) {
	// Create test FASTQ file
	testFile := filepath.Join(t.TempDir(), "sample_R1.fastq")
	fastqContent := `@Illumina_NovaSeq_Run123
GATTTGGGGTTCAAAGCAGTATCGATCAAATAGTAAATCCATTTGTTCAACTCACAGTTT
+
IIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIII
@Illumina_NovaSeq_Run123
GCGCGCGCGCGCGCGCGCGCGCGC
+
IIIIIIIIIIIIIIIIIIIIIIIII
`
	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create root command
	rootCmd := cli.NewRootCmd("test")

	// Capture output
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)

	// Create output file
	outputFile := filepath.Join(t.TempDir(), "output.json")

	// Run metadata extract command with output file
	rootCmd.SetArgs([]string{"metadata", "extract", testFile, "--format", "json", "--output", outputFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("metadata extract failed: %v", err)
	}

	// Read output file
	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Parse JSON output
	var metadata map[string]interface{}
	if err := json.Unmarshal(outputData, &metadata); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, string(outputData))
	}

	// Verify metadata fields
	if metadata["format"] != "FASTQ" {
		t.Errorf("format = %v, want FASTQ", metadata["format"])
	}

	if metadata["total_reads"] != float64(2) { // JSON numbers are float64
		t.Errorf("total_reads = %v, want 2", metadata["total_reads"])
	}

	if metadata["is_paired_end"] != true {
		t.Errorf("is_paired_end = %v, want true", metadata["is_paired_end"])
	}
}

// TestCLI_MetadataValidate tests the metadata validate command with presets
func TestCLI_MetadataValidate(t *testing.T) {
	// Create test FASTQ file
	testFile := filepath.Join(t.TempDir(), "illumina_sample.fastq")
	fastqContent := `@Illumina
ACGTACGTACGTACGT
+
IIIIIIIIIIIIIIII
`
	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create root command
	rootCmd := cli.NewRootCmd("test")

	// Run metadata validate command - should succeed with valid file
	rootCmd.SetArgs([]string{"metadata", "validate", testFile, "--preset", "generic-sequencing"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("metadata validate failed: %v", err)
	}

	// Command completed successfully, which means validation passed
	// (The command would return error if validation failed critically)
}

// TestCLI_MetadataListPresets tests listing available presets
func TestCLI_MetadataListPresets(t *testing.T) {
	// Create root command
	rootCmd := cli.NewRootCmd("test")

	// Run metadata preset list command - should succeed
	rootCmd.SetArgs([]string{"metadata", "preset", "list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("metadata preset list failed: %v", err)
	}

	// Command completed successfully, which means presets were listed
	// Note: Output goes to stdout and cannot be easily captured in tests
}

// TestCLI_DOIPrepare tests the DOI prepare command
func TestCLI_DOIPrepare(t *testing.T) {
	tempDir := t.TempDir()

	// Create test FASTQ file
	testFile := filepath.Join(tempDir, "research_data.fastq")
	fastqContent := `@SEQ_ID
ACGTACGTACGTACGT
+
IIIIIIIIIIIIIIII
`
	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create enrichment file with title and authors
	enrichFile := filepath.Join(tempDir, "enrich.json")
	enrichContent := `{
		"title": "Test Dataset",
		"authors": [
			{
				"name": "Test Researcher",
				"affiliation": "Test Lab"
			}
		],
		"description": "A test dataset for integration testing"
	}`
	if err := os.WriteFile(enrichFile, []byte(enrichContent), 0644); err != nil {
		t.Fatalf("Failed to create enrichment file: %v", err)
	}

	// Create root command
	rootCmd := cli.NewRootCmd("test")

	// Capture output
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)

	// Run DOI prepare command - should succeed with enrichment
	rootCmd.SetArgs([]string{
		"doi", "prepare",
		testFile,
		"--enrich", enrichFile,
		"--publisher", "Test Lab",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("doi prepare failed: %v", err)
	}

	// Command completed successfully, which means DOI preparation succeeded
}

// TestCLI_DOIValidate tests the DOI validate command
func TestCLI_DOIValidate(t *testing.T) {
	// Create test FASTQ file
	testFile := filepath.Join(t.TempDir(), "data.fastq")
	fastqContent := `@SEQ
ACGT
+
IIII
`
	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create root command
	rootCmd := cli.NewRootCmd("test")

	// Run DOI validate command - expect it to fail since file lacks proper metadata
	rootCmd.SetArgs([]string{
		"doi", "validate",
		testFile,
	})

	err := rootCmd.Execute()
	// Validation should fail because the file doesn't have required metadata (authors, etc.)
	if err == nil {
		t.Error("doi validate should fail for file without proper metadata")
	}

	// The error message should indicate validation failure
	if !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("Error should mention 'validation failed', got: %v", err)
	}
}

// TestCLI_VersionCommand tests the version command
func TestCLI_VersionCommand(t *testing.T) {
	// The version command uses fmt.Printf which goes directly to os.Stdout
	// and cannot be captured via cobra's SetOut. We test this differently
	// by verifying the command runs without error
	rootCmd := cli.NewRootCmd("1.2.3-test")

	// Capture any errors
	var stderr bytes.Buffer
	rootCmd.SetErr(&stderr)

	// Run version command
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	// Version command should complete without errors
	if stderr.Len() > 0 {
		t.Errorf("Version command produced errors: %s", stderr.String())
	}
}

// TestCLI_InvalidFile tests error handling for non-existent files
func TestCLI_InvalidFile(t *testing.T) {
	// Create root command
	rootCmd := cli.NewRootCmd("test")

	// Capture output
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)

	// Run metadata extract on non-existent file
	nonExistentFile := filepath.Join(t.TempDir(), "does-not-exist.fastq")
	rootCmd.SetArgs([]string{"metadata", "extract", nonExistentFile})

	// Should return error
	err := rootCmd.Execute()
	if err == nil {
		t.Error("metadata extract should fail for non-existent file")
	}
}

// TestCLI_OutputFormats tests different output formats
func TestCLI_OutputFormats(t *testing.T) {
	tempDir := t.TempDir()

	// Create test FASTQ file
	testFile := filepath.Join(tempDir, "sample.fastq")
	fastqContent := `@SEQ
ACGTACGT
+
IIIIIIII
`
	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	formats := []string{"json", "yaml", "table"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			// Create root command
			rootCmd := cli.NewRootCmd("test")

			// Create output file for this format
			outputFile := filepath.Join(tempDir, "output_"+format+".txt")

			// Run metadata extract with specific format
			rootCmd.SetArgs([]string{"metadata", "extract", testFile, "--format", format, "--output", outputFile})

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("metadata extract failed for %s format: %v", format, err)
			}

			// Read and verify output file
			outputData, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			output := string(outputData)
			if output == "" {
				t.Errorf("No output for %s format", format)
			}

			// Verify format-specific content
			if format == "json" && !strings.Contains(output, "{") {
				t.Errorf("%s format should contain JSON brackets", format)
			}
			if format == "yaml" && !strings.Contains(output, ":") {
				t.Errorf("%s format should contain YAML colons", format)
			}
		})
	}
}

// TestCLI_HelpCommand tests help output
func TestCLI_HelpCommand(t *testing.T) {
	// Create root command
	rootCmd := cli.NewRootCmd("test")

	// Capture output
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)

	// Run help command
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("help command failed: %v", err)
	}

	// Verify output contains expected commands
	output := stdout.String()
	expectedCommands := []string{"metadata", "doi", "version", "sync", "watch", "config"}

	for _, cmd := range expectedCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("Help output should mention '%s' command\nOutput: %s", cmd, output)
		}
	}
}
