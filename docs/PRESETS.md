# Instrument Preset Guide

**Status:** v0.2.0 Documentation
**Audience:** Lab managers, data curators, researchers validating metadata

## Overview

Instrument presets provide pre-configured validation rules and metadata templates for specific scientific instruments. They ensure your metadata meets field-specific standards and improves reproducibility for common instruments.

**Key Benefits:**

- **Standardization**: Consistent metadata across instruments and labs
- **Validation**: Automatic checking of required and optional fields
- **Quality Scoring**: Objective assessment of metadata completeness (0-100 scale)
- **Time Savings**: No need to manually specify field requirements
- **Compliance**: Ensures metadata meets publisher and repository requirements

**Available Presets (v0.2.0):**

- **Microscopy**: Zeiss LSM 880, Zeiss LSM 900, Zeiss LSM 980, Generic Microscopy
- **Sequencing**: Illumina NovaSeq, Illumina MiSeq, Illumina NextSeq, Generic Sequencing

## Quick Start

### Validate with Preset

```bash
# Extract and validate FASTQ file with Illumina NovaSeq preset
cicada metadata extract sample_R1.fastq.gz --preset illumina-novaseq

# Validate existing metadata file
cicada metadata validate metadata.json --preset illumina-novaseq
```

### List Available Presets

```bash
# List all presets
cicada metadata preset list

# Find Illumina presets
cicada metadata preset list --manufacturer Illumina

# Find microscopy presets
cicada metadata preset list --type microscopy
```

### Check Quality Score

```bash
# Extract with preset validation to see quality score
cicada metadata extract sample.fastq --preset illumina-novaseq --format yaml

# Output includes:
# quality_score: 78
# validation_errors: []
# validation_warnings: ["Optional field 'description' not provided"]
```

## Available Presets

### Illumina NovaSeq

**Preset ID:** `illumina-novaseq`

High-throughput sequencing platform for large-scale genomics projects.

**Manufacturer:** Illumina
**Instrument Type:** Sequencing
**Platform:** NovaSeq 6000

**Required Fields (60% of quality score):**
- `format` - File format (e.g., "FASTQ")
- `instrument_manufacturer` - "Illumina"
- `instrument_model` - "NovaSeq 6000" or compatible
- `sequencing_platform` - "Illumina"
- `total_reads` - Number of reads in file
- `read_length` - Length of sequencing reads (bp)

**Optional Fields (40% of quality score):**
- `run_id` - Sequencing run identifier
- `flowcell_id` - Flowcell identifier
- `lane` - Lane number (1-4 for NovaSeq)
- `barcode` - Sample barcode/index sequence
- `quality_encoding` - Quality score encoding (e.g., "Phred+33")
- `gc_content` - GC percentage
- `mean_quality_score` - Average quality score
- `is_paired_end` - Paired-end sequencing flag

**Best For:**
- Whole genome sequencing
- Large RNA-seq projects
- High-depth targeted sequencing

### Illumina MiSeq

**Preset ID:** `illumina-miseq`

Desktop sequencer for targeted gene sequencing and small genome applications.

**Manufacturer:** Illumina
**Instrument Type:** Sequencing
**Platform:** MiSeq

**Required Fields:**
- `format`
- `instrument_manufacturer` - "Illumina"
- `instrument_model` - "MiSeq"
- `sequencing_platform` - "Illumina"
- `total_reads`
- `read_length`

**Optional Fields:**
- `run_id`
- `flowcell_id`
- `barcode`
- `quality_encoding`
- `gc_content`
- `mean_quality_score`
- `is_paired_end`

**Best For:**
- Amplicon sequencing (16S, targeted genes)
- Small genome sequencing
- Validation studies

### Illumina NextSeq

**Preset ID:** `illumina-nextseq`

Mid-throughput sequencer balancing speed and data output.

**Manufacturer:** Illumina
**Instrument Type:** Sequencing
**Platform:** NextSeq 500/550/1000/2000

**Required Fields:**
- `format`
- `instrument_manufacturer` - "Illumina"
- `instrument_model` - "NextSeq" (any model)
- `sequencing_platform` - "Illumina"
- `total_reads`
- `read_length`

**Optional Fields:**
- `run_id`
- `flowcell_id`
- `barcode`
- `quality_encoding`
- `gc_content`
- `mean_quality_score`
- `is_paired_end`

**Best For:**
- RNA-seq
- Exome sequencing
- Targeted panels

### Zeiss LSM 880

