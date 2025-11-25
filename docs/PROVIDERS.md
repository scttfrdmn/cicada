# DOI Provider Setup Guide

> **Note:** Provider integration is an **optional advanced feature** for labs that need to publish datasets. Most Cicada usage involves core data management features (storage, sync, metadata extraction). This guide is only relevant if you need to mint DOIs for dataset publication.
>
> **Current Status (v0.2.0):** Framework for provider integration is implemented. Live API integration planned for v0.4.0+.

**Status:** v0.2.0 Documentation (Framework Ready, Live APIs in v0.4.0+)
**Audience:** Lab managers, data managers, researchers publishing datasets

## Overview

For labs that need to publish datasets, Cicada will support multiple DOI registration providers. Each provider offers different features, pricing, and repository integration.

**Supported Providers (v0.2.0):**

- **DataCite**: Direct DOI registration service for institutions
- **Zenodo**: Free repository with DOI assignment (CERN)

**Planned Providers (v0.3.0+):**
- Dryad
- Figshare
- Mendeley Data

**Key Concepts:**

- **DOI (Digital Object Identifier)**: Persistent identifier for datasets
- **Sandbox Environment**: Testing environment with fake DOIs for development
- **Production Environment**: Live environment creating real, permanent DOIs
- **Repository**: Storage service hosting dataset files
- **Metadata Schema**: DataCite Metadata Schema v4.5

## Quick Start

### Option 1: Zenodo (Recommended for Small Labs)

**Advantages:**
- Free (unlimited datasets)
- Integrated repository (storage + DOI)
- No institutional membership required
- Sandbox available for testing

**Setup Time:** 15 minutes

1. **Create Sandbox Account** (for testing):
   ```bash
   # Visit https://sandbox.zenodo.org
   # Sign up with GitHub, ORCID, or email
   ```

2. **Generate API Token**:
   - Settings → Applications → Personal access tokens
   - New token → Name: "Cicada" → Scopes: `deposit:write`, `deposit:actions`
   - Copy token (shown once)

3. **Configure Cicada**:
   ```bash
   cicada config set provider zenodo-sandbox
   cicada config set zenodo.token YOUR_TOKEN_HERE
   ```

4. **Test Publication**:
   ```bash
   cicada doi prepare sample.fastq \
     --enrich metadata.yaml \
     --provider zenodo-sandbox \
     --upload
   ```

5. **Production Setup** (when ready):
   - Create account at https://zenodo.org (same process)
   - Generate production token
   - Configure: `cicada config set provider zenodo`

### Option 2: DataCite (For Institutions)

**Advantages:**
- Direct DOI control
- Custom DOI prefix (e.g., 10.12345/...)
- Metadata-only registration (host files separately)
- Higher credibility with some journals

**Prerequisites:**
- Institutional DataCite membership ($$$)
- Repository ID from your institution
- API credentials from institution

**Setup Time:** 1-2 hours (including institutional approval)

1. **Get Institutional Credentials**:
   - Contact your library or IT department
   - Request: Repository ID, Username, Password
   - Ask about sandbox vs production access

2. **Configure Cicada**:
   ```bash
   cicada config set provider datacite-sandbox
   cicada config set datacite.repository_id YOUR_REPO_ID
   cicada config set datacite.username YOUR_USERNAME
   cicada config set datacite.password YOUR_PASSWORD
   ```

3. **Test Registration**:
   ```bash
   cicada doi prepare sample.fastq \
     --enrich metadata.yaml \
     --provider datacite-sandbox
   ```

## DataCite Setup

### Understanding DataCite

DataCite is a global DOI registration agency focused on research data. Unlike Zenodo, DataCite provides only DOI registration—you must host files separately.

**When to Use DataCite:**
- Your institution has DataCite membership
- You need custom DOI prefixes
- You have separate file storage (institutional repository, S3)
- You require direct control over DOI metadata

**Cost:**
- Institutional membership: $5,000-10,000/year (varies)
- Per-DOI: Usually unlimited within membership

### Sandbox Environment

DataCite provides a test environment for development and validation.

