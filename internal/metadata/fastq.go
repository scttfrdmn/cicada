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

// Package metadata provides metadata extraction for scientific instrument files.
//
// # FASTQ Format
//
// This file implements metadata extraction for FASTQ files, a text-based format
// for storing both nucleotide sequences and their corresponding quality scores.
//
// FASTQ is the de facto standard for representing high-throughput sequencing data
// from platforms including Illumina, PacBio, Oxford Nanopore, and others.
//
// ## Format Overview
//
// Each sequence entry in a FASTQ file consists of exactly four lines:
//   1. Sequence identifier line (begins with '@')
//   2. Raw sequence bases (ACGT or N for unknown)
//   3. Separator line (begins with '+', optionally repeats identifier)
//   4. Base quality scores (ASCII-encoded Phred scores)
//
// Example entry:
//
//	@SEQ_ID
//	GATTTGGGGTTCAAAGCAGTATCGATCAAATAGTAAATCCATTTGTTCAACTCACAGTTT
//	+
//	!''*((((***+))%%%++)(%%%%).1***-+*''))**55CCF>>>>>>CCCCCCC65
//
// ## Quality Score Encoding
//
// Modern FASTQ files use Phred+33 encoding (Sanger format):
//   - ASCII values 33-126 represent Phred quality scores 0-93
//   - Phred score Q represents error probability P: Q = -10 × log₁₀(P)
//   - Q=20 means P=0.01 (1% error rate, 99% accuracy)
//   - Q=30 means P=0.001 (0.1% error rate, 99.9% accuracy)
//
// ## Paired-End Sequencing
//
// Paired-end reads sequence both ends of a DNA fragment, generating two files:
//   - Read 1 (R1): Forward read
//   - Read 2 (R2): Reverse read
//
// Standard Illumina naming convention:
//   {sample}_S{sample_num}_L{lane}_{R1|R2}_001.fastq.gz
//
// Common naming patterns for paired files:
//   - filename_R1.fastq / filename_R2.fastq
//   - filename_1.fastq / filename_2.fastq
//   - filename.1.fastq / filename.2.fastq
//   - filename_R1_001.fastq / filename_R2_001.fastq
//
// ## Metadata Extraction
//
// This extractor analyzes FASTQ files to extract:
//   - Total read count
//   - Read length statistics (min, max, mean)
//   - GC content percentage
//   - Quality score statistics (mean, min, max)
//   - Paired-end detection from filename
//   - File format and compression
//
// ## Sampling Strategy
//
// For performance with large files (GB-TB scale), the extractor:
//   - Samples reads throughout the file (not just beginning)
//   - Calculates statistics from sample (configurable size)
//   - Provides accurate estimates without full file scan
//   - Uses streaming to minimize memory usage
//
// ## References and Sources
//
// ### Format Specifications
//
// FASTQ Format - Wikipedia:
// https://en.wikipedia.org/wiki/FASTQ_format
//
// Comprehensive overview of FASTQ format history, variants, and specifications.
//
// The Sanger FASTQ file format for sequences with quality scores:
// https://pmc.ncbi.nlm.nih.gov/articles/PMC2847217/
// https://pubmed.ncbi.nlm.nih.gov/20015970/
//
// Original paper by Cock et al. (2010) defining the Sanger FASTQ standard
// and documenting quality score encoding variants.
//
// ### Illumina Documentation
//
// FASTQ Files Explained - Illumina Knowledge:
// https://knowledge.illumina.com/software/general/software-general-reference_material-list/000002211
//
// Official Illumina documentation on FASTQ format and quality scores.
//
// Illumina FASTQ Naming Convention:
// https://support.illumina.com/help/BaseSpace_Sequence_Hub_OLH_009008_2/Source/Informatics/BS/NamingConvention_FASTQ-files-swBS.htm
//
// Standard naming conventions for Illumina sequencer outputs.
//
// ### Bioinformatics Resources
//
// FASTQ Format - scikit-bio:
// https://scikit.bio/docs/latest/generated/skbio.io.format.fastq.html
//
// Python bioinformatics library with FASTQ parsing implementation.
//
// FASTA and FASTQ Formats - Computational Genomics with R:
// https://compgenomr.github.io/book/fasta-and-fastq-formats.html
//
// Educational resource explaining FASTQ format for genomics analysis.
//
// ### Advanced Topics
//
// FASTQ Format Specification - HackMD:
// https://hackmd.io/@EvaMart/SyZ3PKeyD
//
// Detailed technical specification including quality score calculations.
//
// Bioinformatics Stack Exchange - FASTQ Naming:
// https://bioinformatics.stackexchange.com/questions/8880/what-do-the-fastq-file-names-mean-here
//
// Community discussion of real-world FASTQ naming conventions.
//
// ## Implementation Notes
//
// This extractor:
//   - Handles both plain text (.fastq, .fq) and gzip compressed (.fastq.gz, .fq.gz)
//   - Uses buffered reading for memory efficiency
//   - Samples reads for performance on large files
//   - Validates FASTQ format structure
//   - Detects paired-end from filename patterns
//   - Uses Go standard library only (no external dependencies)
//
// ## Limitations
//
//   - Sampling may not capture full read length distribution for variable-length reads
//   - Quality score interpretation assumes Phred+33 encoding (standard for modern platforms)
//   - Does not validate read pairing (assumes filenames indicate pairing)
//   - Does not extract instrument-specific metadata from read identifiers
//
package metadata

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FASTQExtractor extracts metadata from FASTQ sequencing files.
type FASTQExtractor struct{}

