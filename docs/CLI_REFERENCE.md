# Cicada CLI Reference

**Last Updated:** 2025-11-25

Complete reference for all Cicada command-line interface commands, flags, and options.

## Table of Contents

1. [Global Flags](#global-flags)
2. [Command Overview](#command-overview)
3. [Storage & Sync Commands](#storage--sync-commands)
4. [Metadata Commands](#metadata-commands)
5. [Watch Commands](#watch-commands)
6. [DOI Commands](#doi-commands-optional)
7. [Configuration Commands](#configuration-commands)
8. [Version Command](#version-command)
9. [Exit Codes](#exit-codes)
10. [Environment Variables](#environment-variables)
11. [Examples](#examples)

---

## Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | | Path to configuration file | `~/.cicada/config.yaml` |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--help` | `-h` | Show help for command | |

**Examples:**
```bash
# Use custom config file
cicada --config /path/to/config.yaml sync /data s3://bucket/data

# Enable verbose output for any command
cicada --verbose metadata extract file.czi

# Show help for any command
cicada sync --help
```

---

## Command Overview

```
cicada
├── sync                 - Sync files between storage locations
├── watch                - Watch directories for changes
│   ├── add              - Add a directory to watch
│   ├── list             - List active watches
│   └── remove           - Remove a watch
├── metadata             - Metadata operations
│   ├── extract          - Extract metadata from files
│   ├── show             - Display metadata
│   ├── validate         - Validate metadata against presets
│   ├── list             - List available extractors
│   └── preset           - Manage presets
│       ├── list         - List available presets
│       ├── show         - Show preset details
│       └── validate     - Validate file against preset
├── doi                  - DOI preparation (optional)
│   ├── prepare          - Prepare DOI metadata
│   ├── validate         - Validate DOI metadata
│   └── submit           - Submit to provider
├── config               - Configuration management
│   ├── show             - Show configuration
│   ├── set              - Set configuration value
│   └── init             - Initialize configuration
└── version              - Show version information
```

---

## Storage & Sync Commands

### `cicada sync`

Synchronize files between local filesystem and cloud storage (S3).

**Usage:**
```bash
cicada sync [flags] SOURCE DESTINATION
```

**Arguments:**
- `SOURCE` - Source location (local path or S3 URI)
- `DESTINATION` - Destination location (local path or S3 URI)

**S3 URI Format:**
```
s3://bucket-name/prefix/path
s3://bucket-name/           # entire bucket
```

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--dry-run` | bool | Show what would be synced without making changes | `false` |
| `--delete` | bool | Delete files in destination not present in source | `false` |
| `--concurrency` | int | Number of parallel transfers | `4` |

**Sync Behavior:**
- Files are compared using checksums (ETag for S3, MD5 for local)
- Only changed or new files are transferred
- Directories are synced recursively
- Metadata is preserved during transfer

**Examples:**

```bash
# Sync local directory to S3
cicada sync /data/lab s3://my-bucket/lab-data

# Sync from S3 to local
cicada sync s3://my-bucket/lab-data /data/lab

# Dry run to preview changes
cicada sync --dry-run /data/lab s3://my-bucket/lab-data

# Sync and delete files not in source (careful!)
cicada sync --delete /data/lab s3://my-bucket/lab-data

# Sync with increased concurrency for faster transfer
cicada sync --concurrency 8 /data/large s3://my-bucket/large

# Bidirectional sync (run both commands)
cicada sync /local/data s3://bucket/data  # upload changes
cicada sync s3://bucket/data /local/data  # download changes
```

**Exit Codes:**
- `0` - Success
- `1` - General error (network, permissions, etc.)
- `2` - Invalid arguments

**Common Errors:**

| Error | Cause | Solution |
|-------|-------|----------|
| `no such file or directory` | Source path doesn't exist | Verify source path exists |
| `access denied` | Missing S3 permissions | Check AWS credentials and IAM permissions |
| `bucket not found` | S3 bucket doesn't exist | Verify bucket name and region |
| `connection timeout` | Network issue | Check internet connectivity |

---

## Metadata Commands

### `cicada metadata extract`

Extract metadata from scientific instrument files.

**Usage:**
```bash
cicada metadata extract [flags] PATH
```

**Arguments:**
- `PATH` - Path to file to extract metadata from

**Flags:**

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--format` | `-f` | string | Output format: json, yaml, table | `json` |
| `--output` | `-o` | string | Output file (default: stdout) | |
| `--extractor` | | string | Force specific extractor | auto-detect |

**Supported File Formats:**

| Format | Extensions | Extractor | Domain |
|--------|-----------|-----------|---------|
| TIFF | `.tif`, `.tiff` | TIFF | Microscopy |
| OME-TIFF | `.ome.tif`, `.ome.tiff` | OME-TIFF | Microscopy |
| Zeiss CZI | `.czi` | Zeiss CZI | Microscopy |
| Nikon ND2 | `.nd2` | Nikon ND2 | Microscopy |
| Leica LIF | `.lif` | Leica LIF | Microscopy |
| FASTQ | `.fastq`, `.fq`, `.fastq.gz` | FASTQ | Sequencing |
| BAM | `.bam` | BAM | Sequencing |
| mzML | `.mzml` | mzML | Mass Spectrometry |
| MGF | `.mgf` | MGF | Mass Spectrometry |
| HDF5 | `.h5`, `.hdf5` | HDF5 | Data Arrays |
| Zarr | `.zarr` | Zarr | Data Arrays |
| DICOM | `.dcm`, `.dicom` | DICOM | Medical Imaging |
| FCS | `.fcs` | FCS | Flow Cytometry |
| Generic | `*` | Generic | Fallback |

**Examples:**

```bash
# Extract metadata as JSON (default)
cicada metadata extract data/experiment001.czi

# Extract and save to file
cicada metadata extract data/experiment001.czi --output metadata.json

# Extract as YAML
cicada metadata extract data/experiment001.czi --format yaml

# Extract in human-readable table format
cicada metadata extract data/experiment001.czi --format table

# Force specific extractor
cicada metadata extract data/image.tif --extractor OME-TIFF

# Extract from multiple files
for file in data/*.czi; do
  cicada metadata extract "$file" --output "${file}.metadata.json"
done
```

**Output Example (JSON):**
```json
{
  "format": "CZI",
  "microscope_manufacturer": "Zeiss",
  "microscope_model": "LSM 880",
  "modality": "confocal",
  "width": 2048,
  "height": 2048,
  "channels": 3,
  "pixel_size_x": 0.125,
  "pixel_size_y": 0.125,
  "pixel_size_z": 0.5,
  "voxel_size_unit": "micrometers",
  "channel_info": [
    {
      "name": "DAPI",
      "index": 0,
      "excitation_wavelength": 405,
      "emission_wavelength": 450
    }
  ],
  "acquisition_date": "2025-01-24T10:30:00Z",
  "operator": "jsmith"
}
```

---

### `cicada metadata show`

Display metadata in human-readable format with better formatting than `extract`.

**Usage:**
```bash
cicada metadata show [flags] PATH
```

**Arguments:**
- `PATH` - Path to file

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--format` | string | Display format: table, json, yaml | `table` |

**Examples:**

```bash
# Show metadata in table format
cicada metadata show data/experiment001.czi

# Show as JSON (similar to extract)
cicada metadata show data/experiment001.czi --format json
```

---

### `cicada metadata validate`

Validate metadata against an instrument preset to ensure data quality.

**Usage:**
```bash
cicada metadata validate [flags] PATH
```

**Arguments:**
- `PATH` - Path to file to validate

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--preset` | string | Preset name to validate against | (required) |
| `--strict` | bool | Treat warnings as errors | `false` |

**Built-in Presets:**

| Preset | Description | Domain |
|--------|-------------|--------|
| `microscopy-confocal` | Confocal microscopy requirements | Microscopy |
| `microscopy-widefield` | Widefield microscopy requirements | Microscopy |
| `sequencing-illumina` | Illumina sequencing requirements | Sequencing |
| `sequencing-pacbio` | PacBio sequencing requirements | Sequencing |
| `mass-spec-proteomics` | Proteomics mass spec requirements | Mass Spec |
| `mass-spec-metabolomics` | Metabolomics mass spec requirements | Mass Spec |
| `flow-cytometry` | Flow cytometry requirements | Flow Cytometry |
| `general-lab` | Basic lab data requirements | General |

**Examples:**

```bash
# Validate confocal microscopy data
cicada metadata validate data/confocal.czi --preset microscopy-confocal

# Validate with strict mode (warnings are errors)
cicada metadata validate data/image.czi --preset microscopy-confocal --strict

# Validate sequencing data
cicada metadata validate data/sample1.fastq --preset sequencing-illumina
```

**Validation Output:**
```
Validating: data/confocal.czi
Preset: microscopy-confocal

✓ Required fields present (12/12)
⚠ Recommended fields missing (2/8):
  - temperature
  - co2_level

✓ Field values valid

Validation: PASSED (2 warnings)
```

**Exit Codes:**
- `0` - Validation passed
- `1` - Validation failed (missing required fields or invalid values)
- `2` - Validation passed with warnings (strict mode: exit 1)

---

### `cicada metadata list`

List all available metadata extractors.

**Usage:**
```bash
cicada metadata list
```

**Example Output:**
```
Available Metadata Extractors:

TIFF              .tif, .tiff
OME-TIFF          .ome.tif, .ome.tiff
Zeiss CZI         .czi
Nikon ND2         .nd2
Leica LIF         .lif
FASTQ             .fastq, .fq, .fastq.gz
BAM               .bam
mzML              .mzml
MGF               .mgf
HDF5              .h5, .hdf5
Zarr              .zarr
DICOM             .dcm, .dicom
FCS               .fcs
Generic           * (fallback)

Total: 14 extractors
```

---

### `cicada metadata preset`

Manage validation presets.

**Subcommands:**
- `list` - List available presets
- `show` - Show preset details
- `validate` - Validate file against preset (alias for `metadata validate`)

#### `cicada metadata preset list`

List all available presets.

**Usage:**
```bash
cicada metadata preset list
```

**Example Output:**
```
Available Presets:

Microscopy:
  microscopy-confocal         Confocal microscopy requirements
  microscopy-widefield        Widefield microscopy requirements

Sequencing:
  sequencing-illumina         Illumina sequencing requirements
  sequencing-pacbio           PacBio sequencing requirements

Mass Spectrometry:
  mass-spec-proteomics        Proteomics mass spec requirements
  mass-spec-metabolomics      Metabolomics mass spec requirements

Flow Cytometry:
  flow-cytometry              Flow cytometry requirements

General:
  general-lab                 Basic lab data requirements

Total: 8 presets
```

#### `cicada metadata preset show`

Show detailed information about a preset.

**Usage:**
```bash
cicada metadata preset show PRESET_NAME
```

**Example:**
```bash
cicada metadata preset show microscopy-confocal
```

**Example Output:**
```
Preset: microscopy-confocal
Description: Confocal microscopy requirements

Required Fields (12):
  - manufacturer
  - model
  - modality
  - width
  - height
  - channels
  - pixel_size_x
  - pixel_size_y
  - magnification
  - numerical_aperture
  - acquisition_date
  - operator

Recommended Fields (8):
  - sample_id
  - experiment_name
  - organism
  - tissue
  - temperature
  - co2_level
  - objective
  - detector_model

Allowed Values:
  modality: confocal, spinning-disk, two-photon
```

---

## Watch Commands

### `cicada watch add`

Start watching a directory for changes and automatically sync to destination.

**Usage:**
```bash
cicada watch add [flags] SOURCE DESTINATION
```

**Arguments:**
- `SOURCE` - Directory to watch (local path)
- `DESTINATION` - Sync destination (local path or S3 URI)

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--debounce` | int | Debounce delay in seconds | `5` |
| `--min-age` | int | Minimum file age before sync (seconds) | `10` |
| `--delete-source` | bool | Delete source files after successful sync | `false` |
| `--sync-on-start` | bool | Perform initial sync when starting | `true` |

**How Watch Works:**

1. **File System Monitoring**: Watches directory for file creation, modification, deletion
2. **Debouncing**: Waits for quiet period (no changes for `debounce` seconds)
3. **Age Check**: Only syncs files older than `min-age` seconds (prevents syncing incomplete writes)
4. **Sync**: Triggers sync engine to transfer files
5. **Optional Cleanup**: Deletes source files if `delete-source` is true

**Examples:**

```bash
# Basic watch (sync to S3)
cicada watch add /data/microscope s3://my-bucket/microscope-data

# Watch with longer debounce (wait 30 seconds)
cicada watch add /data/output s3://bucket/data --debounce 30

# Watch with file deletion after sync (use carefully!)
cicada watch add /data/temp s3://bucket/archive --delete-source

# Watch without initial sync
cicada watch add /data/new s3://bucket/new --sync-on-start=false

# Watch with increased minimum age (wait for files to stabilize)
cicada watch add /data/large-files s3://bucket/files --min-age 60
```

**Watch Configuration:**
Watches are saved to `~/.cicada/config.yaml` and automatically restored on restart.

---

### `cicada watch list`

List all active watches.

**Usage:**
```bash
cicada watch list
```

**Example Output:**
```
Active watches: 2

Watch: /data/microscope-1704067200
  Source: /data/microscope
  Destination: s3://my-bucket/microscope-data
  Active: true
  Started: 2025-01-24T10:00:00Z
  Last sync: 2025-01-24T10:15:23Z
  Files synced: 45
  Bytes synced: 1048576000
  Errors: 0

Watch: /data/sequencer-1704067800
  Source: /data/sequencer
  Destination: s3://my-bucket/sequencing-data
  Active: true
  Started: 2025-01-24T10:10:00Z
  Last sync: 2025-01-24T10:20:05Z
  Files synced: 12
  Bytes synced: 524288000
  Errors: 2
  Last error: failed to sync file: connection timeout
```

---

### `cicada watch remove`

Stop and remove a watch.

**Usage:**
```bash
cicada watch remove WATCH_ID
```

**Arguments:**
- `WATCH_ID` - Watch ID from `watch list` command

**Examples:**

```bash
# Remove a specific watch
cicada watch remove /data/microscope-1704067200

# Remove all watches (requires manual deletion from config)
# Edit ~/.cicada/config.yaml and remove watch entries
```

---

## DOI Commands (Optional)

> **Note:** DOI preparation is an optional advanced feature for labs that need to publish datasets. Most Cicada usage involves core data management features (storage, sync, metadata extraction).

### `cicada doi prepare`

Prepare DOI metadata for dataset publication.

**Usage:**
```bash
cicada doi prepare [flags] PATH
```

**Arguments:**
- `PATH` - Path to dataset directory

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--title` | string | Dataset title | (required) |
| `--creator` | string | Creator name(s), comma-separated | (required) |
| `--publisher` | string | Publisher name | (required) |
| `--year` | int | Publication year | current year |
| `--resource-type` | string | Resource type (Dataset, Software, etc.) | `Dataset` |
| `--output` | string | Output file for DOI metadata | `doi-metadata.json` |

**Examples:**

```bash
# Prepare DOI metadata
cicada doi prepare /data/experiment-001 \
  --title "Confocal Microscopy of Neural Tissue" \
  --creator "Smith, John; Doe, Jane" \
  --publisher "University Lab" \
  --year 2025

# Specify output file
cicada doi prepare /data/experiment-001 \
  --title "RNA-Seq Analysis" \
  --creator "Jones, Alice" \
  --publisher "Research Institute" \
  --output metadata/doi.json
```

---

### `cicada doi validate`

Validate DOI metadata against DataCite schema.

**Usage:**
```bash
cicada doi validate [flags] PATH
```

**Arguments:**
- `PATH` - Path to DOI metadata file

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--schema` | string | DataCite schema version | `4.5` |
| `--strict` | bool | Strict validation mode | `false` |

**Examples:**

```bash
# Validate DOI metadata
cicada doi validate doi-metadata.json

# Strict validation
cicada doi validate doi-metadata.json --strict
```

---

### `cicada doi submit`

Submit DOI to provider (DataCite, Zenodo).

**Usage:**
```bash
cicada doi submit [flags] PATH
```

**Arguments:**
- `PATH` - Path to DOI metadata file

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--provider` | string | Provider: datacite, zenodo | (required) |
| `--prefix` | string | DOI prefix (for DataCite) | from config |
| `--test` | bool | Use test/sandbox environment | `false` |

**Examples:**

```bash
# Submit to DataCite (requires configuration)
cicada doi submit doi-metadata.json --provider datacite

# Submit to Zenodo sandbox for testing
cicada doi submit doi-metadata.json --provider zenodo --test

# Submit with custom DOI prefix
cicada doi submit doi-metadata.json --provider datacite --prefix 10.5555
```

---

## Configuration Commands

### `cicada config show`

Display current configuration.

**Usage:**
```bash
cicada config show [flags]
```

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--format` | string | Output format: yaml, json, table | `yaml` |

**Examples:**

```bash
# Show configuration as YAML
cicada config show

# Show as JSON
cicada config show --format json

# Show as table
cicada config show --format table
```

---

### `cicada config set`

Set configuration value.

**Usage:**
```bash
cicada config set KEY VALUE
```

**Arguments:**
- `KEY` - Configuration key (dot notation for nested)
- `VALUE` - Value to set

**Common Keys:**

| Key | Description | Example Value |
|-----|-------------|---------------|
| `aws.profile` | AWS profile name | `default` |
| `aws.region` | AWS region | `us-west-2` |
| `sync.concurrency` | Default sync concurrency | `4` |
| `sync.delete` | Default delete behavior | `false` |
| `settings.verbose` | Verbose logging | `true` |
| `settings.log_file` | Log file path | `~/.cicada/cicada.log` |

**Examples:**

```bash
# Set AWS profile
cicada config set aws.profile research-account

# Set AWS region
cicada config set aws.region us-east-1

# Set default concurrency
cicada config set sync.concurrency 8

# Enable verbose logging
cicada config set settings.verbose true
```

---

### `cicada config init`

Initialize configuration file with defaults.

**Usage:**
```bash
cicada config init [flags]
```

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--force` | bool | Overwrite existing configuration | `false` |

**Examples:**

```bash
# Initialize configuration
cicada config init

# Force overwrite existing config
cicada config init --force
```

**Created File:** `~/.cicada/config.yaml`

**Default Configuration:**
```yaml
version: "1"

aws:
  profile: default
  region: ""

sync:
  concurrency: 4
  delete: false
  exclude:
    - .git/**
    - .DS_Store
    - "*.tmp"
    - "*.swp"

watches: []

settings:
  verbose: false
  log_file: ""
  check_updates: true
```

---

## Version Command

### `cicada version`

Display version information.

**Usage:**
```bash
cicada version
```

**Example Output:**
```
Cicada v0.2.0
Commit: a1b2c3d
Built: 2025-01-24T10:00:00Z
Built by: goreleaser
Go version: go1.21.5
```

---

## Exit Codes

Cicada uses standard exit codes:

| Code | Meaning | Examples |
|------|---------|----------|
| `0` | Success | Command completed successfully |
| `1` | General error | Network error, file not found, permission denied |
| `2` | Invalid arguments | Wrong number of arguments, invalid flag values |
| `3` | Validation failed | Metadata validation failed (strict mode) |

**Using Exit Codes in Scripts:**
```bash
# Check if sync succeeded
if cicada sync /data s3://bucket/data; then
  echo "Sync successful"
else
  echo "Sync failed with code $?"
  exit 1
fi

# Validate and fail script if validation fails
cicada metadata validate file.czi --preset microscopy-confocal --strict || exit 1
```

---

## Environment Variables

Cicada respects these environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `CICADA_CONFIG` | Path to configuration file | `/etc/cicada/config.yaml` |
| `AWS_PROFILE` | AWS profile to use | `research-account` |
| `AWS_REGION` | AWS region | `us-west-2` |
| `AWS_ACCESS_KEY_ID` | AWS access key | (from AWS credentials) |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | (from AWS credentials) |
| `HOME` | Home directory for config | `/home/user` |

**Precedence (highest to lowest):**
1. Command-line flags
2. Environment variables
3. Configuration file (`~/.cicada/config.yaml`)
4. Built-in defaults

**Examples:**

```bash
# Use custom config file
export CICADA_CONFIG=/path/to/config.yaml
cicada sync /data s3://bucket/data

# Use specific AWS profile
export AWS_PROFILE=research-account
cicada sync /data s3://bucket/data

# Override AWS region
export AWS_REGION=eu-west-1
cicada sync /data s3://bucket/data
```

---

## Examples

### Common Workflows

#### Daily Data Backup
```bash
# Sync today's data to S3
cicada sync /data/$(date +%Y-%m-%d) s3://lab-backups/$(date +%Y-%m-%d)

# Extract and validate metadata
for file in /data/$(date +%Y-%m-%d)/*.czi; do
  cicada metadata extract "$file" --output "$file.metadata.json"
  cicada metadata validate "$file" --preset microscopy-confocal
done
```

#### Automated Watch Setup
```bash
# Set up watches for multiple instruments
cicada watch add /data/confocal s3://lab-data/confocal --debounce 10
cicada watch add /data/sequencer s3://lab-data/sequencing --debounce 30
cicada watch add /data/mass-spec s3://lab-data/mass-spec --debounce 15

# List watches to verify
cicada watch list
```

#### Batch Metadata Extraction
```bash
# Extract metadata from all CZI files
find /data -name "*.czi" -exec cicada metadata extract {} --output {}.metadata.json \;

# Extract and validate
for file in /data/*.czi; do
  cicada metadata extract "$file" | \
    cicada metadata validate --preset microscopy-confocal
done
```

#### Data Quality Check
```bash
# Validate all files in a directory
for file in /data/experiment-001/*; do
  echo "Validating $file"
  if cicada metadata validate "$file" --preset general-lab --strict; then
    echo "✓ $file passed"
  else
    echo "✗ $file failed"
  fi
done
```

#### Sync with Progress Monitoring
```bash
# Sync with verbose output and logging
cicada --verbose sync /data/large-dataset s3://bucket/dataset 2>&1 | \
  tee sync-log-$(date +%Y%m%d-%H%M%S).log
```

---

## Related Documentation

- [User Guide](USER_GUIDE.md) - Getting started and workflows
- [Configuration Guide](CONFIGURATION.md) - Configuration options
- [Architecture](ARCHITECTURE.md) - System architecture
- [Metadata System](METADATA_SYSTEM.md) - Metadata details

---

**Contributing:** Found an error or want to suggest improvements? See [CONTRIBUTING.md](CONTRIBUTING.md)

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
