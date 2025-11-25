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
	"os"
	"path/filepath"
	"testing"
)

// TestMetadataExtractCmd tests the metadata extract command.
// Note: Commands write directly to stdout/stderr, so we only verify execution.
func TestMetadataExtractCmd(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		args []string
	}{
		{"JSON format", []string{"extract", testFile, "--format", "json"}},
		{"YAML format", []string{"extract", testFile, "--format", "yaml"}},
		{"Table format", []string{"extract", testFile, "--format", "table"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewMetadataCmd()
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Command failed: %v", err)
			}
		})
	}

	t.Run("Extract to output file", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "metadata.json")
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"extract", testFile, "--format", "json", "--output", outputFile})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Error("Output file was not created")
		}
	})

	t.Run("Extract non-existent file", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"extract", "/nonexistent/file.txt"})

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("Extract with invalid format", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"extract", testFile, "--format", "invalid"})

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for invalid format")
		}
	})
}

// TestMetadataShowCmd tests the metadata show command.
func TestMetadataShowCmd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		args []string
	}{
		{"Table format", []string{"show", testFile, "--format", "table"}},
		{"JSON format", []string{"show", testFile, "--format", "json"}},
		{"YAML format", []string{"show", testFile, "--format", "yaml"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewMetadataCmd()
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Command failed: %v", err)
			}
		})
	}

	t.Run("Show non-existent file", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"show", "/nonexistent/file.txt"})

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

// TestMetadataValidateCmd tests the metadata validate command.
func TestMetadataValidateCmd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("Validate file", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"validate", testFile})

		// May fail validation but shouldn't crash
		_ = cmd.Execute()
	})

	t.Run("Validate non-existent file", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"validate", "/nonexistent/file.txt"})

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("Validate multiple files", func(t *testing.T) {
		testFile2 := filepath.Join(tmpDir, "test2.txt")
		if err := os.WriteFile(testFile2, []byte("test content 2"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"validate", testFile, testFile2})

		// May fail validation but shouldn't crash
		_ = cmd.Execute()
	})
}

// TestMetadataListCmd tests the metadata list command.
func TestMetadataListCmd(t *testing.T) {
	t.Run("List extractors", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"list"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}
	})
}

// TestMetadataPresetCmd tests the metadata preset commands.
func TestMetadataPresetCmd(t *testing.T) {
	t.Run("List presets", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"preset", "list"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}
	})

	t.Run("List presets with JSON format", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"preset", "list", "--format", "json"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}
	})

	t.Run("List presets with YAML format", func(t *testing.T) {
		cmd := NewMetadataCmd()
		cmd.SetArgs([]string{"preset", "list", "--format", "yaml"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}
	})
}

// TestFormatAsTable tests the table formatting function.
func TestFormatAsTable(t *testing.T) {
	t.Run("Format simple data", func(t *testing.T) {
		data := map[string]interface{}{
			"format":    "TXT",
			"file_name": "test.txt",
		}

		output := formatAsTable(data)
		if len(output) == 0 {
			t.Error("Expected non-empty table output")
		}
	})

	t.Run("Format with numeric values", func(t *testing.T) {
		data := map[string]interface{}{
			"width":  1024,
			"height": 768,
		}

		output := formatAsTable(data)
		if len(output) == 0 {
			t.Error("Expected non-empty table output")
		}
	})
}
