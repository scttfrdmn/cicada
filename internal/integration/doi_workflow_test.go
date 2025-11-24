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

	"github.com/scttfrdmn/cicada/internal/doi"
	"github.com/scttfrdmn/cicada/internal/metadata"
)

// TestDOIWorkflow_EndToEnd tests the complete DOI preparation workflow
func TestDOIWorkflow_EndToEnd(t *testing.T) {
	// Create a test FASTQ file
	testFile := filepath.Join(t.TempDir(), "sample_R1.fastq")
	fastqContent := `@SEQ_ID_1
GATTTGGGGTTCAAAGCAGTATCGATCAAATAGTAAATCCATTTGTTCAACTCACAGTTT
+
!''*((((***+))%%%++)(%%%%).1***-+*''))**55CCF>>>>>>CCCCCCC65
@SEQ_ID_2
GCGCGCGCGCGCGCGCGCGCGCGC
+
IIIIIIIIIIIIIIIIIIIIIIIII
@SEQ_ID_3
ATATATATATATAT
+
HHHHHHHHHHHHHH
`
	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Step 1: Extract metadata
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	extracted, err := registry.Extract(testFile)
	if err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Step 2: Create DOI workflow
	config := &doi.WorkflowConfig{
		Publisher:          "Test Lab",
		License:            "CC-BY-4.0",
		MinQualityScore:    60.0,
		RequireRealAuthors: true,
		RequireDescription: true,
	}

	providerRegistry := doi.NewProviderRegistry()
	providerRegistry.Register(doi.NewDisabledProvider())
	if err := providerRegistry.SetActive("disabled"); err != nil {
		t.Fatalf("SetActive() failed: %v", err)
	}

	workflow := doi.NewDOIWorkflow(config, providerRegistry)

	// Step 3: Prepare metadata for DOI
	prepReq := &doi.PrepareRequest{
		FilePath: testFile,
		Metadata: extracted,
	}

	result, err := workflow.Prepare(prepReq)
	if err != nil {
		t.Fatalf("Prepare() failed: %v", err)
	}

	// Verify dataset was created
	if result.Dataset == nil {
		t.Fatal("Dataset should not be nil")
	}

	// Verify basic dataset fields
	if result.Dataset.Title == "" {
		t.Error("Dataset title should not be empty")
	}

	if len(result.Dataset.Authors) == 0 {
		t.Error("Dataset should have at least one author")
	}

	if result.Dataset.Publisher != "Test Lab" {
		t.Errorf("Publisher = %v, want Test Lab", result.Dataset.Publisher)
	}

	if result.Dataset.ResourceType != "Dataset" {
		t.Errorf("ResourceType = %v, want Dataset", result.Dataset.ResourceType)
	}

	// Verify validation was performed
	if result.Validation == nil {
		t.Fatal("Validation should not be nil")
	}

	// Verify validation includes required checks
	if result.Validation.Score < 0 || result.Validation.Score > 100 {
		t.Errorf("Validation score = %.1f, should be between 0 and 100", result.Validation.Score)
	}

	// Verify quality level
	qualityLevel := doi.GetQualityLevel(result.Validation.Score)
	if qualityLevel == "" {
		t.Error("Quality level should not be empty")
	}
}

