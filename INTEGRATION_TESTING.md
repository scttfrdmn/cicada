# Integration Testing Guide

This document describes the integration testing approach for Cicada v0.2.0 metadata and DOI functionality.

## Overview

Integration tests verify end-to-end functionality with real files and actual data processing. Unlike unit tests that mock dependencies, integration tests exercise the complete system to ensure all components work together correctly.

## Test Organization

Integration tests are located in `internal/integration/` and organized by functional area:

```
internal/integration/
├── metadata_extraction_test.go    # Metadata extraction workflows
├── doi_workflow_test.go           # DOI preparation and validation
├── preset_integration_test.go     # Preset system integration
└── cli_integration_test.go        # CLI command testing
```

## Test Coverage

### 1. Metadata Extraction Tests (`metadata_extraction_test.go`)

**Tests:** 6 tests covering end-to-end metadata extraction

- **TestMetadataExtraction_FASTQ**: Full FASTQ extraction with quality metrics
  - Creates real FASTQ file with sequences and quality scores
  - Verifies: format, read counts, base counts, GC content, quality scores, pairing detection

- **TestMetadataExtraction_FASTQGzipped**: Gzip compression handling
  - Tests transparent decompression of .fastq.gz files
  - Verifies compression detection and correct extraction

- **TestMetadataExtraction_AutoDetection**: Format auto-detection
  - Tests .fastq, .fq, and .fastq.gz extensions
  - Verifies correct extractor selection

- **TestMetadataExtraction_InvalidFiles**: Error handling
  - Tests empty files, malformed FASTQ, missing files
  - Verifies graceful error reporting

- **TestMetadataExtraction_LargeFile**: Performance with large datasets
  - Creates 50,000 read file, verifies 10,000 read sampling
  - Ensures efficient processing of large files

- **TestMetadataExtraction_ConcurrentExtractions**: Thread safety
  - Extracts 10 files concurrently
  - Verifies no race conditions or data corruption

### 2. DOI Workflow Tests (`doi_workflow_test.go`)

**Tests:** 6 tests covering DOI preparation and validation

- **TestDOIWorkflow_EndToEnd**: Complete extract → map → validate workflow
  - Tests full pipeline from file to DOI-ready metadata
  - Verifies dataset creation, validation, quality scoring

- **TestDOIWorkflow_WithEnrichment**: User-provided metadata enrichment
  - Tests ORCID, affiliation, description enhancement
  - Verifies quality score improvement with enrichment

- **TestDOIWorkflow_ValidationStrictness**: Lenient vs strict validation
  - Tests different validation modes
  - Verifies appropriate error/warning levels

- **TestDOIWorkflow_MultipleFileTypes**: Different file formats
  - Tests single-end and paired-end FASTQ
  - Verifies format-specific handling

- **TestDOIWorkflow_ValidateMetadata**: Validation-only workflow
  - Tests validation without DOI preparation
  - Verifies quality assessment accuracy

- **TestDOIWorkflow_QualityScoring**: Score calculation
  - Compares minimal vs rich metadata scoring
  - Verifies score ranges and improvements

### 3. Preset System Tests (`preset_integration_test.go`)

**Tests:** 8 tests covering preset validation and management

- **TestPresetIntegration_MetadataValidation**: Validate extracted metadata against presets
  - Tests Illumina NovaSeq preset validation
  - Verifies field presence, quality scoring

- **TestPresetIntegration_FindPresets**: Search presets by criteria
  - Tests filtering by manufacturer and instrument type
  - Verifies correct preset matches

- **TestPresetIntegration_AllPresets**: All 8 default presets
  - Tests each preset (Zeiss LSM 880/900/980, Illumina NovaSeq/MiSeq/NextSeq, Generic)
  - Verifies completeness and validation behavior

- **TestPresetIntegration_ExtractAndValidate**: Full extract → validate workflow
  - Tests integration of extraction and preset validation
  - Verifies seamless workflow

