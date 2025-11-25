# Common Workflows

**Last Updated:** 2025-11-25

Practical workflows and patterns for using Cicada in research environments.

---

## Table of Contents

1. [Microscopy Workflows](#microscopy-workflows)
2. [Sequencing Workflows](#sequencing-workflows)
3. [Mass Spectrometry Workflows](#mass-spectrometry-workflows)
4. [Collaborative Workflows](#collaborative-workflows)
5. [Backup and Archival Workflows](#backup-and-archival-workflows)
6. [Integration Workflows](#integration-workflows)
7. [Data Quality Workflows](#data-quality-workflows)

---

## Microscopy Workflows

### Workflow 1: Daily Microscope Data Collection

**Scenario**: Transfer data from microscope workstation to lab storage daily, extract metadata, and backup to S3.

**Setup:**

```bash
# On microscope workstation
# Add to ~/.config/cicada/config.yaml

sync:
  local_root: /Volumes/Microscope/Export
  remote_root: ~/lab-storage/microscopy

watch:
  enabled: true
  paths:
    - /Volumes/Microscope/Export
  on_change:
    - action: sync
      destination: local
      target: ~/lab-storage/microscopy
    - action: extract_metadata

# Start watch daemon
cicada watch start
```

**Daily Workflow:**

```bash
# 1. Acquire images on microscope (automatic export)
# 2. Cicada watches and automatically:
#    - Copies to lab storage
#    - Extracts metadata
# 3. At end of day, backup to S3

cicada sync ~/lab-storage/microscopy s3://lab-bucket/microscopy/$(date +%Y-%m)

# 4. Verify metadata extraction
cicada metadata list ~/lab-storage/microscopy --date today

# 5. Generate daily report
cicada metadata export ~/lab-storage/microscopy \
  --filter "date=$(date +%Y-%m-%d)" \
  --format csv > reports/microscopy_$(date +%Y%m%d).csv
```

**Automation Script** (`~/scripts/microscopy-daily.sh`):

```bash
#!/bin/bash
# Daily microscopy data workflow

DATE=$(date +%Y-%m-%d)
MONTH=$(date +%Y-%m)
LOCAL_DIR=~/lab-storage/microscopy
S3_BUCKET=s3://lab-bucket/microscopy

# Verify watch daemon is running
if ! cicada watch status > /dev/null 2>&1; then
    echo "Starting watch daemon..."
    cicada watch start
fi

# Evening backup to S3
echo "Backing up today's data to S3..."
cicada sync $LOCAL_DIR/$DATE $S3_BUCKET/$MONTH/$DATE --progress

# Generate metadata report
echo "Generating metadata report..."
cicada metadata export $LOCAL_DIR \
  --filter "date=$DATE" \
  --format csv > ~/reports/microscopy_${DATE}.csv

# Verify sync
cicada s3 ls $S3_BUCKET/$MONTH/$DATE --recursive --summarize

echo "Daily backup complete!"
```

**Schedule with cron:**

```bash
# Run at 6 PM daily
0 18 * * * ~/scripts/microscopy-daily.sh >> ~/logs/microscopy-backup.log 2>&1
```

### Workflow 2: Multi-User Microscope Facility

**Scenario**: Core facility with multiple users and microscopes.

**Directory Structure:**

```
~/facility-data/
├── microscope-1/
│   ├── user-a/
│   ├── user-b/
│   └── user-c/
├── microscope-2/
│   ├── user-a/
│   └── user-d/
└── archive/
    ├── 2025-01/
    ├── 2025-02/
    └── 2025-03/
```

**Configuration** (`~/.config/cicada/config.yaml`):

```yaml
sync:
  local_root: ~/facility-data
  remote_root: s3://facility-bucket/data

metadata:
  auto_extract: true
  presets:
    microscope_1:
      instrument: "Nikon Ti2-E"
      location: "Room 101"
    microscope_2:
      instrument: "Zeiss LSM 900"
      location: "Room 102"

watch:
  enabled: true
  paths:
    - ~/facility-data/microscope-1
    - ~/facility-data/microscope-2
  on_change:
    - action: sync
      destination: s3
    - action: notify_user
```

**User Data Collection:**

```bash
# When user starts session
cicada metadata add ~/facility-data/microscope-1/user-a \
  --field user_id=user-a \
  --field user_name="Alice Smith" \
  --field project="Cell division study" \
  --field grant_number="R01-GM123456"

# Microscope exports data to user directory
# Cicada automatically:
#   - Extracts metadata
#   - Applies user metadata
#   - Uploads to S3: s3://facility-bucket/data/microscope-1/user-a/

# At end of session, generate report
cicada metadata export ~/facility-data/microscope-1/user-a \
  --date today \
  --format json > session_report_$(date +%Y%m%d).json
```

**Monthly Archival:**

```bash
#!/bin/bash
# Archive last month's data

LAST_MONTH=$(date -d "last month" +%Y-%m)
ARCHIVE_DIR=~/facility-data/archive/$LAST_MONTH

# Move data to archive
for microscope in microscope-1 microscope-2; do
    for user_dir in ~/facility-data/$microscope/*; do
        user=$(basename $user_dir)

        # Find files from last month
        find $user_dir -type f -newermt "${LAST_MONTH}-01" ! -newermt "$(date +%Y-%m)-01" \
          -exec mv {} $ARCHIVE_DIR/$microscope/$user/ \;
    done
done

# Upload archive to Glacier
cicada sync $ARCHIVE_DIR s3://facility-bucket/archive/$LAST_MONTH \
  --storage-class GLACIER

# Generate monthly report
cicada metadata export $ARCHIVE_DIR \
  --format csv > ~/reports/facility_${LAST_MONTH}.csv
```

### Workflow 3: Time-Lapse Experiment Management

**Scenario**: Long-running time-lapse experiments with real-time monitoring.

**Setup:**

```bash
# Create experiment directory
mkdir -p ~/experiments/timelapse_2025-11-25_cell-division
cd ~/experiments/timelapse_2025-11-25_cell-division

# Add experiment metadata
cat > experiment.yaml <<EOF
experiment_id: EXP_2025_11_25_001
name: "Cell division time-lapse"
investigator: "Alice Smith"
start_time: "2025-11-25T08:00:00Z"
duration_hours: 48
interval_minutes: 5
sample:
  cell_line: "HeLa"
  treatment: "Control"
  temperature: 37
  co2_percent: 5
microscope:
  name: "Nikon Ti2-E"
  objective: "20x Phase"
  channels:
    - PHASE
    - GFP
EOF

# Start watching for new images
cicada watch . --sync-to s3://lab-bucket/experiments/$(basename $(pwd))
```

**During Experiment:**

```bash
# Check progress in real-time
cicada metadata list . --watch --refresh 60s

# Verify latest images
cicada metadata show --latest 5

# Monitor disk usage
watch -n 60 'cicada s3 du s3://lab-bucket/experiments/timelapse_2025-11-25_cell-division'

# Check for imaging problems
cicada metadata list . --filter "status=error" | mail -s "Imaging errors detected" lab@example.com
```

**Post-Experiment Analysis:**

```bash
# Extract all time-points
cicada metadata export . --format json > timelapse_metadata.json

# Calculate actual intervals
python3 << 'PYTHON'
import json
import pandas as pd
from datetime import datetime

with open('timelapse_metadata.json') as f:
    data = json.load(f)

df = pd.DataFrame(data)
df['timestamp'] = pd.to_datetime(df['acquisition_timestamp'])
df = df.sort_values('timestamp')
df['interval'] = df['timestamp'].diff().dt.total_seconds() / 60

print(f"Planned interval: 5 minutes")
print(f"Actual mean interval: {df['interval'].mean():.2f} minutes")
print(f"Interval std dev: {df['interval'].std():.2f} minutes")
print(f"Total timepoints: {len(df)}")
print(f"Expected timepoints: {48 * 60 / 5}")
PYTHON

# Archive completed experiment
cicada sync . s3://lab-bucket/experiments/completed/$(basename $(pwd)) \
  --storage-class STANDARD_IA
```

---

## Sequencing Workflows

### Workflow 1: Illumina Sequencing Run Management

**Scenario**: Managing Illumina sequencing runs from BaseSpace to lab storage.

**Setup:**

```bash
# Directory structure
mkdir -p ~/sequencing/{raw,fastqc,aligned,analyzed}

# Configure Cicada
cat >> ~/.config/cicada/config.yaml << 'EOF'
sync:
  local_root: ~/sequencing
  remote_root: s3://lab-bucket/sequencing

metadata:
  auto_extract: true
  formats:
    - fastq
    - bam
    - vcf

watch:
  paths:
    - ~/sequencing/raw
  on_change:
    - action: extract_metadata
    - action: run_fastqc
    - action: sync
      destination: s3
EOF
```

**Workflow:**

```bash
# 1. Download from BaseSpace (or receive from core facility)
bs download run -i 123456789 -o ~/sequencing/raw/run_20251125

# 2. Cicada automatically extracts metadata
# (Watch daemon detects new files)

# 3. Manually verify metadata
cicada metadata list ~/sequencing/raw/run_20251125

# Example output:
# sample_001_R1.fastq.gz: RNA-seq, 150bp, 42.3M reads
# sample_001_R2.fastq.gz: RNA-seq, 150bp, 42.3M reads
# sample_002_R1.fastq.gz: RNA-seq, 150bp, 38.1M reads
# sample_002_R2.fastq.gz: RNA-seq, 150bp, 38.1M reads

# 4. Add sample metadata
cicada metadata add ~/sequencing/raw/run_20251125/sample_001_*.fastq.gz \
  --field sample_id=SAMPLE001 \
  --field condition=control \
  --field replicate=1 \
  --field experiment=RNA_SEQ_2025_11

# 5. Generate run summary
cicada metadata export ~/sequencing/raw/run_20251125 \
  --format csv > run_20251125_summary.csv

# 6. Sync to S3 (happens automatically via watch)
# Or manually:
cicada sync ~/sequencing/raw/run_20251125 \
  s3://lab-bucket/sequencing/raw/run_20251125 \
  --progress

# 7. Run QC pipeline (see Integration Workflows)
```

### Workflow 2: Multi-Sample RNA-Seq Project

**Scenario**: Organizing and tracking samples across an RNA-seq experiment.

**Project Setup:**

```bash
# Create project structure
PROJECT=~/projects/rnaseq_treatment_comparison
mkdir -p $PROJECT/{raw,fastqc,trimmed,aligned,counts,metadata}

# Create sample sheet
cat > $PROJECT/metadata/sample_sheet.csv << 'EOF'
sample_id,condition,replicate,fastq_r1,fastq_r2
CTRL_1,control,1,sample_001_R1.fastq.gz,sample_001_R2.fastq.gz
CTRL_2,control,2,sample_002_R1.fastq.gz,sample_002_R2.fastq.gz
CTRL_3,control,3,sample_003_R1.fastq.gz,sample_003_R2.fastq.gz
TREAT_1,treatment,1,sample_004_R1.fastq.gz,sample_004_R2.fastq.gz
TREAT_2,treatment,2,sample_005_R1.fastq.gz,sample_005_R2.fastq.gz
TREAT_3,treatment,3,sample_006_R1.fastq.gz,sample_006_R2.fastq.gz
EOF

# Import sample metadata into Cicada
python3 << 'PYTHON'
import csv
import subprocess

with open('metadata/sample_sheet.csv') as f:
    reader = csv.DictReader(f)
    for row in reader:
        for read in ['fastq_r1', 'fastq_r2']:
            fastq = f"raw/{row[read]}"
            cmd = [
                'cicada', 'metadata', 'add', fastq,
                '--field', f"sample_id={row['sample_id']}",
                '--field', f"condition={row['condition']}",
                '--field', f"replicate={row['replicate']}",
                '--field', f"project=rnaseq_treatment_comparison"
            ]
            subprocess.run(cmd)
PYTHON
```

**Analysis Workflow:**

```bash
# 1. Quality control
for fastq in $PROJECT/raw/*.fastq.gz; do
    fastqc -o $PROJECT/fastqc $fastq
done

# 2. Track QC metadata
cicada metadata add $PROJECT/fastqc/*.html \
  --field analysis_step=quality_control \
  --field tool=fastqc \
  --field version=$(fastqc --version)

# 3. Trimming
for sample in CTRL_1 CTRL_2 CTRL_3 TREAT_1 TREAT_2 TREAT_3; do
    # Get fastq files for this sample
    R1=$(cicada metadata list $PROJECT/raw --filter "sample_id=$sample" | grep R1)
    R2=$(cicada metadata list $PROJECT/raw --filter "sample_id=$sample" | grep R2)

    # Trim adapters
    trim_galore --paired -o $PROJECT/trimmed $R1 $R2

    # Add metadata to trimmed files
    cicada metadata add $PROJECT/trimmed/${sample}_*.fq.gz \
      --field sample_id=$sample \
      --field analysis_step=trimming \
      --field tool=trim_galore
done

# 4. Alignment
for sample in CTRL_1 CTRL_2 CTRL_3 TREAT_1 TREAT_2 TREAT_3; do
    STAR --runThreadN 8 \
      --genomeDir /data/genomes/human/STAR \
      --readFilesIn $PROJECT/trimmed/${sample}_R1.fq.gz \
                     $PROJECT/trimmed/${sample}_R2.fq.gz \
      --readFilesCommand zcat \
      --outFileNamePrefix $PROJECT/aligned/${sample}_ \
      --outSAMtype BAM SortedByCoordinate

    # Index BAM
    samtools index $PROJECT/aligned/${sample}_Aligned.sortedByCoord.out.bam

    # Add metadata
    cicada metadata add $PROJECT/aligned/${sample}_*.bam \
      --field sample_id=$sample \
      --field analysis_step=alignment \
      --field tool=STAR \
      --field genome=GRCh38
done

# 5. Generate final project report
cicada metadata export $PROJECT \
  --recursive \
  --format json > $PROJECT/metadata/project_metadata.json

# 6. Backup entire project to S3
cicada sync $PROJECT s3://lab-bucket/projects/rnaseq_treatment_comparison
```

### Workflow 3: Nanopore Sequencing Real-Time Processing

**Scenario**: Processing Nanopore sequencing data as it's generated.

**Setup:**

```bash
# Create directories
mkdir -p ~/nanopore/{raw,basecalled,qc,aligned}

# Start MinKNOW sequencing run
# Output directory: ~/nanopore/raw/run_20251125

# Watch for new FAST5 files and basecall
cicada watch ~/nanopore/raw/run_20251125 \
  --on-change "guppy_basecaller -i {} -s ~/nanopore/basecalled/run_20251125 -c dna_r9.4.1_450bps_fast.cfg" \
  --sync-to s3://lab-bucket/nanopore/run_20251125
```

**Real-Time Monitoring:**

```bash
# Monitor read count in real-time
watch -n 30 'cicada metadata list ~/nanopore/basecalled/run_20251125 --format summary'

# Check quality metrics
python3 << 'PYTHON'
import json
import subprocess

result = subprocess.run(
    ['cicada', 'metadata', 'export', '~/nanopore/basecalled/run_20251125', '--format', 'json'],
    capture_output=True, text=True
)

data = json.loads(result.stdout)
reads = [d for d in data if d['type'] == 'fastq']

total_reads = sum(r['metadata']['read_count'] for r in reads)
total_bases = sum(r['metadata']['total_bases'] for r in reads)
mean_quality = sum(r['metadata']['mean_quality'] * r['metadata']['read_count'] for r in reads) / total_reads

print(f"Total reads: {total_reads:,}")
print(f"Total bases: {total_bases:,} ({total_bases/1e9:.2f} Gb)")
print(f"Mean quality: {mean_quality:.2f}")
print(f"N50: {calculate_n50([r['metadata']['read_lengths'] for r in reads])}")
PYTHON
```

---

## Mass Spectrometry Workflows

### Workflow 1: Proteomics Data Management

**Scenario**: LC-MS/MS proteomics data from core facility to analysis.

**Setup:**

```bash
# Create project structure
PROJECT=~/proteomics/project_2025_11_protein_expression
mkdir -p $PROJECT/{raw,search_results,reports}

# Download raw files from instrument PC
scp instrument@192.168.1.100:/data/export/*.raw $PROJECT/raw/

# Extract metadata
cicada metadata extract $PROJECT/raw --recursive

# Add experimental metadata
for sample in sample_*.raw; do
    sample_id=$(basename $sample .raw)

    cicada metadata add $PROJECT/raw/$sample \
      --field project=protein_expression_2025_11 \
      --field sample_id=$sample_id \
      --field instrument="Thermo Q Exactive HF" \
      --field method="DDA_120min_gradient" \
      --field date=$(date +%Y-%m-%d)
done
```

**Database Search:**

```bash
# Run MaxQuant or similar
maxquant maxquant_config.xml

# Add search results metadata
cicada metadata add $PROJECT/search_results/combined/txt/*.txt \
  --field analysis_step=database_search \
  --field tool=MaxQuant \
  --field version=2.0.3.0 \
  --field database=UniProt_Human_2025_11 \
  --field fdr_threshold=0.01

# Generate summary report
cicada metadata export $PROJECT \
  --recursive \
  --format csv > $PROJECT/reports/project_summary.csv

# Backup to S3
cicada sync $PROJECT s3://lab-bucket/proteomics/$(basename $PROJECT)
```

### Workflow 2: Metabolomics Sample Batch Processing

**Scenario**: Process batch of metabolomics samples with QC samples.

**Batch Processing Script:**

```bash
#!/bin/bash
# Process metabolomics batch

BATCH_DIR=$1
BATCH_ID=$(basename $BATCH_DIR)

echo "Processing batch: $BATCH_ID"

# 1. Extract metadata from all mzML files
echo "Extracting metadata..."
cicada metadata extract $BATCH_DIR/*.mzML

# 2. Identify QC samples
QC_SAMPLES=$(cicada metadata list $BATCH_DIR --filter "sample_type=QC" --format path)

# 3. Check QC metrics
echo "Checking QC samples..."
for qc in $QC_SAMPLES; do
    metrics=$(cicada metadata show $qc --field qc_metrics)
    # Validate metrics
    if [ "$?" -ne 0 ]; then
        echo "WARNING: QC failed for $qc"
        exit 1
    fi
done

# 4. Process biological samples
echo "Processing biological samples..."
SAMPLES=$(cicada metadata list $BATCH_DIR --filter "sample_type=biological" --format path)

for sample in $SAMPLES; do
    # Peak picking
    msconvert $sample --filter "peakPicking true" -o processed/

    # Add processing metadata
    cicada metadata add processed/$(basename $sample) \
      --field processing_step=peak_picking \
      --field batch_id=$BATCH_ID \
      --field processing_date=$(date +%Y-%m-%d)
done

# 5. Upload batch to S3
echo "Uploading to S3..."
cicada sync $BATCH_DIR s3://lab-bucket/metabolomics/batches/$BATCH_ID

# 6. Generate batch report
cicada metadata export $BATCH_DIR \
  --format csv > reports/batch_${BATCH_ID}_report.csv

echo "Batch processing complete!"
```

---

## Collaborative Workflows

### Workflow 1: Multi-Lab Collaboration

**Scenario**: Share data between multiple labs via S3.

**Setup - Lab A (Data Producer):**

```bash
# Lab A configuration
cat > ~/.config/cicada/config.yaml << 'EOF'
sync:
  local_root: ~/lab-data
  remote_root: s3://collaboration-bucket/lab-a-data

sharing:
  enabled: true
  bucket: s3://collaboration-bucket
  access:
    lab-b: read-write
    lab-c: read-only
EOF

# Upload data with metadata
cicada sync ~/lab-data/experiment-001 \
  s3://collaboration-bucket/lab-a-data/experiment-001

# Share specific dataset
cicada share create s3://collaboration-bucket/lab-a-data/experiment-001 \
  --grant lab-b:read-write \
  --grant lab-c:read-only \
  --expires 30d
```

**Setup - Lab B (Collaborator):**

```bash
# Lab B configuration
cat > ~/.config/cicada/config.yaml << 'EOF'
sync:
  remote_roots:
    - s3://collaboration-bucket/lab-a-data
    - s3://collaboration-bucket/lab-b-data
  local_root: ~/collaboration-data

subscriptions:
  - bucket: s3://collaboration-bucket/lab-a-data
    filters:
      - experiment-001
    auto_sync: true
    sync_interval: 1h
EOF

# Subscribe to Lab A's data
cicada subscribe s3://collaboration-bucket/lab-a-data/experiment-001

# Automatic sync (via watch daemon)
cicada watch start

# Or manual sync
cicada sync s3://collaboration-bucket/lab-a-data/experiment-001 \
  ~/collaboration-data/lab-a/experiment-001
```

**Contribution Workflow:**

```bash
# Lab B adds analysis results
mkdir -p ~/collaboration-data/lab-a/experiment-001/analysis-lab-b

# Run analysis
python analysis.py \
  --input ~/collaboration-data/lab-a/experiment-001/raw \
  --output ~/collaboration-data/lab-a/experiment-001/analysis-lab-b

# Add metadata
cicada metadata add ~/collaboration-data/lab-a/experiment-001/analysis-lab-b \
  --field contributor="Lab B" \
  --field analysis_type="differential_expression" \
  --field tool="DESeq2" \
  --field date=$(date +%Y-%m-%d)

# Upload contribution
cicada sync ~/collaboration-data/lab-a/experiment-001/analysis-lab-b \
  s3://collaboration-bucket/lab-a-data/experiment-001/analysis-lab-b
```

### Workflow 2: PI-PostDoc-Student Hierarchy

**Scenario**: Organize data access for research group members.

**Directory Structure:**

```
s3://group-bucket/
├── raw-data/              # PI and PostDocs: read-write, Students: read-only
├── processed-data/        # PostDocs: read-write, Students: read-only
├── student-projects/      # Each student has their own folder
│   ├── alice/            # Alice: read-write
│   ├── bob/              # Bob: read-write
│   └── carol/            # Carol: read-write
└── publications/          # Everyone: read-only
```

**PI Configuration:**

```yaml
# PI: Full access to everything
sync:
  remote_root: s3://group-bucket
  local_root: ~/group-data

access_control:
  enabled: true
  users:
    postdoc_dave:
      permissions: read-write
      paths:
        - raw-data
        - processed-data
    student_alice:
      permissions:
        raw-data: read-only
        processed-data: read-only
        student-projects/alice: read-write
    student_bob:
      permissions:
        raw-data: read-only
        processed-data: read-only
        student-projects/bob: read-write
```

**Student Workflow (Alice):**

```bash
# Alice's config
cat > ~/.config/cicada/config.yaml << 'EOF'
sync:
  remote_root: s3://group-bucket
  local_root: ~/group-data

subscriptions:
  - path: s3://group-bucket/raw-data
    mode: read-only
    auto_sync: true
  - path: s3://group-bucket/student-projects/alice
    mode: read-write
    auto_sync: true
EOF

# Download latest raw data
cicada sync s3://group-bucket/raw-data ~/group-data/raw-data

# Work on analysis
python my_analysis.py \
  --input ~/group-data/raw-data \
  --output ~/group-data/student-projects/alice/analysis

# Upload results
cicada sync ~/group-data/student-projects/alice \
  s3://group-bucket/student-projects/alice

# Request PI review
cicada notify --to pi@example.com \
  --subject "Analysis complete: Alice's project" \
  --message "Results uploaded to student-projects/alice/analysis" \
  --path s3://group-bucket/student-projects/alice/analysis
```

---

## Backup and Archival Workflows

### Workflow 1: Tiered Storage Strategy

**Scenario**: Implement hot/warm/cold storage tiers based on data age.

**Strategy:**

- **Hot (STANDARD)**: Current projects, accessed daily
- **Warm (STANDARD_IA)**: Recent projects, accessed monthly
- **Cold (GLACIER)**: Old projects, accessed rarely
- **Archive (DEEP_ARCHIVE)**: Published data, rarely accessed

**Implementation:**

```bash
# Configure lifecycle policies
cicada s3 lifecycle create s3://lab-bucket \
  --rule-name "Tier to Standard IA" \
  --transition-days 30 \
  --transition-class STANDARD_IA \
  --prefix "projects/"

cicada s3 lifecycle create s3://lab-bucket \
  --rule-name "Tier to Glacier" \
  --transition-days 90 \
  --transition-class GLACIER \
  --prefix "projects/"

cicada s3 lifecycle create s3://lab-bucket \
  --rule-name "Archive published data" \
  --transition-days 180 \
  --transition-class DEEP_ARCHIVE \
  --prefix "publications/"

# Monitor storage classes
cicada s3 storage-class s3://lab-bucket --report monthly
```

**Manual Archival:**

```bash
# Archive completed project to Glacier
cicada sync ~/completed-projects/project-2024-01 \
  s3://lab-bucket/archive/2024/project-2024-01 \
  --storage-class GLACIER

# Restore from Glacier when needed
cicada s3 restore s3://lab-bucket/archive/2024/project-2024-01 \
  --days 7  # Available for 7 days

# Wait for restore (takes 3-5 hours)
cicada s3 restore-status s3://lab-bucket/archive/2024/project-2024-01

# Download restored data
cicada sync s3://lab-bucket/archive/2024/project-2024-01 \
  ~/restored/project-2024-01
```

### Workflow 2: Incremental Backup Strategy

**Scenario**: Daily incremental backups with weekly full backups.

**Backup Script** (`~/scripts/daily-backup.sh`):

```bash
#!/bin/bash
# Daily incremental backup

DATE=$(date +%Y-%m-%d)
DOW=$(date +%u)  # 1-7 (Monday-Sunday)

LOCAL_ROOT=~/lab-data
S3_BUCKET=s3://lab-bucket/backups

if [ "$DOW" -eq "7" ]; then
    # Sunday: Full backup
    echo "Performing full backup..."
    cicada sync $LOCAL_ROOT $S3_BUCKET/full/$DATE --progress

    # Delete incremental backups older than this full backup
    cicada s3 rm $S3_BUCKET/incremental/ --recursive --older-than 7days
else
    # Monday-Saturday: Incremental backup
    echo "Performing incremental backup..."

    # Find files modified in last 24 hours
    find $LOCAL_ROOT -type f -mtime -1 > /tmp/changed_files.txt

    # Sync only changed files
    cicada sync $LOCAL_ROOT $S3_BUCKET/incremental/$DATE \
      --files-from /tmp/changed_files.txt \
      --progress
fi

# Verify backup
cicada s3 verify $S3_BUCKET/$([ "$DOW" -eq "7" ] && echo "full" || echo "incremental")/$DATE

# Generate backup report
cicada s3 du $S3_BUCKET --by-date > ~/reports/backup_${DATE}.txt

# Email report
mail -s "Backup complete: $DATE" admin@lab.example.com < ~/reports/backup_${DATE}.txt
```

**Restore Process:**

```bash
# Restore from latest full backup
LATEST_FULL=$(cicada s3 ls s3://lab-bucket/backups/full --sort date | tail -1)
cicada sync s3://lab-bucket/backups/full/$LATEST_FULL ~/restored-data

# Apply incremental backups since full backup
for incremental in $(cicada s3 ls s3://lab-bucket/backups/incremental \
                      --after $LATEST_FULL --sort date); do
    cicada sync s3://lab-bucket/backups/incremental/$incremental ~/restored-data
done

echo "Restore complete!"
```

---

## Integration Workflows

### Workflow 1: Snakemake Pipeline Integration

**Scenario**: Integrate Cicada with Snakemake for reproducible workflows.

**Snakefile:**

```python
# Snakefile with Cicada integration

import subprocess
import json

# Helper function to run Cicada commands
def cicada(cmd):
    result = subprocess.run(
        f"cicada {cmd}",
        shell=True, capture_output=True, text=True
    )
    return result.stdout

# Download input data from S3
rule fetch_data:
    output:
        "data/raw/{sample}.fastq.gz"
    params:
        s3_path="s3://lab-bucket/sequencing/raw/{sample}.fastq.gz"
    shell:
        "cicada sync {params.s3_path} {output}"

# Run analysis
rule align:
    input:
        "data/raw/{sample}.fastq.gz"
    output:
        "data/aligned/{sample}.bam"
    params:
        metadata=lambda wildcards: cicada(f"metadata show data/raw/{wildcards.sample}.fastq.gz --json")
    threads: 8
    shell:
        """
        # Parse metadata to get genome
        genome=$(echo '{params.metadata}' | jq -r '.metadata.genome')

        # Align
        bwa mem -t {threads} /data/genomes/$genome {input} | samtools sort -o {output}

        # Add metadata
        cicada metadata add {output} \
          --field analysis_step=alignment \
          --field tool=bwa \
          --field genome=$genome
        """

# Upload results
rule upload_results:
    input:
        "data/aligned/{sample}.bam"
    params:
        s3_path="s3://lab-bucket/analysis/aligned/{sample}.bam"
    shell:
        "cicada sync {input} {params.s3_path}"

# All samples
SAMPLES = ["sample1", "sample2", "sample3"]

rule all:
    input:
        expand("data/aligned/{sample}.bam", sample=SAMPLES)
```

**Run Pipeline:**

```bash
# Run Snakemake with Cicada integration
snakemake --cores 8 all

# All input data is automatically fetched from S3
# All results are automatically uploaded back to S3
# Metadata is preserved throughout the pipeline
```

### Workflow 2: Nextflow with Cicada

**Scenario**: Use Cicada for data management in Nextflow pipelines.

**nextflow.config:**

```groovy
// Nextflow configuration with Cicada

process {
    // Before each process, fetch required data
    beforeScript = 'cicada sync s3://lab-bucket/data/$PWD/input .'

    // After each process, upload results
    afterScript = 'cicada sync . s3://lab-bucket/data/$PWD/output'
}
```

**main.nf:**

```groovy
#!/usr/bin/env nextflow

// Fetch sample list from Cicada metadata
process get_samples {
    output:
    path 'samples.txt'

    script:
    """
    cicada metadata list s3://lab-bucket/sequencing/raw \
      --filter "project=rnaseq_2025" \
      --format list > samples.txt
    """
}

// Process each sample
process analyze {
    input:
    val sample from samples.txt.splitText()

    output:
    path "${sample}.results.txt"

    script:
    """
    # Cicada automatically synced input via beforeScript

    # Run analysis
    analyze_sample.py $sample > ${sample}.results.txt

    # Add metadata
    cicada metadata add ${sample}.results.txt \
      --field sample=$sample \
      --field analysis_step=process \
      --field pipeline=nextflow

    # Cicada will automatically sync output via afterScript
    """
}

workflow {
    get_samples | analyze
}
```

---

## Data Quality Workflows

### Workflow 1: Automated Quality Checks

**Scenario**: Automatically validate data quality on upload.

**Quality Check Script** (`~/.config/cicada/hooks/on_upload.sh`):

```bash
#!/bin/bash
# Automated quality checks on file upload

FILE=$1
FILETYPE=$2

echo "Running quality checks on $FILE..."

case $FILETYPE in
    microscopy)
        # Check image dimensions
        WIDTH=$(cicada metadata show $FILE --field "metadata.dimensions.width")
        HEIGHT=$(cicada metadata show $FILE --field "metadata.dimensions.height")

        if [ "$WIDTH" -lt 512 ] || [ "$HEIGHT" -lt 512 ]; then
            echo "WARNING: Image resolution below minimum (512x512): ${WIDTH}x${HEIGHT}"
            cicada metadata add $FILE --field "qc_status=warning" --field "qc_reason=low_resolution"
        fi
        ;;

    sequencing)
        # Check read count
        READS=$(cicada metadata show $FILE --field "metadata.read_count")

        if [ "$READS" -lt 1000000 ]; then
            echo "WARNING: Low read count: $READS"
            cicada metadata add $FILE --field "qc_status=warning" --field "qc_reason=low_read_count"
        fi

        # Check quality scores
        MEAN_Q=$(cicada metadata show $FILE --field "metadata.mean_quality")

        if (( $(echo "$MEAN_Q < 30" | bc -l) )); then
            echo "WARNING: Low quality score: $MEAN_Q"
            cicada metadata add $FILE --field "qc_status=warning" --field "qc_reason=low_quality"
        fi
        ;;

    mass_spec)
        # Check acquisition time
        DURATION=$(cicada metadata show $FILE --field "metadata.acquisition_duration_sec")

        if [ "$DURATION" -lt 600 ]; then
            echo "WARNING: Short acquisition time: ${DURATION}s"
            cicada metadata add $FILE --field "qc_status=warning" --field "qc_reason=short_acquisition"
        fi
        ;;
esac

# Check file integrity
if ! cicada verify $FILE; then
    echo "ERROR: File integrity check failed"
    cicada metadata add $FILE --field "qc_status=error" --field "qc_reason=integrity_failed"
    exit 1
fi

echo "Quality checks complete"
cicada metadata add $FILE --field "qc_status=passed" --field "qc_timestamp=$(date -Iseconds)"
```

**Configure Hook:**

```yaml
# ~/.config/cicada/config.yaml
hooks:
  on_upload:
    enabled: true
    script: ~/.config/cicada/hooks/on_upload.sh
    timeout: 300  # 5 minutes
```

### Workflow 2: Periodic Data Validation

**Scenario**: Regular validation of stored data integrity.

**Validation Script** (`~/scripts/validate-data.sh`):

```bash
#!/bin/bash
# Periodic data validation

LOCAL_ROOT=~/lab-data
S3_BUCKET=s3://lab-bucket/data

echo "Starting data validation: $(date)"

# 1. Check for missing metadata
echo "Checking for files without metadata..."
cicada metadata list $LOCAL_ROOT --filter "metadata_missing=true" > /tmp/missing_metadata.txt

if [ -s /tmp/missing_metadata.txt ]; then
    echo "WARNING: Found $(wc -l < /tmp/missing_metadata.txt) files without metadata"

    # Try to extract metadata
    while read file; do
        echo "  Extracting metadata for $file..."
        cicada metadata extract "$file"
    done < /tmp/missing_metadata.txt
fi

# 2. Verify local-S3 consistency
echo "Verifying local-S3 consistency..."
cicada sync --dry-run --checksum $LOCAL_ROOT $S3_BUCKET > /tmp/sync_diff.txt

if grep -q "would sync" /tmp/sync_diff.txt; then
    echo "WARNING: Local and S3 are out of sync"
    cat /tmp/sync_diff.txt

    # Optionally auto-sync
    # cicada sync $LOCAL_ROOT $S3_BUCKET
fi

# 3. Check for corrupted files
echo "Checking file integrity..."
find $LOCAL_ROOT -type f -name "*.fastq.gz" -exec bash -c '
    if ! gzip -t "$1" 2>/dev/null; then
        echo "ERROR: Corrupted file: $1"
        cicada metadata add "$1" --field "qc_status=error" --field "qc_reason=corrupted"
    fi
' _ {} \;

# 4. Validate metadata schema
echo "Validating metadata schema..."
cicada metadata validate $LOCAL_ROOT --schema ~/.config/cicada/schemas/default.json

# 5. Generate validation report
echo "Generating validation report..."
cat > ~/reports/validation_$(date +%Y%m%d).txt << EOF
Data Validation Report
======================
Date: $(date)

Files without metadata: $(wc -l < /tmp/missing_metadata.txt)
Sync differences: $(grep -c "would sync" /tmp/sync_diff.txt)
Corrupted files: $(grep -c "ERROR: Corrupted" /tmp/validation.log)
Schema violations: $(cicada metadata validate $LOCAL_ROOT --count-errors)

Details:
$(cat /tmp/validation.log)
EOF

echo "Validation complete. Report: ~/reports/validation_$(date +%Y%m%d).txt"

# Email report
mail -s "Data Validation Report $(date +%Y-%m-%d)" admin@lab.example.com \
  < ~/reports/validation_$(date +%Y%m%d).txt
```

**Schedule with cron:**

```bash
# Run validation every Monday at 2 AM
0 2 * * 1 ~/scripts/validate-data.sh >> ~/logs/validation.log 2>&1
```

---

## Quick Reference

### Common Command Patterns

```bash
# Daily data collection
cicada watch ~/data --sync-to s3://bucket/data

# Project setup
cicada metadata add ~/project --field project_id=PRJ001 --recursive

# Backup
cicada sync ~/data s3://backup/$(date +%Y-%m-%d)

# Collaboration
cicada share create s3://shared/project --grant user:read-write

# Quality check
cicada metadata list ~/data --filter "qc_status=warning"

# Archive
cicada sync ~/old-data s3://archive --storage-class GLACIER

# Restore
cicada s3 restore s3://archive/project --days 7

# Generate report
cicada metadata export ~/data --format csv > report.csv
```

### Workflow Templates

See the `examples/workflows/` directory for complete workflow templates:
- `microscopy-daily.sh`
- `sequencing-rnaseq.sh`
- `mass-spec-batch.sh`
- `collaborative-analysis.sh`
- `backup-strategy.sh`

---

**Related Documentation:**
- [Getting Started](GETTING_STARTED.md)
- [Troubleshooting](TROUBLESHOOTING.md)
- [Advanced Topics](ADVANCED.md)
- [Integration Guide](INTEGRATIONS.md)