**Sandbox Details:**
- **URL**: https://api.test.datacite.org
- **DOIs**: Test prefix `10.82041` (not resolvable outside sandbox)
- **Purpose**: Testing workflows without creating real DOIs
- **Persistence**: Data may be deleted periodically

### Production Environment

**Production Details:**
- **URL**: https://api.datacite.org
- **DOIs**: Your institution's prefix (e.g., `10.12345`)
- **Purpose**: Creating permanent, citable DOIs
- **Persistence**: DOIs are permanent and cannot be deleted

### Step-by-Step Setup

#### 1. Obtain Credentials

Contact your institution's library or research data management team:

**Request Template:**
```
Subject: DataCite API Access for Research Data Management

Hello,

I'm setting up Cicada (https://github.com/scttfrdmn/cicada) to automate
DOI registration for our lab's research datasets.

Could you please provide:
1. DataCite Repository ID
2. API username
3. API password
4. Sandbox credentials (for testing)
5. DOI prefix assigned to our lab/department

Thank you,
[Your Name]
```

**What You'll Receive:**
- Repository ID: `INST.DEPT` (e.g., `MIT.BIO`)
- Username: Usually email or `repoID.username`
- Password: Generated password
- DOI Prefix: `10.XXXXX`

#### 2. Test Sandbox Access

Verify credentials work:

```bash
# Set sandbox configuration
export DATACITE_REPO_ID="INST.DEPT"
export DATACITE_USERNAME="your_username"
export DATACITE_PASSWORD="your_password"

# Test API access
curl -u "$DATACITE_USERNAME:$DATACITE_PASSWORD" \
  https://api.test.datacite.org/dois \
  | jq '.'

# Should return JSON with your test DOIs
```

#### 3. Configure Cicada (Sandbox)

```bash
# Configure for testing
cicada config set provider datacite-sandbox
cicada config set datacite.repository_id "$DATACITE_REPO_ID"
cicada config set datacite.username "$DATACITE_USERNAME"
cicada config set datacite.password "$DATACITE_PASSWORD"

# Verify configuration
cicada config list
```

**Config File Location:** `~/.config/cicada/config.yaml`

**Config File Contents:**
```yaml
provider: datacite-sandbox
datacite:
  repository_id: INST.DEPT
  username: your_username
  password: your_password
  sandbox: true
```

#### 4. Test DOI Creation

Create a test DOI with sample data:

```bash
# Create test FASTQ file
echo "@SEQ_ID
ACGTACGTACGTACGT
+
IIIIIIIIIIIIIIII" > test.fastq

# Create enrichment metadata
cat > enrich.yaml <<EOF
title: "Test Dataset for DataCite Integration"
authors:
  - name: Test Researcher
    orcid: 0000-0001-2345-6789
    affiliation: Test University
description: "This is a test dataset for validating DataCite integration"
publisher: Test University
publication_year: 2025
EOF

# Prepare DOI (draft state)
cicada doi prepare test.fastq \
  --enrich enrich.yaml \
  --provider datacite-sandbox \
  --output doi.json

# Review DOI
cat doi.json | jq '.doi'
# Output: "10.82041/test-dataset-12345"
```

#### 5. Publish DOI (Sandbox)

**⚠️ Important:** In sandbox, DOIs are temporary. In production, they're permanent.

```bash
# Publish DOI (makes it findable)
cicada doi publish test.fastq \
  --metadata doi.json \
  --provider datacite-sandbox

# Verify publication
cicada doi status 10.82041/test-dataset-12345 \
  --provider datacite-sandbox

# Output:
# DOI: 10.82041/test-dataset-12345
# State: findable
# URL: https://handle.test.datacite.org/10.82041/test-dataset-12345
```

#### 6. Configure Production

**⚠️ Only after successful sandbox testing**

```bash
# Switch to production
cicada config set provider datacite
cicada config set datacite.sandbox false

# Production credentials (may be same as sandbox)
cicada config set datacite.repository_id "$PROD_REPO_ID"
cicada config set datacite.username "$PROD_USERNAME"
cicada config set datacite.password "$PROD_PASSWORD"

# Verify production access
cicada doi list --provider datacite
```

### DataCite Workflows

#### Draft → Review → Publish

