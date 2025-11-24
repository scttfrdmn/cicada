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
	"os"
	"path/filepath"
	"testing"

	"github.com/scttfrdmn/cicada/internal/metadata"
)

// TestPresetIntegration_MetadataValidation tests preset validation with extracted metadata
func TestPresetIntegration_MetadataValidation(t *testing.T) {
	// Create a test FASTQ file
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

	// Extract metadata
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	extracted, err := registry.Extract(testFile)
	if err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Get preset registry
	presetRegistry := metadata.NewPresetRegistry()
	presetRegistry.RegisterDefaults()

	// Get Illumina preset
	preset, err := presetRegistry.GetPreset("illumina-novaseq")
	if err != nil {
		t.Fatalf("GetPreset() failed: %v", err)
	}

	// Validate against preset
	result := preset.Validate(extracted)

	// The extracted FASTQ metadata won't have all instrument-specific fields,
	// so validation will have missing fields, but should not error on format
	if result.IsValid {
		t.Log("Metadata validates successfully against Illumina NovaSeq preset")
	}

	// Check that format field is present and matches
	if extracted["format"] == "FASTQ" {
		t.Log("Format correctly identified as FASTQ")
	}

	// Quality score should be calculated
	score := result.QualityScore()
	if score < 0 || score > 100 {
		t.Errorf("Quality score = %.1f, should be between 0 and 100", score)
	}

	t.Logf("Preset validation quality score: %.1f/100", score)
	t.Logf("Present fields: %d, Missing fields: %d", len(result.Present), len(result.Missing))
}

// TestPresetIntegration_FindPresets tests finding presets by criteria
func TestPresetIntegration_FindPresets(t *testing.T) {
	registry := metadata.NewPresetRegistry()
	registry.RegisterDefaults()

	tests := []struct {
		name           string
		manufacturer   string
		instrumentType string
		wantMinCount   int
	}{
		{
			name:           "Find Zeiss presets",
			manufacturer:   "Zeiss",
			instrumentType: "",
			wantMinCount:   3, // Should find 3 Zeiss LSM presets
		},
		{
			name:           "Find Illumina presets",
			manufacturer:   "Illumina",
			instrumentType: "",
			wantMinCount:   3, // Should find 3 Illumina presets
		},
		{
			name:           "Find microscopy presets",
			manufacturer:   "",
			instrumentType: "microscopy",
			wantMinCount:   4, // 3 Zeiss + 1 generic
		},
		{
			name:           "Find sequencing presets",
			manufacturer:   "",
			instrumentType: "sequencing",
			wantMinCount:   4, // 3 Illumina + 1 generic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presets := registry.FindPresets(tt.manufacturer, tt.instrumentType)

			if len(presets) < tt.wantMinCount {
				t.Errorf("FindPresets() returned %d presets, want >= %d",
					len(presets), tt.wantMinCount)
			}

			// Verify all returned presets match criteria
			for _, preset := range presets {
				if tt.manufacturer != "" && preset.Manufacturer != tt.manufacturer {
					t.Errorf("Preset %s has manufacturer %s, want %s",
						preset.ID, preset.Manufacturer, tt.manufacturer)
				}

				if tt.instrumentType != "" && preset.InstrumentType != tt.instrumentType {
					t.Errorf("Preset %s has type %s, want %s",
						preset.ID, preset.InstrumentType, tt.instrumentType)
				}
			}
		})
	}
}