**Preset ID:** `zeiss-lsm-880`

Laser scanning confocal microscope with spectral detection.

**Manufacturer:** Zeiss
**Instrument Type:** Microscopy
**Platform:** LSM 880

**Required Fields (60% of quality score):**
- `format` - Image format (e.g., "CZI", "TIFF")
- `instrument_manufacturer` - "Zeiss"
- `instrument_model` - "LSM 880"
- `microscopy_type` - "Confocal"
- `image_width` - Width in pixels
- `image_height` - Height in pixels

**Optional Fields (40% of quality score):**
- `bit_depth` - Bits per pixel
- `pixel_size_x` - X dimension pixel size (μm)
- `pixel_size_y` - Y dimension pixel size (μm)
- `pixel_size_z` - Z dimension pixel size (μm)
- `channels` - Number of imaging channels
- `z_slices` - Number of Z-stack slices
- `time_points` - Number of time points
- `objective_magnification` - Objective lens magnification
- `objective_na` - Numerical aperture
- `excitation_wavelengths` - Laser wavelengths (nm)
- `emission_wavelengths` - Detection wavelengths (nm)

**Best For:**
- Live cell imaging
- Multi-channel fluorescence
- Z-stack acquisition

### Zeiss LSM 900

**Preset ID:** `zeiss-lsm-900`

Advanced confocal with Airyscan 2 for super-resolution imaging.

**Manufacturer:** Zeiss
**Instrument Type:** Microscopy
**Platform:** LSM 900

**Required Fields:**
- `format`
- `instrument_manufacturer` - "Zeiss"
- `instrument_model` - "LSM 900"
- `microscopy_type` - "Confocal"
- `image_width`
- `image_height`

**Optional Fields:**
- `bit_depth`
- `pixel_size_x`
- `pixel_size_y`
- `pixel_size_z`
- `channels`
- `z_slices`
- `time_points`
- `objective_magnification`
- `objective_na`
- `excitation_wavelengths`
- `emission_wavelengths`
- `airyscan_mode` - "SR" (super-resolution) or "FAST"

**Best For:**
- Super-resolution imaging
- Thick tissue sections
- High-speed live imaging

### Zeiss LSM 980

**Preset ID:** `zeiss-lsm-980`

Flagship confocal with multiphoton and spectral unmixing capabilities.

**Manufacturer:** Zeiss
**Instrument Type:** Microscopy
**Platform:** LSM 980

**Required Fields:**
- `format`
- `instrument_manufacturer` - "Zeiss"
- `instrument_model` - "LSM 980"
- `microscopy_type` - "Confocal"
- `image_width`
- `image_height`

**Optional Fields:**
- `bit_depth`
- `pixel_size_x`
- `pixel_size_y`
- `pixel_size_z`
- `channels`
- `z_slices`
- `time_points`
- `objective_magnification`
- `objective_na`
- `excitation_wavelengths`
- `emission_wavelengths`
- `airyscan_mode`
- `multiphoton` - Boolean flag for multiphoton mode

**Best For:**
- Deep tissue imaging
- Spectral imaging
- Multiphoton microscopy

### Generic Sequencing

**Preset ID:** `generic-sequencing`

Universal preset for any sequencing platform.

**Manufacturer:** Any
**Instrument Type:** Sequencing
**Platform:** Generic

**Required Fields:**
- `format`
- `sequencing_platform` - Platform name
- `total_reads`

**Optional Fields:**
- `read_length`
- `is_paired_end`
- `quality_encoding`
- `gc_content`
- `instrument_manufacturer`
- `instrument_model`

**Best For:**
- Non-Illumina platforms (PacBio, Oxford Nanopore, MGI)
- Legacy data
- Unknown instrument details

### Generic Microscopy

**Preset ID:** `generic-microscopy`

Universal preset for any microscopy platform.

**Manufacturer:** Any
**Instrument Type:** Microscopy
**Platform:** Generic

**Required Fields:**
- `format`
- `instrument_type` - "microscopy"
- `image_width`
- `image_height`

**Optional Fields:**
- `bit_depth`
- `pixel_size_x`
- `pixel_size_y`
- `channels`
- `microscopy_type`
- `instrument_manufacturer`
- `instrument_model`

**Best For:**
- Non-Zeiss microscopes (Nikon, Olympus, Leica)
- Widefield/epifluorescence
- Unknown instrument details

## Command Reference

### List Presets

```bash
cicada metadata preset list [flags]
```

