# Cicada v0.2.0 User Scenarios - Metadata & DOI Features

This document provides persona-based walkthroughs for **v0.2.0**, which adds automated metadata extraction, instrument preset validation, and DOI preparation capabilities to Cicada. These scenarios demonstrate how researchers can extract rich metadata from scientific instrument files, validate against instrument specifications, and prepare datasets for DOI minting.

## Table of Contents

1. [Postdoc: Metadata Extraction from Sequencing Data](#scenario-1-postdoc-metadata-extraction)
2. [Graduate Student: DOI Preparation for Publication](#scenario-2-graduate-student-doi-preparation)
3. [Lab Manager: Preset Validation](#scenario-3-lab-manager-preset-validation)
4. [Data Curator: Metadata Enrichment](#scenario-4-data-curator-metadata-enrichment)

---

## Scenario 1: Postdoc - Metadata Extraction from Sequencing Data

### Persona: Dr. Sarah Chen

**Background**:
- Postdoctoral researcher in genomics
- Generates ~200 FASTQ files per month from Illumina NovaSeq
- Needs to track read counts, quality scores, and sample information
- Must document metadata for data repository submission
- Technical level: Comfortable with command line, basic bioinformatics

**Pain Points**:
- Manually tracking file metadata is tedious and error-prone
- Need to extract quality metrics from FASTQ files
- Repository submission requires comprehensive metadata
- Lost track of which samples were already processed

**Goals**:
- Automatically extract metadata from FASTQ files
- Generate metadata reports for lab records
- Prepare data for repository submission
- Track metadata alongside uploaded files in S3

---

### Version Info
- ✅ **v0.2.0**: Metadata extraction and DOI features

---

### Day 1: First Metadata Extraction (5 minutes)

**Step 1: Install Cicada v0.2.0**

Sarah already has Cicada v0.1.0 installed, so she updates:

```bash
# Update to v0.2.0
brew upgrade cicada

# Verify version
cicada version
```

**Output**:
```
cicada version 0.2.0
  Go version: go1.23.4
  OS/Arch: darwin/arm64
```

**What Sarah thinks**: *"Great! Now I have the metadata features. Let me try extracting metadata from one of my FASTQ files."*

---

**Step 2: Extract Metadata from Single FASTQ File**

Sarah has a completed sequencing run:

```bash
# Check the file
ls -lh /data/sequencing/sample_A_R1.fastq.gz

# Extract metadata
cicada metadata extract /data/sequencing/sample_A_R1.fastq.gz
```

**Output**:
```json
{
  "format": "FASTQ",
  "compression": "gzip",
  "file_name": "/data/sequencing/sample_A_R1.fastq.gz",
  "file_size": 2147483648,
  "total_reads": 45623891,
  "total_bases": 6843583650,
  "mean_read_length": 150,
  "min_read_length": 150,
  "max_read_length": 150,
  "gc_content_percent": 42.3,
  "mean_quality_score": 36.8,
  "min_quality_score": 12,
  "max_quality_score": 41,
  "is_paired_end": true,
  "read_pair": "R1",
  "instrument_type": "sequencing",
  "data_type": "nucleotide_sequence",
  "extractor_name": "fastq",
  "schema_name": "fastq_v1"
}
```

**What Sarah thinks**: *"Wow! It extracted all the key metrics automatically - 45M reads with mean quality 36.8. This is exactly what I need for my records."*

---

**Step 3: Extract Metadata in Human-Readable Format**

```bash
# Get table format for easier reading
cicada metadata extract /data/sequencing/sample_A_R1.fastq.gz --format table
```

**Output**:
```
Metadata for: sample_A_R1.fastq.gz
================================

File Information:
  format                : FASTQ
  compression           : gzip
  file_size             : 2.0 GB
  extractor_name        : fastq
  schema_name           : fastq_v1

Sequence Statistics:
  total_reads           : 45,623,891
  total_bases           : 6,843,583,650
  mean_read_length      : 150
  min_read_length       : 150
  max_read_length       : 150
  gc_content_percent    : 42.3%

Quality Metrics:
  mean_quality_score    : 36.8
  min_quality_score     : 12
  max_quality_score     : 41

Sequencing Info:
  instrument_type       : sequencing
  data_type             : nucleotide_sequence
  is_paired_end         : true
  read_pair             : R1
```

**What Sarah thinks**: *"Perfect! This is much easier to read. Now let me extract metadata from all my files."*

---

### Day 1: Batch Metadata Extraction (10 minutes)

**Step 4: Extract Metadata from All FASTQ Files**

Sarah has 48 FASTQ files (24 samples, paired-end):

```bash
# Create metadata directory
mkdir -p /data/sequencing/metadata

# Extract metadata for each file and save to JSON
for file in /data/sequencing/*.fastq.gz; do
  basename=$(basename $file .fastq.gz)
  cicada metadata extract $file --format json --output /data/sequencing/metadata/${basename}.json
done

# Count processed files
ls /data/sequencing/metadata/*.json | wc -l
```

**Output**:
```
Metadata extracted to /data/sequencing/metadata/sample_A_R1.json
Metadata extracted to /data/sequencing/metadata/sample_A_R2.json
Metadata extracted to /data/sequencing/metadata/sample_B_R1.json
...
48
```

**What Sarah thinks**: *"Excellent! Now I have metadata for all 48 files. Let me create a summary report."*

---

**Step 5: Create Summary Statistics**

Sarah creates a simple script to summarize metadata:

```bash
#!/bin/bash
# summarize_metadata.sh

echo "Sequencing Run Summary"
echo "======================"
echo ""

total_reads=0
total_bases=0
count=0

for json in /data/sequencing/metadata/*.json; do
  reads=$(jq '.total_reads' $json)
  bases=$(jq '.total_bases' $json)

  total_reads=$((total_reads + reads))
  total_bases=$((total_bases + bases))
  count=$((count + 1))
done

echo "Total files: $count"
echo "Total reads: $(numfmt --grouping $total_reads)"
echo "Total bases: $(numfmt --grouping $total_bases)"
echo "Average reads per file: $(numfmt --grouping $((total_reads / count)))"
echo "Average quality: $(jq -s 'map(.mean_quality_score) | add / length' /data/sequencing/metadata/*.json)"
```

**Output**:
```
Sequencing Run Summary
======================

Total files: 48
Total reads: 2,189,946,768
Total bases: 328,492,015,200
Average reads per file: 45,623,891
Average quality: 36.7
```

**What Sarah thinks**: *"Perfect! I can include this summary in my lab notebook and grant reports."*

---

### Week 1: Validating with Instrument Presets (5 minutes)

**Step 6: Validate Against Illumina Preset**

Sarah wants to verify her FASTQ files meet Illumina NovaSeq specifications:

```bash
# Validate one file against Illumina preset
cicada metadata validate /data/sequencing/sample_A_R1.fastq.gz \
  --preset illumina-novaseq
```

**Output**:
```
✓ /data/sequencing/sample_A_R1.fastq.gz: valid (FASTQ)
     Quality Score: 100.0/100

Validation Results:
  Present Fields (8):
    ✓ format
    ✓ instrument_type
    ✓ data_type
    ✓ total_reads
    ✓ total_bases
    ✓ mean_quality_score
    ✓ is_paired_end
    ✓ read_pair

  Missing Optional Fields (0):
    (All optional fields present)

  Errors (0):
    No errors

Summary: Excellent quality metadata - all fields present
```

**What Sarah thinks**: *"Perfect! My data validates successfully against the Illumina NovaSeq preset. This confirms the metadata is complete."*

---

**Step 7: List Available Presets**

```bash
# See what other presets are available
cicada metadata preset list
```

**Output**:
```
Available Instrument Presets:

  Illumina NovaSeq
    ID: illumina-novaseq
    Manufacturer: Illumina
    Type: sequencing
    Models: NovaSeq 6000, NovaSeq X, NovaSeq X Plus
    Formats: .fastq, .fq, .fastq.gz, .fq.gz

  Illumina MiSeq
    ID: illumina-miseq
    Manufacturer: Illumina
    Type: sequencing
    Models: MiSeq
    Formats: .fastq, .fq, .fastq.gz, .fq.gz

  Illumina NextSeq
    ID: illumina-nextseq
    Manufacturer: Illumina
    Type: sequencing
    Models: NextSeq 500, NextSeq 550, NextSeq 1000, NextSeq 2000
    Formats: .fastq, .fq, .fastq.gz, .fq.gz

  Generic Sequencing
    ID: generic-sequencing
    Manufacturer: Various
    Type: sequencing
    Formats: .fastq, .fq, .fastq.gz, .fq.gz

  Zeiss LSM 880
    ID: zeiss-lsm-880
    Manufacturer: Zeiss
    Type: microscopy
    Models: LSM 880
    Formats: .czi

  Zeiss LSM 900
    ID: zeiss-lsm-900
    Manufacturer: Zeiss
    Type: microscopy
    Models: LSM 900
    Formats: .czi

  Zeiss LSM 980
    ID: zeiss-lsm-980
    Manufacturer: Zeiss
    Type: microscopy
    Models: LSM 980
    Formats: .czi

  Generic Microscopy
    ID: generic-microscopy
    Manufacturer: Various
    Type: microscopy
    Formats: .tif, .tiff, .czi, .nd2, .lif, .ome.tif, .ome.tiff

Total: 8 presets
```

**What Sarah thinks**: *"Great selection! These cover all the instruments in our lab."*

---

### Month 1: Key Benefits for Sarah

✅ **Automated extraction**: No more manual metadata tracking
✅ **Quality metrics**: Instant access to read counts, GC content, quality scores
✅ **Validation**: Verify metadata completeness against instrument specs
✅ **Batch processing**: Extract metadata from hundreds of files easily
✅ **Reporting**: Generate summaries for grants and publications

**What Sarah experiences**: *"Metadata extraction has saved me hours of manual work. I can now generate comprehensive reports in seconds, and I'm confident my data documentation is complete and accurate."*

---

## Scenario 2: Graduate Student - DOI Preparation for Publication

### Persona: Marcus Johnson

**Background**:
- 5th year PhD student in molecular biology
- About to publish first paper
- Journal requires data deposition with DOI
- Never registered a DOI before
- Technical level: Basic command line, learning as he goes

**Pain Points**:
- Don't know how to prepare data for DOI registration
- Not sure what metadata is required
- Worried about missing required fields
- Repository forms are confusing

**Goals**:
- Prepare sequencing data for DOI registration
- Ensure metadata meets DataCite requirements
- Get step-by-step guidance on what's needed
- Submit to repository with confidence

---

### Week 1: Understanding DOI Requirements (10 minutes)

**Step 1: Initial DOI Validation**

Marcus has raw sequencing data from his experiments:

```bash
# Check what's needed for DOI
cicada doi validate /data/experiment/sample_001.fastq.gz
```

**Output**:
```
DOI Validation Results
======================

File: sample_001.fastq.gz

✗ NOT READY for DOI minting

Quality Score: 47.0/100 (Moderate)

Present Fields (8):
  ✓ identifier (auto-generated)
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
  ✗ version
  ✗ language
  ✗ dates
  ✗ related_identifiers
  ✗ contributors
  ✗ funding_references
  ✗ geo_locations

Errors (1):
  • creator must be specified (currently set to 'Unknown Creator')

Warnings (5):
  • publisher should be specified (currently set to 'Unknown Publisher')
  • landing page URL is recommended
  • author ORCIDs are recommended for attribution
  • author affiliations are recommended
  • Quality score 47.0 is below minimum threshold 60.0

Recommendations:
  Fix all errors before minting DOI:
    - creator must be specified (currently set to 'Unknown Creator')
  Metadata quality is moderate. To improve:
    - Add missing recommended fields
    - Enhance description with methods and context
    - Include author ORCIDs for proper attribution
```

**What Marcus thinks**: *"Okay, I need to add author information and improve the description. Let me create an enrichment file."*

---

**Step 2: Create Metadata Enrichment File**

Marcus creates a file with additional metadata:

```bash
# Create enrichment file
cat > /data/experiment/enrichment.yaml <<EOF
title: "RNA-seq data for CRISPR knockout of gene X in human cell lines"

authors:
  - name: Marcus Johnson
    orcid: 0000-0002-1234-5678
    affiliation: Department of Molecular Biology, State University
  - name: Dr. Jennifer Lee
    orcid: 0000-0003-9876-5432
    affiliation: Department of Molecular Biology, State University
    role: supervisor

description: |
  This dataset contains RNA-seq data from CRISPR knockout experiments targeting
  gene X in HEK293 cells. Three biological replicates were sequenced for both
  control and knockout conditions using Illumina NovaSeq 6000. Reads were
  quality-filtered and aligned to the human genome (GRCh38).

  Experimental conditions:
  - Cell line: HEK293
  - Treatment: CRISPR knockout of gene X
  - Sequencing: Illumina NovaSeq 6000, 2x150bp paired-end
  - Biological replicates: 3 per condition

keywords:
  - RNA-seq
  - CRISPR
  - gene knockout
  - HEK293
  - transcriptomics

publisher: State University Genomics Core
license: CC-BY-4.0

funding_references:
  - funder_name: National Science Foundation
    award_number: NSF-1234567

related_identifiers:
  - identifier: "10.1234/journal.2025.123"
    relation: IsSupplementTo
    type: DOI
EOF
```

**What Marcus thinks**: *"I've added all my information. Let me see how this improves the DOI readiness."*

---

**Step 3: Prepare DOI with Enrichment**

```bash
# Prepare DOI with enrichment
cicada doi prepare /data/experiment/sample_001.fastq.gz \
  --enrich /data/experiment/enrichment.yaml \
  --publisher "State University Genomics Core"
```

**Output**:
```
DOI Preparation Results
=======================

File: sample_001.fastq.gz

Dataset Information:
  Title: RNA-seq data for CRISPR knockout of gene X in human cell lines
  Authors: 2
    1. Marcus Johnson (ORCID: 0000-0002-1234-5678)
       Department of Molecular Biology, State University
    2. Dr. Jennifer Lee (ORCID: 0000-0003-9876-5432)
       Department of Molecular Biology, State University
       Role: supervisor
  Publisher: State University Genomics Core
  Resource Type: Dataset
  Keywords: 5
  Related Publications: 1

Validation:
  ✓ Ready for DOI minting
  Quality Score: 85.0/100 (Excellent)
  Errors: 0
  Warnings: 1

Warnings:
  • landing page URL is recommended (will be added by repository)

Recommendations:
  Metadata quality is excellent. Ready for DOI registration.
  Optional improvements:
    - Add temporal coverage if applicable
    - Include spatial/geographical context if relevant
```

**What Marcus thinks**: *"Excellent! Quality score jumped from 47 to 85. I'm ready to submit this to the repository."*

---

**Step 4: Save DOI-Ready Metadata**

```bash
# Save prepared metadata to file for repository submission
cicada doi prepare /data/experiment/sample_001.fastq.gz \
  --enrich /data/experiment/enrichment.yaml \
  --publisher "State University Genomics Core" \
  --format json \
  --output /data/experiment/datacite_metadata.json
```

**Output**:
```
DOI preparation complete
Metadata saved to /data/experiment/datacite_metadata.json

This file contains DataCite-compliant metadata ready for repository submission.
```

**View the generated metadata**:
```bash
cat /data/experiment/datacite_metadata.json
```

**Output** (abbreviated):
```json
{
  "identifier": {
    "identifier": "dataset-2025-001",
    "identifierType": "local"
  },
  "creators": [
    {
      "name": "Johnson, Marcus",
      "nameType": "Personal",
      "givenName": "Marcus",
      "familyName": "Johnson",
      "nameIdentifiers": [
        {
          "nameIdentifier": "https://orcid.org/0000-0002-1234-5678",
          "nameIdentifierScheme": "ORCID"
        }
      ],
      "affiliation": ["Department of Molecular Biology, State University"]
    },
    {
      "name": "Lee, Jennifer",
      "nameType": "Personal",
      "givenName": "Jennifer",
      "familyName": "Lee",
      "nameIdentifiers": [
        {
          "nameIdentifier": "https://orcid.org/0000-0003-9876-5432",
          "nameIdentifierScheme": "ORCID"
        }
      ],
      "affiliation": ["Department of Molecular Biology, State University"]
    }
  ],
  "titles": [
    {
      "title": "RNA-seq data for CRISPR knockout of gene X in human cell lines"
    }
  ],
  "publisher": "State University Genomics Core",
  "publicationYear": 2025,
  "resourceType": {
    "resourceType": "Dataset",
    "resourceTypeGeneral": "Dataset"
  },
  "subjects": [
    {"subject": "RNA-seq"},
    {"subject": "CRISPR"},
    {"subject": "gene knockout"},
    {"subject": "HEK293"},
    {"subject": "transcriptomics"}
  ],
  "fundingReferences": [
    {
      "funderName": "National Science Foundation",
      "awardNumber": "NSF-1234567"
    }
  ],
  "relatedIdentifiers": [
    {
      "relatedIdentifier": "10.1234/journal.2025.123",
      "relatedIdentifierType": "DOI",
      "relationType": "IsSupplementTo"
    }
  ],
  "descriptions": [
    {
      "description": "This dataset contains RNA-seq data from CRISPR knockout experiments...",
      "descriptionType": "Abstract"
    }
  ],
  "rightsList": [
    {
      "rights": "Creative Commons Attribution 4.0 International",
      "rightsURI": "https://creativecommons.org/licenses/by/4.0/",
      "rightsIdentifier": "CC-BY-4.0"
    }
  ]
}
```

**What Marcus thinks**: *"Perfect! This metadata file is exactly what I need for repository submission. Now I can upload to Zenodo or Dryad with confidence."*

---

### Week 2: Repository Submission

**Step 5: Upload to Repository**

Marcus uploads his data to Zenodo using the prepared metadata:

1. **Upload files** to Zenodo via web interface
2. **Import metadata** from `datacite_metadata.json`
3. **Add landing page URL** (provided by Zenodo)
4. **Publish** to get DOI

**Result**: Dataset published with DOI `10.5281/zenodo.1234567`

---

### Month 1: Key Benefits for Marcus

✅ **Validation**: Know exactly what metadata is required
✅ **Guidance**: Clear recommendations for improvement
✅ **Quality scoring**: Track metadata completeness (47 → 85%)
✅ **DataCite compliance**: Automated mapping to repository schema
✅ **Confidence**: Submit knowing metadata is complete and correct

**What Marcus experiences**: *"Cicada made DOI preparation so much easier. Instead of guessing what metadata I needed, I got clear guidance every step of the way. My data is now properly documented and citable."*

---

## Scenario 3: Lab Manager - Preset Validation

### Persona: Dr. Emily Rodriguez

**Background**:
- Lab manager for multi-instrument core facility
- Manages Zeiss microscopes, Illumina sequencers, flow cytometer
- 50+ users generating diverse data types
- Needs to ensure consistent metadata quality
- Technical level: System administrator, scripting experience

**Pain Points**:
- Users don't document metadata consistently
- Hard to verify instrument-specific metadata is complete
- Need to enforce metadata standards across instruments
- Manual validation is time-consuming

**Goals**:
- Automated validation of instrument metadata
- Enforce metadata standards for each instrument
- Generate compliance reports for users
- Integrate validation into data upload workflow

---

### Week 1: Setting Up Preset Validation (15 minutes)

**Step 1: Create Validation Script**

Emily creates a script to validate all data before S3 upload:

```bash
#!/bin/bash
# validate_before_upload.sh
# Validates metadata before uploading to S3

set -euo pipefail

FILE="$1"
INSTRUMENT_TYPE="$2"

echo "=== Validating ${FILE} ==="

# Determine preset based on instrument
case "${INSTRUMENT_TYPE}" in
  "zeiss-880")
    PRESET="zeiss-lsm-880"
    ;;
  "zeiss-900")
    PRESET="zeiss-lsm-900"
    ;;
  "illumina-novaseq")
    PRESET="illumina-novaseq"
    ;;
  "illumina-miseq")
    PRESET="illumina-miseq"
    ;;
  *)
    echo "Unknown instrument type: ${INSTRUMENT_TYPE}"
    exit 1
    ;;
esac

echo "Using preset: ${PRESET}"

# Validate metadata
cicada metadata validate "${FILE}" --preset "${PRESET}"

if [ $? -eq 0 ]; then
  echo "✓ Validation passed - ready for upload"
  exit 0
else
  echo "✗ Validation failed - fix metadata before upload"
  exit 1
fi
```

**What Emily thinks**: *"This ensures all data is validated before upload. Users will get immediate feedback if metadata is incomplete."*

---

**Step 2: Test Validation with Good Data**

```bash
# Test with valid Illumina FASTQ
./validate_before_upload.sh \
  /data/illumina/sample_A_R1.fastq.gz \
  illumina-novaseq
```

**Output**:
```
=== Validating /data/illumina/sample_A_R1.fastq.gz ===
Using preset: illumina-novaseq

✓ /data/illumina/sample_A_R1.fastq.gz: valid (FASTQ)
     Quality Score: 100.0/100

Validation Results:
  Present Fields (8): All required fields present
  Missing Optional Fields (0): All optional fields present
  Errors (0): No errors

✓ Validation passed - ready for upload
```

**What Emily thinks**: *"Perfect! Validation passed. Now let me test with incomplete data."*

---

**Step 3: Test Validation with Incomplete Data**

```bash
# Test with file missing metadata
./validate_before_upload.sh \
  /data/illumina/incomplete.fastq \
  illumina-novaseq
```

**Output**:
```
=== Validating /data/illumina/incomplete.fastq ===
Using preset: illumina-novaseq

✗ /data/illumina/incomplete.fastq: invalid (FASTQ)
     Quality Score: 62.5/100

Validation Results:
  Present Fields (5):
    ✓ format
    ✓ instrument_type
    ✓ data_type
    ✓ total_reads
    ✓ total_bases

  Missing Required Fields (3):
    ✗ mean_quality_score
    ✗ is_paired_end
    ✗ read_pair

  Errors (3):
    • missing required field: mean_quality_score
    • missing required field: is_paired_end
    • missing required field: read_pair

✗ Validation failed - fix metadata before upload
```

**What Emily thinks**: *"Good! The validation caught the missing fields. Users will know exactly what needs to be fixed."*

---

**Step 4: Integrate with Upload Workflow**

Emily updates the Cicada watch configuration to include validation:

```yaml
# /etc/cicada/illumina-watch-with-validation.yaml

watches:
  - source: /data/illumina/output
    destination: s3://facility-data/illumina
    debounce: 60
    min_age: 120
    enabled: true
    hooks:
      pre_sync: |
        # Validate before upload
        for file in ${CHANGED_FILES}; do
          if [[ $file == *.fastq* ]]; then
            cicada metadata validate "$file" --preset illumina-novaseq || exit 1
          fi
        done
```

**What Emily thinks**: *"Now validation happens automatically before any upload. Users will get immediate feedback if their data has metadata issues."*

---

### Week 2: Generating Compliance Reports (10 minutes)

**Step 5: Create Validation Report Script**

```bash
#!/bin/bash
# generate_validation_report.sh
# Generate weekly metadata compliance report

OUTPUT_DIR="/data/reports"
mkdir -p "${OUTPUT_DIR}"

REPORT_FILE="${OUTPUT_DIR}/metadata_compliance_$(date +%Y-%m-%d).md"

cat > "${REPORT_FILE}" <<EOF
# Metadata Compliance Report
Generated: $(date)

## Illumina NovaSeq Files
EOF

# Validate all Illumina files
for file in /data/illumina/archive/*.fastq.gz; do
  echo "Validating: $(basename $file)"

  result=$(cicada metadata validate "$file" --preset illumina-novaseq 2>&1)

  if echo "$result" | grep -q "✓.*valid"; then
    echo "- ✓ $(basename $file): PASS" >> "${REPORT_FILE}"
  else
    echo "- ✗ $(basename $file): FAIL" >> "${REPORT_FILE}"
    echo "  - Issues: $(echo "$result" | grep "Error\|Warning" | head -3)" >> "${REPORT_FILE}"
  fi
done

cat >> "${REPORT_FILE}" <<EOF

## Zeiss LSM 880 Files
EOF

# Validate Zeiss files
for file in /data/zeiss/archive/*.czi; do
  echo "Validating: $(basename $file)"

  result=$(cicada metadata validate "$file" --preset zeiss-lsm-880 2>&1)

  if echo "$result" | grep -q "✓.*valid"; then
    echo "- ✓ $(basename $file): PASS" >> "${REPORT_FILE}"
  else
    echo "- ✗ $(basename $file): FAIL" >> "${REPORT_FILE}"
    echo "  - Issues: $(echo "$result" | grep "Error\|Warning" | head -3)" >> "${REPORT_FILE}"
  fi
done

echo "Report generated: ${REPORT_FILE}"
```

**Example Report Output**:
```markdown
# Metadata Compliance Report
Generated: Fri Nov 29 2025 10:00:00

## Illumina NovaSeq Files
- ✓ sample_001_R1.fastq.gz: PASS
- ✓ sample_001_R2.fastq.gz: PASS
- ✓ sample_002_R1.fastq.gz: PASS
- ✗ sample_003_R1.fastq.gz: FAIL
  - Issues: missing required field: mean_quality_score
- ✓ sample_004_R1.fastq.gz: PASS

## Zeiss LSM 880 Files
- ✓ experiment_001.czi: PASS
- ✓ experiment_002.czi: PASS
- ✗ experiment_003.czi: FAIL
  - Issues: missing required field: image_width
```

**What Emily thinks**: *"Perfect! Now I can send weekly compliance reports to users showing which files need metadata fixes."*

---

### Month 1: Key Benefits for Emily

✅ **Automated validation**: No more manual metadata checking
✅ **Instrument-specific**: Different presets for different instruments
✅ **Pre-upload validation**: Catch issues before data is uploaded
✅ **Compliance reporting**: Track metadata quality over time
✅ **User feedback**: Clear error messages guide users to fix issues

**What Emily experiences**: *"Preset validation has transformed our metadata quality. Users get immediate feedback when metadata is incomplete, and I can generate compliance reports effortlessly. Our data is now consistently well-documented."*

---

## Scenario 4: Data Curator - Metadata Enrichment

### Persona: Dr. Thomas Liu

**Background**:
- Data curator for university research data repository
- Receives datasets from 100+ research groups
- Responsible for DOI minting and long-term preservation
- Must ensure metadata meets repository standards
- Technical level: Expert in metadata standards, METS/MODS, Dublin Core, DataCite

**Pain Points**:
- Researchers submit datasets with minimal metadata
- Manual enrichment is time-consuming (15-30 min per dataset)
- Need to track metadata quality improvements
- Repository standards require comprehensive metadata

**Goals**:
- Streamline metadata enrichment workflow
- Ensure all datasets meet repository standards
- Track metadata quality scores
- Reduce time spent on metadata cleanup

---

### Week 1: Metadata Quality Assessment (15 minutes)

**Step 1: Initial Quality Check**

Thomas receives a new dataset submission:

```bash
# Assess initial metadata quality
cicada doi validate /repository/submissions/chen-2025-001/sample.fastq.gz
```

**Output**:
```
DOI Validation Results
======================

File: sample.fastq.gz

✗ NOT READY for DOI minting

Quality Score: 52.0/100 (Moderate)

Present Fields (10):
  ✓ identifier
  ✓ title
  ✓ publisher
  ✓ publication_year
  ✓ resource_type
  ✓ description (basic)
  ✓ license
  ✓ keywords (limited)
  ✓ format
  ✓ file_size

Missing Recommended Fields (7):
  ✗ author ORCIDs
  ✗ author affiliations
  ✗ funding information
  ✗ related publications
  ✗ temporal coverage
  ✗ methods description
  ✗ version information

Warnings (4):
  • description is minimal, consider adding methodology
  • only 2 keywords provided, recommend 5-10
  • author ORCIDs improve attribution and discovery
  • funding information recommended for compliance
```

**What Thomas thinks**: *"Score of 52/100. I need to enrich this with ORCIDs, funding info, and better description. Let me contact the researcher."*

---

**Step 2: Create Enrichment Template**

Thomas generates a template for the researcher to fill out:

```bash
# Generate enrichment template from existing metadata
cicada doi prepare /repository/submissions/chen-2025-001/sample.fastq.gz \
  --format yaml \
  --output /repository/submissions/chen-2025-001/enrichment_template.yaml
```

**Generated template** (abbreviated):
```yaml
title: "Existing title from file metadata"

authors:
  - name: "Unknown Creator"  # ← NEEDS UPDATE
    # orcid: "0000-0000-0000-0000"  # ← ADD ORCID
    # affiliation: "Institution Name"  # ← ADD AFFILIATION

description: |
  Basic description from file.

  # ← ADD: Detailed methodology
  # ← ADD: Experimental conditions
  # ← ADD: Data processing steps

keywords:
  - keyword1
  - keyword2
  # ← ADD: More keywords (5-10 total recommended)

# ← ADD: Funding information
# funding_references:
#   - funder_name: "Funding Agency"
#     award_number: "Grant-12345"

# ← ADD: Related publications
# related_identifiers:
#   - identifier: "10.1234/journal.2025.123"
#     relation: "IsSupplementTo"
#     type: "DOI"
```

**What Thomas thinks**: *"I'll send this template to Dr. Chen with clear instructions on what to add."*

---

**Step 3: Researcher Provides Enrichment**

Dr. Chen fills out the enrichment file:

```yaml
title: "Whole genome sequencing of antibiotic-resistant E. coli strains"

authors:
  - name: Dr. Sarah Chen
    orcid: 0000-0002-1234-5678
    affiliation: Department of Microbiology, State University
  - name: Dr. James Park
    orcid: 0000-0003-9876-5432
    affiliation: Department of Bioinformatics, State University

description: |
  This dataset contains whole genome sequencing data from 24 clinical isolates
  of antibiotic-resistant E. coli collected from hospital patients between
  2023-2024. Samples were sequenced using Illumina NovaSeq 6000 platform with
  2x150bp paired-end reads.

  Methodology:
  - DNA extraction: Qiagen DNeasy Blood & Tissue Kit
  - Library prep: Illumina DNA Prep
  - Sequencing: NovaSeq 6000, S4 flow cell
  - Coverage: Average 100x
  - Quality filtering: fastp v0.23.0 (Q30)
  - Assembly: SPAdes v3.15.5

  Data includes:
  - Raw FASTQ files (R1 and R2 for each sample)
  - Quality control reports
  - Assembly statistics

keywords:
  - whole genome sequencing
  - Escherichia coli
  - antibiotic resistance
  - antimicrobial resistance genes
  - genomic epidemiology
  - bacterial genomics
  - clinical isolates

funding_references:
  - funder_name: National Institutes of Health
    award_number: R01AI123456
  - funder_name: State University Research Foundation
    award_number: SURF-2024-789

related_identifiers:
  - identifier: "10.1234/journal.microbiology.2025.456"
    relation: IsSupplementTo
    type: DOI

temporal_coverage:
  start: "2023-01-01"
  end: "2024-12-31"

license: CC-BY-4.0
```

**What Thomas thinks**: *"Excellent enrichment! Now let me validate the improved metadata."*

---

**Step 4: Validate Enriched Metadata**

```bash
# Prepare DOI with enrichment
cicada doi prepare /repository/submissions/chen-2025-001/sample.fastq.gz \
  --enrich /repository/submissions/chen-2025-001/enrichment.yaml \
  --publisher "State University Research Data Repository"
```

**Output**:
```
DOI Preparation Results
=======================

File: sample.fastq.gz

Dataset Information:
  Title: Whole genome sequencing of antibiotic-resistant E. coli strains
  Authors: 2 (both with ORCIDs)
    1. Dr. Sarah Chen (ORCID: 0000-0002-1234-5678)
       Department of Microbiology, State University
    2. Dr. James Park (ORCID: 0000-0003-9876-5432)
       Department of Bioinformatics, State University
  Publisher: State University Research Data Repository
  Resource Type: Dataset
  Keywords: 7
  Funding: 2 grants
  Related Publications: 1
  Temporal Coverage: 2023-01-01 to 2024-12-31

Validation:
  ✓ Ready for DOI minting
  Quality Score: 94.0/100 (Excellent)
  Errors: 0
  Warnings: 0

Summary:
  Excellent metadata quality. All required and recommended fields present.
  Ready for DOI registration and publication.

Quality Improvement: +42 points (52 → 94)
```

**What Thomas thinks**: *"Perfect! Quality improved from 52 to 94. This is now ready for DOI minting and repository publication."*

---

**Step 5: Export DataCite Metadata**

```bash
# Export final DataCite-compliant metadata
cicada doi prepare /repository/submissions/chen-2025-001/sample.fastq.gz \
  --enrich /repository/submissions/chen-2025-001/enrichment.yaml \
  --publisher "State University Research Data Repository" \
  --format json \
  --output /repository/submissions/chen-2025-001/datacite_final.json
```

**What Thomas thinks**: *"Now I can upload this to DataCite and mint the DOI. The whole enrichment process took only 10 minutes instead of 30."*

---

### Month 1: Workflow Automation (20 minutes)

**Step 6: Create Automated Enrichment Workflow**

Thomas creates a script to streamline the process:

```bash
#!/bin/bash
# enrich_submission.sh
# Automated metadata enrichment workflow

SUBMISSION_ID="$1"
SUBMISSION_DIR="/repository/submissions/${SUBMISSION_ID}"
DATA_FILE="${SUBMISSION_DIR}/$(ls ${SUBMISSION_DIR}/*.{fastq,fastq.gz,czi,tif} 2>/dev/null | head -1)"

echo "=== Processing Submission: ${SUBMISSION_ID} ==="

# Step 1: Initial quality check
echo "Step 1: Assessing initial metadata quality..."
cicada doi validate "${DATA_FILE}" > "${SUBMISSION_DIR}/initial_quality.txt"

initial_score=$(grep "Quality Score:" "${SUBMISSION_DIR}/initial_quality.txt" | awk '{print $3}')
echo "Initial quality score: ${initial_score}"

# Step 2: Generate enrichment template
echo "Step 2: Generating enrichment template..."
cicada doi prepare "${DATA_FILE}" \
  --format yaml \
  --output "${SUBMISSION_DIR}/enrichment_template.yaml"

echo "Template saved to: ${SUBMISSION_DIR}/enrichment_template.yaml"
echo ""
echo "→ Please have researcher fill out enrichment_template.yaml"
echo "→ When complete, run: ./finalize_submission.sh ${SUBMISSION_ID}"
```

**Finalization script**:
```bash
#!/bin/bash
# finalize_submission.sh
# Finalize metadata and mint DOI

SUBMISSION_ID="$1"
SUBMISSION_DIR="/repository/submissions/${SUBMISSION_ID}"
DATA_FILE="${SUBMISSION_DIR}/$(ls ${SUBMISSION_DIR}/*.{fastq,fastq.gz,czi,tif} 2>/dev/null | head -1)"
ENRICHMENT="${SUBMISSION_DIR}/enrichment.yaml"

echo "=== Finalizing Submission: ${SUBMISSION_ID} ==="

# Validate enriched metadata
echo "Validating enriched metadata..."
cicada doi prepare "${DATA_FILE}" \
  --enrich "${ENRICHMENT}" \
  --publisher "State University Research Data Repository"

if [ $? -ne 0 ]; then
  echo "✗ Validation failed - check enrichment file"
  exit 1
fi

# Export DataCite metadata
echo "Exporting DataCite metadata..."
cicada doi prepare "${DATA_FILE}" \
  --enrich "${ENRICHMENT}" \
  --publisher "State University Research Data Repository" \
  --format json \
  --output "${SUBMISSION_DIR}/datacite.json"

echo "✓ Ready for DOI minting"
echo "DataCite metadata: ${SUBMISSION_DIR}/datacite.json"
```

**What Thomas thinks**: *"This workflow automates the repetitive parts while guiding researchers through enrichment. Time savings: 20 minutes per dataset."*

---

### Month 1: Key Benefits for Thomas

✅ **Quality scoring**: Track metadata improvements (52 → 94)
✅ **Template generation**: Auto-generate enrichment templates
✅ **Validation workflow**: Ensure completeness before DOI minting
✅ **Time savings**: Reduce enrichment time from 30 to 10 minutes
✅ **Compliance**: All datasets meet DataCite and repository standards

**What Thomas experiences**: *"Cicada's metadata enrichment workflow has doubled my throughput. I can now process 20-30 datasets per week instead of 10-15, and the quality is consistently higher. Researchers appreciate the clear guidance on what metadata to provide."*

---

## Summary: v0.2.0 Common Patterns

### All Personas Benefit From:

1. **Automated Metadata Extraction**
   - Extract comprehensive metadata from scientific files
   - Support for FASTQ, with more formats coming
   - Rich metadata: quality metrics, read counts, compression, pairing

2. **Instrument Preset Validation**
   - Validate against instrument-specific requirements
   - 8 default presets (Illumina, Zeiss, generic)
   - Clear error messages and recommendations

3. **DOI Preparation**
   - Assess DOI readiness with quality scoring
   - DataCite Schema v4.5 compliance
   - Metadata enrichment via YAML/JSON

4. **Quality Scoring**
   - 0-100 scale with clear thresholds
   - Track improvements over time
   - Identify missing fields

5. **Flexible Output**
   - JSON, YAML, or human-readable table
   - Save to files or stdout
   - Integration-ready formats

---

## Getting Started with v0.2.0

**Upgrade from v0.1.0**:
```bash
# macOS
brew upgrade cicada

# Linux
wget https://github.com/scttfrdmn/cicada/releases/download/v0.2.0/cicada_0.2.0_Linux_x86_64.tar.gz
tar -xzf cicada_0.2.0_Linux_x86_64.tar.gz
sudo mv cicada /usr/local/bin/
```

**Quick Start**:
```bash
# Extract metadata
cicada metadata extract your-file.fastq.gz

# Validate against preset
cicada metadata validate your-file.fastq.gz --preset illumina-novaseq

# Prepare for DOI
cicada doi prepare your-file.fastq.gz --enrich metadata.yaml
```

**Documentation**:
- [Metadata Extraction Guide](METADATA_EXTRACTION.md)
- [DOI Workflow Guide](DOI_WORKFLOW.md)
- [Instrument Preset Guide](PRESETS.md)
- [Provider Setup](PROVIDERS.md)
- [Integration Testing](../INTEGRATION_TESTING.md)

---

## Roadmap: What's Next

### v0.3.0 (Planned)
- Additional format extractors (CZI, OME-TIFF, TIFF)
- More instrument presets
- Custom preset creation
- Metadata versioning

### v0.4.0 (Planned)
- Direct DataCite/Zenodo API integration
- Automated DOI minting
- Metadata history tracking
- Repository integration

See [ROADMAP.md](../planning/ROADMAP.md) for full roadmap.
---

## Scenario 5: Small Lab - Complete Adoption Journey

### Lab Profile: Thompson & Kumar Lab

**Lab Composition**:
- **2 PIs**: Dr. Rachel Thompson (Cell Biology), Dr. Arun Kumar (Bioinformatics)
- **Research Staff**: 1 lab manager (Lisa), 1 research technician (David)
- **Postdoc**: Dr. Maria Santos (molecular imaging)
- **Graduate Students**: 3 PhD students (Alex, Jordan, Sam)
- **Undergraduate Students**: 2 research assistants (Taylor, Morgan)

**Infrastructure**:
- **Instruments**:
  - Zeiss LSM 880 confocal microscope (~3 TB/month)
  - Illumina MiSeq sequencer (~2 TB/month)
- **Current Data**: ~12 TB accumulated over 2 years
- **Current Workflow**:
  - Local storage on instrument workstations (rapidly filling up)
  - Manual copying to external drives for backup
  - Analysis on Terra/Dnanexus (PaaS with 100GB free tier)
  - Data regularly deleted to free space

**Pain Points**:
- 💾 **Storage Crisis**: Workstations at 95% capacity, deleting old data
- 🔄 **Manual Workflows**: Copying files manually is time-consuming
- 💰 **PaaS Costs**: Paying $500/month for Terra storage, but still limited
- 📊 **Poor Organization**: Data scattered across drives, hard to find files
- 🔍 **No Metadata**: Can't search for experiments, relying on folder names
- 👥 **Collaboration Issues**: Hard to share data between lab members
- 📄 **Publication Delays**: Scrambling to find/organize data for papers

**Budget**: $2,000/year for data management

**Goals**:
- Centralize all data in S3 (cheaper than PaaS storage)
- Automate instrument data upload
- Extract and track metadata
- Prepare datasets for DOI/publication
- Enable lab-wide data access
- Reduce manual data management time

---

### Month 1: Assessment and Planning

**Week 1: Lab Meeting - Understanding the Problem**

The lab meets to discuss their data crisis:

**Dr. Thompson**: *"We're at 95% capacity on both workstations. We've been deleting old data, but what if we need it later?"*

**Lisa (Lab Manager)**: *"I spend 3-4 hours per week manually copying files to external drives. We have 8 different drives now and I'm not even sure what's on half of them."*

**Alex (PhD Student)**: *"I can't find my imaging data from last year. I think it's on one of the external drives, but which one?"*

**Dr. Kumar**: *"We're paying $500/month for Terra storage but it's still not enough. AWS S3 would be much cheaper - about $30/month for 12 TB."*

**Decision**: Adopt Cicada for centralized S3 storage and automation.

---

**Week 2: Calculating Costs and Benefits**

Lisa creates a cost comparison:

**Current Costs (Annual)**:
- Terra storage: $500/month × 12 = **$6,000/year**
- External drives: $200 × 4/year = **$800/year**
- Lisa's manual labor: 4 hours/week × 50 weeks × $30/hour = **$6,000/year**
- **Total: $12,800/year**

**Projected Costs with Cicada (Annual)**:
- AWS S3 storage (15 TB): ~$350/year
- AWS data transfer: ~$200/year
- Cicada: Free (open source)
- Lisa's time reduced to 1 hour/week: $1,500/year
- **Total: $2,050/year**

**Savings: $10,750/year (84% reduction)**

**Dr. Kumar**: *"This pays for itself in the first month. Let's start implementation."*

---

### Month 2: Initial Implementation (v0.1.0 - Storage & Sync)

**Week 1: AWS Setup**

Lisa sets up AWS infrastructure:

```bash
# Create S3 bucket
aws s3 mb s3://thompson-kumar-lab

# Enable versioning (data recovery)
aws s3api put-bucket-versioning \
  --bucket thompson-kumar-lab \
  --versioning-configuration Status=Enabled

# Set up lifecycle policy (transition to Glacier after 1 year)
cat > lifecycle.json <<'LIFECYCLE'
{
  "Rules": [{
    "Id": "Archive old data",
    "Status": "Enabled",
    "Transitions": [{
      "Days": 365,
      "StorageClass": "GLACIER"
    }]
  }]
}
LIFECYCLE

aws s3api put-bucket-lifecycle-configuration \
  --bucket thompson-kumar-lab \
  --lifecycle-configuration file://lifecycle.json
```

**Lisa's Note**: *"Versioning protects against accidental deletions. Glacier archival will save us money on old data."*

---

**Week 2: Install Cicada on Instruments**

**Microscope Workstation (Ubuntu Linux)**:
```bash
# Install Cicada
wget https://github.com/scttfrdmn/cicada/releases/download/v0.1.0/cicada_0.1.0_Linux_x86_64.tar.gz
tar -xzf cicada_0.1.0_Linux_x86_64.tar.gz
sudo mv cicada /usr/local/bin/

# Initialize config
cicada config init
cicada config set aws.profile lab
cicada config set aws.region us-west-2

# Test sync
cicada sync --dry-run \
  /data/zeiss/2025-01/ \
  s3://thompson-kumar-lab/microscopy/2025-01/
```

**Sequencer Workstation (Windows with WSL)**:
```bash
# Same installation process
# Point to different S3 prefix
cicada sync --dry-run \
  /mnt/d/illumina/runs/ \
  s3://thompson-kumar-lab/sequencing/runs/
```

---

**Week 3: Migrate Historical Data**

Lisa runs a migration to upload all existing data:

```bash
#!/bin/bash
# migrate_historical_data.sh

echo "=== Historical Data Migration ==="

# Microscopy data (6 TB)
echo "Migrating microscopy data..."
cicada sync \
  /data/zeiss/archive/ \
  s3://thompson-kumar-lab/microscopy/archive/

# Sequencing data (6 TB)
echo "Migrating sequencing data..."
cicada sync \
  /mnt/d/illumina/archive/ \
  s3://thompson-kumar-lab/sequencing/archive/

# Check total uploaded
aws s3 ls s3://thompson-kumar-lab/ --recursive --summarize
```

**Result**: 12 TB migrated over 3 days (4 TB/day)

**Cost**: $0 for upload (AWS doesn't charge for data ingress)

---

**Week 4: Set Up Automatic Watching**

**Microscope - Auto-upload new imaging data**:
```bash
# Watch microscope output directory
cicada watch add \
  --debounce 30 \
  --min-age 60 \
  /data/zeiss/output \
  s3://thompson-kumar-lab/microscopy/live

# Create systemd service for persistence
sudo systemctl enable cicada-watch
sudo systemctl start cicada-watch
```

**Sequencer - Auto-upload new runs**:
```bash
# Watch sequencer output
cicada watch add \
  --debounce 60 \
  --min-age 300 \
  /mnt/d/illumina/output \
  s3://thompson-kumar-lab/sequencing/live

# Windows Task Scheduler for persistence
# (Run cicada watch start on login)
```

**Dr. Santos (Postdoc)**: *"I just finished an imaging session and the files are already in S3. I didn't have to do anything!"*

---

### Month 3: First Results and Cleanup

**Successes**:
✅ All 12 TB of historical data safely in S3
✅ Automatic upload working for both instruments
✅ No manual file copying needed
✅ Can delete local copies to free space (verified in S3 first)

**Storage Freed**:
- Microscope workstation: 5 TB freed (from 95% to 35% full)
- Sequencer workstation: 4 TB freed (from 90% to 40% full)
- External drives: 8 drives now in storage (no longer needed)

**Time Savings**:
- Lisa: 3 hours/week → 0.5 hours/week (monitoring only)
- **Saved: 2.5 hours/week = 125 hours/year = $3,750/year**

**Cost Reality Check**:
- Month 1 AWS bill: $35 (12 TB storage + retrieval testing)
- **Projected annual: ~$420 (vs $6,800 for Terra)**

**Lab Reaction**:

**Dr. Thompson**: *"This is transformative. We're not deleting data anymore, everyone can access files, and it's costing us less than one month of Terra storage."*

**Alex (PhD Student)**: *"I can finally find my old data! It's all organized by date and instrument in S3."*

---

### Month 4: Adding Metadata (v0.2.0)

**Week 1: Upgrade to v0.2.0**

Lisa upgrades both workstations:

```bash
# Microscope workstation
wget https://github.com/scttfrdmn/cicada/releases/download/v0.2.0/cicada_0.2.0_Linux_x86_64.tar.gz
tar -xzf cicada_0.2.0_Linux_x86_64.tar.gz
sudo mv cicada /usr/local/bin/
cicada version  # Verify v0.2.0
```

**What's New**: Metadata extraction, preset validation, DOI preparation

---

**Week 2: Extract Metadata from Sequencing Data**

Jordan (PhD student) needs to document their RNA-seq data:

```bash
# Extract metadata from all FASTQ files
for file in /mnt/d/illumina/project-A/*.fastq.gz; do
  basename=$(basename $file .fastq.gz)
  cicada metadata extract $file \
    --format json \
    --output /mnt/d/illumina/project-A/metadata/${basename}.json
done

# Create summary report
cat > summary.sh <<'SCRIPT'
#!/bin/bash
total_reads=0
for json in metadata/*.json; do
  reads=$(jq '.total_reads' $json)
  total_reads=$((total_reads + reads))
done
echo "Total reads: $(numfmt --grouping $total_reads)"
echo "Total files: $(ls metadata/*.json | wc -l)"
echo "Mean quality: $(jq -s 'map(.mean_quality_score) | add / length' metadata/*.json)"
SCRIPT

chmod +x summary.sh
./summary.sh
```

**Output**:
```
Total reads: 324,567,891
Total files: 24
Mean quality: 37.2
```

**Jordan**: *"Perfect! Now I have all the statistics for my methods section."*

---

**Week 3: Validate Microscopy Data**

Maria (postdoc) validates imaging data against Zeiss preset:

```bash
# Validate confocal imaging files
cicada metadata validate \
  /data/zeiss/experiment-001.czi \
  --preset zeiss-lsm-880
```

**Output**:
```
✓ experiment-001.czi: valid (CZI)
     Quality Score: 85.0/100

Validation Results:
  Present Fields (12): All required fields present
  Missing Optional Fields (3):
    • objective_na (recommended)
    • pixel_size_z (recommended for 3D imaging)
    • acquisition_date (helpful for tracking)

Recommendations:
  Good quality metadata. Consider adding:
    - Objective numerical aperture
    - Z-step size for 3D reconstructions
    - Acquisition timestamp
```

**Maria**: *"The metadata looks good, but I should add the objective NA and z-step size for better documentation."*

---

**Week 4: Create Lab-Wide Metadata Script**

Lisa creates a standard script for all lab members:

```bash
#!/bin/bash
# extract_and_upload_metadata.sh
# Extract metadata and upload alongside data

FILE="$1"

echo "Processing: $(basename $FILE)"

# Extract metadata
METADATA_FILE="${FILE}.metadata.json"
cicada metadata extract "$FILE" --format json --output "$METADATA_FILE"

# Upload both file and metadata
cicada sync "$FILE" s3://thompson-kumar-lab/data/
cicada sync "$METADATA_FILE" s3://thompson-kumar-lab/metadata/

echo "✓ Data and metadata uploaded"
```

**Usage**:
```bash
# Any lab member can use this
./extract_and_upload_metadata.sh experiment_001.fastq.gz
```

**Lab Adoption**: All 9 lab members now use this script for their experiments

---

### Month 6: Publication Preparation

**Scenario**: Sam (PhD student) is preparing first paper

**Week 1: Dataset Preparation**

Sam needs to prepare sequencing data for repository submission. Creates enrichment metadata file with:
- Title and detailed description
- Authors with ORCIDs and affiliations
- Keywords
- Funding information
- Related publication DOI (bioRxiv preprint)

---

**Week 2: DOI Preparation**

Sam prepares dataset for DOI registration:

```bash
# Check DOI readiness
cicada doi prepare /data/rnaseq/sample_001.fastq.gz \
  --enrich paper_dataset_enrichment.yaml \
  --publisher "State University Research Data Repository"
```

**Output**:
```
DOI Preparation Results
=======================

File: sample_001.fastq.gz

Dataset Information:
  Title: RNA-seq analysis of stress response in yeast under oxidative conditions
  Authors: 2 (both with ORCIDs)
  Quality Score: 91.0/100 (Excellent)

Summary: Ready for repository submission
```

**Sam**: *"Score of 91! I'm ready to submit to the repository."*

---

**Week 3: Export Metadata for Repository**

Sam exports DataCite metadata and uploads to Dryad, gets DOI, and cites in paper.

**Reviewers' Response**: *"Excellent data availability and documentation. Data is well-organized with comprehensive metadata."*

---

### Month 9: Lab-Wide Benefits Assessment

**Lab Meeting - 9-Month Review**

Lisa presents results:

**Data Management**:
- ✅ **15 TB** now in S3 (12 TB historical + 3 TB new data)
- ✅ **Zero data loss** (vs. multiple losses before)
- ✅ **100%** of new data automatically backed up within 1 hour
- ✅ **3 datasets** published with DOIs

**Cost Savings**:
| Item | Before | After | Savings |
|------|--------|-------|---------|
| Storage (Terra) | $6,000/yr | $0 | $6,000 |
| S3 Storage | $0 | $450/yr | -$450 |
| External drives | $800/yr | $0 | $800 |
| Labor (Lisa) | $6,000/yr | $1,500/yr | $4,500 |
| **Total** | **$12,800/yr** | **$1,950/yr** | **$10,850/yr** |

**ROI**: Saving **$10,850/year** (85% reduction)

---

### Year 2: Lab Transformation

**New Lab Culture**:
- Every experiment automatically documented with metadata
- All data backed up within minutes of generation
- Easy data discovery
- Confident publication submissions with DOI-ready datasets
- New students onboard in 15 minutes

**Unexpected Benefits**:
1. **Grant Writing**: Can easily quantify data generation
2. **Compliance**: NIH data sharing plan compliance is trivial
3. **Collaboration**: Easy to share data with external collaborators
4. **Student Training**: Students learn best practices
5. **Lab Reputation**: Known for excellent data documentation

---

## Key Takeaways: Small Lab Adoption

### Success Factors

1. **Start Small**: Begin with v0.1.0, add features gradually
2. **Clear ROI**: Calculate cost savings to justify adoption
3. **Lab Buy-in**: Get PIs and all members on board
4. **Standardize**: Create lab-wide scripts for common tasks
5. **Document**: Maintain simple guides for lab members

### Timeline

- **Month 1**: Planning and setup
- **Month 2**: v0.1.0 implementation (storage/sync)
- **Month 3**: Verify success, clean up local storage
- **Month 4**: Add v0.2.0 features (metadata/DOI)
- **Month 6**: First publication with DOI
- **Month 12**: Full lab transformation

### Costs (Annual)

| Phase | Cost | Benefit |
|-------|------|---------|
| **Before Cicada** | $12,800 | Constant data crisis |
| **After Cicada** | $1,950 | 15TB organized, automated, DOI-ready |
| **Savings** | **$10,850** | 85% reduction + peace of mind |

### Realistic Expectations

**What Cicada Does Today (v0.2.0)**:
- ✅ Automated S3 backup
- ✅ FASTQ metadata extraction
- ✅ Preset validation (Illumina, Zeiss)
- ✅ DOI preparation
- ✅ Multi-format output (JSON, YAML, table)

**What Requires Workarounds**:
- ⚠️ CZI/microscopy metadata (use manual documentation until v0.3.0)
- ⚠️ Direct DOI minting (use Dryad/Zenodo web interface with prepared metadata)
- ⚠️ Real-time collaboration (use presigned URLs until access control in v0.5.0)

**Key Message**: Even without every feature, Cicada provides immediate value. As new versions release, labs can adopt additional capabilities incrementally.

---

## Small Lab Adoption Checklist

### Prerequisites
- [ ] AWS account with S3 access
- [ ] 2-4 hours for initial setup
- [ ] Lab-wide agreement to adopt
- [ ] Basic command line comfort (at least one person)

### Month 1: Assessment
- [ ] Document current data volumes
- [ ] Calculate current costs (storage, labor)
- [ ] Identify pain points
- [ ] Present Cicada proposal to PIs

### Month 2: Implementation (v0.1.0)
- [ ] Set up AWS S3 bucket with versioning
- [ ] Install Cicada on instrument workstations
- [ ] Test sync with small dataset
- [ ] Migrate historical data
- [ ] Set up automatic watches
- [ ] Create systemd/launchd services

### Month 3: Validation
- [ ] Verify all data in S3
- [ ] Test data retrieval
- [ ] Delete local copies (keep local cache)
- [ ] Document lab procedures
- [ ] Train all lab members

### Month 4: Metadata (v0.2.0)
- [ ] Upgrade to v0.2.0
- [ ] Extract metadata from key datasets
- [ ] Create metadata extraction scripts
- [ ] Validate against presets
- [ ] Document metadata workflows

### Month 6: Publication
- [ ] Prepare dataset for DOI
- [ ] Create enrichment metadata
- [ ] Validate DOI readiness
- [ ] Submit to repository
- [ ] Update paper materials section

### Month 12: Review
- [ ] Calculate actual costs and savings
- [ ] Assess time savings
- [ ] Survey lab member satisfaction
- [ ] Plan for future Cicada versions
- [ ] Share success with department/university
