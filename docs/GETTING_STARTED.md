# Getting Started with Cicada

**Last Updated:** 2025-11-25

A comprehensive tutorial for getting started with Cicada, the small lab data commons platform.

---

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [First-Time Setup](#first-time-setup)
4. [Your First Sync](#your-first-sync)
5. [Adding Metadata](#adding-metadata)
6. [Setting Up Watch Mode](#setting-up-watch-mode)
7. [Working with AWS S3](#working-with-aws-s3)
8. [Next Steps](#next-steps)

---

## Introduction

Cicada is a lightweight data commons platform designed for small research labs. It helps you:

- **Organize** research data with consistent structure
- **Sync** data between local storage and AWS S3
- **Extract** metadata from scientific file formats automatically
- **Monitor** directories for changes and sync automatically
- **Preserve** data provenance and research context

This guide will walk you through installing Cicada, configuring it for your lab, and performing common tasks.

### What You'll Learn

By the end of this tutorial, you'll be able to:

1. Install and configure Cicada
2. Sync data between local storage and S3
3. Extract and view metadata from scientific files
4. Set up automatic file monitoring
5. Organize your research data effectively

### Prerequisites

- macOS, Linux, or Windows machine
- Basic command-line familiarity
- (Optional) AWS account with S3 access for cloud sync

---

## Installation

### Option 1: Download Pre-Built Binary (Recommended)

Download the latest release from GitHub:

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/scttfrdmn/cicada/releases/latest/download/cicada-darwin-arm64
chmod +x cicada-darwin-arm64
sudo mv cicada-darwin-arm64 /usr/local/bin/cicada

# macOS (Intel)
curl -LO https://github.com/scttfrdmn/cicada/releases/latest/download/cicada-darwin-amd64
chmod +x cicada-darwin-amd64
sudo mv cicada-darwin-amd64 /usr/local/bin/cicada

# Linux
curl -LO https://github.com/scttfrdmn/cicada/releases/latest/download/cicada-linux-amd64
chmod +x cicada-linux-amd64
sudo mv cicada-linux-amd64 /usr/local/bin/cicada

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/scttfrdmn/cicada/releases/latest/download/cicada-windows-amd64.exe" -OutFile "cicada.exe"
```

### Option 2: Install from Source

If you have Go 1.21 or later installed:

```bash
# Clone the repository
git clone https://github.com/scttfrdmn/cicada.git
cd cicada

# Build and install
make install
```

### Verify Installation

```bash
cicada version
```

You should see output like:

```
Cicada v0.2.0
Commit: abc1234
Built: 2025-11-25T10:30:00Z
```

---

## First-Time Setup

### 1. Create Your Lab Directory Structure

Cicada works best with an organized directory structure. Here's a recommended layout:

```bash
# Create your lab's data directory
mkdir -p ~/lab-data
cd ~/lab-data

# Create subdirectories for different data types
mkdir -p raw/{microscopy,sequencing,mass-spec}
mkdir -p processed
mkdir -p metadata
mkdir -p analysis
```

This structure separates:
- **raw/**: Original, unmodified data files
- **processed/**: Analysis results and derived data
- **metadata/**: Extracted metadata and documentation
- **analysis/**: Scripts and notebooks

### 2. Initialize Cicada Configuration

Create a basic configuration file:

```bash
cicada config init
```

This creates `~/.config/cicada/config.yaml` with default settings:

```yaml
version: "1.0"

sync:
  local_root: "/Users/username/lab-data"
  remote_root: "s3://my-lab-bucket/data"
  concurrency: 4
  delete_mode: false

metadata:
  auto_extract: true
  formats:
    - microscopy
    - sequencing
    - mass-spec
    - imaging

watch:
  enabled: false
  scan_interval: 60s
  debounce_delay: 5s
```

### 3. Customize Your Configuration

Edit the config file to match your setup:

```bash
# Edit with your preferred editor
nano ~/.config/cicada/config.yaml

# Or use Cicada's config set command
cicada config set sync.local_root ~/lab-data
cicada config set metadata.auto_extract true
```

Key settings to configure:

- **sync.local_root**: Path to your local data directory
- **sync.remote_root**: S3 bucket path (if using cloud sync)
- **sync.concurrency**: Number of concurrent uploads (2-8 recommended)
- **metadata.auto_extract**: Automatically extract metadata on sync

### 4. Verify Configuration

```bash
cicada config show
```

Review the output to ensure your settings are correct.

---

## Your First Sync

Now let's sync some data! We'll start with local-only operations before moving to S3.

### Example 1: Local Directory Sync

Let's say you have microscopy data in one location and want to organize it:

```bash
# Source: Raw data from microscope
# Destination: Organized lab directory

cicada sync /path/to/microscope/export ~/lab-data/raw/microscopy
```

Cicada will:
1. Copy files from the microscope export directory
2. Extract metadata from supported formats
3. Create `.cicada/metadata/` directory with extracted metadata
4. Preserve file timestamps and structure

**Output:**

```
Starting sync: /path/to/microscope/export → ~/lab-data/raw/microscopy
Scanning source directory...
Found 15 files (2.3 GB)

Syncing files:
[████████████████████] 15/15 files (2.3 GB) - 100%

Extracting metadata:
✓ image_001.nd2 → microscopy metadata extracted
✓ image_002.nd2 → microscopy metadata extracted
✓ image_003.tif → imaging metadata extracted
...

Sync complete!
Duration: 45s
Files synced: 15
Metadata extracted: 12
Errors: 0
```

### Example 2: Sync with Dry Run

Before syncing large datasets, test with `--dry-run`:

```bash
cicada sync --dry-run /source/data ~/lab-data/raw
```

This shows what *would* happen without actually copying files:

```
DRY RUN MODE - No files will be modified

Would sync:
  + /source/data/experiment1.fastq.gz → ~/lab-data/raw/sequencing/
  + /source/data/experiment2.fastq.gz → ~/lab-data/raw/sequencing/
  + /source/data/sample.mzML → ~/lab-data/raw/mass-spec/

Summary:
  Files to copy: 3
  Total size: 4.5 GB
  Estimated duration: 2m 30s
```

### Example 3: Incremental Sync

Cicada only syncs changed files by default:

```bash
# First sync
cicada sync /source ~/lab-data/raw
# Copies all files

# Add more files to source
# Second sync
cicada sync /source ~/lab-data/raw
# Only copies new/changed files
```

Cicada uses file size and modification time to detect changes efficiently.

---

## Adding Metadata

Cicada automatically extracts metadata from scientific file formats, but you can also extract metadata manually or add custom metadata.

### View Extracted Metadata

After syncing, metadata is stored in `.cicada/metadata/`:

```bash
# View metadata for a specific file
cicada metadata show ~/lab-data/raw/microscopy/image_001.nd2
```

**Output:**

```json
{
  "file": "image_001.nd2",
  "type": "microscopy",
  "instrument": "nikon_nd2",
  "extracted_at": "2025-11-25T10:45:00Z",
  "metadata": {
    "dimensions": {
      "width": 2048,
      "height": 2048,
      "channels": 3,
      "z_slices": 20,
      "time_points": 1
    },
    "microscope": {
      "manufacturer": "Nikon",
      "model": "Ti2-E",
      "objective": "60x Oil",
      "na": 1.4
    },
    "acquisition": {
      "timestamp": "2025-11-24T15:30:00Z",
      "exposure_ms": 100,
      "binning": "1x1"
    },
    "channels": [
      {"name": "DAPI", "wavelength": 405, "power": 10},
      {"name": "FITC", "wavelength": 488, "power": 15},
      {"name": "TRITC", "wavelength": 561, "power": 20}
    ]
  }
}
```

### Extract Metadata Manually

Extract metadata without syncing:

```bash
# Single file
cicada metadata extract /path/to/file.nd2

# Directory (recursive)
cicada metadata extract /path/to/directory --recursive
```

### Search Metadata

Find files by metadata attributes:

```bash
# Find all Nikon microscopy files
cicada metadata list ~/lab-data --filter "instrument=nikon_nd2"

# Find files from specific date
cicada metadata list ~/lab-data --filter "date=2025-11-24"

# Find sequencing files with specific read length
cicada metadata list ~/lab-data --filter "type=sequencing,read_length=150"
```

### Export Metadata

Export metadata for analysis:

```bash
# Export all metadata as JSON
cicada metadata export ~/lab-data --format json > metadata.json

# Export as CSV for spreadsheet analysis
cicada metadata export ~/lab-data --format csv > metadata.csv

# Export specific fields only
cicada metadata export ~/lab-data --fields "file,type,instrument,date" --format csv
```

### Add Custom Metadata

Add your own metadata fields:

```bash
# Add metadata to a file
cicada metadata add ~/lab-data/raw/microscopy/image_001.nd2 \
  --field experiment_id=EXP001 \
  --field sample_type="HeLa cells" \
  --field treatment="Vehicle control" \
  --field replicate=1

# Add metadata to multiple files
cicada metadata add ~/lab-data/raw/microscopy/*.nd2 \
  --field experiment_id=EXP001 \
  --field project="Cell division study"
```

Custom metadata is merged with extracted metadata:

```json
{
  "file": "image_001.nd2",
  "type": "microscopy",
  "instrument": "nikon_nd2",
  "custom_metadata": {
    "experiment_id": "EXP001",
    "sample_type": "HeLa cells",
    "treatment": "Vehicle control",
    "replicate": 1,
    "project": "Cell division study"
  },
  "metadata": {
    ...
  }
}
```

---

## Setting Up Watch Mode

Watch mode automatically monitors directories for changes and syncs them in real-time.

### Enable Watch Daemon

Start the watch daemon:

```bash
cicada watch start
```

**Output:**

```
Starting Cicada watch daemon...
Daemon started (PID: 12345)
Monitoring: ~/lab-data/raw

Watching for changes in:
  - ~/lab-data/raw/microscopy
  - ~/lab-data/raw/sequencing
  - ~/lab-data/raw/mass-spec

Press Ctrl+C to stop (daemon will continue in background)
```

### Configure Watch Behavior

Edit watch settings in `~/.config/cicada/config.yaml`:

```yaml
watch:
  enabled: true
  scan_interval: 60s        # How often to scan directories
  debounce_delay: 5s        # Wait time after last change
  sync_on_change: true      # Automatically sync changes

  # Which directories to watch
  paths:
    - ~/lab-data/raw/microscopy
    - ~/lab-data/raw/sequencing
    - ~/lab-data/raw/mass-spec

  # What to do on changes
  on_change:
    - action: sync
      source: local
      destination: s3
    - action: extract_metadata
```

### Watch a Specific Directory

Monitor a single directory temporarily:

```bash
# Watch and sync to S3 when files change
cicada watch ~/lab-data/raw/microscopy --sync-to s3://my-bucket/data

# Watch and extract metadata only (no sync)
cicada watch ~/lab-data/raw --metadata-only

# Watch with verbose output
cicada watch ~/lab-data/raw --verbose
```

**Interactive Output:**

```
Watching: ~/lab-data/raw/microscopy
Press Ctrl+C to stop

[15:30:45] File added: image_010.nd2 (1.2 GB)
[15:30:45] Waiting 5s for additional changes...
[15:30:50] Change detected - starting sync
[15:30:51] Extracting metadata from image_010.nd2...
[15:31:05] Uploading to S3: image_010.nd2 (1.2 GB)
[15:31:58] Upload complete (53s, 23 MB/s)
[15:31:58] Ready - watching for changes
```

### Check Watch Status

```bash
# Check if daemon is running
cicada watch status

# View recent activity
cicada watch log --tail 50

# Stop the daemon
cicada watch stop
```

---

## Working with AWS S3

Cicada provides seamless integration with AWS S3 for cloud backup and collaboration.

### Prerequisites

1. **AWS Account**: Sign up at https://aws.amazon.com
2. **S3 Bucket**: Create a bucket in your preferred region
3. **AWS Credentials**: Configure access keys

### Configure AWS Credentials

**Option 1: AWS CLI (Recommended)**

```bash
# Install AWS CLI
brew install awscli  # macOS
# or
pip install awscli   # Python

# Configure credentials
aws configure

# Enter your credentials:
AWS Access Key ID: AKIAIOSFODNN7EXAMPLE
AWS Secret Access Key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
Default region: us-east-1
Default output format: json
```

**Option 2: Environment Variables**

```bash
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=us-east-1
```

**Option 3: Cicada Config File**

```yaml
# ~/.config/cicada/config.yaml
aws:
  region: us-east-1
  access_key_id: AKIAIOSFODNN7EXAMPLE
  secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

⚠️ **Security Note**: Never commit AWS credentials to version control!

### Create an S3 Bucket

```bash
# Create bucket using AWS CLI
aws s3 mb s3://my-lab-data-bucket --region us-east-1

# Or use Cicada
cicada s3 create-bucket my-lab-data-bucket --region us-east-1
```

**Bucket Naming Best Practices:**
- Use lowercase letters, numbers, and hyphens only
- Include your lab name: `smith-lab-microscopy-data`
- Include purpose: `-backup`, `-archive`, `-raw-data`
- Be specific: Easier to manage multiple buckets

### Sync Local Data to S3

Upload your local data to S3:

```bash
# Upload entire directory
cicada sync ~/lab-data s3://my-lab-data-bucket/data

# Upload specific subdirectory
cicada sync ~/lab-data/raw/microscopy s3://my-lab-data-bucket/microscopy

# Upload with progress bar
cicada sync ~/lab-data s3://my-lab-data-bucket/data --progress
```

**Output:**

```
Starting sync: ~/lab-data → s3://my-lab-data-bucket/data
Scanning source directory...
Found 156 files (45.2 GB)

Uploading files:
[████████░░░░░░░░░░░░] 78/156 files (22.1 GB) - 50%
Current: image_078.nd2 (1.2 GB)
Speed: 18.5 MB/s
ETA: 18m 30s

...

Sync complete!
Duration: 32m 15s
Files uploaded: 156
Data transferred: 45.2 GB
Average speed: 23.4 MB/s
Errors: 0
```

### Sync S3 Data to Local

Download data from S3 to your local machine:

```bash
# Download entire bucket
cicada sync s3://my-lab-data-bucket/data ~/lab-data

# Download specific prefix
cicada sync s3://my-lab-data-bucket/microscopy ~/lab-data/raw/microscopy

# Download only new/changed files
cicada sync s3://my-lab-data-bucket/data ~/lab-data --incremental
```

### Bidirectional Sync

Keep local and S3 in sync:

```bash
# Sync in both directions (like rsync)
cicada sync ~/lab-data s3://my-lab-data-bucket/data --bidirectional

# With delete (removes files not in source)
cicada sync ~/lab-data s3://my-lab-data-bucket/data --bidirectional --delete
```

⚠️ **Warning**: `--delete` flag removes files! Use `--dry-run` first.

### Configure Automatic S3 Sync

Set up automatic syncing in your config:

```yaml
# ~/.config/cicada/config.yaml
sync:
  local_root: ~/lab-data
  remote_root: s3://my-lab-data-bucket/data
  mode: bidirectional
  schedule: "0 */6 * * *"  # Every 6 hours
  concurrency: 8

watch:
  enabled: true
  on_change:
    - action: sync
      destination: s3
      delay: 5m  # Wait 5 min after changes stop
```

This configuration:
1. Syncs every 6 hours automatically
2. Watches for local changes
3. Uploads to S3 after 5 minutes of inactivity
4. Uses 8 concurrent uploads for speed

### Monitor S3 Storage Costs

```bash
# Check bucket size
cicada s3 du s3://my-lab-data-bucket

# Check storage class distribution
cicada s3 storage-class s3://my-lab-data-bucket

# Estimate monthly costs
cicada s3 cost-estimate s3://my-lab-data-bucket
```

**Output:**

```
Bucket: s3://my-lab-data-bucket
Total size: 45.2 GB
Files: 156

Storage class distribution:
  STANDARD: 45.2 GB (100%)
  STANDARD_IA: 0 B (0%)
  GLACIER: 0 B (0%)

Estimated monthly cost (us-east-1):
  Storage: $1.04 (45.2 GB × $0.023/GB)
  Requests: $0.02 (estimated)
  Data transfer: $0.00 (downloads not included)
  Total: ~$1.06/month
```

### Use Lifecycle Policies

Reduce costs by moving old data to cheaper storage:

```bash
# Create lifecycle policy
cicada s3 lifecycle create s3://my-lab-data-bucket \
  --transition-ia 30 \
  --transition-glacier 90 \
  --expiration 365

# This policy:
# - After 30 days → STANDARD_IA (cheaper storage)
# - After 90 days → GLACIER (archive storage)
# - After 365 days → Delete (optional)
```

**Cost Comparison:**

| Storage Class | Cost per GB/month | Retrieval Time |
|--------------|-------------------|----------------|
| STANDARD | $0.023 | Instant |
| STANDARD_IA | $0.0125 | Instant |
| GLACIER | $0.004 | 3-5 hours |
| DEEP_ARCHIVE | $0.00099 | 12 hours |

---

## Next Steps

Congratulations! You've completed the Cicada getting started tutorial. Here's what you've learned:

✅ Install and configure Cicada
✅ Sync data locally and to S3
✅ Extract and view metadata
✅ Set up automatic file watching
✅ Work with AWS S3 for cloud backup

### Continue Learning

1. **Common Workflows**: Learn specific workflows for your research
   - See: [Common Workflows Guide](WORKFLOWS.md)

2. **Advanced Topics**: Dive deeper into Cicada features
   - See: [Advanced Topics](ADVANCED.md)

3. **Integration**: Integrate Cicada with your existing tools
   - See: [Integration Guide](INTEGRATIONS.md)

4. **Troubleshooting**: Solve common problems
   - See: [Troubleshooting Guide](TROUBLESHOOTING.md)

5. **API Reference**: Use Cicada programmatically
   - See: [API Reference](API.md)

### Get Help

- **Documentation**: https://scttfrdmn.github.io/cicada
- **GitHub Issues**: https://github.com/scttfrdmn/cicada/issues
- **Discussions**: https://github.com/scttfrdmn/cicada/discussions

### Share Feedback

We'd love to hear from you! Share your experience:

```bash
# Submit feedback
cicada feedback send

# Report a bug
cicada bug-report
```

---

## Quick Reference

### Essential Commands

```bash
# Configuration
cicada config init                    # Create config file
cicada config show                    # View current config
cicada config set <key> <value>       # Update setting

# Syncing
cicada sync <source> <destination>    # Sync files
cicada sync --dry-run <src> <dst>     # Preview sync
cicada sync --delete <src> <dst>      # Delete extra files

# Metadata
cicada metadata extract <path>        # Extract metadata
cicada metadata show <file>           # View metadata
cicada metadata list <path>           # List all metadata
cicada metadata export <path>         # Export metadata

# Watch mode
cicada watch start                    # Start daemon
cicada watch status                   # Check status
cicada watch stop                     # Stop daemon

# S3 operations
cicada s3 ls s3://bucket              # List S3 contents
cicada s3 du s3://bucket              # Check bucket size
cicada s3 create-bucket <name>        # Create bucket
```

### Configuration File Location

```
~/.config/cicada/config.yaml          # macOS/Linux
%APPDATA%\cicada\config.yaml          # Windows
```

### Getting Help

```bash
cicada help                           # General help
cicada <command> --help               # Command-specific help
cicada version                        # Version info
```

---

Happy data organizing! 🎉