**Flags:**
- `--manufacturer, -m` - Filter by manufacturer (e.g., "Illumina", "Zeiss")
- `--type, -t` - Filter by instrument type ("sequencing", "microscopy")
- `--format, -f` - Output format: json, yaml, or table (default: table)

**Examples:**

```bash
# List all presets in table format
cicada metadata preset list

# List Illumina presets
cicada metadata preset list --manufacturer Illumina

# List microscopy presets as JSON
cicada metadata preset list --type microscopy --format json
```

### Show Preset Details

```bash
cicada metadata preset show <preset-id> [flags]
```

**Flags:**
- `--format, -f` - Output format: json, yaml, or table (default: yaml)

**Examples:**

```bash
# Show Illumina NovaSeq preset details
cicada metadata preset show illumina-novaseq

# Show as JSON
cicada metadata preset show zeiss-lsm-880 --format json
```

### Validate with Preset

```bash
cicada metadata validate <metadata-file> --preset <preset-id> [flags]
```

**Flags:**
- `--preset, -p` - Preset ID to validate against
- `--strict` - Treat warnings as errors
- `--format, -f` - Output format: json, yaml, or table (default: table)

**Examples:**

```bash
# Validate metadata file
cicada metadata validate metadata.json --preset illumina-novaseq

# Strict validation (warnings become errors)
cicada metadata validate metadata.json --preset illumina-novaseq --strict
```

### Extract with Preset

```bash
cicada metadata extract <file> --preset <preset-id> [flags]
```

Automatically validates extracted metadata against preset requirements.

**Examples:**

```bash
# Extract and validate in one step
cicada metadata extract sample.fastq.gz --preset illumina-novaseq

# Extract with preset and save to file
cicada metadata extract sample.fastq.gz \
  --preset illumina-novaseq \
  --format json \
  --output metadata.json
```

## Quality Scoring

Presets use a 0-100 quality score based on field completeness:

- **60 points**: Required fields (all must be present)
- **40 points**: Optional fields (proportional to completeness)

### Score Calculation

```
Score = (Required Fields × 60) + (Optional Fields × 40)

Where:
- Required Fields = 1.0 if all present, 0.0 otherwise
- Optional Fields = (present optional fields) / (total optional fields)
```

### Score Interpretation

| Score | Interpretation | Publication Ready? |
|-------|----------------|-------------------|
| 90-100 | Excellent | ✅ Yes - High quality |
| 75-89 | Good | ✅ Yes - Acceptable |
| 60-74 | Adequate | ⚠️ Maybe - Missing recommended fields |
| 0-59 | Insufficient | ❌ No - Missing required fields |

### Example Score Progression

**Minimal Metadata (Score: 60)**
```yaml
format: FASTQ
instrument_manufacturer: Illumina
instrument_model: NovaSeq 6000
sequencing_platform: Illumina
total_reads: 1000000
read_length: 150
```
- ✅ All 6 required fields present: 60 points
- ❌ 0 of 8 optional fields: 0 points
- **Total: 60/100**

**Basic Metadata (Score: 80)**
```yaml
format: FASTQ
instrument_manufacturer: Illumina
instrument_model: NovaSeq 6000
sequencing_platform: Illumina
total_reads: 1000000
read_length: 150
run_id: RUN_123
quality_encoding: Phred+33
gc_content: 52.3
is_paired_end: true
```
- ✅ All 6 required fields: 60 points
- ✅ 4 of 8 optional fields: 20 points
- **Total: 80/100**

**Complete Metadata (Score: 100)**
```yaml
format: FASTQ
instrument_manufacturer: Illumina
instrument_model: NovaSeq 6000
sequencing_platform: Illumina
total_reads: 1000000
read_length: 150
run_id: RUN_123
flowcell_id: FLOWCELL_456
lane: 2
barcode: AGTCACTA
quality_encoding: Phred+33
gc_content: 52.3
mean_quality_score: 36.8
is_paired_end: true
```
- ✅ All 6 required fields: 60 points
- ✅ All 8 optional fields: 40 points
- **Total: 100/100**

## Validation Workflow

### Step 1: Extract Metadata

```bash
# Extract without validation
cicada metadata extract sample_R1.fastq.gz --output metadata.json
```

### Step 2: Choose Preset

```bash
# List available presets
cicada metadata preset list --type sequencing

# Show preset requirements
cicada metadata preset show illumina-novaseq
```

### Step 3: Validate