```bash
# Step 1: Create draft DOI
cicada doi prepare dataset.fastq \
  --enrich metadata.yaml \
  --provider datacite \
  --state draft \
  --output doi.json

# Step 2: Review metadata
cat doi.json | jq '.'

# Step 3: Update if needed
cicada doi update 10.12345/dataset-001 \
  --metadata updated_metadata.yaml \
  --provider datacite

# Step 4: Publish when ready
cicada doi publish 10.12345/dataset-001 \
  --provider datacite
```

#### Metadata-Only Registration

DataCite doesn't host files—provide URL to where files are stored:

```bash
# Prepare DOI with file URL
cat > enrich.yaml <<EOF
title: "My Dataset"
authors:
  - name: Researcher Name
description: "Dataset description"
url: "https://mylab.edu/datasets/dataset-001"
EOF

cicada doi prepare dataset.fastq \
  --enrich enrich.yaml \
  --provider datacite \
  --no-upload  # Don't try to upload files
```

## Zenodo Setup

### Understanding Zenodo

Zenodo is a free, open-access repository operated by CERN. It provides both file storage and DOI registration.

**When to Use Zenodo:**
- Your lab doesn't have DataCite membership
- You want free, unlimited dataset hosting
- You need a simple, integrated solution
- You publish openly and want broad discoverability

**Cost:**
- **Free** (up to 50 GB per dataset)
- Larger datasets: Contact Zenodo for approval

### Sandbox Environment

Zenodo provides a complete sandbox for testing.

**Sandbox Details:**
- **URL**: https://sandbox.zenodo.org
- **DOIs**: Test DOIs `10.5072/zenodo.XXXXXX`
- **Purpose**: Full-featured testing environment
- **Storage**: Same as production (50 GB per dataset)
- **Persistence**: Data retained indefinitely

### Production Environment

**Production Details:**
- **URL**: https://zenodo.org
- **DOIs**: Real DOIs `10.5281/zenodo.XXXXXX`
- **Indexing**: Indexed by Google Scholar, OpenAIRE, DataCite
- **Persistence**: Permanent (files and metadata cannot be deleted)

### Step-by-Step Setup

#### 1. Create Sandbox Account

Visit https://sandbox.zenodo.org:

1. Click **Sign up** (top right)
2. Choose authentication:
   - **GitHub** (recommended for developers)
   - **ORCID** (recommended for researchers)
   - **Email** (basic signup)
3. Complete profile:
   - Full name
   - Affiliation
   - ORCID (if not using ORCID login)

#### 2. Generate API Token (Sandbox)

1. Click profile icon → **Applications**
2. Scroll to **Personal access tokens**
3. Click **New token**
4. Configure:
   - **Name**: `Cicada Testing`
   - **Scopes**:
     - ✅ `deposit:write` (create and update deposits)
     - ✅ `deposit:actions` (publish deposits)
5. Click **Create**
6. **Copy token immediately** (shown once)

**Token Format:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...` (long string)

#### 3. Test API Access

```bash
# Set token
export ZENODO_TOKEN="your_token_here"

# List your deposits
curl "https://sandbox.zenodo.org/api/deposit/depositions" \
  -H "Authorization: Bearer $ZENODO_TOKEN" \
  | jq '.'

# Should return empty array: []
```

#### 4. Configure Cicada (Sandbox)

```bash
# Configure sandbox
cicada config set provider zenodo-sandbox
cicada config set zenodo.token "$ZENODO_TOKEN"
cicada config set zenodo.sandbox true

# Verify
cicada config list
```

**Config File:** `~/.config/cicada/config.yaml`
```yaml
provider: zenodo-sandbox
zenodo:
  token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
  sandbox: true
```

#### 5. Create Test Upload

```bash
# Create test data
echo "@SEQ_ID
ACGTACGTACGTACGT
+
IIIIIIIIIIIIIIII" > test_R1.fastq

# Create metadata
cat > metadata.yaml <<EOF
title: "Test Sequencing Dataset"
authors:
  - name: Your Name
    orcid: 0000-0001-2345-6789
    affiliation: Your University
description: |
  Test dataset for Zenodo integration testing.
  This is a small FASTQ file for validation purposes.
