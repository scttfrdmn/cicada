# Advanced Topics

**Last Updated:** 2025-11-25

Advanced features and techniques for power users of Cicada.

---

## Table of Contents

1. [Advanced Sync Strategies](#advanced-sync-strategies)
2. [Custom Metadata Extractors](#custom-metadata-extractors)
3. [Hooks and Automation](#hooks-and-automation)
4. [Performance Optimization](#performance-optimization)
5. [Security Best Practices](#security-best-practices)
6. [Advanced AWS S3 Features](#advanced-aws-s3-features)
7. [Programmatic Usage](#programmatic-usage)
8. [Multi-Site Deployments](#multi-site-deployments)
9. [Data Governance](#data-governance)
10. [Advanced Monitoring](#advanced-monitoring)

---

## Advanced Sync Strategies

### Bidirectional Sync with Conflict Resolution

Handle conflicts when syncing between multiple locations.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
sync:
  mode: bidirectional
  conflict_resolution: newest  # Options: newest, local, remote, ask, merge

  # Conflict handling strategies
  conflicts:
    strategy: newest           # Use newest file by modification time
    backup_conflicts: true     # Keep copies of conflicted files
    backup_dir: .cicada/conflicts/
```

**Manual Conflict Resolution:**

```bash
# Sync and detect conflicts
cicada sync ~/data s3://bucket/data --detect-conflicts

# List conflicts
cicada conflicts list

# Resolve conflicts interactively
cicada conflicts resolve --interactive

# Example conflict:
# local:  file.txt (2025-11-25 10:00, size: 1234)
# remote: file.txt (2025-11-25 11:00, size: 1456)
# Choose: [l]ocal, [r]emote, [b]oth, [m]erge, [s]kip?
```

**Automated Conflict Resolution Script:**

```bash
#!/bin/bash
# Advanced conflict resolution

cicada sync ~/data s3://bucket/data --detect-conflicts --output json > conflicts.json

# Parse conflicts
python3 << 'PYTHON'
import json

with open('conflicts.json') as f:
    conflicts = json.load(f)

for conflict in conflicts:
    local_time = conflict['local']['modified']
    remote_time = conflict['remote']['modified']
    path = conflict['path']

    # Custom logic: prefer remote for data files, local for analysis
    if path.startswith('raw/'):
        print(f"Using remote: {path}")
        action = 'remote'
    elif path.startswith('analysis/'):
        print(f"Using local: {path}")
        action = 'local'
    else:
        # Use newest for everything else
        action = 'remote' if remote_time > local_time else 'local'

    # Apply resolution
    subprocess.run(['cicada', 'conflicts', 'resolve', path, '--use', action])
PYTHON
```

### Selective Sync with Advanced Filters

Sync only specific files based on complex criteria.

**Filter Configuration:**

```yaml
# ~/.config/cicada/config.yaml
sync:
  filters:
    # Include patterns
    include:
      - "*.fastq.gz"
      - "*.bam"
      - "**/*.nd2"

    # Exclude patterns
    exclude:
      - "*.tmp"
      - "**/.DS_Store"
      - "**/cache/**"

    # Size filters
    min_size: 1KB
    max_size: 10GB

    # Time filters
    modified_after: "2025-01-01"
    modified_before: "2025-12-31"

    # Metadata filters
    metadata_filters:
      - field: "qc_status"
        value: "passed"
      - field: "instrument"
        value: "nikon_nd2"
```

**Dynamic Filtering:**

```bash
# Sync only files from last 7 days
cicada sync ~/data s3://bucket/data \
  --filter "mtime:-7d"

# Sync only large files
cicada sync ~/data s3://bucket/data \
  --filter "size:>1GB"

# Sync only files with specific metadata
cicada sync ~/data s3://bucket/data \
  --metadata-filter "experiment_id=EXP001"

# Combine multiple filters
cicada sync ~/data s3://bucket/data \
  --filter "mtime:-30d" \
  --filter "size:<100MB" \
  --metadata-filter "qc_status=passed"
```

**Complex Filter Expression:**

```bash
# Boolean expressions
cicada sync ~/data s3://bucket/data \
  --filter-expr "(size > 1GB AND mtime < 30d) OR (type = microscopy AND qc_status = passed)"
```

### Sync Chains and Dependencies

Define sync workflows with dependencies.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
sync:
  workflows:
    daily_backup:
      steps:
        - name: sync_raw_data
          source: /data/microscope/export
          destination: ~/lab-data/raw/microscopy

        - name: extract_metadata
          depends_on: sync_raw_data
          command: cicada metadata extract ~/lab-data/raw/microscopy

        - name: backup_to_s3
          depends_on: extract_metadata
          source: ~/lab-data/raw/microscopy
          destination: s3://lab-bucket/microscopy/$(date +%Y-%m-%d)
          storage_class: STANDARD

        - name: archive_old
          depends_on: backup_to_s3
          source: ~/lab-data/raw/microscopy
          destination: s3://lab-bucket/archive/microscopy
          filter: "mtime:+90d"
          storage_class: GLACIER
```

**Execute Workflow:**

```bash
# Run workflow
cicada workflow run daily_backup

# Run specific step
cicada workflow run daily_backup --step backup_to_s3

# Run in parallel (where possible)
cicada workflow run daily_backup --parallel

# Dry run
cicada workflow run daily_backup --dry-run
```

---

## Custom Metadata Extractors

### Creating Custom Extractors

Build custom extractors for proprietary file formats.

**Extractor Interface:**

```go
// internal/metadata/extractor.go
package metadata

type Extractor interface {
    // Check if this extractor can handle the file
    CanHandle(filename string) bool

    // Extract metadata from file
    Extract(filepath string) (map[string]interface{}, error)

    // Extract from io.Reader (for streaming)
    ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error)

    // Extractor name
    Name() string

    // Supported file extensions
    SupportedFormats() []string
}
```

**Example: Custom Microscope Extractor:**

```go
// custom_extractors/custom_microscope.go
package custom

import (
    "encoding/binary"
    "io"
    "os"
    "path/filepath"
    "strings"
)

type CustomMicroscopeExtractor struct{}

func (e *CustomMicroscopeExtractor) Name() string {
    return "custom_microscope"
}

func (e *CustomMicroscopeExtractor) SupportedFormats() []string {
    return []string{".czi", ".lif"}
}

func (e *CustomMicroscopeExtractor) CanHandle(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    return ext == ".czi" || ext == ".lif"
}

func (e *CustomMicroscopeExtractor) Extract(filepath string) (map[string]interface{}, error) {
    f, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    return e.ExtractFromReader(f, filepath)
}

func (e *CustomMicroscopeExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
    // Read file header
    header := make([]byte, 256)
    if _, err := io.ReadFull(r, header); err != nil {
        return nil, err
    }

    metadata := make(map[string]interface{})

    // Parse custom format
    // This is format-specific - example for demonstration
    metadata["format"] = "custom_microscope"
    metadata["version"] = binary.LittleEndian.Uint32(header[0:4])

    // Extract dimensions
    metadata["dimensions"] = map[string]interface{}{
        "width":      binary.LittleEndian.Uint32(header[8:12]),
        "height":     binary.LittleEndian.Uint32(header[12:16]),
        "channels":   binary.LittleEndian.Uint16(header[16:18]),
        "z_slices":   binary.LittleEndian.Uint16(header[18:20]),
        "time_points": binary.LittleEndian.Uint16(header[20:22]),
    }

    // Extract acquisition metadata
    metadata["acquisition"] = map[string]interface{}{
        "timestamp":   parseTimestamp(header[24:32]),
        "exposure_ms": float64(binary.LittleEndian.Uint32(header[32:36])) / 1000.0,
        "gain":        binary.LittleEndian.Uint16(header[36:38]),
    }

    // Extract instrument info
    metadata["instrument"] = map[string]interface{}{
        "manufacturer": readNullTermString(header[64:128]),
        "model":        readNullTermString(header[128:192]),
        "serial":       readNullTermString(header[192:256]),
    }

    return metadata, nil
}

// Helper functions
func parseTimestamp(data []byte) string {
    // Parse timestamp from bytes
    // Format-specific implementation
    return "2025-11-25T10:00:00Z"  // Example
}

func readNullTermString(data []byte) string {
    for i, b := range data {
        if b == 0 {
            return string(data[:i])
        }
    }
    return string(data)
}
```

**Register Custom Extractor:**

```go
// cmd/cicada/main.go
package main

import (
    "github.com/scttfrdmn/cicada/internal/metadata"
    "your-package/custom_extractors"
)

func init() {
    // Register custom extractor
    registry := metadata.GetRegistry()
    registry.Register(&custom.CustomMicroscopeExtractor{})
}
```

**Python Plugin System:**

For easier development, use Python plugins:

```python
# ~/.config/cicada/extractors/custom_microscope.py

import struct
from datetime import datetime

class CustomMicroscopeExtractor:
    """Custom extractor for proprietary microscope format"""

    def name(self):
        return "custom_microscope"

    def supported_formats(self):
        return [".czi", ".lif"]

    def can_handle(self, filename):
        return filename.lower().endswith(('.czi', '.lif'))

    def extract(self, filepath):
        """Extract metadata from file"""
        with open(filepath, 'rb') as f:
            # Read header
            header = f.read(256)

            # Parse format-specific fields
            metadata = {
                'format': 'custom_microscope',
                'version': struct.unpack('<I', header[0:4])[0],
                'dimensions': {
                    'width': struct.unpack('<I', header[8:12])[0],
                    'height': struct.unpack('<I', header[12:16])[0],
                    'channels': struct.unpack('<H', header[16:18])[0],
                    'z_slices': struct.unpack('<H', header[18:20])[0],
                    'time_points': struct.unpack('<H', header[20:22])[0],
                },
                'acquisition': {
                    'timestamp': self.parse_timestamp(header[24:32]),
                    'exposure_ms': struct.unpack('<I', header[32:36])[0] / 1000.0,
                    'gain': struct.unpack('<H', header[36:38])[0],
                },
                'instrument': {
                    'manufacturer': header[64:128].decode('utf-8').split('\0')[0],
                    'model': header[128:192].decode('utf-8').split('\0')[0],
                    'serial': header[192:256].decode('utf-8').split('\0')[0],
                }
            }

            return metadata

    def parse_timestamp(self, data):
        """Parse timestamp from bytes"""
        # Format-specific implementation
        return datetime.now().isoformat()

# Register extractor
def register():
    return CustomMicroscopeExtractor()
```

**Enable Python Plugins:**

```yaml
# ~/.config/cicada/config.yaml
metadata:
  plugins:
    enabled: true
    paths:
      - ~/.config/cicada/extractors
    python: true
```

### Metadata Transformation Pipelines

Transform and enrich metadata during extraction.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
metadata:
  pipelines:
    microscopy_enrichment:
      input_types:
        - microscopy
      steps:
        - name: convert_units
          command: |
            # Convert exposure from ms to s
            .metadata.acquisition.exposure_s = .metadata.acquisition.exposure_ms / 1000

        - name: calculate_resolution
          command: |
            # Calculate pixel size in µm
            .metadata.resolution_um = .metadata.pixel_size_nm / 1000

        - name: add_quality_score
          command: |
            # Calculate quality score based on multiple factors
            .metadata.quality_score = (
              .metadata.acquisition.exposure_ms * 0.3 +
              .metadata.dimensions.width * 0.002 +
              .metadata.dimensions.height * 0.002
            )

        - name: classify_experiment_type
          script: ~/.config/cicada/scripts/classify_experiment.py
```

**Transformation Script Example:**

```python
# ~/.config/cicada/scripts/classify_experiment.py
import sys
import json

def classify_experiment(metadata):
    """Classify experiment type based on metadata"""

    # Extract relevant fields
    channels = metadata.get('metadata', {}).get('channels', [])
    z_slices = metadata.get('metadata', {}).get('dimensions', {}).get('z_slices', 0)
    time_points = metadata.get('metadata', {}).get('dimensions', {}).get('time_points', 0)

    # Classify
    if time_points > 1:
        exp_type = "time-lapse"
    elif z_slices > 1:
        exp_type = "z-stack"
    elif len(channels) > 3:
        exp_type = "multi-channel"
    else:
        exp_type = "standard"

    # Add classification
    metadata['experiment_type'] = exp_type

    # Add complexity score
    complexity = len(channels) * z_slices * time_points
    metadata['complexity_score'] = complexity

    return metadata

if __name__ == '__main__':
    # Read metadata from stdin
    metadata = json.load(sys.stdin)

    # Transform
    metadata = classify_experiment(metadata)

    # Output to stdout
    json.dump(metadata, sys.stdout, indent=2)
```

---

## Hooks and Automation

### Pre/Post Sync Hooks

Execute custom commands before and after sync operations.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
hooks:
  pre_sync:
    enabled: true
    script: ~/.config/cicada/hooks/pre_sync.sh
    timeout: 60

  post_sync:
    enabled: true
    script: ~/.config/cicada/hooks/post_sync.sh
    timeout: 300

  on_error:
    enabled: true
    script: ~/.config/cicada/hooks/on_error.sh
```

**Pre-Sync Hook Example:**

```bash
#!/bin/bash
# ~/.config/cicada/hooks/pre_sync.sh

# Arguments: source destination operation
SOURCE=$1
DEST=$2
OPERATION=$3

echo "Pre-sync hook: $OPERATION from $SOURCE to $DEST"

# Check disk space
AVAILABLE=$(df -k "$SOURCE" | tail -1 | awk '{print $4}')
REQUIRED=$(du -sk "$SOURCE" | awk '{print $1}')

if [ "$AVAILABLE" -lt "$REQUIRED" ]; then
    echo "ERROR: Insufficient disk space"
    exit 1
fi

# Verify AWS credentials if syncing to S3
if [[ "$DEST" == s3://* ]]; then
    if ! aws sts get-caller-identity > /dev/null 2>&1; then
        echo "ERROR: AWS credentials not configured"
        exit 1
    fi
fi

# Create snapshot of current state
cicada snapshot create "$SOURCE" --name "pre-sync-$(date +%Y%m%d-%H%M%S)"

echo "Pre-sync checks passed"
exit 0
```

**Post-Sync Hook Example:**

```bash
#!/bin/bash
# ~/.config/cicada/hooks/post_sync.sh

SOURCE=$1
DEST=$2
OPERATION=$3
STATUS=$4  # success or error

echo "Post-sync hook: $OPERATION completed with status: $STATUS"

if [ "$STATUS" = "success" ]; then
    # Generate sync report
    cicada sync-report --format html > ~/reports/sync_$(date +%Y%m%d-%H%M%S).html

    # Send notification
    curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
      -d "{\"text\": \"Sync completed: $SOURCE → $DEST\"}"

    # Clean up old snapshots
    cicada snapshot clean --keep 5

    # Update log
    echo "$(date): $OPERATION $SOURCE → $DEST [SUCCESS]" >> ~/logs/sync.log
else
    # Send error notification
    curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
      -d "{\"text\": \"❌ Sync failed: $SOURCE → $DEST\"}"

    # Update log
    echo "$(date): $OPERATION $SOURCE → $DEST [FAILED]" >> ~/logs/sync.log
fi
```

### Event-Driven Automation

Respond to file system events with custom actions.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
automation:
  events:
    file_created:
      - trigger:
          path: ~/lab-data/raw/microscopy/*.nd2
        actions:
          - extract_metadata: true
          - run_script: ~/.config/cicada/scripts/process_microscopy.sh
          - sync_to: s3://lab-bucket/microscopy/
          - notify: slack

    file_modified:
      - trigger:
          path: ~/lab-data/analysis/**/*.ipynb
        actions:
          - backup_to: s3://lab-bucket/backups/analysis/
          - run_notebook: true  # Re-execute notebook
          - commit_to_git: true

    metadata_extracted:
      - trigger:
          type: microscopy
          filter: "qc_status=passed"
        actions:
          - run_pipeline: microscopy_analysis
          - add_to_queue: high_priority
```

**Processing Script Example:**

```bash
#!/bin/bash
# ~/.config/cicada/scripts/process_microscopy.sh

FILE=$1
FILETYPE=$2

echo "Processing microscopy file: $FILE"

# Extract metadata
metadata=$(cicada metadata show "$FILE" --json)

# Get dimensions
width=$(echo "$metadata" | jq -r '.metadata.dimensions.width')
height=$(echo "$metadata" | jq -r '.metadata.dimensions.height')

# Generate thumbnail
convert "$FILE[0]" -resize 512x512 "${FILE%.nd2}_thumb.jpg"

# Run quality control
python3 ~/scripts/qc_microscopy.py "$FILE"

# Add QC results to metadata
qc_result=$?
if [ $qc_result -eq 0 ]; then
    cicada metadata add "$FILE" --field "qc_status=passed"
else
    cicada metadata add "$FILE" --field "qc_status=failed"
fi

# Sync processed files
cicada sync "$(dirname "$FILE")" s3://lab-bucket/processed/
```

---

## Performance Optimization

### Parallel Processing

Maximize throughput with parallel operations.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
performance:
  # Sync parallelism
  sync:
    concurrency: 8              # Concurrent file transfers
    part_size: 10MB             # Multipart upload size
    max_retries: 3              # Retry failed transfers

  # Metadata extraction parallelism
  metadata:
    concurrency: 4              # Concurrent extractions
    batch_size: 100             # Files per batch
    worker_pool: 8              # Worker threads

  # IO optimization
  io:
    buffer_size: 4MB            # Read/write buffer
    use_direct_io: false        # Bypass OS cache
    prefetch: true              # Prefetch next files
```

**Benchmark and Tune:**

```bash
# Benchmark different concurrency settings
for concurrency in 2 4 8 16; do
    echo "Testing concurrency: $concurrency"
    time cicada sync ~/test-data s3://test-bucket/data \
      --concurrency $concurrency \
      --progress 2>&1 | tee benchmark_$concurrency.log
done

# Analyze results
python3 << 'PYTHON'
import re
import glob

for log in glob.glob('benchmark_*.log'):
    concurrency = re.search(r'benchmark_(\d+)', log).group(1)
    with open(log) as f:
        content = f.read()
        duration = re.search(r'Duration: (\d+)s', content)
        throughput = re.search(r'(\d+\.?\d*) MB/s', content)

        print(f"Concurrency {concurrency}: {duration.group(1)}s, {throughput.group(1)} MB/s")
PYTHON
```

### Caching Strategies

Reduce redundant operations with intelligent caching.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
cache:
  enabled: true

  # Metadata cache
  metadata:
    enabled: true
    ttl: 3600                   # Cache for 1 hour
    max_size: 1GB               # Maximum cache size
    strategy: lru               # LRU, LFU, or FIFO

  # File listing cache
  listings:
    enabled: true
    ttl: 300                    # Cache for 5 minutes
    s3_ttl: 600                 # S3 listings cache longer

  # Checksum cache
  checksums:
    enabled: true
    ttl: 86400                  # Cache for 24 hours
    persistent: true            # Persist across restarts
```

**Cache Management:**

```bash
# View cache statistics
cicada cache stats

# Clear specific cache
cicada cache clear metadata
cicada cache clear listings
cicada cache clear checksums

# Clear all caches
cicada cache clear --all

# Preload cache
cicada cache preload ~/lab-data --recursive
```

### Compression and Deduplication

Reduce storage and transfer costs.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
optimization:
  # Compression
  compression:
    enabled: true
    algorithm: zstd             # zstd, gzip, lz4
    level: 3                    # 1-22 for zstd
    min_size: 1MB               # Don't compress small files

  # Deduplication
  deduplication:
    enabled: true
    method: content_hash        # content_hash or rolling_hash
    chunk_size: 4MB             # Chunk size for dedup
    cache: true                 # Cache chunk hashes
```

**Usage:**

```bash
# Sync with compression
cicada sync ~/data s3://bucket/data --compress

# Sync with deduplication
cicada sync ~/data s3://bucket/data --deduplicate

# Analyze savings
cicada analyze --compression ~/data
cicada analyze --deduplication ~/data
```

---

## Security Best Practices

### Encryption

Encrypt data at rest and in transit.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
security:
  # Encryption at rest
  encryption:
    enabled: true
    algorithm: AES-256-GCM
    key_source: kms            # kms, keyring, or file

    # AWS KMS
    kms:
      key_id: arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012
      region: us-east-1

  # Encryption in transit
  tls:
    enabled: true
    min_version: "1.2"
    cipher_suites:
      - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

**Client-Side Encryption:**

```bash
# Encrypt before upload
cicada sync ~/data s3://bucket/data \
  --encrypt \
  --kms-key arn:aws:kms:us-east-1:123456789012:key/...

# Decrypt on download
cicada sync s3://bucket/data ~/data \
  --decrypt \
  --kms-key arn:aws:kms:us-east-1:123456789012:key/...
```

### Access Control

Fine-grained access control for multi-user environments.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
access_control:
  enabled: true

  # Define roles
  roles:
    admin:
      permissions:
        - read
        - write
        - delete
        - configure
      paths: ["**"]

    scientist:
      permissions:
        - read
        - write
      paths:
        - "data/**"
        - "analysis/**"

    student:
      permissions:
        - read
      paths:
        - "data/**"
        - "published/**"

  # Assign users to roles
  users:
    alice@example.com: admin
    bob@example.com: scientist
    carol@example.com: student
```

**Audit Logging:**

```yaml
# ~/.config/cicada/config.yaml
audit:
  enabled: true
  log_file: ~/.config/cicada/audit.log
  log_level: info

  # What to log
  events:
    - sync
    - metadata_access
    - configuration_changes
    - access_denials

  # Log format
  format: json
  include_metadata: true
```

**View Audit Logs:**

```bash
# View recent audit events
cicada audit log --tail 100

# Filter by event type
cicada audit log --filter "event=sync"

# Filter by user
cicada audit log --filter "user=alice@example.com"

# Export audit log
cicada audit export --since "2025-11-01" --format csv > audit_november.csv
```

---

## Advanced AWS S3 Features

### S3 Batch Operations

Perform operations on millions of objects efficiently.

**Configuration:**

```bash
# Create S3 Batch job manifest
cicada s3 batch create-manifest s3://bucket/data > manifest.csv

# Create batch job: change storage class
cicada s3 batch create-job \
  --manifest manifest.csv \
  --operation set-storage-class \
  --storage-class GLACIER \
  --role arn:aws:iam::123456789012:role/BatchOperationRole

# Monitor batch job
cicada s3 batch status job-id-12345

# Get results
cicada s3 batch results job-id-12345
```

### S3 Select

Query data directly in S3 without downloading.

**Example:**

```bash
# Query CSV file in S3
cicada s3 select s3://bucket/data/results.csv \
  --query "SELECT * FROM s3object WHERE quality_score > 0.9" \
  --input-format csv \
  --output-format json

# Query JSON file
cicada s3 select s3://bucket/metadata/samples.json \
  --query "SELECT * FROM s3object[*].metadata WHERE instrument = 'nikon_nd2'"

# Query Parquet file
cicada s3 select s3://bucket/analysis/results.parquet \
  --query "SELECT sample_id, expression_level FROM s3object WHERE p_value < 0.05"
```

### S3 Intelligent-Tiering

Automatically move data between access tiers.

**Configuration:**

```bash
# Enable Intelligent-Tiering
cicada s3 tiering enable s3://bucket/data

# Configure archive access tiers
cicada s3 tiering configure s3://bucket/data \
  --archive-access-days 90 \
  --deep-archive-access-days 180

# Monitor tiering
cicada s3 tiering status s3://bucket/data

# View cost savings
cicada s3 tiering savings s3://bucket/data --since 30d
```

### Cross-Region Replication

Replicate data across AWS regions for disaster recovery.

**Configuration:**

```bash
# Enable versioning (required for replication)
aws s3api put-bucket-versioning \
  --bucket source-bucket \
  --versioning-configuration Status=Enabled

aws s3api put-bucket-versioning \
  --bucket destination-bucket \
  --versioning-configuration Status=Enabled

# Create replication configuration
cat > replication.json << 'EOF'
{
  "Role": "arn:aws:iam::123456789012:role/ReplicationRole",
  "Rules": [
    {
      "Status": "Enabled",
      "Priority": 1,
      "Filter": {
        "Prefix": "data/"
      },
      "Destination": {
        "Bucket": "arn:aws:s3:::destination-bucket",
        "ReplicationTime": {
          "Status": "Enabled",
          "Time": {
            "Minutes": 15
          }
        },
        "Metrics": {
          "Status": "Enabled"
        }
      }
    }
  ]
}
EOF

# Apply replication configuration
aws s3api put-bucket-replication \
  --bucket source-bucket \
  --replication-configuration file://replication.json

# Monitor replication
cicada s3 replication status s3://source-bucket
```

---

## Programmatic Usage

### Go API

Use Cicada programmatically in Go applications.

**Example Application:**

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scttfrdmn/cicada/pkg/sync"
    "github.com/scttfrdmn/cicada/pkg/metadata"
    "github.com/scttfrdmn/cicada/pkg/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load("~/.config/cicada/config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // Create sync engine
    engine := sync.NewEngine(cfg)

    // Set up progress callback
    engine.SetProgressFunc(func(update sync.ProgressUpdate) {
        fmt.Printf("[%s] %s: %d/%d bytes\\n",
            update.Operation, update.Path,
            update.BytesDone, update.BytesTotal)
    })

    // Perform sync
    ctx := context.Background()
    result, err := engine.Sync(ctx, "/local/data", "s3://bucket/data")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Synced %d files (%d bytes) in %s\\n",
        result.FileCount, result.BytesTransferred, result.Duration)

    // Extract metadata
    registry := metadata.GetRegistry()
    for _, file := range result.Files {
        meta, err := registry.Extract(file.Path)
        if err != nil {
            log.Printf("Failed to extract metadata from %s: %v", file.Path, err)
            continue
        }

        fmt.Printf("Metadata for %s: %+v\\n", file.Path, meta)
    }
}
```

### Python API

Use Cicada from Python scripts.

**Example Script:**

```python
#!/usr/bin/env python3
# Cicada Python API example

import cicada
from datetime import datetime, timedelta

# Initialize client
client = cicada.Client(config_path='~/.config/cicada/config.yaml')

# Sync data
result = client.sync(
    source='~/lab-data',
    destination='s3://lab-bucket/data',
    progress=True,
    filters={'mtime': '-7d'}  # Last 7 days only
)

print(f"Synced {result.file_count} files ({result.bytes_transferred} bytes)")

# Extract and query metadata
metadata_list = client.metadata.extract('~/lab-data', recursive=True)

# Filter by criteria
microscopy_files = [
    m for m in metadata_list
    if m.type == 'microscopy' and m.metadata.get('qc_status') == 'passed'
]

print(f"Found {len(microscopy_files)} microscopy files that passed QC")

# Export metadata
client.metadata.export(
    path='~/lab-data',
    format='csv',
    output='metadata_export.csv',
    fields=['file', 'type', 'instrument', 'date', 'qc_status']
)

# Watch for changes
def on_change(event):
    print(f"File changed: {event.path}")
    # Extract metadata and sync
    metadata = client.metadata.extract(event.path)
    client.sync(event.path, 's3://lab-bucket/data')

watcher = client.watch('~/lab-data', on_change=on_change)
watcher.start()

# Keep running
try:
    watcher.join()
except KeyboardInterrupt:
    watcher.stop()
```

### REST API

Use Cicada via HTTP REST API.

**Start API Server:**

```bash
# Start Cicada API server
cicada api start --port 8080 --auth-token YOUR_TOKEN

# Or in config
cicada config set api.enabled true
cicada config set api.port 8080
cicada config set api.auth_token YOUR_TOKEN
```

**API Usage:**

```bash
# Sync via API
curl -X POST http://localhost:8080/api/v1/sync \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "/local/data",
    "destination": "s3://bucket/data",
    "options": {
      "concurrency": 4,
      "delete": false
    }
  }'

# Get metadata
curl http://localhost:8080/api/v1/metadata?path=/local/data/file.nd2 \
  -H "Authorization: Bearer YOUR_TOKEN"

# List files
curl http://localhost:8080/api/v1/list?path=s3://bucket/data \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get sync status
curl http://localhost:8080/api/v1/sync/status/job-id \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Multi-Site Deployments

### Hub-and-Spoke Architecture

Central data repository with multiple satellite labs.

**Hub Configuration (Central Site):**

```yaml
# Hub: Central data repository
sync:
  role: hub
  remote_root: s3://central-hub/data

  # Accept connections from spoke sites
  spoke_sites:
    - name: lab-west
      url: https://lab-west.example.com
      auth_token: TOKEN_WEST

    - name: lab-east
      url: https://lab-east.example.com
      auth_token: TOKEN_EAST

  # Sync schedule from spokes
  schedule:
    - site: lab-west
      interval: 1h
      path: microscopy/**

    - site: lab-east
      interval: 6h
      path: sequencing/**
```

**Spoke Configuration (Satellite Lab):**

```yaml
# Spoke: Satellite laboratory
sync:
  role: spoke
  local_root: ~/lab-data
  hub_url: https://central-hub.example.com
  hub_auth_token: HUB_TOKEN

  # What to sync to hub
  sync_to_hub:
    paths:
      - microscopy/**
      - sequencing/**
    schedule: "0 */1 * * *"  # Every hour
    compression: true
```

**Deploy Hub:**

```bash
# On central server
cicada hub start --config hub-config.yaml

# Monitor hub
cicada hub status
cicada hub sites
```

**Deploy Spoke:**

```bash
# On satellite lab server
cicada spoke start --config spoke-config.yaml

# Test connectivity to hub
cicada spoke test-connection

# Manual sync to hub
cicada spoke sync-now
```

### Mesh Architecture

Peer-to-peer data sharing between labs.

**Configuration:**

```yaml
# Mesh participant configuration
sync:
  role: mesh
  node_id: lab-alpha

  # Other mesh participants
  peers:
    - node_id: lab-beta
      url: https://lab-beta.example.com
      shared_paths:
        - shared/projects/**
        - shared/protocols/**

    - node_id: lab-gamma
      url: https://lab-gamma.example.com
      shared_paths:
        - shared/reference-data/**

  # Sync strategy
  mesh:
    strategy: eventual_consistency
    conflict_resolution: newest
    discovery: multicast         # or dns, manual
```

---

## Data Governance

### Retention Policies

Implement data retention and deletion policies.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
governance:
  retention:
    enabled: true

    policies:
      # Raw data: keep 2 years
      - name: raw_data_retention
        paths:
          - "raw/**"
        retain_days: 730
        action: move_to_glacier
        notify: true

      # Analysis results: keep 1 year
      - name: analysis_retention
        paths:
          - "analysis/**"
        retain_days: 365
        action: delete
        require_approval: true

      # Published data: keep forever
      - name: published_data
        paths:
          - "published/**"
        retain_days: null  # forever
        action: none
```

**Compliance Tracking:**

```bash
# Check compliance status
cicada governance status

# List files subject to deletion
cicada governance review --action delete

# Apply retention policies
cicada governance apply --dry-run
cicada governance apply --confirm

# Generate compliance report
cicada governance report --since 2025-01-01 --format pdf > compliance_report.pdf
```

### Data Provenance

Track data lineage and transformations.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
provenance:
  enabled: true

  # Track these events
  track:
    - file_creation
    - file_modification
    - file_deletion
    - metadata_changes
    - transformations
    - access

  # Store provenance data
  storage:
    backend: postgresql
    connection: "postgresql://user:pass@localhost/provenance"
```

**Query Provenance:**

```bash
# Get provenance for file
cicada provenance show ~/data/processed/results.csv

# Output:
# File: results.csv
# Created: 2025-11-25 10:00:00
# Derived from:
#   - raw/experiment_001.fastq.gz (accessed 2025-11-25 09:30:00)
#   - raw/experiment_002.fastq.gz (accessed 2025-11-25 09:30:00)
# Transformations:
#   1. alignment (STAR v2.7.10, 2025-11-25 09:35:00)
#   2. quantification (featureCounts v2.0.1, 2025-11-25 09:45:00)
#   3. normalization (DESeq2 v1.34.0, 2025-11-25 09:55:00)
# Accessed by:
#   - alice@example.com (2025-11-25 10:15:00)
#   - bob@example.com (2025-11-25 14:30:00)

# Generate lineage graph
cicada provenance graph ~/data/processed/results.csv --output graph.svg
```

---

## Advanced Monitoring

### Metrics Collection

Collect and analyze Cicada metrics.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
monitoring:
  metrics:
    enabled: true
    backend: prometheus           # prometheus, influxdb, cloudwatch

    # Prometheus configuration
    prometheus:
      port: 9090
      path: /metrics

    # What to collect
    collect:
      - sync_operations
      - transfer_speeds
      - metadata_extractions
      - errors
      - cache_hit_rate
      - storage_usage
```

**Prometheus Metrics:**

```bash
# Start metrics server
cicada metrics start

# View metrics
curl http://localhost:9090/metrics

# Example metrics:
# cicada_sync_operations_total{status="success"} 1234
# cicada_sync_operations_total{status="error"} 5
# cicada_transfer_bytes_total 123456789012
# cicada_transfer_speed_bytes_per_second 12345678
# cicada_metadata_extractions_total{type="microscopy"} 456
# cicada_cache_hit_rate 0.85
```

**Grafana Dashboard:**

```bash
# Import Cicada Grafana dashboard
grafana-cli plugins install cicada-datasource

# Dashboard JSON available at:
# https://grafana.com/grafana/dashboards/cicada-monitoring
```

### Alerting

Set up alerts for important events.

**Configuration:**

```yaml
# ~/.config/cicada/config.yaml
alerting:
  enabled: true

  # Alert destinations
  destinations:
    email:
      enabled: true
      smtp_server: smtp.gmail.com:587
      from: cicada@example.com
      to:
        - admin@example.com

    slack:
      enabled: true
      webhook_url: https://hooks.slack.com/services/YOUR/WEBHOOK/URL

    pagerduty:
      enabled: true
      integration_key: YOUR_INTEGRATION_KEY

  # Alert rules
  rules:
    - name: sync_failures
      condition: sync_error_rate > 0.1
      severity: critical
      message: "High sync failure rate detected"

    - name: low_disk_space
      condition: disk_usage > 0.9
      severity: warning
      message: "Low disk space: {{ .disk_usage_percent }}% used"

    - name: metadata_extraction_failed
      condition: metadata_extraction_error
      severity: warning
      message: "Metadata extraction failed for {{ .file }}"
```

---

**Related Documentation:**
- [Getting Started](GETTING_STARTED.md)
- [Common Workflows](WORKFLOWS.md)
- [Troubleshooting](TROUBLESHOOTING.md)
- [API Reference](API.md)
- [Development Guide](DEVELOPMENT.md)