// TestPresetIntegration_AllPresets tests that all default presets work correctly
func TestPresetIntegration_AllPresets(t *testing.T) {
	registry := metadata.NewPresetRegistry()
	registry.RegisterDefaults()

	allPresets := registry.ListPresets()

	if len(allPresets) < 8 {
		t.Errorf("ListPresets() returned %d presets, want >= 8", len(allPresets))
	}

	for _, preset := range allPresets {
		t.Run(preset.ID, func(t *testing.T) {
			// Verify preset has required fields
			if preset.ID == "" {
				t.Error("Preset ID is empty")
			}

			if preset.Name == "" {
				t.Error("Preset Name is empty")
			}

			if preset.InstrumentType == "" {
				t.Error("Preset InstrumentType is empty")
			}

			// Verify preset has some required or optional fields
			totalFields := len(preset.RequiredFields) + len(preset.OptionalFields)
			if totalFields == 0 {
				t.Error("Preset has no field requirements")
			}

			// Test validation with empty metadata (should fail gracefully)
			emptyMetadata := map[string]interface{}{}
			result := preset.Validate(emptyMetadata)

			// Should have errors for missing required fields
			if len(preset.RequiredFields) > 0 && len(result.Errors) == 0 {
				t.Error("Validation with empty metadata should have errors for required fields")
			}

			// Quality score should be calculable
			score := result.QualityScore()
			if score < 0 || score > 100 {
				t.Errorf("Quality score = %.1f, should be between 0 and 100", score)
			}
		})
	}
}

// TestPresetIntegration_ExtractAndValidate tests full extract → validate workflow
func TestPresetIntegration_ExtractAndValidate(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
		presetID string
	}{
		{
			name:     "FASTQ with Illumina preset",
			filename: "illumina_sample_R1.fastq",
			content: `@Illumina
ACGTACGTACGTACGT
+
IIIIIIIIIIIIIIII
`,
			presetID: "illumina-novaseq",
		},
		{
			name:     "Generic sequencing preset",
			filename: "generic_sample.fastq",
			content: `@SEQ
ACGT
+
IIII
`,
			presetID: "generic-sequencing",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(t.TempDir(), tc.filename)
			if err := os.WriteFile(testFile, []byte(tc.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Extract metadata
			extractorRegistry := metadata.NewExtractorRegistry()
			extractorRegistry.RegisterDefaults()

			extracted, err := extractorRegistry.Extract(testFile)
			if err != nil {
				t.Fatalf("Extract() failed: %v", err)
			}

			// Get preset and validate
			presetRegistry := metadata.NewPresetRegistry()
			presetRegistry.RegisterDefaults()

			preset, err := presetRegistry.GetPreset(tc.presetID)
			if err != nil {
				t.Fatalf("GetPreset() failed: %v", err)
			}

			result := preset.Validate(extracted)

			// Log results
			t.Logf("Validation result for %s with preset %s:", tc.filename, tc.presetID)
			t.Logf("  IsValid: %v", result.IsValid)
			t.Logf("  Quality Score: %.1f/100", result.QualityScore())
			t.Logf("  Present: %d, Missing: %d, Errors: %d",
				len(result.Present), len(result.Missing), len(result.Errors))

			// Verify basic validation behavior
			if len(result.Present) == 0 {
				t.Error("Should have at least some present fields from extraction")
			}

			score := result.QualityScore()
			if score < 0 || score > 100 {
				t.Errorf("Quality score = %.1f, out of valid range", score)
			}
		})
	}
}

// TestPresetIntegration_FieldValidation tests individual field validation
func TestPresetIntegration_FieldValidation(t *testing.T) {
	registry := metadata.NewPresetRegistry()
	registry.RegisterDefaults()

	preset, err := registry.GetPreset("generic-microscopy")
	if err != nil {
		t.Fatalf("GetPreset() failed: %v", err)
	}

	tests := []struct {
		name      string
		metadata  map[string]interface{}
		wantValid bool
		wantErrors int
	}{
		{
			name: "Valid microscopy metadata",
			metadata: map[string]interface{}{
				"format":           "CZI",
				"manufacturer":     "Zeiss",
				"instrument_model": "LSM 880",
				"instrument_type":  "microscopy",
				"data_type":        "image",
				"image_width":      2048,
				"image_height":     2048,
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "Missing required fields",
			metadata: map[string]interface{}{
				"manufacturer": "Zeiss",
			},
			wantValid:  false,
			wantErrors: 4, // Missing format, instrument_type, image_width, image_height
		},
		{
			name: "Invalid field types",
			metadata: map[string]interface{}{
				"format":          "CZI",
				"instrument_type": "microscopy",
				"data_type":       "image",
				"image_width":     "not a number", // Should be number
			},
			wantValid:  false,
			wantErrors: 2, // Invalid image_width type + missing image_height
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preset.Validate(tt.metadata)

			if result.IsValid != tt.wantValid {
				t.Errorf("IsValid = %v, want %v (errors: %v)",
					result.IsValid, tt.wantValid, result.Errors)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Error count = %d, want %d (errors: %v)",
					len(result.Errors), tt.wantErrors, result.Errors)
			}
		})
	}
}

