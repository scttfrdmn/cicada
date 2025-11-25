# Troubleshooting Guide

**Last Updated:** 2025-11-25

Solutions to common problems and errors when using Cicada.

---

## Table of Contents

1. [Installation Issues](#installation-issues)
2. [Configuration Problems](#configuration-problems)
3. [Sync Errors](#sync-errors)
4. [AWS S3 Issues](#aws-s3-issues)
5. [Metadata Extraction Issues](#metadata-extraction-issues)
6. [Watch Daemon Problems](#watch-daemon-problems)
7. [Performance Issues](#performance-issues)
8. [Data Integrity Issues](#data-integrity-issues)
9. [Common Error Messages](#common-error-messages)
10. [Getting Help](#getting-help)

---

## Installation Issues

### Problem: "cicada: command not found"

**Symptoms:**
```bash
$ cicada version
-bash: cicada: command not found
```

**Solutions:**

1. **Verify installation location:**
```bash
# Check if binary exists
ls -l /usr/local/bin/cicada

# If not found, search for it
find ~ -name cicada -type f 2>/dev/null
```

2. **Add to PATH:**
```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH="/usr/local/bin:$PATH"

# Reload shell config
source ~/.bashrc  # or source ~/.zshrc
```

3. **Reinstall:**
```bash
# Download latest version
curl -LO https://github.com/scttfrdmn/cicada/releases/latest/download/cicada-$(uname -s)-$(uname -m)

# Make executable
chmod +x cicada-*

# Move to PATH
sudo mv cicada-* /usr/local/bin/cicada
```

### Problem: "Permission denied" when running cicada

**Symptoms:**
```bash
$ cicada version
-bash: /usr/local/bin/cicada: Permission denied
```

**Solution:**
```bash
# Make executable
chmod +x /usr/local/bin/cicada

# Or reinstall with correct permissions
sudo install -m 755 cicada /usr/local/bin/cicada
```

### Problem: Build fails from source

**Symptoms:**
```bash
$ make build
go: module requires Go 1.21 or later
```

**Solutions:**

1. **Update Go version:**
```bash
# Check current version
go version

# Update Go (macOS)
brew upgrade go

# Update Go (Linux)
sudo rm -rf /usr/local/go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
```

2. **Check build dependencies:**
```bash
# Install required tools
go install golang.org/x/tools/cmd/goimports@latest

# Clean and rebuild
make clean
make build
```

### Problem: macOS quarantine warning

**Symptoms:**
```
"cicada" cannot be opened because it is from an unidentified developer
```

**Solution:**
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/cicada

# Or allow in System Preferences
# System Preferences → Security & Privacy → Allow
```

---

## Configuration Problems

### Problem: Config file not found

**Symptoms:**
```bash
$ cicada config show
Error: configuration file not found
```

**Solutions:**

1. **Create config file:**
```bash
# Create default config
cicada config init

# Verify location
cicada config path
```

2. **Check config locations (searched in order):**
```bash
# 1. Current directory
./cicada.yaml

# 2. User config directory
~/.config/cicada/config.yaml

# 3. System config directory (Linux)
/etc/cicada/config.yaml

# 4. Environment variable
export CICADA_CONFIG=/path/to/config.yaml
```

3. **Specify config explicitly:**
```bash
cicada --config /path/to/config.yaml sync ...
```

### Problem: Invalid configuration

**Symptoms:**
```bash
$ cicada sync
Error: invalid configuration: sync.local_root is required
```

**Solution:**

1. **Validate config:**
```bash
# Check config syntax
cicada config validate

# View current config
cicada config show
```

2. **Fix required fields:**
```bash
# Set missing required fields
cicada config set sync.local_root ~/lab-data
cicada config set sync.remote_root s3://my-bucket/data

# Verify
cicada config validate
```

3. **Reset to defaults:**
```bash
# Backup existing config
cp ~/.config/cicada/config.yaml ~/.config/cicada/config.yaml.bak

# Create fresh config
cicada config init --force
```

### Problem: AWS credentials not found

**Symptoms:**
```bash
Error: NoCredentialProviders: no valid providers in chain
```

**Solutions:**

1. **Configure AWS credentials:**
```bash
# Option 1: AWS CLI
aws configure

# Option 2: Environment variables
export AWS_ACCESS_KEY_ID=your_key_id
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1

# Option 3: Cicada config
cicada config set aws.access_key_id your_key_id
cicada config set aws.secret_access_key your_secret_key
cicada config set aws.region us-east-1
```

2. **Verify credentials:**
```bash
# Test AWS access
aws s3 ls

# Or with Cicada
cicada s3 ls s3://your-bucket
```

3. **Check credential chain:**
```bash
# Cicada checks credentials in this order:
# 1. Environment variables
# 2. Cicada config file
# 3. AWS credentials file (~/.aws/credentials)
# 4. IAM role (EC2/ECS)

# Debug credential resolution
cicada --debug s3 ls s3://your-bucket
```

---

## Sync Errors

### Problem: Sync hangs or appears stuck

**Symptoms:**
- Sync command runs but shows no progress
- Process appears frozen

**Solutions:**

1. **Check network connectivity:**
```bash
# Test S3 connectivity
aws s3 ls s3://your-bucket

# Check network
ping s3.amazonaws.com
```

2. **Use verbose mode:**
```bash
# See what's happening
cicada sync --verbose ~/data s3://bucket/data

# Or debug mode
cicada --debug sync ~/data s3://bucket/data
```

3. **Reduce concurrency:**
```bash
# Lower concurrent transfers
cicada config set sync.concurrency 2

# Or use flag
cicada sync --concurrency 2 ~/data s3://bucket/data
```

4. **Check for large files:**
```bash
# Find large files
find ~/data -type f -size +1G -ls

# Sync large files separately
cicada sync ~/data/large-files s3://bucket/data \
  --concurrency 1 \
  --timeout 3600
```

### Problem: "No space left on device"

**Symptoms:**
```bash
Error: write error: no space left on device
```

**Solutions:**

1. **Check disk space:**
```bash
# Check available space
df -h

# Find large directories
du -sh ~/lab-data/*
```

2. **Clean up temporary files:**
```bash
# Remove Cicada temp files
rm -rf ~/.cache/cicada/tmp/*

# Clean up old metadata cache
cicada cache clean
```

3. **Use streaming mode (no local copy):**
```bash
# Stream directly to S3 without local temp file
cicada sync source s3://bucket/data --streaming
```

4. **Sync in batches:**
```bash
# Sync subdirectories one at a time
for dir in ~/lab-data/*/; do
    cicada sync "$dir" s3://bucket/$(basename "$dir")
    # Clean up after each sync if needed
done
```

### Problem: Sync fails with "Access Denied"

**Symptoms:**
```bash
Error: AccessDenied: Access Denied
Status Code: 403
```

**Solutions:**

1. **Check S3 bucket permissions:**
```bash
# Verify bucket policy
aws s3api get-bucket-policy --bucket your-bucket

# Check IAM permissions
aws iam get-user-policy --user-name your-user --policy-name your-policy
```

2. **Verify IAM permissions needed:**
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
      "Resource": "arn:aws:s3:::your-bucket"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::your-bucket/*"
    }
  ]
}
```

3. **Check bucket ownership:**
```bash
# Verify you own the bucket
aws s3api list-buckets | grep your-bucket
```

4. **Test with different credentials:**
```bash
# Use different AWS profile
export AWS_PROFILE=different-profile
cicada sync ~/data s3://bucket/data
```

### Problem: Sync deletes files unexpectedly

**Symptoms:**
- Files missing after sync
- Warning messages about deletions

**Solutions:**

1. **Always test with --dry-run first:**
```bash
# Preview changes before syncing
cicada sync --dry-run ~/data s3://bucket/data

# Check what would be deleted
cicada sync --dry-run --delete ~/data s3://bucket/data | grep "DELETE"
```

2. **Understand delete modes:**
```bash
# Default: no deletions
cicada sync ~/data s3://bucket/data

# Delete extra files in destination
cicada sync --delete ~/data s3://bucket/data

# Only delete files explicitly removed from source
cicada sync --delete-after ~/data s3://bucket/data
```

3. **Enable version versioning (S3):**
```bash
# Enable S3 versioning for safety
aws s3api put-bucket-versioning \
  --bucket your-bucket \
  --versioning-configuration Status=Enabled

# Recover deleted files
aws s3api list-object-versions \
  --bucket your-bucket \
  --prefix path/to/file
```

4. **Use backup before delete:**
```bash
# Create backup before syncing with delete
cicada backup create ~/data

# Then sync with delete
cicada sync --delete ~/data s3://bucket/data

# Restore if needed
cicada backup restore ~/data
```

---

## AWS S3 Issues

### Problem: S3 upload speed is slow

**Symptoms:**
- Upload taking hours for large files
- Transfer speed < 1 MB/s

**Solutions:**

1. **Increase concurrency:**
```bash
# Use more concurrent uploads
cicada config set sync.concurrency 8

# Or use flag
cicada sync --concurrency 8 ~/data s3://bucket/data
```

2. **Check network bandwidth:**
```bash
# Test S3 transfer speed
aws s3 cp /dev/zero s3://your-bucket/test-file --expected-size 1048576000

# Test network speed to AWS
curl -o /dev/null https://s3.amazonaws.com/speed-test-file
```

3. **Use multipart uploads for large files:**
```bash
# Configure multipart threshold (default: 100 MB)
cicada config set aws.multipart_threshold 50MB
cicada config set aws.multipart_chunk_size 10MB
```

4. **Choose closer AWS region:**
```bash
# Use region closer to you
cicada config set aws.region us-west-2  # If you're on West Coast

# Create bucket in closer region
aws s3 mb s3://your-bucket --region us-west-2
```

5. **Compress files before upload:**
```bash
# Compress large files
tar -czf data.tar.gz ~/data

# Upload compressed file
cicada sync data.tar.gz s3://bucket/data/
```

### Problem: S3 costs higher than expected

**Symptoms:**
- Unexpected AWS bill
- High request costs

**Solutions:**

1. **Check current costs:**
```bash
# Estimate costs
cicada s3 cost-estimate s3://your-bucket

# View storage breakdown
cicada s3 storage-class s3://your-bucket
```

2. **Implement lifecycle policies:**
```bash
# Move old data to cheaper storage
cicada s3 lifecycle create s3://your-bucket \
  --transition-days 30 \
  --transition-class STANDARD_IA

# Archive old data
cicada s3 lifecycle create s3://your-bucket \
  --transition-days 90 \
  --transition-class GLACIER
```

3. **Reduce request costs:**
```bash
# Batch operations instead of many small requests
# Use larger files when possible

# Enable metadata caching
cicada config set cache.enabled true
cicada config set cache.ttl 3600
```

4. **Delete old data:**
```bash
# Find old files
cicada s3 ls s3://your-bucket --older-than 365days

# Delete old data
cicada s3 rm s3://your-bucket/old-data --recursive --older-than 365days

# Or use lifecycle expiration
cicada s3 lifecycle create s3://your-bucket \
  --expiration-days 365 \
  --prefix "temp/"
```

### Problem: S3 bucket not found

**Symptoms:**
```bash
Error: NoSuchBucket: The specified bucket does not exist
```

**Solutions:**

1. **Verify bucket name:**
```bash
# List all your buckets
aws s3 ls

# Check bucket name in config
cicada config show | grep remote_root
```

2. **Create bucket:**
```bash
# Create bucket
aws s3 mb s3://your-bucket

# Or with Cicada
cicada s3 create-bucket your-bucket
```

3. **Check region:**
```bash
# Bucket might be in different region
aws s3api get-bucket-location --bucket your-bucket

# Set correct region
cicada config set aws.region us-west-2
```

### Problem: "RequestTimeTooSkewed" error

**Symptoms:**
```bash
Error: RequestTimeTooSkewed: The difference between the request time and the server's time is too large
```

**Solution:**
```bash
# Sync system clock (macOS)
sudo sntp -sS time.apple.com

# Sync system clock (Linux)
sudo ntpdate -s time.nist.gov

# Enable automatic time sync (Linux)
sudo timedatectl set-ntp true
```

---

## Metadata Extraction Issues

### Problem: Metadata not extracted

**Symptoms:**
- `.cicada/metadata/` directory empty
- "No metadata found" messages

**Solutions:**

1. **Verify auto-extraction is enabled:**
```bash
# Check config
cicada config show | grep auto_extract

# Enable if needed
cicada config set metadata.auto_extract true
```

2. **Manually extract metadata:**
```bash
# Extract from specific file
cicada metadata extract ~/data/file.nd2

# Extract from directory
cicada metadata extract ~/data --recursive
```

3. **Check file format support:**
```bash
# List supported formats
cicada metadata formats

# Check if file type is supported
file ~/data/file.nd2
```

4. **Verify file permissions:**
```bash
# Check file is readable
ls -l ~/data/file.nd2

# Fix permissions if needed
chmod 644 ~/data/file.nd2
```

### Problem: Incorrect or corrupt metadata

**Symptoms:**
- Missing metadata fields
- Incorrect values
- JSON parsing errors

**Solutions:**

1. **Re-extract metadata:**
```bash
# Force re-extraction
cicada metadata extract ~/data/file.nd2 --force

# Delete old metadata and re-extract
rm -rf .cicada/metadata/
cicada metadata extract ~/data --recursive
```

2. **Validate metadata:**
```bash
# Validate against schema
cicada metadata validate ~/data

# Check specific file
cicada metadata show ~/data/file.nd2
```

3. **Check file integrity:**
```bash
# Verify file is not corrupted
cicada verify ~/data/file.nd2

# Check file size
ls -lh ~/data/file.nd2
```

4. **Update extractors:**
```bash
# Update Cicada to latest version
cicada update

# Or rebuild from source
cd ~/cicada
git pull
make install
```

### Problem: Custom metadata not persisting

**Symptoms:**
- Added metadata fields disappear
- Metadata resets after sync

**Solutions:**

1. **Use correct add command:**
```bash
# Add custom metadata (persists)
cicada metadata add ~/data/file.nd2 \
  --field experiment_id=EXP001 \
  --field sample="HeLa cells"

# Verify added
cicada metadata show ~/data/file.nd2 | grep experiment_id
```

2. **Check metadata storage:**
```bash
# Verify metadata file exists
ls -la .cicada/metadata/

# Check metadata file content
cat .cicada/metadata/$(basename ~/data/file.nd2).json
```

3. **Ensure metadata is synced:**
```bash
# Sync metadata directory
cicada sync .cicada/metadata/ s3://bucket/.cicada/metadata/

# Or enable auto-sync
cicada config set metadata.auto_sync true
```

---

## Watch Daemon Problems

### Problem: Watch daemon not starting

**Symptoms:**
```bash
$ cicada watch start
Error: failed to start daemon: address already in use
```

**Solutions:**

1. **Check if daemon is already running:**
```bash
# Check daemon status
cicada watch status

# View daemon logs
cicada watch log --tail 50
```

2. **Kill existing daemon:**
```bash
# Stop daemon
cicada watch stop

# Force kill if needed
pkill -f "cicada watch"

# Restart
cicada watch start
```

3. **Check port availability:**
```bash
# Default daemon port: 8765
lsof -i :8765

# Use different port
cicada config set watch.port 8766
cicada watch start
```

4. **Fix permissions:**
```bash
# Check daemon directory permissions
ls -ld ~/.config/cicada/watch/

# Fix if needed
chmod 755 ~/.config/cicada/watch/
```

### Problem: Watch daemon crashes or stops

**Symptoms:**
- Daemon stops unexpectedly
- No auto-sync happening

**Solutions:**

1. **Check daemon logs:**
```bash
# View recent logs
cicada watch log --tail 100

# Watch logs in real-time
cicada watch log --follow
```

2. **Increase resource limits:**
```bash
# Check current limits
ulimit -a

# Increase file descriptors
ulimit -n 4096

# Make permanent (add to ~/.bashrc)
echo "ulimit -n 4096" >> ~/.bashrc
```

3. **Reduce watched paths:**
```bash
# Edit config to watch fewer directories
cicada config set watch.paths ~/lab-data/raw

# Restart daemon
cicada watch restart
```

4. **Run in foreground for debugging:**
```bash
# Run in foreground to see errors
cicada watch --foreground --verbose
```

### Problem: Changes not detected

**Symptoms:**
- Files added but not synced
- No response to file modifications

**Solutions:**

1. **Check scan interval:**
```bash
# Reduce scan interval
cicada config set watch.scan_interval 30s

# Restart daemon
cicada watch restart
```

2. **Verify watched paths:**
```bash
# Check configured paths
cicada config show | grep watch.paths

# Ensure path is absolute
cicada config set watch.paths /Users/username/lab-data/raw
```

3. **Check debounce delay:**
```bash
# Adjust debounce delay
cicada config set watch.debounce_delay 3s

# Lower = faster response, higher = fewer false triggers
```

4. **Test manually:**
```bash
# Create test file
echo "test" > ~/lab-data/raw/test.txt

# Check if detected
cicada watch log --tail 10 | grep test.txt
```

---

## Performance Issues

### Problem: Slow metadata extraction

**Symptoms:**
- Extraction taking minutes per file
- High CPU usage during extraction

**Solutions:**

1. **Enable parallel extraction:**
```bash
# Use multiple cores
cicada config set metadata.concurrency 4

# Or use flag
cicada metadata extract ~/data --concurrency 4
```

2. **Cache metadata:**
```bash
# Enable metadata caching
cicada config set cache.enabled true
cicada config set cache.metadata true
```

3. **Skip heavy extraction for large files:**
```bash
# Set file size limit
cicada config set metadata.max_file_size 1GB

# Files larger than limit won't be fully extracted
```

4. **Use faster extractor variants:**
```bash
# Use quick extraction mode (less detailed)
cicada metadata extract ~/data --quick
```

### Problem: High memory usage

**Symptoms:**
- Cicada using several GB of RAM
- System slowdown during sync

**Solutions:**

1. **Reduce concurrency:**
```bash
# Lower concurrent operations
cicada config set sync.concurrency 2
```

2. **Enable streaming mode:**
```bash
# Stream files instead of loading into memory
cicada config set sync.streaming true
```

3. **Clear metadata cache:**
```bash
# Clear cache
cicada cache clear

# Reduce cache size
cicada config set cache.max_size 100MB
```

4. **Process in batches:**
```bash
# Sync subdirectories separately
for dir in ~/data/*/; do
    cicada sync "$dir" s3://bucket/$(basename "$dir")
done
```

### Problem: Slow initial sync

**Symptoms:**
- First sync taking very long
- Scanning phase taking hours

**Solutions:**

1. **Use --fast flag:**
```bash
# Skip detailed scanning
cicada sync --fast ~/data s3://bucket/data
```

2. **Exclude unnecessary files:**
```bash
# Create .cicadaignore file
cat > ~/data/.cicadaignore << 'EOF'
*.tmp
*.swp
.DS_Store
thumbs.db
EOF
```

3. **Sync in chunks:**
```bash
# Sync by subdirectory
cicada sync ~/data/2025-11 s3://bucket/data/2025-11
cicada sync ~/data/2025-10 s3://bucket/data/2025-10
```

4. **Pre-build file list:**
```bash
# Create file list first
find ~/data -type f > files.txt

# Sync using file list
cicada sync --files-from files.txt ~/data s3://bucket/data
```

---

## Data Integrity Issues

### Problem: Checksum mismatches

**Symptoms:**
```bash
Error: checksum mismatch: expected abc123, got def456
```

**Solutions:**

1. **Re-sync with checksum verification:**
```bash
# Force checksum verification
cicada sync --checksum ~/data s3://bucket/data
```

2. **Verify file integrity:**
```bash
# Check local file
cicada verify ~/data/file.dat

# Check S3 file
cicada s3 verify s3://bucket/data/file.dat
```

3. **Re-upload corrupted files:**
```bash
# Find files with checksum errors
cicada sync --dry-run --checksum ~/data s3://bucket/data | grep "checksum mismatch"

# Re-upload specific files
cicada sync --force ~/data/corrupted-file.dat s3://bucket/data/
```

4. **Enable ETag verification:**
```bash
# Use S3 ETags for integrity checking
cicada config set sync.verify_etag true
```

### Problem: Files corrupted after transfer

**Symptoms:**
- Files can't be opened after download
- Different file size than original

**Solutions:**

1. **Always verify after transfer:**
```bash
# Sync with verification
cicada sync ~/data s3://bucket/data --verify

# Verify existing files
cicada verify ~/data --recursive
```

2. **Check for network issues:**
```bash
# Test network stability
ping -c 100 s3.amazonaws.com

# If unstable, use retry
cicada sync ~/data s3://bucket/data --retry 3
```

3. **Use S3 transfer acceleration:**
```bash
# Enable transfer acceleration
aws s3api put-bucket-accelerate-configuration \
  --bucket your-bucket \
  --accelerate-configuration Status=Enabled

# Use in Cicada
cicada config set aws.use_accelerate true
```

---

## Common Error Messages

### "context deadline exceeded"

**Cause:** Operation timed out

**Solutions:**
```bash
# Increase timeout
cicada sync --timeout 3600 ~/data s3://bucket/data

# Or in config
cicada config set sync.timeout 3600
```

### "too many open files"

**Cause:** File descriptor limit reached

**Solution:**
```bash
# Increase limit
ulimit -n 4096

# Make permanent
echo "ulimit -n 4096" >> ~/.bashrc

# Reduce concurrency
cicada config set sync.concurrency 2
```

### "connection reset by peer"

**Cause:** Network connection lost

**Solutions:**
```bash
# Enable retries
cicada sync --retry 3 ~/data s3://bucket/data

# Check network stability
ping -c 100 s3.amazonaws.com

# Use longer timeout
cicada sync --timeout 600 ~/data s3://bucket/data
```

### "invalid character in metadata"

**Cause:** Special characters in metadata fields

**Solution:**
```bash
# Use proper escaping
cicada metadata add file.txt --field 'note=Test "quoted" text'

# Or use file input
echo 'note: Test "quoted" text' > metadata.yaml
cicada metadata add file.txt --from-file metadata.yaml
```

### "bucket is in a different region"

**Cause:** Bucket region mismatch

**Solution:**
```bash
# Check bucket region
aws s3api get-bucket-location --bucket your-bucket

# Update config
cicada config set aws.region us-west-2
```

---

## Getting Help

### Enable Debug Mode

```bash
# Run with debug output
cicada --debug sync ~/data s3://bucket/data

# Save debug output to file
cicada --debug sync ~/data s3://bucket/data 2> debug.log
```

### Check Version and System Info

```bash
# Check Cicada version
cicada version

# System information
cicada system-info
```

### Generate Diagnostic Report

```bash
# Create diagnostic report
cicada diagnose > cicada-diagnostic.txt

# This includes:
# - Version info
# - Configuration
# - Recent errors
# - System information
```

### Search Documentation

```bash
# Search docs for topic
cicada help search "metadata"

# View specific topic
cicada help metadata extract
```

### Community Support

1. **GitHub Issues**: https://github.com/scttfrdmn/cicada/issues
   - Search existing issues
   - Create new issue with diagnostic report

2. **GitHub Discussions**: https://github.com/scttfrdmn/cicada/discussions
   - Ask questions
   - Share solutions

3. **Documentation**: https://scttfrdmn.github.io/cicada
   - Complete guides
   - API reference
   - Examples

### Reporting Bugs

When reporting bugs, include:

1. **Cicada version:**
```bash
cicada version
```

2. **Command that failed:**
```bash
cicada --debug <your-command> 2> error.log
```

3. **Configuration (redact credentials):**
```bash
cicada config show | sed 's/secret.*/***REDACTED***/g'
```

4. **System information:**
```bash
cicada system-info
```

5. **Error logs:**
```bash
cicada watch log --tail 100
```

---

## Quick Troubleshooting Checklist

```bash
# 1. Check version
cicada version

# 2. Validate configuration
cicada config validate

# 3. Test AWS credentials
aws s3 ls

# 4. Test network
ping s3.amazonaws.com

# 5. Check disk space
df -h

# 6. View recent logs
cicada watch log --tail 50

# 7. Run diagnostic
cicada diagnose

# 8. Try with debug mode
cicada --debug <your-command>
```

---

**Related Documentation:**
- [Getting Started](GETTING_STARTED.md)
- [Common Workflows](WORKFLOWS.md)
- [Advanced Topics](ADVANCED.md)
- [Configuration Reference](CONFIGURATION.md)
