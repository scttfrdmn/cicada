# Storage Backends Guide

**Last Updated:** 2025-01-24

Complete guide to Cicada's storage backend system for local filesystem and AWS S3.

## Table of Contents

1. [Storage Overview](#storage-overview)
2. [Local Filesystem Backend](#local-filesystem-backend)
3. [AWS S3 Backend](#aws-s3-backend)
4. [Backend Configuration](#backend-configuration)
5. [Performance Characteristics](#performance-characteristics)
6. [Cost Analysis](#cost-analysis)
7. [Migration Between Backends](#migration-between-backends)
8. [Backup and Recovery](#backup-and-recovery)
9. [S3 Object Tagging](#s3-object-tagging)
10. [Storage Classes](#storage-classes)
11. [Best Practices](#best-practices)
12. [Troubleshooting](#troubleshooting)

---

## Storage Overview

Cicada provides a unified storage abstraction that works consistently across different backend types.

### Supported Backends

| Backend | Status | Use Case |
|---------|--------|----------|
| **Local Filesystem** | ✅ Available | Local storage, NFS mounts, network shares |
| **AWS S3** | ✅ Available | Cloud storage, archival, backup |

### Backend Interface

All backends implement a common interface:

```go
type Backend interface {
    List(ctx context.Context, prefix string) ([]FileInfo, error)
    Read(ctx context.Context, path string) (io.ReadCloser, error)
    Write(ctx context.Context, path string, r io.Reader, size int64) error
    Delete(ctx context.Context, path string) error
    Stat(ctx context.Context, path string) (*FileInfo, error)
    Close() error
}
```

This ensures consistent behavior regardless of storage type.

---

## Local Filesystem Backend

### Overview

The local backend uses standard filesystem operations for data storage.

**Use Cases:**
- Local data analysis
- NFS/CIFS mounted network storage
- Staging area for cloud upload
- High-performance local storage (NVMe, SSD)

### Path Format

```bash
# Absolute paths
/data/microscopy
/home/user/lab-data
/mnt/nfs/shared

# Relative paths (not recommended)
./data
../shared/data
```

**Best Practice:** Always use absolute paths for clarity and reliability.

### Configuration

No special configuration required. Cicada uses the local filesystem directly.

**Example:**
```bash
# Sync from local to local
cicada sync /data/source /backup/destination

# Sync local to S3
cicada sync /data/lab s3://bucket/lab-data
```

### Features

| Feature | Support | Notes |
|---------|---------|-------|
| **Read** | ✅ | Standard file read operations |
| **Write** | ✅ | Standard file write operations |
| **Delete** | ✅ | Standard file deletion |
| **List** | ✅ | Directory listing with recursion |
| **Stat** | ✅ | File metadata (size, mtime) |
| **ETag/Checksum** | ✅ | MD5 hash computed on demand |
| **Streaming** | ✅ | No memory buffering of entire file |
| **Concurrent Access** | ✅ | OS-level locking |

### Performance

**Throughput:**
- Limited by disk I/O speed
- NVMe SSD: 3-7 GB/s read, 2-5 GB/s write
- SATA SSD: 500-600 MB/s read/write
- HDD: 100-200 MB/s read/write
- NFS: 100-1000 MB/s (network-dependent)

**Latency:**
- NVMe: <100µs
- SSD: ~1ms
- HDD: 5-15ms
- NFS: 1-10ms (network-dependent)

**Concurrency:**
- Limited by filesystem and OS
- Generally excellent for reads
- Write contention depends on filesystem

### Limitations

- No built-in versioning
- No automatic redundancy
- Local hardware failures affect availability
- Storage capacity limited by hardware

### Best Practices

**1. Use Fast Storage for Active Work**
```bash
# Keep working data on fast local storage
/data/active-projects/     # NVMe SSD
/data/archive/             # Slower HDD or cloud
```

**2. Regular Backups**
```bash
# Daily backup to S3
0 2 * * * cicada sync /data/important s3://backups/daily/$(date +\%Y-\%m-\%d)
```

**3. Monitor Disk Space**
```bash
# Check available space
df -h /data

# Set up alerts when <10% free
```

**4. Use Network Storage Carefully**
```bash
# NFS mounts may be slower
cicada sync --concurrency 4 /mnt/nfs/data s3://bucket/data

# Consider caching layer for frequent access
```

---

## AWS S3 Backend

### Overview

The S3 backend provides cloud storage with high durability, availability, and scalability.

**Key Benefits:**
- 99.999999999% (11 9's) durability
- Unlimited storage capacity
- Pay-per-use pricing
- Geographic redundancy
- Versioning support
- Lifecycle policies

**Use Cases:**
- Long-term data archival
- Backup and disaster recovery
- Collaborative data sharing
- Cost-effective bulk storage
- Compliance and regulatory requirements

### S3 URI Format

```
s3://bucket-name/prefix/path/to/file
```

**Components:**
- `s3://` - Protocol identifier
- `bucket-name` - S3 bucket name (globally unique)
- `prefix/path/to/file` - Object key (path within bucket)

**Examples:**
```bash
s3://lab-data/                         # Entire bucket
s3://lab-data/microscopy/              # Prefix
s3://lab-data/microscopy/2025/         # Nested prefix
s3://lab-data/microscopy/2025/exp001.czi  # Specific object
```

### Bucket Creation

**AWS CLI:**
```bash
# Create bucket
aws s3 mb s3://lab-data --region us-west-2

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket lab-data \
  --versioning-configuration Status=Enabled

# Set lifecycle policy (optional)
aws s3api put-bucket-lifecycle-configuration \
  --bucket lab-data \
  --lifecycle-configuration file://lifecycle.json
```

**Cicada automatically creates buckets:** No (must exist before use)

### Configuration

**AWS Credentials:**
```bash
# ~/.aws/credentials
[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

[research-account]
aws_access_key_id = AKIAI44QH8DHBEXAMPLE
aws_secret_access_key = je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
```

**Cicada Configuration:**
```yaml
# ~/.cicada/config.yaml
aws:
  profile: research-account
  region: us-west-2
```

**Environment Variables:**
```bash
export AWS_PROFILE=research-account
export AWS_REGION=us-west-2
```

### Features

| Feature | Support | Notes |
|---------|---------|-------|
| **Read** | ✅ | GetObject API |
| **Write** | ✅ | PutObject API (multipart for large files) |
| **Delete** | ✅ | DeleteObject API |
| **List** | ✅ | ListObjectsV2 API with pagination |
| **Stat** | ✅ | HeadObject API |
| **ETag/Checksum** | ✅ | Native S3 ETag |
| **Streaming** | ✅ | Chunked upload/download |
| **Versioning** | ✅ | S3 versioning (if enabled) |
| **Encryption** | ✅ | Server-side encryption (SSE-S3, SSE-KMS) |
| **Tagging** | ✅ | Object tags for metadata |
| **Storage Classes** | ✅ | STANDARD, IA, GLACIER, etc. |

### Performance

**Throughput:**
- Single connection: 50-100 MB/s
- Multiple connections: 1-10 GB/s (aggregate)
- Scales with concurrency

**Latency:**
- First byte: 100-200ms (depends on region)
- Subsequent bytes: streaming at line rate

**Concurrency:**
- Highly parallelizable
- No contention issues
- Scales to thousands of requests/second

**Request Pricing:**
- GET: $0.0004 per 1,000 requests
- PUT: $0.005 per 1,000 requests
- LIST: $0.005 per 1,000 requests

### Multipart Upload

For files >5GB or improved performance:

**Automatic in Cicada:**
- Files >100MB use multipart upload automatically
- 5MB part size (configurable)
- Parallel part uploads
- Resume on failure

**Manual (AWS CLI):**
```bash
# Initiate multipart upload
aws s3api create-multipart-upload \
  --bucket lab-data \
  --key large-file.dat

# Upload parts (parallel)
aws s3api upload-part \
  --bucket lab-data \
  --key large-file.dat \
  --part-number 1 \
  --body part1.dat \
  --upload-id <upload-id>

# Complete upload
aws s3api complete-multipart-upload \
  --bucket lab-data \
  --key large-file.dat \
  --upload-id <upload-id>
```

### Encryption

**Server-Side Encryption (SSE):**

**SSE-S3 (Default):**
```bash
# Enabled by default in Cicada
# Uses AWS-managed keys
# No configuration needed
```

**SSE-KMS:**
```yaml
# Future support planned
aws:
  s3:
    encryption: kms
    kms_key_id: arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012
```

**SSE-C (Customer Keys):**
```
# Not supported (requires key management)
```

**Client-Side Encryption:**
```bash
# Encrypt before upload
openssl enc -aes-256-cbc -in file.dat -out file.dat.enc
cicada sync file.dat.enc s3://bucket/encrypted/

# Decrypt after download
openssl enc -d -aes-256-cbc -in file.dat.enc -out file.dat
```

### IAM Permissions

**Minimum Required Permissions:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetBucketLocation"
      ],
      "Resource": "arn:aws:s3:::lab-data"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:GetObjectTagging",
        "s3:PutObjectTagging"
      ],
      "Resource": "arn:aws:s3:::lab-data/*"
    }
  ]
}
```

**Read-Only Access:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": "arn:aws:s3:::lab-data"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::lab-data/*"
    }
  ]
}
```

**Prefix-Limited Access:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::lab-data",
      "Condition": {
        "StringLike": {
          "s3:prefix": "user-alice/*"
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::lab-data/user-alice/*"
    }
  ]
}
```

---

## Backend Configuration

### Path Resolution

Cicada automatically detects backend type:

```bash
# Local filesystem (no prefix)
/data/lab
./relative/path
~/home/user/data

# S3 (s3:// prefix)
s3://bucket/path
s3://bucket/
s3://bucket
```

### Mixed Backend Operations

**Sync between different backends:**

```bash
# Local to S3
cicada sync /data/microscopy s3://lab-data/microscopy

# S3 to local
cicada sync s3://lab-data/microscopy /data/microscopy

# S3 to S3 (same or different accounts)
cicada sync s3://source-bucket/data s3://dest-bucket/data

# Local to local (copy)
cicada sync /data/source /backup/destination
```

---

## Performance Characteristics

### Throughput Comparison

| Backend | Read Speed | Write Speed | Latency |
|---------|------------|-------------|---------|
| **NVMe Local** | 3-7 GB/s | 2-5 GB/s | <1ms |
| **SATA SSD Local** | 500-600 MB/s | 500-600 MB/s | ~1ms |
| **HDD Local** | 100-200 MB/s | 100-200 MB/s | 10-15ms |
| **NFS (Gigabit)** | 100-125 MB/s | 100-125 MB/s | 1-5ms |
| **NFS (10 Gbit)** | 500-1200 MB/s | 500-1200 MB/s | 1-5ms |
| **S3 (single)** | 50-100 MB/s | 50-100 MB/s | 100-200ms |
| **S3 (parallel x8)** | 400-800 MB/s | 400-800 MB/s | 100-200ms |
| **S3 (parallel x16)** | 800-1600 MB/s | 800-1600 MB/s | 100-200ms |

### Optimization Strategies

**For Small Files (<1MB):**
```bash
# Increase concurrency
cicada sync --concurrency 16 /data/many-small-files s3://bucket/data

# Batch operations reduce overhead
```

**For Large Files (>100MB):**
```bash
# Moderate concurrency (bandwidth-limited)
cicada sync --concurrency 4 /data/large-files s3://bucket/data

# Multipart upload automatically used
```

**For Mixed Workloads:**
```bash
# Balance concurrency
cicada sync --concurrency 8 /data/mixed s3://bucket/data
```

### Benchmark Examples

**Local NVMe to S3 (1000 files, 10MB each):**
```
Concurrency: 8
Total: 10 GB
Time: 120 seconds
Throughput: 85 MB/s
```

**S3 to Local SSD (100 files, 100MB each):**
```
Concurrency: 4
Total: 10 GB
Time: 180 seconds
Throughput: 56 MB/s
```

**Local to Local (10,000 files, 1MB each):**
```
Concurrency: 16
Total: 10 GB
Time: 15 seconds
Throughput: 680 MB/s
```

---

## Cost Analysis

### AWS S3 Pricing (US East, as of 2025)

**Storage:**
| Storage Class | Price (per GB/month) | Use Case |
|---------------|---------------------|----------|
| Standard | $0.023 | Frequently accessed |
| Intelligent-Tiering | $0.023 + $0.0025 (monitoring) | Unknown access patterns |
| Standard-IA | $0.0125 | Infrequently accessed |
| One Zone-IA | $0.01 | Infrequent, single-AZ |
| Glacier Instant | $0.004 | Archive, instant retrieval |
| Glacier Flexible | $0.0036 | Archive, 1-5 min retrieval |
| Glacier Deep Archive | $0.00099 | Long-term archive, 12-hour retrieval |

**Requests:**
- PUT/POST: $0.005 per 1,000 requests
- GET: $0.0004 per 1,000 requests
- LIST: $0.005 per 1,000 requests

**Data Transfer:**
- Upload: Free
- Download: $0.09/GB (first 10TB/month)
- Same-region transfer: Free

### Cost Example

**Scenario:** 10TB microscopy data, monthly access

**Option 1: Standard Storage**
```
Storage: 10,000 GB × $0.023 = $230/month
Requests: 100,000 GET × $0.0004/1000 = $0.04/month
Transfer: 100 GB × $0.09 = $9/month
Total: $239.04/month
Annual: $2,868.48
```

**Option 2: Standard-IA Storage**
```
Storage: 10,000 GB × $0.0125 = $125/month
Requests: 100,000 GET × $0.001/1000 = $0.10/month
Retrieval: 100 GB × $0.01 = $1/month
Transfer: 100 GB × $0.09 = $9/month
Total: $135.10/month
Annual: $1,621.20
Savings: $1,247.28/year (43% less)
```

**Option 3: Glacier Deep Archive**
```
Storage: 10,000 GB × $0.00099 = $9.90/month
Requests: 100 GET × $0.05/1000 = $0.005/month
Retrieval: 100 GB × $0.02 = $2/month (12-hour)
Transfer: 100 GB × $0.09 = $9/month
Total: $20.905/month
Annual: $250.86
Savings: $2,617.62/year (91% less)
```

### Cost Optimization

**1. Use Lifecycle Policies:**
```json
{
  "Rules": [
    {
      "Id": "Archive old data",
      "Status": "Enabled",
      "Filter": {
        "Prefix": "microscopy/"
      },
      "Transitions": [
        {
          "Days": 90,
          "StorageClass": "STANDARD_IA"
        },
        {
          "Days": 365,
          "StorageClass": "GLACIER"
        }
      ]
    }
  ]
}
```

**2. Intelligent-Tiering:**
```bash
# Enable on bucket
aws s3api put-bucket-intelligent-tiering-configuration \
  --bucket lab-data \
  --id auto-tier \
  --intelligent-tiering-configuration file://tiering-config.json
```

**3. Request Optimization:**
```bash
# List once, cache results
aws s3 ls s3://lab-data/ --recursive > file-list.txt

# Use prefix filters
cicada sync /data s3://bucket/specific-prefix/
```

**4. Compression:**
```bash
# Compress before upload (if applicable)
tar -czf archive.tar.gz /data/text-files
cicada sync archive.tar.gz s3://bucket/compressed/

# Savings: 70-90% for text/log files
```

---

## Migration Between Backends

### Local to S3

```bash
# Initial upload
cicada sync /data/lab s3://lab-data/initial-upload/

# Keep in sync
cicada watch add /data/lab s3://lab-data/live/
```

### S3 to Local

```bash
# Download all data
cicada sync s3://lab-data/archive /local/archive

# Selective download
cicada sync s3://lab-data/2025/ /local/2025-data
```

### S3 to S3 (Cross-Region)

```bash
# Migrate to different region
cicada sync s3://us-east-bucket/data s3://eu-west-bucket/data

# Note: Cross-region transfer fees apply
```

### S3 to S3 (Cross-Account)

```yaml
# Source account config
aws:
  profile: source-account

# Destination account config
aws:
  profile: dest-account
```

```bash
# Two-step migration
cicada sync s3://source-bucket/data /tmp/migration
cicada --config dest-config.yaml sync /tmp/migration s3://dest-bucket/data
```

### Storage Class Migration

**Using AWS CLI:**
```bash
# Copy with new storage class
aws s3 cp s3://bucket/old/ s3://bucket/new/ \
  --storage-class GLACIER \
  --recursive
```

**Using Lifecycle Policy:**
```json
{
  "Rules": [
    {
      "Id": "Move to IA after 30 days",
      "Status": "Enabled",
      "Transitions": [
        {
          "Days": 30,
          "StorageClass": "STANDARD_IA"
        }
      ]
    }
  ]
}
```

---

## Backup and Recovery

### Backup Strategies

**Strategy 1: 3-2-1 Backup**
- 3 copies of data
- 2 different media types
- 1 off-site copy

**Implementation:**
```bash
# Primary: Local NVMe
/data/primary/

# Backup 1: Local HDD
cicada sync /data/primary /backup/local-hdd/

# Backup 2: S3 (off-site)
cicada sync /data/primary s3://backups/off-site/
```

**Strategy 2: Incremental Backup**
```bash
# Daily incrementals
0 2 * * * cicada sync /data s3://backups/daily/$(date +\%Y-\%m-\%d)

# Weekly full backup
0 2 * * 0 cicada sync /data s3://backups/weekly/$(date +\%Y-W\%W)
```

**Strategy 3: Versioned Backup**
```bash
# Enable S3 versioning
aws s3api put-bucket-versioning \
  --bucket backups \
  --versioning-configuration Status=Enabled

# Sync (versions preserved)
cicada sync /data s3://backups/versioned/
```

### Recovery Procedures

**Full Recovery:**
```bash
# Download entire backup
cicada sync s3://backups/latest/ /data/recovery/
```

**Point-in-Time Recovery:**
```bash
# Restore from specific date
cicada sync s3://backups/daily/2025-01-15/ /data/recovery/
```

**Selective Recovery:**
```bash
# Restore specific files
cicada sync s3://backups/latest/experiment-001/ /data/recovery/experiment-001/
```

**Version Recovery (S3):**
```bash
# List versions
aws s3api list-object-versions --bucket backups --prefix file.dat

# Restore specific version
aws s3api get-object \
  --bucket backups \
  --key file.dat \
  --version-id <version-id> \
  file.dat.recovered
```

---

## S3 Object Tagging

### Metadata Storage

Cicada can store metadata as S3 object tags (up to 10 tags per object).

**Example Tags:**
```
format: CZI
instrument: Zeiss-LSM-880
date: 2025-01-24
sample_id: EXP001
operator: jsmith
```

**Implementation:**
```bash
# Extract metadata and tag S3 objects
cicada metadata extract file.czi | \
  cicada s3 tag s3://bucket/file.czi --from-metadata
```

### Tag Limitations

**S3 Tag Constraints:**
- Maximum 10 tags per object
- Key length: 128 characters max
- Value length: 256 characters max
- Case-sensitive
- UTF-8 encoding

**Workaround for Large Metadata:**
- Use sidecar JSON files: `file.czi.metadata.json`
- Store summary in tags, full metadata in sidecar

---

## Storage Classes

### S3 Storage Classes

| Class | Availability | Min Duration | Retrieval | Use Case |
|-------|--------------|--------------|-----------|----------|
| **Standard** | 99.99% | None | Instant | Frequent access |
| **Intelligent-Tiering** | 99.9% | None | Instant | Unknown pattern |
| **Standard-IA** | 99.9% | 30 days | Instant | Infrequent access |
| **One Zone-IA** | 99.5% | 30 days | Instant | Infrequent, non-critical |
| **Glacier Instant** | 99.9% | 90 days | Instant | Archive, fast retrieval |
| **Glacier Flexible** | 99.99% | 90 days | 1-5 min | Archive, flexible retrieval |
| **Deep Archive** | 99.99% | 180 days | 12 hours | Long-term archive |

### Choosing Storage Class

**Decision Tree:**

```
How often accessed?
├─ Daily/Weekly → STANDARD
├─ Monthly → STANDARD_IA
├─ Quarterly → GLACIER_IR (Instant Retrieval)
├─ Annually → GLACIER
└─ Rarely/Compliance → DEEP_ARCHIVE
```

**By Data Type:**

| Data Type | Recommended Class | Reasoning |
|-----------|------------------|-----------|
| Active experiments | STANDARD | Frequent access |
| Published datasets | INTELLIGENT_TIERING | Unknown access |
| Completed projects | STANDARD_IA | Occasional reference |
| Old experiments | GLACIER | Rare access |
| Regulatory archives | DEEP_ARCHIVE | Long-term retention |

---

## Best Practices

### 1. Use Appropriate Storage Tiers

```bash
# Hot data: local fast storage
/data/active/ → NVMe SSD

# Warm data: S3 Standard
cicada sync /data/recent s3://lab-data/standard/

# Cold data: S3 Glacier
aws s3 cp /data/old s3://lab-data/archive/ --storage-class GLACIER --recursive
```

### 2. Enable Versioning for Critical Data

```bash
# Enable versioning
aws s3api put-bucket-versioning \
  --bucket critical-data \
  --versioning-configuration Status=Enabled

# Recover from accidental deletion
aws s3api list-object-versions --bucket critical-data
```

### 3. Implement Lifecycle Policies

```json
{
  "Rules": [
    {
      "Id": "Auto-tier data",
      "Status": "Enabled",
      "Transitions": [
        {"Days": 30, "StorageClass": "STANDARD_IA"},
        {"Days": 90, "StorageClass": "GLACIER"},
        {"Days": 365, "StorageClass": "DEEP_ARCHIVE"}
      ]
    }
  ]
}
```

### 4. Monitor Costs

```bash
# AWS Cost Explorer API
aws ce get-cost-and-usage \
  --time-period Start=2025-01-01,End=2025-01-31 \
  --granularity MONTHLY \
  --metrics BlendedCost \
  --group-by Type=SERVICE
```

### 5. Use Bucket Policies for Access Control

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:user/lab-user"
      },
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::lab-data/*"
    }
  ]
}
```

---

## Troubleshooting

### S3 Access Denied

**Symptoms:**
```
Error: AccessDenied: Access Denied
```

**Solutions:**
1. Check IAM permissions
2. Verify bucket policy
3. Check S3 Block Public Access settings
4. Verify AWS credentials

```bash
# Test AWS CLI access
aws s3 ls s3://lab-data/ --profile default
```

### Slow Transfer Speeds

**Symptoms:**
- Transfers slower than expected
- High latency

**Solutions:**
1. Increase concurrency
2. Check network bandwidth
3. Verify S3 region (use nearest)
4. Consider AWS Direct Connect for large datasets

```bash
# Increase concurrency
cicada sync --concurrency 16 /data s3://bucket/data
```

### Multipart Upload Failures

**Symptoms:**
```
Error: UploadPartFailure: Failed to upload part
```

**Solutions:**
1. Check network stability
2. Verify file permissions
3. Retry upload

```bash
# Abort incomplete multipart uploads
aws s3api list-multipart-uploads --bucket lab-data
aws s3api abort-multipart-upload \
  --bucket lab-data \
  --key file.dat \
  --upload-id <upload-id>
```

---

## Related Documentation

- [Architecture](ARCHITECTURE.md) - Storage architecture details
- [Configuration](CONFIGURATION.md) - Storage configuration
- [CLI Reference](CLI_REFERENCE.md) - Sync commands

---

**Contributing:** Found an error or want to suggest improvements? See [CONTRIBUTING.md](CONTRIBUTING.md)

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