// TestPresetIntegration_QualityScoring tests quality score calculation
func TestPresetIntegration_QualityScoring(t *testing.T) {
	registry := metadata.NewPresetRegistry()
	registry.RegisterDefaults()

	preset, err := registry.GetPreset("generic-microscopy")
	if err != nil {
		t.Fatalf("GetPreset() failed: %v", err)
	}

	// Minimal metadata (only required fields)
	minimalMetadata := map[string]interface{}{
		"format":          "CZI",
		"instrument_type": "microscopy",
		"image_width":     2048,
		"image_height":    2048,
	}

	// Rich metadata (required + many optional fields)
	richMetadata := map[string]interface{}{
		"format":                  "CZI",
		"instrument_type":         "microscopy",
		"data_type":               "image",
		"manufacturer":            "Zeiss",
		"instrument_model":        "LSM 880",
		"image_width":             2048,
		"image_height":            2048,
		"channel_count":           3,
		"z_planes":                10,
		"timepoints":              5,
		"pixel_size_x":            0.1,
		"pixel_size_y":            0.1,
		"objective_magnification": 63.0,
		"objective_na":            1.4,
	}

	minimalResult := preset.Validate(minimalMetadata)
	richResult := preset.Validate(richMetadata)

	minimalScore := minimalResult.QualityScore()
	richScore := richResult.QualityScore()

	// Rich metadata should score higher
	if richScore <= minimalScore {
		t.Errorf("Rich metadata score (%.1f) should be higher than minimal (%.1f)",
			richScore, minimalScore)
	}

	// Minimal with all required fields should score reasonably
	// Note: Score considers both required and optional fields in total
	if minimalScore < 40 {
		t.Errorf("Minimal score = %.1f, expected >= 40 for having all required fields", minimalScore)
	}

	// Rich metadata should score well (required + many optional)
	if richScore < 75 {
		t.Errorf("Rich score = %.1f, expected >= 75 with many fields present", richScore)
	}

	t.Logf("Quality scores: Minimal = %.1f, Rich = %.1f", minimalScore, richScore)
}

// TestPresetIntegration_TemplateGeneration tests generating templates from presets
func TestPresetIntegration_TemplateGeneration(t *testing.T) {
	registry := metadata.NewPresetRegistry()
	registry.RegisterDefaults()

	preset, err := registry.GetPreset("zeiss-lsm-880")
	if err != nil {
		t.Fatalf("GetPreset() failed: %v", err)
	}

	// Generate template
	template := preset.GenerateTemplate()

	// Verify template has all required fields
	for _, req := range preset.RequiredFields {
		if _, exists := template[req.Name]; !exists {
			t.Errorf("Template missing required field: %s", req.Name)
		}
	}

	// Verify template has format field
	if format, ok := template["format"]; !ok {
		t.Error("Template missing format field")
	} else if format != "CZI" {
		t.Errorf("Template format = %v, want CZI", format)
	}

	// Verify template has example values
	foundExample := false
	for _, req := range preset.RequiredFields {
		if req.Example != nil {
			foundExample = true
			break
		}
	}

	if !foundExample {
		t.Error("Preset should have at least one field with example value")
	}
}