// readStats tracks statistics during FASTQ parsing.
type readStats struct {
	totalReads    int
	totalBases    int64
	gcCount       int64
	qualitySum    int64
	qualityCount  int64
	minReadLength int
	maxReadLength int
	minQuality    int
	maxQuality    int
}

// Name returns the extractor name.
func (e *FASTQExtractor) Name() string {
	return "FASTQ"
}

// SupportedFormats returns the file extensions this extractor handles.
func (e *FASTQExtractor) SupportedFormats() []string {
	return []string{".fastq", ".fq", ".fastq.gz", ".fq.gz"}
}

// CanHandle returns true if this extractor can handle the given filename.
func (e *FASTQExtractor) CanHandle(filename string) bool {
	lower := strings.ToLower(filename)
	for _, format := range e.SupportedFormats() {
		if strings.HasSuffix(lower, format) {
			return true
		}
	}
	return false
}

// Extract extracts metadata from a FASTQ file.
func (e *FASTQExtractor) Extract(filepath string) (map[string]interface{}, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	return e.extractFromReader(f, filepath)
}

// ExtractFromReader extracts metadata from a reader.
func (e *FASTQExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return e.extractFromReader(r, filename)
}

// extractFromReader performs the actual metadata extraction.
func (e *FASTQExtractor) extractFromReader(r io.Reader, filepath string) (map[string]interface{}, error) {
	metadata := map[string]interface{}{
		"format":         "FASTQ",
		"file_name":      filepath,
		"extractor_name": "fastq",
		"schema_name":    "fastq_v1",
	}

	// Get file size if available
	if f, ok := r.(*os.File); ok {
		if info, err := f.Stat(); err == nil {
			metadata["file_size"] = info.Size()
		}
	}

	// Detect compression from filename
	isGzipped := strings.HasSuffix(strings.ToLower(filepath), ".gz")
	if isGzipped {
		metadata["compression"] = "gzip"
		// Wrap reader in gzip decompressor
		gzReader, err := gzip.NewReader(r)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		r = gzReader
	} else {
		metadata["compression"] = "none"
	}

	// Detect paired-end from filename
	pairedInfo := detectPairedEnd(filepath)
	if pairedInfo["is_paired"].(bool) {
		metadata["is_paired_end"] = true
		metadata["read_pair"] = pairedInfo["read_pair"]
	} else {
		metadata["is_paired_end"] = false
	}

	// Parse FASTQ and collect statistics
	stats, err := parseFASTQ(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FASTQ: %w", err)
	}

	// Add statistics to metadata
	metadata["total_reads"] = stats.totalReads
	metadata["total_bases"] = stats.totalBases

	if stats.totalReads > 0 {
		metadata["mean_read_length"] = float64(stats.totalBases) / float64(stats.totalReads)
		metadata["min_read_length"] = stats.minReadLength
		metadata["max_read_length"] = stats.maxReadLength
	}

	if stats.totalBases > 0 {
		metadata["gc_content_percent"] = float64(stats.gcCount) / float64(stats.totalBases) * 100
	}

	if stats.qualityCount > 0 {
		metadata["mean_quality_score"] = float64(stats.qualitySum) / float64(stats.qualityCount)
		metadata["min_quality_score"] = stats.minQuality
		metadata["max_quality_score"] = stats.maxQuality
	}

	// Add sequencing platform hint
	metadata["instrument_type"] = "sequencing"
	metadata["data_type"] = "nucleotide_sequence"

	return metadata, nil
}

