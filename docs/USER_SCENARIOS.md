# Cicada User Scenarios

Cicada is a **dormant data commons platform for academic research labs**, providing federated storage, access control, compute-to-data capabilities, and collaboration primitives. Like a cicada, it lies dormant (consuming minimal resources) until needed, then emerges powerfully for data-intensive work.

This document provides persona-based walkthroughs for **v0.1.0**, which implements the foundational storage and sync layer. These scenarios demonstrate how researchers can leverage Cicada's core capabilities to automate instrument data uploads, backup research data, integrate with analysis pipelines, and ensure data governance compliance.

## Table of Contents

1. [Lab Researcher: Instrument Data Upload](#scenario-1-lab-researcher-instrument-data-upload)
2. [PhD Student: Research Data Backup](#scenario-2-phd-student-research-data-backup)
3. [Bioinformatician: Analysis Pipeline](#scenario-3-bioinformatician-analysis-pipeline)
4. [Lab Manager: Data Governance](#scenario-4-lab-manager-data-governance)

---

## Scenario 1: Lab Researcher - Instrument Data Upload

### Persona: Dr. Maria Rodriguez

**Background**:
- Postdoctoral researcher in neuroscience
- Uses a Zeiss confocal microscope that writes ~50GB/day of imaging data
- Lab policy: All raw data must be in S3 within 24 hours
- Technical level: Comfortable with command line basics
- Works at the microscope workstation (Ubuntu Linux)

**Pain Points**:
- Manually copying files to S3 is tedious and error-prone
- Forgot to upload data twice → got reminders from PI
- Large file transfers sometimes fail halfway through
- Needs to verify all files uploaded correctly

**Goals**:
- Automatic upload of new microscope files
- Reliable transfers that handle network interruptions
- Confirmation that files are safely in S3

---

### Version Info
- ✅ **v0.1.0 (Current)**: All features shown are available today

---

### Day 1: Initial Setup (5 minutes)

**Step 1: Install Cicada**

Maria installs Cicada on the microscope workstation:

```bash
# Download from GitHub releases
wget https://github.com/scttfrdmn/cicada/releases/download/v0.1.0/cicada_0.1.0_Linux_x86_64.tar.gz

# Extract and install
tar -xzf cicada_0.1.0_Linux_x86_64.tar.gz
sudo mv cicada /usr/local/bin/

# Verify installation
cicada version
```

**Output**:
```
cicada version 0.1.0 (commit: f6a1f79, built: 2025-11-23T21:36:10Z, by: goreleaser)
  Go version: go1.23.4
  OS/Arch: linux/amd64
```

**What Maria thinks**: *"Good, it's installed. Now I need to set up AWS access."*

---

**Step 2: Configure AWS Credentials**

The lab already has an AWS profile configured on the workstation:

```bash
# Verify AWS credentials work
aws s3 ls s3://rodriguez-lab-data/
```

**Output**:
```
                           PRE microscopy/
                           PRE sequencing/
                           PRE analysis/
```

**What Maria thinks**: *"Perfect, AWS access is working. Now to configure Cicada."*

---

**Step 3: Initialize Cicada Configuration**

```bash
# Create default config
cicada config init
```

**Output**:
```
✓ Created configuration directory: /home/maria/.cicada
✓ Created configuration file: /home/maria/.cicada/config.yaml

Default configuration created. Customize with:
  cicada config set <key> <value>
```

**View the created config**:
```bash
cat ~/.cicada/config.yaml
```

**Generated config**:
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

**What Maria thinks**: *"This looks reasonable. I should update the AWS profile to match our lab's setup."*

---

**Step 4: Configure for Lab Environment**

```bash
# Set the lab's AWS profile
cicada config set aws.profile rodriguez-lab

# Set the correct region
cicada config set aws.region us-west-2

# Verify configuration
cicada config list
```

**Output**:
```
version: "1"

aws:
  profile: rodriguez-lab
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

**What Maria thinks**: *"Great! Now I'm ready to test uploading some files."*

---

### Day 1: Testing Manual Sync (2 minutes)

**Step 5: Test Manual Upload**

Maria has a completed imaging session from yesterday that needs uploading:

```bash
# Check what files are in the directory
ls -lh /mnt/zeiss/2025-11-22_neurons/

# Dry run first to see what would be uploaded
cicada sync --dry-run \
  /mnt/zeiss/2025-11-22_neurons \
  s3://rodriguez-lab-data/microscopy/2025-11-22_neurons
```

**Output**:
```
Syncing from /mnt/zeiss/2025-11-22_neurons to s3://rodriguez-lab-data/microscopy/2025-11-22_neurons

Scanning source files...
Found 127 files (48.3 GB)

Scanning destination...
Found 0 files

Plan:
  Upload: 127 files (48.3 GB)
  Delete: 0 files
  Skip: 0 files (already synced)

DRY RUN - no changes made
```

**What Maria thinks**: *"Perfect! It found all 127 files. Now let me do the actual upload."*

---

**Step 6: Perform Actual Upload**

```bash
# Real upload (remove --dry-run)
cicada sync \
  /mnt/zeiss/2025-11-22_neurons \
  s3://rodriguez-lab-data/microscopy/2025-11-22_neurons
```

**Output** (abbreviated):
```
Syncing from /mnt/zeiss/2025-11-22_neurons to s3://rodriguez-lab-data/microscopy/2025-11-22_neurons

Scanning source files...
Found 127 files (48.3 GB)

Scanning destination...
Found 0 files

Uploading files... (4 concurrent transfers)

✓ image_001.czi (387 MB) - 12s
✓ image_002.czi (391 MB) - 11s
✓ image_003.czi (385 MB) - 12s
✓ image_004.czi (389 MB) - 12s
...
✓ metadata.xml (42 KB) - 0s

Completed: 127 files uploaded (48.3 GB) in 18m 32s
```

**What Maria thinks**: *"Excellent! 48 GB uploaded in under 20 minutes. Much faster than my manual copies. Now let me set up automatic watching."*

---

### Day 1: Setting Up Automatic Watching (3 minutes)

**Step 7: Configure Watch for Microscope Output**

The microscope writes files to `/mnt/zeiss/output/` as imaging sessions complete.

```bash
# Add a watch with appropriate delays
cicada watch add \
  --debounce 30 \
  --min-age 60 \
  /mnt/zeiss/output \
  s3://rodriguez-lab-data/microscopy/live
```

**Parameters explained**:
- `--debounce 30`: Wait 30 seconds after last file change before syncing (handles burst writes)
- `--min-age 60`: Only sync files older than 60 seconds (prevents syncing incomplete files)

**Output**:
```
✓ Watch created: /mnt/zeiss/output-1732397890

Configuration:
  Source: /mnt/zeiss/output
  Destination: s3://rodriguez-lab-data/microscopy/live
  Debounce: 30s
  Min-age: 60s
  Delete source: false

Performing initial sync...
Found 3 files (1.2 GB)
✓ Synced 3 files in 2m 15s

Watch is now active. Files will be automatically synced to S3.
Press Ctrl+C to stop, or let it run in the background.
```

**What Maria thinks**: *"Perfect! Now any new files will automatically upload. Let me verify it's working."*

---

**Step 8: Verify Watch is Active**

```bash
# List all active watches
cicada watch list
```

**Output**:
```
Active watches:

ID: /mnt/zeiss/output-1732397890
  Source: /mnt/zeiss/output
  Destination: s3://rodriguez-lab-data/microscopy/live
  Status: Running
  Debounce: 30s
  Min-age: 60s
  Last sync: 2 minutes ago (3 files, 1.2 GB)
  Total syncs: 1
```

**What Maria thinks**: *"Great! It's running and already did the initial sync. Now let me test it with a new file."*

---

### Day 1: Testing Automatic Sync (5 minutes)

**Step 9: Test with New File**

Maria copies a test file to simulate the microscope writing data:

```bash
# Copy a test file to the watched directory
cp /home/maria/test-image.czi /mnt/zeiss/output/

# Wait 90 seconds (min-age 60s + debounce 30s)
# ... Maria gets a coffee ...

# Check S3 to see if file appeared
aws s3 ls s3://rodriguez-lab-data/microscopy/live/
```

**Output**:
```
2025-11-23 14:23:15  387234561 image_001.czi
2025-11-23 14:23:18  391847234 image_002.czi
2025-11-23 14:23:21  385123456 image_003.czi
2025-11-23 14:28:42  389456789 test-image.czi  ← New file!
```

**What Maria thinks**: *"Success! The file was automatically uploaded. This is going to save me so much time."*

---

**Step 10: Check Watch Status Again**

```bash
cicada watch list
```

**Output**:
```
Active watches:

ID: /mnt/zeiss/output-1732397890
  Source: /mnt/zeiss/output
  Destination: s3://rodriguez-lab-data/microscopy/live
  Status: Running
  Debounce: 30s
  Min-age: 60s
  Last sync: 1 minute ago (1 file, 389 MB)
  Total syncs: 2
```

**What Maria thinks**: *"Perfect! It tracked the sync. Now I should make sure this starts automatically when the computer reboots."*

---

### Day 1: Making Watch Persistent (2 minutes)

**Step 11: Verify Configuration Persistence**

Cicada automatically saves watch configuration:

```bash
# Check the saved configuration
cat ~/.cicada/config.yaml
```

**Updated config (watches section)**:
```yaml
version: "1"

aws:
  profile: rodriguez-lab
  region: us-west-2

sync:
  concurrency: 4
  delete: false
  exclude:
    - .git/**
    - .DS_Store
    - "*.tmp"
    - "*.swp"

watches:
  - id: /mnt/zeiss/output-1732397890
    source: /mnt/zeiss/output
    destination: s3://rodriguez-lab-data/microscopy/live
    debounce: 30
    min_age: 60
    delete_source: false
    sync_on_start: true
    enabled: true

settings:
  verbose: false
  check_updates: true
```

**What Maria thinks**: *"Good, it's saved in the config. Now I need to set up a systemd service so it starts on boot."*

---

**Step 12: Create Systemd Service**

Maria creates a systemd service to run Cicada watches automatically:

```bash
# Create service file
sudo nano /etc/systemd/system/cicada-watch.service
```

**Service file content**:
```ini
[Unit]
Description=Cicada File Watcher
After=network.target

[Service]
Type=simple
User=maria
WorkingDirectory=/home/maria
ExecStart=/usr/local/bin/cicada watch start
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Enable and start the service**:
```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable cicada-watch

# Start the service now
sudo systemctl start cicada-watch

# Check status
sudo systemctl status cicada-watch
```

**Output**:
```
● cicada-watch.service - Cicada File Watcher
   Loaded: loaded (/etc/systemd/system/cicada-watch.service; enabled)
   Active: active (running) since Sat 2025-11-23 14:35:12 PST; 5s ago
 Main PID: 12345 (cicada)
   CGroup: /system.slice/cicada-watch.service
           └─12345 /usr/local/bin/cicada watch start

Nov 23 14:35:12 zeiss-workstation systemd[1]: Started Cicada File Watcher.
Nov 23 14:35:13 zeiss-workstation cicada[12345]: ✓ Watch loaded: /mnt/zeiss/output
Nov 23 14:35:13 zeiss-workstation cicada[12345]: Watching for file changes...
```

**What Maria thinks**: *"Perfect! Now it will automatically start watching whenever the computer boots. Setup complete!"*

---

### Week 1: Daily Usage

**Maria's typical workflow now**:

1. **Morning**: Arrives at lab, checks that Cicada is running
   ```bash
   sudo systemctl status cicada-watch
   cicada watch list
   ```

2. **During imaging**: Works normally at the microscope
   - Files automatically upload as sessions complete
   - No manual intervention needed

3. **End of day**: Verifies all data is in S3
   ```bash
   # Check S3 to confirm files are there
   aws s3 ls s3://rodriguez-lab-data/microscopy/live/ --recursive
   ```

4. **Weekly**: Reviews sync statistics
   ```bash
   cicada watch list
   ```
   **Example output**:
   ```
   ID: /mnt/zeiss/output-1732397890
     Status: Running
     Last sync: 15 minutes ago (8 files, 3.2 GB)
     Total syncs: 47
     Total data: 241 GB
   ```

**What Maria experiences**: *"I barely think about data uploads anymore. I just check once a day that Cicada is running, and all my microscope data is automatically backed up to S3. This has saved me hours of work every week."*

---

### Key Benefits for Maria

✅ **Time savings**: No manual uploads (saves ~2 hours/week)
✅ **Reliability**: Never forgets to upload data
✅ **Peace of mind**: Automatic verification via watch status
✅ **Fast transfers**: Concurrent uploads maximize bandwidth
✅ **Simple setup**: 15-minute one-time configuration

---

## Scenario 2: PhD Student - Research Data Backup

### Persona: Alex Thompson

**Background**:
- 3rd year PhD student in bioinformatics
- Works on RNA-seq analysis with ~500GB of data
- University provides S3 storage for research data
- Technical level: Proficient programmer, uses command line daily
- Works from laptop (macOS) and university HPC cluster

**Pain Points**:
- Lost a week of analysis when laptop hard drive failed
- Manually copying files to S3 is tedious
- Not sure if backups are complete or up-to-date
- Running low on laptop disk space

**Goals**:
- Reliable automated backups to S3
- Verify all important data is backed up
- Free up local disk space by moving old data to S3
- Easy access to data from both laptop and HPC

---

### Day 1: Installing on macOS (2 minutes)

**Step 1: Install via Homebrew**

```bash
# Add the Cicada tap
brew install scttfrdmn/tap/cicada

# Verify installation
cicada version
```

**Output**:
```
cicada version 0.1.0 (commit: f6a1f79, built: 2025-11-23T21:36:10Z, by: goreleaser)
  Go version: go1.23.4
  OS/Arch: darwin/arm64
```

**What Alex thinks**: *"Nice, Homebrew makes this easy. Now to set up AWS."*

---

**Step 2: Configure AWS Credentials**

```bash
# Alex already has AWS CLI configured for university account
aws configure list

# Test S3 access
aws s3 ls s3://alex-phd-research/
```

**Output**:
```
                           PRE rnaseq-analysis/
                           PRE reference-genomes/
                           PRE publications/
```

**What Alex thinks**: *"Good, S3 access works. Let me initialize Cicada."*

---

**Step 3: Initialize Configuration**

```bash
# Create config
cicada config init

# Set university AWS profile
cicada config set aws.profile university

# Set region
cicada config set aws.region us-east-1

# Verify
cicada config get aws.profile
cicada config get aws.region
```

**Output**:
```
university
us-east-1
```

**What Alex thinks**: *"Configuration looks good. Now let me back up my current analysis."*

---

### Day 1: First Backup - Current Analysis (10 minutes)

**Step 4: Preview Backup**

Alex wants to back up the current RNA-seq analysis project:

```bash
# See what will be uploaded (dry run)
cicada sync --dry-run \
  ~/research/rnaseq-analysis \
  s3://alex-phd-research/rnaseq-analysis
```

**Output**:
```
Syncing from /Users/alex/research/rnaseq-analysis to s3://alex-phd-research/rnaseq-analysis

Scanning source files...
Found 1,847 files (67.3 GB)

Scanning destination...
Found 1,203 files (45.2 GB)

Plan:
  Upload: 644 new files (22.1 GB)
  Update: 0 modified files
  Delete: 0 files
  Skip: 1,203 files (already synced)

DRY RUN - no changes made
```

**What Alex thinks**: *"Good! It recognizes that some files are already there and only needs to upload 644 new ones. Let me do the real sync."*

---

**Step 5: Perform Backup**

```bash
# Real backup
cicada sync \
  ~/research/rnaseq-analysis \
  s3://alex-phd-research/rnaseq-analysis
```

**Output**:
```
Syncing from /Users/alex/research/rnaseq-analysis to s3://alex-phd-research/rnaseq-analysis

Scanning source files...
Found 1,847 files (67.3 GB)

Scanning destination...
Found 1,203 files (45.2 GB)

Uploading files... (4 concurrent transfers)

✓ results/DE_analysis.csv (2.3 MB) - 1s
✓ results/counts_matrix.tsv (45 MB) - 3s
✓ plots/volcano_plot.png (387 KB) - 0s
✓ data/sample_A_01.fastq.gz (891 MB) - 24s
...

Completed: 644 files uploaded (22.1 GB) in 8m 15s
  Skipped: 1,203 files (already in sync)
```

**What Alex thinks**: *"Excellent! Only 8 minutes to upload 22 GB of new data. And it intelligently skipped the files that were already there."*

---

### Day 1: Setting Up Continuous Backup (5 minutes)

**Step 6: Configure Watch for Active Project**

```bash
# Watch the analysis directory
cicada watch add \
  --debounce 60 \
  --min-age 120 \
  ~/research/rnaseq-analysis \
  s3://alex-phd-research/rnaseq-analysis
```

**Parameters**:
- `--debounce 60`: Wait 1 minute after changes (allows for batch edits)
- `--min-age 120`: Only backup files older than 2 minutes (prevents backing up temp files)

**Output**:
```
✓ Watch created: /Users/alex/research/rnaseq-analysis-1732398234

Configuration:
  Source: /Users/alex/research/rnaseq-analysis
  Destination: s3://alex-phd-research/rnaseq-analysis
  Debounce: 60s
  Min-age: 120s

Performing initial sync...
All files already synced (1,847 files in sync)

Watch is now active. Changes will be automatically backed up to S3.
```

**What Alex thinks**: *"Perfect! Now any changes I make will automatically get backed up. But wait, I need this to run even when I close my terminal..."*

---

**Step 7: Run Watch in Background**

Alex uses macOS LaunchAgent to run Cicada in the background:

```bash
# Create LaunchAgent directory if needed
mkdir -p ~/Library/LaunchAgents

# Create plist file
cat > ~/Library/LaunchAgents/com.cicada.watch.plist <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.cicada.watch</string>
    <key>ProgramArguments</key>
    <array>
        <string>/opt/homebrew/bin/cicada</string>
        <string>watch</string>
        <string>start</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/alex/Library/Logs/cicada-watch.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/alex/Library/Logs/cicada-watch-error.log</string>
</dict>
</plist>
EOF

# Load the LaunchAgent
launchctl load ~/Library/LaunchAgents/com.cicada.watch.plist

# Verify it's running
launchctl list | grep cicada
```

**Output**:
```
12346	0	com.cicada.watch
```

**What Alex thinks**: *"Great! Now Cicada will run in the background and automatically backup my work."*

---

### Week 1: Recovering from Laptop Issue

**Day 5: Laptop Problem**

Alex's laptop runs out of disk space:

```bash
df -h /Users/alex
```

**Output**:
```
Filesystem      Size   Used  Avail Capacity
/dev/disk1s1   500GB  485GB   15GB    97%
```

**What Alex thinks**: *"I'm almost out of space! Let me move some old data to S3 and delete the local copies."*

---

**Step 8: Archive Old Data**

```bash
# Sync old reference genomes to S3
cicada sync \
  ~/research/reference-genomes \
  s3://alex-phd-research/reference-genomes

# Verify all files are in S3
aws s3 ls s3://alex-phd-research/reference-genomes/ --recursive --summarize
```

**Output**:
```
...
2025-11-28 10:23:15  1234567890 hg38.fa
2025-11-28 10:23:18  2345678901 hg38.fa.fai
...

Total Objects: 47
   Total Size: 45.8 GB
```

**Verify checksums match**:
```bash
# Use sync with --dry-run to verify everything matches
cicada sync --dry-run \
  ~/research/reference-genomes \
  s3://alex-phd-research/reference-genomes
```

**Output**:
```
Syncing from /Users/alex/research/reference-genomes to s3://alex-phd-research/reference-genomes

Scanning source files...
Found 47 files (45.8 GB)

Scanning destination...
Found 47 files (45.8 GB)

Plan:
  Upload: 0 files
  Delete: 0 files
  Skip: 47 files (already synced) ← All files match!

DRY RUN - no changes made
```

**What Alex thinks**: *"Perfect! All 47 files are verified to be in S3 with matching checksums. Safe to delete locally."*

---

**Step 9: Free Up Local Disk Space**

```bash
# Delete local copies (keeping only S3)
rm -rf ~/research/reference-genomes

# Check disk space
df -h /Users/alex
```

**Output**:
```
Filesystem      Size   Used  Avail Capacity
/dev/disk1s1   500GB  439GB   61GB    88%
```

**What Alex thinks**: *"Freed up 46 GB! And I can always get the files back from S3 when needed."*

---

### Week 2: Working on HPC Cluster

**Step 10: Access Data from HPC**

Alex needs to run analysis on the university HPC cluster:

```bash
# SSH to HPC cluster
ssh alex@hpc.university.edu

# Install Cicada on HPC (from downloaded binary)
wget https://github.com/scttfrdmn/cicada/releases/download/v0.1.0/cicada_0.1.0_Linux_x86_64.tar.gz
tar -xzf cicada_0.1.0_Linux_x86_64.tar.gz
mv cicada ~/bin/

# Initialize config
cicada config init
cicada config set aws.profile university
cicada config set aws.region us-east-1

# Download data from S3
cicada sync \
  s3://alex-phd-research/rnaseq-analysis \
  ~/scratch/rnaseq-analysis
```

**Output**:
```
Syncing from s3://alex-phd-research/rnaseq-analysis to /home/alex/scratch/rnaseq-analysis

Scanning source...
Found 1,847 files (67.3 GB)

Scanning destination...
Found 0 files

Downloading files... (4 concurrent transfers)

✓ data/sample_A_01.fastq.gz (891 MB) - 18s
✓ data/sample_A_02.fastq.gz (867 MB) - 17s
...

Completed: 1,847 files downloaded (67.3 GB) in 12m 42s
```

**What Alex thinks**: *"Perfect! All my data is now on the HPC cluster. I can run my analysis and then sync the results back to S3."*

---

**Step 11: Sync Results Back to S3**

After running analysis on HPC:

```bash
# Upload new results back to S3
cicada sync \
  ~/scratch/rnaseq-analysis/results \
  s3://alex-phd-research/rnaseq-analysis/results
```

**Output**:
```
Syncing from /home/alex/scratch/rnaseq-analysis/results to s3://alex-phd-research/rnaseq-analysis/results

Scanning source files...
Found 234 files (8.7 GB)

Scanning destination...
Found 189 files (6.2 GB)

Uploading files... (4 concurrent transfers)

✓ DE_results_updated.csv (3.4 MB) - 1s
✓ pathway_analysis.txt (892 KB) - 0s
...

Completed: 45 new files uploaded (2.5 GB) in 2m 18s
  Skipped: 189 files (already in sync)
```

**What Alex thinks**: *"Excellent! My new results are now in S3 and will sync to my laptop automatically via the watch I set up."*

---

### Month 1: Key Benefits for Alex

✅ **Data safety**: No more data loss from hardware failures
✅ **Disk space management**: Can archive old data to S3, free up local space
✅ **Multi-device workflow**: Access data from laptop and HPC seamlessly
✅ **Automatic versioning**: S3 versioning enabled, can recover old file versions
✅ **Peace of mind**: Continuous automatic backup of active work

**What Alex experiences**: *"Cicada has completely changed how I manage my research data. I never worry about losing work anymore, and I can easily work across my laptop and the HPC cluster. The automatic backup means I focus on research, not file management."*

---

## Scenario 3: Bioinformatician - Analysis Pipeline

### Persona: Dr. James Park

**Background**:
- Bioinformatics core facility manager
- Runs high-throughput analysis pipelines (WGS, RNA-seq, ChIP-seq)
- Processes ~50 TB/month of sequencing data
- Technical level: Expert Linux admin and pipeline developer
- Infrastructure: 10-node analysis cluster + S3 for long-term storage

**Pain Points**:
- Need to move completed analysis results to S3 for archival
- Want to automatically clean up local storage after S3 upload
- Must verify data integrity before deleting local copies
- Current scripts are fragile and require monitoring

**Goals**:
- Reliable pipeline integration for S3 archival
- Automatic cleanup of local storage after successful upload
- Strong verification that uploads completed successfully
- Integration with Nextflow/Snakemake pipelines

---

### Implementation: Pipeline Integration

**Step 1: Install Cicada on Analysis Nodes**

```bash
# Install on all analysis nodes
for node in node{01..10}; do
  ssh $node "
    wget -q https://github.com/scttfrdmn/cicada/releases/download/v0.1.0/cicada_0.1.0_Linux_x86_64.tar.gz
    tar -xzf cicada_0.1.0_Linux_x86_64.tar.gz
    sudo mv cicada /usr/local/bin/
    cicada version
  "
done
```

**What James thinks**: *"Quick installation across all nodes. Now to integrate with our pipeline."*

---

**Step 2: Create Pipeline Archive Script**

James creates a script for Nextflow integration:

```bash
#!/bin/bash
# archive_to_s3.sh - Archive completed analysis to S3

set -euo pipefail

PROJECT_ID="$1"
LOCAL_DIR="/data/analysis/${PROJECT_ID}"
S3_PREFIX="s3://bioinformatics-archive/projects/${PROJECT_ID}"

echo "=== Archiving project ${PROJECT_ID} to S3 ==="

# Step 1: Dry run to validate
echo "Step 1: Validating files..."
cicada sync --dry-run "${LOCAL_DIR}" "${S3_PREFIX}"

# Step 2: Perform upload
echo "Step 2: Uploading to S3..."
cicada sync --concurrency 16 "${LOCAL_DIR}" "${S3_PREFIX}"

# Step 3: Verify upload with dry-run (should show everything in sync)
echo "Step 3: Verifying upload..."
cicada sync --dry-run "${LOCAL_DIR}" "${S3_PREFIX}" | grep "Skip: .* files (already synced)"

if [ $? -eq 0 ]; then
  echo "✓ Upload verified - all files match S3"

  # Step 4: Clean up local storage (optional)
  if [ "${DELETE_LOCAL:-false}" = "true" ]; then
    echo "Step 4: Cleaning up local files..."
    rm -rf "${LOCAL_DIR}"
    echo "✓ Local files deleted"
  fi
else
  echo "✗ Upload verification failed - NOT deleting local files"
  exit 1
fi

echo "=== Archive complete for ${PROJECT_ID} ==="
```

**What James thinks**: *"This gives me the safety I need - verify before delete, and fail loudly if anything goes wrong."*

---

**Step 3: Integrate with Nextflow Pipeline**

```groovy
// nextflow.config
process {
  withName: archive_results {
    executor = 'local'
    memory = '2 GB'
    cpus = 1
  }
}
```

```groovy
// main.nf
process archive_results {
  publishDir "${params.outdir}", mode: 'copy'

  input:
  path analysis_dir

  output:
  path "archive.log"

  script:
  """
  /usr/local/bin/archive_to_s3.sh ${params.project_id} > archive.log 2>&1
  """
}

workflow {
  // ... analysis steps ...

  // Archive to S3 after analysis completes
  archive_results(analysis_output)
}
```

**What James thinks**: *"Perfect integration with our existing Nextflow pipelines."*

---

**Step 4: Test Pipeline Archival**

```bash
# Run test project
nextflow run analysis_pipeline.nf \
  --project_id TEST_001 \
  --input samples.csv \
  --outdir /data/analysis/TEST_001

# Check archive log
cat /data/analysis/TEST_001/archive.log
```

**Output**:
```
=== Archiving project TEST_001 to S3 ===
Step 1: Validating files...
Plan:
  Upload: 1,247 files (124 GB)
  Skip: 0 files

Step 2: Uploading to S3...
Completed: 1,247 files uploaded (124 GB) in 15m 32s

Step 3: Verifying upload...
Skip: 1,247 files (already synced)
✓ Upload verified - all files match S3

=== Archive complete for TEST_001 ===
```

**What James thinks**: *"Excellent! Safe, verified archival. Now let me set up automatic cleanup for completed projects."*

---

### Key Benefits for James

✅ **Pipeline integration**: Easy to call from Nextflow/Snakemake
✅ **Verification**: Dry-run validates uploads before cleanup
✅ **Performance**: Configurable concurrency (using 16 for large transfers)
✅ **Reliability**: Script fails safely if verification doesn't pass
✅ **Automation**: No manual intervention needed

---

## Scenario 4: Lab Manager - Data Governance

### Persona: Dr. Lisa Zhang

**Background**:
- Lab manager for 25-person research lab
- Responsible for data compliance and backup policies
- Lab generates ~10 TB/month across 5 instruments
- Technical level: Comfortable with command line, focuses on policy
- Must ensure all data is backed up within 24 hours per NIH guidelines

**Pain Points**:
- No centralized visibility into what's backed up
- Researchers sometimes forget to upload data
- Hard to audit compliance with backup policies
- Manual verification is time-consuming

**Goals**:
- Automated backup for all lab instruments
- Audit trail showing all data is backed up
- Compliance with institutional data policies
- Easy monitoring and reporting

---

### Implementation: Lab-Wide Deployment

**Step 1: Instrument-Specific Configurations**

Lisa creates standard configurations for each instrument type:

```bash
# Microscope workstation config
cat > /etc/cicada/microscope-watch.yaml <<EOF
# Cicada watch configuration for Zeiss microscope
# Auto-generated by lab manager

watches:
  - source: /mnt/zeiss/output
    destination: s3://lab-data-backup/microscopy/zeiss
    debounce: 30
    min_age: 60
    enabled: true
EOF

# Sequencer workstation config
cat > /etc/cicada/sequencer-watch.yaml <<EOF
# Cicada watch configuration for Illumina sequencer
# Auto-generated by lab manager

watches:
  - source: /data/illumina/output
    destination: s3://lab-data-backup/sequencing/illumina
    debounce: 60
    min_age: 300
    enabled: true
EOF
```

**What Lisa thinks**: *"Standardized configurations make it easy to deploy across all instruments."*

---

**Step 2: Monitoring Script**

Lisa creates a monitoring dashboard:

```bash
#!/bin/bash
# lab_backup_status.sh - Check backup status across all instruments

echo "=== Lab Backup Status Report ==="
echo "Generated: $(date)"
echo ""

# List of instrument workstations
INSTRUMENTS=(
  "zeiss-microscope-1"
  "zeiss-microscope-2"
  "illumina-sequencer"
  "flow-cytometer"
  "mass-spec"
)

for instrument in "${INSTRUMENTS[@]}"; do
  echo "--- ${instrument} ---"

  # SSH to instrument and check Cicada status
  ssh admin@${instrument} "
    if systemctl is-active cicada-watch >/dev/null 2>&1; then
      echo 'Status: ✓ Running'
      cicada watch list | tail -n +3
    else
      echo 'Status: ✗ NOT RUNNING'
    fi
  "

  echo ""
done

echo "=== Report Complete ==="
```

**What Lisa thinks**: *"This gives me a quick overview of all instrument backups."*

---

**Step 3: Weekly Backup Report**

Example output from monitoring script:

```
=== Lab Backup Status Report ===
Generated: Fri Nov 29 09:00:00 PST 2025

--- zeiss-microscope-1 ---
Status: ✓ Running
ID: /mnt/zeiss/output-1732397890
  Source: /mnt/zeiss/output
  Destination: s3://lab-data-backup/microscopy/zeiss
  Status: Running
  Last sync: 15 minutes ago (8 files, 3.2 GB)
  Total syncs: 342
  Total data: 1.8 TB

--- zeiss-microscope-2 ---
Status: ✓ Running
ID: /mnt/zeiss/output-1732398123
  Last sync: 2 hours ago (12 files, 5.1 GB)
  Total syncs: 298
  Total data: 1.5 TB

--- illumina-sequencer ---
Status: ✓ Running
ID: /data/illumina/output-1732398456
  Last sync: 30 minutes ago (1 file, 42 GB)
  Total syncs: 156
  Total data: 6.7 TB

--- flow-cytometer ---
Status: ✓ Running
ID: /data/cytometer/output-1732398789
  Last sync: 1 hour ago (45 files, 890 MB)
  Total syncs: 421
  Total data: 387 GB

--- mass-spec ---
Status: ✗ NOT RUNNING

=== Report Complete ===
```

**What Lisa thinks**: *"Most instruments are working great. Need to check on the mass-spec workstation."*

---

### Key Benefits for Lisa

✅ **Centralized monitoring**: Single script monitors all instruments
✅ **Compliance**: Automated backups ensure 24-hour policy compliance
✅ **Audit trail**: Watch statistics provide backup verification
✅ **Alerting**: Can integrate monitoring with Slack/email alerts
✅ **Scalability**: Easy to add new instruments with standard configs

---

## Summary: Common Patterns

### All Personas Benefit From:

1. **Simple Installation**
   - Homebrew (macOS), apt/yum (Linux), or direct binary download
   - Single binary, no dependencies

2. **Flexible Configuration**
   - YAML config files
   - Environment-specific settings
   - AWS profile integration

3. **Reliable Transfers**
   - MD5/ETag verification
   - Concurrent transfers
   - Resume capability (via S3 multipart for large files)

4. **Automatic Watching**
   - File system monitoring
   - Debouncing and min-age filtering
   - Persistent configuration

5. **Verification**
   - Dry-run mode
   - Built-in integrity checking
   - Clear status reporting

---

## Getting Started

Choose the scenario that best matches your use case and follow the walkthrough. All features shown are available in **Cicada v0.1.0**.

**Installation**:
- macOS: `brew install scttfrdmn/tap/cicada`
- Linux/Windows: Download from [releases](https://github.com/scttfrdmn/cicada/releases)

**Documentation**:
- [README.md](../README.md) - Full documentation
- [VISION.md](../VISION.md) - Future roadmap
- [GitHub Issues](https://github.com/scttfrdmn/cicada/issues) - Report problems or request features