keywords:
  - test
  - sequencing
  - FASTQ
upload_type: dataset
access_right: open
license: cc-by-4.0
EOF

# Prepare and upload
cicada doi prepare test_R1.fastq \
  --enrich metadata.yaml \
  --provider zenodo-sandbox \
  --upload \
  --output zenodo_response.json

# Check response
cat zenodo_response.json | jq '{doi, title, state, links}'
```

**Response:**
```json
{
  "doi": "10.5072/zenodo.123456",
  "title": "Test Sequencing Dataset",
  "state": "unsubmitted",
  "links": {
    "record": "https://sandbox.zenodo.org/record/123456",
    "bucket": "https://sandbox.zenodo.org/api/files/...",
    "publish": "https://sandbox.zenodo.org/api/deposit/depositions/123456/actions/publish"
  }
}
```

#### 6. Publish Deposition

```bash
# Publish to make DOI active
cicada doi publish 10.5072/zenodo.123456 \
  --provider zenodo-sandbox

# Visit record in browser
open "https://sandbox.zenodo.org/record/123456"
```

#### 7. Production Setup

**⚠️ Only after successful sandbox testing**

1. **Create Production Account**: https://zenodo.org (same process)
2. **Generate Production Token**: Same steps as sandbox
3. **Configure Cicada**:
   ```bash
   cicada config set provider zenodo
   cicado config set zenodo.token "$PROD_TOKEN"
   cicada config set zenodo.sandbox false
   ```

### Zenodo Workflows

#### Complete Upload Workflow

```bash
# Step 1: Prepare files
cicada metadata extract sample_R1.fastq.gz --output metadata.json
cicada metadata extract sample_R2.fastq.gz --output metadata.json --append

# Step 2: Create enrichment
cat > enrich.yaml <<EOF
title: "Bacterial genome sequencing - Sample 042"
authors:
  - name: Dr. Jane Smith
    orcid: 0000-0002-1234-5678
    affiliation: Department of Microbiology, State University
description: |
  Whole genome sequencing of antibiotic-resistant E. coli strain 042.
  Sequenced on Illumina NovaSeq 6000, 2x150bp paired-end reads.
keywords:
  - whole genome sequencing
  - E. coli
  - antibiotic resistance
upload_type: dataset
access_right: open
license: cc-by-4.0
EOF

# Step 3: Prepare DOI with both files
cicada doi prepare sample_R1.fastq.gz sample_R2.fastq.gz \
  --enrich enrich.yaml \
  --provider zenodo \
  --upload \
  --output doi_response.json

# Step 4: Review
cat doi_response.json | jq '.'

# Step 5: Publish
doi=$(jq -r '.doi' doi_response.json)
cicada doi publish "$doi" --provider zenodo

# Step 6: Get URL
echo "Dataset published: https://zenodo.org/doi/$doi"
```

#### Update Published Record

```bash
# Create new version
cicada doi version 10.5281/zenodo.123456 \
  --provider zenodo \
  --output new_version.json

# Upload new files
new_doi=$(jq -r '.doi' new_version.json)
cicada doi upload "$new_doi" updated_data.fastq \
  --provider zenodo

# Publish new version
cicada doi publish "$new_doi" --provider zenodo
```

## Provider Comparison

| Feature | DataCite | Zenodo | Dryad (v0.3.0) | Figshare (v0.3.0) |
|---------|----------|--------|----------------|-------------------|
| **Cost** | Institutional membership | Free | $120/dataset | Free (limited) |
| **Storage** | None (metadata only) | 50 GB/dataset | Unlimited | 20 GB free |
| **DOI Prefix** | Custom (10.XXXXX) | Fixed (10.5281) | Fixed (10.5061) | Fixed (10.6084) |
| **Sandbox** | ✅ Yes | ✅ Yes | ✅ Yes | ❌ No |
| **File Hosting** | ❌ No | ✅ Yes | ✅ Yes | ✅ Yes |
| **API Access** | ✅ Full | ✅ Full | ✅ Full | ✅ Full |
| **Setup Time** | 1-2 hours | 15 minutes | 30 minutes | 30 minutes |
| **Best For** | Institutions | Small labs | Publications | Figures/datasets |

### Choosing a Provider

**Choose DataCite if:**
- Your institution has membership
- You need custom DOI prefix
- You host files separately (S3, institutional repository)
- You want direct control over metadata

**Choose Zenodo if:**
- You want free, simple solution
- You need integrated file hosting
- You publish open access
- You're a small lab or individual researcher

**Choose Dryad if (v0.3.0):**
- You're publishing in journals requiring Dryad (e.g., Evolution, Am Nat)
- You have large datasets (> 50 GB)
- You want curated, journal-integrated submission

**Choose Figshare if (v0.3.0):**
- You have figures, posters, presentations (not just data)
- You want institutional Figshare integration
- You need file versioning and previews

## Configuration Management

### Configuration File Location

- **Linux/macOS**: `~/.config/cicada/config.yaml`
- **Windows**: `%APPDATA%\cicada\config.yaml`

### Configuration Commands

```bash
# View all settings
cicada config list

