# Metadata Extraction Guide

This guide covers Cicada's automated metadata extraction capabilities for scientific instrument files. Learn how to extract comprehensive metadata, validate against instrument specifications, and integrate metadata extraction into your workflows.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Supported Formats](#supported-formats)
- [Command Reference](#command-reference)
- [Output Formats](#output-formats)
- [Batch Processing](#batch-processing)
- [Validation with Presets](#validation-with-presets)
- [Integration Examples](#integration-examples)
- [Metadata Schema](#metadata-schema)
- [Troubleshooting](#troubleshooting)

---

## Overview

Cicada automatically extracts rich metadata from scientific instrument files, including:

- **File information**: Format, size, compression
- **Content statistics**: Read counts, base counts, sequence lengths
- **Quality metrics**: Quality scores, GC content
- **Instrument details**: Type, pairing information
- **Format-specific fields**: Depends on file type

### Benefits

✅ **Automated**: No manual metadata tracking
✅ **Comprehensive**: Extracts all available metadata
✅ **Fast**: Processes large files efficiently with sampling
✅ **Accurate**: Validates metadata during extraction
✅ **Flexible**: Multiple output formats (JSON, YAML, table)

---

## Quick Start

### Extract Metadata from a Single File

```bash
# Basic extraction (JSON to stdout)
cicada metadata extract sample.fastq.gz

# Human-readable table format
cicada metadata extract sample.fastq.gz --format table

# Save to file
cicada metadata extract sample.fastq.gz --output metadata.json
```

### Example Output

```json
{
  "format": "FASTQ",
  "compression": "gzip",
  "file_name": "/data/sample.fastq.gz",
  "file_size": 2147483648,
  "total_reads": 45623891,
  "total_bases": 6843583650,
  "mean_read_length": 150,
  "gc_content_percent": 42.3,
  "mean_quality_score": 36.8,
  "is_paired_end": true,
  "read_pair": "R1",
  "instrument_type": "sequencing",
  "data_type": "nucleotide_sequence"
}
```

---

## Supported Formats

### Current (v0.2.0)

| Format | Extensions | Description | Metadata Extracted |
|--------|-----------|-------------|-------------------|
| **FASTQ** | `.fastq`, `.fq`, `.fastq.gz`, `.fq.gz` | Nucleotide sequencing data | Read counts, quality scores, GC content, pairing info |

### Planned (v0.3.0+)

| Format | Extensions | Description | Status |
|--------|-----------|-------------|--------|
| **CZI** | `.czi` | Zeiss microscopy images | Planned v0.3.0 |
| **OME-TIFF** | `.ome.tif`, `.ome.tiff` | Open Microscopy Environment | Planned v0.3.0 |
| **TIFF** | `.tif`, `.tiff` | Generic microscopy images | Planned v0.3.0 |
| **ND2** | `.nd2` | Nikon microscopy images | Planned v0.4.0 |
| **LIF** | `.lif` | Leica microscopy images | Planned v0.4.0 |

---

## Command Reference

### `cicada metadata extract`

Extract metadata from a scientific instrument file.

#### Syntax

```bash
cicada metadata extract <path> [flags]
```

#### Arguments

- `<path>` - Path to the file to extract metadata from (required)

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--format` | `-f` | string | `json` | Output format: `json`, `yaml`, or `table` |
| `--output` | `-o` | string | stdout | Output file path (omit for stdout) |
| `--extractor` | | string | auto | Force specific extractor (e.g., `fastq`) |

#### Examples

```bash
# Basic usage
cicada metadata extract data.fastq.gz

# Save to file
cicada metadata extract data.fastq.gz -o metadata.json

# YAML format
cicada metadata extract data.fastq.gz --format yaml

# Table format for human reading
cicada metadata extract data.fastq.gz --format table

# Force specific extractor
cicada metadata extract ambiguous.txt --extractor fastq
```

---

## Output Formats

### JSON Format

Structured data ideal for programmatic processing:

```json
{
  "format": "FASTQ",
  "file_size": 2147483648,
  "total_reads": 45623891,
  "mean_quality_score": 36.8
}
```

**Use cases**:
- Parsing with `jq` or scripts
- Database import
- API integration
- Automated workflows

### YAML Format

Human-readable structured format:

```yaml
format: FASTQ
file_size: 2147483648
total_reads: 45623891
mean_quality_score: 36.8
```

**Use cases**:
- Configuration files
- Documentation
- Manual review
- Git-trackable metadata

### Table Format

Readable format for terminal display:

```
Metadata for: sample.fastq.gz
================================

File Information:
  format                : FASTQ
  compression           : gzip
  file_size             : 2.0 GB

Sequence Statistics:
  total_reads           : 45,623,891
  total_bases           : 6,843,583,650
  mean_read_length      : 150

Quality Metrics:
  mean_quality_score    : 36.8
  min_quality_score     : 12
  max_quality_score     : 41
```

**Use cases**:
- Quick inspection
- Lab notebook entries
- Reports
- Presentations

---

## Batch Processing

### Extract Metadata from Multiple Files

#### Shell Loop

```bash
# Extract metadata for all FASTQ files
for file in *.fastq.gz; do
  basename=$(basename "$file" .fastq.gz)
  cicada metadata extract "$file" \
    --format json \
    --output "metadata/${basename}.json"
done
```

#### Parallel Processing

```bash
# Use GNU parallel for speed
find . -name "*.fastq.gz" | parallel \
  'cicada metadata extract {} --output {}.metadata.json'
```

#### Create Summary Report

```bash
#!/bin/bash
# summarize_metadata.sh - Create summary from metadata files

echo "Dataset Summary Report"
echo "======================"
echo ""

total_reads=0
total_bases=0
file_count=0

for json in metadata/*.json; do
  reads=$(jq '.total_reads' "$json")
  bases=$(jq '.total_bases' "$json")

  total_reads=$((total_reads + reads))
  total_bases=$((total_bases + bases))
  file_count=$((file_count + 1))
done

echo "Total files: $file_count"
echo "Total reads: $(numfmt --grouping $total_reads)"
echo "Total bases: $(numfmt --grouping $total_bases)"
echo "Average reads per file: $(numfmt --grouping $((total_reads / file_count)))"
echo "Average quality: $(jq -s 'map(.mean_quality_score) | add / length' metadata/*.json)"
```

**Output**:
```
Dataset Summary Report
======================

Total files: 48
Total reads: 2,189,946,768
Total bases: 328,492,015,200
Average reads per file: 45,623,891
Average quality: 37.2
```

### Integration with Cicada Sync

Automatically extract metadata when uploading to S3:

```bash
#!/bin/bash
# sync_with_metadata.sh - Upload files and metadata together

FILE="$1"

echo "Processing: $(basename "$FILE")"

# Extract metadata
METADATA_FILE="${FILE}.metadata.json"
cicada metadata extract "$FILE" \
  --format json \
  --output "$METADATA_FILE"

# Upload both file and metadata to S3
cicada sync "$FILE" s3://my-bucket/data/
cicada sync "$METADATA_FILE" s3://my-bucket/metadata/

echo "✓ File and metadata uploaded"
```

**Usage**:
```bash
./sync_with_metadata.sh sample_R1.fastq.gz
```

---

## Validation with Presets

Combine extraction with preset validation to ensure metadata completeness.

### Validate Against Instrument Preset

```bash
# Extract and validate in one command
cicada metadata validate sample.fastq.gz --preset illumina-novaseq
```

**Output**:
```
✓ sample.fastq.gz: valid (FASTQ)
     Quality Score: 100.0/100

Validation Results:
  Present Fields (8):
    ✓ format
    ✓ instrument_type
    ✓ total_reads
    ✓ total_bases
    ✓ mean_quality_score
    ✓ gc_content_percent
    ✓ is_paired_end
    ✓ read_pair

  Missing Optional Fields (0):
    (All optional fields present)

  Errors (0):
    No errors
```

### Available Presets

```bash
# List all available presets
cicada metadata preset list
```

See [PRESETS.md](PRESETS.md) for detailed preset documentation.

---

## Integration Examples

### Nextflow Pipeline Integration

```groovy
// Extract metadata as part of pipeline
process extract_metadata {
  publishDir "${params.outdir}/metadata", mode: 'copy'

  input:
  path fastq_file

  output:
  path "${fastq_file}.metadata.json"

  script:
  """
  cicada metadata extract ${fastq_file} \
    --format json \
    --output ${fastq_file}.metadata.json
  """
}

workflow {
  // ... other processes ...

  fastq_files = Channel.fromPath(params.input)
  metadata = extract_metadata(fastq_files)

  // Use metadata in downstream processes
  metadata.view()
}
```

### Snakemake Pipeline Integration

```python
# Snakefile
rule extract_metadata:
    input:
        fastq="{sample}.fastq.gz"
    output:
        metadata="{sample}.metadata.json"
    shell:
        """
        cicada metadata extract {input.fastq} \
            --format json \
            --output {output.metadata}
        """

rule summarize_metadata:
    input:
        metadata=expand("{sample}.metadata.json", sample=SAMPLES)
    output:
        summary="summary.txt"
    run:
        import json
        total_reads = 0
        for f in input.metadata:
            with open(f) as fh:
                meta = json.load(fh)
                total_reads += meta['total_reads']

        with open(output.summary, 'w') as out:
            out.write(f"Total reads: {total_reads}\n")
```

### Python Script Integration

```python
#!/usr/bin/env python3
"""Extract and analyze metadata using Cicada."""

import subprocess
import json
from pathlib import Path

def extract_metadata(fastq_path):
    """Extract metadata using Cicada CLI."""
    result = subprocess.run(
        ['cicada', 'metadata', 'extract', str(fastq_path), '--format', 'json'],
        capture_output=True,
        text=True,
        check=True
    )
    return json.loads(result.stdout)

def analyze_quality(metadata):
    """Analyze quality metrics."""
    mean_q = metadata['mean_quality_score']

    if mean_q >= 35:
        return "Excellent"
    elif mean_q >= 30:
        return "Good"
    elif mean_q >= 25:
        return "Acceptable"
    else:
        return "Poor"

def main():
    """Process all FASTQ files in directory."""
    fastq_files = Path('.').glob('*.fastq.gz')

    results = []
    for fastq in fastq_files:
        meta = extract_metadata(fastq)
        quality = analyze_quality(meta)

        results.append({
            'file': str(fastq),
            'reads': meta['total_reads'],
            'quality': quality,
            'quality_score': meta['mean_quality_score']
        })

    # Print summary
    print(f"Processed {len(results)} files")
    for r in results:
        print(f"  {r['file']}: {r['reads']:,} reads, {r['quality']} quality ({r['quality_score']:.1f})")

if __name__ == '__main__':
    main()
```

**Output**:
```
Processed 3 files
  sample1.fastq.gz: 45,623,891 reads, Excellent quality (37.2)
  sample2.fastq.gz: 42,891,234 reads, Excellent quality (36.8)
  sample3.fastq.gz: 48,234,567 reads, Good quality (34.5)
```

---

## Metadata Schema

### Common Fields (All Formats)

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `format` | string | File format | `"FASTQ"` |
| `file_name` | string | Full file path | `"/data/sample.fastq.gz"` |
| `file_size` | integer | File size in bytes | `2147483648` |
| `compression` | string | Compression type | `"gzip"`, `"none"` |
| `instrument_type` | string | Instrument category | `"sequencing"`, `"microscopy"` |
| `data_type` | string | Data category | `"nucleotide_sequence"`, `"image"` |
| `extractor_name` | string | Extractor used | `"fastq"` |
| `schema_name` | string | Metadata schema version | `"fastq_v1"` |

### FASTQ-Specific Fields

| Field | Type | Description | Range/Values |
|-------|------|-------------|--------------|
| `total_reads` | integer | Number of reads | > 0 |
| `total_bases` | integer | Total nucleotides | > 0 |
| `mean_read_length` | integer | Average read length | > 0 |
| `min_read_length` | integer | Shortest read | > 0 |
| `max_read_length` | integer | Longest read | > 0 |
| `gc_content_percent` | float | GC content percentage | 0-100 |
| `mean_quality_score` | float | Average Phred score | 0-93 |
| `min_quality_score` | integer | Minimum Phred score | 0-93 |
| `max_quality_score` | integer | Maximum Phred score | 0-93 |
| `is_paired_end` | boolean | Paired-end detected | `true`, `false` |
| `read_pair` | string | Read pair identifier | `"R1"`, `"R2"`, `"1"`, `"2"` |

### Quality Score Interpretation

FASTQ uses Phred+33 quality scores:

| Score Range | Quality | Interpretation |
|-------------|---------|----------------|
| 40-41 | Excellent | Error rate < 0.01% |
| 30-39 | Very Good | Error rate < 0.1% |
| 20-29 | Good | Error rate < 1% |
| 10-19 | Marginal | Error rate < 10% |
| 0-9 | Poor | Error rate > 10% |

---

## Troubleshooting

### File Not Recognized

**Problem**: `Error: no extractor found for file`

**Solutions**:
1. Check file extension is supported:
   ```bash
   # List supported formats
   cicada metadata list
   ```

2. Force specific extractor:
   ```bash
   cicada metadata extract file.txt --extractor fastq
   ```

3. Verify file is not corrupted:
   ```bash
   # For gzipped files
   gunzip -t file.fastq.gz

   # For plain files
   head -4 file.fastq
   ```

### Invalid File Format

**Problem**: `Error: invalid FASTQ format`

**Common causes**:
- File is empty
- Missing `@` header line
- Missing `+` separator line
- Sequence and quality lengths don't match

**Solution**: Validate FASTQ format
```bash
# Check first record
head -4 file.fastq

# Expected format:
# @HEADER
# SEQUENCE
# +
# QUALITY
```

### Large File Processing

**Problem**: Extraction taking too long

**Solution**: Cicada automatically samples large files
- Processes up to 10,000 reads by default
- Provides representative statistics
- Fast even for 100+ GB files

**Manual control** (future feature):
```bash
# Not yet implemented, but planned:
cicada metadata extract file.fastq.gz --sample-size 5000
```

### Gzip Errors

**Problem**: `Error: gzip: invalid compressed data`

**Solutions**:
1. Verify file integrity:
   ```bash
   gunzip -t file.fastq.gz
   ```

2. Re-download file if corrupted

3. Check disk space:
   ```bash
   df -h /path/to/file
   ```

### Permission Denied

**Problem**: `Error: permission denied`

**Solutions**:
1. Check file permissions:
   ```bash
   ls -l file.fastq.gz
   ```

2. Add read permission:
   ```bash
   chmod +r file.fastq.gz
   ```

3. Check directory permissions:
   ```bash
   ls -ld $(dirname file.fastq.gz)
   ```

---

## Best Practices

### 1. Extract Metadata Early

Extract metadata as soon as data is generated:
- Easier to track while experiment is fresh
- Can detect quality issues early
- Enables immediate data discovery

### 2. Store Metadata with Data

Keep metadata files alongside data files:
```
project/
├── data/
│   ├── sample_R1.fastq.gz
│   ├── sample_R2.fastq.gz
├── metadata/
│   ├── sample_R1.metadata.json
│   └── sample_R2.metadata.json
```

Or use sidecar files:
```
project/
├── sample_R1.fastq.gz
├── sample_R1.fastq.gz.metadata.json
├── sample_R2.fastq.gz
└── sample_R2.fastq.gz.metadata.json
```

### 3. Validate Against Presets

Always validate against instrument presets to ensure completeness:
```bash
cicada metadata validate file.fastq.gz --preset illumina-novaseq
```

### 4. Version Your Metadata

Track metadata schema versions:
- Check `schema_name` field in output
- Document which version was used
- Update extraction when schemas change

### 5. Automate Extraction

Integrate metadata extraction into your workflows:
- Add to upload scripts
- Include in pipeline stages
- Trigger on file system events

---

## Performance

### Benchmarks

Typical extraction performance (single-threaded):

| File Size | Reads | Time | Throughput |
|-----------|-------|------|------------|
| 100 MB | 500K | 0.5s | 200 MB/s |
| 1 GB | 5M | 2s | 500 MB/s |
| 10 GB | 50M | 5s | 2 GB/s (sampled) |
| 100 GB | 500M | 5s | 20 GB/s (sampled) |

**Note**: Large files are automatically sampled (10,000 reads), providing fast representative statistics.

### Optimization Tips

1. **Use local storage**: Extract from local files, not network mounts
2. **SSD recommended**: Faster I/O for large files
3. **Batch processing**: Use parallel processing for multiple files
4. **Sampling**: Automatic for large files, no configuration needed

---

## Next Steps

- **Validate metadata**: See [PRESETS.md](PRESETS.md)
- **Prepare for DOI**: See [DOI_WORKFLOW.md](DOI_WORKFLOW.md)
- **User scenarios**: See [USER_SCENARIOS_v0.2.0.md](USER_SCENARIOS_v0.2.0.md)
- **Integration testing**: See [../INTEGRATION_TESTING.md](../INTEGRATION_TESTING.md)

---

## Support

**Questions or Issues?**
- 📖 Full documentation: [README.md](../README.md)
- 🐛 Report bugs: [GitHub Issues](https://github.com/scttfrdmn/cicada/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/scttfrdmn/cicada/discussions)
- 🚀 Feature requests: [GitHub Issues](https://github.com/scttfrdmn/cicada/issues)