- **TestPresetIntegration_FieldValidation**: Field-level validation
  - Tests valid metadata, missing fields, invalid types
  - Verifies error reporting accuracy

- **TestPresetIntegration_QualityScoring**: Preset quality scoring
  - Compares minimal (required only) vs rich (required + optional)
  - Verifies scoring algorithm

- **TestPresetIntegration_TemplateGeneration**: Generate metadata templates
  - Tests template creation from presets
  - Verifies all required fields present

### 4. CLI Command Tests (`cli_integration_test.go`)

**Tests:** 9 tests covering command-line interface

- **TestCLI_MetadataExtract**: `cicada metadata extract` command
  - Tests JSON output format
  - Verifies metadata extraction via CLI

- **TestCLI_MetadataValidate**: `cicada metadata validate` command
  - Tests preset-based validation
  - Verifies validation results

- **TestCLI_MetadataListPresets**: `cicada metadata preset list` command
  - Tests preset listing functionality
  - Verifies command execution

- **TestCLI_DOIPrepare**: `cicada doi prepare` command
  - Tests DOI preparation with enrichment file
  - Verifies successful preparation

- **TestCLI_DOIValidate**: `cicada doi validate` command
  - Tests validation failure for incomplete metadata
  - Verifies error reporting

- **TestCLI_VersionCommand**: `cicada version` command
  - Tests version display
  - Verifies command execution

- **TestCLI_InvalidFile**: Error handling for non-existent files
  - Tests appropriate error messages
  - Verifies graceful failure

- **TestCLI_OutputFormats**: JSON, YAML, and table formats
  - Tests all three output formats
  - Verifies format-specific content

- **TestCLI_HelpCommand**: `cicada --help` command
  - Tests help text generation
  - Verifies all commands listed

## Testing Principles

### 1. Real Data, No Mocks

All integration tests use real files with actual scientific data:

```go
fastqContent := `@Illumina_NovaSeq_Run123
GATTTGGGGTTCAAAGCAGTATCGATCAAATAGTAAATCCATTTGTTCAACTCACAGTTT
+
IIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIII
@Illumina_NovaSeq_Run123
GCGCGCGCGCGCGCGCGCGCGCGC
+
IIIIIIIIIIIIIIIIIIIIIIIII
`
```

- Uses valid FASTQ format with Phred+33 quality scores
- Creates actual sequences with realistic GC content
- Tests actual file I/O and parsing

### 2. Temporary Directories

All test files are created in temporary directories:

```go
testFile := filepath.Join(t.TempDir(), "sample.fastq")
```

- Automatically cleaned up after tests
- No pollution of source tree
- Safe for concurrent execution

### 3. Comprehensive Verification

Tests verify multiple aspects of functionality:

```go
// Verify format
if extracted["format"] != "FASTQ" { ... }

// Verify counts
if extracted["total_reads"] != 2 { ... }

// Verify calculations
if gcContent < 0 || gcContent > 100 { ... }
```

### 4. Error Path Testing

Tests verify both success and failure scenarios:

```go
// Test should fail for invalid file
err := registry.Extract(invalidFile)
if err == nil {
    t.Error("Extract() should have returned error")
}
```

## Running Integration Tests

### Run All Integration Tests

```bash
go test ./internal/integration/...
```

### Run Specific Test Suite

```bash
# Metadata extraction tests
go test ./internal/integration -run TestMetadataExtraction

# DOI workflow tests
go test ./internal/integration -run TestDOIWorkflow

# Preset tests
go test ./internal/integration -run TestPresetIntegration

# CLI tests
go test ./internal/integration -run TestCLI
```

### Run Individual Test

```bash
go test ./internal/integration -run TestDOIWorkflow_EndToEnd
```

### Verbose Output

```bash
go test -v ./internal/integration/...
```

### Skip Long-Running Tests