# Set individual values
cicada config set provider zenodo
cicada config set zenodo.token "your_token"

# Get specific value
cicada config get provider

# Reset to defaults
cicada config reset

# Edit config file directly
vim ~/.config/cicada/config.yaml
```

### Multi-Provider Configuration

Configure multiple providers and switch as needed:

**config.yaml:**
```yaml
# Active provider
provider: zenodo

# DataCite (production)
datacite:
  repository_id: MIT.BIO
  username: mit_bio_user
  password: secret_password
  sandbox: false

# DataCite (sandbox)
datacite_sandbox:
  repository_id: MIT.BIO
  username: mit_bio_user
  password: sandbox_password
  sandbox: true

# Zenodo (production)
zenodo:
  token: prod_token_here
  sandbox: false

# Zenodo (sandbox)
zenodo_sandbox:
  token: sandbox_token_here
  sandbox: true
```

**Switch providers:**
```bash
# Use Zenodo sandbox
cicada config set provider zenodo-sandbox

# Use DataCite production
cicada config set provider datacite

# Use Zenodo production
cicada config set provider zenodo
```

### Environment Variables

Override config file with environment variables:

```bash
# Provider
export CICADA_PROVIDER=zenodo-sandbox

# Zenodo
export ZENODO_TOKEN=your_token
export ZENODO_SANDBOX=true

# DataCite
export DATACITE_REPOSITORY_ID=INST.DEPT
export DATACITE_USERNAME=username
export DATACITE_PASSWORD=password
export DATACITE_SANDBOX=true

# Run command
cicada doi prepare data.fastq --enrich metadata.yaml --upload
```

Environment variables take precedence over config file.

## Testing Workflows

### Pre-Production Testing Checklist

Before publishing real DOIs, test thoroughly in sandbox:

- [ ] **API Access**: Verify credentials work
  ```bash
  cicada doi list --provider zenodo-sandbox
  ```

- [ ] **Metadata Validation**: Check quality score
  ```bash
  cicada doi validate data.fastq --enrich metadata.yaml
  ```

- [ ] **Draft Creation**: Create draft DOI
  ```bash
  cicada doi prepare data.fastq --enrich metadata.yaml --provider zenodo-sandbox
  ```

- [ ] **File Upload**: Upload test files (Zenodo only)
  ```bash
  cicada doi prepare data.fastq --enrich metadata.yaml --upload --provider zenodo-sandbox
  ```

- [ ] **Publication**: Publish test DOI
  ```bash
  cicada doi publish 10.5072/zenodo.123456 --provider zenodo-sandbox
  ```

- [ ] **URL Resolution**: Verify DOI resolves
  ```bash
  open "https://sandbox.zenodo.org/doi/10.5072/zenodo.123456"
  ```

- [ ] **Metadata Display**: Check all fields display correctly

- [ ] **File Download**: Download files and verify integrity

### Automated Testing

Create test script for CI/CD:

```bash
#!/bin/bash
# test_doi_workflow.sh

set -e  # Exit on error

PROVIDER="zenodo-sandbox"
TEST_FILE="test_data.fastq"
METADATA="test_metadata.yaml"

echo "Testing DOI workflow with $PROVIDER..."

