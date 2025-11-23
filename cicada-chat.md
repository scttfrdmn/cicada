
New chat
Think about a small  academic lab (8-10) people, none particularly technical when it comes to computers. They have a pile of data collected over years, some organized, some not, stored on USB drives, laptops, maybe a small NAS that is not been managed in years or is beyond fixing or bothering with anymore. Backups are performed never or very rarely and with no real schedule - basically just a copy.
I want you to think about an AWS based tool/deployment that allows a group like this to make sense of their data and protect it with sensible tiering (immediate tiering for old data). Think also about streamlining their computational workflows as an enhancement. Right now they manually move data from an instrument with limited storage to a laptop, then manually move data to a SaaS platform to do analysis, again with limited storage, then move data back manually and store it... somewhere.
One challenge, potentially, is that while AWS S3 is the ideal choice for storing research data many persona's like this are not familar with object based storage and are used to POSIX based filesystems. 
Anyway, interested to hear your thoughts

Show more
This is a really interesting challenge that I see all the time in academic research. Let me think through a practical solution that balances AWS capabilities with the realities of non-technical users.

The Core Problem
You've identified the key tension: S3 is perfect for this use case (durable, cheap, lifecycle policies, versioning), but researchers think in files and folders, not objects and buckets. They need something that "just works" without becoming a second job.

Proposed Architecture
Storage Layer:

AWS S3 as the backend with intelligent tiering and lifecycle policies
AWS Storage Gateway (File Gateway) as the bridge - this is the critical piece. It presents S3 as an NFS/SMB share that looks exactly like a network drive to users
AWS Backup for point-in-time snapshots and disaster recovery
Access Layer:

File Gateway mounted on lab workstations as a network drive (e.g., Z:\LabData)
Supports both Windows and Mac clients
Local cache on the gateway appliance for frequently accessed files
Intelligence Layer:

S3 Lifecycle policies that automatically move data to Infrequent Access after 90 days, then Glacier after 1 year
S3 Versioning enabled to protect against accidental deletions
AWS Macie (optional) to discover and classify sensitive data
Compute Enhancement:

AWS Transfer Family (SFTP endpoint) for instrument-to-cloud direct transfers
AWS Batch or EC2 spot instances for computational workflows
Step Functions to orchestrate multi-step analysis pipelines
What This Looks Like in Practice
Day 1 (Migration):

Users copy existing data from USB drives/laptops to Z:\LabData\Archive
Gateway uploads to S3 in the background
Lifecycle policy immediately evaluates age and tiers old data
Daily Workflow:

Instrument collects data → researcher exports to Z:\LabData\RawData\2025-11\experiment_name
Data automatically syncs to S3 (with versioning)
Analysis happens either:
Locally with cached files (gateway is smart about this)
In AWS using Batch jobs that read directly from S3
Via SaaS platform, but now with automated sync scripts instead of manual uploads
Backup Reality:

It just happens. S3 has 11 nines of durability
AWS Backup creates daily snapshots for 30 days (configurable)
Cross-region replication (optional) for disaster recovery
Cost Considerations
For a 10TB dataset (rough estimate):

S3 Standard: ~$230/month for active data (500GB)
S3 Intelligent-Tiering: ~$50/month for older data (9.5TB after auto-tiering)
File Gateway: ~$125/month (medium appliance on EC2)
Data transfer: Varies, but ingress is free
Total: ~$400-500/month for comprehensive, hands-off storage and backup
Compare this to: hardware NAS ($3-5K upfront) + inevitable failure + someone's time managing it + risk of data loss.

Implementation Approach
Phase 1: Storage Foundation (Week 1-2)

Deploy File Gateway as EC2 instance or hardware appliance
Create S3 bucket with versioning and lifecycle policies
Mount gateway as network drive on 2-3 pilot workstations
Migrate one project's data as proof of concept
Phase 2: Workflow Integration (Week 3-4)

Set up automated instrument data transfer (scripts or Transfer Family)
Create simple shell scripts or Python scripts for common tasks
Document the new "save data here" workflow with visual guides
Phase 3: Compute Enhancement (Optional, Month 2+)

Identify 1-2 computational bottlenecks
Create AWS Batch job definitions for common analyses
Build simple web form or CLI tool to submit jobs
Key Success Factors
Make it invisible: The best technology is the one users don't think about. Z:\LabData should just work.
No new concepts: Users never need to hear "S3 bucket" or "object storage." It's just "the lab drive."
Fail-safe defaults: Lifecycle policies, versioning, and backups happen automatically. Users can't opt out or forget.
One champion: You need one semi-technical person (grad student, postdoc, lab manager) who understands the system and can troubleshoot basics.
Document the happy path: Create a one-page guide with screenshots showing exactly where to save data and how to access it.
Alternative: Fully Managed Option
If even File Gateway feels too complex, consider AWS FSx for Lustre or FSx for NetApp ONTAP - these are fully managed file systems that integrate with S3 but handle all the complexity. Trade-off is higher cost (~$1000-1500/month for similar capacity) but zero management overhead.

The Migration Carrot
Pitch this to the PI as: "We can make it so your lab never loses data again, old data costs almost nothing to store, and we can run bigger analyses faster - all for less than buying another NAS that will fail in 3 years."

What aspects would you like me to dig deeper into? The workflow automation piece could be really powerful for streamlining that instrument→analysis→storage cycle.

Now I am going to add a twist - academic labs have very little money. That is why I mentioned creating a tool that manages a lot of this - a new tell we will develop potentially. Another aspect of groups like this is their activity is very bursty. Having a File Gateway running all the time and an appliance is overkill in more ways than one - for example. 
This should definitely "just work" but now think and be clever about how a group like this could have a tool that deploys, manages, monitors this for as little as possible in terms of cost.

Ah, now we're talking! This is the real constraint that makes it interesting. Let me rethink this as a serverless, on-demand, pay-per-use system that acknowledges the bursty reality of academic research.

The Actual Reality
Most data transfer happens in bursts: after experiments, before paper submissions, when a hard drive is dying
Compute needs spike around deadlines, then nothing for weeks
Budget is maybe $50-100/month, not $500
"Always-on" anything is wasteful for this use pattern
Revised Architecture: The "Dormant Until Needed" Approach
Storage Layer (Always On, But Cheap)
S3 with Intelligent-Tiering from Day 1 - this is your only "always on" component
S3 versioning - because it's essentially free insurance
No File Gateway running 24/7
Access Layer (On-Demand)
Option A: Mountpoint for Amazon S3 (New, Linux/Mac)

Free, open-source POSIX interface to S3
Users run it locally: mount-s3 lab-bucket ~/LabData
Only runs when they need it, on their own machines
No EC2 costs at all
Option B: Rclone (Cross-platform)

Free, mature, works everywhere
Can mount S3 as local filesystem: rclone mount s3:lab-bucket ~/LabData
Also excellent for sync operations
Again, runs locally - no infrastructure costs
Option C: AWS File Gateway, But On-Demand

Deploy as EC2 instance with your tool
Only runs when actively transferring data
Spin up for bulk transfers, shut down after
Cost: ~$0.10/hour × maybe 20 hours/month = $2/month instead of $125
The Tool You'd Build
Let me describe "LabSync" - a CLI tool and optional simple web UI that manages everything:

bash
# Install once
pip install labsync
labsync init --lab-name "smithlab"

# Daily usage - researchers don't think about AWS
labsync mount                    # Mounts S3 locally via rclone/mountpoint
labsync upload /path/to/data     # Smart upload with deduplication
labsync sync-instrument          # Auto-detects instrument data, uploads
labsync compute --script analyze.py --input data/  # Runs job in cloud
labsync status                   # Shows storage costs, what's tiered where
```

**What LabSync Actually Does:**

1. **Intelligent Mounting**
   - Detects if mountpoint-s3 or rclone is installed, uses the free option
   - Creates local mount point automatically
   - Handles AWS credentials via AWS SSO or long-term credentials (simplified setup)

2. **Smart Uploads**
   - Checksums files to avoid re-uploading
   - Automatically sets S3 storage class based on file age/size
   - Progress bars, resume on failure
   - Automatically tags files with metadata (date, instrument type, project)

3. **On-Demand File Gateway**
   - For bulk migrations (hundreds of GB), spins up EC2 File Gateway
   - CloudFormation template that deploys everything
   - Automatically shuts down after 2 hours of inactivity
   - Email notification when complete: "Your data is uploaded, gateway is shutting down"

4. **Burst Compute**
   - Detects common analysis scripts (Python, R, MATLAB)
   - Creates Docker container on-the-fly
   - Submits to AWS Batch using spot instances
   - Automatically pulls input from S3, pushes results back
   - Cost: pennies per job vs hours of laptop time

5. **Cost Monitoring**
   - Daily email: "This month: $12.47 (3.2TB stored, 45GB transferred)"
   - Warns if costs spike: "Unusual activity detected: 500GB downloaded yesterday"
   - Shows cost breakdown: "You could save $8/month by archiving data older than 2 years"

6. **Lifecycle Management**
   - Automatically sets S3 lifecycle policies on bucket creation
   - Default: Intelligent-Tiering immediately
   - Option to auto-archive to Glacier after 1 year
   - Never deletes anything (academics are paranoid about this)

## Real-World Cost Breakdown

**10TB dataset, bursty usage:**

| Component | Monthly Cost |
|-----------|-------------|
| S3 Intelligent-Tiering (avg) | ~$80 |
| S3 API requests | ~$2 |
| Data transfer OUT (minimal) | ~$5 |
| EC2 for compute (20 spot hours) | ~$3 |
| File Gateway (4 hours/month) | ~$0.40 |
| **Total** | **~$90/month** |

**Compared to alternatives:**
- Dropbox for 10TB: $240/month
- Backblaze B2: ~$50-60/month (but none of the compute benefits)
- Local NAS: $0/month until it dies and you lose everything

## LabSync Implementation Details

**Architecture:**
```
┌─────────────────┐
│  Researcher's   │
│    Laptop       │
│                 │
│  ┌───────────┐  │
│  │ LabSync   │  │──┐
│  │   CLI     │  │  │ Mounts S3 locally via
│  └───────────┘  │  │ rclone/mountpoint
└─────────────────┘  │
                     │
                     ▼
         ┌────────────────────┐
         │   AWS Account      │
         │                    │
         │  ┌──────────────┐  │
         │  │  S3 Bucket   │  │ ◄── Lifecycle policies
         │  │  (Primary)   │  │     Intelligent-Tiering
         │  └──────────────┘  │     Versioning enabled
         │                    │
         │  ┌──────────────┐  │
         │  │  Lambda      │  │ ◄── Triggered by uploads
         │  │  (Metadata)  │  │     Tag extraction
         │  └──────────────┘  │     Cost tracking
         │                    │
         │  ┌──────────────┐  │
         │  │  Batch       │  │ ◄── On-demand compute
         │  │  (Spot)      │  │     Auto-shutdown
         │  └──────────────┘  │
         │                    │
         │  ┌──────────────┐  │
         │  │  CloudWatch  │  │ ◄── Cost alerts
         │  │  (Alerts)    │  │     Usage monitoring
         │  └──────────────┘  │
         └────────────────────┘
Key Features for Cost Control:

Aggressive Intelligent-Tiering
Start files in Infrequent Access tier if older than 30 days at upload
Let S3 automatically move to Archive tiers
Deduplication Detection
Before upload, hash file and check if already in S3
Prevents uploading same dataset multiple times
"File already exists in bucket, skipping upload"
Partial File Sync
Only upload changed files
Use S3 ETags to detect differences
Multipart upload for large files with resume capability
Spot-Only Compute
Never use on-demand EC2
Batch jobs configured for spot instances only
Accept occasional interruptions (rare for short jobs)
Automatic Cleanup
Lambda function runs weekly: delete incomplete multipart uploads
These accumulate and cost money if not cleaned up
LabSync does this automatically
The Deployment Tool Component
labsync deploy sets up the entire AWS infrastructure:

bash
labsync deploy --budget 100  # Sets up cost alerts at $100/month
What it creates (via CloudFormation/Terraform):

S3 bucket with sensible defaults
IAM roles with least-privilege access
CloudWatch cost alerts
Lambda function for metadata extraction
Batch compute environment (disabled until first use)
EventBridge rules for automated tiering
Cost of infrastructure itself: $0 - everything is pay-per-use

The "Instrument Sync" Feature
This is where you really save researchers time:

bash
# Run once to configure
labsync add-instrument --type "zeiss-microscope" --path "/Volumes/Instrument/Data"

# Then just:
labsync sync-instrument
What it does:

Watches instrument export folder
Automatically uploads new files when they appear
Adds metadata tags (instrument type, date, user)
Can run on a schedule (cron job) or manually
Progress bar shows upload status
Alternative for Windows users:

Simple Windows Task Scheduler task runs labsync sync-instrument daily
Or: LabSync runs as background service (optional)
The Compute Enhancement
Simple case:

bash
labsync compute --script process_images.py --input "s3://bucket/raw/exp_001/"
What happens:

LabSync packages script + dependencies into Docker image
Pushes to ECR (Elastic Container Registry)
Submits Batch job with spot instances
Job runs, outputs to S3
Email notification when complete
Total cost: $0.50 for what would take 6 hours on laptop
Advanced case:

yaml
# labsync.yaml in project directory
compute:
  environment: python:3.11
  dependencies:
    - numpy
    - scikit-image
    - pytorch
  resources:
    vcpus: 8
    memory: 32GB
    gpu: false
  spot: true
Then just: labsync compute --config labsync.yaml

Migration Path
Week 1: Initial Setup

bash
pip install labsync
labsync init --lab-name "smithlab"
labsync mount
# Copy old data from USB drives to mounted folder
Week 2: Instrument Integration

bash
labsync add-instrument --type microscope --path /path/to/instrument
# Set up daily sync (LabSync creates cron job)
Week 3: First Compute Job

bash
labsync compute --script analyze.py --input "recent_data/"
# See how much faster/cheaper it is
Open Source & Community
Since this is for academics, LabSync should be:

Open source (MIT license)
Documented with examples from real labs
Community plugins for specific instruments
Templates for common analysis workflows
Other labs contribute instrument integrations, analysis scripts, cost optimization tips.

The Business Model (If You Want One)
Core tool: Free and open source
Hosted web UI: $10/month per lab (optional, just makes it easier)
Support tier: $50/month for setup help and custom integrations
Enterprise features: SSO, compliance reporting, multi-lab organizations
But honestly, just open-sourcing it and having happy users might be the goal.

Real-World Example
Lab scenario:

5TB of microscopy images collected over 3 years
New experiments generate 50GB/week during active periods
Analysis jobs run once a month, take 8 hours on laptop
Before LabSync:

Data scattered across laptops and dying external drives
No backups
Analysis blocks laptop for a day
No idea what data they even have
After LabSync:

One-time bulk upload: 5TB → S3 Intelligent-Tiering
Weekly syncs during active research: automatic
Monthly analysis: labsync compute, done in 1 hour on spot instances
Cost: ~$60/month for storage, $2/month for compute
Data is searchable, tagged, safe
Technical Challenges to Solve
AWS Credentials for Non-Technical Users
LabSync wizard: "Enter your AWS access key once, we'll store it securely"
Or: AWS SSO integration for universities that have it
Or: Provide pre-configured CloudFormation that creates limited IAM user
Windows/Mac/Linux Compatibility
Python CLI works everywhere
Graceful fallback if mountpoint-s3 isn't available
Windows Service vs macOS LaunchAgent vs Linux systemd
Large File Handling
Multipart upload with resume
Progress indication that doesn't spam terminal
Bandwidth throttling option (don't saturate university network)
Offline Operation
Local cache of frequently accessed files
Queue uploads when internet is spotty
Sync when connection restored
The Pitch to Academic Labs
"LabSync gives you enterprise-grade data storage and compute for $50-100/month. No hardware to maintain, no data to lose, scales with your needs. When you're not using it, you're barely paying for it."

Would you want me to sketch out the actual Python package structure for LabSync, or dive deeper into specific components like the on-demand File Gateway orchestration or the instrument sync feature?

I am thinking this should be called Cicada - because the reality is, it lies dormant most of the time. Implemented in Go. One nice feature would be something that allows:
1. periodcially checking a folder (or instrument) for new data and automatically triggering a data movement
2. rsync like behavior
3. Not only workflow execution (using say, snakemake, etc) but visualization/workstation capability

Cicada - I love it! The metaphor is perfect: emerges when needed, dormant otherwise, and when it does emerge, it's remarkably efficient.

Let me think through a Go implementation that really embraces this concept.

Why Go Is Perfect For This
Single binary deployment (no Python dependency hell for academics)
Cross-platform compilation (build once, run on Windows/Mac/Linux)
Excellent concurrency for file watching and parallel uploads
Fast, low memory footprint (can run as background daemon without being noticed)
Great libraries for S3 (aws-sdk-go-v2), file watching (fsnotify), and sync operations
Cicada Architecture
cicada/
├── cmd/
│   ├── cicada/           # Main CLI entry point
│   ├── cicada-daemon/    # Background service for watching
│   └── cicada-gateway/   # On-demand file gateway orchestrator
├── internal/
│   ├── sync/            # rsync-like engine
│   ├── watch/           # File system watcher
│   ├── storage/         # S3 operations with intelligent tiering
│   ├── compute/         # Workflow execution (Batch, Fargate)
│   ├── workstation/     # Remote desktop/visualization
│   ├── gateway/         # File Gateway lifecycle management
│   └── config/          # Configuration management
└── pkg/
    ├── instrument/      # Pluggable instrument adapters
    └── workflow/        # Snakemake, Nextflow, CWL integrations
Core Features
1. Intelligent File Watching & Auto-Sync
bash
# Set up a watch on instrument folder
cicada watch add microscope-1 \
  --path /Volumes/ZeissMicroscope/Export \
  --sync-on-new \
  --min-age 5m \
  --sync-schedule "0 2 * * *"  # Also sync at 2 AM daily
What this does:

Watches folder for new files using fsnotify
Waits 5 minutes after file creation (ensures writes are complete)
Automatically syncs to S3 when stable
Also runs full sync at 2 AM to catch anything missed
Runs as daemon in background (minimal CPU/memory)
Smart behaviors:

go
type WatchConfig struct {
    Path          string
    MinAge        time.Duration  // Don't sync files newer than this
    IgnorePattern []string       // .tmp, .partial, etc.
    OnComplete    []string       // Commands to run after sync
    Checksum      bool           // Verify uploads with checksums
    DeleteSource  bool           // Delete local after successful upload
    Bandwidth     int64          // Throttle to avoid saturating network
}
Daemon mode:

bash
# Start daemon (runs in background)
cicada daemon start

# Check status
cicada daemon status
# Output: 
# Cicada daemon running (PID: 1234)
# Watching: 2 locations
#   - /Volumes/Microscope/Export → s3://lab-bucket/raw/microscope/
#   - /data/sequencer/output → s3://lab-bucket/raw/sequencing/
# Last sync: 3 files, 2.3 GB, 12 minutes ago
# Next scheduled sync: today at 2:00 AM

# View logs
cicada daemon logs --follow
2. rsync-Like Sync Engine
bash
# Basic sync
cicada sync /local/data s3://lab-bucket/data

# Dry run to see what would be transferred
cicada sync /local/data s3://lab-bucket/data --dry-run

# Bidirectional sync
cicada sync /local/data s3://lab-bucket/data --bidirectional

# Exclude patterns
cicada sync /local/data s3://lab-bucket/data \
  --exclude "*.tmp" \
  --exclude ".DS_Store" \
  --delete  # Delete remote files not in source
Smart sync features:

Checksum-based comparison (not just timestamp/size)
Incremental uploads with multipart for large files
Resume capability for interrupted transfers
Parallel transfers (configurable concurrency)
Bandwidth limiting (crucial for not disrupting lab network)
Deduplication detection (don't upload if identical file exists)
Implementation detail:

go
type SyncEngine struct {
    Source      string
    Destination string
    Options     SyncOptions
}

type SyncOptions struct {
    DryRun      bool
    Delete      bool
    Checksum    bool
    Exclude     []string
    Include     []string
    Concurrency int
    Bandwidth   int64  // bytes per second limit
}

// Smart delta detection
func (s *SyncEngine) buildTransferPlan() (*TransferPlan, error) {
    // 1. Walk source tree
    // 2. Build manifest with checksums
    // 3. Compare with S3 ETags
    // 4. Generate minimal transfer list
    // 5. Estimate time/cost
}
```

**Progress output:**
```
Scanning source... 1,247 files (15.3 GB)
Comparing with remote... 1,198 files unchanged
Transfer plan:
  Upload: 49 files (873 MB)
  Delete: 3 files (remote only)
  Estimated time: 8m 32s
  Estimated cost: $0.02

Proceed? [Y/n] y

Uploading: [=====>    ] 45% (2.3 GB/s) ETA: 3m 12s
  ✓ experiment_001/image_0045.tif (128 MB)
  ↻ experiment_002/data.zarr (multipart, 2.1 GB)
  ⏸ experiment_003/video.avi (queued)
3. Workflow Execution with Snakemake/Nextflow
This is where Cicada gets really powerful - bridging the gap between local workflows and cloud compute:

bash
# Execute a Snakemake workflow in AWS
cicada workflow run snakemake \
  --snakefile Snakefile \
  --config input=s3://lab-bucket/raw/experiment_123/ \
  --cores 32 \
  --memory 64GB \
  --spot

# Or with a local workflow file
cicada workflow run nextflow \
  --workflow pipeline.nf \
  --input "s3://lab-bucket/raw/**.fastq.gz" \
  --outdir s3://lab-bucket/results/run_456/
What Cicada does:

Parses workflow file to understand dependencies
Packages workflow + dependencies into container
Provisions AWS Batch compute environment (if not exists)
Submits job graph to Batch using spot instances
Monitors progress, streams logs
Notifies on completion
Tears down compute environment after idle period
Local workflow development:

bash
# Test workflow locally first
cicada workflow run snakemake \
  --snakefile Snakefile \
  --local \
  --cores 4

# When ready, run in cloud with same command (remove --local)
Workflow config file (for repeated runs):

yaml
# cicada-workflow.yaml
name: image-processing-pipeline
engine: snakemake
workflow: Snakefile

compute:
  type: batch
  instance_types: [c5.4xlarge, c5.9xlarge]
  spot: true
  max_vcpus: 256
  min_vcpus: 0  # Scale to zero when not in use

storage:
  input: s3://lab-bucket/raw/
  output: s3://lab-bucket/processed/
  working: s3://lab-bucket/scratch/  # Auto-delete after 7 days

notifications:
  email: lab@university.edu
  slack: "#lab-notifications"

cost_limit: 50  # Abort if estimated cost exceeds $50
Then just:

bash
cicada workflow run --config cicada-workflow.yaml
4. Remote Workstation/Visualization
This is the killer feature - many academic workflows need GUI tools (ImageJ, MATLAB, Paraview, etc.):

bash
# Launch a Linux workstation with GPU for visualization
cicada workstation launch \
  --name viz-session \
  --instance g4dn.xlarge \
  --image nvidia-opengl \
  --storage s3://lab-bucket/data/experiment_123/ \
  --spot

# Output:
# Launching workstation... 
# Instance starting: i-0abc123def456
# Waiting for connection... ready!
# 
# Connect via:
#   Web:  https://viz-session-abc123.cicada.cloud
#   VNC:  vnc://54.123.45.67:5901
#   SSH:  ssh cicada@54.123.45.67
#
# Auto-shutdown: 2 hours of inactivity
# Cost: ~$0.50/hour (spot pricing)
What you get:

Full Linux desktop (XFCE or similar, lightweight)
Pre-installed scientific tools (ImageJ, ParaView, Python with viz libraries)
S3 data mounted via s3fs or goofys (appears as local folder)
GPU acceleration for rendering
Browser-based access (no VNC client needed)
Implementation using:

AWS EC2 with GPU instances (g4dn family)
DCV (AWS's remote desktop protocol) or noVNC for browser access
CloudFormation template for quick provisioning
Auto-shutdown Lambda that monitors idle time
Pre-built images:

bash
cicada workstation images

# Output:
# Available workstation images:
#   - basic-linux      Ubuntu 22.04, basic tools
#   - imagej           ImageJ, FIJI, common plugins
#   - matlab           MATLAB R2024a (requires license)
#   - paraview         ParaView, VTK visualization
#   - rstudio          RStudio Server, tidyverse
#   - jupyter          JupyterLab, scipy stack
#   - custom           Build your own from Dockerfile
Custom image:

bash
# Create custom workstation image
cicada workstation build \
  --from imagej \
  --add-package napari \
  --add-package cellpose \
  --name microscopy-workstation

# Use it
cicada workstation launch --image microscopy-workstation
Session management:

bash
# List active sessions
cicada workstation list

# Reconnect to existing session
cicada workstation connect viz-session

# Extend auto-shutdown timer
cicada workstation extend viz-session --hours 4

# Manual shutdown (to save costs)
cicada workstation stop viz-session

# Resume a stopped session (retains data on EBS)
cicada workstation start viz-session
Snapshot feature:

bash
# Save session state for later
cicada workstation snapshot viz-session \
  --name "before-batch-processing"

# Restore later
cicada workstation launch \
  --from-snapshot "before-batch-processing"
The On-Demand File Gateway
For those bulk migrations where rclone/mountpoint aren't ideal:

bash
cicada gateway start

# Output:
# Deploying File Gateway...
# CloudFormation stack: cicada-gateway-abc123
# EC2 instance: i-0def789ghi012 (m5.large)
# Status: Configuring... (2-3 minutes)
# 
# Gateway ready!
# Mount on Mac/Linux:
#   mount -t nfs 10.0.1.45:/lab-bucket ~/LabData
# 
# Mount on Windows:
#   net use Z: \\10.0.1.45\lab-bucket
#
# Auto-shutdown: 2 hours of inactivity
# Cost while running: ~$0.10/hour

# Check status
cicada gateway status

# Manual shutdown
cicada gateway stop
Smart gateway management:

Automatically provisions via CloudFormation
Uses spot instances when possible
Monitors activity, auto-shuts down when idle
Can be scheduled: cicada gateway start --schedule "Mon-Fri 9AM-5PM"
Email notification when shutting down
Configuration & Setup
Initial Setup Wizard
bash
cicada init

# Interactive wizard:
# 
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
# Step 4: Optional Features
#   Enable compute workflows? [Y/n] y
#   Enable remote workstations? [Y/n] y
#   Enable file gateway? [y/N] n
# 
# Setup complete! 
# 
# Next steps:
#   1. Start daemon: cicada daemon start
#   2. Add watch location: cicada watch add --path /path/to/data
#   3. Test sync: cicada sync /local/folder s3://smith-lab-data-20241122/test/
Configuration File
yaml
# ~/.cicada/config.yaml
lab:
  name: smith-lab
  bucket: smith-lab-data-20241122
  region: us-west-2

storage:
  intelligent_tiering: true
  versioning: true
  lifecycle:
    - rule: archive-old-data
      days_to_glacier: 365
      enabled: true

sync:
  default_concurrency: 10
  bandwidth_limit: 100MB  # Don't saturate network
  checksum_verify: true
  exclude_patterns:
    - "*.tmp"
    - ".DS_Store"
    - "Thumbs.db"

compute:
  spot_preferred: true
  max_vcpus: 256
  instance_types: [c5.xlarge, c5.2xlarge, c5.4xlarge]

workstation:
  default_instance: g4dn.xlarge
  auto_shutdown_hours: 2
  allow_spot: true

notifications:
  email: pi@university.edu
  slack_webhook: https://hooks.slack.com/...

cost:
  budget_limit: 100
  alert_threshold: 80
Instrument Adapters (Pluggable)
go
// pkg/instrument/adapter.go
type InstrumentAdapter interface {
    Name() string
    Detect(path string) bool
    ParseMetadata(file string) (Metadata, error)
    ShouldSync(file string) bool
}

// Example: Zeiss microscope adapter
type ZeissMicroscope struct{}

func (z *ZeissMicroscope) ParseMetadata(file string) (Metadata, error) {
    // Read CZI file metadata
    // Extract: magnification, timestamp, operator, etc.
    return metadata, nil
}

func (z *ZeissMicroscope) ShouldSync(file string) bool {
    // Only sync .czi files, ignore thumbnails
    return strings.HasSuffix(file, ".czi")
}
User creates custom adapter:

bash
# Generate adapter template
cicada instrument create-adapter \
  --name custom-hplc \
  --template generic

# Edit generated file
# ~/.cicada/instruments/custom-hplc.go

# Build and register
cicada instrument install custom-hplc
Cost Tracking & Reporting
bash
cicada cost report

# Output:
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
#   Data ingress:             $0.00   (125 GB)
#   Data egress:              $2.15   (24 GB)
# 
# Total: $84.50 / $100.00 budget
# 
# Trends:
#   ↓ 12% vs last month
#   Storage growing ~50 GB/week
#   Compute usage steady
# 
# Recommendations:
#   ⚡ Archive data older than 2 years → save $15/month
#   💡 Use smaller workstation instance → save $5/month
Cost prediction:

bash
cicada cost predict --action "archive data older than 2 years"

# Estimated savings: $15.20/month
# One-time glacier deep archive cost: $2.50
# Payback period: 1 month
# 
# Apply this recommendation? [Y/n]
Example Workflows
Scenario 1: Microscopy Lab Daily Workflow
Setup (once):

bash
cicada daemon start
cicada watch add microscope \
  --path /Volumes/Microscope/Export \
  --sync-on-new \
  --min-age 10m \
  --delete-source  # Free up instrument storage
Daily operation:

Researcher saves images from microscope → instrument folder
Cicada daemon detects new files after 10 minutes
Automatically uploads to S3 with metadata tags
Deletes from local instrument (frees space)
Email notification: "Uploaded experiment_456: 45 files, 8.2 GB"
Weekly analysis:

bash
cicada workflow run snakemake \
  --snakefile ~/pipelines/cell-segmentation.smk \
  --config experiment=experiment_456 \
  --spot
Monthly visualization:

bash
cicada workstation launch \
  --image napari-workstation \
  --data s3://bucket/processed/experiment_456/
Scenario 2: Sequencing Core
Automated pipeline:

yaml
# .cicada/auto-pipeline.yaml
name: auto-fastq-processing
trigger:
  watch: /data/sequencer/BaseCalls/
  pattern: "*.fastq.gz"
  min_files: 2  # Wait for R1 and R2

workflow:
  engine: nextflow
  pipeline: nf-core/rnaseq
  params:
    genome: GRCh38
    outdir: s3://lab-bucket/results/${RUN_ID}/

compute:
  spot: true
  max_cpus: 64

notifications:
  on_complete:
    email: bioinformatics@lab.edu
    message: "Sequencing run ${RUN_ID} complete"
bash
cicada workflow enable auto-pipeline
Now sequencer outputs are automatically processed when complete.

Scenario 3: Collaborative Analysis
PI shares workstation:

bash
# PI launches workstation with shared access
cicada workstation launch \
  --name group-analysis \
  --instance r5.2xlarge \
  --data s3://bucket/project_x/ \
  --users "postdoc@lab.edu,gradstudent@lab.edu" \
  --duration 8h

# Team members connect
cicada workstation connect group-analysis
Everyone sees the same filesystem, can run analyses together, results save to S3.

Implementation Priorities
Phase 1: MVP (4-6 weeks)

Basic CLI structure
S3 sync engine with rsync-like features
File watching and auto-sync
Daemon mode
Simple cost tracking
Phase 2: Compute (4 weeks)

Workflow execution (Snakemake first)
AWS Batch integration
Spot instance management
Phase 3: Workstations (4 weeks)

EC2 provisioning for remote desktop
Browser-based access
Auto-shutdown logic
Pre-built images
Phase 4: Polish (ongoing)

More instrument adapters
Better cost optimization
Community contributions
Documentation and tutorials
Go Package Structure
go
// cmd/cicada/main.go
package main

import (
    "github.com/lab/cicada/internal/cli"
)

func main() {
    cli.Execute()
}

// internal/sync/engine.go
package sync

type Engine struct {
    source      storage.Backend
    destination storage.Backend
    options     Options
}

func (e *Engine) Sync(ctx context.Context) error {
    plan, err := e.buildPlan(ctx)
    if err != nil {
        return err
    }
    
    return e.executePlan(ctx, plan)
}

// internal/watch/watcher.go
package watch

type Watcher struct {
    locations map[string]*WatchConfig
    syncer    *sync.Engine
}

func (w *Watcher) Start(ctx context.Context) error {
    // Use fsnotify to watch configured paths
    // Debounce events
    // Trigger syncs
}
Would you like me to dive deeper into any specific component? For example:

The rsync-like sync algorithm implementation
File Gateway orchestration (CloudFormation templates)
Workstation provisioning and DCV setup
The daemon's event processing logic
Cost calculation and prediction algorithms
This is basically a mini data commons platform isn't it?
What about users and groups in the lab? IAM mappings? more for organizing projects than anything else. They are probably not sophisticated enough for things like SSO (but could maybe support using Globus Auth easily)
The one thing left out is a GUI front end - could be web based off of the local daemon?
Lean & mean, simple, does what this persona intuits

You're absolutely right - this is a mini data commons platform! And you've identified the missing pieces that would make it actually usable for this persona.

Let me think through this more holistically...

What This Really Is
It's a lightweight, academic-focused data commons that provides:

Federated storage (everyone's data in one logical place)
Access control (who can see/use what)
Compute-to-data (bring analysis to the data, not data to analysis)
Collaboration primitives (shared workspaces, reproducible workflows)
But unlike traditional data commons platforms (which are heavyweight, require IT staff, cost $$$), this is:

Self-service: PI installs it, it just works
Cost-aware: Designed for shoestring budgets
Opportunistic: Uses resources only when needed
Users, Groups, and Projects
Academic labs naturally organize around projects, not just individuals. Here's how Cicada should model this:

Conceptual Model
Lab (Organization)
├── Members (Users)
│   ├── PI (admin)
│   ├── Postdocs
│   ├── Grad Students
│   ├── Undergrads
│   └── Collaborators (external)
├── Groups
│   ├── #protein-structure (PI, 2 postdocs, 3 grad students)
│   ├── #metabolism (PI, 1 postdoc, 2 grad students)
│   └── #methods-dev (postdoc lead, 4 students)
└── Projects
    ├── NIH-R01-2024 (specific grant/paper)
    ├── student-thesis-jane
    └── exploratory-experiments
Storage Organization
s3://smith-lab-data/
├── shared/              # Lab-wide access
│   ├── protocols/
│   ├── reagents/
│   └── instruments/
├── projects/
│   ├── NIH-R01-2024/   # Project-level access
│   │   ├── raw/
│   │   ├── processed/
│   │   └── results/
│   └── student-thesis-jane/
├── groups/
│   ├── protein-structure/  # Group-level access
│   └── metabolism/
└── users/
    ├── jsmith/         # Personal scratch space
    └── agarcia/
User Management
bash
# Initialize lab (PI does this once)
cicada lab init --name smith-lab

# Add lab members
cicada user add jsmith@university.edu \
  --role postdoc \
  --groups protein-structure,methods-dev

cicada user add agarcia@university.edu \
  --role grad-student \
  --groups metabolism

# Add external collaborator (limited access)
cicada user add external@otheruniv.edu \
  --role collaborator \
  --projects NIH-R01-2024

# List users
cicada user list

# Output:
# Lab Members (5):
#   👤 pi@university.edu (PI, admin)
#   👤 jsmith@university.edu (Postdoc)
#      Groups: protein-structure, methods-dev
#   👤 agarcia@university.edu (Grad Student)
#      Groups: metabolism
#   👤 external@otheruniv.edu (Collaborator)
#      Projects: NIH-R01-2024
Project Management
bash
# Create project
cicada project create NIH-R01-2024 \
  --description "Protein structure determination" \
  --members pi@university.edu,jsmith@university.edu \
  --budget 200  # $200/month cloud spending limit

# Add data to project
cicada sync /local/data s3://smith-lab-data/projects/NIH-R01-2024/raw/

# Grant collaborator access
cicada project add-member NIH-R01-2024 external@otheruniv.edu \
  --permission read-only

# Project dashboard
cicada project info NIH-R01-2024

# Output:
# Project: NIH-R01-2024
# Description: Protein structure determination
# Created: 2024-03-15
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
#   • 2 hours ago: jsmith uploaded 45 files (8.2 GB)
#   • 1 day ago: workflow completed (cell-segmentation)
#   • 3 days ago: external accessed dataset_v3
# 
# Compute (this month):
#   Workflows: 12 runs, $8.40
#   Workstations: 3 sessions, $4.20
IAM Mapping Strategy
This is where it gets clever. We need to balance security with usability for non-technical users.

Approach: Managed IAM with Simple Abstraction
Backend (automatic, invisible to users):

Cicada creates IAM roles and policies in AWS
Uses fine-grained S3 bucket policies
Everything follows principle of least privilege
Frontend (what users see):

Simple permissions: admin, read-write, read-only
Scoped to: lab-wide, group, project, or specific paths
Implementation
bash
# What user types:
cicada project add-member NIH-R01-2024 jsmith@university.edu --permission read-write
What Cicada does behind the scenes:

Creates IAM policy:
json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:PutObject", "s3:ListBucket"],
      "Resource": [
        "arn:aws:s3:::smith-lab-data/projects/NIH-R01-2024/*",
        "arn:aws:s3:::smith-lab-data/shared/*"
      ]
    }
  ]
}
```

2. Creates/updates IAM user for jsmith@university.edu
3. Attaches policy
4. Generates credentials
5. Emails jsmith with setup instructions

**For the user:**
```
Subject: You've been added to Cicada project NIH-R01-2024

Hi Jane,

You've been granted access to the NIH-R01-2024 project in Smith Lab's Cicada.

Setup (one-time):
  1. Download Cicada: https://cicada.sh/download
  2. Run: cicada login jsmith@university.edu
  3. Enter the temporary password: [generated-password]
  4. You'll be prompted to set a new password

That's it! You can now access project data:
  cicada sync s3://smith-lab-data/projects/NIH-R01-2024/ ~/local-data/

Questions? Reply to this email or ask your PI.
Authentication Options
Option 1: Cicada-Managed IAM (default, simplest)

Cicada creates IAM users automatically
Stores credentials securely in local keychain
Users never see AWS console
Option 2: Globus Auth Integration

bash
cicada lab configure-auth --provider globus

# Now users can login with institutional credentials
cicada login --provider globus
# Opens browser → logs in via university SSO → token stored
Behind the scenes:

Cicada uses Globus OAuth2 flow
Maps Globus identity to IAM role
Uses AWS STS to generate temporary credentials
Credentials auto-refresh
Option 3: Bring Your Own AWS Account

bash
# For labs with existing AWS accounts
cicada lab init --use-existing-aws \
  --account-id 123456789012 \
  --admin-role arn:aws:iam::123456789012:role/CicadaAdmin
```

## The Web UI

This is critical for adoption. Many academics are more comfortable with GUIs than CLIs.

### Architecture
```
┌─────────────────────────────────────┐
│  Browser (any device)               │
│  http://localhost:7878              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Cicada Daemon (local machine)      │
│  - HTTP API server                  │
│  - WebSocket for real-time updates  │
│  - Serves static web UI             │
└──────────────┬──────────────────────┘
               │
               ▼
         ┌─────────┐
         │   AWS   │
         └─────────┘
Key insight: The web UI is served locally by the daemon. No external server needed. Works offline (for browsing local data).

Starting the Web UI
bash
# Start daemon with web UI
cicada daemon start --web

# Output:
# Cicada daemon started
# Web UI: http://localhost:7878
# API: http://localhost:7878/api
# 
# Opening browser...

# Or just browse to it manually
open http://localhost:7878
```

### Web UI Features

**Dashboard (Home)**
```
┌─────────────────────────────────────────────────────────────┐
│ Smith Lab Data Commons                    [jsmith] [Settings]│
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  📊 Overview                                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │   Storage   │ │   Compute   │ │    Cost     │           │
│  │   12.3 TB   │ │  3 active   │ │  $78 / $100 │           │
│  │   ↑ 2.1 GB  │ │  workflows  │ │   this month│           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
│                                                               │
│  🔔 Recent Activity                                          │
│  • 2m ago  - Sync completed: microscope-1 (45 files)        │
│  • 15m ago - Workflow started: cell-segmentation             │
│  • 1h ago  - jsmith uploaded to NIH-R01-2024/raw/           │
│                                                               │
│  📁 Quick Access                                             │
│  → My Projects (3)                                           │
│  → Shared Data                                               │
│  → Recent Files                                              │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**File Browser**
```
┌─────────────────────────────────────────────────────────────┐
│ 📁 smith-lab-data / projects / NIH-R01-2024 / raw           │
├─────────────────────────────────────────────────────────────┤
│ [↑ Up] [📤 Upload] [⬇️ Download Selected] [🔄 Sync Folder]  │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ ☐ 📁 experiment_001/                    2024-11-15  8.2 GB  │
│ ☐ 📁 experiment_002/                    2024-11-18  12.1 GB │
│ ☐ 📄 metadata.csv                       2024-11-15  2.3 MB  │
│ ☐ 📄 README.md                          2024-11-14  1.2 KB  │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│ Preview: metadata.csv                                        │
│ ┌───────────────────────────────────────────────────────┐   │
│ │ experiment_id, date, operator, magnification          │   │
│ │ exp_001, 2024-11-15, jsmith, 40x                      │   │
│ │ exp_002, 2024-11-18, agarcia, 63x                     │   │
│ └───────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

**Drag & Drop Upload**
- Drag files from desktop → browser
- Shows upload progress with pause/resume
- Auto-extracts metadata from known file types
- Suggests project/folder based on filename patterns

**Data Sync Manager**
```
┌─────────────────────────────────────────────────────────────┐
│ Sync Manager                              [+ Add Watch Path] │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ 🔄 Active Watches                                            │
│                                                               │
│ ┌─ microscope-1 ────────────────────────────────────────┐   │
│ │ 📂 /Volumes/Microscope/Export                          │   │
│ │ ☁️  s3://smith-lab-data/raw/microscope/                │   │
│ │                                                         │   │
│ │ Status: ✓ Synced (2 minutes ago)                       │   │
│ │ Last sync: 45 files, 8.2 GB                            │   │
│ │ Next sync: on new files (watching...)                  │   │
│ │                                                         │   │
│ │ [⚙️ Settings] [⏸️ Pause] [📊 History]                    │   │
│ └─────────────────────────────────────────────────────────┘   │
│                                                               │
│ ┌─ sequencer-output ────────────────────────────────────┐   │
│ │ 📂 /data/sequencer/BaseCalls/                          │   │
│ │ ☁️  s3://smith-lab-data/raw/sequencing/                │   │
│ │                                                         │   │
│ │ Status: 🟡 Syncing... (2.1 GB / 5.8 GB)                │   │
│ │ [████████░░░░] 36% - ETA: 8 minutes                    │   │
│ │                                                         │   │
│ │ [⏸️ Pause] [📊 History]                                  │   │
│ └─────────────────────────────────────────────────────────┘   │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Workflow Builder/Runner**
```
┌─────────────────────────────────────────────────────────────┐
│ Workflows                                  [+ New Workflow]  │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ 📋 cell-segmentation                                         │
│    Snakemake pipeline for automated cell segmentation        │
│                                                               │
│    Input folder:   [Browse] NIH-R01-2024/raw/exp_002/       │
│    Output folder:  [Browse] NIH-R01-2024/processed/         │
│                                                               │
│    Compute:  ☑ Use spot instances  (~70% cheaper)           │
│             [32] cores  [64 GB] memory                       │
│                                                               │
│    Estimated: ~45 minutes, ~$2.30                            │
│                                                               │
│    [▶️ Run Workflow]  [💾 Save Config]                       │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│ Recent Runs                                                  │
│ • exp_002 - Running... 12/45 steps ████░░░░ 26%             │
│ • exp_001 - ✓ Completed (42m, $2.15) - 1 hour ago           │
│ • test_run - ✗ Failed (see logs) - 2 days ago               │
└─────────────────────────────────────────────────────────────┘
```

**Workstation Launcher**
```
┌─────────────────────────────────────────────────────────────┐
│ Remote Workstations                      [🖥️ Launch New]     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ 🟢 viz-session (active)                                      │
│    g4dn.xlarge (GPU) • Running 1h 23m                        │
│    Cost so far: $0.68 (spot)                                 │
│    Auto-shutdown: in 37 minutes                              │
│                                                               │
│    [🌐 Open in Browser] [🔌 Connect via VNC] [⏸️ Stop]       │
│    [⏰ Extend 2 Hours] [📸 Snapshot]                          │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│ Stopped Sessions                                             │
│ ⚫ analysis-20241118 (stopped)                               │
│    Cost: $3.45 (6h 42m)                                      │
│    [▶️ Resume] [🗑️ Delete]                                    │
└─────────────────────────────────────────────────────────────┘
```

**Project Management**
```
┌─────────────────────────────────────────────────────────────┐
│ Project: NIH-R01-2024                            [⚙️ Settings]│
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ 📊 Overview                                                  │
│   Storage: 2.3 TB ($45/month)                                │
│   Compute: $12.40 this month                                 │
│   Members: 3                                                 │
│                                                               │
│ 👥 Members                                  [+ Add Member]   │
│   • pi@university.edu (Admin)                                │
│   • jsmith@university.edu (Read-Write)                       │
│   • external@otheruniv.edu (Read-Only)                       │
│                                                               │
│ 📁 Data Organization                                         │
│   /raw/                  1.8 TB   (original data)            │
│   /processed/            450 GB   (analysis outputs)         │
│   /results/              50 GB    (figures, publications)    │
│                                                               │
│ 🔔 Activity Feed                                             │
│   • 2m ago - jsmith: uploaded 45 files to /raw/              │
│   • 1h ago - workflow completed: cell-segmentation           │
│   • 3h ago - external: downloaded dataset_v3                 │
│                                                               │
│ 📈 Storage Trends                                            │
│   [Graph showing growth over time]                           │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Cost Dashboard**
```
┌─────────────────────────────────────────────────────────────┐
│ Cost Management                 Budget: $78 / $100 (78%)    │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ This Month                                                   │
│ ┌───────────────────────────────────────────────────────┐   │
│ │                                                         │   │
│ │  $80 ┤                                              ●  │   │
│ │  $60 ┤                                    ●         ╱   │   │
│ │  $40 ┤                     ●        ╱                   │   │
│ │  $20 ┤         ●     ╱                                   │   │
│ │   $0 ┼──●──────────────────────────────────────────────│   │
│ │      Nov 1        Nov 10       Nov 20       Nov 30      │   │
│ └───────────────────────────────────────────────────────┘   │
│                                                               │
│ Breakdown                                                    │
│   Storage:      $62.30  (80%)  [████████  ]                 │
│   Compute:      $12.40  (16%)  [██        ]                 │
│   Transfer:     $3.15   (4%)   [          ]                 │
│                                                               │
│ 💡 Recommendations                                           │
│   ⚡ Archive data older than 2 years → save $15/month        │
│   💾 Enable compression on logs → save $3/month              │
│   🎯 Use smaller workstation instances → save $5/month       │
│                                                               │
│   [Apply All] [Customize]                                    │
│                                                               │
└─────────────────────────────────────────────────────────────┘
Web UI Technology
Backend API (Go):

go
// internal/webui/server.go
package webui

type Server struct {
    daemon *daemon.Daemon
    router *gin.Engine  // or chi, or stdlib http
}

func (s *Server) Start(port int) error {
    // Serve static files (embedded in binary)
    s.router.StaticFS("/", http.FS(embeddedUI))
    
    // API endpoints
    api := s.router.Group("/api")
    {
        api.GET("/status", s.getStatus)
        api.GET("/projects", s.listProjects)
        api.POST("/sync", s.startSync)
        api.GET("/files/*path", s.listFiles)
        // ... etc
    }
    
    // WebSocket for real-time updates
    s.router.GET("/ws", s.websocketHandler)
    
    return s.router.Run(fmt.Sprintf(":%d", port))
}
Frontend (Simple & Fast):

Vanilla JS + htmx for interactivity (no build step, fast loading)
Tailwind CSS for styling (via CDN, minimal bundle)
Chart.js for visualizations
WebSocket for real-time updates (progress bars, sync status)
Or slightly more sophisticated:

Vue.js or Svelte (lightweight, good for this use case)
Tauri if we want a native desktop app later
Embedded in Binary:

go
//go:embed ui/dist/*
var embeddedUI embed.FS

// Users just run: cicada daemon start --web
// No separate installation needed!
Simple, Intuitive Workflow
Let's walk through what a grad student experiences:

Day 1: Onboarding
PI sends email:

"Hey Jane, we're using Cicada now for lab data. I've added you to the NIH-R01-2024 project.

Setup:

Go to https://cicada.sh/install
Run the installer
Login with your email
That's it! Let me know if you have questions."

Jane's experience:

bash
# Downloads and runs installer (auto-detects OS)
curl -sSL https://cicada.sh/install | sh

# First run
cicada login jane@university.edu
# Opens browser → institutional login (Globus) → done

cicada daemon start --web
# Browser opens to http://localhost:7878
# Shows welcome tour
```

**Jane sees in browser:**
```
Welcome to Smith Lab's Cicada!

You have access to:
  • NIH-R01-2024 (Read-Write)
  • Shared Lab Data (Read-Only)
  • Your Personal Space

Quick Start:
  📤 Upload your data → Drag files to the browser
  💻 Run an analysis → Workflows tab
  🖥️ Need a powerful computer? → Launch a workstation

Your PI's message: "Please upload your microscopy data 
from last week to NIH-R01-2024/raw/jane/"
Day 2: Daily Work
Jane collects microscopy data:

Saves images from microscope to usual folder
Cicada daemon (running in background) notices new files
After 10 minutes (files stable), auto-uploads to S3
Jane gets notification: "Uploaded 45 files from microscope to NIH-R01-2024/raw/jane/experiment_001"
Jane opens web UI:

Sees her files in the project
Clicks "Run Workflow" → selects "cell-segmentation"
Picks input folder, clicks Run
Goes to lunch
Gets email: "Workflow complete! Results in NIH-R01-2024/processed/jane/experiment_001/"
Jane needs to visualize results:

Clicks "Launch Workstation"
Selects "Napari + GPU" template
Clicks Launch
2 minutes later, browser opens with remote desktop
Napari is already open, data is already mounted
Analyzes data, saves figures to S3 automatically
Closes browser tab
Workstation auto-shuts down after 2 hours
Week Later: Sharing with Collaborator
PI wants to share data with external collaborator:

In web UI:

Goes to NIH-R01-2024 project settings
Clicks "Add Member"
Enters: collaborator@otheruniv.edu
Selects: Read-Only access
Clicks Add
Collaborator receives email:

"You've been granted access to NIH-R01-2024 in Smith Lab's Cicada..."
[Setup instructions]

Collaborator can now:

Browse data in web UI
Download specific files
Cannot modify or delete anything
The "Just Works" Factor
The key to success with this persona:

Zero AWS knowledge required
Never mention "S3 bucket" or "IAM role"
It's just "lab data" and "projects"
Familiar concepts
Files and folders (not objects and keys)
Projects and members (not policies and principals)
Upload/download (not PUT/GET)
Automatic everything
Backups happen (versioning + lifecycle)
Costs optimize themselves (intelligent tiering)
Old data archives automatically
Idle resources shut down
Visual feedback
Progress bars for uploads
Real-time sync status
Cost trending graphs
Activity feeds
Helpful, not patronizing
Sensible defaults
Recommendations (with one-click apply)
Learn from usage patterns
Implementation Roadmap
Phase 1: Core + CLI (6 weeks)

Storage sync engine
File watching
Basic user/project management
IAM automation
Phase 2: Web UI (4 weeks)

File browser
Upload/download
Sync manager
Project dashboard
Phase 3: Collaboration (3 weeks)

User invitations
Access control UI
Activity feeds
Sharing links
Phase 4: Compute (4 weeks)

Workflow runner
Workstation launcher
Progress tracking
Phase 5: Polish (ongoing)

Cost optimization suggestions
Usage analytics
Documentation
Video tutorials
The Pitch
To PIs:

"Cicada is a data management platform for academic labs. It costs $50-100/month, requires no IT staff, and ensures your lab never loses data again. It's like Dropbox, but designed for research data and computational workflows."

To Grad Students:

"Cicada automatically backs up your data, runs your analyses faster (on cloud GPUs), and makes it easy to share with collaborators. Install it once, forget about it."

To Funding Agencies:

"Cicada makes research data FAIR (Findable, Accessible, Interoperable, Reusable) by default, with minimal burden on researchers. It's open source and can be deployed at any institution."

This really is a mini data commons - but one that respects the constraints and culture of academic research. Thoughts on any of these pieces?

Think a bit more about metadata and FAIR principals
Also, data protection (NIST 800-171 at a minimum - technical)
Data sharing, web portal/presence options

Excellent points - these are critical for making Cicada actually useful for serious research (and fundable). Let me think through each dimension...

Metadata & FAIR Principles
Academic data is useless without context. Cicada needs to make metadata automatic and effortless, not an afterthought.

Automatic Metadata Extraction
At Upload Time:

go
// Cicada automatically extracts metadata from common formats
type MetadataExtractor interface {
    CanHandle(filename string) bool
    Extract(filepath string) (Metadata, error)
}

// Built-in extractors:
- TIFF/OME-TIFF (microscopy): dimensions, channels, timestamps, instrument settings
- DICOM (medical imaging): patient info, scan parameters
- FASTQ (sequencing): read counts, quality scores
- HDF5/Zarr: embedded attributes
- Generic: file size, checksums, creation time, EXIF data
User Experience:

bash
# Upload with auto-extracted metadata
cicada upload microscopy_data/experiment_001.tif

# Output:
# Uploading experiment_001.tif...
# ✓ Uploaded (2.3 GB)
# 
# Extracted metadata:
#   Instrument: Zeiss LSM 980
#   Magnification: 63x
#   Dimensions: 2048x2048x45 (XYZ)
#   Channels: 4 (DAPI, GFP, RFP, Cy5)
#   Timestamp: 2024-11-18 14:32:15
#   Operator: jsmith
#   Experiment: protein_localization
# 
# Add additional metadata? [y/N] y
# Sample ID: WT-strain-01
# Treatment: control
# Notes: baseline condition
# 
# Metadata saved ✓
```

**Web UI - Metadata Enhancement:**
```
┌─────────────────────────────────────────────────────────────┐
│ File: experiment_001.tif                                     │
├─────────────────────────────────────────────────────────────┤
│ 📋 Metadata                                    [✏️ Edit]      │
│                                                               │
│ Automatically Extracted:                                     │
│   📸 Instrument: Zeiss LSM 980                               │
│   🔬 Magnification: 63x                                      │
│   📐 Dimensions: 2048x2048x45 (XYZ)                          │
│   🎨 Channels: DAPI, GFP, RFP, Cy5                           │
│   📅 Acquired: 2024-11-18 14:32:15                           │
│   👤 Operator: jsmith                                        │
│                                                               │
│ User-Added:                                                  │
│   🧬 Sample ID: WT-strain-01                                 │
│   💊 Treatment: control                                      │
│   📝 Notes: baseline condition for comparison                │
│   🏷️ Tags: control, wild-type, baseline                      │
│                                                               │
│ Provenance:                                                  │
│   📤 Uploaded by: jsmith@university.edu                      │
│   📍 Upload source: /Volumes/Microscope/Export               │
│   🔗 Related files: experiment_002.tif, experiment_003.tif   │
│   📊 Used in workflows: cell-segmentation (run_456)          │
│                                                               │
│ [💾 Save] [📋 Copy as JSON] [📝 Edit Schema]                 │
└─────────────────────────────────────────────────────────────┘
FAIR Compliance Engine
Findable:

bash
# Rich search across all metadata fields
cicada search "protein localization" \
  --instrument "Zeiss LSM 980" \
  --date-after 2024-11-01 \
  --operator jsmith

# Persistent identifiers (DOIs via DataCite)
cicada dataset publish NIH-R01-2024/final_data/ \
  --mint-doi \
  --title "Protein localization in yeast under stress" \
  --authors "Smith J, Garcia A, Chen L"

# Output:
# Published dataset:
#   DOI: 10.12345/cicada.smith-lab.2024.001
#   URL: https://data.cicada.sh/smith-lab/doi/10.12345/...
#   Citation: Smith J, et al. (2024). Protein localization...
Accessible:

bash
# Multiple access methods
- Web portal (with embargos/access controls)
- Direct S3 access (for authenticated users)
- HTTPS downloads (with signed URLs)
- Globus endpoints (for large transfers)
- OAI-PMH harvesting (for repositories)

# Embargo support
cicada dataset publish NIH-R01-2024/final_data/ \
  --embargo-until 2025-06-01 \
  --embargo-type metadata-only  # or 'full'

# During embargo: metadata visible, data access restricted
Interoperable:

bash
# Export metadata in standard formats
cicada metadata export experiment_001.tif \
  --format datacite  # or: dublin-core, schema-org, json-ld

# Outputs DataCite XML:
<?xml version="1.0" encoding="UTF-8"?>
<resource xmlns="http://datacite.org/schema/kernel-4">
  <identifier identifierType="DOI">10.12345/...</identifier>
  <creators>
    <creator><creatorName>Smith, Jane</creatorName></creator>
  </creators>
  <titles>
    <title>Protein localization experiment 001</title>
  </titles>
  <subjects>
    <subject>Cell Biology</subject>
    <subject>Fluorescence Microscopy</subject>
  </subjects>
  ...
</resource>

# Vocabulary mappings (ontologies)
cicada metadata annotate experiment_001.tif \
  --organism "NCBITaxon:4932" \  # Yeast
  --assay "OBI:0000725" \         # Fluorescence microscopy
  --sample-type "CL:0000003"     # Native cell
Reusable:

bash
# Comprehensive provenance tracking
cicada provenance show results/analysis_final.csv

# Output (as graph):
results/analysis_final.csv
  ← cell-segmentation workflow (v2.1.0, 2024-11-20)
      ← processed/experiment_001_segmented.tif
          ← processing script normalize.py (v1.3, 2024-11-19)
              ← raw/experiment_001.tif
                  ← Zeiss LSM 980 (SN: 12345)
                  ← Operator: jsmith
                  ← Protocol: standard_imaging_v3.pdf

# License specification
cicada dataset publish --license CC-BY-4.0
cicada dataset publish --license CC0  # Public domain

# Computational environment capture
cicada workflow run --capture-environment
# Saves: Docker image, conda env, pip freeze, system info
Metadata Schema Management
Lab-Specific Schemas:

bash
# Define custom metadata schema for your lab
cicada schema create microscopy-experiment

# Interactive schema builder in web UI:
┌─────────────────────────────────────────────────────────────┐
│ Create Metadata Schema: microscopy-experiment               │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ Required Fields:                                             │
│   [+] Sample ID         (text, unique)                       │
│   [+] Experiment Date   (date)                               │
│   [+] Operator          (user, auto-filled)                  │
│   [+] Magnification     (number, units: x)                   │
│                                                               │
│ Optional Fields:                                             │
│   [+] Treatment         (text, controlled vocabulary)        │
│       Values: control, drug-A, drug-B, heat-shock            │
│   [+] Replicate Number  (integer, 1-10)                      │
│   [+] Notes             (long text)                          │
│                                                               │
│ Computed Fields:                                             │
│   [+] Upload Date       (auto, datetime)                     │
│   [+] File Checksum     (auto, SHA256)                       │
│                                                               │
│ [💾 Save Schema] [📋 Import from Template]                   │
└─────────────────────────────────────────────────────────────┘

# Apply schema to folder
cicada schema apply microscopy-experiment \
  --path NIH-R01-2024/raw/microscopy/

# Now all uploads to that path require schema compliance
cicada upload image.tif NIH-R01-2024/raw/microscopy/
# Prompts for required fields if not auto-extracted
Template Library:

bash
cicada schema list-templates

# Output:
# Available Templates:
#   - microscopy-basic (fluorescence imaging)
#   - microscopy-advanced (includes instrument settings)
#   - sequencing-rnaseq (RNA-seq experiments)
#   - sequencing-chipseq (ChIP-seq experiments)
#   - proteomics-ms (mass spectrometry)
#   - hplc (chromatography)
#   - generic-experiment (minimal fields)
# 
# Community templates: 52 available
# Search: cicada schema search [keyword]
Data Protection (NIST 800-171)
For labs handling sensitive data (patient info, CUI, export-controlled), Cicada needs to be compliant out-of-the-box.

NIST 800-171 Requirements Mapping
Access Control (AC):

yaml
# Cicada enforces automatically:
- AC.1.001: Limit access to authorized users ✓
  → IAM policies, per-project/per-file access
  
- AC.1.002: Limit access to authorized transactions ✓
  → Audit logs of all data operations
  
- AC.1.003: Control public information ✓
  → All data private by default, explicit sharing
  
- AC.2.005: Separate duties ✓
  → Role-based access (admin, read-write, read-only)
  
- AC.2.007: Employ least privilege ✓
  → Principle enforced in IAM policy generation
  
- AC.2.013: Monitor/control remote access ✓
  → All access via authenticated API/web UI
Audit & Accountability (AU):

bash
# Complete audit trail
cicada audit log --user jsmith --last 30d

# Output:
# Audit Log: jsmith@university.edu (Last 30 days)
# 
# 2024-11-20 14:32:15 - ACCESS - Read file: patient_data/subject_001.dcm
#   Source IP: 192.168.1.45
#   User Agent: Cicada CLI v1.2.0
#   Auth Method: Globus + MFA
# 
# 2024-11-20 14:35:22 - MODIFY - Upload file: patient_data/subject_002.dcm
#   Source IP: 192.168.1.45
#   File Hash: sha256:abc123...
#   Encrypted: AES-256
# 
# 2024-11-19 09:15:33 - SHARE - Granted access to external@otheruniv.edu
#   Project: Clinical-Study-2024
#   Permission: Read-Only
#   Approved by: pi@university.edu

# Export audit logs for compliance
cicada audit export --format csv --year 2024
# Saved to: audit_2024.csv (for compliance reporting)
Identification & Authentication (IA):

bash
# Multi-factor authentication
cicada login jsmith@university.edu
# Username: jsmith@university.edu
# Password: ********
# 2FA Code (from authenticator app): 123456
# ✓ Authenticated

# Or via institutional SSO + MFA
cicada login --provider globus
# Redirects to university SSO (which has MFA)

# Session management
cicada session list
# Active Sessions:
#   - laptop (192.168.1.45) - expires in 7 days
#   - workstation (10.0.1.23) - expires in 1 hour
# 
# Revoke session:
cicada session revoke laptop
Media Protection (MP):

bash
# Encryption at rest (all S3 data)
- Server-side encryption (SSE-S3 or SSE-KMS)
- Optional: Customer-managed keys (CMK)

# Encryption in transit (all transfers)
- TLS 1.3 for all API calls
- Encrypted multipart uploads
- Verified checksums

# Secure deletion
cicada delete patient_data/subject_001.dcm \
  --secure-delete \
  --overwrite-passes 3

# S3 object versions + metadata also deleted
# Audit log records deletion + verifies unrecoverability
System & Communications Protection (SC):

bash
# Network segmentation
cicada lab configure --network-mode private

# Creates:
- VPC with private subnets
- VPC endpoints for S3 (no internet traversal)
- Security groups (least-privilege)
- Network ACLs

# Data classification
cicada classify patient_data/ --level CUI

# Automatically applies:
- Stricter access controls
- Mandatory encryption
- Enhanced audit logging
- DLP (data loss prevention) scanning
System & Information Integrity (SI):

bash
# Malware scanning on upload
cicada upload data.zip --scan

# Output:
# Scanning data.zip for malware...
# ✓ Clean (scanned with ClamAV)
# Uploading...

# Integrity verification
cicada verify s3://lab-bucket/data/experiment.tif

# Output:
# Verifying integrity...
# Local SHA256:  abc123...
# Remote SHA256: abc123...
# S3 ETag:       def456...
# ✓ File integrity verified

# Flaw remediation (automatic updates)
cicada update --check
# New version available: v1.3.0 (security fixes)
# Update now? [Y/n]
Compliance Mode
Enable NIST 800-171 Mode:

bash
cicada lab configure --compliance nist-800-171

# This enforces:
✓ Multi-factor authentication required
✓ Session timeouts (1 hour default)
✓ All data encrypted (in transit + at rest)
✓ Customer-managed encryption keys (KMS)
✓ Enhanced audit logging (immutable)
✓ Regular access reviews (monthly)
✓ Automated vulnerability scanning
✓ Network isolation (VPC with endpoints)
✓ Principle of least privilege
✓ Incident response procedures

# Generate compliance report
cicada compliance report --standard nist-800-171

# Output: PDF report mapping Cicada controls to requirements
Additional Compliance Profiles:

bash
# HIPAA (for clinical research)
cicada lab configure --compliance hipaa

# GDPR (for EU data subjects)
cicada lab configure --compliance gdpr
# Enables: data subject rights, consent tracking, geographic restrictions

# ITAR/EAR (export control)
cicada lab configure --compliance itar
# Enables: US-only storage, person authentication, access restrictions
Data Loss Prevention (DLP)
Sensitive Data Detection:

bash
# Scan for PII/PHI automatically
cicada scan patient_data/ --detect-sensitive

# Output:
# Scanning 1,247 files...
# 
# ⚠️  Found sensitive data:
#   - 15 files contain SSN patterns
#   - 8 files contain credit card numbers
#   - 142 files contain dates of birth
#   - 23 files contain IP addresses
# 
# Recommendations:
#   1. Apply stricter access controls
#   2. Enable field-level encryption
#   3. Redact sensitive fields before sharing
# 
# Apply recommendations? [Y/n]

# Automatic redaction for sharing
cicada share patient_data/study_results.csv \
  --redact SSN,DOB,PatientName \
  --recipient external@otheruniv.edu
Data Sharing & Web Portal
This is where Cicada becomes a research data repository for the lab.

Public Data Portal
Enable Public Portal:

bash
cicada portal enable --domain smithlab.cicada.sh

# Output:
# Setting up public data portal...
# ✓ SSL certificate provisioned
# ✓ CloudFront distribution created
# ✓ Custom domain configured
# 
# Your portal is live at:
#   https://smithlab.cicada.sh
# 
# Customize appearance:
#   cicada portal customize --logo lab_logo.png
```

**Portal Features:**
```
┌─────────────────────────────────────────────────────────────┐
│ Smith Lab Data Repository                    🔍 [Search]     │
│ https://smithlab.cicada.sh                                   │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ 🏠 Welcome                                                   │
│                                                               │
│ The Smith Lab studies protein localization dynamics in       │
│ yeast. We share our data openly to advance scientific        │
│ discovery.                                                   │
│                                                               │
│ 📊 Browse Datasets (12)                                      │
│                                                               │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ 📁 Protein Localization Under Heat Shock                 │ │
│ │    DOI: 10.12345/smith.2024.001                          │ │
│ │    Published: 2024-11-15 | Size: 2.3 TB                  │ │
│ │    License: CC-BY-4.0                                     │ │
│ │    Downloads: 47 | Citations: 3                          │ │
│ │    [📥 Download] [📋 Cite] [ℹ️ Details]                   │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                               │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ 📁 Time-Lapse Microscopy: Cell Division                  │ │
│ │    DOI: 10.12345/smith.2024.002                          │ │
│ │    Published: 2024-10-22 | Size: 856 GB                  │ │
│ │    License: CC0 (Public Domain)                          │ │
│ │    Downloads: 123 | Citations: 8                         │ │
│ │    [📥 Download] [📋 Cite] [ℹ️ Details]                   │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                               │
│ 📑 Publications                                              │
│   - Smith et al. (2024) Nature Cell Biology                 │
│   - Garcia & Smith (2024) Cell Reports                      │
│                                                               │
│ 👥 Lab Members | 📧 Contact | 🔬 Protocols                   │
└─────────────────────────────────────────────────────────────┘
```

**Dataset Landing Page:**
```
┌─────────────────────────────────────────────────────────────┐
│ Protein Localization Under Heat Shock                       │
│ DOI: 10.12345/smith.2024.001                                │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ 📊 Overview                                                  │
│   Description: Time-lapse fluorescence microscopy of        │
│                protein localization dynamics in S.           │
│                cerevisiae during heat shock stress           │
│   Authors: Smith J, Garcia A, Chen L, Williams R            │
│   Published: 2024-11-15                                      │
│   License: CC-BY-4.0                                         │
│   Size: 2.3 TB (1,247 files)                                │
│                                                               │
│ 📋 Citation                                                  │
│   Smith J, Garcia A, Chen L, Williams R (2024). Protein     │
│   localization under heat shock. Smith Lab Data Repository. │
│   https://doi.org/10.12345/smith.2024.001                   │
│   [📋 Copy] [📄 BibTeX] [📄 RIS]                             │
│                                                               │
│ 📥 Download Options                                          │
│   [⬇️ Download All (2.3 TB)]                                 │
│   [📂 Browse Files]                                          │
│   [🌐 Globus Transfer] (faster for large datasets)           │
│   [💻 AWS S3 Access] (for computational analysis)            │
│                                                               │
│ 🔬 Methods                                                   │
│   Organism: Saccharomyces cerevisiae (NCBITaxon:4932)       │
│   Assay: Fluorescence microscopy (OBI:0000725)              │
│   Instrument: Zeiss LSM 980                                  │
│   Magnification: 63x                                         │
│   Time points: 100 (every 30 seconds)                       │
│   Replicates: 3 biological, 3 technical                     │
│                                                               │
│ 📊 Files (sample)                                            │
│   raw/                                  1.8 TB               │
│     experiment_001.czi                  128 MB               │
│     experiment_002.czi                  131 MB               │
│     ...                                                      │
│   processed/                            450 GB               │
│     segmented/                                               │
│     tracked/                                                 │
│   analysis/                             50 GB                │
│     figures/                                                 │
│     statistics/                                              │
│   README.md                             15 KB                │
│   protocols/imaging_protocol.pdf        2.1 MB               │
│                                                               │
│ 🔗 Related                                                   │
│   Publication: Smith et al. (2024) Nature Cell Biology      │
│   DOI: 10.1038/s41556-024-xxxxx                             │
│   Code: github.com/smith-lab/heat-shock-analysis            │
│   Protocol: protocols.io/view/heat-shock-imaging            │
│                                                               │
│ 📈 Usage Statistics                                          │
│   Downloads: 47                                              │
│   Citations: 3                                               │
│   Views: 892                                                 │
│                                                               │
└─────────────────────────────────────────────────────────────┘
Flexible Sharing Models
1. Public (Open Data):

bash
cicada dataset publish NIH-R01-2024/final_data/ \
  --visibility public \
  --license CC-BY-4.0

# Anyone can download, must cite
2. Embargo (Delayed Release):

bash
cicada dataset publish NIH-R01-2024/final_data/ \
  --visibility public \
  --embargo-until 2025-06-01 \
  --embargo-type metadata-only

# Metadata visible now, data downloads on 2025-06-01
# Journal accepts this for data availability statements
3. Controlled Access (Request-Based):

bash
cicada dataset publish sensitive_data/ \
  --visibility controlled \
  --approval-required \
  --approvers pi@university.edu

# Portal shows:
# "Access to this dataset requires approval.
#  Reason for access: [text field]
#  [Request Access]"

# PI receives email, approves/denies
cicada access-requests list
cicada access-requests approve req-123 --expires 1y
4. Private Link (Reviewers/Collaborators):

bash
cicada share NIH-R01-2024/final_data/ \
  --create-link \
  --expires 90d \
  --password optional-password

# Output:
# Shareable link created:
#   https://smithlab.cicada.sh/s/a8f3k29s
#   Password: dolphin-keyboard-87 (optional)
#   Expires: 2025-02-20
# 
# This link can be shared with reviewers/collaborators.
# No login required.
5. Federated (Cross-Institution):

bash
# Share with other Cicada instances
cicada share NIH-R01-2024/final_data/ \
  --federate johnslab.cicada.sh \
  --permission read-only

# Their users can access: cicada mount smithlab.cicada.sh/NIH-R01-2024/
Data Discovery & Indexing
Make Lab Data Discoverable:

bash
# Submit to public data registries
cicada portal index --registry re3data
cicada portal index --registry datacite
cicada portal index --registry google-dataset-search

# Portal generates schema.org JSON-LD automatically:
<script type="application/ld+json">
{
  "@context": "https://schema.org/",
  "@type": "Dataset",
  "name": "Protein Localization Under Heat Shock",
  "description": "Time-lapse fluorescence microscopy...",
  "identifier": "https://doi.org/10.12345/smith.2024.001",
  "creator": [
    {"@type": "Person", "name": "Jane Smith"},
    {"@type": "Person", "name": "Alex Garcia"}
  ],
  "license": "https://creativecommons.org/licenses/by/4.0/",
  "distribution": {
    "@type": "DataDownload",
    "contentUrl": "https://smithlab.cicada.sh/datasets/001",
    "encodingFormat": "application/zip"
  },
  ...
}
</script>

# Now discoverable via Google Dataset Search!
Portal Customization
bash
# Branding
cicada portal customize \
  --logo lab_logo.png \
  --color-primary "#1a5490" \
  --title "Smith Lab Data Repository" \
  --description "Open data from protein localization research"

# Add pages
cicada portal page create about \
  --title "About the Lab" \
  --content about.md

cicada portal page create protocols \
  --title "Experimental Protocols" \
  --content protocols.md

# Navigation
cicada portal menu add "Publications" https://lab.university.edu/pubs
cicada portal menu add "GitHub" https://github.com/smith-lab

# Analytics
cicada portal analytics --provider plausible
# Track downloads, views, citations (privacy-friendly)
API for Data Access
Public API for programmatic access:

bash
# REST API (auto-generated)
GET https://smithlab.cicada.sh/api/datasets
GET https://smithlab.cicada.sh/api/datasets/001
GET https://smithlab.cicada.sh/api/datasets/001/files

# Download via API
curl -O https://smithlab.cicada.sh/api/datasets/001/download

# OAI-PMH (for harvesting by repositories)
https://smithlab.cicada.sh/oai?verb=ListRecords&metadataPrefix=oai_dc
Integration with External Tools:

python
# Python client library
from cicada import Dataset

# Browse public datasets
datasets = Dataset.list_public(lab="smithlab")

# Download specific dataset
ds = Dataset.from_doi("10.12345/smith.2024.001")
ds.download(path="./local_data/")

# Load data directly
import pandas as pd
df = pd.read_csv(ds.get_file("analysis/results.csv"))
Putting It All Together
The Complete Workflow:

Researcher collects data
Automatic upload + metadata extraction
Schema validation ensures completeness
Provenance tracking starts
Analysis
Workflows capture computational environment
Input/output relationships recorded
Reproducible by design
Collaboration
Secure sharing with collaborators
Access controls + audit logs
NIST 800-171 compliant if needed
Publication
One command to publish dataset
DOI minted, portal page created
Indexed in Google Dataset Search
Data availability statement satisfied
Long-term Preservation
Intelligent tiering (cost-optimized)
Versioning (never lose data)
Checksum verification (integrity)
Glacier for archival
The User Experience:

bash
# Day 1: Setup
cicada lab init --name smith-lab
cicada portal enable --domain smithlab.cicada.sh

# Day 2-365: Daily work
# (automatic sync, metadata, provenance)

# Day 366: Publication time
cicada dataset prepare NIH-R01-2024/final_data/ \
  --title "Protein localization under heat shock" \
  --authors "Smith J, Garcia A, Chen L" \
  --description description.md \
  --keywords "protein localization, heat shock, yeast" \
  --license CC-BY-4.0

cicada dataset validate NIH-R01-2024/final_data/
# ✓ All files have metadata
# ✓ README present
# ✓ Protocols documented
# ✓ Code available (GitHub linked)
# ✓ FAIR principles satisfied

cicada dataset publish NIH-R01-2024/final_data/ \
  --mint-doi \
  --notify-coauthors

# Output:
# Dataset published! 
# DOI: 10.12345/smith.2024.001
# URL: https://smithlab.cicada.sh/datasets/001
# Citation: Smith J, et al. (2024)...
# 
# Data availability statement (copy for manuscript):
# "Data are available at https://doi.org/10.12345/smith.2024.001
#  under a CC-BY-4.0 license."
This makes Cicada not just a storage system, but a complete research data lifecycle platform - from collection to publication to preservation. And critically, it does this while being:

Affordable ($50-100/month)
Easy to use (web UI + CLI)
Compliant (NIST 800-171, HIPAA, etc.)
FAIR by default
Thoughts on this expanded vision?

DOI generation?
The metadata could be even more flexible - the group I was speaking to does, for example, a lot of imaging and sequencing. So, basically support for domain specific metadata they may label themselves.

DOI Generation
You're right - DOI minting is critical for making datasets citable and fundable. Let me think through the implementation...

DataCite Integration
How It Works:

bash
# One-time setup (PI does this)
cicada lab configure-doi --provider datacite

# Interactive setup:
# 
# DataCite DOI Configuration
# 
# Option 1: Use Cicada's DataCite Account (Easiest)
#   Cost: $1 per DOI minted
#   Prefix: 10.12345 (Cicada's prefix)
#   Example DOI: 10.12345/cicada.smith-lab.2024.001
#   [Select Option 1]
# 
# Option 2: Use Your Institution's DataCite Account
#   You'll need: Repository ID, Password
#   Prefix: Your institution's prefix (e.g., 10.5061)
#   Example DOI: 10.5061/smith-lab.2024.001
#   [Select Option 2]
# 
# Option 3: Use Zenodo (Free, but less control)
#   Free DOIs via Zenodo integration
#   Metadata synced to Zenodo
#   [Select Option 3]

# Choose option 1 (Cicada managed)
# Selected: Cicada Managed DOIs
# ✓ Configuration saved
Minting a DOI:

bash
cicada dataset publish NIH-R01-2024/final_data/ \
  --mint-doi \
  --title "Protein localization dynamics under heat shock" \
  --authors "Smith J, Garcia A, Chen L" \
  --description "High-resolution time-lapse microscopy..." \
  --license CC-BY-4.0

# Behind the scenes:
# 1. Validates dataset completeness
# 2. Generates DataCite metadata XML
# 3. Calls DataCite API to reserve DOI
# 4. Creates landing page
# 5. Registers DOI with DataCite
# 6. Updates portal

# Output:
# Minting DOI...
# ✓ DOI registered: 10.12345/cicada.smith-lab.2024.001
# 
# Landing page: https://smithlab.cicada.sh/doi/10.12345/cicada.smith-lab.2024.001
# Citation: Smith J, Garcia A, Chen L (2024). Protein localization 
#   dynamics under heat shock. Smith Lab Data Repository. 
#   https://doi.org/10.12345/cicada.smith-lab.2024.001
# 
# Cost: $1.00 (charged to lab account)
DOI Management:

bash
# List all DOIs
cicada doi list

# Output:
# Lab DOIs (3):
#   10.12345/cicada.smith-lab.2024.001
#     Title: Protein localization under heat shock
#     Status: Public
#     Downloads: 47
#     Citations: 3 (detected via Crossref)
# 
#   10.12345/cicada.smith-lab.2024.002
#     Title: Cell division time-lapse microscopy
#     Status: Embargo (until 2025-06-01)
#     Downloads: 0
# 
#   10.12345/cicada.smith-lab.2023.005
#     Title: RNA-seq data: stress response
#     Status: Public
#     Downloads: 123
#     Citations: 8

# Update DOI metadata (if needed)
cicada doi update 10.12345/cicada.smith-lab.2024.001 \
  --add-author "Williams R" \
  --version 1.1

# Versions create new DOI with relationship:
# 10.12345/cicada.smith-lab.2024.001.v1 (original)
# 10.12345/cicada.smith-lab.2024.001.v2 (updated)
# Concept DOI: 10.12345/cicada.smith-lab.2024.001 (always latest)
Cost Structure:

bash
# Cicada managed DOIs
- $1 per DOI minted
- Included in lab bill
- No annual fees
- No minimum commitment

# Institution managed DOIs
- Uses your institution's account
- You pay institution's rates
- Cicada handles technical integration

# Zenodo integration
- Free
- But data lives on Zenodo, not your S3
- Good for small datasets (<50GB)
DataCite Metadata Generation:

go
// internal/doi/datacite.go
func (d *DOIManager) GenerateDataCiteXML(dataset *Dataset) (string, error) {
    // Required fields
    metadata := DataCiteMetadata{
        Identifier: dataset.DOI,
        Creators:   dataset.Authors,
        Titles:     []string{dataset.Title},
        Publisher:  dataset.Lab.Name,
        PublicationYear: time.Now().Year(),
        ResourceType: "Dataset",
        
        // Optional but recommended
        Subjects:     dataset.Keywords,
        Contributors: dataset.Contributors,
        Dates:        dataset.Dates,
        Language:     "en",
        AlternateIdentifiers: dataset.AlternateIDs,
        RelatedIdentifiers: dataset.RelatedDOIs,
        Sizes:        []string{dataset.Size},
        Formats:      dataset.FileFormats,
        Version:      dataset.Version,
        RightsList:   []string{dataset.License},
        Descriptions: []Description{
            {Type: "Abstract", Value: dataset.Description},
            {Type: "Methods", Value: dataset.Methods},
        },
        GeoLocations: dataset.GeoLocations,
        FundingReferences: dataset.Grants,
    }
    
    return metadata.ToXML()
}
```

## Flexible Domain-Specific Metadata

This is where Cicada becomes really powerful - letting researchers define their own metadata schemas while maintaining interoperability.

### Metadata Architecture
```
┌─────────────────────────────────────────┐
│         Core Metadata (Required)        │
│  - Title, Authors, Date, License        │
│  - Basic Dublin Core/DataCite fields    │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│    Domain Extensions (Optional)         │
│  - microscopy                           │
│  - sequencing                           │
│  - proteomics                           │
│  - custom (user-defined)                │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│   Lab-Specific Fields (Customizable)    │
│  - Your experimental variables          │
│  - Your controlled vocabularies         │
│  - Your workflow parameters             │
└─────────────────────────────────────────┘
Creating Custom Metadata Schemas
Scenario: Imaging + Sequencing Lab

bash
# Start with a template
cicada metadata schema create microscopy-experiment \
  --template fluorescence-microscopy

# Opens an editor (or web UI):
Schema Definition (YAML):

yaml
# microscopy-experiment.yaml
schema_version: "1.0"
name: microscopy-experiment
description: Fluorescence microscopy experimental metadata
domain: imaging

# Inherit from standard schemas
extends:
  - cicada://core/experiment  # Basic experiment fields
  - cicada://domains/microscopy/fluorescence  # Standard microscopy fields

# Core fields (always required)
required_fields:
  - sample_id
  - experiment_date
  - operator
  - protocol_id

# Domain-specific fields (microscopy)
microscopy:
  instrument:
    type: string
    required: true
    vocabulary: [Zeiss LSM 980, Nikon A1R, Leica SP8, Other]
    
  magnification:
    type: number
    required: true
    units: x
    range: [1, 100]
    
  objective:
    type: object
    fields:
      na: 
        type: number
        description: Numerical aperture
      immersion:
        type: string
        vocabulary: [air, oil, water, glycerol]
  
  channels:
    type: array
    items:
      type: object
      fields:
        name: string
        wavelength_ex: {type: number, units: nm}
        wavelength_em: {type: number, units: nm}
        fluorophore: string
        exposure_ms: number
        laser_power_percent: number
  
  dimensions:
    type: object
    required: true
    fields:
      x_pixels: integer
      y_pixels: integer
      z_slices: integer
      time_points: integer
      pixel_size_um: number
  
  acquisition_mode:
    type: string
    vocabulary: [single-plane, z-stack, time-lapse, z-stack-time-lapse]

# Lab-specific fields (your custom additions)
experimental:
  sample:
    strain:
      type: string
      description: Yeast strain identifier
      pattern: "^[A-Z]{2,4}[0-9]{3,4}$"
      examples: [BY4741, W303, YPH499]
    
    genotype:
      type: string
      description: Full genotype string
    
    plasmid:
      type: array
      items: string
      description: Plasmid IDs
  
  treatment:
    condition:
      type: string
      required: true
      vocabulary: [control, heat-shock, drug-A, drug-B, starvation]
    
    duration_min:
      type: number
      required_if: condition != "control"
      units: minutes
    
    concentration:
      type: number
      required_if: condition in [drug-A, drug-B]
      units: uM
  
  replication:
    biological_replicate:
      type: integer
      range: [1, 10]
    
    technical_replicate:
      type: integer
      range: [1, 5]

# Computed/automatic fields
computed:
  file_checksum:
    type: string
    auto: sha256
  
  upload_date:
    type: datetime
    auto: now
  
  file_size_bytes:
    type: integer
    auto: filesize

# Ontology mappings (for interoperability)
ontology_mappings:
  organism: NCBITaxon:4932  # S. cerevisiae
  assay_type: OBI:0000725   # fluorescence microscopy assay
  sample_type: CL:0000003   # native cell

# Validation rules
validation:
  - rule: if z_slices > 1 then acquisition_mode in [z-stack, z-stack-time-lapse]
    message: "Z-slices specified but acquisition mode doesn't match"
  
  - rule: if treatment.condition != "control" then treatment.duration_min required
    message: "Treatment duration required for non-control conditions"
  
  - rule: channels.length >= 1
    message: "At least one channel required"

# Search/faceting (for web portal)
facets:
  - field: microscopy.instrument
    label: "Instrument"
  - field: experimental.treatment.condition
    label: "Treatment"
  - field: experimental.sample.strain
    label: "Strain"
  - field: operator
    label: "Operator"
Similar schema for sequencing:

yaml
# sequencing-experiment.yaml
schema_version: "1.0"
name: sequencing-experiment
description: NGS experiment metadata
domain: sequencing

extends:
  - cicada://core/experiment
  - cicada://domains/sequencing/rnaseq

sequencing:
  platform:
    type: string
    required: true
    vocabulary: [Illumina, PacBio, ONT, BGI]
  
  instrument:
    type: string
    vocabulary: [NovaSeq 6000, NextSeq 2000, HiSeq 4000, ...]
  
  library_prep:
    protocol:
      type: string
      required: true
    
    kit:
      type: string
    
    insert_size_bp:
      type: integer
    
    selection:
      type: string
      vocabulary: [PCR, PolyA, rRNA-depletion, size-selection, ...]
  
  run_parameters:
    read_length:
      type: integer
      units: bp
    
    read_type:
      type: string
      vocabulary: [single-end, paired-end]
    
    flow_cell_id:
      type: string
    
    lane:
      type: integer
      range: [1, 8]

experimental:
  sample:
    tissue:
      type: string
      ontology: UBERON
    
    cell_type:
      type: string
      ontology: CL
    
    developmental_stage:
      type: string
  
  treatment:
    # Same as microscopy schema
  
  extraction:
    method:
      type: string
    
    rna_integrity_number:
      type: number
      range: [1, 10]
    
    concentration_ng_ul:
      type: number

quality_control:
  total_reads:
    type: integer
    auto: fastq_count
  
  percent_q30:
    type: number
    auto: fastq_quality
  
  percent_aligned:
    type: number
Using Custom Schemas
Apply Schema to Data:

bash
# Apply schema to folder
cicada metadata schema apply microscopy-experiment \
  --path NIH-R01-2024/raw/microscopy/

# Now uploads require this schema
cicada upload experiment_001.czi NIH-R01-2024/raw/microscopy/

# Cicada extracts what it can from file, prompts for rest:
```

**Interactive Metadata Entry (CLI):**
```
Uploading: experiment_001.czi (2.3 GB)

Auto-extracted metadata:
  ✓ Instrument: Zeiss LSM 980
  ✓ Magnification: 63x
  ✓ Dimensions: 2048x2048x45x100
  ✓ Channels: 4 (DAPI, GFP, RFP, Cy5)
  ✓ Acquisition mode: z-stack-time-lapse

Please provide additional metadata:

Sample Information:
  Sample ID: [WT-001]
  Strain: [BY4741]
  Genotype: [MATa his3Δ1 leu2Δ0 met15Δ0 ura3Δ0]
  Plasmid (optional): [pRS413-GFP-NLS]

Treatment:
  Condition: [1) control  2) heat-shock  3) drug-A  4) drug-B]
  Choose: [2]
  Duration (minutes): [30]

