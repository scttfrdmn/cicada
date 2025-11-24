# Cicada

**Dormant data commons for academic research** - Lightweight, cost-effective platform providing federated storage, access control, and compute-to-data capabilities. Like a cicada, it lies dormant (consuming minimal resources) until needed, then emerges powerfully for data-intensive work.

[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/cicada)](https://goreportcard.com/report/github.com/scttfrdmn/cicada)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/scttfrdmn/cicada)](go.mod)

> **Current Version**: v0.2.0 - Adds comprehensive metadata extraction and DOI preparation capabilities. See **[CHANGELOG.md](CHANGELOG.md)** for release details and **[planning/PROJECT-SUMMARY.md](planning/PROJECT-SUMMARY.md)** for the complete data commons vision.

## Documentation

- 📖 **[User Scenarios v0.2.0](docs/USER_SCENARIOS_v0.2.0.md)** - Detailed persona-based walkthroughs:
  - Postdoc extracting metadata from sequencing data
  - Graduate student preparing datasets for publication with DOI
  - Lab manager validating compliance with instrument presets
  - Data curator enriching metadata for repositories
  - Small lab complete adoption journey (target users)
- 📚 **User Guides**:
  - [Metadata Extraction Guide](docs/METADATA_EXTRACTION.md) - Extract metadata from scientific files
  - [DOI Workflow Guide](docs/DOI_WORKFLOW.md) - Prepare datasets for DOI registration
  - [Instrument Presets Guide](docs/PRESETS.md) - Validate metadata with instrument-specific presets
  - [Provider Setup Guide](docs/PROVIDERS.md) - Configure DataCite/Zenodo integration
- 🚀 **Quick Start** (below) - Get started in 5 minutes
- 📋 **[Full Documentation](#usage)** - Complete command reference

## Features

### Storage & Sync (v0.1.0)
- ✅ **Fast S3 Sync**: Concurrent transfers with MD5/ETag comparison
- ✅ **File Watching**: Auto-sync directories on file changes
- ✅ **Two-way Sync**: Local ↔ S3 in both directions
- ✅ **Smart Transfers**: Only sync changed files
- ✅ **Dry Run Mode**: Preview changes before syncing

### Metadata & DOI (v0.2.0)
- ✅ **Metadata Extraction**: Automatic extraction from FASTQ files (sequencing data)
- ✅ **Instrument Presets**: 8 built-in presets (Illumina, Zeiss, generic)
- ✅ **DOI Preparation**: DataCite Schema v4.5 compliance validation
- ✅ **Quality Scoring**: 0-100 scale for metadata completeness
- ✅ **Provider Integration**: DataCite and Zenodo support (API structure ready)

### Cross-Cutting
- ✅ **Configurable**: YAML configuration with extensible design
- ✅ **Cross-platform**: Linux, macOS, Windows
- ✅ **Performant**: < 1 ms metadata extraction, 4-8x concurrent speedup

## Quick Start

### Installation

```bash
# From source (requires Go 1.23+)
git clone https://github.com/scttfrdmn/cicada.git
cd cicada
make install

# Or download pre-built binary from releases
```

### Basic Usage

#### Storage & Sync (v0.1.0)

```bash
# Configure AWS credentials
aws configure

# One-time sync: local to S3
cicada sync /local/data s3://my-bucket/data

# One-time sync: S3 to local
cicada sync s3://my-bucket/data /local/data

# Watch directory and auto-sync to S3
cicada watch add /data/microscope s3://my-bucket/microscope
cicada watch list

# Preview changes without syncing
cicada sync --dry-run /local/data s3://my-bucket/data

# Delete files in destination not in source
cicada sync --delete /local/data s3://my-bucket/data
```

#### Metadata Extraction (v0.2.0)

```bash
# Extract metadata from FASTQ file
cicada metadata extract sample_R1.fastq.gz

# Extract with preset validation
cicada metadata extract sample_R1.fastq.gz --preset illumina-novaseq

# Save metadata to file
cicada metadata extract sample_R1.fastq.gz --format json --output metadata.json

# List available instrument presets
cicada metadata preset list

# Show preset details
cicada metadata preset show illumina-novaseq
```

#### DOI Preparation (v0.2.0)

```bash
# Validate metadata for DOI readiness
cicada doi validate sample.fastq

# Prepare dataset for DOI registration with enrichment
cicada doi prepare sample_R1.fastq sample_R2.fastq \
  --enrich metadata.yaml \
  --publisher "My Lab" \
  --output doi-ready.json

# Check quality score
jq '.validation.score' doi-ready.json
```

## Installation

### Requirements

- Go 1.23+ (for building from source)
- AWS credentials configured (`~/.aws/credentials` or environment variables)
- S3 bucket with appropriate permissions

### From Source

```bash
# Clone repository
git clone https://github.com/scttfrdmn/cicada.git
cd cicada

# Build and install
make install

# Verify installation
cicada version
```

### AWS Setup

Cicada needs AWS credentials with S3 permissions:

```bash
# Configure AWS CLI
aws configure

# Or set environment variables
export AWS_PROFILE=myprofile
export AWS_REGION=us-west-2
```

**Required S3 Permissions**:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": [
        "arn:aws:s3:::your-bucket",
        "arn:aws:s3:::your-bucket/*"
      ]
    }
  ]
}
```

## Configuration

Initialize configuration:

```bash
cicada config init
```

This creates `~/.cicada/config.yaml`:

```yaml
version: "1"

