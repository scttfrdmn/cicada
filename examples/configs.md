# Example Cicada Configuration Files

This directory contains example configuration files for various Cicada components.

## Main Configuration (~/.cicada/config.yaml)

```yaml
# config.yaml
version: "1.0"

# Lab information
lab:
  name: smith-lab
  institution: University Example
  bucket: smith-lab-data-20241122
  region: us-west-2

# Storage configuration
storage:
  intelligent_tiering: true
  versioning: true
  lifecycle:
    - rule: archive-old-data
      days_to_glacier: 365
      enabled: true
    - rule: delete-temp-files
      prefix: scratch/
      days_to_expiration: 7
      enabled: true

# Sync defaults
sync:
  default_concurrency: 10
  bandwidth_limit: 100MB  # megabytes per second
  checksum_verify: true
  exclude_patterns:
    - "*.tmp"
    - "*.swp"
    - ".DS_Store"
    - "Thumbs.db"
    - "desktop.ini"

# Compute configuration
compute:
  spot_preferred: true
  max_vcpus: 256
  min_vcpus: 0  # Scale to zero when idle
  instance_types:
    - c5.xlarge
    - c5.2xlarge
    - c5.4xlarge
  timeout_hours: 24

# Workstation defaults
workstation:
  default_instance: g4dn.xlarge
  auto_shutdown_hours: 2
  allow_spot: true
  default_image: basic-linux

# Notifications
notifications:
  email: pi@university.edu
  slack_webhook: https://hooks.slack.com/services/XXX/YYY/ZZZ
  notify_on:
    - sync_complete
    - workflow_complete
    - workstation_shutdown
    - cost_alert
    - error

# Cost management
cost:
  budget_limit: 100  # USD per month
  alert_threshold: 80  # Percent of budget
  currency: USD
  daily_reports: false
  monthly_reports: true

# Compliance
compliance:
  mode: none  # Options: none, nist-800-171, hipaa, gdpr
  audit_log_retention_days: 365
  require_mfa: false
  session_timeout_hours: 8

# Advanced
advanced:
  log_level: info  # debug, info, warn, error
  log_file: ~/.cicada/logs/cicada.log
  log_rotation: daily
  telemetry: true  # Anonymous usage statistics
  auto_update_check: true
```

---

## Workflow Configuration (cicada-workflow.yaml)

```yaml
# cicada-workflow.yaml - Cell segmentation pipeline
name: cell-segmentation-pipeline
version: "1.0"
engine: snakemake
workflow: Snakefile

# Compute resources
compute:
  type: batch
  instance_types:
    - c5.4xlarge
    - c5.9xlarge
  spot: true
  max_vcpus: 256
  min_vcpus: 0

# Storage locations
storage:
  input: s3://smith-lab-data/raw/microscopy/
  output: s3://smith-lab-data/processed/
  working: s3://smith-lab-data/scratch/
  working_retention_days: 7  # Auto-delete after 7 days

# Workflow parameters
parameters:
  min_cell_size: 100
  max_cell_size: 5000
  threshold_method: otsu

# Notifications
notifications:
  email: researcher@university.edu
  slack: "#lab-notifications"
  on_complete: true
  on_failure: true

# Cost controls
cost:
  max_cost: 50  # USD
  abort_on_exceed: true

# Environment
environment:
  docker_image: smithlab/cellseg:latest
  env_vars:
    NUMBA_NUM_THREADS: 8
```

---

## Auto-Pipeline Configuration (auto-pipeline.yaml)

```yaml
# auto-pipeline.yaml - Automatic processing of new sequencing data
name: auto-fastq-processing
version: "1.0"

# Trigger configuration
trigger:
  watch: /data/sequencer/BaseCalls/
  pattern: "*.fastq.gz"
  min_files: 2  # Wait for R1 and R2
  stable_duration: 5m  # Wait 5 minutes after last write

# Workflow
workflow:
  engine: nextflow
  pipeline: nf-core/rnaseq
  version: 3.12.0
  profile: aws_batch

# Parameters
params:
  genome: GRCh38
  outdir: s3://smith-lab-data/results/${RUN_ID}/
  email: bioinformatics@lab.edu

# Compute
compute:
  spot: true
  max_cpus: 64
  max_memory: 128GB

# Notifications
notifications:
  on_complete:
    email: bioinformatics@lab.edu
    message: "Sequencing run ${RUN_ID} complete. Results: ${OUTDIR}"
  on_failure:
    email: bioinformatics@lab.edu
    slack: "#bioinformatics-alerts"
```

---

## Metadata Schema (fluorescence-microscopy.yaml)