# Step 1: Validate
echo "1. Validating metadata..."
cicada doi validate "$TEST_FILE" --enrich "$METADATA"

# Step 2: Prepare
echo "2. Creating draft DOI..."
cicada doi prepare "$TEST_FILE" \
  --enrich "$METADATA" \
  --provider "$PROVIDER" \
  --upload \
  --output response.json

# Step 3: Extract DOI
DOI=$(jq -r '.doi' response.json)
echo "Created DOI: $DOI"

# Step 4: Publish
echo "3. Publishing DOI..."
cicada doi publish "$DOI" --provider "$PROVIDER"

# Step 5: Verify
echo "4. Verifying publication..."
cicada doi status "$DOI" --provider "$PROVIDER"

echo "✅ Test workflow completed successfully"
```

## Publishing Workflows

### Workflow 1: Pre-Publication Dataset

Publish dataset before paper submission:

```bash
# Step 1: Extract and validate metadata
cicada metadata extract data_R1.fastq.gz --preset illumina-novaseq
cicada metadata extract data_R2.fastq.gz --preset illumina-novaseq

# Step 2: Create comprehensive metadata
cat > publication_metadata.yaml <<EOF
title: "Genome-wide association study of antibiotic resistance in E. coli"
authors:
  - name: Dr. Jane Smith
    orcid: 0000-0002-1234-5678
    affiliation: State University
  - name: Dr. John Doe
    orcid: 0000-0003-9876-5432
    affiliation: State University
description: |
  Raw sequencing data supporting our manuscript "Mechanisms of antibiotic
  resistance evolution in E. coli" submitted to Nature Microbiology.

  Data includes whole genome sequencing of 50 clinical isolates...
keywords:
  - GWAS
  - E. coli
  - antibiotic resistance
  - whole genome sequencing
related_identifiers:
  - identifier: "10.1101/2025.01.123456"  # bioRxiv preprint
    relation: IsSupplementTo
    type: DOI
EOF

# Step 3: Prepare DOI
cicada doi prepare data_*.fastq.gz \
  --enrich publication_metadata.yaml \
  --provider zenodo \
  --upload \
  --output dataset_doi.json

# Step 4: Review metadata
cat dataset_doi.json | jq '.'

# Step 5: Publish
doi=$(jq -r '.doi' dataset_doi.json)
cicada doi publish "$doi" --provider zenodo

# Step 6: Include DOI in manuscript
echo "Data Availability: Dataset available at https://zenodo.org/doi/$doi"
```

### Workflow 2: Supplementary Data for Published Paper

Add DOI after paper acceptance:

```bash
# Prepare metadata with paper DOI
cat > supplementary_metadata.yaml <<EOF
title: "Supplementary Data: [Paper Title]"
authors:
  - name: Paper Author 1
  - name: Paper Author 2
description: "Supplementary dataset for our paper published in [Journal]"
related_identifiers:
  - identifier: "10.1234/journal.2025.5678"  # Paper DOI
    relation: IsSupplementTo
    type: DOI
publication_year: 2025
EOF

# Create DOI
cicada doi prepare supplementary_data.zip \
  --enrich supplementary_metadata.yaml \
  --provider zenodo \
  --upload

# Publish immediately (paper already accepted)
cicada doi publish [DOI] --provider zenodo
```

### Workflow 3: Dataset Series

Create related datasets with version control:

```bash
# Version 1: Raw data
cicada doi prepare raw_data_v1.tar.gz \
  --enrich metadata_v1.yaml \
  --provider zenodo \
  --upload \
  --output v1_doi.json

v1_doi=$(jq -r '.doi' v1_doi.json)

# Version 2: Processed data (link to v1)
cat > metadata_v2.yaml <<EOF
title: "Processed Data - Version 2"
description: "Processed version of raw data"
related_identifiers:
  - identifier: "$v1_doi"
    relation: IsNewVersionOf
    type: DOI
EOF

cicada doi prepare processed_data_v2.tar.gz \
  --enrich metadata_v2.yaml \
  --provider zenodo \
  --upload
```

## Troubleshooting

### Authentication Failures

**Error:** `401 Unauthorized`

**Zenodo Solutions:**
```bash
# Check token is set
cicada config get zenodo.token

