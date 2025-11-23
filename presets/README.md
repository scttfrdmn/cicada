# Cicada Instrument Presets

Pre-configured settings for common laboratory instruments to simplify Cicada setup and ensure best practices.

## Available Presets

### Microscopy

- **[zeiss-confocal.yaml](microscopy/zeiss-confocal.yaml)** - Zeiss LSM confocal microscopes (LSM 880, 900, 980)
  - File format: CZI
  - Typical file size: 500MB - 5GB
  - Debounce: 30s, Min-age: 60s

### Sequencing

- **[illumina-novaseq.yaml](sequencing/illumina-novaseq.yaml)** - Illumina NovaSeq 6000
  - File format: FASTQ/FASTQ.GZ
  - Typical run size: 500GB - 1.5TB
  - Debounce: 60s, Min-age: 300s (5 min)

## Using Presets

### Method 1: Interactive Setup (Recommended)

```bash
# Launch interactive wizard
cicada instrument setup

# Wizard will:
# 1. Ask for instrument type (microscopy, sequencing, etc.)
# 2. Show available presets
# 3. Configure paths and destination
# 4. Set up watch automatically
```

### Method 2: Direct Preset Application

```bash
# Apply specific preset
cicada instrument setup zeiss-confocal \
  --path /mnt/zeiss/output \
  --destination s3://lab-data/microscopy

# This automatically:
# - Configures debounce and min-age
# - Enables metadata extraction
# - Sets up file validation
# - Creates watch with optimized settings
```

### Method 3: Auto-Detection

```bash
# Let Cicada detect the instrument type
cicada instrument detect /mnt/zeiss/output

# Output:
# Detected: Zeiss Confocal Microscope (confidence: high)
# Preset: zeiss-confocal
#
# Apply this preset? [y/N]: y
```

### Method 4: Manual Configuration with Preset

```bash
# Use preset with watch command
cicada watch add \
  --preset zeiss-confocal \
  /mnt/zeiss/output \
  s3://lab-data/microscopy
```

## Preset Structure

Each preset YAML file contains:

```yaml
# Identity
name: "Human-readable instrument name"
id: "unique-preset-id"
version: "1.0"
category: microscopy|sequencing|mass-spec|flow-cytometry

# Detection rules
detection:
  file_extensions: [".czi", ".tif"]
  magic_bytes:
    offset: 0
    signature: "ZISRAWFILE"

# Sync configuration
sync:
  debounce_seconds: 30
  min_age_seconds: 60
  concurrency: 4
  exclude_patterns: []

# Metadata extraction
metadata:
  enabled: true
  extractor: "zeiss-czi"
  storage:
    s3_tags: true
    sidecar_json: true
    catalog: true

# Validation
validation:
  enabled: true
  checks:
    - type: file_integrity
    - type: minimum_file_size
      value: 1048576

# Additional settings
s3:
  storage_class: STANDARD
  versioning: true

tags:
  instrument_type: microscopy
  manufacturer: zeiss
```

## Creating Custom Presets

### Option 1: Copy and Modify

```bash
# Copy an existing preset
cp presets/microscopy/zeiss-confocal.yaml \
   presets/microscopy/my-custom-preset.yaml

# Edit to customize
nano presets/microscopy/my-custom-preset.yaml

# Use your custom preset
cicada instrument setup my-custom-preset \
  --path /data/instrument \
  --destination s3://lab-data/
```

### Option 2: Export from Existing Watch

```bash
# Set up a watch manually
cicada watch add \
  --debounce 45 \
  --min-age 90 \
  /data/custom-instrument \
  s3://lab-data/custom

# Export configuration as preset
cicada instrument export my-instrument-preset \
  --from-watch /data/custom-instrument \
  --output presets/custom/my-instrument.yaml
```

### Option 3: Interactive Creation

```bash
# Launch preset creator
cicada instrument create

# Wizard will ask:
# - Instrument name and type
# - File extensions and patterns
# - Sync settings (debounce, min-age, concurrency)
# - Metadata extraction preferences
# - Validation rules
# - Output location
```

## Preset Guidelines

When creating custom presets, follow these guidelines:

### Debounce Settings

- **Small files (< 10MB)**: 5-15 seconds
- **Medium files (10-100MB)**: 15-30 seconds
- **Large files (100MB-1GB)**: 30-60 seconds
- **Very large files (> 1GB)**: 60-120 seconds

### Min-Age Settings

- **Fast-writing instruments**: 30-60 seconds
- **Slow-writing instruments** (compression, large files): 60-300 seconds
- **Network-mounted storage**: Add 60+ seconds buffer

### Concurrency Settings

- **Many small files**: 8-16 concurrent transfers
- **Few large files**: 2-4 concurrent transfers
- **Mixed workload**: 4-8 concurrent transfers
- **Limited bandwidth**: 2-4 concurrent transfers

### Validation

Always enable validation for:
- Corrupt file detection
- Incomplete write detection
- Format compliance

### S3 Storage Classes

- **STANDARD**: Active data, frequent access
- **STANDARD_IA**: Infrequent access, cost optimization
- **GLACIER**: Long-term archival, lowest cost
- **INTELLIGENT_TIERING**: Automatic optimization

## Preset Testing

Test your preset before production use:

```bash
# Test with dry-run
cicada watch add \
  --preset my-custom-preset \
  --dry-run \
  /test/data \
  s3://test-bucket/test

# Monitor for one hour
cicada watch list

# Check logs for any issues
cicada watch status /test/data-123456

# If successful, apply to production
cicada instrument setup my-custom-preset \
  --path /production/data \
  --destination s3://prod-bucket/data
```

## Contributing Presets

Have a preset for a common lab instrument? Contribute it!

1. Create preset following the structure above
2. Test thoroughly with real instrument data
3. Document any instrument-specific quirks
4. Submit pull request to: https://github.com/scttfrdmn/cicada

Popular instruments we'd love presets for:
- Nikon Ti2/A1 confocal microscopes
- Leica SP8 confocal
- Olympus FV3000
- Illumina MiSeq/NextSeq
- Oxford Nanopore MinION
- Thermo Orbitrap mass spec
- BD FACSAria flow cytometer
- And many more!

## Preset Versioning

Presets follow semantic versioning:

- **Major version** (1.0 → 2.0): Breaking changes, incompatible settings
- **Minor version** (1.0 → 1.1): New features, backward compatible
- **Patch version** (1.0.0 → 1.0.1): Bug fixes, minor improvements

## Support

Questions about presets?

- GitHub Issues: https://github.com/scttfrdmn/cicada/issues
- Documentation: https://github.com/scttfrdmn/cicada/docs
- Discussions: https://github.com/scttfrdmn/cicada/discussions