// TestDOIWorkflow_WithEnrichment tests DOI preparation with user enrichment
func TestDOIWorkflow_WithEnrichment(t *testing.T) {
	// Create test file
	testFile := filepath.Join(t.TempDir(), "sample.fastq")
	fastqContent := `@SEQ
ACGTACGT
+
IIIIIIII
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

	// Create workflow
	config := &doi.WorkflowConfig{
		Publisher: "Generic Lab",
		License:   "CC0",
	}

	providerRegistry := doi.NewProviderRegistry()
	providerRegistry.Register(doi.NewDisabledProvider())
	_ = providerRegistry.SetActive("disabled")

	workflow := doi.NewDOIWorkflow(config, providerRegistry)

	// Prepare with enrichment
	enrichment := map[string]interface{}{
		"title":       "Enhanced Title",
		"description": "This is a comprehensive description of the dataset",
		"authors": []interface{}{
			map[string]interface{}{
				"name":        "Dr. Jane Doe",
				"given_name":  "Jane",
				"family_name": "Doe",
				"orcid":       "0000-0001-2345-6789",
				"affiliation": "University Lab",
			},
		},
		"keywords": []string{"genomics", "research"},
	}

	prepReq := &doi.PrepareRequest{
		FilePath:   testFile,
		Metadata:   extracted,
		Enrichment: enrichment,
	}

	result, err := workflow.Prepare(prepReq)
	if err != nil {
		t.Fatalf("Prepare() failed: %v", err)
	}

	// Verify enrichment was applied
	if result.Dataset.Title != "Enhanced Title" {
		t.Errorf("Title = %v, want Enhanced Title", result.Dataset.Title)
	}

	if result.Dataset.Description != "This is a comprehensive description of the dataset" {
		t.Errorf("Description not enriched correctly")
	}

	if len(result.Dataset.Authors) != 1 {
		t.Errorf("Authors count = %d, want 1", len(result.Dataset.Authors))
	} else {
		if result.Dataset.Authors[0].Name != "Dr. Jane Doe" {
			t.Errorf("Author name = %v, want Dr. Jane Doe", result.Dataset.Authors[0].Name)
		}
		if result.Dataset.Authors[0].ORCID != "0000-0001-2345-6789" {
			t.Errorf("Author ORCID not enriched")
		}
	}

	// Verify quality score improved with enrichment
	if result.Validation.Score < 70 {
		t.Errorf("Enriched dataset score = %.1f, expected > 70", result.Validation.Score)
	}
}

// TestDOIWorkflow_ValidationStrictness tests different validation configurations
func TestDOIWorkflow_ValidationStrictness(t *testing.T) {
	// Create test file
	testFile := filepath.Join(t.TempDir(), "sample.fastq")
	fastqContent := `@SEQ
ACGT
+
IIII
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

	tests := []struct {
		name               string
		requireRealAuthors bool
		requireDescription bool
		minQualityScore    float64
		wantReady          bool
	}{
		{
			name:               "Lenient validation",
			requireRealAuthors: false,
			requireDescription: false,
			minQualityScore:    40.0,
			wantReady:          true,
		},
		{
			name:               "Strict validation - require real authors",
			requireRealAuthors: true,
			requireDescription: false,
			minQualityScore:    40.0,
			wantReady:          false, // Will fail due to "Unknown Creator"
		},
		{
			name:               "Strict validation - high quality score",
			requireRealAuthors: false,
			requireDescription: false,
			minQualityScore:    80.0,
			wantReady:          false, // Will fail due to low quality
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &doi.WorkflowConfig{
				Publisher:          "Test Lab",
				RequireRealAuthors: tt.requireRealAuthors,
				RequireDescription: tt.requireDescription,
				MinQualityScore:    tt.minQualityScore,
			}

			providerRegistry := doi.NewProviderRegistry()
			providerRegistry.Register(doi.NewDisabledProvider())
			_ = providerRegistry.SetActive("disabled")

			workflow := doi.NewDOIWorkflow(config, providerRegistry)

			prepReq := &doi.PrepareRequest{
				FilePath: testFile,
				Metadata: extracted,
			}

			result, err := workflow.Prepare(prepReq)
			if err != nil {
				t.Fatalf("Prepare() failed: %v", err)
			}

			if result.Validation.IsReady != tt.wantReady {
				t.Errorf("IsReady = %v, want %v (score: %.1f, errors: %d)",
					result.Validation.IsReady, tt.wantReady,
					result.Validation.Score, len(result.Validation.Errors))
			}
		})
	}
}

// TestDOIWorkflow_MultipleFileTypes tests workflow with different file formats
func TestDOIWorkflow_MultipleFileTypes(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		content      string
		wantFormat   string
		minAuthors   int
		minKeywords  int
	}{
		{
			name:     "FASTQ file",
			filename: "sample_R1.fastq",
			content: `@SEQ
ACGTACGT
+
IIIIIIII
`,
			wantFormat:  "FASTQ",
			minAuthors:  1,
			minKeywords: 3, // genomics, sequencing, etc.
		},
		{
			name:     "FASTQ paired-end R2",
			filename: "sample_R2.fastq",
			content: `@SEQ
ACGTACGT
+
IIIIIIII
`,
			wantFormat:  "FASTQ",
			minAuthors:  1,
			minKeywords: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(t.TempDir(), tt.filename)
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Extract and prepare
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			extracted, err := registry.Extract(testFile)
			if err != nil {
				t.Fatalf("Extract() failed: %v", err)
			}

			// Verify format was detected correctly
			if extracted["format"] != tt.wantFormat {
				t.Errorf("Extracted format = %v, want %v", extracted["format"], tt.wantFormat)
			}

			// Create workflow and prepare
			config := &doi.WorkflowConfig{
				Publisher:          "Test Lab",
				RequireRealAuthors: false,
			}

			providerRegistry := doi.NewProviderRegistry()
			providerRegistry.Register(doi.NewDisabledProvider())
			_ = providerRegistry.SetActive("disabled")

			workflow := doi.NewDOIWorkflow(config, providerRegistry)

			prepReq := &doi.PrepareRequest{
				FilePath: testFile,
				Metadata: extracted,
			}

			result, err := workflow.Prepare(prepReq)
			if err != nil {
				t.Fatalf("Prepare() failed: %v", err)
			}

			// Verify mapping created appropriate fields
			if len(result.Dataset.Authors) < tt.minAuthors {
				t.Errorf("Authors count = %d, want >= %d", len(result.Dataset.Authors), tt.minAuthors)
			}

			if len(result.Dataset.Keywords) < tt.minKeywords {
				t.Errorf("Keywords count = %d, want >= %d", len(result.Dataset.Keywords), tt.minKeywords)
			}
		})
	}
}

