# Cicada Configuration Guide

**Last Updated:** 2025-11-25

Complete guide to configuring Cicada for your lab's data commons.

## Table of Contents

1. [Configuration Overview](#configuration-overview)
2. [Configuration File Format](#configuration-file-format)
3. [Configuration Sections](#configuration-sections)
4. [AWS Configuration](#aws-configuration)
5. [Sync Configuration](#sync-configuration)
6. [Watch Configuration](#watch-configuration)
7. [Settings Configuration](#settings-configuration)
8. [Environment Variables](#environment-variables)
9. [Configuration Precedence](#configuration-precedence)
10. [Profile Management](#profile-management)
11. [Example Configurations](#example-configurations)
12. [Troubleshooting](#troubleshooting)

---

## Configuration Overview

Cicada uses a YAML configuration file to store settings for AWS credentials, sync behavior, watch configurations, and global settings.

**Default Location:** `~/.cicada/config.yaml`

**Configuration File Structure:**
```yaml
version: "1"
aws: { }
sync: { }
watches: [ ]
settings: { }
```

### Creating Configuration

**Initialize with defaults:**
```bash
cicada config init
```

**Manually create:**
```bash
mkdir -p ~/.cicada
touch ~/.cicada/config.yaml
```

### Viewing Configuration

```bash
# Show current configuration
cicada config show

# Show as JSON
cicada config show --format json

# Show as table
cicada config show --format table
```

---

## Configuration File Format

Cicada uses YAML format for configuration files.

**File Location:** `~/.cicada/config.yaml`

**Basic Structure:**
```yaml
# Configuration format version
version: "1"

# AWS settings for S3 access
aws:
  profile: default
  region: us-west-2
  endpoint: ""

# Default sync options
sync:
  concurrency: 4
  delete: false
  exclude:
    - .git/**
    - .DS_Store
    - "*.tmp"

# Watch configurations
watches:
  - id: microscopy-watch
    source: /data/microscopy
    destination: s3://lab-data/microscopy
    enabled: true

# Global settings
settings:
  verbose: false
  log_file: ~/.cicada/cicada.log
  check_updates: true
```

### YAML Syntax Notes

- Use 2 spaces for indentation (not tabs)
- Lists use `-` prefix
- Strings with special characters need quotes
- Boolean values: `true` or `false` (lowercase)
- Comments start with `#`

---

## Configuration Sections

### 1. Version

**Purpose:** Configuration format version for future compatibility

**Format:**
```yaml
version: "1"
```

**Values:**
- `"1"` - Current version (v0.1.0 - v0.2.0+)

---

### 2. AWS Configuration

Settings for AWS S3 access.

**Format:**
```yaml
aws:
  profile: default
  region: us-west-2
  endpoint: ""
```

**Fields:**

| Field | Type | Description | Default | Required |
|-------|------|-------------|---------|----------|
| `profile` | string | AWS profile from `~/.aws/credentials` | `default` | No |
| `region` | string | AWS region override | auto-detect | No |
| `endpoint` | string | Custom S3 endpoint (for S3-compatible services) | | No |

**See:** [AWS Configuration](#aws-configuration) section for details

---

### 3. Sync Configuration

Default sync behavior settings.

**Format:**
```yaml
sync:
  concurrency: 4
  delete: false
  exclude:
    - .git/**
    - .DS_Store
```

**Fields:**

| Field | Type | Description | Default | Required |
|-------|------|-------------|---------|----------|
| `concurrency` | int | Number of parallel transfers | `4` | No |
| `delete` | bool | Delete files not in source | `false` | No |
| `exclude` | []string | File patterns to exclude | `[]` | No |

**See:** [Sync Configuration](#sync-configuration) section for details

---

### 4. Watch Configuration

Directory watch configurations.

**Format:**
```yaml
watches:
  - id: microscopy-watch
    source: /data/microscopy
    destination: s3://lab-data/microscopy
    debounce_seconds: 10
    min_age_seconds: 60
    delete_source: false
    sync_on_start: true
    exclude:
      - "*.tmp"
    enabled: true
```

**Fields:**

| Field | Type | Description | Default | Required |
|-------|------|-------------|---------|----------|
| `id` | string | Unique watch identifier | | Yes |
| `source` | string | Directory to watch | | Yes |
| `destination` | string | Sync destination | | Yes |
| `debounce_seconds` | int | Delay after last event | `5` | No |
| `min_age_seconds` | int | Minimum file age before sync | `10` | No |
| `delete_source` | bool | Delete after successful sync | `false` | No |
| `sync_on_start` | bool | Initial sync on start | `true` | No |
| `exclude` | []string | Exclude patterns | `[]` | No |
| `enabled` | bool | Watch enabled/disabled | `true` | No |

**See:** [Watch Configuration](#watch-configuration) section for details

---

### 5. Settings Configuration

Global application settings.

**Format:**
```yaml
settings:
  verbose: false
  log_file: ~/.cicada/cicada.log
  check_updates: true
```

**Fields:**

| Field | Type | Description | Default | Required |
|-------|------|-------------|---------|----------|
| `verbose` | bool | Enable verbose logging | `false` | No |
| `log_file` | string | Log file path | | No |
| `check_updates` | bool | Check for updates on startup | `true` | No |

**See:** [Settings Configuration](#settings-configuration) section for details

---

## AWS Configuration

### AWS Profile

Cicada uses AWS profiles from `~/.aws/credentials` for authentication.

**AWS Credentials File Location:**
```
~/.aws/credentials
```

**Format:**
```ini
[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

[research-account]
aws_access_key_id = AKIAI44QH8DHBEXAMPLE
aws_secret_access_key = je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
```

**Cicada Configuration:**
```yaml
aws:
  profile: research-account  # Use 'research-account' profile
```

**Command-line Override:**
```bash
export AWS_PROFILE=research-account
cicada sync /data s3://bucket/data
```

---

### AWS Region

Specify AWS region for S3 operations.

**Configuration:**
```yaml
aws:
  region: us-west-2
```

**Supported Regions:**
- `us-east-1` - US East (N. Virginia)
- `us-west-2` - US West (Oregon)
- `eu-west-1` - Europe (Ireland)
- `eu-central-1` - Europe (Frankfurt)
- `ap-northeast-1` - Asia Pacific (Tokyo)
- [Full list](https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)

**Auto-detection:**
If not specified, Cicada will:
1. Check `AWS_REGION` environment variable
2. Check `~/.aws/config` for default region
3. Auto-detect from bucket location

**Command-line Override:**
```bash
export AWS_REGION=eu-west-1
cicada sync /data s3://bucket/data
```

---

### Custom S3 Endpoint

Use S3-compatible services (MinIO, DigitalOcean Spaces, Wasabi, etc.)

**Configuration:**
```yaml
aws:
  profile: default
  endpoint: https://nyc3.digitaloceanspaces.com  # DigitalOcean Spaces
  region: nyc3
```

**Examples:**

**MinIO:**
```yaml
aws:
  endpoint: http://localhost:9000
  region: us-east-1
```

**DigitalOcean Spaces:**
```yaml
aws:
  endpoint: https://nyc3.digitaloceanspaces.com
  region: nyc3
```

**Wasabi:**
```yaml
aws:
  endpoint: https://s3.wasabisys.com
  region: us-east-1
```

**Backblaze B2:**
```yaml
aws:
  endpoint: https://s3.us-west-000.backblazeb2.com
  region: us-west-000
```

---

## Sync Configuration

### Concurrency

Number of parallel file transfers.

**Configuration:**
```yaml
sync:
  concurrency: 4  # 4 parallel transfers
```

**Recommendations:**

| Use Case | Concurrency | Reasoning |
|----------|-------------|-----------|
| Small files (<1MB) | 8-16 | Higher concurrency reduces overhead |
| Large files (>100MB) | 2-4 | Bandwidth-limited, fewer transfers |
| Limited bandwidth | 2-4 | Prevent network saturation |
| High-speed network | 8-16 | Maximize throughput |
| Many small files | 16-32 | Parallelize I/O operations |

**Command-line Override:**
```bash
cicada sync --concurrency 8 /data s3://bucket/data
```

---

### Delete Behavior

Whether to delete files in destination not present in source.

**Configuration:**
```yaml
sync:
  delete: false  # Do not delete (safe default)
```

**Values:**
- `true` - Delete files in destination not in source (mirror sync)
- `false` - Keep all files in destination (additive sync)

**⚠️ Warning:** Enabling `delete: true` can result in data loss if source is incomplete or corrupted.

**Use Cases:**

**Additive Sync (delete: false):**
- Backup/archive scenarios
- Multi-source to single destination
- Incremental data collection

**Mirror Sync (delete: true):**
- Exact replica required
- Source is authoritative
- Space management needed

**Command-line Override:**
```bash
cicada sync --delete /data s3://bucket/data
```

---

### Exclude Patterns

File patterns to exclude from sync.

**Configuration:**
```yaml
sync:
  exclude:
    - .git/**           # Git directories
    - .DS_Store         # macOS metadata
    - "*.tmp"           # Temporary files
    - "*.swp"           # Vim swap files
    - "*.bak"           # Backup files
    - Thumbs.db         # Windows thumbnails
    - "~*"              # Editor temp files
```

**Pattern Syntax:**

| Pattern | Matches | Example |
|---------|---------|---------|
| `*.txt` | Files ending in .txt | `file.txt`, `data.txt` |
| `temp*` | Files starting with temp | `temp.log`, `temp123` |
| `dir/` | Specific directory | `dir/file.txt` |
| `**/cache/` | cache directory anywhere | `a/cache/`, `b/c/cache/` |
| `**/*.log` | .log files anywhere | `a/b/c/file.log` |
| `[Tt]emp*` | Temp or temp prefix | `Temp.txt`, `temp.log` |

**Examples:**

**Exclude all hidden files:**
```yaml
sync:
  exclude:
    - ".*"         # Hidden files/dirs
    - "**/.*"      # Hidden anywhere
```

**Exclude by file type:**
```yaml
sync:
  exclude:
    - "*.tmp"
    - "*.cache"
    - "*.lock"
```

**Exclude by directory:**
```yaml
sync:
  exclude:
    - build/**
    - dist/**
    - node_modules/**
```

---

## Watch Configuration

### Watch ID

Unique identifier for each watch.

**Format:**
```yaml
watches:
  - id: microscopy-lab-01
```

**ID Guidelines:**
- Use descriptive names
- Include instrument or location
- Avoid spaces (use hyphens or underscores)
- Must be unique across all watches

**Examples:**
- `microscopy-confocal-01`
- `sequencer-illumina-novaseq`
- `mass-spec-proteomics`
- `instrument-room-201`

---

### Watch Source and Destination

**Configuration:**
```yaml
watches:
  - id: watch-01
    source: /data/microscopy          # Local directory
    destination: s3://lab-data/microscopy  # S3 bucket
```

**Source:**
- Must be local directory path
- Will be monitored for file changes
- Use absolute paths for clarity

**Destination:**
- Can be local path or S3 URI
- Files will be synced here
- Created if doesn't exist

---

### Debounce Delay

Time to wait after last file change before syncing.

**Configuration:**
```yaml
watches:
  - debounce_seconds: 10  # Wait 10 seconds
```

**Purpose:**
- Prevent syncing partial writes
- Coalesce rapid changes
- Reduce sync frequency

**Recommendations:**

| Scenario | Delay | Reasoning |
|----------|-------|-----------|
| Small files, quick writes | 5-10s | Files complete quickly |
| Large files, slow writes | 30-60s | Files take time to write |
| Batch processing | 60-300s | Wait for batch completion |
| Real-time needs | 1-5s | Minimize latency |

---

### Minimum File Age

Only sync files older than this duration.

**Configuration:**
```yaml
watches:
  - min_age_seconds: 60  # Only sync files >60s old
```

**Purpose:**
- Ensure files are completely written
- Prevent syncing incomplete data
- Wait for application to close file

**Recommendations:**

| File Size | Min Age | Reasoning |
|-----------|---------|-----------|
| <10MB | 10-30s | Quick writes |
| 10-100MB | 30-60s | Moderate writes |
| >100MB | 60-300s | Slow writes |
| Database dumps | 300-600s | Very slow, ensure completion |

**Example:**
```yaml
# Large microscopy files (100MB-1GB)
watches:
  - id: confocal-imaging
    min_age_seconds: 120  # Wait 2 minutes
```

---

### Delete Source After Sync

Delete source files after successful sync (move instead of copy).

**Configuration:**
```yaml
watches:
  - delete_source: false  # Keep source files
```

**Values:**
- `true` - Delete source after successful sync (MOVE operation)
- `false` - Keep source files (COPY operation)

**⚠️ Warning:** Enabling `delete_source: true` will DELETE source files. Use only when:
- Source is temporary/staging directory
- Destination is reliable (S3 with versioning)
- You have backups
- You understand the risk

**Use Cases:**

**Keep Source (delete_source: false):**
- Local analysis still needed
- Backup scenario
- Multiple destinations
- Source is primary copy

**Delete Source (delete_source: true):**
- Limited local storage
- Automatic cleanup needed
- Source is temporary
- Cloud is primary storage

---

### Sync On Start

Perform initial sync when watch starts.

**Configuration:**
```yaml
watches:
  - sync_on_start: true  # Sync existing files on start
```

**Values:**
- `true` - Sync existing files when watch starts
- `false` - Only sync new changes after watch starts

**Use Cases:**

**Sync On Start (true):**
- Catch up on missed changes
- Ensure destination is current
- Initial setup of watch

**Don't Sync On Start (false):**
- Only watch for new changes
- Destination already current
- Avoid initial sync overhead

---

### Watch Exclude Patterns

Exclude specific files from this watch.

**Configuration:**
```yaml
watches:
  - id: watch-01
    exclude:
      - "*.tmp"
      - ".DS_Store"
      - "*.lock"
```

**Inherits:** Global `sync.exclude` patterns are also applied

**Watch-specific overrides:**
```yaml
sync:
  exclude:
    - "*.tmp"      # Global exclude

watches:
  - id: watch-01
    exclude:
      - "*.draft"  # Additional exclude for this watch
```

---

### Watch Enable/Disable

Enable or disable a watch without removing it.

**Configuration:**
```yaml
watches:
  - id: microscopy-watch
    enabled: true  # Watch is active

  - id: old-sequencer
    enabled: false  # Watch is disabled
```

**Use Cases:**
- Temporarily disable watch
- Maintenance periods
- Seasonal instruments
- Testing configurations

---

## Settings Configuration

### Verbose Logging

Enable detailed logging output.

**Configuration:**
```yaml
settings:
  verbose: true
```

**Effect:**
- Shows all operations
- Displays file paths
- Reports progress
- Includes timing information

**Command-line Override:**
```bash
cicada --verbose sync /data s3://bucket/data
```

---

### Log File

Write logs to file.

**Configuration:**
```yaml
settings:
  log_file: ~/.cicada/cicada.log
```

**Log Rotation:**
Cicada does not automatically rotate logs. Use external tools:

**Linux/macOS (logrotate):**
```
~/.cicada/cicada.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
}
```

**Manual rotation:**
```bash
# Rotate logs manually
mv ~/.cicada/cicada.log ~/.cicada/cicada.log.$(date +%Y%m%d)
gzip ~/.cicada/cicada.log.*
```

---

### Check Updates

Check for new Cicada versions on startup.

**Configuration:**
```yaml
settings:
  check_updates: true
```

**Behavior:**
- Checks GitHub releases on startup
- Non-blocking (doesn't delay commands)
- Displays notification if update available

**Disable:**
```yaml
settings:
  check_updates: false
```

---

## Environment Variables

Environment variables override configuration file values.

### Cicada Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CICADA_CONFIG` | Path to config file | `/etc/cicada/config.yaml` |
| `CICADA_VERBOSE` | Enable verbose output | `true` |

### AWS Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `AWS_PROFILE` | AWS profile name | `research-account` |
| `AWS_REGION` | AWS region | `us-west-2` |
| `AWS_ACCESS_KEY_ID` | AWS access key | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | (secret) |
| `AWS_SESSION_TOKEN` | Temporary session token | (for assumed roles) |
| `AWS_ENDPOINT_URL` | Custom S3 endpoint | `https://s3.wasabisys.com` |

### System Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `HOME` | User home directory | `/home/user` |
| `TMPDIR` | Temporary directory | `/tmp` |

---

## Configuration Precedence

Configuration values are resolved in this order (highest to lowest):

1. **Command-line flags**
   ```bash
   cicada sync --concurrency 8 /data s3://bucket/data
   ```

2. **Environment variables**
   ```bash
   export AWS_PROFILE=research-account
   cicada sync /data s3://bucket/data
   ```

3. **Configuration file**
   ```yaml
   aws:
     profile: default
   ```

4. **Built-in defaults**
   ```
   concurrency: 4
   delete: false
   ```

**Example:**
```yaml
# config.yaml
sync:
  concurrency: 4
```

```bash
# Environment variable overrides config
export CICADA_CONCURRENCY=8

# Command-line flag overrides both
cicada sync --concurrency 16 /data s3://bucket/data

# Effective value: 16 (command-line wins)
```

---

## Profile Management

### Multiple Environments

Manage different configurations for different environments.

**Directory Structure:**
```
~/.cicada/
├── config.yaml              # Default
├── config.production.yaml   # Production
├── config.development.yaml  # Development
└── config.testing.yaml      # Testing
```

**Usage:**
```bash
# Use production config
cicada --config ~/.cicada/config.production.yaml sync /data s3://prod-bucket/data

# Use development config
cicada --config ~/.cicada/config.development.yaml sync /data s3://dev-bucket/data
```

**Environment variable:**
```bash
export CICADA_CONFIG=~/.cicada/config.production.yaml
cicada sync /data s3://prod-bucket/data
```

---

### Lab-Specific Configurations

**Scenario:** Multiple labs sharing compute resources

**Lab A Configuration:**
```yaml
# /etc/cicada/config.lab-a.yaml
version: "1"
aws:
  profile: lab-a-account
sync:
  concurrency: 4
watches:
  - id: lab-a-microscope
    source: /data/lab-a/microscopy
    destination: s3://lab-a-data/microscopy
```

**Lab B Configuration:**
```yaml
# /etc/cicada/config.lab-b.yaml
version: "1"
aws:
  profile: lab-b-account
sync:
  concurrency: 8
watches:
  - id: lab-b-sequencer
    source: /data/lab-b/sequencing
    destination: s3://lab-b-data/sequencing
```

**Usage:**
```bash
# Lab A operations
cicada --config /etc/cicada/config.lab-a.yaml sync /data/lab-a s3://lab-a-data/

# Lab B operations
cicada --config /etc/cicada/config.lab-b.yaml sync /data/lab-b s3://lab-b-data/
```

---

## Example Configurations

### Minimal Configuration

```yaml
version: "1"

aws:
  profile: default

sync:
  concurrency: 4
  delete: false

watches: []

settings:
  verbose: false
```

---

### Small Lab Configuration

```yaml
version: "1"

# Use default AWS account
aws:
  profile: default
  region: us-west-2

# Conservative sync settings
sync:
  concurrency: 4
  delete: false
  exclude:
    - .git/**
    - .DS_Store
    - "*.tmp"
    - "*.swp"

# Single microscope watch
watches:
  - id: confocal-microscope
    source: /data/microscopy
    destination: s3://lab-backups/microscopy
    debounce_seconds: 10
    min_age_seconds: 60
    delete_source: false
    sync_on_start: true
    enabled: true

settings:
  verbose: false
  log_file: ~/.cicada/cicada.log
  check_updates: true
```

---

### Multi-Instrument Lab Configuration

```yaml
version: "1"

aws:
  profile: research-account
  region: us-east-1

sync:
  concurrency: 8
  delete: false
  exclude:
    - .git/**
    - .DS_Store
    - "*.tmp"
    - "~*"

# Multiple instrument watches
watches:
  # Confocal microscope
  - id: confocal-zeiss-lsm880
    source: /data/confocal
    destination: s3://lab-data/microscopy/confocal
    debounce_seconds: 15
    min_age_seconds: 120
    enabled: true

  # Illumina sequencer
  - id: sequencer-novaseq-6000
    source: /data/sequencing
    destination: s3://lab-data/sequencing/illumina
    debounce_seconds: 30
    min_age_seconds: 300
    enabled: true

  # Mass spectrometer
  - id: mass-spec-q-exactive
    source: /data/mass-spec
    destination: s3://lab-data/proteomics
    debounce_seconds: 20
    min_age_seconds: 180
    enabled: true

  # Flow cytometer
  - id: flow-cytometer-bd-aria
    source: /data/flow-cytometry
    destination: s3://lab-data/flow
    debounce_seconds: 10
    min_age_seconds: 60
    enabled: true

settings:
  verbose: true
  log_file: /var/log/cicada/cicada.log
  check_updates: true
```

---

### High-Performance Configuration

```yaml
version: "1"

aws:
  profile: default
  region: us-west-2

# Aggressive performance settings
sync:
  concurrency: 16         # High parallelism
  delete: false
  exclude:
    - "*.tmp"

watches:
  - id: high-throughput-sequencing
    source: /nvme/fast-storage/sequencing
    destination: s3://archive/sequencing
    debounce_seconds: 5   # Quick response
    min_age_seconds: 30   # Minimal delay
    delete_source: true   # Clean up local storage
    sync_on_start: false  # Only new data
    enabled: true

settings:
  verbose: false          # Reduce overhead
  log_file: /dev/null     # No logging
  check_updates: false    # No network calls
```

---

### Development/Testing Configuration

```yaml
version: "1"

aws:
  profile: default
  endpoint: http://localhost:9000  # Local MinIO
  region: us-east-1

sync:
  concurrency: 2
  delete: false

watches:
  - id: test-watch
    source: /tmp/test-data
    destination: s3://test-bucket/data
    debounce_seconds: 2
    min_age_seconds: 5
    enabled: true

settings:
  verbose: true
  log_file: /tmp/cicada-test.log
  check_updates: false
```

---

## Troubleshooting

### Configuration Not Loading

**Symptoms:**
- Settings from config file not applied
- Using default values instead

**Solutions:**

1. **Check file location:**
   ```bash
   ls -la ~/.cicada/config.yaml
   ```

2. **Verify YAML syntax:**
   ```bash
   # Use a YAML validator
   python3 -c "import yaml; yaml.safe_load(open('~/.cicada/config.yaml'))"
   ```

3. **Check file permissions:**
   ```bash
   chmod 644 ~/.cicada/config.yaml
   ```

4. **Use explicit path:**
   ```bash
   cicada --config ~/.cicada/config.yaml sync /data s3://bucket/data
   ```

---

### AWS Credentials Not Working

**Symptoms:**
- "AccessDenied" errors
- "InvalidAccessKeyId" errors
- "SignatureDoesNotMatch" errors

**Solutions:**

1. **Verify credentials file:**
   ```bash
   cat ~/.aws/credentials
   ```

2. **Check profile name:**
   ```bash
   # List available profiles
   grep '^\[' ~/.aws/credentials
   ```

3. **Test AWS CLI:**
   ```bash
   aws s3 ls --profile default
   ```

4. **Check environment variables:**
   ```bash
   echo $AWS_PROFILE
   echo $AWS_REGION
   ```

5. **Verify IAM permissions:**
   - Ensure user/role has S3 permissions
   - Required: `s3:ListBucket`, `s3:GetObject`, `s3:PutObject`

---

### Watches Not Starting

**Symptoms:**
- Watches don't start automatically
- "watch not found" errors

**Solutions:**

1. **Check watch configuration:**
   ```bash
   cicada config show
   ```

2. **Verify enabled flag:**
   ```yaml
   watches:
     - id: my-watch
       enabled: true  # Must be true
   ```

3. **Check source directory exists:**
   ```bash
   ls -la /data/microscopy
   ```

4. **Verify permissions:**
   ```bash
   # Cicada needs read access to source
   ls -la /data/microscopy
   ```

---

### Exclude Patterns Not Working

**Symptoms:**
- Files still syncing despite exclude patterns
- Wrong files excluded

**Solutions:**

1. **Check pattern syntax:**
   ```yaml
   sync:
     exclude:
       - "*.tmp"      # Correct (quotes)
       - *.tmp        # May not work (no quotes)
   ```

2. **Test patterns:**
   ```bash
   # Dry run to see what would sync
   cicada sync --dry-run /data s3://bucket/data
   ```

3. **Use absolute patterns:**
   ```yaml
   exclude:
     - "**/cache/*"   # Match cache anywhere
   ```

---

## Related Documentation

- [CLI Reference](CLI_REFERENCE.md) - Command-line interface
- [User Guide](USER_GUIDE.md) - Getting started
- [Architecture](ARCHITECTURE.md) - System architecture

---

**Contributing:** Found an error or want to suggest improvements? See [CONTRIBUTING.md](CONTRIBUTING.md)

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
