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
	"compress/gzip"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/scttfrdmn/cicada/internal/doi"
	"github.com/scttfrdmn/cicada/internal/metadata"
)

// BenchmarkMetadataExtraction_SmallFASTQ benchmarks extraction of small FASTQ file (10 reads)
func BenchmarkMetadataExtraction_SmallFASTQ(b *testing.B) {
	// Create small FASTQ file (10 reads)
	testFile := filepath.Join(b.TempDir(), "small.fastq")
	fastqContent := ""
	for i := 0; i < 10; i++ {
		fastqContent += "@SEQ_ID_" + string(rune(48+i%10)) + "\n"
		fastqContent += "ACGTACGTACGTACGT\n"
		fastqContent += "+\n"
		fastqContent += "IIIIIIIIIIIIIIII\n"
	}
	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	// Setup
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Extract(testFile)
		if err != nil {
			b.Fatalf("Extract() failed: %v", err)
		}
	}
}

// BenchmarkMetadataExtraction_MediumFASTQ benchmarks extraction of medium FASTQ file (1000 reads)
func BenchmarkMetadataExtraction_MediumFASTQ(b *testing.B) {
	// Create medium FASTQ file (1000 reads)
	testFile := filepath.Join(b.TempDir(), "medium.fastq")
	f, err := os.Create(testFile)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}
	for i := 0; i < 1000; i++ {
		if _, err := f.WriteString("@SEQ_ID_" + string(rune(48+i%10)) + "\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
		if _, err := f.WriteString("ACGTACGTACGTACGT\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
		if _, err := f.WriteString("+\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
		if _, err := f.WriteString("IIIIIIIIIIIIIIII\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
	}
	if err := f.Close(); err != nil {
		b.Fatalf("Failed to close test file: %v", err)
	}

	// Setup
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Extract(testFile)
		if err != nil {
			b.Fatalf("Extract() failed: %v", err)
		}
	}
}

// BenchmarkMetadataExtraction_LargeFASTQ benchmarks extraction of large FASTQ file (10,000 reads)
func BenchmarkMetadataExtraction_LargeFASTQ(b *testing.B) {
	// Create large FASTQ file (10,000 reads - will be sampled)
	testFile := filepath.Join(b.TempDir(), "large.fastq")
	f, err := os.Create(testFile)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}
	for i := 0; i < 10000; i++ {
		if _, err := f.WriteString("@SEQ_ID_" + string(rune(48+i%10)) + "\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
		if _, err := f.WriteString("ACGTACGTACGTACGT\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
		if _, err := f.WriteString("+\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
		if _, err := f.WriteString("IIIIIIIIIIIIIIII\n"); err != nil {
			b.Fatalf("Failed to write to test file: %v", err)
		}
	}
	if err := f.Close(); err != nil {
		b.Fatalf("Failed to close test file: %v", err)
	}

	// Setup
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Extract(testFile)
		if err != nil {
			b.Fatalf("Extract() failed: %v", err)
		}
	}
}

// BenchmarkMetadataExtraction_GzipFASTQ benchmarks extraction of gzipped FASTQ file
func BenchmarkMetadataExtraction_GzipFASTQ(b *testing.B) {
	// Create gzipped FASTQ file (1000 reads)
	testFile := filepath.Join(b.TempDir(), "compressed.fastq.gz")
	f, err := os.Create(testFile)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}
	gzWriter := gzip.NewWriter(f)
	for i := 0; i < 1000; i++ {
		if _, err := gzWriter.Write([]byte("@SEQ_ID_" + string(rune(48+i%10)) + "\n")); err != nil {
			b.Fatalf("Failed to write to gzip writer: %v", err)
		}
		if _, err := gzWriter.Write([]byte("ACGTACGTACGTACGT\n")); err != nil {
			b.Fatalf("Failed to write to gzip writer: %v", err)
		}
		if _, err := gzWriter.Write([]byte("+\n")); err != nil {
			b.Fatalf("Failed to write to gzip writer: %v", err)
		}
		if _, err := gzWriter.Write([]byte("IIIIIIIIIIIIIIII\n")); err != nil {
			b.Fatalf("Failed to write to gzip writer: %v", err)
		}
	}
	if err := gzWriter.Close(); err != nil {
		b.Fatalf("Failed to close gzip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		b.Fatalf("Failed to close test file: %v", err)
	}

	// Setup
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Extract(testFile)
		if err != nil {
			b.Fatalf("Extract() failed: %v", err)
		}
	}
}

// BenchmarkMetadataExtraction_Concurrent benchmarks concurrent extraction
func BenchmarkMetadataExtraction_Concurrent(b *testing.B) {
	// Create 10 test files
	tempDir := b.TempDir()
	files := make([]string, 10)
	for i := 0; i < 10; i++ {
		files[i] = filepath.Join(tempDir, "sample"+string(rune(48+i))+".fastq")
		fastqContent := "@SEQ_ID\nACGTACGTACGTACGT\n+\nIIIIIIIIIIIIIIII\n"
		if err := os.WriteFile(files[i], []byte(fastqContent), 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Setup
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for _, file := range files {
			wg.Add(1)
			go func(f string) {
				defer wg.Done()
				_, err := registry.Extract(f)
				if err != nil {
					b.Errorf("Extract() failed: %v", err)
				}
			}(file)
		}
		wg.Wait()
	}
}

// BenchmarkPresetValidation benchmarks preset validation
func BenchmarkPresetValidation(b *testing.B) {
	// Setup test metadata
	testMetadata := map[string]interface{}{
		"format":                   "FASTQ",
		"instrument_manufacturer":  "Illumina",
		"instrument_model":         "NovaSeq 6000",
		"sequencing_platform":      "Illumina",
		"total_reads":              1000000,
		"read_length":              150,
		"run_id":                   "RUN_123",
		"flowcell_id":              "FLOWCELL_456",
		"lane":                     2,
		"barcode":                  "AGTCACTA",
		"quality_encoding":         "Phred+33",
		"gc_content":               52.3,
		"mean_quality_score":       36.8,
		"is_paired_end":            true,
	}

	// Setup preset registry
	presetRegistry := metadata.NewPresetRegistry()
	presetRegistry.RegisterDefaults()

	// Get preset
	preset, err := presetRegistry.GetPreset("illumina-novaseq")
	if err != nil {
		b.Fatalf("GetPreset() failed: %v", err)
	}

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = preset.Validate(testMetadata)
	}
}

// BenchmarkDOIWorkflow_EndToEnd benchmarks complete DOI workflow
func BenchmarkDOIWorkflow_EndToEnd(b *testing.B) {
	// Create test FASTQ file
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "research_data.fastq")
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
		b.Fatalf("Failed to create test file: %v", err)
	}

	// Setup enrichment
	enrichmentData := map[string]interface{}{
		"title": "Benchmark Dataset",
		"authors": []map[string]interface{}{
			{
				"name":        "Test Researcher",
				"orcid":       "0000-0001-2345-6789",
				"affiliation": "Test Lab",
			},
		},
		"description": "A benchmark dataset for performance testing",
		"publisher":   "Test Lab",
	}

	// Setup DOI workflow
	config := &doi.WorkflowConfig{
		Publisher:          "Test Lab",
		License:            "CC-BY-4.0",
		MinQualityScore:    60.0,
		RequireRealAuthors: true,
		RequireDescription: true,
	}
	providerRegistry := doi.NewProviderRegistry()
	workflow := doi.NewDOIWorkflow(config, providerRegistry)

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Extract metadata
		registry := metadata.NewExtractorRegistry()
		registry.RegisterDefaults()
		extracted, err := registry.Extract(testFile)
		if err != nil {
			b.Fatalf("Extract() failed: %v", err)
		}

		// Prepare DOI (map + validate)
		prepReq := &doi.PrepareRequest{
			FilePath:   testFile,
			Metadata:   extracted,
			Enrichment: enrichmentData,
		}
		_, err = workflow.Prepare(prepReq)
		if err != nil {
			b.Fatalf("Prepare() failed: %v", err)
		}
	}
}

