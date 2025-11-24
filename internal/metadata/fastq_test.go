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

package metadata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestFASTQExtractor_Name(t *testing.T) {
	extractor := &FASTQExtractor{}
	if got := extractor.Name(); got != "FASTQ" {
		t.Errorf("Name() = %v, want %v", got, "FASTQ")
	}
}

func TestFASTQExtractor_SupportedFormats(t *testing.T) {
	extractor := &FASTQExtractor{}
	formats := extractor.SupportedFormats()
	expected := []string{".fastq", ".fq", ".fastq.gz", ".fq.gz"}

	if len(formats) != len(expected) {
		t.Errorf("SupportedFormats() returned %d formats, want %d", len(formats), len(expected))
	}

	for i, format := range expected {
		if formats[i] != format {
			t.Errorf("SupportedFormats()[%d] = %v, want %v", i, formats[i], format)
		}
	}
}

func TestFASTQExtractor_CanHandle(t *testing.T) {
	extractor := &FASTQExtractor{}

	tests := []struct {
		filename string
		want     bool
	}{
		{"sample.fastq", true},
		{"sample.fq", true},
		{"sample.fastq.gz", true},
		{"sample.fq.gz", true},
		{"SAMPLE.FASTQ", true}, // Case insensitive
		{"sample_R1.fastq", true},
		{"sample_R2.fq.gz", true},
		{"sample.txt", false},
		{"sample.fasta", false},
		{"sample.bam", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if got := extractor.CanHandle(tt.filename); got != tt.want {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestFASTQExtractor_ExtractFromReader(t *testing.T) {
	extractor := &FASTQExtractor{}

	// Create sample FASTQ data with 3 reads
	fastqData := `@SEQ_ID_1
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

	reader := strings.NewReader(fastqData)
	metadata, err := extractor.ExtractFromReader(reader, "test.fastq")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	// Verify basic fields
	if metadata["format"] != "FASTQ" {
		t.Errorf("format = %v, want FASTQ", metadata["format"])
	}

	if metadata["extractor_name"] != "fastq" {
		t.Errorf("extractor_name = %v, want fastq", metadata["extractor_name"])
	}

	if metadata["schema_name"] != "fastq_v1" {
		t.Errorf("schema_name = %v, want fastq_v1", metadata["schema_name"])
	}

	// Verify read counts
	if metadata["total_reads"] != 3 {
		t.Errorf("total_reads = %v, want 3", metadata["total_reads"])
	}

	// Verify total bases (60 + 24 + 14 = 98)
	if metadata["total_bases"] != int64(98) {
		t.Errorf("total_bases = %v, want 98", metadata["total_bases"])
	}

	// Verify mean read length
	meanLength, ok := metadata["mean_read_length"].(float64)
	if !ok {
		t.Errorf("mean_read_length is not float64")
	}
	expectedMean := 98.0 / 3.0
	if meanLength != expectedMean {
		t.Errorf("mean_read_length = %v, want %v", meanLength, expectedMean)
	}

	// Verify min/max read length
	if metadata["min_read_length"] != 14 {
		t.Errorf("min_read_length = %v, want 14", metadata["min_read_length"])
	}

	if metadata["max_read_length"] != 60 {
		t.Errorf("max_read_length = %v, want 60", metadata["max_read_length"])
	}

	// Verify GC content exists
	if _, ok := metadata["gc_content_percent"]; !ok {
		t.Errorf("gc_content_percent not found in metadata")
	}

	// Verify quality scores exist
	if _, ok := metadata["mean_quality_score"]; !ok {
		t.Errorf("mean_quality_score not found in metadata")
	}

	// Verify instrument type
	if metadata["instrument_type"] != "sequencing" {
		t.Errorf("instrument_type = %v, want sequencing", metadata["instrument_type"])
	}

	if metadata["data_type"] != "nucleotide_sequence" {
		t.Errorf("data_type = %v, want nucleotide_sequence", metadata["data_type"])
	}

	// Verify compression
	if metadata["compression"] != "none" {
		t.Errorf("compression = %v, want none", metadata["compression"])
	}
}

func TestFASTQExtractor_GCContent(t *testing.T) {
	extractor := &FASTQExtractor{}

	// Create FASTQ with known GC content (50% - 4 G, 4 C out of 8 bases)
	fastqData := `@SEQ_ID
GCGCATAA
+
IIIIIIII
`

	reader := strings.NewReader(fastqData)
	metadata, err := extractor.ExtractFromReader(reader, "test.fastq")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	gcContent, ok := metadata["gc_content_percent"].(float64)
	if !ok {
		t.Fatalf("gc_content_percent is not float64")
	}

	// 4 GC bases out of 8 total = 50%
	expected := 50.0
	if gcContent != expected {
		t.Errorf("gc_content_percent = %v, want %v", gcContent, expected)
	}
}

func TestFASTQExtractor_QualityScores(t *testing.T) {
	extractor := &FASTQExtractor{}

	// Create FASTQ with known quality scores
	// Using 'I' (ASCII 73) = Phred score 40 (73 - 33 = 40)
	fastqData := `@SEQ_ID
ACGT
+
IIII
`

	reader := strings.NewReader(fastqData)
	metadata, err := extractor.ExtractFromReader(reader, "test.fastq")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	meanQuality, ok := metadata["mean_quality_score"].(float64)
	if !ok {
		t.Fatalf("mean_quality_score is not float64")
	}

	// All bases have quality 40
	expected := 40.0
	if meanQuality != expected {
		t.Errorf("mean_quality_score = %v, want %v", meanQuality, expected)
	}

	if metadata["min_quality_score"] != 40 {
		t.Errorf("min_quality_score = %v, want 40", metadata["min_quality_score"])
	}

	if metadata["max_quality_score"] != 40 {
		t.Errorf("max_quality_score = %v, want 40", metadata["max_quality_score"])
	}
}

func TestFASTQExtractor_PairedEnd(t *testing.T) {
	extractor := &FASTQExtractor{}

	tests := []struct {
		filename string
		isPaired bool
		readPair string
		gzipped  bool
	}{
		{"sample_R1.fastq", true, "R1", false},
		{"sample_R2.fastq", true, "R2", false},
		{"sample_R1_001.fastq.gz", true, "R1", true},
		{"sample_R2_001.fastq.gz", true, "R2", true},
		{"sample_1.fastq", true, "1", false},
		{"sample_2.fastq", true, "2", false},
		{"sample.1.fastq", true, "1", false},
		{"sample.2.fastq", true, "2", false},
		{"sample.fastq", false, "", false},
		{"single_end.fq", false, "", false},
	}

	// Create minimal valid FASTQ data
	fastqData := `@SEQ_ID
ACGT
+
IIII
`

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			var reader io.Reader
			if tt.gzipped {
				// Gzip the data for .gz tests
				var buf bytes.Buffer
				gzWriter := gzip.NewWriter(&buf)
				if _, err := gzWriter.Write([]byte(fastqData)); err != nil {
					t.Fatalf("Failed to write gzip data: %v", err)
				}
				if err := gzWriter.Close(); err != nil {
					t.Fatalf("Failed to close gzip writer: %v", err)
				}
				reader = bytes.NewReader(buf.Bytes())
			} else {
				reader = strings.NewReader(fastqData)
			}

			metadata, err := extractor.ExtractFromReader(reader, tt.filename)
			if err != nil {
				t.Fatalf("ExtractFromReader() error = %v", err)
			}

			isPaired, ok := metadata["is_paired_end"].(bool)
			if !ok {
				t.Fatalf("is_paired_end is not bool")
			}

			if isPaired != tt.isPaired {
				t.Errorf("is_paired_end = %v, want %v", isPaired, tt.isPaired)
			}

			if tt.isPaired {
				readPair, ok := metadata["read_pair"].(string)
				if !ok {
					t.Fatalf("read_pair is not string")
				}
				if readPair != tt.readPair {
					t.Errorf("read_pair = %v, want %v", readPair, tt.readPair)
				}
			} else {
				if _, exists := metadata["read_pair"]; exists {
					t.Errorf("read_pair should not exist for single-end file")
				}
			}
		})
	}
}

func TestFASTQExtractor_Gzip(t *testing.T) {
	extractor := &FASTQExtractor{}

	// Create sample FASTQ data
	fastqData := `@SEQ_ID
ACGTACGT
+
IIIIIIII
`

	// Compress data with gzip
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write([]byte(fastqData)); err != nil {
		t.Fatalf("Failed to write gzip data: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	metadata, err := extractor.ExtractFromReader(reader, "test.fastq.gz")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	// Verify compression detected
	if metadata["compression"] != "gzip" {
		t.Errorf("compression = %v, want gzip", metadata["compression"])
	}

	// Verify data was parsed correctly
	if metadata["total_reads"] != 1 {
		t.Errorf("total_reads = %v, want 1", metadata["total_reads"])
	}

	if metadata["total_bases"] != int64(8) {
		t.Errorf("total_bases = %v, want 8", metadata["total_bases"])
	}
}

func TestFASTQExtractor_InvalidFormat(t *testing.T) {
	extractor := &FASTQExtractor{}

	tests := []struct {
		name string
		data string
	}{
		{
			name: "missing @ in identifier",
			data: `SEQ_ID
ACGT
+
IIII
`,
		},
		{
			name: "missing + in separator",
			data: `@SEQ_ID
ACGT
SEQ_ID
IIII
`,
		},
		{
			name: "empty file",
			data: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.data)
			_, err := extractor.ExtractFromReader(reader, "test.fastq")
			if err == nil {
				t.Errorf("ExtractFromReader() should return error for invalid format")
			}
		})
	}
}

func TestFASTQExtractor_LargeFile(t *testing.T) {
	extractor := &FASTQExtractor{}

	// Create FASTQ with many reads (should sample only first 10,000)
	var buf strings.Builder
	for i := 0; i < 20000; i++ {
		buf.WriteString("@SEQ_ID_")
		buf.WriteString(fmt.Sprintf("%d", i))
		buf.WriteString("\n")
		buf.WriteString("ACGTACGTACGTACGT\n")
		buf.WriteString("+\n")
		buf.WriteString("IIIIIIIIIIIIIIII\n")
	}

	metadata, err := extractor.ExtractFromReader(strings.NewReader(buf.String()), "large.fastq")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	// Should sample only first 10,000 reads
	totalReads := metadata["total_reads"].(int)
	if totalReads != 10000 {
		t.Errorf("total_reads = %v, want 10000 (sampled)", totalReads)
	}

	// Total bases should be 10,000 * 16 = 160,000
	totalBases := metadata["total_bases"].(int64)
	if totalBases != 160000 {
		t.Errorf("total_bases = %v, want 160000", totalBases)
	}
}

func TestDetectPairedEnd(t *testing.T) {
	tests := []struct {
		filename   string
		wantPaired bool
		wantPair   string
	}{
		{"sample_R1.fastq", true, "R1"},
		{"sample_R2.fastq", true, "R2"},
		{"sample_R1_001.fastq", true, "R1"},
		{"sample_R2_001.fastq", true, "R2"},
		{"sample_1.fastq", true, "1"},
		{"sample_2.fastq", true, "2"},
		{"sample.1.fastq", true, "1"},
		{"sample.2.fastq", true, "2"},
		{"sample.fastq", false, ""},
		{"single.fq", false, ""},
		{"Sample_R1.FASTQ", true, "R1"}, // Case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := detectPairedEnd(tt.filename)

			isPaired, ok := result["is_paired"].(bool)
			if !ok {
				t.Fatalf("is_paired is not bool")
			}

			if isPaired != tt.wantPaired {
				t.Errorf("is_paired = %v, want %v", isPaired, tt.wantPaired)
			}

			if tt.wantPaired {
				readPair, ok := result["read_pair"].(string)
				if !ok {
					t.Fatalf("read_pair is not string")
				}
				if readPair != tt.wantPair {
					t.Errorf("read_pair = %v, want %v", readPair, tt.wantPair)
				}
			}
		})
	}
}