```bash
# Validate extracted metadata
cicada metadata validate metadata.json --preset illumina-novaseq
```

**Output:**
```
Validation Results
==================

Preset: illumina-novaseq (Illumina NovaSeq 6000)
Quality Score: 80/100

✅ Required Fields (6/6)
  ✓ format
  ✓ instrument_manufacturer
  ✓ instrument_model
  ✓ sequencing_platform
  ✓ total_reads
  ✓ read_length

⚠️ Optional Fields (4/8)
  ✓ run_id
  ✓ quality_encoding
  ✓ gc_content
  ✓ is_paired_end
  ✗ flowcell_id - Not provided
  ✗ lane - Not provided
  ✗ barcode - Not provided
  ✗ mean_quality_score - Not provided

Recommendation: Consider adding optional fields to improve quality score.
```

### Step 4: Enrich Metadata (Optional)

Create enrichment file with missing fields:

**enrich.yaml:**
```yaml
flowcell_id: FLOWCELL_456
lane: 2
barcode: AGTCACTA
mean_quality_score: 36.8
```

```bash
# Re-validate with enrichment
cicada metadata validate metadata.json \
  --preset illumina-novaseq \
  --enrich enrich.yaml
```

**New Output:**
```
Quality Score: 100/100
✅ All fields present
```

## Finding Presets

### By Manufacturer

```bash
# Find all Illumina presets
cicada metadata preset list --manufacturer Illumina

# Output:
# ID                  Manufacturer  Platform        Type
# illumina-novaseq    Illumina      NovaSeq 6000    sequencing
# illumina-miseq      Illumina      MiSeq           sequencing
# illumina-nextseq    Illumina      NextSeq         sequencing
```

### By Instrument Type

```bash
# Find all microscopy presets
cicada metadata preset list --type microscopy

# Output:
# ID                  Manufacturer  Platform    Type
# zeiss-lsm-880       Zeiss         LSM 880     microscopy
# zeiss-lsm-900       Zeiss         LSM 900     microscopy
# zeiss-lsm-980       Zeiss         LSM 980     microscopy
# generic-microscopy  Any           Generic     microscopy
```

### Programmatically (Future: v0.3.0+)

**Python example:**
```python
from cicada import PresetRegistry

registry = PresetRegistry()

# Find by manufacturer
illumina_presets = registry.find(manufacturer="Illumina")

# Find by type
microscopy_presets = registry.find(instrument_type="microscopy")

# Find by both
zeiss_microscopes = registry.find(
    manufacturer="Zeiss",
    instrument_type="microscopy"
)

# Get specific preset
preset = registry.get("illumina-novaseq")
print(f"Required fields: {preset.required_fields}")
```

## Integration Examples

### Nextflow Pipeline

```nextflow
process validate_metadata {
    input:
    path fastq
    val preset_id

    output:
    path "metadata.json"

    script:
    """
    # Extract and validate with preset
    cicada metadata extract ${fastq} \
      --preset ${preset_id} \
      --format json \
      --output metadata.json

    # Check quality score
    score=\$(jq -r '.quality_score' metadata.json)
    if (( \$(echo "\$score < 75" | bc -l) )); then
        echo "Warning: Quality score \$score is below recommended threshold"
        exit 1
    fi
    """
}

workflow {
    fastqs = Channel.fromPath("data/*.fastq.gz")
    preset = "illumina-novaseq"

    validate_metadata(fastqs, preset)
}
```

### Snakemake Workflow

```python
rule validate_metadata:
    input:
        fastq="data/{sample}.fastq.gz"
    output:
        metadata="metadata/{sample}.json",
        report="reports/{sample}_validation.txt"
    params:
        preset="illumina-novaseq",
        min_score=75
    shell:
        """
        # Extract with preset
        cicada metadata extract {input.fastq} \
          --preset {params.preset} \
          --format json \
          --output {output.metadata}

        # Generate validation report
        cicada metadata validate {output.metadata} \
          --preset {params.preset} \
          --format table > {output.report}

        # Check quality score
        score=$(jq -r '.quality_score' {output.metadata})
        if [ "$score" -lt "{params.min_score}" ]; then
            echo "Quality score $score below threshold {params.min_score}"
            exit 1
        fi
        """
```

### Python Script