// BenchmarkDOIWorkflow_MetadataMapping benchmarks just metadata mapping
func BenchmarkDOIWorkflow_MetadataMapping(b *testing.B) {
	// Setup test data
	extracted := map[string]interface{}{
		"format":                   "FASTQ",
		"total_reads":              1000000,
		"total_bases":              150000000,
		"gc_content":               52.3,
		"mean_quality_score":       36.8,
		"read_length":              150,
		"is_paired_end":            true,
		"instrument_manufacturer":  "Illumina",
		"instrument_model":         "NovaSeq 6000",
		"sequencing_platform":      "Illumina",
	}

	// Setup mapper
	mapper := doi.NewMetadataMapper("Test Lab", "CC-BY-4.0", "https://example.com")

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mapper.MapToDataset(extracted, "test.fastq")
		if err != nil {
			b.Fatalf("MapToDataset() failed: %v", err)
		}
	}
}

// BenchmarkDOIWorkflow_Validation benchmarks DOI metadata validation
func BenchmarkDOIWorkflow_Validation(b *testing.B) {
	// Setup test dataset
	dataset := &doi.Dataset{
		Title:           "Test Dataset",
		Authors: []doi.Author{
			{
				Name:        "Jane Smith",
				Affiliation: "Test University",
			},
		},
		Publisher:       "Test University",
		PublicationYear: 2025,
		ResourceType:    "Dataset",
		Description:     "Test description for validation benchmark",
	}

	// Setup validator
	validator := doi.NewDOIReadinessValidator()
	validator.MinQualityScore = 60.0
	validator.RequireRealAuthors = true
	validator.RequireDescription = true

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.Validate(dataset)
	}
}

// BenchmarkPresetRegistry_FindPresets benchmarks finding presets
func BenchmarkPresetRegistry_FindPresets(b *testing.B) {
	// Setup
	presetRegistry := metadata.NewPresetRegistry()
	presetRegistry.RegisterDefaults()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = presetRegistry.FindPresets("Illumina", "sequencing")
	}
}

// BenchmarkPresetRegistry_ListPresets benchmarks listing all presets
func BenchmarkPresetRegistry_ListPresets(b *testing.B) {
	// Setup
	presetRegistry := metadata.NewPresetRegistry()
	presetRegistry.RegisterDefaults()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = presetRegistry.ListPresets()
	}
}
