# Cicada Project Summary

**Version**: 1.0
**Date**: 2025-11-22
**Status**: Pre-implementation - Specification Complete
**Author**: Scott Friedman
**License**: Apache 2.0
**Copyright**: © 2025 Scott Friedman

---

## Table of Contents

1. [Overview](#overview)
2. [Target Users & Problem Space](#target-users--problem-space)
3. [Core Design Principles](#core-design-principles)
4. [Key Features](#key-features)
5. [Technical Architecture](#technical-architecture)
6. [Technology Stack](#technology-stack)
7. [Document Structure](#document-structure)
8. [Development Plan](#development-plan)
9. [Key Innovations](#key-innovations)
10. [Cost Model](#cost-model)

---

## Overview

**Cicada** is a comprehensive specification for a "dormant data commons" platform designed specifically for academic research labs. The name reflects its core design philosophy - like a cicada, it lies dormant most of the time (consuming minimal resources/costs), but emerges powerfully when needed.

### What is Cicada?

Cicada is a lightweight, cost-effective data commons platform that provides:
- **Federated storage**: Everyone's data in one logical place
- **Access control**: Fine-grained permissions for users, groups, and projects
- **Compute-to-data**: Bring analysis to the data, not data to analysis
- **Collaboration primitives**: Shared workspaces, reproducible workflows
- **Data publication**: DOI minting and public data portals

Unlike traditional data commons platforms (which are heavyweight, require IT staff, and cost $$$), Cicada is:
- **Self-service**: PI installs it, it just works
- **Cost-aware**: Designed for shoestring budgets
- **Opportunistic**: Uses resources only when needed

---

## Data Organization Philosophy

Cicada provides **guided flexibility** for data organization:

### Default Recommended Structure

Cicada suggests a standard structure during setup but does not enforce it:

```
s3://lab-bucket/
├── shared/                    # Lab-wide accessible data
│   ├── protocols/            # SOPs, methods
│   ├── reagents/             # Reagent info, vendors
│   └── instruments/          # Instrument manuals, configs
│
├── projects/                 # Project-based organization
│   ├── grant-nih-r01-2024/
│   │   ├── raw/             # Original instrument data
│   │   ├── processed/       # Analysis outputs
│   │   ├── results/         # Final figures, tables
│   │   └── metadata/        # Project-specific metadata
│   └── paper-nature-2024/
│
├── groups/                   # Research group/theme data
│   ├── protein-structure/
│   └── metabolism/
│
└── users/                    # Personal scratch/work areas
    ├── jsmith/
    └── agarcia/
```

### Flexibility Options

**Researchers can**:
- ✅ Use their own folder structure
- ✅ Sync to any S3 path they choose
- ✅ Organize by date, experiment, instrument, or any scheme
- ✅ Create custom project structures
- ✅ Mix and match organizational approaches

**Cicada provides tools for any structure**:
- Metadata tagging (organizes logically, regardless of folder structure)
- Search across all data (find by metadata, not just folder names)
- Virtual collections (group files from different locations)
- Flexible access controls (can apply to any path pattern)

**Example alternative structures**:
```
# Chronological
/2024/11/experiment-001/
/2024/11/experiment-002/

# By instrument
/zeiss-microscope/2024-11-22/
/sequencer/run-456/

# By student
/jane-thesis/microscopy/
/john-project/proteomics/

# Completely flat
/all-data/
```

### Best Practices (Optional)

Cicada suggests but doesn't require:
1. Separate `raw/` from `processed/` data
2. Use project-based top-level organization for grants
3. Store shared resources in a common location
4. Use consistent date formats (YYYY-MM-DD)
5. Include README files in project folders

**Key Insight**: Metadata and search matter more than folder structure. As long as files have good metadata, researchers can organize however makes sense to them.

---

## Target Users & Problem Space

### Primary Persona

**Dr. Sarah Chen** - Assistant Professor, Cell Biology
- Lab of 8 people (1 postdoc, 5 grad students, 2 undergrads)
- Generates ~500GB/month (microscopy, some sequencing)
- Budget: $50-75/month for cloud storage
- Technical skill: Can use command line, but prefers GUI
- Pain points:
  - Data on dying hard drives
  - No organized backup system
  - Students leaving with data on their laptops
  - Can't share data easily with collaborators
  - Journal requires data availability statements
  - Manual workflows for moving data between instruments, analysis platforms, and storage

### Current State Problems

Small academic labs (8-10 people) typically struggle with:
- **Limited technical expertise**: Not computer scientists
- **Tight budgets**: $50-100/month maximum
- **Minimal IT support**: No dedicated staff
- **Scattered data**: USB drives, laptops, aging NAS devices
- **Poor backup strategy**: Never or rarely, no schedule
- **Manual workflows**: Tedious data movement between:
  - Instruments with limited storage
  - Laptops for temporary analysis
  - SaaS platforms with limited storage
  - Final storage locations (USB drives, NAS)

### What Cicada Solves

- ✅ **Automatic backup**: $80/month, cheaper than replacing drives
- ✅ **Organized metadata**: Find experiments from 2 years ago
- ✅ **Team collaboration**: Everyone has appropriate access
- ✅ **Easy sharing**: Generate DOI, public links for publications
- ✅ **Grant compliance**: FAIR principles, data sharing requirements
- ✅ **Streamlined workflows**: Automated data movement and processing

---

## Core Design Principles

### 1. Dormant by Design 🦗

All resources spin up on-demand and shut down automatically when idle:
- File gateways run only during bulk transfers (~$0.10/hour)
- Workstations auto-shutdown after 2 hours of inactivity
- Compute jobs use spot instances (70% cheaper)
- No "always-on" infrastructure except S3 storage

### 2. Cost-Conscious 💰

Architected for ~$50-100/month budgets:
- S3 Intelligent-Tiering (automatic cost optimization)
- Spot instances for compute (70% cheaper)
- On-demand resources (only pay when using)
- Aggressive lifecycle policies (old data → cheaper storage tiers)
- Cost tracking and budget alerts built-in

### 3. Zero AWS Knowledge Required 🎯

Abstracts all cloud complexity:
- Users never see "S3 bucket" or "object storage"
- It's just "the lab drive" or "lab data"
- CLI and Web UI use familiar concepts (folders, files, projects)
- Automatic IAM policy management
- Setup wizard handles AWS configuration

### 4. FAIR by Default 📊

Implements research data standards:
- **F**indable: Metadata, search, DOIs
- **A**ccessible: Access controls, public portal
- **I**nteroperable: Standard formats, ontology integration
- **R**eusable: Provenance tracking, workflow capture

### 5. Domain-Flexible 🔬

Supports diverse research fields:
- Custom metadata schemas (YAML-based)
- Pluggable instrument adapters
- Domain-specific file format extractors
- Community-contributed schemas and workflows

### 6. Open & Community-Driven 🌐

Open source from day one:
- MIT License for maximum compatibility
- Community contributions encouraged
- Domain-specific schemas shared
- Example workflows and configurations

---

## Key Features

### 1. Intelligent Data Management

#### Sync Engine
- **rsync-like functionality**: Delta detection, checksum-based comparison
- **Multipart uploads**: Resume capability for large files
- **Parallel transfers**: Configurable concurrency
- **Bandwidth throttling**: Don't saturate lab network
- **Deduplication**: Detect identical files, skip re-upload
- **Progress reporting**: Bytes, files, ETA with visual progress bars

#### File Watching
- **Automatic sync**: Monitor instrument folders, upload new files
- **Debouncing**: Wait for writes to complete before syncing
- **Age-based filtering**: Don't sync files that are too new
- **Pattern-based filtering**: Ignore .tmp, .partial files
- **Scheduled syncs**: Cron-like scheduling for batch operations
- **Delete synchronization**: Optionally delete local files after successful upload

#### Storage Management
- **Intelligent-Tiering**: Automatic cost optimization
- **Versioning**: Protect against accidental deletions
- **Lifecycle policies**: Auto-archive old data to Glacier
- **Compression**: Automatic for compatible formats

### 2. Metadata & FAIR Compliance

#### Schema System
- **YAML-based schemas**: Human-readable, version-controlled
- **Schema inheritance**: Extend base schemas for specific needs
- **Field types**: string, number, integer, boolean, array, object, date, datetime
- **Validation**: Required fields, controlled vocabularies, regex patterns
- **Ontology integration**: Map to standard ontologies (EDAM, OBI, etc.)
- **Quality scoring**: 0-100 score for metadata completeness/richness

#### Automatic Extraction
Format-specific extractors for:
- **Microscopy**: TIFF, OME-TIFF, Zeiss CZI, Nikon ND2, Leica LIF
- **Sequencing**: FASTQ, BAM, SAM
- **Mass Spectrometry**: mzML, MGF
- **Medical Imaging**: DICOM
- **Flow Cytometry**: FCS
- **Generic**: HDF5, Zarr, NetCDF
- **Fallback**: Basic file metadata (size, dates, checksums)

#### Export Formats
- **DataCite**: For DOI registration
- **ISA-Tab**: For submission to repositories
- **JSON-LD**: For semantic web compatibility
- **BagIt**: For preservation
- **Custom**: User-defined export templates

### 3. Workflow Execution

#### Workflow Engines
- **Snakemake**: Python-based workflows
- **Nextflow**: JVM-based, popular in bioinformatics
- **CWL**: Common Workflow Language (portable)
- **Custom scripts**: Python, R, shell scripts

#### AWS Batch Integration
- **Spot instances**: 70% cost savings
- **Auto-scaling**: 0 to 256+ vCPUs as needed
- **Job dependencies**: DAG-based execution
- **Retry logic**: Automatic retry on spot interruption
- **Cost limits**: Abort if exceeding budget

#### Features
- **Local testing**: Test workflows locally before cloud execution
- **Environment capture**: Reproducible with Docker containers
- **Progress monitoring**: Real-time logs and status
- **Result storage**: Outputs back to S3 automatically
- **Notifications**: Email/Slack on completion or failure

### 4. Remote Workstations

#### AWS NICE DCV Integration
- **High-performance remote desktop**: Hardware-accelerated graphics
- **H.264 streaming**: Efficient bandwidth usage
- **Browser-based access**: No client software required
- **GPU support**: NVIDIA GPUs for visualization and compute
- **Multi-user sessions**: Collaborative workspaces
- **Audio support**: For video playback and analysis

#### Instance Types
- **g4dn.xlarge**: NVIDIA T4 GPU (~$0.50/hour spot)
- **g4dn.2xlarge**: More GPU memory for large datasets
- **r5.2xlarge**: High memory for intensive analysis
- **c5.4xlarge**: Compute-optimized for CPU-intensive work

#### Pre-built Images
- **basic-linux**: Ubuntu 22.04 with essential tools
- **imagej**: ImageJ/FIJI with common plugins
- **matlab**: MATLAB with common toolboxes (requires license)
- **paraview**: ParaView for scientific visualization
- **rstudio**: RStudio Server with tidyverse
- **jupyter**: JupyterLab with scipy stack
- **napari**: Python visualization for multi-dimensional images
- **custom**: Build your own from Dockerfile

#### Management Features
- **Auto-shutdown**: After configurable idle time (default: 2 hours)
- **Session persistence**: Stop/start sessions, data retained on EBS
- **Snapshots**: Save session state for later restoration
- **Cost tracking**: Real-time cost display
- **Extend timer**: Add time before auto-shutdown

### 5. Collaboration & Sharing

#### User Management
- **Roles**: PI (admin), postdoc, grad student, undergrad, collaborator
- **Groups**: Organize by research theme (#protein-structure, #metabolism)
- **Projects**: Grant-specific or paper-specific data organization
- **External collaborators**: Read-only or time-limited access

#### IAM Automation
- **Automatic policy creation**: Cicada generates IAM policies
- **Least privilege**: Users only access what they need
- **Path-based permissions**: Grant access to specific S3 prefixes
- **Credential management**: Stored securely in OS keychain

#### Authentication Options
1. **Cicada-managed IAM** (default): Automatic user creation
2. **Globus Auth**: Institutional SSO via OAuth2
3. **AWS IAM Identity Center**: For institutions with existing AWS SSO
4. **Bring-your-own**: Use existing AWS accounts

#### Public Data Portal
- **Static site generation**: Fast, cheap (CloudFront CDN)
- **Search and browse**: Faceted search by metadata fields
- **Download links**: Direct S3 presigned URLs
- **DOI landing pages**: Rich metadata display
- **Usage analytics**: Track dataset access and downloads
- **Citation export**: BibTeX, RIS formats

#### DOI Minting
- **DataCite**: Primary DOI provider (institutional memberships common)
- **Zenodo**: Alternative provider (free for researchers, powered by DataCite)
- **Automatic metadata**: Converts Cicada metadata to DOI schema
- **Landing pages**: Rich dataset descriptions with download links
- **Versioning**: Support for dataset versions with related DOIs

### 6. Compliance & Security

#### NIST 800-171 Mode (Controlled Unclassified Information)

**Required for**:
- Government contracts and defense research (CUI)
- **NIH controlled-access genomic data** (effective January 25, 2025)
- Federal research requiring CUI protection

**NIH Mandate**: As of January 25, 2025, all users and developers of [NIH controlled-access genomic data must comply with NIST 800-171](https://grants.nih.gov/grants/guide/notice-files/NOT-OD-24-157.html) cybersecurity requirements. Institutions must attest compliance when requesting access to NIH genomic data repositories.

**Implementation**:
- **Access controls**: Multi-factor authentication (MFA) required
- **Audit logging**: All data access logged to CloudWatch Logs
- **Encryption**: At rest (S3-SSE, EBS encryption) and in transit (TLS 1.2+)
- **Incident response**: Automated alerts on suspicious activity
- **Configuration management**: Enforce security baselines via AWS Config
- **Media protection**: Secure data sanitization on deletion
- **Physical protection**: AWS data center controls (inherited)
- **System monitoring**: Continuous monitoring with CloudWatch and GuardDuty
- **110 security requirements**: Full implementation of NIST SP 800-171 Rev. 3

#### NIST 800-53 Mode (HIPAA/FISMA Compliance)
For labs handling PHI (Protected Health Information) or requiring FISMA compliance:
- **HIPAA-eligible services only**: Enforced service allowlist
  - ✅ S3, EBS, EC2, Lambda, CloudWatch Logs, KMS, IAM, Secrets Manager
  - ❌ Non-HIPAA services blocked by policy
- **Business Associate Agreement (BAA)**: Required with AWS
- **PHI encryption**:
  - AWS KMS with customer-managed keys (CMK)
  - Encryption at rest for all storage (S3, EBS, RDS if used)
  - TLS 1.2+ for all data in transit
- **Access controls** (NIST 800-53 AC family):
  - Multi-factor authentication mandatory
  - Role-based access control (RBAC)
  - Least privilege principle enforced
  - Session timeouts and re-authentication
- **Audit and accountability** (AU family):
  - Comprehensive audit logging to CloudWatch Logs
  - Log retention: minimum 6 years (configurable)
  - Audit log encryption and integrity protection
  - Automated audit review and reporting
- **Identification and authentication** (IA family):
  - Strong password policies
  - MFA for all users
  - Session management and timeout
- **System and communications protection** (SC family):
  - Network segmentation with VPC
  - Security groups and NACLs
  - Boundary protection with AWS WAF (optional)
- **Risk assessment** (RA family):
  - Vulnerability scanning with Inspector
  - Threat detection with GuardDuty
  - Automated compliance checking
- **Incident response** (IR family):
  - Automated security incident detection
  - Notification workflows
  - Incident tracking and remediation
- **Data retention and disposal**:
  - Configurable retention policies
  - Secure deletion with verification
  - Media sanitization procedures
- **De-identification tools**:
  - Automated PHI detection and redaction
  - Safe harbor and expert determination methods
  - Audit trail of de-identification operations

#### Research Security Program Compliance (OSTP Guidelines)

**⚠️ Important**: OSTP research security requirements are **separate and independent** from the NIH NIST 800-171 requirement described above. These are two distinct federal mandates with different frameworks and timelines.

**OSTP July 2024 Directive**: The White House Office of Science and Technology Policy [released final guidelines](https://bidenwhitehouse.archives.gov/ostp/news-updates/2024/07/09/white-house-office-of-science-and-technology-policy-releases-guidelines-for-research-security-programs-at-covered-institutions/) on July 9, 2024, requiring covered research institutions to establish formal research security programs with four required elements:

1. **Cybersecurity** (aligned with [NIST IR 8481](https://csrc.nist.gov/publications/detail/nistir/8481/final))
2. **Foreign travel security**
3. **Research security training**
4. **Export control training**

**Cascading Timeline** (often misunderstood):

The OSTP directive creates a **three-tier cascading timeline**:

1. **Federal Agencies** (~6 months from July 9, 2024):
   - Deadline: **~January 2025**
   - Requirement: Agencies must provide OSTP/OMB their plans for updating policies

2. **Updated Agency Policies Take Effect** (~6 months after plans submitted):
   - Deadline: **~July 2025**
   - Requirement: Agencies must implement updated policies incorporating OSTP guidelines

3. **Covered Institutions Must Comply** (up to 18 months after policies take effect):
   - Deadline: **~December 2026 / January 2027**
   - Requirement: Institutions must have compliant research security programs

**Total Timeline**: Approximately **30 months** from OSTP guideline release (July 2024) to institutional compliance deadline (end of 2026).

**Why Cascading?**: This gives institutions time to prepare while agencies update their grant requirements. Institutions are not expected to comply until after agencies have finalized their policies.

**Cicada's Cybersecurity Support**:

Cicada helps institutions meet the [OSTP cybersecurity requirements](https://www.ropesgray.com/en/insights/alerts/2024/07/final-issuance-of-federal-guidelines-for-security-in-scientific-research-impact-on-universities) by providing:

- ✅ **Access control**: MFA, RBAC, session management
- ✅ **Audit logging**: Comprehensive activity tracking
- ✅ **Data protection**: Encryption at rest and in transit
- ✅ **Incident detection**: Automated threat monitoring (GuardDuty)
- ✅ **Configuration management**: AWS Config rules and compliance checks
- ✅ **Vulnerability management**: Automated scanning (Inspector)
- ✅ **Compliance reporting**: Generate reports for institutional security offices

**Note**: The OSTP guidelines reference [NIST IR 8481](https://www.nist.gov/blogs/manufacturing-innovation-blog/what-nist-sp-800-171-and-who-needs-follow-it-0) (research cybersecurity findings), which provides recommendations but does not mandate a specific framework like NIST 800-171. Cicada's security features align with NIST IR 8481 recommendations while providing optional NIST 800-171 compliance for NIH genomic data requirements.

**For detailed control implementation**: See [COMPLIANCE-CROSSWALK.md](COMPLIANCE-CROSSWALK.md) for a complete mapping of who implements what controls (Cicada, AWS, or Institution).

#### Federal Compliance Landscape Summary

**Two Separate Federal Mandates** (not to be confused):

| Requirement | Source | Effective Date | Applies To | Framework | Cicada Support |
|------------|--------|----------------|------------|-----------|----------------|
| **NIH Genomic Data** | NIH NOT-OD-24-157 | Jan 25, 2025 | NIH controlled-access genomic data users | NIST SP 800-171 (110 controls) | NIST 800-171 mode |
| **Research Security Program** | OSTP Directive | ~End 2026 | All covered research institutions | NIST IR 8481 (recommendations) | Standard mode + security features |

**Key Differences**:
- **NIH requirement**: Specific to genomic data, mandatory NIST 800-171 compliance
- **OSTP requirement**: Broad research security program (4 elements), cybersecurity references NIST IR 8481 (not 800-171)
- **Timelines**: NIH is January 2025 (imminent); OSTP is end of 2026 (cascading)
- **Frameworks**: Different NIST publications (800-171 vs. IR 8481)

#### Data Loss Prevention (DLP)
- **Sensitive data scanning**: Detect PII, credentials, PHI
- **Prevent uploads**: Block files with sensitive data patterns
- **Alerts**: Notify admins of potential data leaks
- **Quarantine**: Move suspicious files to secure location

#### Compliance Mode Selection Guide

| Compliance Mode | Use Case | Key Requirements | BAA Required | Effective Date |
|----------------|----------|------------------|--------------|----------------|
| **Standard** | General research data, OSTP guidelines compliance | Basic security, encryption, MFA | No | Now |
| **NIST 800-171** | Government contracts, CUI, defense research, **NIH genomic data** | 110 security controls, MFA, audit logs, encryption | No | NIH: Jan 25, 2025 |
| **NIST 800-53** | HIPAA (PHI), FISMA, clinical trials, patient data | HIPAA-eligible services only, 6-year logs, CMK encryption | Yes | Now |

**How to Enable**:
```bash
# During setup
cicada init --compliance-mode nist-800-53

# After setup
cicada config set compliance.mode nist-800-53
cicada compliance validate  # Check configuration

# View compliance status
cicada compliance status
# Output:
#   Compliance Mode: NIST 800-53 (HIPAA)
#   BAA Status: ⚠️  User-attested (signed: 2024-01-15)
#                  ⓘ AWS does not provide BAA status via API
#                  ⓘ You must maintain a signed BAA with AWS
#   Encryption: ✓ KMS CMK enabled (key: arn:aws:kms:...)
#   MFA: ✓ Enforced for all users
#   Audit Logs: ✓ Encrypted, 6-year retention
#   HIPAA Services: ✓ Only eligible services in use
#   Non-Eligible Services: ✓ Blocked via SCP/policy
#   Score: 98/100 (2 recommendations)
```

**Important Notes**:
- NIST 800-53 mode requires a signed BAA with AWS before handling PHI
- HIPAA mode blocks non-eligible services (users will receive clear error messages)
- Compliance mode cannot be downgraded (e.g., from NIST 800-53 to Standard) without re-initialization
- All compliance modes include encryption, but NIST 800-53 requires customer-managed keys (CMK)

#### BAA Verification Limitations

**AWS does not provide an API to verify BAA status.** This is an AWS limitation, not a Cicada limitation.

**Cicada's Approach**:

1. **User Attestation** (during NIST 800-53 setup):
   ```bash
   cicada init --compliance-mode nist-800-53

   # Interactive prompts:
   #
   # ⚠️  HIPAA Compliance Requirement
   #
   # To use NIST 800-53 mode for PHI, you MUST:
   # 1. Sign a Business Associate Agreement (BAA) with AWS
   # 2. Download it from AWS Artifact: https://console.aws.amazon.com/artifact
   #
   # Have you signed a BAA with AWS for account 123456789012? [y/N]: y
   # BAA signed date (YYYY-MM-DD): 2024-01-15
   #
   # ✓ Attestation recorded. You are responsible for maintaining
   #   an active BAA with AWS while handling PHI.
   ```

2. **What Cicada DOES Verify Automatically**:
   - ✅ Only HIPAA-eligible services are in use (via AWS Config rules)
   - ✅ KMS customer-managed keys are enabled
   - ✅ All data is encrypted at rest and in transit
   - ✅ Audit logging is configured correctly
   - ✅ MFA is enforced for all users
   - ✅ 6-year log retention is active
   - ✅ Service Control Policies (SCPs) block non-eligible services

3. **Ongoing Reminders**:
   - Monthly reminder: "Ensure your AWS BAA is still active"
   - Annual prompt: "Confirm BAA renewal date"
   - Startup warning if BAA attestation is >1 year old

4. **Documentation Generation**:
   ```bash
   cicada compliance report --format pdf
   # Generates report including:
   # - User-attested BAA date
   # - Technical controls verification
   # - Service usage audit (only eligible services)
   # - Encryption status
   # - Audit log completeness
   ```

**Why This Matters**:
- IT auditors will ask for proof of BAA
- Cicada can prove technical controls are in place
- But ultimate BAA compliance responsibility remains with the customer
- This is the same approach used by HIPAA-compliant SaaS products

**Best Practice**:
- Download BAA from AWS Artifact when enabling NIST 800-53 mode
- Store signed BAA with institutional compliance office
- Set calendar reminder for annual BAA review
- Include BAA in institutional HIPAA documentation

---

## Technical Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    User Interfaces                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐ │
│  │   CLI    │  │  Web UI  │  │  Public Data Portal  │ │
│  └────┬─────┘  └────┬─────┘  └──────────┬───────────┘ │
└───────┼─────────────┼────────────────────┼─────────────┘
        │             │                    │
        └─────────────┼────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────────┐
│                  Cicada Daemon (Local)                  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Core Services                                    │  │
│  │  • Sync Engine        • Metadata Manager          │  │
│  │  • File Watcher       • Schema Validator          │  │
│  │  • Workflow Orchestrator                          │  │
│  │  • Workstation Manager                            │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  HTTP Server (Port 7878)                          │  │
│  │  • REST API           • WebSocket (real-time)     │  │
│  │  • Static Files       • SSE (progress streams)    │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    AWS Infrastructure                   │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐ │
│  │  S3 Storage │  │  AWS Batch   │  │  EC2/Fargate  │ │
│  │  • Intelligent│  │  • Spot      │  │  • Workstations│ │
│  │    Tiering   │  │    Instances │  │  • Gateway    │ │
│  │  • Versioning│  │  • Auto-scale│  │  • On-demand  │ │
│  └─────────────┘  └──────────────┘  └───────────────┘ │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐ │
│  │  IAM/SSO    │  │  CloudWatch  │  │  Lambda       │ │
│  │  • Policies │  │  • Logs      │  │  • Triggers   │ │
│  │  • Roles    │  │  • Alerts    │  │  • Cleanup    │ │
│  └─────────────┘  └──────────────┘  └───────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Component Details

#### Local Daemon (Go)
Runs on researcher's workstation:
- **Background service**: systemd (Linux), launchd (macOS), Windows Service
- **Web server**: Serves UI and API on localhost:7878
- **File watcher**: Monitors instrument folders for new data
- **Sync engine**: Handles uploads/downloads to S3
- **Metadata manager**: Extracts, validates, indexes metadata
- **Workflow orchestrator**: Submits jobs to AWS Batch
- **Workstation manager**: Launches and manages EC2 instances

#### Web UI (AWS Cloudscape)
Modern, accessible interface served locally:
- **Framework**: AWS Cloudscape Design System
  - Professional, consistent UI components
  - Accessibility (WCAG 2.1 AA compliant)
  - Responsive design for desktop and tablet
  - Dark mode support
  - Built-in data tables, forms, modals
- **Build**: Vite for fast development and builds
- **State Management**: Zustand or Pinia (lightweight)
- **Real-time Updates**: WebSocket connection to daemon
- **Features**:
  - Dashboard with activity feed
  - File browser with drag-drop upload
  - Sync manager with live progress
  - Workflow builder and monitor
  - Workstation launcher
  - Cost dashboard with charts
  - Project and user management
  - Settings and preferences

#### AWS Services (On-Demand)

**Storage (Always On, But Cheap)**:
- **S3**: Primary storage with Intelligent-Tiering (HIPAA-eligible ✓)
- **S3 Glacier**: Long-term archival of old data (HIPAA-eligible ✓)
- **EBS**: Block storage for workstations (HIPAA-eligible ✓)
- **DynamoDB**: Metadata index for fast search (HIPAA-eligible ✓, optional)

**Compute (On-Demand)**:
- **AWS Batch**: Workflow execution with spot instances (HIPAA-eligible ✓)
- **EC2**: Workstations with NICE DCV, file gateway (HIPAA-eligible ✓)
- **Fargate**: Serverless container execution (HIPAA-eligible ✓, alternative to Batch)
- **Lambda**: Event triggers, cleanup tasks, API endpoints (HIPAA-eligible ✓)
- **ECR**: Container registry for workflow images (HIPAA-eligible ✓)

**Access & Security**:
- **IAM**: User policies, roles, service accounts (HIPAA-eligible ✓)
- **IAM Identity Center (AWS SSO)**: Enterprise authentication (HIPAA-eligible ✓)
- **Secrets Manager**: Credentials storage (HIPAA-eligible ✓)
- **KMS**: Encryption key management (HIPAA-eligible ✓)
- **AWS Config**: Compliance monitoring (HIPAA-eligible ✓)
- **GuardDuty**: Threat detection (HIPAA-eligible ✓)
- **Inspector**: Vulnerability scanning (HIPAA-eligible ✓)

**Monitoring & Alerts**:
- **CloudWatch Logs**: Audit logging (HIPAA-eligible ✓)
- **CloudWatch Metrics**: Performance monitoring (HIPAA-eligible ✓)
- **Cost Explorer API**: Cost tracking and prediction (HIPAA-eligible ✓)
- **EventBridge**: Scheduled tasks, event routing (HIPAA-eligible ✓)

**Networking**:
- **VPC**: Virtual private cloud (HIPAA-eligible ✓)
- **Security Groups**: Firewall rules (HIPAA-eligible ✓)
- **CloudFront**: CDN for public portal (HIPAA-eligible ✓)

**Deployment**:
- **CloudFormation**: Infrastructure as code (HIPAA-eligible ✓)
- **EC2 Launch Templates**: Pre-configured workstation images (HIPAA-eligible ✓)
- **Systems Manager**: Parameter store, automation (HIPAA-eligible ✓)

**Services Avoided in HIPAA Mode**:
- ❌ **SQS** (not HIPAA-eligible) - use EventBridge or Lambda instead
- ❌ **SNS** (not HIPAA-eligible for PHI) - use SES or in-app notifications
- ❌ **ElastiCache** (not HIPAA-eligible) - use in-memory caching or DynamoDB
- ❌ Any third-party integrations without signed BAAs

---

## Technology Stack

### Backend (Go 1.21+)

**Core Libraries**:
```go
require (
    // CLI framework
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.0

    // AWS SDK
    github.com/aws/aws-sdk-go-v2 v1.24.0
    github.com/aws/aws-sdk-go-v2/service/s3 v1.44.0
    github.com/aws/aws-sdk-go-v2/service/batch v1.30.0
    github.com/aws/aws-sdk-go-v2/service/ec2 v1.141.0

    // File operations
    github.com/fsnotify/fsnotify v1.7.0

    // HTTP server
    github.com/go-chi/chi/v5 v5.0.11
    github.com/gorilla/websocket v1.5.1

    // Database
    github.com/mattn/go-sqlite3 v1.14.19

    // Utilities
    github.com/cheggaaa/pb/v3 v3.1.4          // Progress bars
    github.com/robfig/cron/v3 v3.0.1          // Scheduling
    gopkg.in/yaml.v3 v3.0.1                   // YAML parsing
)
```

**Project Structure**:
```
cicada/
├── cmd/
│   ├── cicada/              # Main CLI entry point
│   ├── cicada-daemon/       # Background daemon service
│   └── cicada-gateway/      # File gateway orchestrator
│
├── internal/
│   ├── sync/               # rsync-like sync engine
│   ├── watch/              # File system watcher
│   ├── storage/            # S3 operations
│   ├── metadata/           # Metadata management
│   ├── workflow/           # Workflow execution
│   ├── workstation/        # Remote desktop management
│   ├── auth/               # Authentication & authorization
│   ├── doi/                # DOI management
│   ├── portal/             # Public data portal
│   ├── compliance/         # NIST 800-171, HIPAA, etc.
│   ├── cost/               # Cost tracking & optimization
│   ├── config/             # Configuration management
│   └── webui/              # Web UI backend
│
├── pkg/
│   ├── instrument/         # Pluggable instrument adapters
│   └── schemas/            # Community metadata schemas
│
├── web/                    # Frontend (Cloudscape)
├── templates/              # CloudFormation templates
├── docs/                   # Documentation
└── examples/               # Example configs and workflows
```

### Frontend (AWS Cloudscape)

**Framework**: AWS Cloudscape Design System
- **Official AWS UI framework**: Consistent with AWS Console experience
- **Components**: 50+ pre-built React components
  - AppLayout (responsive app shell)
  - Table (sortable, filterable, selectable)
  - Form, Input, Select, DatePicker
  - Modal, Alert, Flashbar (notifications)
  - ProgressBar, Spinner
  - Charts (via CloudWatch integration)
  - CodeEditor (for workflow editing)
- **Accessibility**: WCAG 2.1 AA compliant out of the box
- **Themes**: Light and dark mode support
- **TypeScript**: Full type definitions included

**Stack**:
```json
{
  "dependencies": {
    "@cloudscape-design/components": "^3.0.0",
    "@cloudscape-design/global-styles": "^1.0.0",
    "@cloudscape-design/design-tokens": "^3.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "zustand": "^4.4.0",
    "axios": "^1.6.0",
    "recharts": "^2.10.0"
  },
  "devDependencies": {
    "vite": "^5.0.0",
    "typescript": "^5.3.0",
    "@vitejs/plugin-react": "^4.2.0"
  }
}
```

**File Structure**:
```
web/
├── src/
│   ├── components/
│   │   ├── dashboard/      # Dashboard widgets
│   │   ├── files/          # File browser components
│   │   ├── sync/           # Sync manager components
│   │   ├── workflows/      # Workflow UI
│   │   ├── workstations/   # Workstation management
│   │   ├── projects/       # Project management
│   │   ├── users/          # User management
│   │   └── common/         # Shared components
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── FileBrowser.tsx
│   │   ├── SyncManager.tsx
│   │   ├── Workflows.tsx
│   │   ├── Workstations.tsx
│   │   ├── Projects.tsx
│   │   └── Settings.tsx
│   ├── stores/
│   │   ├── appStore.ts
│   │   ├── fileStore.ts
│   │   └── syncStore.ts
│   ├── lib/
│   │   ├── api.ts          # API client
│   │   └── websocket.ts    # WebSocket client
│   ├── App.tsx
│   └── main.tsx
├── public/
└── index.html
```

### AWS Services

**Storage**:
- S3 with Intelligent-Tiering
- S3 Glacier for archival
- EBS for workstation storage

**Compute**:
- AWS Batch (spot instances)
- EC2 (g4dn, r5, c5 families)
- Fargate (optional)

**Remote Desktop**:
- **NICE DCV**: High-performance remote desktop
  - GPU acceleration with NVIDIA drivers
  - H.264 hardware encoding
  - Audio support
  - USB redirection
  - Multi-monitor support
  - Copy/paste between local and remote
  - File transfer via web client

**Networking**:
- VPC with private subnets
- Security Groups
- NAT Gateway (optional)

**Management**:
- IAM
- CloudWatch
- Cost Explorer
- EventBridge
- CloudFormation

---

## Document Structure

The specification is exceptionally well-organized with ~8,300 lines across 8 files:

### Core Documents

1. **cicada-chat.md** (~3,300 lines)
   - Original design conversation
   - Shows evolution of ideas
   - Rationale for design decisions
   - Use case exploration

2. **INDEX.md** (~290 lines)
   - Navigation guide for all specs
   - Document summaries
   - Development workflow guide
   - File statistics

3. **QUICKSTART.md** (~350 lines)
   - Fastest path to understanding
   - Key design decisions
   - Domain examples (12 research fields)
   - Implementation tips
   - Success metrics

4. **README.md** (~370 lines)
   - Project overview
   - System architecture
   - Complete project structure
   - Technology stack
   - Getting started for developers

5. **ROADMAP.md** (~1,800 lines)
   - Detailed 26-week implementation plan
   - 6 development phases
   - Technical specifications per component
   - Test requirements
   - Acceptance criteria

### User Documentation

6. **docs/cli-reference.md** (~1,200 lines)
   - Complete CLI command reference
   - Installation and setup
   - All commands with examples
   - Usage patterns

7. **docs/domain-schemas.md** (~500 lines)
   - 12+ research domain examples
   - Microscopy, genomics, proteomics, etc.
   - Schema structures per domain
   - Metadata field specifications

### Code Specifications

8. **internal/metadata/schema.go** (~700 lines)
   - Complete Go code structures
   - Schema data types
   - Validation engine
   - Quality scoring
   - Manager implementation

9. **internal/metadata/extractor.go** (~900 lines)
   - File format extractors
   - TIFF, OME-TIFF, CZI, ND2, LIF
   - FASTQ, BAM, mzML, FCS
   - HDF5, Zarr, DICOM

10. **internal/doi/datacite.go** (~800 lines)
    - DOI minting implementation
    - DataCite API client
    - Metadata schema structures
    - XML generation

### Configuration Examples

11. **examples/configs.md** (~800 lines)
    - Main config (config.yaml)
    - Workflow configs
    - Metadata schema examples
    - Watch configurations
    - Project configs
    - DOI configs
    - Compliance configs

---

## Development Plan

### 26-Week Roadmap

#### Phase 1: Core Storage & Sync (Weeks 1-6)

**Priorities**: CRITICAL
- S3 sync engine with rsync-like delta detection
- File watching with debouncing and scheduling
- Basic CLI with colored output and progress bars
- Daemon service (systemd, launchd, Windows Service)
- Cost tracking via CloudWatch/Cost Explorer

**Milestone**: Basic data management working
- Can sync data to S3
- Daemon watches folders
- Basic cost tracking

**Acceptance Criteria**:
- Sync 10GB in <15 minutes on 100Mbps
- Handle network interruptions gracefully
- Detect new files within 5 seconds

---

#### Phase 2: Metadata & FAIR (Weeks 7-10)

**Priorities**: CRITICAL
- YAML-based metadata schema system
- Automatic extraction from file formats
- Validation engine with quality scoring
- Search and discovery
- Export formats (DataCite, ISA-Tab)

**Milestone**: FAIR-compliant data management
- Metadata extracted automatically
- Validation with quality scores
- Searchable metadata index

**Domain Schemas to Implement**:
- Microscopy (fluorescence, confocal)
- Sequencing (RNA-seq, DNA-seq)
- Proteomics (mass spec)
- Flow cytometry
- Chromatography

---

#### Phase 3: Web UI & User Management (Weeks 11-14)

**Priorities**: CRITICAL
- Web UI using AWS Cloudscape Design System
- User/group/project management
- IAM policy automation
- Globus Auth integration
- File browser with drag-drop upload

**Milestone**: Complete user experience
- GUI for all core features
- Multi-user collaboration
- Project-based organization

**UI Pages to Build**:
- Dashboard (overview, activity feed)
- File Browser (browse, upload, download)
- Sync Manager (watch configs, progress)
- Project Management (members, permissions)
- User Settings (preferences, credentials)

---

#### Phase 4: Compute & Workflows (Weeks 15-18)

**Priorities**: HIGH
- Workflow execution (Snakemake first)
- AWS Batch integration
- Spot instance management
- Cost-aware execution
- Environment capture (Docker)

**Milestone**: Cloud compute integration
- Run workflows in AWS
- Automatic scaling
- Cost tracking per workflow

**Workflow Engines**:
- Snakemake (priority 1)
- Nextflow (priority 2)
- CWL (priority 3)
- Custom scripts

---

#### Phase 5: Workstations & Portal (Weeks 19-22)

**Priorities**: HIGH
- Remote workstation launcher
- NICE DCV integration
- Pre-built images (ImageJ, ParaView, etc.)
- Public data portal
- DOI minting via DataCite

**Milestone**: Complete feature set
- GPU workstations available
- Public data sharing
- DOI-minted datasets

**Workstation Images**:
- basic-linux (Ubuntu 22.04)
- imagej (FIJI with plugins)
- paraview (3D visualization)
- rstudio (R environment)
- jupyter (Python notebooks)
- napari (multi-dimensional imaging)

---

#### Phase 6: Compliance & Polish (Weeks 23-26)

**Priorities**: MEDIUM
- **NIST 800-171 compliance mode** (CUI/government contracts)
  - 110 security requirements implementation
  - MFA enforcement, audit logging, encryption
  - AWS Config rules for continuous compliance
  - Compliance reporting dashboard
- **NIST 800-53 compliance mode** (HIPAA/FISMA)
  - HIPAA-eligible services enforcement
  - BAA requirement validation
  - KMS with customer-managed keys (CMK)
  - 6-year audit log retention
  - PHI detection and de-identification tools
  - Comprehensive security control implementation
- **Audit logging enhancements**
  - CloudWatch Logs integration
  - Audit log encryption and integrity
  - Automated audit review and alerting
  - Compliance report generation
- **DLP scanning**
  - PII/PHI detection engines
  - Automated redaction tools
  - Upload prevention for sensitive data
- **Documentation**
  - User guides for all features
  - Compliance setup guides
  - API documentation
  - Video tutorials
- **Community schemas**
  - 12+ domain schemas completed
  - Contribution guidelines
  - Schema validation tools

**Milestone**: Production ready
- Compliant with NIST 800-171 and NIST 800-53
- HIPAA-ready with BAA-eligible services only
- Complete documentation
- Example workflows and schemas
- Community contribution pathways

---

### Post-Launch Roadmap

**Months 7-9**: Community Building
- Gather feedback from beta users
- Add domain-specific schemas
- Build instrument adapter library
- Create tutorial videos
- Write academic paper about platform

**Months 10-12**: Advanced Features
- Multi-region support
- Private PyPI/CRAN package hosting
- Notebook hosting (JupyterHub)
- Data streaming analytics
- ML model training integration

---

## Key Innovations

### 1. Cost Model: 10x Cheaper Than Alternatives

**Typical 10TB Lab Cost**: ~$90/month
- S3 Intelligent-Tiering: ~$80/month (avg)
- Compute (spot): ~$3/month (20 hours)
- Transfer: ~$5/month
- File Gateway: ~$0.40/month (4 hours)
- API requests: ~$2/month

**Compared to**:
- Dropbox 10TB: $240/month (2.7x more expensive)
- Traditional data commons (iRODS/Dataverse): $500-1000/month (5-10x more)
- Local NAS: $0/month until it fails and you lose everything

**How We Achieve This**:
- ✅ Intelligent-Tiering automatically moves to cheaper storage tiers
- ✅ Spot instances provide 70% compute savings
- ✅ On-demand infrastructure (no always-on costs)
- ✅ Aggressive lifecycle policies
- ✅ Deduplication prevents redundant uploads
- ✅ Multipart resume prevents re-uploading on failure

### 2. User Experience: Zero AWS Knowledge

**Traditional Approach**:
```bash
# User needs to know:
aws s3 sync /local/data s3://my-bucket/data --storage-class INTELLIGENT_TIERING
aws batch submit-job --job-name analysis --job-definition arn:aws:...
```

**Cicada Approach**:
```bash
# User-friendly commands:
cicada sync /local/data lab-data/raw/
cicada workflow run analysis.smk --input recent-data/
```

**Web UI**: Point-and-click for everything
- Drag files → auto-upload with metadata extraction
- Click "Run Workflow" → AWS Batch job submitted
- Click "Launch Workstation" → EC2 with DCV ready in 3 minutes

### 3. Domain Flexibility: Pluggable Schemas

**Problem**: Every research domain has unique metadata needs
- Microscopy needs: magnification, objective, wavelengths
- Sequencing needs: read length, library prep, organism
- Mass spec needs: ionization method, m/z range, resolution

**Solution**: YAML-based schema system
```yaml
# schemas/fluorescence-microscopy.yaml
name: fluorescence-microscopy
extends: [microscopy-base]
fields:
  excitation_wavelength:
    type: number
    units: nm
    range: {min: 300, max: 800}
    required: true
  emission_wavelength:
    type: number
    units: nm
    range: {min: 300, max: 800}
  objective:
    type: string
    vocabulary: [4x, 10x, 20x, 40x, 63x, 100x]
```

**Community-Driven**:
- Users contribute schemas for their domains
- Share via GitHub repo
- Install: `cicada schema install community/cryo-em`

### 4. Hybrid Access: CLI + GUI

**Power Users** (postdocs, bioinformaticians):
- CLI for automation and scripting
- Integrate into existing workflows
- SSH to servers, use from HPC login nodes

**Non-Technical Users** (PIs, students):
- Web UI for all operations
- Drag-and-drop file upload
- Visual workflow builder
- Point-and-click workstation launcher

**Both**:
- Same underlying API
- Consistent behavior
- Real-time updates across interfaces

### 5. Compute-to-Data: Stop Moving Large Files

**Traditional Workflow**:
1. Download 100GB from instrument → laptop (30 min)
2. Process on laptop (8 hours)
3. Upload results to storage (10 min)
4. Repeat for each experiment

**Cicada Workflow**:
1. Instrument → S3 (automatic, background)
2. Process in AWS Batch (2 hours, $2)
3. Results → S3 (automatic)
4. Download only final results (1GB)

**Benefits**:
- ✅ Faster (data doesn't move, compute does)
- ✅ Cheaper (spot instances)
- ✅ Reproducible (Docker containers)
- ✅ Laptop-free (no blocking researcher's computer)

### 6. Dormant Design: Pay Only When Using

**Example**: Lab generates 50GB/week data
- **Active Research Period** (3 months):
  - Weekly instrument syncs: automatic
  - Monthly workflow runs: $2 each
  - Occasional workstation use: $5/week
  - **Cost**: ~$100/month

- **Quiet Period** (paper writing, conferences):
  - No new data
  - No workflows
  - No workstations
  - **Cost**: ~$80/month (storage only)

**Traditional Data Commons**: Same cost whether using it or not (~$500/month)

---

## Cost Model

### Detailed Breakdown

#### 10TB Dataset Example

**Storage Costs** (varies by data age):
```
Recent data (0-90 days):      500 GB @ $0.023/GB = $11.50/mo
Frequent Access (90-180 days): 1 TB @ $0.019/GB  = $19.00/mo
Infrequent (180-365 days):    3 TB @ $0.013/GB  = $39.00/mo
Archive (1-2 years):          3 TB @ $0.005/GB  = $15.00/mo
Deep Archive (2+ years):      2.5 TB @ $0.002/GB = $5.00/mo
───────────────────────────────────────────────────────────
Total Storage:                                    $89.50/mo
```

**Compute Costs** (bursty usage):
```
Workflow runs:        20 vCPU-hours @ $0.017 spot = $0.34
                      10 GB-hours RAM @ $0.002    = $0.02
Workstation:          8 hours g4dn.xlarge spot   = $4.00
                      @ $0.50/hour
File Gateway:         2 hours m5.large           = $0.20
                      @ $0.10/hour
───────────────────────────────────────────────────────────
Total Compute:                                    $4.56/mo
```

**Data Transfer**:
```
Ingress:              Free (125 GB uploaded)     = $0.00
Egress:               24 GB @ $0.09/GB           = $2.16
───────────────────────────────────────────────────────────
Total Transfer:                                   $2.16/mo
```

**API Requests**:
```
PUT/POST:             100,000 requests @ $0.005  = $0.50
GET/LIST:             500,000 requests @ $0.0004 = $0.20
───────────────────────────────────────────────────────────
Total API:                                        $0.70/mo
```

**Grand Total: ~$97/month**

### Cost Optimization Tips

**Built-in to Cicada**:
1. ✅ Intelligent-Tiering enabled by default
2. ✅ Lifecycle policies (1yr → Glacier, 2yr → Deep Archive)
3. ✅ Deduplication (don't upload same file twice)
4. ✅ Spot instances only (70% cheaper)
5. ✅ Auto-shutdown idle resources
6. ✅ Multipart resume (don't re-upload on failure)
7. ✅ Compression for text/logs

**User Actions**:
1. Archive old projects: `cicada project archive old-data`
2. Delete unneeded files: `cicada files cleanup --older-than 3y`
3. Use smaller instances: Change default workstation size
4. Limit egress: Download only final results, not raw data

### Cost Alerts

**Automatic Alerts**:
- Approaching budget limit (80%, 90%, 100%)
- Unusual activity (10x normal)
- Large downloads (egress costs)
- Idle resources running

**Recommendations**:
- "Archive data older than 2 years → save $15/month"
- "Use c5.xlarge instead of c5.2xlarge → save $3/month"
- "Enable compression on logs → save $2/month"

---

## Supported Research Domains

Cicada includes pre-built metadata schemas and file format extractors for diverse research fields:

### Imaging Sciences
- **Microscopy**: Fluorescence, confocal, super-resolution (STORM, PALM), light sheet
- **Medical Imaging**: CT, MRI, PET, ultrasound (DICOM support)
- **Electron Microscopy**: TEM, SEM, cryo-EM
- **X-ray**: Crystallography, micro-CT

### Omics
- **Genomics**: DNA-seq, RNA-seq, ChIP-seq, ATAC-seq, single-cell
- **Proteomics**: Mass spectrometry (LC-MS, MALDI-TOF)
- **Metabolomics**: NMR, GC-MS
- **Transcriptomics**: Microarray, RNA-seq

### Analytical Chemistry
- **Chromatography**: HPLC, UHPLC, GC, LC-MS
- **Spectroscopy**: NMR, IR, UV-Vis, Raman, fluorescence
- **Mass Spectrometry**: ESI, MALDI, TOF, Orbitrap

### Cell Biology
- **Flow Cytometry**: Multi-parameter cell analysis, FACS
- **Cell Culture**: Live-cell imaging, high-content screening
- **Cell Tracking**: Video microscopy, automated tracking

### Other Sciences
- **Behavioral Studies**: Video recordings, motion tracking, psychometrics
- **Environmental Science**: Field sampling, sensor networks, climate data
- **Materials Science**: AFM, SEM, mechanical testing, spectroscopy
- **Clinical Trials**: Patient data (HIPAA-compliant), electronic health records
- **Neuroscience**: Electrophysiology, calcium imaging, MRI, behavior

Each domain includes:
1. ✅ Metadata schema (YAML)
2. ✅ File format extractors (Go)
3. ✅ Validation rules
4. ✅ Search facets
5. ✅ Export templates
6. ✅ Example datasets
7. ✅ Documentation

---

## Success Metrics

### Technical Metrics

**Performance**:
- Setup time: <10 minutes from install to first sync
- Sync speed: 10GB in <15 minutes (100Mbps connection)
- Metadata extraction: <1 second per file
- Search latency: <200ms for 100K files
- UI responsiveness: <100ms for interactions

**Reliability**:
- Uptime: >99.9% (daemon stability)
- Data durability: 99.999999999% (S3 SLA)
- Resume success rate: >99% (interrupted transfers)
- Spot interruption handling: 100% (automatic retry)

**Quality**:
- Test coverage: >80% (unit + integration)
- Security scan: 0 high/critical vulnerabilities
- Accessibility: WCAG 2.1 AA compliance
- Documentation: 100% of features documented

### User Metrics

**Adoption**:
- Weekly active users: Growing month-over-month
- Data under management: >1 PB across all users
- Projects created: >1000
- DOIs minted: >100 in first year

**Engagement**:
- Daily syncs: >50% of users
- Workflow runs: >100 per week
- Workstation sessions: >50 per week
- Support ticket rate: <5% of users

**Satisfaction**:
- Net Promoter Score (NPS): >50
- User retention: >90% after 3 months
- GitHub stars: >1000 in first year
- Conference presentations: >5 in first year

### Business Metrics (if SaaS)

**Financial**:
- Cost per user (overhead): <$5/month
- Revenue per user: $10-20/month
- Gross margin: >60%
- Customer acquisition cost: <$100

**Growth**:
- User growth: 20% month-over-month
- Churn rate: <5% monthly
- Expansion revenue: 15% from upsells
- Enterprise deals: >5 in first year

---

## Getting Started

### For End Users

**Installation**:
```bash
# macOS
brew install cicada

# Linux
curl -sSL https://install.cicada.sh | sh

# Windows
winget install cicada
```

**First-Time Setup**:
```bash
# Interactive setup wizard
cicada init

# Walks through:
# 1. AWS credentials
# 2. Create S3 bucket
# 3. Configure cost alerts
# 4. Optional features (compute, workstations)
```

**Basic Usage**:
```bash
# Start daemon with web UI
cicada daemon start --web
# Opens browser to http://localhost:7878

# Add watch folder
cicada watch add microscope \
  --path /Volumes/Microscope/Export \
  --sync-on-new

# Manual sync
cicada sync /local/data lab-data/project-x/

# Check status
cicada status

# View costs
cicada cost report
```

### For Developers

**Prerequisites**:
```bash
# Go 1.21+
go version

# Node.js 20+
node --version

# AWS CLI configured
aws sts get-caller-identity

# Docker (for testing)
docker --version
```

**Clone and Build**:
```bash
git clone https://github.com/your-org/cicada.git
cd cicada

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run locally
./bin/cicada daemon start --dev

# Build web UI
cd web
npm install
npm run dev
```

**Development Workflow**:
1. Create feature branch
2. Write tests first (TDD encouraged)
3. Implement feature
4. Run linters: `make lint`
5. Run tests: `make test`
6. Manual testing with `--dev` mode
7. Create PR with description and tests

---

## Contributing

### Open Source Model

**License**: Apache License 2.0

Copyright © 2025 Scott Friedman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

**Why Apache 2.0?**:
- Explicit patent grant protection
- Well-understood by institutions and legal teams
- Compatible with most open source licenses
- Suitable for enterprise and research adoption

**Repository Structure**:
- Main repo: Core Cicada code
- Schema registry: Metadata schemas for research domains
- Instrument adapters: Pluggable instrument support
- Workflow templates: Example Snakemake/Nextflow pipelines

### Community

- **GitHub**: https://github.com/your-org/cicada
- **Discussions**: GitHub Discussions for Q&A
- **Slack**: cicada-data.slack.com for real-time chat
- **Twitter**: @cicada_data for announcements
- **Blog**: Technical deep-dives and user stories
- **Conferences**: FORCE11, RDA, domain-specific meetings

---

## Roadmap Beyond MVP

### Phase 7: Enhanced Collaboration (Months 7-9)

- **Real-time collaboration**: Multiple users in same workstation
- **Notebook hosting**: JupyterHub integration
- **Shared workflows**: Team workflow library
- **Comments & annotations**: On datasets and files
- **Activity feed**: Lab-wide or project-specific

### Phase 8: Advanced Analytics (Months 10-12)

- **Data streaming**: Real-time data ingestion from instruments
- **ML integration**: Model training on AWS SageMaker
- **Automated QC**: Quality control pipelines
- **Dashboard builder**: Custom data visualizations
- **Alerts**: Threshold-based notifications

### Phase 9: Enterprise Features (Year 2)

- **Multi-region**: Replicate data across regions
- **Private package hosting**: PyPI, CRAN mirrors
- **Advanced compliance**: FedRAMP, ISO 27001
- **SSO integrations**: Okta, Azure AD, Google Workspace
- **Multi-tenancy**: Shared infrastructure for institutions
- **White-labeling**: Custom branding for universities
- **Globus integration**: For institutions with Globus subscriptions that include S3 connector
  - Use Globus for large file transfers (optimized for high-throughput)
  - Leverage existing institutional Globus endpoints
  - Integrate Globus Auth for unified authentication
  - Globus sharing and access management alongside Cicada's built-in controls

### Phase 10: Ecosystem (Year 2+)

- **Plugin system**: Third-party integrations
- **Marketplace**: Buy/sell schemas, workflows, adapters
- **Academic partnerships**: Integrate with institutional systems
- **Cloud provider agnostic**: Support Azure, GCP
- **On-premises**: Deploy in university data centers

---

## Comparison to Alternatives

### vs. Dropbox/Google Drive

| Feature | Cicada | Dropbox/Drive |
|---------|--------|---------------|
| Cost (10TB) | $90/month | $240/month |
| Compute integration | ✅ Yes (AWS Batch) | ❌ No |
| Metadata management | ✅ Rich, domain-specific | ❌ Limited to filenames |
| DOI minting | ✅ Yes (DataCite/Zenodo) | ❌ No |
| Compliance modes | ✅ NIST 800-171, NIST 800-53 (HIPAA) | ⚠️ Limited, no HIPAA BAA |
| Workflow execution | ✅ Snakemake, Nextflow, CWL | ❌ No |
| Cost optimization | ✅ Automatic tiering | ❌ Fixed cost |
| HIPAA-eligible | ✅ Yes (BAA available) | ❌ No BAA for PHI |

### vs. Traditional Data Commons (iRODS, Dataverse)

| Feature | Cicada | iRODS/Dataverse |
|---------|--------|-----------------|
| Setup time | 10 minutes | Days/weeks |
| IT staff required | No | Yes |
| Cost (10TB) | $90/month | $500-1000/month |
| User experience | Modern GUI (Cloudscape) | Complex |
| Cloud-native | ✅ Yes | ❌ Usually on-prem |
| Auto-scaling | ✅ Yes | ❌ No |
| Self-service | ✅ Yes | ⚠️ Limited |
| HIPAA compliance | ✅ NIST 800-53 mode | ⚠️ Varies by deployment |
| Compliance automation | ✅ Built-in validation | ❌ Manual |

### vs. rclone/AWS CLI

| Feature | Cicada | rclone/AWS CLI |
|---------|--------|----------------|
| Metadata management | ✅ Rich schemas | ❌ Manual |
| User-friendly GUI | ✅ Yes | ❌ CLI only |
| Workflow integration | ✅ Built-in | ⚠️ Manual |
| Cost management | ✅ Tracking + alerts | ❌ Manual |
| Collaboration | ✅ Users/groups/projects | ❌ Manual IAM |
| Setup complexity | Low | High |

---

## Testimonials (Future)

*Space reserved for user testimonials after beta testing*

> "Cicada saved our lab when our NAS failed. All our data was already backed up to S3 without us thinking about it."
> — Dr. Sarah Chen, Cell Biology

> "We cut our data storage costs by 60% and got compute capabilities we never had before."
> — PI, Genomics Lab

> "The metadata search is incredible. I can find experiments from 2 years ago in seconds."
> — Postdoc, Chemistry

---

## FAQ

**Q: Do I need to know AWS to use Cicada?**
A: No! Cicada abstracts all AWS complexity. You just work with familiar concepts like folders, files, and projects.

**Q: What if I already have data in S3?**
A: Cicada can import existing S3 buckets and add metadata on top of your existing data.

**Q: Can I use my university's AWS account?**
A: Yes! Cicada supports bring-your-own-AWS-account mode.

**Q: Is my data secure?**
A: Yes. Data is encrypted at rest (S3/EBS) and in transit (TLS). Access is controlled via IAM policies.

**Q: Can I use Cicada for HIPAA-compliant research with patient data?**
A: Yes! Enable NIST 800-53 compliance mode which:
- Uses only HIPAA-eligible AWS services
- Requires a Business Associate Agreement (BAA) with AWS
- Enforces KMS encryption with customer-managed keys
- Maintains 6-year audit logs
- Implements comprehensive NIST 800-53 security controls
- Provides PHI detection and de-identification tools

You must sign a BAA with AWS before storing any PHI in Cicada.

**Q: What's the difference between NIST 800-171 and NIST 800-53 compliance modes?**
A:
- **NIST 800-171**: For Controlled Unclassified Information (CUI), government contracts, defense research. 110 security requirements.
- **NIST 800-53**: For HIPAA/PHI and FISMA compliance. More comprehensive controls (~900 controls across 20 families). Required for clinical trials and patient data.

Choose based on your data type: CUI → 800-171, PHI → 800-53.

**Q: How does Cicada verify that I have a BAA with AWS?**
A: AWS does not provide an API to verify BAA status - this is an AWS limitation affecting all HIPAA-compliant software.

Cicada's approach:
- **User attestation**: You confirm BAA during setup and provide signed date
- **Technical verification**: Cicada automatically verifies that only HIPAA-eligible services are in use, encryption is enabled, audit logs are configured, etc.
- **Ongoing reminders**: Monthly reminders to ensure BAA remains active
- **Compliance reports**: Generate reports showing technical controls are in place

This is the same approach used by enterprise HIPAA-compliant SaaS products. You are responsible for maintaining an active BAA with AWS.

**Q: What happens if Cicada development stops?**
A: Your data is in S3 with standard metadata files. You can access it directly via AWS tools even without Cicada.

**Q: Can I run Cicada on-premises?**
A: Not in MVP. Year 2+ roadmap includes on-premises deployment for institutions with data sovereignty requirements.

**Q: Does Cicada work with Azure or GCP?**
A: Not initially. AWS-only in MVP. Multi-cloud support in Year 2+ roadmap.

**Q: How do I get support?**
A: GitHub Discussions for community support. Optional paid support for institutions.

---

## Citation

If you use Cicada in your research, please cite:

```bibtex
@software{cicada2024,
  title = {Cicada: Dormant Data Commons for Academic Research},
  author = {Your Name and Contributors},
  year = {2024},
  url = {https://github.com/your-org/cicada},
  doi = {10.5281/zenodo.XXXXX}
}
```

---

## Acknowledgments

This project draws inspiration from:
- **iRODS**: Pioneering data management for research
- **Globus**: User-friendly file transfer and sharing
- **DataCite**: Making data citable and discoverable
- **rclone**: Excellent cloud storage synchronization
- **Snakemake/Nextflow**: Reproducible workflow systems
- **FAIR principles**: Guiding modern research data management

Special thanks to the research community for feedback during the design phase.

---

**Last Updated**: 2024-11-22
**Project Status**: Pre-implementation - Specifications Complete
**Next Steps**: Begin Phase 1 development

🦗 **Ready to build!**