```yaml
# fluorescence-microscopy.yaml
schema_version: "1.0"
name: fluorescence-microscopy
description: Metadata for fluorescence microscopy experiments
domain: imaging

extends:
  - cicada://core/experiment

required_fields:
  - sample_id
  - experiment_date
  - operator
  - microscope_type
  - objective_magnification

# Instrument configuration
instrument:
  microscope_type:
    type: string
    required: true
    vocabulary:
      - Wide-field
      - Confocal
      - Two-photon
      - Super-resolution
  
  model:
    type: string
    examples: [LSM 980, A1R HD25, SP8]
  
  objective:
    magnification:
      type: number
      required: true
      vocabulary: [10, 20, 40, 63, 100]
    
    numerical_aperture:
      type: number
      range: [0.1, 1.5]

# Acquisition parameters
acquisition:
  dimensions:
    x_pixels:
      type: integer
      required: true
    y_pixels:
      type: integer
      required: true
    z_slices:
      type: integer
      default: 1
    time_points:
      type: integer
      default: 1

# Channels
channels:
  type: array
  required: true
  items:
    type: object
    fields:
      name:
        type: string
      fluorophore:
        type: string
      excitation_wavelength:
        type: number
        units: nanometers
      emission_wavelength:
        type: number
        units: nanometers

# Sample information
sample:
  organism:
    scientific_name:
      type: string
      examples: [Homo sapiens, Mus musculus]
    ncbi_taxon:
      type: string
      pattern: "^NCBITaxon:[0-9]+$"
  
  cell_line:
    type: string
    examples: [HeLa, HEK293, U2OS]

# Experimental conditions
conditions:
  temperature_celsius:
    type: number
    default: 22
  
  treatment:
    type: array
    items:
      type: object
      fields:
        agent: string
        concentration:
          value: number
          units: string
        duration:
          value: number
          units: string

# Validation rules
validation:
  - rule: if z_slices > 1 then acquisition_mode contains "z-stack"
    message: "Z-stack mode required for multiple z-slices"
```

---

## Watch Configuration (.cicada/watches.yaml)

```yaml
# watches.yaml - Persistent watch configurations
watches:
  - name: microscope-1
    path: /Volumes/ZeissMicroscope/Export
    destination: s3://smith-lab-data/raw/microscopy/
    enabled: true
    
    # Trigger settings
    trigger:
      on_new_file: true
      min_age: 5m  # Wait 5 minutes after file creation
      schedule: null  # No scheduled sync
    
    # Sync options
    sync:
      delete_source: true
      checksum_verify: true
      concurrency: 5
    
    # Metadata
    metadata:
      schema: fluorescence-microscopy
      auto_extract: true
      inherit_from_path:
        operator: from_username
        experiment_date: from_file_mtime
  
  - name: sequencer
    path: /data/sequencer/output
    destination: s3://smith-lab-data/raw/sequencing/
    enabled: true
    
    trigger:
      on_new_file: false
      schedule: "0 2 * * *"  # Daily at 2 AM
    
    sync:
      delete_source: false
      include_patterns:
        - "*.fastq.gz"
        - "*.bam"
      exclude_patterns:
        - "*_failed_*"
    
    metadata:
      schema: rnaseq-experiment
      template: rnaseq-standard
```

---

## Project Configuration (project-NIH-R01.yaml)

```yaml
# project-NIH-R01.yaml
name: NIH-R01-2024
description: Protein structure determination study
created: 2024-01-15
status: active

# Members
members:
  - email: pi@university.edu
    role: admin
    added: 2024-01-15
  
  - email: jsmith@university.edu
    role: read-write
    added: 2024-01-20
    groups: [protein-structure, methods-dev]
  
  - email: external@otheruniv.edu
    role: read-only
    added: 2024-03-01
    expiration: 2025-01-01

# Storage
storage:
  paths:
    - s3://smith-lab-data/projects/NIH-R01-2024/
  quota: 5TB
  warn_at: 4TB

# Cost
budget:
  monthly_limit: 200
  alert_threshold: 80
  cost_center: NIH-R01-12345

# Metadata
metadata:
  grant_number: R01-GM123456
  pi: Dr. Jane Smith
  institution: University Example
  start_date: 2024-01-01
  end_date: 2027-12-31
  keywords:
    - protein structure
    - X-ray crystallography
    - biophysics
```

---

## DOI Configuration (doi-config.yaml)

```yaml
# doi-config.yaml
provider: datacite
test_mode: false

# DataCite credentials
datacite:
  repository_id: EXAMPLE.SMITHLAB
  password: ${DATACITE_PASSWORD}  # From environment variable
  prefix: 10.12345
  cost_per_doi: 1.00

# Default metadata
defaults:
  publisher: Smith Lab Data Repository
  resource_type: Dataset
  language: en
  
  authors:
    - name: Smith, Jane
      orcid: 0000-0001-2345-6789
      affiliation: University Example
  
  funding:
    - funder: National Institutes of Health
      funder_id: "10.13039/100000002"
      award: R01-GM123456

# Landing page template
landing_page:
  template: default
  custom_css: /path/to/custom.css
  logo: /path/to/lab-logo.png
  footer: |
    Data published by Smith Lab at University Example.
    Questions? Contact: data@smithlab.org
```