aws:
  profile: default
  region: us-west-2

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
  check_updates: true
```

### Configuration Commands

```bash
# Set AWS profile
cicada config set aws.profile myprofile

# Set default concurrency
cicada config set sync.concurrency 8

# Get configuration value
cicada config get aws.region

# List all configuration
cicada config list
```

## Usage

### Sync Command

Sync files between local filesystem and S3:

```bash
# Basic sync
cicada sync <source> <destination>

# Options
  --dry-run        Preview changes without syncing
  --delete         Delete files in destination not in source
  --concurrency N  Number of concurrent transfers (default: 4)
  --verbose        Show detailed output
```

**Examples**:

```bash
# Upload to S3
cicada sync /data/experiment s3://my-bucket/experiment

# Download from S3
cicada sync s3://my-bucket/experiment /data/experiment

# Dry run (preview only)
cicada sync --dry-run /data s3://my-bucket/data

# Sync with deletion
cicada sync --delete /data s3://my-bucket/data

# Increase concurrency for large transfers
cicada sync --concurrency 16 /large-dataset s3://my-bucket/data
```

### Watch Command

Monitor directories and automatically sync changes:

```bash
# Add a watch
cicada watch add <source> <destination> [options]

# Options
  --debounce N      Seconds to wait before syncing (default: 5)
  --min-age N       Minimum file age before sync (default: 10)
  --delete-source   Delete source files after successful sync
  --no-sync-on-start  Don't perform initial sync

# List active watches
cicada watch list

# Remove a watch
cicada watch remove <watch-id>
```

**Examples**:

```bash
# Watch microscope data directory
cicada watch add /data/microscope s3://lab-data/microscope

# Custom debounce and min-age
cicada watch add \
  --debounce 10 \
  --min-age 30 \
  /data/sequencer s3://lab-data/sequencing

# Move files to S3 (delete after upload)
cicada watch add \
  --delete-source \
  /data/completed s3://lab-archive/data

# List all watches
cicada watch list

# Remove a watch
cicada watch remove /data/microscope-1234567890
```

**Watch Behavior**:

- Initial sync on start (unless `--no-sync-on-start`)
- Debouncing: Groups rapid file changes to avoid sync storms
- Min-age filter: Only syncs files older than threshold (prevents syncing partial writes)
- Exclude patterns: Respects global exclude patterns from config
- Persistence: Watches are saved to config and restored on startup

## Use Cases

### Instrument Data Upload

Automatically upload data as instruments write files:

```bash
# Zeiss microscope auto-upload
cicada watch add \
  --debounce 30 \
  --min-age 60 \
  /mnt/zeiss/output s3://lab-data/microscopy/zeiss

