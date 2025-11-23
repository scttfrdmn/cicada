# Cicada v0.2.0 - Metadata & Intelligence

**Target**: Q1 2026
**Focus**: Instrument-aware data management with metadata extraction and validation
**Project Management**: [GitHub Projects](https://github.com/scttfrdmn/cicada/projects)

---

## Vision

Cicada is a **dormant data commons platform for academic research labs**. Version 0.1.0 established the foundational storage and sync layer. Version 0.2.0 builds on that foundation by adding **instrument awareness and metadata intelligence** - critical capabilities for a true data commons.

This release transforms Cicada into an intelligent research data pipeline that:
- Automatically extracts and preserves instrument metadata during data ingestion
- Understands diverse scientific instruments through pluggable extractors
- Validates file integrity and metadata completeness
- Enables FAIR-compliant data management with rich, searchable metadata
- Supports DOI minting for data publication (DataCite/Zenodo)

These capabilities move Cicada closer to the full data commons vision: federated storage with metadata, access control, compute-to-data, collaboration primitives, and data publication.

---

## Core Features

### 1. Metadata Extraction Framework

**Goal**: Automatically extract and preserve instrument metadata during sync operations.

#### 1.1 Metadata Extractor Interface

```go
type Extractor interface {
    // CanHandle returns true if extractor supports this file
    CanHandle(filename string) bool

    // Extract extracts metadata from file path
    Extract(filepath string) (*Metadata, error)

    // ExtractFromReader extracts from stream (for S3→Local)
    ExtractFromReader(r io.Reader, filename string) (*Metadata, error)

    // Validate checks file integrity
    Validate(filepath string) error

    // Name returns extractor name
    Name() string

    // SupportedFormats returns file extensions
    SupportedFormats() []string
}

type Metadata struct {
    // Common fields (all instruments)
    Format           string                 `json:"format"`
    InstrumentType   string                 `json:"instrument_type"`
    InstrumentModel  string                 `json:"instrument_model"`
    Manufacturer     string                 `json:"manufacturer"`
    AcquisitionDate  time.Time              `json:"acquisition_date"`
    Operator         string                 `json:"operator,omitempty"`
    FileSize         int64                  `json:"file_size"`
    Checksum         string                 `json:"checksum"`

    // Instrument-specific fields
    Microscopy       *MicroscopyMetadata    `json:"microscopy,omitempty"`
    Sequencing       *SequencingMetadata    `json:"sequencing,omitempty"`
    MassSpec         *MassSpecMetadata      `json:"mass_spec,omitempty"`
    FlowCytometry    *FlowCytometryMetadata `json:"flow_cytometry,omitempty"`

    // Custom fields
    Custom           map[string]interface{} `json:"custom,omitempty"`

    // Extraction metadata
    ExtractedAt      time.Time              `json:"extracted_at"`
    ExtractorVersion string                 `json:"extractor_version"`
}

type MicroscopyMetadata struct {
    Objective        ObjectiveInfo          `json:"objective"`
    Channels         []ChannelInfo          `json:"channels"`
    Dimensions       ImageDimensions        `json:"dimensions"`
    PixelSize        PixelSize              `json:"pixel_size"`
    ExposureTimes    []float64              `json:"exposure_times_ms"`
    IlluminationType string                 `json:"illumination_type"` // confocal, widefield, etc.
}

type ObjectiveInfo struct {
    Magnification float64 `json:"magnification"`
    NA            float64 `json:"numerical_aperture"`
    Immersion     string  `json:"immersion"` // oil, water, air
    WorkingDist   float64 `json:"working_distance_um,omitempty"`
}

type ChannelInfo struct {
    Name             string  `json:"name"`
    Wavelength       int     `json:"wavelength_nm"`
    ExposureTime     float64 `json:"exposure_time_ms"`
    EmissionFilter   string  `json:"emission_filter,omitempty"`
    ExcitationFilter string  `json:"excitation_filter,omitempty"`
}

type ImageDimensions struct {
    Width  int `json:"width"`
    Height int `json:"height"`
    Depth  int `json:"depth"`  // Z-stack
    Time   int `json:"time"`   // Time series
    Channels int `json:"channels"`
}

type PixelSize struct {
    X float64 `json:"x_um"`
    Y float64 `json:"y_um"`
    Z float64 `json:"z_um,omitempty"`
}

type SequencingMetadata struct {
    Platform       string `json:"platform"`        // Illumina, Nanopore, PacBio
    Instrument     string `json:"instrument"`      // NovaSeq 6000, MinION, etc.
    RunID          string `json:"run_id"`
    FlowcellID     string `json:"flowcell_id"`
    ReadLength     int    `json:"read_length"`
    IsPaired       bool   `json:"is_paired"`
    ReadCount      int64  `json:"read_count,omitempty"`
    QualityEncoding string `json:"quality_encoding"` // Phred+33, Phred+64
}
```

#### 1.2 Priority Extractors (v0.2.0)

**Microscopy**:
1. **Zeiss CZI** (Confocal microscopy)
   - Library: github.com/ome/bioformats (via CGo) or custom parser
   - Extracts: Instrument, objective, channels, dimensions, pixel size
   - Validation: CZI file structure integrity

2. **OME-TIFF** (Open Microscopy Environment)
   - Library: encoding/xml for OME-XML parsing
   - Extracts: Full OME metadata specification
   - Validation: OME-XML schema compliance

**Sequencing**:
3. **FASTQ** (Raw sequencing data)
   - Library: Custom parser
   - Extracts: Read count, quality encoding, sequence length distribution
   - Validation: FASTQ format integrity, quality score range

**Future** (v0.3.0+):
- Nikon ND2
- Leica LIF
- BAM/CRAM
- mzML (mass spec)
- FCS (flow cytometry)

#### 1.3 CLI Integration

```bash
# Enable metadata extraction during sync
cicada sync --extract-metadata /data/microscope s3://lab-data/

# Extract metadata from existing files
cicada metadata extract /data/microscope/*.czi

# Show metadata for a file
cicada metadata show /data/microscope/sample_001.czi

# Validate file integrity
cicada metadata validate /data/microscope/*.czi

# Export metadata catalog
cicada metadata export s3://lab-data/ --format json > catalog.json
```

#### 1.4 Storage Options

**Option A: S3 Object Metadata** (Default)
- Store as S3 object tags (10 tags max)
- Queryable via S3 API
- Limitations: Limited to 10 key-value pairs

**Option B: Sidecar JSON Files**
- Store as `.metadata.json` alongside each file
- Full metadata preservation
- Easy to read/process

**Option C: Central Catalog**
- Store in S3 as `s3://bucket/.cicada/metadata-catalog.json`
- Searchable index of all files
- Enables cross-file queries

**Recommendation**: Hybrid approach
- Critical fields → S3 object tags (instrument, date, format)
- Full metadata → Sidecar JSON
- Searchable index → Central catalog (updated incrementally)

---

### 2. Instrument Presets System

**Goal**: Simplify configuration for common lab instruments with pre-configured settings.

#### 2.1 Preset Definition

```yaml
# presets/zeiss-confocal.yaml
name: "Zeiss Confocal Microscope"
description: "Zeiss LSM series confocal microscopes (LSM 880, 900, 980)"
version: "1.0"

# Instrument detection
detection:
  file_extensions:
    - .czi
  file_patterns:
    - "*.czi"
  magic_bytes:
    offset: 0
    bytes: "ZISRAWFILE"  # CZI file signature

# Sync configuration
sync:
  debounce_seconds: 30      # CZI files are large, need time to write
  min_age_seconds: 60       # Wait for complete write
  concurrency: 4            # Balance speed vs system load

  # Suggested exclude patterns
  exclude_patterns:
    - "*.tmp"
    - "*.partial"
    - "*_preview.jpg"       # Zeiss preview images
    - "Experiment.czexp"    # Experiment metadata file

# Metadata extraction
metadata:
  enabled: true
  extractor: "zeiss-czi"
  extract_on_sync: true

  # Fields to extract
  fields:
    - instrument_model
    - objective
    - channels
    - dimensions
    - pixel_size
    - acquisition_date
    - operator

  # Fields to store as S3 tags (max 10)
  s3_tags:
    - instrument_type: microscopy
    - instrument_model: "{metadata.instrument_model}"
    - format: czi
    - acquisition_date: "{metadata.acquisition_date}"
    - magnification: "{metadata.microscopy.objective.magnification}x"

# Validation
validation:
  enabled: true
  checks:
    - file_integrity      # Verify CZI structure
    - minimum_file_size: 1048576  # 1 MB minimum
    - complete_metadata   # Ensure required metadata present

# Notifications (optional)
notifications:
  on_sync_complete:
    message: "Synced {file_count} CZI files ({total_size} GB)"
  on_validation_failure:
    message: "⚠️  Validation failed: {file_path}"
    action: skip  # skip, retry, fail
```

#### 2.2 Preset Library Structure

```
presets/
├── microscopy/
│   ├── zeiss-confocal.yaml
│   ├── zeiss-lightsheet.yaml
│   ├── nikon-widefield.yaml
│   ├── leica-confocal.yaml
│   └── olympus-spinning-disk.yaml
├── sequencing/
│   ├── illumina-novaseq.yaml
│   ├── illumina-miseq.yaml
│   ├── nanopore-minion.yaml
│   └── pacbio-sequel.yaml
├── mass-spec/
│   ├── thermo-orbitrap.yaml
│   └── bruker-maldi.yaml
├── flow-cytometry/
│   ├── bd-facs.yaml
│   └── beckman-cytoflex.yaml
└── generic/
    ├── large-files.yaml      # Generic preset for large files
    ├── many-small-files.yaml # Generic preset for many small files
    └── custom-template.yaml  # Template for users to customize
```

#### 2.3 CLI Commands

```bash
# List available presets
cicada instrument list

# Show preset details
cicada instrument show zeiss-confocal

# Interactive setup wizard
cicada instrument setup

# Apply specific preset
cicada instrument setup zeiss-confocal \
  --path /mnt/zeiss/output \
  --destination s3://lab-data/microscopy

# Auto-detect instrument type
cicada instrument detect /mnt/zeiss/output

# Create custom preset from existing watch
cicada instrument export my-custom-preset \
  --from-watch /mnt/zeiss/output-123456
```

#### 2.4 Auto-Detection Logic

```go
// Auto-detect instrument type from files
func DetectInstrument(path string) (*Preset, error) {
    // 1. Check file extensions
    files, _ := ioutil.ReadDir(path)
    extensions := countExtensions(files)

    // 2. Read magic bytes from sample files
    magicBytes := readMagicBytes(files[0:5])

    // 3. Match against preset detection rules
    for _, preset := range presets {
        if preset.Matches(extensions, magicBytes) {
            return preset, nil
        }
    }

    return nil, ErrNoPresetMatch
}
```

---

### 3. Pluggable DOI Provider System

**Goal**: Support multiple DOI providers (DataCite, Zenodo) or disable DOI minting entirely.

#### 3.1 Provider Interface

```go
// DOIProvider interface for pluggable DOI minting
type DOIProvider interface {
    // Name returns provider name (datacite, zenodo, etc.)
    Name() string

    // Mint creates a new DOI for a dataset
    Mint(ctx context.Context, dataset *Dataset) (*DOI, error)

    // Update updates DOI metadata
    Update(ctx context.Context, doi string, dataset *Dataset) error

    // Get retrieves DOI information
    Get(ctx context.Context, doi string) (*DOI, error)

    // Validate checks if dataset metadata is valid for this provider
    Validate(dataset *Dataset) error

    // EstimateCost returns estimated cost for minting (some providers charge)
    EstimateCost(dataset *Dataset) (float64, string, error)
}

// DOIConfig for provider selection
type DOIConfig struct {
    Provider     string                 `yaml:"provider"` // datacite, zenodo, off
    Enabled      bool                   `yaml:"enabled"`

    // DataCite specific
    DataCite     *DataCiteConfig        `yaml:"datacite,omitempty"`

    // Zenodo specific
    Zenodo       *ZenodoConfig          `yaml:"zenodo,omitempty"`

    // Common settings
    DefaultLicense   string             `yaml:"default_license"`
    DefaultPublisher string             `yaml:"default_publisher"`
    AutoPublish      bool               `yaml:"auto_publish"` // Auto-mint DOI on publish command
}

type DataCiteConfig struct {
    RepositoryID string  `yaml:"repository_id"`
    Password     string  `yaml:"password"`
    Prefix       string  `yaml:"prefix"`        // e.g., "10.12345"
    TestMode     bool    `yaml:"test_mode"`     // Use test environment
    BaseURL      string  `yaml:"base_url"`
}

type ZenodoConfig struct {
    AccessToken  string  `yaml:"access_token"`
    Sandbox      bool    `yaml:"sandbox"`       // Use sandbox environment
    Community    string  `yaml:"community"`     // Zenodo community ID
}
```

#### 3.2 Provider Implementations

**DataCite Provider** (Priority):
```go
type DataCiteProvider struct {
    client *DataCiteClient
    config *DataCiteConfig
}

func (p *DataCiteProvider) Mint(ctx context.Context, dataset *Dataset) (*DOI, error) {
    // 1. Validate dataset metadata
    if err := p.Validate(dataset); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Generate DOI suffix
    suffix := generateDOISuffix(dataset)
    doiString := fmt.Sprintf("%s/%s", p.config.Prefix, suffix)

    // 3. Prepare DataCite metadata XML
    metadata := buildDataCiteMetadata(dataset, doiString)

    // 4. Register with DataCite
    resp, err := p.client.CreateDOI(ctx, metadata)
    if err != nil {
        return nil, fmt.Errorf("datacite registration failed: %w", err)
    }

    return &DOI{
        DOI:       doiString,
        URL:       dataset.URL,
        State:     "draft",
        Provider:  "datacite",
        CreatedAt: time.Now(),
    }, nil
}
```

**Zenodo Provider** (Future):
```go
type ZenodoProvider struct {
    client *ZenodoClient
    config *ZenodoConfig
}

func (p *ZenodoProvider) Mint(ctx context.Context, dataset *Dataset) (*DOI, error) {
    // Zenodo auto-generates DOI
    // 1. Create deposition
    // 2. Upload files (optional)
    // 3. Publish to get DOI

    deposition, err := p.client.CreateDeposition(ctx, dataset)
    if err != nil {
        return nil, err
    }

    return &DOI{
        DOI:      deposition.DOI,
        URL:      deposition.URL,
        Provider: "zenodo",
    }, nil
}
```

**Disabled Provider**:
```go
type DisabledProvider struct{}

func (p *DisabledProvider) Mint(ctx context.Context, dataset *Dataset) (*DOI, error) {
    return nil, fmt.Errorf("DOI minting is disabled")
}

func (p *DisabledProvider) Name() string {
    return "disabled"
}
```

#### 3.3 Configuration

```yaml
# ~/.cicada/config.yaml
version: "1"

aws:
  profile: default
  region: us-west-2

# DOI configuration
doi:
  provider: datacite  # datacite, zenodo, off
  enabled: true

  # Provider-specific configs
  datacite:
    repository_id: "INST.LAB"
    password: "${DATACITE_PASSWORD}"  # From environment
    prefix: "10.12345"
    test_mode: false

  zenodo:
    access_token: "${ZENODO_TOKEN}"
    sandbox: false
    community: "my-institution"

  # Common settings
  default_license: "CC-BY-4.0"
  default_publisher: "Rodriguez Lab, University"
  auto_publish: false  # Require explicit publish command
```

#### 3.4 CLI Commands

```bash
# Configure DOI provider
cicada config set doi.provider datacite
cicada config set doi.datacite.repository_id INST.LAB
cicada config set doi.datacite.prefix 10.12345

# Disable DOI minting
cicada config set doi.provider off

# Publish dataset with DOI
cicada publish s3://lab-data/experiment-2025-11 \
  --title "Neuronal differentiation RNA-seq" \
  --authors "Maria Rodriguez, Alex Thompson" \
  --license CC-BY-4.0 \
  --description "RNA-seq time series of neuronal differentiation"

# Dry run (preview DOI metadata without minting)
cicada publish --dry-run s3://lab-data/experiment-2025-11 \
  --title "..."

# Update existing DOI
cicada publish --update 10.12345/exp.2025.789 \
  --add-author "James Park"

# List published datasets
cicada publish list

# Show DOI details
cicada publish show 10.12345/exp.2025.789
```

---

## Implementation Plan

### Phase 1: Core Infrastructure (Weeks 1-2)

1. **Metadata Package Refactor**
   - Define `Metadata` struct with common fields
   - Implement `Extractor` interface
   - Create registry system
   - Add basic validation framework

2. **Storage Layer**
   - S3 object tagging integration
   - Sidecar JSON writer
   - Central catalog management
   - Metadata query interface

### Phase 2: First Extractor - Zeiss CZI (Weeks 3-4)

1. **CZI Parser**
   - Research CZI file format
   - Evaluate libraries (bioformats, pylibCZI bindings)
   - Implement extractor
   - Add validation

2. **Testing**
   - Unit tests with sample CZI files
   - Integration tests with S3
   - Performance benchmarks

### Phase 3: Additional Extractors (Weeks 5-6)

1. **OME-TIFF Extractor**
   - XML parsing for OME metadata
   - Schema validation

2. **FASTQ Extractor**
   - Format parsing
   - Quality encoding detection
   - Read count estimation

### Phase 4: Instrument Presets (Weeks 7-8)

1. **Preset System**
   - YAML preset format
   - Preset loader
   - Auto-detection logic
   - CLI commands

2. **Initial Presets**
   - Zeiss confocal
   - Illumina sequencers
   - Generic templates

### Phase 5: DOI Provider System (Weeks 9-10)

1. **Provider Interface**
   - Define interface
   - Implement provider registry
   - Configuration system

2. **DataCite Provider**
   - API client
   - Metadata mapping
   - Error handling

3. **CLI Integration**
   - Publish commands
   - Interactive workflows

### Phase 6: Documentation & Testing (Weeks 11-12)

1. **Documentation**
   - Update USER_SCENARIOS with metadata examples
   - API documentation
   - Provider setup guides

2. **Integration Testing**
   - End-to-end workflows
   - Multi-instrument tests
   - DOI minting tests (test mode)

3. **Performance Testing**
   - Large file metadata extraction
   - Concurrent extraction performance
   - S3 tagging overhead

---

## Success Metrics

### Technical Metrics
- **Metadata extraction**: < 5s for typical files (< 1GB)
- **Sync overhead**: < 10% performance impact with metadata enabled
- **Validation accuracy**: 99%+ detection of corrupt files
- **Preset accuracy**: 95%+ correct auto-detection

### User Metrics
- **Setup time reduction**: 5 min → 1 min with presets
- **Metadata completeness**: 80%+ of required fields extracted
- **DOI minting success**: 99%+ successful mints
- **User satisfaction**: Positive feedback from beta users

---

## Testing Strategy

### Unit Tests
- Each extractor independently tested
- Sample files for each format
- Edge cases (corrupt files, partial metadata)

### Integration Tests
- Full sync with metadata extraction
- S3 storage of metadata
- DOI minting (test mode)

### User Acceptance Testing
- Beta test with 3-5 labs
- Different instrument types
- Real-world workloads

---

## Documentation Deliverables

1. **User Guide**: Metadata extraction and DOI minting
2. **Developer Guide**: Creating custom extractors
3. **Administrator Guide**: DOI provider setup
4. **Preset Library**: Documentation for each preset
5. **Updated USER_SCENARIOS**: Metadata-aware walkthroughs

---

## Risks & Mitigations

### Risk 1: CZI Parser Complexity
**Mitigation**: Use bioformats library (battle-tested) or partner with OME team

### Risk 2: S3 Tagging Limitations
**Mitigation**: Hybrid approach with sidecar files for full metadata

### Risk 3: DOI Provider Changes
**Mitigation**: Abstract interface, versioned API clients

### Risk 4: Performance Impact
**Mitigation**: Optional metadata extraction, async processing, caching

### Risk 5: Legal/Compliance
**Mitigation**: Clear licensing terms, terms of service for DOI minting

---

## Future Considerations (v0.3.0+)

- **More extractors**: Nikon ND2, Leica LIF, BAM, mzML, FCS
- **Metadata search**: Query interface for finding files by metadata
- **Batch operations**: Extract metadata from entire S3 buckets
- **Visualization**: Web UI for browsing metadata
- **AI/ML integration**: Automatic quality assessment, anomaly detection
- **Workflow integration**: Integration with pipeline tools (Nextflow, Snakemake)

---

## Success Criteria for v0.2.0 Release

- ✅ 3+ extractors fully implemented (CZI, OME-TIFF, FASTQ)
- ✅ Metadata extraction working in sync workflows
- ✅ 5+ instrument presets available
- ✅ DataCite provider fully functional
- ✅ Documentation complete
- ✅ 80%+ test coverage for new features
- ✅ Beta tested by 3+ labs
- ✅ No regressions in v0.1.0 functionality