---

## Compliance Configuration (compliance-nist.yaml)

```yaml
# compliance-nist.yaml - NIST 800-171 configuration
compliance_mode: nist-800-171
version: rev2

# Access control
access_control:
  require_mfa: true
  session_timeout_hours: 1
  max_failed_logins: 3
  password_min_length: 12
  password_complexity: true

# Audit logging
audit:
  enabled: true
  retention_days: 365
  log_all_access: true
  log_all_modifications: true
  immutable: true  # Write to write-once storage
  export_to:
    - s3://audit-logs-bucket/
    - syslog://siem.example.com

# Encryption
encryption:
  at_rest:
    enabled: true
    kms_key_id: arn:aws:kms:us-west-2:123456789012:key/abc123
  
  in_transit:
    tls_min_version: "1.3"
    require_tls: true

# Data classification
classification:
  default_level: CUI
  auto_detect_sensitive: true
  sensitive_patterns:
    - ssn
    - credit_card
    - dob
    - patient_id

# Network controls
network:
  vpc_only: true
  no_internet_access: true
  vpc_endpoints: true
  security_groups:
    - sg-abc123

# Incident response
incident_response:
  notify:
    - security@university.edu
  
  auto_actions:
    - suspend_user_on_anomaly
    - alert_on_bulk_download
```

---

## User Preferences (~/.cicada/preferences.yaml)

```yaml
# preferences.yaml - User-specific settings
user:
  email: jsmith@university.edu
  name: Jane Smith
  orcid: 0000-0001-2345-6789

# UI preferences
ui:
  theme: auto  # light, dark, auto
  date_format: ISO8601
  file_size_units: binary  # binary (GiB) or decimal (GB)
  timezone: America/Los_Angeles

# Default metadata
metadata_defaults:
  operator: Jane Smith
  institution: University Example
  lab: Smith Lab

# Favorite locations
favorites:
  - name: My Experiments
    path: s3://smith-lab-data/users/jsmith/
  
  - name: Shared Data
    path: s3://smith-lab-data/shared/
  
  - name: Current Project
    path: s3://smith-lab-data/projects/NIH-R01-2024/

# CLI preferences
cli:
  color: auto
  progress_bars: true
  confirm_destructive: true
  editor: vim

# Shortcuts
shortcuts:
  sync_my_data: cicada sync ~/data/ s3://smith-lab-data/users/jsmith/
  launch_jupyter: cicada workstation launch --image jupyter
```

---

## Environment Variables (.env)

```bash
# .env - Environment variables for sensitive data
# DO NOT commit this file to git!

# AWS credentials (if not using AWS CLI config)
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...
AWS_DEFAULT_REGION=us-west-2

# DataCite credentials
DATACITE_PASSWORD=...

# Slack webhook
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...

# Globus Auth (if using)
GLOBUS_CLIENT_ID=...
GLOBUS_CLIENT_SECRET=...

# Optional: Encryption keys
ENCRYPTION_KEY=...

# Optional: Database connection (if using remote DB)
DB_CONNECTION_STRING=...
```

---

## Deployment Configuration (docker-compose.yaml)

```yaml
# docker-compose.yaml - For running Cicada as a service
version: '3.8'

services:
  cicada-daemon:
    image: cicada:latest
    container_name: cicada-daemon
    restart: unless-stopped
    
    volumes:
      - ~/.cicada:/root/.cicada
      - /data:/data  # Mount local data directories
    
    environment:
      - AWS_REGION=us-west-2
      - CICADA_WEB_PORT=7878
    
    env_file:
      - .env
    
    ports:
      - "7878:7878"  # Web UI
    
    networks:
      - cicada-net
    
    healthcheck:
      test: ["CMD", "cicada", "status"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  cicada-net:
    driver: bridge
```

---

## Systemd Service (cicada.service)

```ini
[Unit]
Description=Cicada Data Management Daemon
After=network.target

[Service]
Type=simple
User=cicada
Group=cicada
WorkingDirectory=/opt/cicada
ExecStart=/usr/local/bin/cicada daemon start --config /etc/cicada/config.yaml
ExecStop=/usr/local/bin/cicada daemon stop
Restart=on-failure
RestartSec=10

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/cicada /var/log/cicada

[Install]
WantedBy=multi-user.target
```

---

These configuration files provide a complete set of examples for all major Cicada components. Users can copy and customize them for their specific needs.
