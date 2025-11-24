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
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/scttfrdmn/cicada/internal/metadata"
)

// TestMetadataExtraction_FASTQ tests end-to-end FASTQ metadata extraction
func TestMetadataExtraction_FASTQ(t *testing.T) {
	// Create a real FASTQ file with actual data
	testFile := filepath.Join(t.TempDir(), "test_sample_R1.fastq")

	fastqContent := `@SEQ_ID_1 test sequence 1
GATTTGGGGTTCAAAGCAGTATCGATCAAATAGTAAATCCATTTGTTCAACTCACAGTTT
+
!''*((((***+))%%%++)(%%%%).1***-+*''))**55CCF>>>>>>CCCCCCC65
@SEQ_ID_2 test sequence 2
GCGCGCGCGCGCGCGCGCGCGCGC
+
IIIIIIIIIIIIIIIIIIIIIIIII
@SEQ_ID_3 test sequence 3
ATATATATATATAT
+
HHHHHHHHHHHHHH
`

	if err := os.WriteFile(testFile, []byte(fastqContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create extractor registry and extract metadata
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	extracted, err := registry.Extract(testFile)
	if err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Verify extracted metadata
	if extracted["format"] != "FASTQ" {
		t.Errorf("format = %v, want FASTQ", extracted["format"])
	}

	if extracted["total_reads"] != 3 {
		t.Errorf("total_reads = %v, want 3", extracted["total_reads"])
	}

	if extracted["total_bases"].(int64) != 98 {
		t.Errorf("total_bases = %v, want 98", extracted["total_bases"])
	}

	// Verify GC content calculation
	if gcContent, ok := extracted["gc_content_percent"].(float64); !ok {
		t.Error("gc_content_percent should be present")
	} else if gcContent < 0 || gcContent > 100 {
		t.Errorf("gc_content_percent = %v, should be between 0 and 100", gcContent)
	}

	// Verify quality scores
	if meanQuality, ok := extracted["mean_quality_score"].(float64); !ok {
		t.Error("mean_quality_score should be present")
	} else if meanQuality < 0 || meanQuality > 93 { // Phred+33 max is 93
		t.Errorf("mean_quality_score = %v, should be between 0 and 93", meanQuality)
	}

	// Verify paired-end detection
	if isPaired, ok := extracted["is_paired_end"].(bool); !ok {
		t.Error("is_paired_end should be present")
	} else if !isPaired {
		t.Error("is_paired_end should be true for _R1 file")
	}

	if readPair, ok := extracted["read_pair"].(string); !ok {
		t.Error("read_pair should be present")
	} else if readPair != "R1" {
		t.Errorf("read_pair = %v, want R1", readPair)
	}
}

// TestMetadataExtraction_FASTQGzipped tests gzipped FASTQ extraction
func TestMetadataExtraction_FASTQGzipped(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test_sample.fastq.gz")

	// Create gzipped FASTQ content
	fastqContent := `@SEQ_ID
ACGTACGT
+
IIIIIIII
`

	gzContent := gzipContent(t, []byte(fastqContent))
	if err := os.WriteFile(testFile, gzContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Extract metadata
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	extracted, err := registry.Extract(testFile)
	if err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Verify decompression worked
	if extracted["compression"] != "gzip" {
		t.Errorf("compression = %v, want gzip", extracted["compression"])
	}

	if extracted["total_reads"] != 1 {
		t.Errorf("total_reads = %v, want 1", extracted["total_reads"])
	}
}

// TestMetadataExtraction_AutoDetection tests format auto-detection
func TestMetadataExtraction_AutoDetection(t *testing.T) {
	tests := []struct {
		name              string
		filename          string
		content           string
		wantFormat        string
		wantExtractorName string
	}{
		{
			name:     "FASTQ with .fastq extension",
			filename: "sample.fastq",
			content: `@SEQ
ACGT
+
IIII
`,
			wantFormat:        "FASTQ",
			wantExtractorName: "FASTQ",
		},
		{
			name:     "FASTQ with .fq extension",
			filename: "sample.fq",
			content: `@SEQ
ACGT
+
IIII
`,
			wantFormat:        "FASTQ",
			wantExtractorName: "FASTQ",
		},
		{
			name:     "FASTQ with .fastq.gz extension",
			filename: "sample.fastq.gz",
			content: `@SEQ
ACGT
+
IIII
`,
			wantFormat:        "FASTQ",
			wantExtractorName: "FASTQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(t.TempDir(), tt.filename)

			content := []byte(tt.content)
			if filepath.Ext(tt.filename) == ".gz" {
				content = gzipContent(t, content)
			}

			if err := os.WriteFile(testFile, content, 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test auto-detection
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			extractor := registry.FindExtractor(tt.filename)
			if extractor == nil {
				t.Fatalf("FindExtractor() returned nil for %s", tt.filename)
			}

			if extractor.Name() != tt.wantExtractorName {
				t.Errorf("FindExtractor().Name() = %v, want %v", extractor.Name(), tt.wantExtractorName)
			}

			// Test extraction
			extracted, err := registry.Extract(testFile)
			if err != nil {
				t.Fatalf("Extract() failed: %v", err)
			}

			if extracted["format"] != tt.wantFormat {
				t.Errorf("format = %v, want %v", extracted["format"], tt.wantFormat)
			}

			// Note: extractor_name in metadata is lowercase, Name() is uppercase
			extractorName, ok := extracted["extractor_name"].(string)
			if !ok {
				t.Error("extractor_name field not found in metadata")
			} else if extractorName != "fastq" {
				t.Errorf("extractor_name = %v, want fastq", extractorName)
			}
		})
	}
}

// TestMetadataExtraction_InvalidFiles tests error handling for invalid files
func TestMetadataExtraction_InvalidFiles(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		wantErr  bool
	}{
		{
			name:     "Empty FASTQ file",
			filename: "empty.fastq",
			content:  "",
			wantErr:  true,
		},
		{
			name:     "Invalid FASTQ format - missing @",
			filename: "invalid.fastq",
			content: `SEQ_ID
ACGT
+
IIII
`,
			wantErr: true,
		},
		{
			name:     "Invalid FASTQ format - missing +",
			filename: "invalid2.fastq",
			content: `@SEQ_ID
ACGT
SEQ_ID
IIII
`,
			wantErr: true,
		},
		{
			name:     "Non-existent file",
			filename: "does-not-exist.fastq",
			content:  "", // Won't be created
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := metadata.NewExtractorRegistry()
			registry.RegisterDefaults()

			testFile := filepath.Join(t.TempDir(), tt.filename)

			// Only create file if content is not empty
			if tt.content != "" {
				if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			_, err := registry.Extract(testFile)

			if tt.wantErr && err == nil {
				t.Error("Extract() should have returned error but didn't")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Extract() unexpected error: %v", err)
			}
		})
	}
}

// TestMetadataExtraction_LargeFile tests performance with large files
func TestMetadataExtraction_LargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file test in short mode")
	}

	testFile := filepath.Join(t.TempDir(), "large.fastq")

	// Create a file with 50,000 reads (should sample 10,000)
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer f.Close()

	for i := 0; i < 50000; i++ {
		f.WriteString("@SEQ_ID_")
		f.WriteString(string(rune(48 + (i % 10)))) // Write digits 0-9
		f.WriteString("\n")
		f.WriteString("ACGTACGTACGTACGT\n")
		f.WriteString("+\n")
		f.WriteString("IIIIIIIIIIIIIIII\n")
	}
	f.Close()

	// Extract metadata
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	extracted, err := registry.Extract(testFile)
	if err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Verify sampling (should stop at 10,000)
	totalReads, ok := extracted["total_reads"].(int)
	if !ok {
		t.Fatal("total_reads not present or wrong type")
	}

	if totalReads != 10000 {
		t.Errorf("total_reads = %v, want 10000 (sampled)", totalReads)
	}

	// Verify file was recognized
	if extracted["format"] != "FASTQ" {
		t.Errorf("format = %v, want FASTQ", extracted["format"])
	}
}

// TestMetadataExtraction_ConcurrentExtractions tests thread safety
func TestMetadataExtraction_ConcurrentExtractions(t *testing.T) {
	// Create multiple test files
	tempDir := t.TempDir()

	files := []string{}
	for i := 0; i < 10; i++ {
		filename := filepath.Join(tempDir, "test_"+string(rune(48+i))+".fastq")
		content := `@SEQ_ID
ACGTACGT
+
IIIIIIII
`
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		files = append(files, filename)
	}

	// Extract concurrently
	registry := metadata.NewExtractorRegistry()
	registry.RegisterDefaults()

	results := make(chan error, len(files))

	for _, file := range files {
		go func(f string) {
			_, err := registry.Extract(f)
			results <- err
		}(file)
	}

	// Collect results
	for i := 0; i < len(files); i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent extraction failed: %v", err)
		}
	}
}

// gzipContent compresses content with gzip
func gzipContent(t *testing.T, content []byte) []byte {
	t.Helper()

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(content); err != nil {
		t.Fatalf("Failed to write gzip data: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}

	return buf.Bytes()
}