Replication:
  Biological replicate: [1]
  Technical replicate: [1]

Operator: [jsmith] (auto-filled)
Protocol ID: [protocols.io/view/heat-shock-imaging-v3]

✓ Metadata validation passed
✓ Uploading to S3...
```

**Web UI Metadata Form:**
```
┌─────────────────────────────────────────────────────────────┐
│ Upload File: experiment_001.czi                              │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ 📋 Metadata Form (microscopy-experiment)                     │
│                                                               │
│ ━━━ Automatically Extracted ━━━                              │
│ ✓ Instrument: Zeiss LSM 980                                  │
│ ✓ Magnification: 63x                                         │
│ ✓ Objective NA: 1.4                                          │
│ ✓ Dimensions: 2048×2048×45×100 (X×Y×Z×T)                     │
│ ✓ Pixel size: 0.13 µm                                        │
│ ✓ Channels:                                                  │
│   • DAPI (ex: 405nm, em: 450nm)                              │
│   • GFP (ex: 488nm, em: 525nm)                               │
│   • RFP (ex: 561nm, em: 610nm)                               │
│   • Cy5 (ex: 633nm, em: 670nm)                               │
│                                                               │
│ ━━━ Sample Information ━━━                                   │
│ Sample ID: [WT-001        ] ⓘ Unique identifier              │
│ Strain:    [BY4741        ] ⓘ From controlled vocabulary     │
│ Genotype:  [MATa his3Δ1...] ⓘ Full genotype                  │
│ Plasmid:   [pRS413-GFP-NLS] ⓘ Optional, can add multiple     │
│            [+ Add another]                                    │
│                                                               │
│ ━━━ Treatment ━━━                                            │
│ Condition: [heat-shock ▼] ⓘ Select from list                 │
│ Duration:  [30] minutes   ⓘ Required for treatments          │
│ Temperature: [42] °C      ⓘ Condition-specific field         │
│                                                               │
│ ━━━ Replication ━━━                                          │
│ Biological replicate: [1]                                    │
│ Technical replicate:  [1]                                    │
│                                                               │
│ ━━━ Experiment Details ━━━                                   │
│ Operator:   [jsmith ▼]    ⓘ Auto-filled from login           │
│ Protocol:   [Browse protocols.io] or [Upload PDF]            │
│ Notes:      [                                        ]        │
│             [                                        ]        │
│                                                               │
│ [💾 Save Metadata] [📋 Save as Template] [❌ Cancel]         │
└─────────────────────────────────────────────────────────────┘
Metadata Templates (Save Time)
Create Reusable Templates:

bash
# Save metadata as template for similar experiments
cicada metadata template save \
  --from experiment_001.czi \
  --name heat-shock-standard \
  --fields strain,genotype,plasmid,treatment.condition,protocol_id

# Use template for next experiment
cicada upload experiment_002.czi \
  --template heat-shock-standard \
  --override sample_id=WT-002 \
  --override replication.biological_replicate=2

# Only prompts for changed fields!
Batch Upload with CSV:

bash
# Create CSV with metadata for bulk upload:
# samples.csv:
# filename,sample_id,strain,treatment,biological_rep,technical_rep
# exp_001.czi,WT-001,BY4741,control,1,1
# exp_002.czi,WT-002,BY4741,control,2,1
# exp_003.czi,WT-003,BY4741,control,3,1
# exp_004.czi,HS-001,BY4741,heat-shock,1,1
# exp_005.czi,HS-002,BY4741,heat-shock,2,1

cicada upload --batch samples.csv \
  --files-dir /path/to/experiments/

# Uploads all files with metadata from CSV
# Validates against schema
# Shows progress bar
Schema Management
Share Schemas Across Labs:

bash
# Publish your schema to community
cicada metadata schema publish microscopy-experiment \
  --visibility public \
  --description "Standard microscopy metadata for yeast imaging"