# Regenerate token at https://zenodo.org (Settings → Applications)
cicada config set zenodo.token "new_token_here"

# Test authentication
curl "https://zenodo.org/api/deposit/depositions" \
  -H "Authorization: Bearer $(cicada config get zenodo.token)"
```

**DataCite Solutions:**
```bash
# Verify credentials with institution
cicada config get datacite.repository_id
cicada config get datacite.username

# Test API access
curl -u "username:password" https://api.datacite.org/dois
```

### Upload Failures

**Error:** `File upload failed: Connection timeout`

**Solutions:**

1. **Check file size** (Zenodo limit: 50 GB)
   ```bash
   ls -lh data.fastq.gz
   ```

2. **Compress files**
   ```bash
   gzip -9 data.fastq  # Maximum compression
   ```

3. **Split large files**
   ```bash
   split -b 10G data.tar.gz data_part_
   ```

4. **Upload separately**
   ```bash
   # Create deposition first
   cicada doi prepare data.fastq --enrich metadata.yaml --no-upload

   # Upload files individually
   cicada doi upload [DOI] data_part_aa
   cicada doi upload [DOI] data_part_ab
   ```

### Validation Failures

**Error:** `Validation failed: Missing required field 'publisher'`

**Solution:**

Add missing fields to enrichment file:

```yaml
# Required by DataCite
publisher: "Your Institution"
publication_year: 2025
resource_type: Dataset

# Recommended
description: "Detailed description of dataset"
subjects: ["Biology", "Genomics"]
```

### DOI Already Exists

**Error:** `DOI 10.5281/zenodo.123456 already exists`

**Solution:**

Either update existing DOI or create new version:

```bash
# Option 1: Update existing
cicada doi update 10.5281/zenodo.123456 \
  --metadata new_metadata.yaml \
  --provider zenodo

# Option 2: Create new version
cicada doi version 10.5281/zenodo.123456 \
  --provider zenodo
```

## Best Practices

### 1. Always Test in Sandbox First

```bash
# ❌ Don't do this first time
cicada doi prepare data.fastq --provider zenodo --upload

# ✅ Do this instead
cicada doi prepare data.fastq --provider zenodo-sandbox --upload
```

### 2. Validate Before Preparation

```bash
# Check quality score first
cicada doi validate data.fastq --enrich metadata.yaml

# If score < 75, add more metadata
# If score >= 80, proceed with preparation
```

### 3. Use Version Control for Metadata

```bash
# Store metadata in git
git add metadata/
git commit -m "Add DOI metadata for Dataset 042"
```

### 4. Document DOIs

Keep a lab registry:

**dois.md:**
```markdown
# Lab DOI Registry

| Dataset | DOI | Published | Paper |
|---------|-----|-----------|-------|
| WGS E. coli 042 | 10.5281/zenodo.123456 | 2025-01-15 | Nature Micro |
| RNA-seq time series | 10.5281/zenodo.123457 | 2025-02-01 | Cell |
```

### 5. Include Data Availability Statements

Standard text for papers:

```
Data Availability: Raw sequencing data have been deposited in Zenodo
under accession code 10.5281/zenodo.123456.
```

### 6. Monitor Usage

Check dataset views and downloads:

```bash
# Zenodo: View statistics at record page
open "https://zenodo.org/record/123456"

# DataCite: Check usage statistics
cicada doi stats 10.12345/dataset-001 --provider datacite
```

## Related Documentation

- **[DOI Workflow Guide](DOI_WORKFLOW.md)**: Preparing metadata for DOI registration
- **[Metadata Extraction Guide](METADATA_EXTRACTION.md)**: Extracting metadata from files
- **[User Scenarios](USER_SCENARIOS_v0.2.0.md)**: Real-world publishing workflows

## Support

For provider-specific issues:

- **Zenodo**: https://zenodo.org/support
- **DataCite**: https://support.datacite.org
- **Cicada**: https://github.com/scttfrdmn/cicada/issues

## Version History

- **v0.2.0** (Current): DataCite and Zenodo support
- **v0.3.0** (Planned): Dryad, Figshare, Mendeley Data support