```python
#!/usr/bin/env python3
import json
import subprocess
import sys

def validate_with_preset(fastq_file, preset_id, min_score=75):
    """Extract and validate FASTQ metadata with preset."""

    # Extract metadata
    result = subprocess.run([
        "cicada", "metadata", "extract", fastq_file,
        "--preset", preset_id,
        "--format", "json"
    ], capture_output=True, text=True)

    if result.returncode != 0:
        print(f"Extraction failed: {result.stderr}", file=sys.stderr)
        return False

    # Parse metadata
    metadata = json.loads(result.stdout)
    score = metadata.get("quality_score", 0)

    # Check score
    if score < min_score:
        print(f"Quality score {score} below threshold {min_score}")
        print(f"Missing fields: {metadata.get('validation_warnings', [])}")
        return False

    print(f"✅ Validation passed (score: {score})")
    return True

# Usage
if __name__ == "__main__":
    success = validate_with_preset(
        fastq_file="sample_R1.fastq.gz",
        preset_id="illumina-novaseq",
        min_score=75
    )
    sys.exit(0 if success else 1)
```

### Bash Script

```bash
#!/bin/bash
# Batch validation with presets

PRESET="illumina-novaseq"
MIN_SCORE=75
OUTPUT_DIR="validated_metadata"

mkdir -p "$OUTPUT_DIR"

for fastq in data/*.fastq.gz; do
    basename=$(basename "$fastq" .fastq.gz)
    echo "Processing $basename..."

    # Extract with preset
    cicada metadata extract "$fastq" \
      --preset "$PRESET" \
      --format json \
      --output "$OUTPUT_DIR/${basename}.json"

    # Check quality score
    score=$(jq -r '.quality_score' "$OUTPUT_DIR/${basename}.json")

    if (( $(echo "$score < $MIN_SCORE" | bc -l) )); then
        echo "❌ $basename: Score $score below threshold"
        # Log warnings
        jq -r '.validation_warnings[]' "$OUTPUT_DIR/${basename}.json" \
          >> "$OUTPUT_DIR/warnings.log"
    else
        echo "✅ $basename: Score $score"
    fi
done

echo "Validation complete. Check $OUTPUT_DIR/warnings.log for issues."
```

## Troubleshooting

### "Preset not found"

**Error:**
```
Error: preset 'illumina-nova' not found
```

**Solution:**
```bash
# List available presets
cicada metadata preset list

# Use exact preset ID
cicada metadata validate data.json --preset illumina-novaseq
```

### Low Quality Score

**Issue:** Quality score below 75

**Diagnosis:**
```bash
# Show validation details
cicada metadata validate metadata.json --preset illumina-novaseq
```

**Solutions:**

1. **Add missing required fields** (if score < 60):
   ```bash
   # Check which required fields are missing
   cicada metadata validate metadata.json --preset illumina-novaseq

   # Add fields via enrichment file
   echo "instrument_model: NovaSeq 6000" > enrich.yaml
   cicada metadata validate metadata.json \
     --preset illumina-novaseq \
     --enrich enrich.yaml
   ```

2. **Add optional fields** (if score 60-89):
   ```yaml
   # enrich.yaml
   run_id: RUN_12345
   flowcell_id: FC_67890
   lane: 2
   barcode: AGTCACTA
   mean_quality_score: 36.5
   ```

### Field Type Mismatch

**Error:**
```
Validation Error: Field 'total_reads' expected type 'integer', got 'string'
```

**Solution:**

Check metadata format:
```bash
# View metadata
cat metadata.json | jq '.total_reads'

# Should be: 1000000 (number)
# Not: "1000000" (string)
```

Fix in enrichment file:
```yaml
# Correct (YAML auto-converts)
total_reads: 1000000

# Incorrect
total_reads: "1000000"
```

### Strict Validation Failures

**Issue:** Validation fails in strict mode but passes in lenient mode

**Command:**
```bash
# Fails
cicada metadata validate metadata.json --preset illumina-novaseq --strict

# Passes
cicada metadata validate metadata.json --preset illumina-novaseq
```

**Solution:**

Strict mode treats warnings as errors. Add missing optional fields:
```bash
# See which optional fields are missing
cicada metadata validate metadata.json --preset illumina-novaseq

# Add recommended fields
cat > enrich.yaml <<EOF
run_id: RUN_123
quality_encoding: Phred+33
is_paired_end: true
EOF

# Re-validate with enrichment
cicada metadata validate metadata.json \
  --preset illumina-novaseq \
  --enrich enrich.yaml \
  --strict
```

### Wrong Preset Selected

**Issue:** Validation fails because wrong preset was used

