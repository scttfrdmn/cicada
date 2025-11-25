# Integration Guide

**Last Updated:** 2025-11-25

Complete guide for integrating Cicada with workflow managers, scripts, and lab systems.

## Table of Contents

1. [Integration Overview](#integration-overview)
2. [Nextflow Pipelines](#nextflow-pipelines)
3. [Snakemake Workflows](#snakemake-workflows)
4. [Python Scripts](#python-scripts)
5. [Bash Scripts](#bash-scripts)
6. [CI/CD Pipelines](#cicd-pipelines)
7. [Lab Information Management Systems](#lab-information-management-systems)
8. [Docker Integration](#docker-integration)
9. [Best Practices](#best-practices)
10. [Example Workflows](#example-workflows)

---

## Integration Overview

Cicada integrates seamlessly with existing lab workflows and systems.

### Integration Points

| System | Integration Method | Use Case |
|--------|-------------------|----------|
| Nextflow | Process executor | Pipeline data management |
| Snakemake | Rule commands | Workflow automation |
| Python | Subprocess calls | Custom scripts |
| Bash | Direct CLI | Shell automation |
| CI/CD | GitHub Actions, GitLab CI | Automated testing/deployment |
| LIMS | API integration | Lab workflow automation |
| Docker | Container execution | Reproducible environments |

### Common Patterns

**1. Data Acquisition → Cicada → Analysis**
```
Instrument → cicada sync → S3 → Analysis Pipeline
```

**2. Analysis → Cicada → Archive**
```
Pipeline Results → cicada metadata extract → cicada sync → S3
```

**3. Watch Mode → Automated Processing**
```
cicada watch → Trigger Pipeline → Process Data
```

---

## Nextflow Pipelines

### Basic Integration

**Process with Cicada:**
```groovy
process downloadFromS3 {
    output:
    path "data/*"

    script:
    """
    cicada sync s3://lab-data/input/ data/
    """
}

process extractMetadata {
    input:
    path dataFile

    output:
    path "${dataFile}.metadata.json"

    script:
    """
    cicada metadata extract ${dataFile} --output ${dataFile}.metadata.json
    """
}

process uploadResults {
    input:
    path results

    script:
    """
    cicada sync ${results} s3://lab-data/results/\$(date +%Y-%m-%d)/
    """
}
```

### Complete Pipeline Example

**Microscopy Image Analysis:**
```groovy
#!/usr/bin/env nextflow

params.input_bucket = "s3://lab-data/microscopy/raw"
params.output_bucket = "s3://lab-data/microscopy/processed"
params.metadata_preset = "microscopy-confocal"

workflow {
    // Download data from S3
    downloadData(params.input_bucket)

    // Extract and validate metadata
    extractMetadata(downloadData.out)
    validateMetadata(extractMetadata.out, params.metadata_preset)

    // Process images
    processImages(downloadData.out)

    // Upload results
    uploadResults(processImages.out, params.output_bucket)
}

process downloadData {
    input:
    val s3_path

    output:
    path "raw_data/*"

    script:
    """
    echo "Downloading from ${s3_path}"
    cicada sync ${s3_path} raw_data/
    """
}

process extractMetadata {
    input:
    path dataFile

    output:
    tuple path(dataFile), path("${dataFile}.metadata.json")

    script:
    """
    cicada metadata extract ${dataFile} --output ${dataFile}.metadata.json --format json
    """
}

process validateMetadata {
    input:
    tuple path(dataFile), path(metadataFile)
    val preset

    output:
    path "${dataFile}.validated"

    script:
    """
    cicada metadata validate ${dataFile} --preset ${preset}
    if [ \$? -eq 0 ]; then
        touch ${dataFile}.validated
    else
        echo "Validation failed for ${dataFile}"
        exit 1
    fi
    """
}

process processImages {
    input:
    path dataFile

    output:
    path "processed/${dataFile.baseName}_processed.tif"

    script:
    """
    mkdir -p processed
    # Your image processing here
    python process_image.py ${dataFile} processed/${dataFile.baseName}_processed.tif
    """
}

process uploadResults {
    publishDir params.output_bucket, mode: 'copy'

    input:
    path processedFile
    val bucket

    script:
    """
    cicada metadata extract ${processedFile} --output ${processedFile}.metadata.json
    cicada sync ${processedFile} ${bucket}
    cicada sync ${processedFile}.metadata.json ${bucket}
    """
}
```

### Running the Pipeline

```bash
# Local execution
nextflow run microscopy_pipeline.nf

# With parameters
nextflow run microscopy_pipeline.nf \
  --input_bucket s3://lab-data/exp001/raw \
  --output_bucket s3://lab-data/exp001/processed \
  --metadata_preset microscopy-confocal

# With configuration file
nextflow run microscopy_pipeline.nf -c nextflow.config
```

**Configuration File (nextflow.config):**
```groovy
params {
    input_bucket = "s3://lab-data/microscopy/raw"
    output_bucket = "s3://lab-data/microscopy/processed"
    metadata_preset = "microscopy-confocal"
}

process {
    executor = 'local'
    cpus = 4
    memory = '16 GB'
}

env {
    AWS_PROFILE = 'research-account'
}
```

---

## Snakemake Workflows

### Basic Integration

**Snakefile:**
```python
rule download_data:
    output:
        "data/raw/{sample}.czi"
    params:
        s3_path = "s3://lab-data/raw/{sample}.czi"
    shell:
        "cicada sync {params.s3_path} {output}"

rule extract_metadata:
    input:
        "data/raw/{sample}.czi"
    output:
        "data/metadata/{sample}.metadata.json"
    shell:
        "cicada metadata extract {input} --output {output}"

rule validate_metadata:
    input:
        data = "data/raw/{sample}.czi",
        metadata = "data/metadata/{sample}.metadata.json"
    output:
        "data/validated/{sample}.validated"
    params:
        preset = "microscopy-confocal"
    shell:
        """
        cicada metadata validate {input.data} --preset {params.preset}
        touch {output}
        """

rule upload_results:
    input:
        "data/processed/{sample}.tif"
    params:
        s3_bucket = "s3://lab-data/processed"
    shell:
        "cicada sync {input} {params.s3_bucket}/{wildcards.sample}.tif"
```

### Complete Workflow Example

**Sequencing Data Pipeline:**
```python
# Snakefile
configfile: "config.yaml"

SAMPLES = config["samples"]
BUCKET = config["s3_bucket"]

rule all:
    input:
        expand("results/{sample}.stats.txt", sample=SAMPLES),
        "results/multiqc_report.html"

rule download_fastq:
    output:
        r1 = "data/raw/{sample}_R1.fastq.gz",
        r2 = "data/raw/{sample}_R2.fastq.gz"
    params:
        s3_prefix = f"{BUCKET}/raw/{{sample}}"
    shell:
        """
        cicada sync {params.s3_prefix}_R1.fastq.gz {output.r1}
        cicada sync {params.s3_prefix}_R2.fastq.gz {output.r2}
        """

rule extract_metadata:
    input:
        "data/raw/{sample}_R1.fastq.gz"
    output:
        "data/metadata/{sample}.metadata.json"
    shell:
        "cicada metadata extract {input} --output {output} --format json"

rule validate_metadata:
    input:
        fastq = "data/raw/{sample}_R1.fastq.gz",
        metadata = "data/metadata/{sample}.metadata.json"
    output:
        "data/validated/{sample}.validated"
    params:
        preset = config["metadata_preset"]
    shell:
        """
        cicada metadata validate {input.fastq} --preset {params.preset} --strict
        touch {output}
        """

rule quality_control:
    input:
        r1 = "data/raw/{sample}_R1.fastq.gz",
        r2 = "data/raw/{sample}_R2.fastq.gz",
        validated = "data/validated/{sample}.validated"
    output:
        "results/{sample}.stats.txt"
    shell:
        "fastqc {input.r1} {input.r2} -o results/"

rule multiqc:
    input:
        expand("results/{sample}.stats.txt", sample=SAMPLES)
    output:
        "results/multiqc_report.html"
    shell:
        "multiqc results/ -o results/"

rule upload_results:
    input:
        stats = "results/{sample}.stats.txt",
        report = "results/multiqc_report.html"
    params:
        s3_results = f"{BUCKET}/results"
    shell:
        """
        cicada sync results/ {params.s3_results}/$(date +%Y-%m-%d)/
        """
```

**Configuration (config.yaml):**
```yaml
samples:
  - sample001
  - sample002
  - sample003

s3_bucket: "s3://lab-data/sequencing"
metadata_preset: "sequencing-illumina"
```

**Running the Workflow:**
```bash
# Dry run
snakemake -n

# Execute
snakemake --cores 8

# With specific config
snakemake --cores 8 --configfile config.yaml

# Generate DAG
snakemake --dag | dot -Tpdf > workflow.pdf
```

---

## Python Scripts

### Using Subprocess

**Basic Example:**
```python
import subprocess
import json
from pathlib import Path

def sync_to_s3(local_path, s3_path):
    """Sync files to S3 using Cicada."""
    result = subprocess.run(
        ["cicada", "sync", local_path, s3_path],
        capture_output=True,
        text=True
    )

    if result.returncode != 0:
        raise RuntimeError(f"Sync failed: {result.stderr}")

    print(f"Synced {local_path} to {s3_path}")
    return result.returncode

def extract_metadata(filepath, output_path=None):
    """Extract metadata from file."""
    cmd = ["cicada", "metadata", "extract", filepath, "--format", "json"]

    if output_path:
        cmd.extend(["--output", output_path])

    result = subprocess.run(cmd, capture_output=True, text=True)

    if result.returncode != 0:
        raise RuntimeError(f"Extraction failed: {result.stderr}")

    # Parse JSON output
    if not output_path:
        return json.loads(result.stdout)
    else:
        with open(output_path) as f:
            return json.load(f)

def validate_metadata(filepath, preset):
    """Validate file metadata against preset."""
    result = subprocess.run(
        ["cicada", "metadata", "validate", filepath, "--preset", preset],
        capture_output=True,
        text=True
    )

    return result.returncode == 0

# Usage
if __name__ == "__main__":
    # Extract metadata
    metadata = extract_metadata("experiment001.czi")
    print(f"Format: {metadata['format']}")
    print(f"Dimensions: {metadata['width']}x{metadata['height']}")

    # Validate
    if validate_metadata("experiment001.czi", "microscopy-confocal"):
        print("Validation passed")

        # Sync to S3
        sync_to_s3("experiment001.czi", "s3://lab-data/validated/")
    else:
        print("Validation failed")
```

### Complete Script Example

**Automated Data Processing:**
```python
#!/usr/bin/env python3
"""
Automated microscopy data processing and upload script.
"""

import subprocess
import json
import logging
from pathlib import Path
from datetime import datetime
import sys

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class CicadaIntegration:
    """Wrapper for Cicada CLI operations."""

    def __init__(self, verbose=False):
        self.verbose = verbose

    def run_command(self, cmd):
        """Execute Cicada command."""
        if self.verbose:
            logger.info(f"Running: {' '.join(cmd)}")

        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True
        )

        if result.returncode != 0:
            logger.error(f"Command failed: {result.stderr}")
            raise RuntimeError(f"Command failed: {' '.join(cmd)}")

        return result

    def extract_metadata(self, filepath):
        """Extract metadata from file."""
        cmd = [
            "cicada", "metadata", "extract",
            str(filepath),
            "--format", "json"
        ]

        result = self.run_command(cmd)
        return json.loads(result.stdout)

    def validate_metadata(self, filepath, preset):
        """Validate metadata against preset."""
        cmd = [
            "cicada", "metadata", "validate",
            str(filepath),
            "--preset", preset,
            "--strict"
        ]

        try:
            self.run_command(cmd)
            return True
        except RuntimeError:
            return False

    def sync_to_s3(self, local_path, s3_path, dry_run=False):
        """Sync files to S3."""
        cmd = ["cicada", "sync", str(local_path), s3_path]

        if dry_run:
            cmd.append("--dry-run")

        result = self.run_command(cmd)
        return result.returncode == 0

def process_microscopy_data(data_dir, s3_bucket, preset="microscopy-confocal"):
    """
    Process microscopy data: extract metadata, validate, and upload.

    Args:
        data_dir: Directory containing microscopy files
        s3_bucket: S3 bucket path for upload
        preset: Metadata validation preset
    """
    cicada = CicadaIntegration(verbose=True)
    data_path = Path(data_dir)

    # Find all CZI files
    czi_files = list(data_path.glob("*.czi"))
    logger.info(f"Found {len(czi_files)} CZI files")

    results = {
        "processed": [],
        "failed": [],
        "skipped": []
    }

    for czi_file in czi_files:
        logger.info(f"Processing: {czi_file.name}")

        try:
            # Extract metadata
            logger.info("Extracting metadata...")
            metadata = cicada.extract_metadata(czi_file)

            # Save metadata to sidecar file
            metadata_file = czi_file.with_suffix(czi_file.suffix + ".metadata.json")
            with open(metadata_file, 'w') as f:
                json.dump(metadata, f, indent=2)

            logger.info(f"Format: {metadata.get('format')}")
            logger.info(f"Dimensions: {metadata.get('width')}x{metadata.get('height')}")

            # Validate metadata
            logger.info("Validating metadata...")
            if not cicada.validate_metadata(czi_file, preset):
                logger.warning(f"Validation failed for {czi_file.name}")
                results["failed"].append(str(czi_file))
                continue

            logger.info("Validation passed")

            # Upload to S3
            logger.info(f"Uploading to {s3_bucket}...")

            # Upload data file
            s3_data_path = f"{s3_bucket}/{datetime.now().strftime('%Y-%m-%d')}/{czi_file.name}"
            cicada.sync_to_s3(czi_file, s3_data_path)

            # Upload metadata file
            s3_metadata_path = f"{s3_bucket}/{datetime.now().strftime('%Y-%m-%d')}/{metadata_file.name}"
            cicada.sync_to_s3(metadata_file, s3_metadata_path)

            logger.info(f"Successfully processed {czi_file.name}")
            results["processed"].append(str(czi_file))

        except Exception as e:
            logger.error(f"Error processing {czi_file.name}: {e}")
            results["failed"].append(str(czi_file))

    # Summary
    logger.info("=" * 60)
    logger.info("Processing Summary:")
    logger.info(f"  Processed: {len(results['processed'])}")
    logger.info(f"  Failed: {len(results['failed'])}")
    logger.info(f"  Skipped: {len(results['skipped'])}")

    return results

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python process_data.py <data_dir> <s3_bucket> [preset]")
        sys.exit(1)

    data_dir = sys.argv[1]
    s3_bucket = sys.argv[2]
    preset = sys.argv[3] if len(sys.argv) > 3 else "microscopy-confocal"

    results = process_microscopy_data(data_dir, s3_bucket, preset)

    # Exit with error if any failed
    sys.exit(1 if results["failed"] else 0)
```

**Usage:**
```bash
python process_data.py /data/microscopy s3://lab-data/processed microscopy-confocal
```

---

## Bash Scripts

### Basic Integration

**Simple Sync Script:**
```bash
#!/bin/bash
# sync_daily_data.sh

set -e  # Exit on error

DATA_DIR="/data/microscope"
S3_BUCKET="s3://lab-data/microscopy"
DATE=$(date +%Y-%m-%d)

echo "Syncing data from $DATE"

# Sync data to S3
cicada sync "$DATA_DIR" "$S3_BUCKET/$DATE/"

echo "Sync complete"
```

### Complete Workflow Script

**Automated Processing Script:**
```bash
#!/bin/bash
# process_and_upload.sh
#
# Process microscopy data, extract metadata, validate, and upload to S3

set -euo pipefail

# Configuration
DATA_DIR="${1:-/data/microscopy}"
S3_BUCKET="${2:-s3://lab-data/processed}"
PRESET="${3:-microscopy-confocal}"
TEMP_DIR="/tmp/cicada-processing-$$"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Cleanup on exit
cleanup() {
    log_info "Cleaning up temporary files..."
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Create temp directory
mkdir -p "$TEMP_DIR"

# Counter variables
PROCESSED=0
FAILED=0
SKIPPED=0

# Process each CZI file
log_info "Scanning directory: $DATA_DIR"

while IFS= read -r -d '' file; do
    filename=$(basename "$file")
    log_info "Processing: $filename"

    # Extract metadata
    log_info "  Extracting metadata..."
    if ! cicada metadata extract "$file" --output "$TEMP_DIR/$filename.metadata.json" --format json; then
        log_error "  Metadata extraction failed"
        ((FAILED++))
        continue
    fi

    # Validate metadata
    log_info "  Validating metadata..."
    if ! cicada metadata validate "$file" --preset "$PRESET" --strict; then
        log_warn "  Validation failed"
        ((FAILED++))
        continue
    fi

    log_info "  Validation passed"

    # Upload to S3
    DATE_PATH=$(date +%Y/%m/%d)
    S3_PATH="$S3_BUCKET/$DATE_PATH"

    log_info "  Uploading to $S3_PATH..."

    # Upload data file
    if ! cicada sync "$file" "$S3_PATH/$filename"; then
        log_error "  Upload failed"
        ((FAILED++))
        continue
    fi

    # Upload metadata
    if ! cicada sync "$TEMP_DIR/$filename.metadata.json" "$S3_PATH/$filename.metadata.json"; then
        log_warn "  Metadata upload failed (data uploaded successfully)"
    fi

    log_info "  ✓ Successfully processed $filename"
    ((PROCESSED++))

done < <(find "$DATA_DIR" -name "*.czi" -print0)

# Summary
echo ""
log_info "========================================="
log_info "Processing Summary:"
log_info "  Processed: $PROCESSED"
log_info "  Failed: $FAILED"
log_info "  Skipped: $SKIPPED"
log_info "========================================="

# Exit with error if any failed
if [ $FAILED -gt 0 ]; then
    exit 1
fi
```

**Usage:**
```bash
# Basic usage
./process_and_upload.sh

# With arguments
./process_and_upload.sh /data/experiments/exp001 s3://lab-data/exp001 microscopy-confocal

# Run as cron job (daily at 2 AM)
0 2 * * * /path/to/process_and_upload.sh /data/microscopy s3://lab-data/daily
```

---

## CI/CD Pipelines

### GitHub Actions

**Workflow File (.github/workflows/data-validation.yml):**
```yaml
name: Data Validation

on:
  push:
    paths:
      - 'data/**'
  pull_request:
    paths:
      - 'data/**'
  workflow_dispatch:

jobs:
  validate-data:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Download Cicada
        run: |
          curl -L https://github.com/scttfrdmn/cicada/releases/latest/download/cicada-linux-amd64 -o cicada
          chmod +x cicada
          sudo mv cicada /usr/local/bin/

      - name: Verify installation
        run: cicada version

      - name: Extract metadata
        run: |
          mkdir -p metadata
          for file in data/*.czi; do
            if [ -f "$file" ]; then
              cicada metadata extract "$file" --output "metadata/$(basename $file).json"
            fi
          done

      - name: Validate metadata
        run: |
          EXIT_CODE=0
          for file in data/*.czi; do
            if [ -f "$file" ]; then
              echo "Validating $file..."
              if ! cicada metadata validate "$file" --preset microscopy-confocal; then
                echo "✗ Validation failed: $file"
                EXIT_CODE=1
              else
                echo "✓ Validation passed: $file"
              fi
            fi
          done
          exit $EXIT_CODE

      - name: Upload metadata artifacts
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: metadata
          path: metadata/

      - name: Configure AWS credentials
        if: github.ref == 'refs/heads/main'
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Upload to S3
        if: github.ref == 'refs/heads/main'
        run: |
          cicada sync data/ s3://lab-data/validated/$(date +%Y-%m-%d)/
          cicada sync metadata/ s3://lab-data/metadata/$(date +%Y-%m-%d)/
```

### GitLab CI

**.gitlab-ci.yml:**
```yaml
stages:
  - validate
  - upload

variables:
  CICADA_VERSION: "0.2.0"
  S3_BUCKET: "s3://lab-data"

before_script:
  - apt-get update && apt-get install -y curl
  - curl -L https://github.com/scttfrdmn/cicada/releases/download/v${CICADA_VERSION}/cicada-linux-amd64 -o cicada
  - chmod +x cicada
  - mv cicada /usr/local/bin/
  - cicada version

validate_metadata:
  stage: validate
  script:
    - |
      for file in data/*.czi; do
        echo "Processing $file"
        cicada metadata extract "$file" --output "$file.metadata.json"
        cicada metadata validate "$file" --preset microscopy-confocal --strict
      done
  artifacts:
    paths:
      - data/*.metadata.json
    expire_in: 1 week

upload_to_s3:
  stage: upload
  only:
    - main
  script:
    - cicada sync data/ ${S3_BUCKET}/$(date +%Y-%m-%d)/
  dependencies:
    - validate_metadata
```

---

## Lab Information Management Systems

### REST API Integration

**Example LIMS Integration:**
```python
import requests
import subprocess
import json

class LIMSCicadaIntegration:
    """Integrate Cicada with LIMS system."""

    def __init__(self, lims_url, api_key):
        self.lims_url = lims_url
        self.api_key = api_key
        self.headers = {
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json"
        }

    def get_samples_for_processing(self):
        """Fetch samples ready for processing from LIMS."""
        response = requests.get(
            f"{self.lims_url}/api/samples",
            headers=self.headers,
            params={"status": "ready_for_processing"}
        )
        return response.json()["samples"]

    def download_sample_data(self, sample_id, local_path):
        """Download sample data using Cicada."""
        sample = self.get_sample_info(sample_id)
        s3_path = sample["s3_location"]

        subprocess.run(
            ["cicada", "sync", s3_path, local_path],
            check=True
        )

        return local_path

    def extract_and_upload_metadata(self, sample_id, file_path):
        """Extract metadata and upload to LIMS."""
        # Extract metadata using Cicada
        result = subprocess.run(
            ["cicada", "metadata", "extract", file_path, "--format", "json"],
            capture_output=True,
            text=True,
            check=True
        )

        metadata = json.loads(result.stdout)

        # Upload metadata to LIMS
        response = requests.post(
            f"{self.lims_url}/api/samples/{sample_id}/metadata",
            headers=self.headers,
            json=metadata
        )

        return response.json()

    def update_sample_status(self, sample_id, status, metadata=None):
        """Update sample status in LIMS."""
        data = {"status": status}
        if metadata:
            data["metadata"] = metadata

        response = requests.patch(
            f"{self.lims_url}/api/samples/{sample_id}",
            headers=self.headers,
            json=data
        )

        return response.json()

# Usage
lims = LIMSCicadaIntegration("https://lims.example.com", "api_key_here")

# Get samples
samples = lims.get_samples_for_processing()

for sample in samples:
    # Download data
    local_path = f"/tmp/{sample['id']}"
    lims.download_sample_data(sample['id'], local_path)

    # Extract metadata
    metadata = lims.extract_and_upload_metadata(sample['id'], local_path)

    # Update status
    lims.update_sample_status(sample['id'], "processed", metadata)
```

---

## Docker Integration

### Dockerfile

```dockerfile
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Cicada
RUN curl -L https://github.com/scttfrdmn/cicada/releases/latest/download/cicada-linux-amd64 -o /usr/local/bin/cicada \
    && chmod +x /usr/local/bin/cicada

# Verify installation
RUN cicada version

# Set working directory
WORKDIR /data

# Default command
CMD ["cicada", "--help"]
```

### Docker Compose

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  cicada-processor:
    build: .
    volumes:
      - ./data:/data
      - ~/.aws:/root/.aws:ro
    environment:
      - AWS_PROFILE=default
    command: >
      bash -c "
        cicada metadata extract /data/*.czi &&
        cicada sync /data/ s3://lab-data/processed/
      "
```

### Usage

```bash
# Build image
docker build -t cicada-processor .

# Run container
docker run -v $(pwd)/data:/data -v ~/.aws:/root/.aws:ro cicada-processor \
  cicada sync /data s3://lab-data/backup/

# With Docker Compose
docker-compose up
```

---

## Best Practices

### 1. Error Handling

```bash
# Always check exit codes
if ! cicada sync /data s3://bucket/data; then
    echo "Sync failed"
    exit 1
fi

# Use set -e in bash scripts
set -e  # Exit on any error
```

### 2. Logging

```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('/var/log/cicada-integration.log'),
        logging.StreamHandler()
    ]
)

logger = logging.getLogger(__name__)
logger.info("Starting Cicada integration")
```

### 3. Configuration Management

```python
# Use configuration files
import yaml

with open('config.yaml') as f:
    config = yaml.safe_load(f)

s3_bucket = config['s3_bucket']
preset = config['metadata_preset']
```

### 4. Testing

```python
# Unit test example
import unittest
from unittest.mock import patch, MagicMock

class TestCicadaIntegration(unittest.TestCase):

    @patch('subprocess.run')
    def test_metadata_extraction(self, mock_run):
        mock_run.return_value = MagicMock(
            returncode=0,
            stdout='{"format": "CZI"}'
        )

        result = extract_metadata("test.czi")
        self.assertEqual(result['format'], 'CZI')
```

### 5. Idempotency

```bash
# Make scripts idempotent
# Check if already processed
if [ -f "data/processed/.${filename}.done" ]; then
    echo "Already processed: $filename"
    exit 0
fi

# Process
cicada sync "$filename" s3://bucket/

# Mark as done
touch "data/processed/.${filename}.done"
```

---

## Example Workflows

### Workflow 1: Automated Nightly Backup

```bash
#!/bin/bash
# nightly-backup.sh

# Configuration
DATA_DIRS=("/data/microscopy" "/data/sequencing" "/data/mass-spec")
S3_BUCKET="s3://lab-backups"
DATE=$(date +%Y-%m-%d)

# Backup each directory
for dir in "${DATA_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo "Backing up $dir..."
        cicada sync "$dir" "$S3_BUCKET/$DATE/$(basename $dir)/"
    fi
done

# Cron: 0 2 * * * /path/to/nightly-backup.sh
```

### Workflow 2: Data Quality Pipeline

```python
# quality_pipeline.py

def quality_control_pipeline(data_dir, s3_bucket):
    """Complete quality control pipeline."""

    for file in Path(data_dir).glob("*.czi"):
        # 1. Extract metadata
        metadata = extract_metadata(file)

        # 2. Validate metadata
        if not validate_metadata(file, "microscopy-confocal"):
            logger.error(f"Failed validation: {file}")
            continue

        # 3. Check image quality
        if not check_image_quality(file, metadata):
            logger.error(f"Poor image quality: {file}")
            continue

        # 4. Upload to S3
        sync_to_s3(file, f"{s3_bucket}/validated/")
        sync_to_s3(f"{file}.metadata.json", f"{s3_bucket}/metadata/")

        # 5. Trigger downstream processing
        trigger_analysis(file, metadata)
```

### Workflow 3: Watch + Process

```bash
# Setup watch mode to trigger processing
cicada watch add /data/microscope s3://lab-data/raw/

# Separate script monitors S3 and processes
# monitor_and_process.sh

aws s3 ls s3://lab-data/raw/ --recursive | while read -r line; do
    file=$(echo $line | awk '{print $4}')

    if [[ $file == *.czi ]]; then
        # Download for processing
        aws s3 cp s3://lab-data/raw/$file /tmp/

        # Process
        python process_image.py /tmp/$file

        # Upload results
        cicada sync /tmp/processed/ s3://lab-data/processed/
    fi
done
```

---

## Related Documentation

- [CLI Reference](CLI_REFERENCE.md) - All CLI commands
- [Architecture](ARCHITECTURE.md) - System architecture
- [Python Examples](../examples/python/) - More Python examples
- [Bash Examples](../examples/bash/) - More Bash examples

---

**Contributing:** Have an integration example to share? See [CONTRIBUTING.md](CONTRIBUTING.md)

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