// TestDOIWorkflow_ValidateMetadata tests the validation-only workflow
func TestDOIWorkflow_ValidateMetadata(t *testing.T) {
	// Create test file
	testFile := filepath.Join(t.TempDir(), "sample.fastq")
	fastqContent := `@SEQ
ACGTACGTACGTACGT
+
IIIIIIIIIIIIIIII
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

	// Create workflow
	config := &doi.WorkflowConfig{
		Publisher:       "Test Lab",
		MinQualityScore: 50.0,
	}

	providerRegistry := doi.NewProviderRegistry()
	providerRegistry.Register(doi.NewDisabledProvider())
	_ = providerRegistry.SetActive("disabled")

	workflow := doi.NewDOIWorkflow(config, providerRegistry)

	// Validate only (no preparation)
	result, err := workflow.ValidateMetadata(extracted, testFile)
	if err != nil {
		t.Fatalf("ValidateMetadata() failed: %v", err)
	}

	// Verify validation results
	if result.Validation == nil {
		t.Fatal("Validation should not be nil")
	}

	if result.Validation.Score < 0 || result.Validation.Score > 100 {
		t.Errorf("Validation score out of range: %.1f", result.Validation.Score)
	}

	// Verify present/missing fields are tracked
	if len(result.Validation.Present) == 0 && len(result.Validation.Missing) == 0 {
		t.Error("Validation should track present and missing fields")
	}

	// Verify recommendations are provided
	recommendations := workflow.GetRecommendations(result.Validation)
	if len(recommendations) == 0 {
		t.Error("GetRecommendations() should return suggestions")
	}
}

// TestDOIWorkflow_QualityScoring tests quality score calculation
func TestDOIWorkflow_QualityScoring(t *testing.T) {
	// Create test file
	testFile := filepath.Join(t.TempDir(), "sample.fastq")
	fastqContent := `@SEQ
ACGT
+
IIII
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

	// Test with minimal enrichment
	minimalEnrichment := map[string]interface{}{
		"authors": []interface{}{
			map[string]interface{}{
				"name":        "Author Name",
				"given_name":  "Author",
				"family_name": "Name",
			},
		},
	}

	// Test with rich enrichment
	richEnrichment := map[string]interface{}{
		"title":       "Comprehensive Dataset Title",
		"description": "This is a detailed description of the dataset providing context and methodology for reproducibility and reuse.",
		"authors": []interface{}{
			map[string]interface{}{
				"name":        "Dr. Jane Doe",
				"given_name":  "Jane",
				"family_name": "Doe",
				"orcid":       "0000-0001-2345-6789",
				"affiliation": "University Research Lab",
			},
		},
		"keywords": []string{"keyword1", "keyword2", "keyword3", "keyword4", "keyword5"},
		"license":  "CC-BY-4.0",
		"url":      "https://example.com/dataset",
		"version":  "1.0",
	}

	config := &doi.WorkflowConfig{
		Publisher: "Test Lab",
	}

	providerRegistry := doi.NewProviderRegistry()
	providerRegistry.Register(doi.NewDisabledProvider())
	_ = providerRegistry.SetActive("disabled")

	workflow := doi.NewDOIWorkflow(config, providerRegistry)

	// Test minimal enrichment
	minimalReq := &doi.PrepareRequest{
		FilePath:   testFile,
		Metadata:   extracted,
		Enrichment: minimalEnrichment,
	}

	minimalResult, err := workflow.Prepare(minimalReq)
	if err != nil {
		t.Fatalf("Prepare() with minimal enrichment failed: %v", err)
	}

	// Test rich enrichment
	richReq := &doi.PrepareRequest{
		FilePath:   testFile,
		Metadata:   extracted,
		Enrichment: richEnrichment,
	}

	richResult, err := workflow.Prepare(richReq)
	if err != nil {
		t.Fatalf("Prepare() with rich enrichment failed: %v", err)
	}

	// Rich enrichment should have higher score
	// Note: Minimal with proper author already scores ~84, so improvement is ~10 points
	scoreDiff := richResult.Validation.Score - minimalResult.Validation.Score
	if scoreDiff < 5 {
		t.Errorf("Rich enrichment score improvement = %.1f, expected >= 5 points\nMinimal: %.1f, Rich: %.1f",
			scoreDiff, minimalResult.Validation.Score, richResult.Validation.Score)
	}

	// Verify absolute scores are in expected ranges
	if minimalResult.Validation.Score < 75 {
		t.Errorf("Minimal enrichment score = %.1f, expected >= 75", minimalResult.Validation.Score)
	}

	if richResult.Validation.Score < 85 {
		t.Errorf("Rich enrichment score = %.1f, expected >= 85", richResult.Validation.Score)
	}

	// Rich enrichment should be ready
	if !richResult.Validation.IsReady {
		t.Errorf("Richly enriched dataset should be ready (score: %.1f, errors: %v)",
			richResult.Validation.Score, richResult.Validation.Errors)
	}
}
