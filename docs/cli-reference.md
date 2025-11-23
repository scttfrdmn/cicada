# Cicada CLI Reference

Complete command-line interface reference with examples for all major workflows.

## Table of Contents

1. [Installation & Setup](#installation--setup)
2. [Basic Commands](#basic-commands)
3. [Data Sync](#data-sync)
4. [Metadata Management](#metadata-management)
5. [Workflows](#workflows)
6. [Workstations](#workstations)
7. [DOI Management](#doi-management)
8. [User & Project Management](#user--project-management)
9. [Cost Management](#cost-management)
10. [Compliance](#compliance)

---

## Installation & Setup

### Install Cicada

```bash
# macOS / Linux (via Homebrew)
brew install cicada

# Or download binary directly
curl -sSL https://cicada.sh/install | sh

# Or build from source
git clone https://github.com/your-org/cicada.git
cd cicada
make install

# Verify installation
cicada version
```

### Initial Setup

```bash
# Interactive setup wizard
cicada init

# Example output:
# Welcome to Cicada! Let's set up your lab's data management.
# 
# Step 1: AWS Configuration
#   Do you have AWS credentials? [Y/n] y
#   AWS Access Key ID: ****************
#   AWS Secret Access Key: ****************
#   Default region [us-east-1]: us-west-2
#   ✓ Credentials validated
# 
# Step 2: Storage Configuration
#   Lab name (for bucket naming): smith-lab
#   Create S3 bucket? [Y/n] y
#   Bucket name [smith-lab-data-20241122]: 
#   Enable versioning (recommended)? [Y/n] y
#   Enable intelligent tiering? [Y/n] y
#   ✓ Bucket created: s3://smith-lab-data-20241122
# 
# Step 3: Cost Controls
#   Monthly budget alert threshold: $100
#   Email for alerts: pi@university.edu
#   ✓ CloudWatch budget alert configured
# 
# Setup complete!

# Or non-interactive
cicada init \
  --lab-name smith-lab \
  --region us-west-2 \
  --bucket smith-lab-data \
  --budget 100 \
  --email pi@university.edu
```

### Configuration

```bash
# View current configuration
cicada config show

# Edit configuration
cicada config edit

# Set specific values
cicada config set sync.concurrency 20
cicada config set cost.budget_limit 150
cicada config set auth.provider globus
```

---

## Basic Commands

### Status & Information

```bash
# Show overall status
cicada status

# Example output:
# Cicada Status
# Lab: smith-lab
# Bucket: smith-lab-data-20241122
# Region: us-west-2
# 
# Storage:
#   Total: 12.3 TB
#   This month: $78.40
# 
# Active Resources:
#   • Daemon: running (PID 12345)
#   • Watches: 2 locations
#   • Workflows: 1 running
#   • Workstations: 0 active
# 
# Recent Activity:
#   • 2m ago: 45 files synced (8.2 GB)
#   • 1h ago: Workflow completed
#   • 3h ago: User added to project

# Check daemon status
cicada daemon status

# View logs
cicada logs --tail 50 --follow
```

---

## Data Sync

### Manual Sync

```bash
# Basic sync (upload local to S3)
cicada sync /local/data s3://lab-bucket/data/

# Sync with options
cicada sync /local/data s3://lab-bucket/data/ \
  --dry-run \
  --checksum \
  --delete \
  --exclude "*.tmp" \
  --exclude ".DS_Store" \
  --concurrency 10 \
  --bandwidth-limit 50MB

# Bidirectional sync
cicada sync /local/data s3://lab-bucket/data/ --bidirectional

# Download from S3
cicada sync s3://lab-bucket/data/ /local/data/

# Show sync statistics
cicada sync --stats s3://lab-bucket/data/

# Example output:
# Sync Statistics: s3://lab-bucket/data/
# Files: 1,247
# Total size: 15.3 GB
# Last modified: 2 hours ago
# Storage class:
#   - Standard: 500 GB
#   - Intelligent-Tiering: 14.8 GB
```

### Automated Watching

```bash
# Add a watch location
cicada watch add microscope-1 \
  --path /Volumes/ZeissMicroscope/Export \
  --destination s3://lab-bucket/raw/microscopy/ \
  --sync-on-new \
  --min-age 5m \
  --delete-source

# With schedule
cicada watch add sequencer \
  --path /data/sequencer/output \
  --destination s3://lab-bucket/raw/sequencing/ \
  --sync-schedule "0 2 * * *"  # 2 AM daily

# List watches
cicada watch list

# Example output:
# Active Watches (2):
# 
#   microscope-1
#     Path: /Volumes/ZeissMicroscope/Export
#     Destination: s3://lab-bucket/raw/microscopy/
#     Trigger: on new files (min age: 5m)
#     Delete source: yes
#     Status: watching (last sync: 3m ago)
# 
#   sequencer
#     Path: /data/sequencer/output
#     Destination: s3://lab-bucket/raw/sequencing/
#     Schedule: daily at 2:00 AM
#     Status: idle (next sync: today 2:00 AM)

# Show watch details
cicada watch status microscope-1

# Pause/resume watch
cicada watch pause microscope-1
cicada watch resume microscope-1

# Remove watch
cicada watch remove microscope-1
```

### Daemon Management

```bash
# Start daemon (required for watches)
cicada daemon start

# Start with web UI
cicada daemon start --web --port 7878

# Stop daemon
cicada daemon stop

# Restart daemon
cicada daemon restart

# View daemon logs
cicada daemon logs --follow
```

---

## Metadata Management

### Schema Management

```bash
# List available schemas
cicada metadata schema list

# Search for schemas
cicada metadata schema search "microscopy"

# Example output:
# Community Schemas (3 results):
#   1. fluorescence-microscopy (smith-lab)
#      ★★★★☆ 23 labs using
#      "Comprehensive metadata for fluorescence microscopy"
# 
#   2. super-resolution-microscopy (chen-lab)
#      ★★★☆☆ 8 labs using
#      "Extended fields for PALM/STORM imaging"
# 
#   3. live-cell-imaging (williams-lab)
#      ★★★★★ 45 labs using
#      "Time-lapse microscopy metadata"

# View schema details
cicada metadata schema show fluorescence-microscopy

# Install community schema
cicada metadata schema install live-cell-imaging --from williams-lab

# Create custom schema
cicada metadata schema create my-experiment \
  --template microscopy-basic \
  --edit

# Apply schema to folder
cicada metadata schema apply fluorescence-microscopy \
  --path s3://lab-bucket/raw/microscopy/

# Validate schema
cicada metadata schema validate my-experiment.yaml
```

### Working with Metadata

```bash
# Upload with automatic metadata extraction
cicada upload experiment_001.czi s3://lab-bucket/raw/microscopy/

# Example interaction:
# Uploading: experiment_001.czi (2.3 GB)
# 
# Auto-extracted metadata:
#   ✓ Instrument: Zeiss LSM 980
#   ✓ Magnification: 63x
#   ✓ Dimensions: 2048x2048x45x100
#   ✓ Channels: 4 (DAPI, GFP, RFP, Cy5)
# 
# Please provide additional metadata:
# Sample ID: [WT-001]
# Strain: [BY4741]
# Treatment: [1) control  2) heat-shock  3) drug-A]
# Choose: [2]
# Duration (min): [30]
# 
# ✓ Metadata saved

# Upload with metadata file
cicada upload experiment_002.czi \
  --metadata metadata.json \
  s3://lab-bucket/raw/microscopy/

# Upload with template
cicada upload experiment_003.czi \
  --template heat-shock-standard \
  --override sample_id=WT-003 \
  s3://lab-bucket/raw/microscopy/

# Batch upload with CSV
cicada upload --batch samples.csv \
  --files-dir /path/to/experiments/ \
  s3://lab-bucket/raw/microscopy/

# View metadata
cicada metadata show s3://lab-bucket/raw/microscopy/experiment_001.czi

# Edit metadata
cicada metadata edit s3://lab-bucket/raw/microscopy/experiment_001.czi

# Export metadata
cicada metadata export s3://lab-bucket/project/final_data/ \
  --format datacite \
  --output metadata/

# Other formats: isa-tab, dats-json, frictionless, ro-crate

# Search by metadata
cicada search \
  --schema fluorescence-microscopy \
  --where "treatment.condition=heat-shock" \
  --where "magnification=63" \
  --where "date>=2024-11-01"

# Quality check
cicada metadata quality-check s3://lab-bucket/project/

# Example output:
# Metadata Quality Report
# Overall Score: 87/100 (Good)
# 
# Completeness: 92/100
#   ✓ 1,234 files have complete metadata
#   ⚠ 45 files missing "replication" fields
# 
# Recommendations:
#   1. Add protocol IDs → [Fix]
#   2. Standardize strain names → [Fix automatically]
```

### Metadata Templates

```bash
# Save current metadata as template
cicada metadata template save \
  --from experiment_001.czi \
  --name heat-shock-standard \
  --fields strain,treatment,protocol_id

# List templates
cicada metadata template list

# Use template
cicada upload experiment_new.czi \
  --template heat-shock-standard \
  s3://lab-bucket/raw/
```

---

## Workflows

### Running Workflows

```bash
# Run Snakemake workflow
cicada workflow run snakemake \
  --snakefile Snakefile \
  --config input=s3://lab-bucket/raw/experiment_123/ \
  --cores 32 \
  --memory 64GB \
  --spot

# Run Nextflow workflow
cicada workflow run nextflow \
  --workflow pipeline.nf \
  --input "s3://lab-bucket/raw/**.fastq.gz" \
  --outdir s3://lab-bucket/results/run_456/

# Run with config file
cicada workflow run --config cicada-workflow.yaml

# Example cicada-workflow.yaml:
# name: image-processing-pipeline
# engine: snakemake
# workflow: Snakefile
# compute:
#   type: batch
#   instance_types: [c5.4xlarge, c5.9xlarge]
#   spot: true
#   max_vcpus: 256
# storage:
#   input: s3://lab-bucket/raw/
#   output: s3://lab-bucket/processed/
# cost_limit: 50

# Test workflow locally first
cicada workflow run snakemake \
  --snakefile Snakefile \
  --local \
  --cores 4

# Monitor running workflows
cicada workflow list

# Example output:
# Active Workflows (2):
# 
#   cell-segmentation (run_456)
#     Started: 15 minutes ago
#     Progress: 12/45 steps (26%)
#     Cost so far: $2.15
#     ETA: 30 minutes
#     [View logs] [Stop]
# 
#   rnaseq-pipeline (run_457)
#     Started: 2 hours ago
#     Status: Queued (waiting for resources)

# View workflow details
cicada workflow status run_456

# View logs
cicada workflow logs run_456 --follow

# Stop workflow
cicada workflow stop run_456

# Enable auto-workflow on new data
cicada workflow enable auto-pipeline

# Example auto-pipeline.yaml:
# name: auto-fastq-processing
# trigger:
#   watch: /data/sequencer/output/
#   pattern: "*.fastq.gz"
#   min_files: 2  # Wait for R1 and R2
# workflow:
#   engine: nextflow
#   pipeline: nf-core/rnaseq
```

---

## Workstations

### Launching Remote Workstations

```bash
# Launch basic workstation
cicada workstation launch viz-session

# Launch with specific configuration
cicada workstation launch viz-session \
  --instance g4dn.xlarge \
  --image napari-workstation \
  --storage s3://lab-bucket/data/experiment_123/ \
  --spot

# Example output:
# Launching workstation...
# Instance starting: i-0abc123def456
# Waiting for connection... ready!
# 
# Connect via:
#   Web:  https://viz-session-abc123.cicada.cloud
#   VNC:  vnc://54.123.45.67:5901
#   SSH:  ssh cicada@54.123.45.67
# 
# Data mounted at: /data
# Auto-shutdown: 2 hours of inactivity
# Cost: ~$0.50/hour (spot)

# List available images
cicada workstation images

# Example output:
# Available Workstation Images:
#   - basic-linux      Ubuntu 22.04, basic tools
#   - imagej           ImageJ, FIJI, common plugins
#   - matlab           MATLAB R2024a (requires license)
#   - napari           Napari, Python scientific stack
#   - paraview         ParaView, VTK visualization
#   - rstudio          RStudio Server, tidyverse
#   - jupyter          JupyterLab, scipy stack

# Build custom image
cicada workstation build \
  --from napari \
  --add-package cellpose \
  --add-package stardist \
  --name my-workstation

# List active sessions
cicada workstation list

# Example output:
# Active Sessions (1):
# 
#   viz-session
#     Instance: g4dn.xlarge (GPU)
#     Running: 1h 23m
#     Cost: $0.68
#     Auto-shutdown: in 37 minutes
#     URL: https://viz-session-abc123.cicada.cloud

# Reconnect to session
cicada workstation connect viz-session

# Extend auto-shutdown
cicada workstation extend viz-session --hours 4

# Stop workstation
cicada workstation stop viz-session

# Resume stopped workstation
cicada workstation start viz-session

# Snapshot workstation
cicada workstation snapshot viz-session \
  --name "before-analysis"

# Launch from snapshot
cicada workstation launch \
  --from-snapshot "before-analysis"
```

---

## DOI Management

### Minting DOIs

```bash
# Prepare dataset for publication
cicada dataset prepare s3://lab-bucket/project/final_data/ \
  --title "Protein localization under heat shock" \
  --authors "Smith J, Garcia A, Chen L" \
  --description description.md \
  --keywords "protein localization, heat shock, yeast" \
  --license CC-BY-4.0

# Validate dataset readiness
cicada dataset validate s3://lab-bucket/project/final_data/

# Example output:
# Dataset Validation
# ✓ All files have metadata
# ✓ README present
# ✓ Protocols documented
# ✓ Code available (GitHub linked)
# ✓ FAIR principles satisfied
# 
# Ready to publish!

# Mint DOI and publish
cicada dataset publish s3://lab-bucket/project/final_data/ \
  --mint-doi \
  --notify-coauthors

# Example output:
# Minting DOI...
# ✓ DOI registered: 10.12345/cicada.smith-lab.2024.001
# 
# Landing page: https://smithlab.cicada.sh/doi/10.12345/...
# Citation: Smith J, Garcia A, Chen L (2024). Protein
#   localization under heat shock. Smith Lab Data Repository.
#   https://doi.org/10.12345/cicada.smith-lab.2024.001
# 
# Data availability statement (copy for manuscript):
# "Data are available at https://doi.org/10.12345/...
#  under a CC-BY-4.0 license."
# 
# Cost: $1.00 (charged to lab account)

# Publish with embargo
cicada dataset publish s3://lab-bucket/project/final_data/ \
  --mint-doi \
  --embargo-until 2025-06-01 \
  --embargo-type metadata-only

# List lab DOIs
cicada doi list

# Example output:
# Lab DOIs (3):
#   10.12345/cicada.smith-lab.2024.001
#     Title: Protein localization under heat shock
#     Status: Public
#     Downloads: 47
#     Citations: 3
# 
#   10.12345/cicada.smith-lab.2024.002
#     Title: Cell division time-lapse
#     Status: Embargo (until 2025-06-01)

# View DOI details
cicada doi show 10.12345/cicada.smith-lab.2024.001

# Update DOI metadata
cicada doi update 10.12345/cicada.smith-lab.2024.001 \
  --add-author "Williams R" \
  --version 1.1

# Configure DOI provider
cicada lab configure-doi --provider datacite
# or --provider zenodo
# or --provider institution
```

---

## User & Project Management

### User Management

```bash
# Add user to lab
cicada user add jsmith@university.edu \
  --role postdoc \
  --groups protein-structure,methods-dev

# List users
cicada user list

# Example output:
# Lab Members (5):
#   👤 pi@university.edu (PI, admin)
#   👤 jsmith@university.edu (Postdoc)
#      Groups: protein-structure, methods-dev
#   👤 agarcia@university.edu (Grad Student)
#      Groups: metabolism

# Remove user
cicada user remove jsmith@university.edu

# Modify user permissions
cicada user update jsmith@university.edu \
  --add-group new-project \
  --role senior-researcher
```

### Project Management

```bash
# Create project
cicada project create NIH-R01-2024 \
  --description "Protein structure determination" \
  --members pi@university.edu,jsmith@university.edu \
  --budget 200

# List projects
cicada project list

# Show project details
cicada project info NIH-R01-2024

# Example output:
# Project: NIH-R01-2024
# Description: Protein structure determination
# 
# Storage:
#   Total: 2.3 TB
#   Cost: $45/month
#   Growth: +50 GB/week
# 
# Members (3):
#   👤 pi@university.edu (admin)
#   👤 jsmith@university.edu (read-write)
#   👤 external@otheruniv.edu (read-only)
# 
# Recent Activity:
#   • 2h ago: jsmith uploaded 45 files
#   • 1d ago: workflow completed
#   • 3d ago: external accessed dataset_v3

# Add member to project
cicada project add-member NIH-R01-2024 external@otheruniv.edu \
  --permission read-only

# Remove member
cicada project remove-member NIH-R01-2024 external@otheruniv.edu
```

---

## Cost Management

```bash
# View cost report
cicada cost report

# Example output:
# Cost Report: November 2024
# 
# Storage:
#   S3 Standard:              $12.30  (500 GB)
#   S3 Intelligent-Tiering:   $45.80  (6.2 TB)
#   S3 Glacier:               $8.20   (2.1 TB)
# 
# Compute:
#   Batch (spot):             $3.45   (28 vCPU-hours)
#   Workstations:             $12.60  (24 instance-hours)
# 
# Transfer:
#   Data egress:              $2.15   (24 GB)
# 
# Total: $84.50 / $100.00 budget
# 
# Recommendations:
#   ⚡ Archive data >2 years → save $15/month
#   💡 Use smaller instances → save $5/month

# Detailed cost breakdown
cicada cost breakdown --month 2024-11

# Cost prediction
cicada cost predict --action "archive data older than 2 years"

# Example output:
# Estimated savings: $15.20/month
# One-time cost: $2.50
# Payback: 1 month

# Set budget alerts
cicada cost set-budget 100 --alert-at 80

# View by project
cicada cost by-project
```

---

## Compliance

```bash
# Enable compliance mode
cicada lab configure --compliance nist-800-171

# Generate compliance report
cicada compliance report --standard nist-800-171

# View audit log
cicada audit log --user jsmith --last 30d

# Export audit logs
cicada audit export --format csv --year 2024

# Scan for sensitive data
cicada scan s3://lab-bucket/data/ --detect-sensitive

# Example output:
# ⚠️  Found sensitive data:
#   - 15 files contain SSN patterns
#   - 8 files contain credit card numbers
#   - 142 files contain dates of birth
# 
# Recommendations:
#   1. Apply stricter access controls
#   2. Enable field-level encryption

# Verify file integrity
cicada verify s3://lab-bucket/data/experiment.tif
```

---

## Advanced Examples

### Complete Workflow Example

```bash
# 1. Initial setup
cicada init --lab-name smith-lab

# 2. Configure automatic syncing
cicada watch add microscope \
  --path /Volumes/Microscope/Export \
  --destination s3://smith-lab-data/raw/microscopy/ \
  --sync-on-new \
  --delete-source

# 3. Upload existing data with metadata
cicada upload /archive/old_data/ \
  s3://smith-lab-data/archive/ \
  --batch metadata.csv

# 4. Run analysis workflow
cicada workflow run snakemake \
  --snakefile analysis.smk \
  --config experiment=exp_123 \
  --spot

# 5. Visualize results
cicada workstation launch \
  --image napari-workstation \
  --data s3://smith-lab-data/processed/exp_123/

# 6. Prepare for publication
cicada dataset prepare s3://smith-lab-data/final/ \
  --title "My Research Dataset" \
  --authors "Smith J, et al"

# 7. Mint DOI and publish
cicada dataset publish s3://smith-lab-data/final/ \
  --mint-doi \
  --license CC-BY-4.0

# 8. Monitor costs
cicada cost report
```

---

## Getting Help

```bash
# General help
cicada --help

# Command-specific help
cicada sync --help
cicada workflow --help

# Version info
cicada version

# Check for updates
cicada update --check
```