# Illumina sequencer
cicada watch add \
  --debounce 60 \
  --min-age 300 \
  /data/sequencer/runs s3://lab-data/sequencing
```

### Data Backup

Regular backup of research data:

```bash
# Daily sync (via cron)
0 2 * * * cicada sync /data/research s3://lab-backup/research

# Continuous backup (watch mode)
cicada watch add /data/active-projects s3://lab-backup/projects
```

### Collaborative Data Sharing

Share data with team via S3:

```bash
# Upload shared data
cicada sync /data/shared s3://team-data/shared

# Download team data
cicada sync s3://team-data/shared /data/team-shared
```

## Performance

Typical performance on modern hardware:

- **Small files** (< 1MB): ~100-200 files/sec
- **Large files** (> 100MB): Limited by network bandwidth
- **Concurrency**: 4 transfers by default (configurable)
- **Memory**: ~50-100MB typical usage

**Optimization tips**:

- Increase concurrency for many small files: `--concurrency 16`
- Use exclude patterns to skip unnecessary files
- Run on machine with good network connectivity to AWS

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/scttfrdmn/cicada.git
cd cicada

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run integration tests (requires AWS credentials)
make test-integration

# Run linters
make lint

# Install locally
make install
```

### Project Structure

```
cicada/
├── cmd/cicada/          # CLI entry point
├── internal/
│   ├── cli/             # CLI commands
│   ├── sync/            # Sync engine (backends, engine)
│   ├── watch/           # File watching system
│   ├── config/          # Configuration management
│   ├── metadata/        # Metadata extraction (future)
│   ├── doi/             # DOI management (future)
│   └── integration/     # Integration tests
├── Makefile             # Build targets
└── README.md
```

### Testing

```bash
# Unit tests
go test ./...

# Integration tests (requires AWS)
AWS_PROFILE=aws AWS_REGION=us-west-2 \
  go test -v -tags=integration ./internal/integration/...

# Test coverage
go test -cover ./...
```

## Troubleshooting

### AWS Credentials Issues

```bash
# Verify AWS configuration
aws sts get-caller-identity

# Test S3 access
aws s3 ls s3://your-bucket

# Use specific profile
AWS_PROFILE=myprofile cicada sync /data s3://bucket/data
```

### Permission Denied

Ensure your AWS user/role has required S3 permissions (see [AWS Setup](#aws-setup))

### Slow Syncs

- Check network connectivity to AWS region
- Increase concurrency: `--concurrency 8`
- Ensure exclude patterns are working (check with `--dry-run --verbose`)

### File Not Syncing

- Check exclude patterns in config
- Verify file age meets `--min-age` threshold (watch mode)
- Enable verbose mode: `--verbose`

## Roadmap

**v0.1.0 (Current)**: Core sync and watch functionality

**v0.2.0**:
- Background daemon service
- Improved path handling
- Multipart uploads for large files
- Bandwidth throttling

**v0.3.0+**:
- Metadata extraction and management
- DOI minting integration
- Public data portal
- Workflow execution support

See [VISION.md](VISION.md) for the complete roadmap and [CHANGELOG.md](CHANGELOG.md) for release history.

## Contributing

Contributions welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.

## Citation

If you use Cicada in your research, please cite:

```bibtex
@software{cicada2025,
  title = {Cicada: Dormant Data Commons for Academic Research},
  author = {Scott Friedman},
  year = {2025},
  url = {https://github.com/scttfrdmn/cicada},
  version = {0.1.0}
}
```

## Acknowledgments

- Built with [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper)
- AWS SDK for Go v2
- Inspired by rsync, rclone, and aws-cli