# Others can discover and use it
cicada metadata schema search "fluorescence microscopy yeast"

# Output:
# Community Schemas (3 results):
#   1. microscopy-experiment (smith-lab)
#      ★★★★☆ 23 labs using this
#      "Standard microscopy metadata for yeast imaging"
#   
#   2. super-resolution-microscopy (chen-lab)
#      ★★★☆☆ 8 labs using this
#      "Extended fields for PALM/STORM imaging"
#   
#   3. live-cell-imaging (williams-lab)
#      ★★★★★ 45 labs using this
#      "Comprehensive time-lapse microscopy metadata"

# Use someone else's schema
cicada metadata schema install live-cell-imaging \
  --from williams-lab
Version Control for Schemas:

bash
# Schemas are versioned
cicada metadata schema version microscopy-experiment

# Output:
# Schema: microscopy-experiment
#   v1.0 (current) - 2024-11-01
#     - Initial release
#   v0.9 - 2024-10-15
#     - Beta version
#   v0.8 - 2024-10-01
#     - Alpha version

# Pin to specific version
cicada metadata schema apply microscopy-experiment@1.0

# Upgrade to newer version (with migration)
cicada metadata schema upgrade microscopy-experiment \
  --to 1.1 \
  --migrate existing_data/
Metadata Validation & Quality Control
Automatic Validation:

go
// internal/metadata/validator.go
type Validator struct {
    schema *Schema
}

func (v *Validator) Validate(metadata Metadata) ValidationResult {
    errors := []ValidationError{}
    warnings := []ValidationWarning{}
    
    // Required fields
    for _, field := range v.schema.RequiredFields {
        if !metadata.Has(field) {
            errors = append(errors, MissingFieldError{Field: field})
        }
    }
    
    // Type checking
    for field, value := range metadata.Fields {
        expectedType := v.schema.GetFieldType(field)
        if !value.MatchesType(expectedType) {
            errors = append(errors, TypeMismatchError{
                Field: field,
                Expected: expectedType,
                Got: value.Type(),
            })
        }
    }
    
    // Vocabulary constraints
    for field, value := range metadata.Fields {
        if vocab := v.schema.GetVocabulary(field); vocab != nil {
            if !vocab.Contains(value) {
                warnings = append(warnings, UnknownValueWarning{
                    Field: field,
                    Value: value,
                    Suggestions: vocab.FindSimilar(value),
                })
            }
        }
    }
    
    // Custom validation rules
    for _, rule := range v.schema.ValidationRules {
        if !rule.Evaluate(metadata) {
            errors = append(errors, RuleViolationError{
                Rule: rule.Name,
                Message: rule.Message,
            })
        }
    }
    
    // Ontology validation
    for field, value := range metadata.Fields {
        if ontology := v.schema.GetOntology(field); ontology != nil {
            if !ontology.IsValidTerm(value) {
                warnings = append(warnings, InvalidOntologyTermWarning{
                    Field: field,
                    Value: value,
                    Ontology: ontology.Name,
                })
            }
        }
    }
    
    return ValidationResult{
        Valid: len(errors) == 0,
        Errors: errors,
        Warnings: warnings,
    }
}
Quality Scoring:

bash
cicada metadata quality-check NIH-R01-2024/

# Output:
# Metadata Quality Report
# 
# Overall Score: 87/100 (Good)
# 
# Completeness: 92/100
#   ✓ 1,234 files have complete required metadata
#   ⚠ 45 files missing optional "replication" fields
#   ⚠ 12 files missing "protocol_id"
# 
# Consistency: 85/100
#   ✓ Controlled vocabularies followed
#   ⚠ 23 files use non-standard strain names
#   ⚠ Date formats inconsistent (3 different formats found)
# 
# Richness: 84/100
#   ✓ Good use of domain-specific fields
#   ⚠ Few files have "notes" or "keywords"
#   ⚠ Ontology mappings could be improved
# 
# Interoperability: 88/100
#   ✓ Standard schema used
#   ✓ Ontology terms present
#   ⚠ Some fields lack units
# 
# Recommendations:
#   1. Standardize strain names → [Fix automatically]
#   2. Add protocol IDs to 12 files → [Batch edit]
#   3. Convert dates to ISO 8601 → [Fix automatically]
Advanced Metadata Features
Metadata Inheritance:

bash
# Set metadata at folder level (applies to all files)
cicada metadata set NIH-R01-2024/raw/microscopy/batch-2024-11/ \
  --inherit \
  experiment_date=2024-11-20 \
  operator=jsmith \
  protocol_id=protocols.io/heat-shock-v3 \
  experimental.sample.strain=BY4741

# Files uploaded to this folder inherit these values
# Can still override per-file
Metadata Search & Discovery:

bash
# Search by metadata
cicada search \
  --schema microscopy-experiment \
  --where "treatment.condition=heat-shock" \
  --where "replication.biological_replicate>=2" \
  --where "microscopy.magnification=63"

# Output:
# Found 23 files matching criteria:
#   experiment_004.czi (heat-shock, rep 2, 63x)
#   experiment_005.czi (heat-shock, rep 2, 63x)
#   experiment_006.czi (heat-shock, rep 3, 63x)
#   ...
# 
# [📥 Download All] [📊 View Metadata Table] [🔗 Create Collection]
Metadata-Driven Analysis:

bash
# Run workflow on all files matching criteria
cicada workflow run cell-segmentation \
  --input-query "treatment=heat-shock AND timepoint<=30" \
  --group-by "biological_replicate"

# Automatically processes matching files
# Groups results by biological replicate
Export & Interoperability
Export to Standard Formats:

bash
# Export metadata for publication
cicada metadata export NIH-R01-2024/ \
  --format isa-tab \
  --output manuscript_supplement/

# Generates ISA-Tab files (standard for omics data)
# i_investigation.txt
# s_sample.txt
# a_assay.txt

# Other formats:
--format datacite-xml
--format dats-json      # DATS (Data Tags Suite)
--format frictionless   # Frictionless Data Package
--format ro-crate       # RO-Crate (Research Object)
This flexible metadata system means researchers can:

Start with sensible defaults (core schema)
Add domain-specific fields (microscopy, sequencing)
Customize for their specific needs (lab fields)
Share schemas with community
Maintain interoperability (ontologies, standard formats)
Publish with confidence (FAIR compliant)
Want me to dive deeper into any of these aspects, or explore how this integrates with the DOI workflow and public portal?

Yes, also I think it would be beneficial to consider a wide range of domains with examples.
I plan to develop this using Claude Code so artifacts that would support that would be very welcome

Cicada: Comprehensive Development Specification
I'll create a series of artifacts that will help you develop this with Claude Code. Let me start with the architecture and then provide domain-specific examples.

Part 1: Core Architecture & Project Structure

Main project README with overview, architecture, and getting started guide
README.md
Part 2: Domain-Specific Metadata Schemas

Failed to create

Comprehensive collection of domain-specific metadata schemas for various research fields
domain-schemas.md
Now let me create the core Go code structures for the metadata system:


Core metadata schema types and structures in Go
schema.go
Now let me create the metadata extractor system:


Metadata extractor interfaces and implementations for automatic metadata extraction from files
extractor.go
Let me now create the DOI management system and then example CLI commands:


DOI minting and management system with DataCite integration
datacite.go
Now let me create CLI examples and usage documentation:


Comprehensive CLI usage examples and command reference
cli-reference.md
Perfect! Now let me create one final comprehensive artifact - a development roadmap with implementation priorities:


Development roadmap with implementation priorities, technical specifications, and testing strategy
ROADMAP.md
Perfect! Let me create one more artifact with the complete file listing you can use as a reference:


List all created specification files


Now let me create a quick summary document:


Quick start guide for developers using these specifications
QUICKSTART.md
Perfect! Let me create one final artifact with example configuration files:


Example configuration files for Cicada
configs.md
Excellent! Now let me create a final summary of everything we've created:


Create master index file


Perfect! Now let me show you what we have:


Show final directory structure


Perfect! Let me create one final summary for you:

Summary
I've created a comprehensive specification package for Cicada - a "dormant" data commons platform for academic research labs. Here's what you have:

📦 Complete Package (10 Files, ~8,300 lines)
Start Here:

INDEX.md - Master index and navigation guide
QUICKSTART.md - Fastest path to understanding the project
Core Documentation: 3. README.md - Project overview and architecture 4. ROADMAP.md - 26-week implementation plan with detailed specs

User Guides: 5. docs/cli-reference.md - Complete CLI command examples 6. docs/domain-schemas.md - Research domain metadata schemas

Code Specifications: 7. internal/metadata/schema.go - Metadata system (700 lines of Go) 8. internal/metadata/extractor.go - File extractors (900 lines of Go) 9. internal/doi/datacite.go - DOI minting (800 lines of Go)

Configuration: 10. examples/configs.md - All config file templates

🎯 Key Features Covered
Core Functionality:

rsync-like S3 sync engine
Automatic file watching & upload
Flexible metadata schemas (12+ domains)
Automatic metadata extraction
FAIR-compliant data management
Advanced Features:

DOI minting (DataCite)
Workflow execution (Snakemake, Nextflow)
Remote GPU workstations
Public data portal
User/project management
Cost tracking & optimization
NIST 800-171 compliance mode
🔬 Domain Coverage
Metadata schemas for 12+ research fields:

Microscopy, Sequencing, Proteomics, Flow Cytometry
Chromatography, Spectroscopy, X-ray Crystallography
Electron Microscopy, Behavioral Studies, Clinical Trials
Environmental Sampling, Materials Science
💡 Design Philosophy
🦗 Dormant by Design - Resources only spin up when needed
💰 Cost-Conscious - $50-100/month typical cost
🎯 User-Friendly - No AWS knowledge required
📊 FAIR by Default - Research data standards
🔒 Compliant - HIPAA/NIST ready

🚀 Ready for Claude Code
All specifications are written to work seamlessly with Claude Code for implementation. Each phase has:

Clear technical specifications
Acceptance criteria
Test requirements
Example code structures
Configuration templates
📈 Development Timeline
Weeks 1-6: Core sync & daemon
Weeks 7-10: Metadata & FAIR
Weeks 11-14: Web UI & users
Weeks 15-18: Compute & workflows
Weeks 19-22: Workstations & portal
Weeks 23-26: Compliance & polish
All files are in /tmp/cicada-spec/ ready for you to use with Claude Code!

What is the full local directory path


Show full path and verify directory contents


The full local directory path is:

/tmp/cicada-spec/
Here's the complete structure with all files:

/tmp/cicada-spec/
├── INDEX.md                      # Master index (start here!)
├── QUICKSTART.md                 # Developer quick start guide
├── README.md                     # Project overview & architecture
├── ROADMAP.md                    # 26-week implementation plan
│
├── docs/
│   ├── cli-reference.md         # Complete CLI examples
│   └── domain-schemas.md        # Research domain metadata schemas
│
├── examples/
│   └── configs.md               # Configuration file templates
│
└── internal/
    ├── metadata/
    │   ├── schema.go            # Metadata system (Go code)
    │   └── extractor.go         # File metadata extractors (Go code)
    └── doi/
        └── datacite.go          # DOI minting system (Go code)
To use these files:
bash
# Copy to your project directory
cp -r /tmp/cicada-spec ~/cicada-project/

# Or access directly
cd /tmp/cicada-spec
ls -R
All files are ready for you to use with Claude Code for implementation!

ls /tmp/cicada-spec
ls: /tmp/cicada-spec: No such file or directory

You're right - the files were created in the Claude environment but aren't accessible to you directly. Let me create a downloadable archive for you:


Create compressed archive of all specification files



Copy archive to outputs directory for download


Perfect! I've created a compressed archive with all the specification files.

Download cicada-spec.tar.gz (39 KB)

To extract and use:
bash
# Extract the archive
tar -xzf cicada-spec.tar.gz

# Navigate to the directory
cd cicada-spec

# Start reading
cat INDEX.md
What's inside:
10 files with ~8,300 lines of specifications
Complete architecture and design documentation
26-week implementation roadmap
Go code specifications ready for Claude Code
CLI examples and configuration templates
Metadata schemas for 12+ research domains
All ready for you to start building Cicada! 🦗

Cicada spec.tar
GZ 


How can I help you today?







