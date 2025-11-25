# DOI Workflow Guide

> **Note:** DOI preparation is an **optional advanced feature** of Cicada's data commons platform. Most labs use Cicada for daily data management, storage, and metadata extraction. DOI support is provided for labs that need to publish datasets to repositories like Zenodo or institutional data repositories.
>
> **For core data management features**, see the [User Guide](USER_GUIDE.md) (coming in v0.3.0) and [Metadata Extraction Guide](METADATA_EXTRACTION.md).

This guide covers preparing datasets for DOI (Digital Object Identifier) registration using Cicada's automated metadata mapping and validation tools. Learn how to assess DOI readiness, enrich metadata, and prepare DataCite-compliant metadata for repository submission.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Understanding DOI Requirements](#understanding-doi-requirements)
- [Command Reference](#command-reference)
- [Metadata Enrichment](#metadata-enrichment)
- [Quality Scoring](#quality-scoring)
- [Validation Modes](#validation-modes)
- [Repository Submission](#repository-submission)
- [Complete Workflow Examples](#complete-workflow-examples)
- [Troubleshooting](#troubleshooting)

---

## Overview

DOIs (Digital Object Identifiers) provide permanent, citable identifiers for research datasets that you wish to publish. This is typically the **final step** in your data lifecycle:

```
1. Generate data → 2. Store & sync → 3. Extract metadata → 4. Organize & manage → 5. Publish (DOI)
   (Instrument)      (Cicada core)     (Cicada core)      (Cicada core)           (This guide)
```

Most research data stays within your lab's data commons for analysis and collaboration. You only need DOI registration when:
- Publishing datasets alongside papers
- Making data publicly available
- Meeting funder requirements for data sharing
- Archiving data in institutional repositories

### DOI Preparation Workflow

Cicada helps prepare datasets for DOI registration by:

- **Extracting** base metadata from instrument files (core platform feature)
- **Mapping** metadata to DataCite Schema v4.5
- **Validating** metadata completeness and quality
- **Enriching** metadata with author information and descriptions
- **Scoring** metadata quality (0-100 scale)
- **Exporting** DataCite-compliant metadata for repository submission

### Benefits

✅ **Leverages existing metadata**: Builds on metadata already extracted for data management
✅ **Guided workflow**: Know exactly what metadata is required for publication
✅ **Quality scoring**: Track metadata completeness (0-100)
✅ **Standards compliance**: DataCite Schema v4.5 compatible
✅ **Repository ready**: Export metadata for Zenodo, Dryad, institutional repositories

---

## Quick Start

### Check DOI Readiness

```bash
# Assess current metadata quality
cicada doi validate sample.fastq.gz
```

**Output**:
```
DOI Validation Results
======================

File: sample.fastq.gz

✗ NOT READY for DOI minting

Quality Score: 47.0/100 (Moderate)

Present Fields (8):
  ✓ identifier
  ✓ title
  ✓ publisher
  ✓ publication_year
  ✓ resource_type
  ✓ description
  ✓ license
  ✓ keywords

Missing Fields (9):
  ✗ real creator names
  ✗ url
  ✗ author ORCIDs
  ✗ author affiliations
  ...

Recommendations:
  Fix errors before minting DOI:
    - Add real author names with ORCIDs
    - Enhance description with methodology
```

### Prepare DOI with Enrichment

```bash
# Prepare with enriched metadata
cicada doi prepare sample.fastq.gz \
  --enrich enrichment.yaml \
  --publisher "University Lab"
```

**Output**:
```
DOI Preparation Results
=======================

File: sample.fastq.gz

Dataset Information:
  Title: My Research Dataset
  Authors: 2 (both with ORCIDs)
  Quality Score: 91.0/100 (Excellent)

✓ Ready for DOI minting
```

---

## Understanding DOI Requirements

### DataCite Required Fields

These fields are **mandatory** for DOI registration:

| Field | Description | How Cicada Provides |
|-------|-------------|---------------------|
| **Identifier** | Unique identifier | Auto-generated or from filename |
| **Creators** | Authors/creators | From enrichment file |
| **Titles** | Dataset title | From enrichment file or filename |
| **Publisher** | Publishing entity | From `--publisher` flag |
| **Publication Year** | Year of publication | Current year or from enrichment |
| **Resource Type** | Type of resource | Auto-detected ("Dataset") |

### DataCite Recommended Fields

These fields are **highly recommended** for discovery and attribution:

| Field | Description | Impact on Quality Score |
|-------|-------------|------------------------|
| **Subjects** | Keywords/topics | +5 points |
| **Contributors** | Other contributors | +3 points |
| **Dates** | Relevant dates | +3 points |
| **Related Identifiers** | Related DOIs/URLs | +5 points |
| **Descriptions** | Detailed descriptions | +10 points |
| **Geo Locations** | Geographic coverage | +3 points |
| **Language** | Primary language | +2 points |
| **Sizes** | Data sizes | Auto-extracted |
| **Formats** | File formats | Auto-extracted |
| **Version** | Dataset version | +2 points |
| **Rights** | License information | +5 points |
| **Funding References** | Grant information | +5 points |

### Quality Score Calculation

- **0-59**: Not ready for DOI (fix errors first)
- **60-79**: Acceptable (minimum for DOI registration)
- **80-89**: Good (recommended for publication)
- **90-100**: Excellent (comprehensive metadata)

---

## Command Reference

### `cicada doi validate`

Validate metadata for DOI readiness without preparing.

#### Syntax

```bash
cicada doi validate <file> [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `table` | Output format: `table`, `json`, `yaml` |
| `--min-score` | float | `60.0` | Minimum quality score threshold |

#### Examples

```bash
# Basic validation
cicada doi validate data.fastq.gz

# JSON output for parsing
cicada doi validate data.fastq.gz --format json

# Require higher quality score
cicada doi validate data.fastq.gz --min-score 80
```

---

### `cicada doi prepare`

Prepare metadata for DOI registration with enrichment.

#### Syntax

```bash
cicada doi prepare <file> [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--enrich` | string | - | Enrichment metadata file (YAML or JSON) |
| `--publisher` | string | - | Publisher name (required) |
| `--license` | string | `CC-BY-4.0` | License identifier |
| `--preset` | string | - | Instrument preset for validation |
| `--format` | string | `table` | Output format: `table`, `json`, `yaml` |
| `--output` | string | stdout | Output file path |

#### Examples

```bash
# Prepare with enrichment
cicada doi prepare data.fastq.gz \
  --enrich metadata.yaml \
  --publisher "State University Lab"

# Save DataCite JSON
cicada doi prepare data.fastq.gz \
  --enrich metadata.yaml \
  --publisher "Lab Name" \
  --format json \
  --output datacite.json

# With preset validation
cicada doi prepare data.fastq.gz \
  --enrich metadata.yaml \
  --publisher "Lab" \
  --preset illumina-novaseq

# Different license
cicada doi prepare data.fastq.gz \
  --enrich metadata.yaml \
  --publisher "Lab" \
  --license CC0-1.0
```

---

## Metadata Enrichment

### Creating an Enrichment File

Enrichment files provide additional metadata not extractable from files.

#### YAML Format (Recommended)

```yaml
# enrichment.yaml

title: "Whole genome sequencing of antibiotic-resistant bacteria"

authors:
  - name: Dr. Jane Smith
    orcid: 0000-0002-1234-5678
    affiliation: Department of Microbiology, State University
  - name: Dr. John Doe
    orcid: 0000-0003-9876-5432
    affiliation: Department of Bioinformatics, State University
    role: supervisor

description: |
  This dataset contains whole genome sequencing data from 50 clinical
  isolates of antibiotic-resistant bacteria collected from hospital
  patients between 2023-2024.

  Methodology:
  - DNA extraction: Qiagen DNeasy Kit
  - Sequencing: Illumina NovaSeq 6000, 2x150bp paired-end
  - Coverage: Average 100x per isolate
  - Quality filtering: fastp v0.23.0 (Q30)

  Data includes raw FASTQ files and quality control reports.

keywords:
  - whole genome sequencing
  - antibiotic resistance
  - bacterial genomics
  - clinical isolates
  - antimicrobial resistance genes

publisher: State University Genomics Core

license: CC-BY-4.0

funding_references:
  - funder_name: National Institutes of Health
    award_number: R01AI123456
  - funder_name: University Research Foundation
    award_number: URF-2024-789

related_identifiers:
  - identifier: "10.1234/journal.2025.456"
    relation: IsSupplementTo
    type: DOI
    description: "Associated publication"

temporal_coverage:
  start: "2023-01-01"
  end: "2024-12-31"

version: "1.0"

language: en
```

#### JSON Format

```json
{
  "title": "My Research Dataset",
  "authors": [
    {
      "name": "Dr. Jane Smith",
      "orcid": "0000-0002-1234-5678",
      "affiliation": "State University"
    }
  ],
  "description": "Dataset description with methodology...",
  "keywords": ["keyword1", "keyword2", "keyword3"],
  "funding_references": [
    {
      "funder_name": "NSF",
      "award_number": "NSF-123456"
    }
  ]
}
```

### Enrichment File Fields

#### Required in Enrichment

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `title` | string | Dataset title | `"RNA-seq analysis of..."` |
| `authors` | array | List of authors | See below |
| `description` | string | Detailed description | Multi-line text |

#### Author Object

```yaml
authors:
  - name: Dr. Jane Smith         # Required
    orcid: 0000-0002-1234-5678  # Highly recommended
    affiliation: State University  # Recommended
    role: principal investigator   # Optional
```

#### Recommended in Enrichment

| Field | Type | Description |
|-------|------|-------------|
| `keywords` | array | 5-10 keywords |
| `funding_references` | array | Funding sources |
| `related_identifiers` | array | Related DOIs/URLs |
| `temporal_coverage` | object | Time period covered |
| `version` | string | Dataset version |
| `language` | string | Primary language (ISO 639-1) |

---

## Quality Scoring

### Score Breakdown

Quality scores are calculated based on metadata completeness:

```
Total Score = Base Score + Optional Score

Base Score (60 points):
  - Required fields present: 60 points
  - Any required field missing: 0 points

Optional Score (40 points):
  - Each optional field present: varies by importance
  - Author ORCIDs: +5 points
  - Detailed description: +10 points
  - Funding information: +5 points
  - Related publications: +5 points
  - Keywords (5+): +5 points
  - ... (other optional fields)
```

### Improving Your Score

| Current Score | Actions to Improve |
|---------------|-------------------|
| **0-40** | Add required fields: title, authors, description |
| **40-60** | Complete all required fields |
| **60-70** | Add author ORCIDs and affiliations |
| **70-80** | Add funding information and keywords |
| **80-90** | Add related publications and detailed methodology |
| **90-100** | Add temporal/spatial coverage, version info |

### Example Score Progression

```bash
# Initial validation (no enrichment)
cicada doi validate data.fastq.gz
# Quality Score: 47.0/100 (Moderate)

# Add basic enrichment (title, authors, description)
cicada doi prepare data.fastq.gz --enrich basic.yaml
# Quality Score: 68.0/100 (Acceptable)

# Add ORCIDs and funding
cicada doi prepare data.fastq.gz --enrich enhanced.yaml
# Quality Score: 85.0/100 (Good)

# Add related publications and full methodology
cicada doi prepare data.fastq.gz --enrich complete.yaml
# Quality Score: 94.0/100 (Excellent)
```

---

## Validation Modes

### Lenient Mode (Default)

Allows DOI preparation with warnings:
- Score must be ≥ 60
- Required fields must be present
- Warnings don't block preparation

**Use when**: Preparing for internal/initial DOI registration

### Strict Mode

Requires high-quality metadata:
- Score must be ≥ 80 (configurable)
- All warnings should be addressed
- Best practices enforced

**Use when**: Preparing for publication or public repositories

```bash
# Strict validation
cicada doi prepare data.fastq.gz \
  --enrich metadata.yaml \
  --publisher "Lab" \
  --min-score 80
```

---

## Repository Submission

### Supported Repositories

Cicada prepares metadata for any DataCite-compatible repository:

| Repository | DOI Provider | Test Environment | Cost |
|------------|--------------|------------------|------|
| **Zenodo** | DataCite | ✅ sandbox.zenodo.org | Free |
| **Dryad** | DataCite | ✅ Test instance | $120/dataset* |
| **Figshare** | DataCite | ✅ Test instance | Free (< 20GB) |
| **Mendeley Data** | DataCite | Limited | Free |
| **Institutional Repositories** | Varies | Varies | Varies |

*Some institutions have subscriptions

### Zenodo Submission Workflow

1. **Prepare metadata**:
   ```bash
   cicada doi prepare data.fastq.gz \
     --enrich enrichment.yaml \
     --publisher "Your Institution" \
     --format json \
     --output datacite.json
   ```

2. **Upload to Zenodo**:
   - Go to https://zenodo.org (or https://sandbox.zenodo.org for testing)
   - Click "New upload"
   - Upload your data files

3. **Import metadata**:
   - Click "Import" → "DataCite JSON"
   - Upload `datacite.json`
   - Review and edit if needed

4. **Add landing page URL**:
   - Will be provided by Zenodo after publication

5. **Publish**:
   - Review all fields
   - Click "Publish"
   - Get your DOI!

### Dryad Submission Workflow

1. **Prepare metadata**:
   ```bash
   cicada doi prepare data.fastq.gz \
     --enrich enrichment.yaml \
     --publisher "Your Institution" \
     --format json \
     --output datacite.json
   ```

2. **Create Dryad submission**:
   - Go to https://datadryad.org
   - Click "Submit Data"
   - Upload files

3. **Fill metadata form**:
   - Copy information from `datacite.json`
   - Dryad web form is user-friendly

4. **Submit for curation**:
   - Dryad curators review submission
   - May request improvements

5. **Publication**:
   - After approval, DOI is minted
   - Dataset becomes publicly available

---

## Complete Workflow Examples

### Example 1: Sequencing Data for Publication

```bash
# Step 1: Extract base metadata
cicada metadata extract sample_R1.fastq.gz --output base_metadata.json

# Step 2: Check initial quality
cicada doi validate sample_R1.fastq.gz

# Output: Quality Score: 47.0/100 - need enrichment

# Step 3: Create enrichment file
cat > enrichment.yaml <<'EOF'
title: "RNA-seq analysis of drug resistance in cancer cells"
authors:
  - name: Dr. Sarah Chen
    orcid: 0000-0002-1234-5678
    affiliation: Cancer Biology Department, State University
description: |
  RNA-seq data from drug-resistant cancer cell lines...
  [detailed methodology]
keywords: [RNA-seq, cancer, drug resistance, cell lines]
funding_references:
  - funder_name: National Cancer Institute
    award_number: CA123456
related_identifiers:
  - identifier: "10.1101/2025.123456"
    relation: IsSupplementTo
    type: DOI
EOF

# Step 4: Prepare DOI
cicada doi prepare sample_R1.fastq.gz \
  --enrich enrichment.yaml \
  --publisher "State University Genomics Core" \
  --format json \
  --output datacite.json

# Output: Quality Score: 91.0/100 - Ready!

# Step 5: Submit to Zenodo
# Upload files and datacite.json via web interface

# Step 6: Get DOI and cite in paper
# DOI: 10.5281/zenodo.123456
```

### Example 2: Batch DOI Preparation

```bash
#!/bin/bash
# prepare_dois.sh - Prepare DOI metadata for multiple files

ENRICHMENT="project_metadata.yaml"
PUBLISHER="University Research Lab"

for file in data/*.fastq.gz; do
  basename=$(basename "$file" .fastq.gz)

  echo "Processing: $basename"

  # Prepare DOI
  cicada doi prepare "$file" \
    --enrich "$ENRICHMENT" \
    --publisher "$PUBLISHER" \
    --format json \
    --output "doi_metadata/${basename}.datacite.json"

  if [ $? -eq 0 ]; then
    echo "  ✓ Ready for DOI"
  else
    echo "  ✗ Failed - check metadata"
  fi
done

echo "DOI preparation complete"
echo "Upload datacite.json files to repository"
```

### Example 3: Incremental Improvement

```bash
# Start with minimal metadata
cat > basic.yaml <<EOF
title: "My Dataset"
authors:
  - name: John Doe
description: "Research data"
EOF

cicada doi prepare data.fastq.gz --enrich basic.yaml --publisher "Lab"
# Score: 68/100 - needs improvement

# Add ORCIDs and better description
cat > improved.yaml <<EOF
title: "Whole genome sequencing of E. coli strains"
authors:
  - name: Dr. John Doe
    orcid: 0000-0002-1234-5678
    affiliation: Microbiology Department
description: |
  Complete methodology:
  - Sample collection and preparation
  - Sequencing platform and parameters
  - Data processing pipeline
keywords: [genomics, bacteria, sequencing]
EOF

cicada doi prepare data.fastq.gz --enrich improved.yaml --publisher "Lab"
# Score: 85/100 - much better!

# Add funding and related work
cat > complete.yaml <<EOF
title: "Whole genome sequencing of E. coli strains"
authors:
  - name: Dr. John Doe
    orcid: 0000-0002-1234-5678
    affiliation: Microbiology Department
description: |
  [Complete methodology as above]
keywords: [genomics, bacteria, sequencing, comparative genomics]
funding_references:
  - funder_name: National Institutes of Health
    award_number: R01GM123456
related_identifiers:
  - identifier: "10.1234/journal.2025.789"
    relation: IsSupplementTo
    type: DOI
version: "1.0"
EOF

cicada doi prepare data.fastq.gz --enrich complete.yaml --publisher "Lab"
# Score: 94/100 - excellent!
```

---

## Troubleshooting

### Low Quality Score

**Problem**: Score is below 60

**Solutions**:
1. Check for required fields:
   ```bash
   cicada doi validate file.fastq.gz --format json | jq '.errors'
   ```

2. Add missing required fields to enrichment file

3. Validate again to see improvement

### Author ORCID Issues

**Problem**: ORCIDs not recognized or invalid

**Solutions**:
1. Verify ORCID format: `0000-0002-1234-5678` (16 digits with dashes)

2. Check ORCID is real: https://orcid.org/0000-0002-1234-5678

3. Use correct field in enrichment:
   ```yaml
   authors:
     - name: Dr. Jane Smith
       orcid: 0000-0002-1234-5678  # No "https://" prefix
   ```

### Validation Fails

**Problem**: `Error: validation failed: 1 errors`

**Cause**: Required field is missing or invalid

**Solution**:
1. Read error message carefully
2. Check enrichment file has the field
3. Verify field format is correct
4. Re-run preparation

### Publisher Not Set

**Problem**: `Error: publisher is required`

**Solution**: Always provide `--publisher` flag:
```bash
cicada doi prepare file.fastq.gz \
  --enrich metadata.yaml \
  --publisher "Your Institution Name"
```

### Related Identifier Format

**Problem**: Related identifier rejected

**Solution**: Use correct format:
```yaml
related_identifiers:
  - identifier: "10.1234/journal.2025.123"  # DOI without https://
    relation: IsSupplementTo  # Use DataCite relation types
    type: DOI  # Use DataCite identifier types
```

**Valid relation types**:
- `IsSupplementTo` - This dataset supplements a publication
- `IsPartOf` - This dataset is part of a larger collection
- `IsCitedBy` - This dataset is cited by a publication
- `Cites` - This dataset cites another work
- See [DataCite documentation](https://datacite-metadata-schema.readthedocs.io/) for full list

---

## Best Practices

### 1. Start Early

Begin DOI preparation during data collection:
- Draft enrichment file as you work
- Easier to document while fresh in mind
- Can identify missing information early

### 2. Use Version Control

Track enrichment files in git:
```bash
git add enrichment.yaml
git commit -m "Add DOI metadata for dataset"
```

### 3. Reuse Enrichment Templates

Create templates for your lab:
```yaml
# lab_template.yaml
publisher: State University Genomics Core
license: CC-BY-4.0
funding_references:
  - funder_name: National Science Foundation
    award_number: NSF-XXXXX  # Update per project
```

### 4. Validate Before Submission

Always validate before repository submission:
```bash
cicada doi prepare data.fastq.gz \
  --enrich metadata.yaml \
  --publisher "Lab" \
  --min-score 80  # Require high quality
```

### 5. Document Everything

Include in your enrichment:
- Detailed methodology
- Data processing steps
- Quality control procedures
- Known limitations
- Related publications

---

## Next Steps

- **Extract metadata**: See [METADATA_EXTRACTION.md](METADATA_EXTRACTION.md)
- **Use presets**: See [PRESETS.md](PRESETS.md)
- **Provider setup**: See [PROVIDERS.md](PROVIDERS.md)
- **User scenarios**: See [USER_SCENARIOS_v0.2.0.md](USER_SCENARIOS_v0.2.0.md)

---

## Support

**Questions or Issues?**
- 📖 Full documentation: [README.md](../README.md)
- 🐛 Report bugs: [GitHub Issues](https://github.com/scttfrdmn/cicada/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/scttfrdmn/cicada/discussions)