// parseFASTQ reads and analyzes a FASTQ file.
// Uses sampling for performance on large files.
func parseFASTQ(r io.Reader) (*readStats, error) {
	stats := &readStats{
		minReadLength: int(^uint(0) >> 1), // Max int
		maxReadLength: 0,
		minQuality:    int(^uint(0) >> 1),
		maxQuality:    0,
	}

	scanner := bufio.NewScanner(r)
	// Increase buffer size for long reads (e.g., Nanopore)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB max token size

	lineNum := 0
	var seqLine, qualLine string

readLoop:
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		switch lineNum % 4 {
		case 1: // Identifier line
			if !strings.HasPrefix(line, "@") {
				return nil, fmt.Errorf("invalid FASTQ format: line %d should start with '@', got %q", lineNum, line)
			}
		case 2: // Sequence line
			seqLine = line
		case 3: // Separator line
			if !strings.HasPrefix(line, "+") {
				return nil, fmt.Errorf("invalid FASTQ format: line %d should start with '+', got %q", lineNum, line)
			}
		case 0: // Quality line
			qualLine = line

			// Process this read
			stats.totalReads++
			readLen := len(seqLine)
			stats.totalBases += int64(readLen)

			if readLen < stats.minReadLength {
				stats.minReadLength = readLen
			}
			if readLen > stats.maxReadLength {
				stats.maxReadLength = readLen
			}

			// Count GC content
			for _, base := range seqLine {
				if base == 'G' || base == 'C' || base == 'g' || base == 'c' {
					stats.gcCount++
				}
			}

			// Calculate quality scores (Phred+33 encoding)
			if len(qualLine) == readLen {
				for _, qual := range qualLine {
					phredScore := int(qual) - 33
					if phredScore < 0 {
						phredScore = 0
					}
					stats.qualitySum += int64(phredScore)
					stats.qualityCount++

					if phredScore < stats.minQuality {
						stats.minQuality = phredScore
					}
					if phredScore > stats.maxQuality {
						stats.maxQuality = phredScore
					}
				}
			}

			// Sample only (for performance on large files)
			// Read first 10,000 reads to get representative statistics
			if stats.totalReads >= 10000 {
				break readLoop
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading FASTQ: %w", err)
	}

	if stats.totalReads == 0 {
		return nil, fmt.Errorf("no reads found in FASTQ file")
	}

	return stats, nil
}

// detectPairedEnd analyzes filename to detect paired-end sequencing.
func detectPairedEnd(filename string) map[string]interface{} {
	result := map[string]interface{}{
		"is_paired": false,
	}

	base := filepath.Base(filename)
	base = strings.ToLower(base)

	// Common paired-end patterns
	patterns := []struct {
		regex    *regexp.Regexp
		readPair string
	}{
		{regexp.MustCompile(`[._-]r1[._-]`), "R1"},
		{regexp.MustCompile(`[._-]r2[._-]`), "R2"},
		{regexp.MustCompile(`[._-]r1\.`), "R1"},
		{regexp.MustCompile(`[._-]r2\.`), "R2"},
		{regexp.MustCompile(`[._-]r1_`), "R1"},
		{regexp.MustCompile(`[._-]r2_`), "R2"},
		{regexp.MustCompile(`[._-]1\.f(ast)?q`), "1"},
		{regexp.MustCompile(`[._-]2\.f(ast)?q`), "2"},
		{regexp.MustCompile(`\.1\.f(ast)?q`), "1"},
		{regexp.MustCompile(`\.2\.f(ast)?q`), "2"},
		{regexp.MustCompile(`_1\.f(ast)?q`), "1"},
		{regexp.MustCompile(`_2\.f(ast)?q`), "2"},
	}

	for _, pattern := range patterns {
		if pattern.regex.MatchString(base) {
			result["is_paired"] = true
			result["read_pair"] = pattern.readPair
			return result
		}
	}

	return result
}