**Diagnosis:**
```bash
# File is from MiSeq, but validated with NovaSeq preset
cicada metadata validate metadata.json --preset illumina-novaseq
# Error: instrument_model "MiSeq" doesn't match preset "NovaSeq 6000"
```

**Solution:**
```bash
# Use correct preset
cicada metadata validate metadata.json --preset illumina-miseq

# Or use generic preset
cicada metadata validate metadata.json --preset generic-sequencing
```

## Best Practices

### 1. Use Specific Presets When Possible

✅ **Recommended:**
```bash
cicada metadata extract data.fastq --preset illumina-novaseq
```

❌ **Avoid (unless necessary):**
```bash
cicada metadata extract data.fastq --preset generic-sequencing
```

Specific presets provide better validation and higher quality scores.

### 2. Validate Early

Validate metadata immediately after extraction:
```bash
# Good workflow
cicada metadata extract sample.fastq --preset illumina-novaseq
# Reviews validation results immediately

# Bad workflow
cicada metadata extract sample.fastq
# ... weeks later ...
cicada doi prepare sample.fastq  # Validation fails at DOI preparation
```

### 3. Target 80+ Quality Scores

Aim for quality scores above 80 for publication:
- **90-100**: Excellent - include all optional fields
- **80-89**: Good - acceptable for most journals
- **70-79**: Adequate - may need additional fields
- **< 70**: Insufficient - add more metadata

### 4. Maintain Consistency

Use the same preset for all files from the same instrument:
```bash
# Define preset once
PRESET="illumina-novaseq"

# Use consistently
for file in *.fastq.gz; do
    cicada metadata extract "$file" --preset "$PRESET"
done
```

### 5. Document Preset Choice

Include preset information in project documentation:
```markdown
## Metadata Standards

- Instrument: Illumina NovaSeq 6000
- Preset: `illumina-novaseq`
- Minimum Quality Score: 80
- Validation: All files validated before DOI preparation
```

### 6. Use Enrichment Files

Create project-wide enrichment files:
```yaml
# project_metadata.yaml
run_id: PROJECT_2025_001
quality_encoding: Phred+33
instrument_manufacturer: Illumina
instrument_model: NovaSeq 6000
sequencing_platform: Illumina
```

Apply to all files:
```bash
cicada metadata extract sample.fastq \
  --preset illumina-novaseq \
  --enrich project_metadata.yaml
```

## Future Enhancements (v0.3.0+)

Planned features for future releases:

### Custom Presets

Create your own presets:
```bash
# Create preset from template
cicada metadata preset create \
  --name custom-nanopore \
  --type sequencing \
  --manufacturer "Oxford Nanopore" \
  --template nanopore-template.yaml

# Use custom preset
cicada metadata validate data.fastq --preset custom-nanopore
```

### Preset Templates

Generate metadata templates:
```bash
# Generate template from preset
cicada metadata preset template illumina-novaseq > template.yaml

# Fill in template
vim template.yaml

# Validate against preset
cicada metadata validate template.yaml --preset illumina-novaseq
```

### Preset Validation Rules

Define custom validation rules:
```yaml
# custom-preset.yaml
name: custom-miseq
manufacturer: Illumina
instrument_model: MiSeq
required_fields:
  - format
  - total_reads
optional_fields:
  - run_id
  - barcode
validation_rules:
  read_length:
    min: 50
    max: 300
  quality_score:
    min: 20
```

### Preset Import/Export

Share presets across teams:
```bash
# Export preset
cicada metadata preset export illumina-novaseq > preset.yaml

# Import preset
cicada metadata preset import preset.yaml

# Share with team
git add presets/illumina-novaseq.yaml
git commit -m "Add NovaSeq preset for lab"
```

## Related Documentation

- **[Metadata Extraction Guide](METADATA_EXTRACTION.md)**: Extracting metadata from files
- **[DOI Workflow Guide](DOI_WORKFLOW.md)**: Preparing metadata for DOI registration
- **[User Scenarios](USER_SCENARIOS_v0.2.0.md)**: Real-world preset usage examples
- **[Integration Testing](../INTEGRATION_TESTING.md)**: Testing preset validation

## Support

For questions or issues with presets:

- **Documentation**: See guides above
- **Issues**: Report problems at https://github.com/scttfrdmn/cicada/issues
- **Feature Requests**: Suggest new presets or improvements

## Version History

- **v0.2.0** (Current): Initial release with 8 default presets
- **v0.3.0** (Planned): Custom presets, templates, advanced validation rules