```bash
go test -short ./internal/integration/...
```

Note: The large file test (`TestMetadataExtraction_LargeFile`) is skipped in short mode.

## Test Data Patterns

### FASTQ Files

```go
// Minimal valid FASTQ
fastqContent := `@SEQ_ID
ACGT
+
IIII
`

// Realistic FASTQ with quality variance
fastqContent := `@SEQ_ID_1 test sequence 1
GATTTGGGGTTCAAAGCAGTATCGATCAAATAGTAAATCCATTTGTTCAACTCACAGTTT
+
!''*((((***+))%%%++)(%%%%).1***-+*''))**55CCF>>>>>>CCCCCCC65
`
```

### Enrichment Data

```json
{
  "title": "Research Dataset Title",
  "authors": [
    {
      "name": "Dr. Researcher Name",
      "orcid": "0000-0001-2345-6789",
      "affiliation": "University Lab"
    }
  ],
  "description": "Detailed description of the dataset",
  "keywords": ["keyword1", "keyword2"]
}
```

## Test Metrics

### Coverage

- **29 integration tests** covering all major workflows
- **100+ test cases** across all scenarios
- **All file formats** supported by v0.2.0 (FASTQ, FASTQ.gz)
- **All CLI commands** tested with real execution

### Performance

- Average test suite runtime: **~0.5 seconds**
- Large file test: **~0.3 seconds** (50,000 reads)
- Concurrent test: **<0.1 seconds** (10 concurrent extractions)

## Adding New Integration Tests

When adding new features, follow this pattern:

1. **Create test file** in `internal/integration/`
2. **Use real data** - no mocks
3. **Test complete workflow** - end-to-end functionality
4. **Verify multiple aspects** - not just success/failure
5. **Test error cases** - invalid input, missing files, etc.
6. **Use temp directories** - clean up after tests
7. **Document test purpose** - clear comments and test names

### Example Template

```go
func TestNewFeature_Workflow(t *testing.T) {
    // Create test data
    testFile := filepath.Join(t.TempDir(), "test.dat")
    content := `... realistic test data ...`
    if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
        t.Fatalf("Failed to create test file: %v", err)
    }

    // Execute workflow
    result, err := ExecuteWorkflow(testFile)
    if err != nil {
        t.Fatalf("Workflow failed: %v", err)
    }

    // Verify results
    if result.Field != expectedValue {
        t.Errorf("Field = %v, want %v", result.Field, expectedValue)
    }
}
```

## Continuous Integration

Integration tests are designed to run in CI/CD environments:

- **Fast execution** - complete suite runs in < 1 second
- **No external dependencies** - self-contained tests
- **Deterministic** - consistent results across runs
- **Portable** - works on all platforms (Linux, macOS, Windows)

## Known Limitations

1. **CLI Output Capture**: Some CLI commands write directly to stdout/stderr using `fmt.Printf`, which cannot be captured by `cmd.SetOut()`. Tests use `--output` flag to write to files instead.

2. **Version Command**: The `version` command output goes to stdout and is not easily captured in tests. Test verifies command executes without error instead of checking output.

3. **Large File Sampling**: The large file test creates 50,000 reads but only 10,000 are processed due to sampling. This is expected behavior but worth noting.

## Future Enhancements

Potential areas for expansion:

- **More file formats**: Add tests for CZI, OME-TIFF when extractors are implemented
- **Provider integration**: Test actual DataCite/Zenodo API calls with sandbox environments
- **Storage integration**: Test S3/GCS upload workflows
- **Watch mode**: Test file watching and auto-sync functionality
- **Performance benchmarks**: Add benchmark tests for large-scale processing

## Questions and Support

For questions about integration testing:

- **Architecture**: See `ARCHITECTURE.md` for system design
- **Contributing**: See `CONTRIBUTING.md` for development guidelines
- **Issues**: Report test failures as bugs in GitHub issues
