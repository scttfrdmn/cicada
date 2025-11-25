# Cicada Architecture

**Last Updated:** 2025-01-24

## Table of Contents

1. [Platform Overview](#platform-overview)
2. [System Architecture](#system-architecture)
3. [Component Architecture](#component-architecture)
4. [Data Flow](#data-flow)
5. [Storage Architecture](#storage-architecture)
6. [Metadata Architecture](#metadata-architecture)
7. [Extensibility Points](#extensibility-points)
8. [Design Principles](#design-principles)

---

## Platform Overview

### What is Cicada?

Cicada is a **small lab data commons platform** designed to provide academic research labs with cost-effective, dormant data management infrastructure. Like its namesake insect, Cicada lies dormant most of the time, consuming minimal resources, but emerges powerfully when needed for data operations.

### Core Capabilities

- **Storage & Sync**: Multi-backend storage with efficient bi-directional synchronization
- **Metadata Management**: Automated extraction from 14 scientific file formats
- **Watch Mode**: Automatic monitoring and synchronization of directories
- **Data Quality**: Validation using instrument-specific presets
- **DOI Preparation**: Optional DataCite integration for dataset publication (advanced feature)

### Target Users

- Academic research labs with limited technical expertise
- Small labs with tight budgets
- Research groups needing federated data storage
- Labs managing diverse scientific data formats

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI Layer (Cobra)                        │
│  cicada [sync|watch|metadata|doi|config|version]                │
└────────────┬────────────────────────────────────────────────────┘
             │
             ├─── Sync Commands ──────► Sync Engine
             ├─── Watch Commands ─────► Watch Manager
             ├─── Metadata Commands ──► Metadata System
             ├─── DOI Commands ───────► DOI Workflow
             └─── Config Commands ────► Configuration System
                         │
         ┌───────────────┴────────────────┐
         │                                 │
┌────────▼─────────┐            ┌─────────▼─────────┐
│  Storage Layer   │            │  Metadata Layer   │
│                  │            │                   │
│  ┌────────────┐  │            │  ┌─────────────┐ │
│  │   Local    │  │            │  │  Extractor  │ │
│  │  Backend   │  │            │  │  Registry   │ │
│  └────────────┘  │            │  └─────────────┘ │
│                  │            │                   │
│  ┌────────────┐  │            │  ┌─────────────┐ │
│  │    S3      │  │            │  │   Preset    │ │
│  │  Backend   │  │            │  │  Validator  │ │
│  └────────────┘  │            │  └─────────────┘ │
└──────────────────┘            └───────────────────┘
         │                                 │
         │                                 │
┌────────▼─────────┐            ┌─────────▼─────────┐
│  Watch Daemon    │            │   DOI Workflow    │
│                  │            │   (Optional)      │
│  ┌────────────┐  │            │                   │
│  │  Watcher   │  │            │  ┌─────────────┐ │
│  │  Manager   │  │            │  │  Provider   │ │
│  └────────────┘  │            │  │  Interface  │ │
│                  │            │  └─────────────┘ │
│  ┌────────────┐  │            │                   │
│  │ Debouncer  │  │            │  ┌─────────────┐ │
│  └────────────┘  │            │  │  Metadata   │ │
└──────────────────┘            │  │   Mapper    │ │
                                 │  └─────────────┘ │
                                 └───────────────────┘
```

### Package Structure

```
cicada/
├── cmd/cicada/           # Main CLI entry point
│   └── main.go           # Application bootstrap
│
├── internal/
│   ├── cli/              # CLI commands (Cobra)
│   │   ├── root.go       # Root command
│   │   ├── sync.go       # Sync commands
│   │   ├── watch.go      # Watch commands
│   │   ├── metadata.go   # Metadata commands
│   │   ├── doi.go        # DOI commands
│   │   ├── config.go     # Config commands
│   │   └── version.go    # Version command
│   │
│   ├── sync/             # Storage and sync engine
│   │   ├── backend.go    # Backend interface
│   │   ├── engine.go     # Sync engine
│   │   ├── local.go      # Local filesystem backend
│   │   └── s3.go         # S3 backend
│   │
│   ├── metadata/         # Metadata extraction system
│   │   ├── extractor.go  # Extractor interface & registry
│   │   ├── types.go      # Instrument-specific types
│   │   ├── preset.go     # Preset validation
│   │   ├── s3tags.go     # S3 metadata tagging
│   │   ├── schema.go     # Metadata schema
│   │   ├── fastq.go      # FASTQ extractor
│   │   ├── ome_tiff.go   # OME-TIFF extractor
│   │   ├── zeiss_czi.go  # Zeiss CZI extractor
│   │   └── ...           # Other extractors
│   │
│   ├── watch/            # Watch daemon
│   │   ├── manager.go    # Watch manager
│   │   ├── watcher.go    # File system watcher
│   │   ├── debouncer.go  # Event debouncing
│   │   └── config.go     # Watch configuration
│   │
│   ├── doi/              # DOI preparation (optional)
│   │   ├── workflow.go   # DOI workflow
│   │   ├── provider.go   # Provider interface
│   │   ├── mapper.go     # Metadata mapping
│   │   ├── validation.go # Validation rules
│   │   ├── datacite_provider.go
│   │   └── zenodo_provider.go
│   │
│   └── config/           # Configuration management
│       ├── config.go     # Config structures
│       └── providers.go  # Provider configuration
│
└── docs/                 # Documentation
```

---

## Component Architecture

### 1. CLI Layer

**Location:** `internal/cli/`

The CLI layer uses [Cobra](https://github.com/spf13/cobra) to provide a command-line interface. Commands are organized hierarchically:

```
cicada
├── sync          - Sync files between storage locations
├── watch         - Watch directories for changes
│   ├── start     - Start watching a directory
│   ├── stop      - Stop a watch
│   ├── list      - List active watches
│   └── status    - Show watch status
├── metadata      - Metadata operations
│   ├── extract   - Extract metadata from files
│   ├── show      - Display metadata
│   ├── validate  - Validate against presets
│   ├── list      - List available extractors
│   └── preset    - Manage presets
├── doi           - DOI preparation (optional)
│   ├── prepare   - Prepare DOI metadata
│   ├── validate  - Validate DOI metadata
│   └── submit    - Submit to provider
├── config        - Configuration management
│   ├── show      - Show current configuration
│   ├── set       - Set configuration values
│   └── init      - Initialize configuration
└── version       - Show version information
```

**Key Files:**
- `root.go` - Root command and global flags (internal/cli/root.go:35)
- `sync.go` - Sync command implementations
- `watch.go` - Watch command implementations
- `metadata.go` - Metadata command implementations

### 2. Storage Backends

**Location:** `internal/sync/`

Storage backends provide a unified interface for different storage types:

```go
// Backend interface (internal/sync/backend.go:34)
type Backend interface {
    List(ctx context.Context, prefix string) ([]FileInfo, error)
    Read(ctx context.Context, path string) (io.ReadCloser, error)
    Write(ctx context.Context, path string, r io.Reader, size int64) error
    Delete(ctx context.Context, path string) error
    Stat(ctx context.Context, path string) (*FileInfo, error)
    Close() error
}
```

**Implementations:**
1. **Local Backend** (`local.go`) - Local filesystem operations
2. **S3 Backend** (`s3.go`) - AWS S3 (and S3-compatible services)

**Planned:**
- Azure Blob Storage backend
- Google Cloud Storage backend

### 3. Sync Engine

**Location:** `internal/sync/engine.go`

The sync engine orchestrates data transfer between backends:

**Key Features:**
- Parallel transfers (configurable concurrency, default: 4)
- Dry-run mode for previewing changes
- Progress reporting via callbacks
- Delete synchronization (optional)
- Smart sync using ETag/checksum comparison
- Size and modification time comparison

**Sync Algorithm:**
1. List files from source and destination
2. Build destination file map for quick lookup
3. Determine files needing sync (missing, modified, or different size)
4. Identify files to delete (if delete mode enabled)
5. Perform parallel transfers using semaphore
6. Execute deletions sequentially

**Configuration:**
```go
type SyncOptions struct {
    DryRun       bool   // Preview without making changes
    Delete       bool   // Remove files not in source
    Concurrency  int    // Parallel transfer count
    ProgressFunc func(ProgressUpdate)
}
```

### 4. Metadata System

**Location:** `internal/metadata/`

The metadata system provides format-agnostic metadata extraction:

#### Extractor Registry

**Architecture:**
```go
// Extractor interface (internal/metadata/extractor.go:14)
type Extractor interface {
    CanHandle(filename string) bool
    Extract(filepath string) (map[string]interface{}, error)
    ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error)
    Name() string
    SupportedFormats() []string
}

// Registry (internal/metadata/extractor.go:32)
type ExtractorRegistry struct {
    extractors []Extractor
}
```

**Registered Extractors (internal/metadata/extractor.go:49):**
1. TIFF - Standard TIFF images
2. OME-TIFF - Open Microscopy Environment TIFF
3. Zeiss CZI - Zeiss confocal microscopy
4. Nikon ND2 - Nikon microscopy
5. Leica LIF - Leica microscopy
6. FASTQ - DNA sequencing reads
7. BAM - Binary alignment/map
8. mzML - Mass spectrometry
9. MGF - Mascot Generic Format
10. HDF5 - Hierarchical Data Format
11. Zarr - Chunked array storage
12. DICOM - Medical imaging
13. FCS - Flow cytometry
14. Generic - Fallback extractor

#### Instrument-Specific Metadata Types

**Location:** `internal/metadata/types.go`

Six specialized metadata structures for different scientific domains:

1. **MicroscopyMetadata** (types.go:24) - Confocal, widefield, TIRF, light sheet
   - Instrument info (manufacturer, model, software)
   - Imaging parameters (magnification, NA, pixel size)
   - Channel information (fluorophores, wavelengths)
   - Acquisition settings (exposure, binning, gain)

2. **SequencingMetadata** (types.go:89) - DNA/RNA sequencing
   - Platform info (Illumina, PacBio, ONT)
   - Run information (run ID, flowcell, lane)
   - Library details (kit, protocol, insert size)
   - Quality metrics (Q30, duplication rate)

3. **MassSpecMetadata** (types.go:140) - Mass spectrometry
   - Instrument type (Orbitrap, Q-TOF, Triple Quad)
   - Ionization mode (ESI, MALDI, APCI)
   - Acquisition parameters (resolution, scan rate)
   - Chromatography details (for LC-MS)

4. **FlowCytometryMetadata** (types.go:182) - Flow cytometry
   - Instrument and software info
   - Acquisition parameters (events, rate)
   - Channel/parameter details
   - Compensation matrix

5. **CryoEMMetadata** (types.go:226) - Cryo-electron microscopy
   - Microscope and detector details
   - Acquisition parameters (voltage, defocus, dose)
   - Movie/frame information

6. **XRayMetadata** (types.go:257) - X-ray crystallography
   - Facility and beamline
   - Acquisition parameters (wavelength, energy)
   - Crystal information (space group, unit cell)

#### Preset Validation

**Location:** `internal/metadata/preset.go`

Presets define required and recommended metadata fields for data quality validation:

```go
type Preset struct {
    Name             string
    Description      string
    RequiredFields   []string   // Must be present
    RecommendedFields []string  // Should be present (warnings)
    AllowedValues    map[string][]string
}
```

**Built-in Presets:**
- Microscopy-confocal
- Microscopy-widefield
- Sequencing-illumina
- Sequencing-pacbio
- Mass-spec-proteomics
- Mass-spec-metabolomics
- Flow-cytometry
- General-lab

### 5. Watch Daemon

**Location:** `internal/watch/`

The watch daemon monitors directories and automatically syncs changes:

**Components:**

1. **Watcher** (watcher.go) - File system monitoring using fsnotify
2. **Debouncer** (debouncer.go) - Event coalescing to prevent duplicate syncs
3. **Manager** (manager.go) - Multi-watch orchestration and persistence

**Watch Flow:**
```
File Change Event
    │
    ▼
┌─────────┐
│ Watcher │ (fsnotify)
└────┬────┘
     │
     ▼
┌─────────┐
│Debouncer│ (wait for quiet period)
└────┬────┘
     │
     ▼
┌─────────┐
│  Sync   │ (trigger sync engine)
└─────────┘
```

**Configuration:**
```go
type Config struct {
    Source          string        // Directory to watch
    Destination     string        // Sync destination
    DebounceDelay   time.Duration // Wait time after last event
    MinAge          time.Duration // Only sync files older than this
    DeleteSource    bool          // Delete after successful sync
    SyncOnStart     bool          // Initial sync when starting
    ExcludePatterns []string      // Files to exclude
}
```

**Persistence:**
Watch configurations are saved to `~/.cicada/config.yaml` and automatically restored on restart.

### 6. DOI Preparation System (Optional)

**Location:** `internal/doi/`

The DOI system prepares datasets for publication with permanent identifiers:

**Components:**

1. **Workflow** (workflow.go) - DOI preparation orchestration
2. **Provider Interface** (provider.go) - Abstraction for DOI services
3. **Metadata Mapper** (mapper.go) - Map Cicada metadata to DataCite schema
4. **Validation** (validation.go) - Validate DataCite compliance

**Provider Implementations:**
- DataCite (datacite_provider.go) - Direct DOI minting
- Zenodo (zenodo_provider.go) - Zenodo repository integration

**Note:** This is an advanced, optional feature. Most labs use Cicada for core data management without DOI preparation.

### 7. Configuration System

**Location:** `internal/config/`

Configuration is managed using [Viper](https://github.com/spf13/viper) and stored in YAML format:

**Configuration Structure:**
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

watches:
  - id: microscopy-data
    source: /data/microscopy
    destination: s3://lab-data/microscopy
    enabled: true
    debounce_seconds: 10
    min_age_seconds: 60

settings:
  verbose: false
  log_file: ~/.cicada/cicada.log
  check_updates: true
```

**Default Location:** `~/.cicada/config.yaml`

---

## Data Flow

### Sync Operation Flow

```
┌──────────┐
│   User   │
│  Command │
└────┬─────┘
     │ cicada sync /local/data s3://bucket/prefix
     ▼
┌────────────┐
│ CLI Parser │
└─────┬──────┘
      │
      ▼
┌─────────────────┐
│ Backend Factory │
└────┬────────┬───┘
     │        │
     ▼        ▼
┌────────┐ ┌────────┐
│ Local  │ │   S3   │
│Backend │ │Backend │
└───┬────┘ └───┬────┘
    │          │
    └────┬─────┘
         ▼
   ┌────────────┐
   │Sync Engine │
   └──────┬─────┘
          │
          ├─► List Source Files
          ├─► List Destination Files
          ├─► Compare (ETag, size, mtime)
          ├─► Determine Changes
          │
          ▼
   ┌──────────────┐
   │   Transfer   │ (parallel, configurable concurrency)
   │  Operations  │
   └──────┬───────┘
          │
          ▼
   ┌──────────────┐
   │   Progress   │
   │  Reporting   │
   └──────────────┘
```

### Metadata Extraction Flow

```
┌──────────┐
│   User   │
│  Command │
└────┬─────┘
     │ cicada metadata extract file.czi
     ▼
┌─────────────────┐
│ CLI Parser      │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│ Extractor       │
│ Registry        │
└────┬────────────┘
     │
     ├─► FindExtractor(file.czi)
     │   └─► ZeissCZIExtractor
     │
     ▼
┌─────────────────┐
│ Extract         │
│ Metadata        │
└────┬────────────┘
     │
     ├─► Parse File Format
     ├─► Extract Embedded Metadata
     ├─► Add File Information
     │
     ▼
┌─────────────────┐
│ Validate        │ (if preset specified)
│ Against Preset  │
└────┬────────────┘
     │
     ├─► Check Required Fields
     ├─► Check Recommended Fields
     ├─► Validate Field Values
     │
     ▼
┌─────────────────┐
│ Store Metadata  │
└────┬────────────┘
     │
     ├─► S3 Object Tags (if S3 backend)
     ├─► Sidecar JSON File
     │
     ▼
┌─────────────────┐
│ Output Results  │
└─────────────────┘
```

### Watch Mode Flow

```
┌──────────┐
│   User   │
│  Command │
└────┬─────┘
     │ cicada watch start microscopy /data/microscopy s3://bucket/data
     ▼
┌─────────────────┐
│ Watch Manager   │
└────┬────────────┘
     │
     ├─► Create Watcher
     ├─► Configure Debouncer
     ├─► Create Sync Engine
     │
     ▼
┌─────────────────┐
│ Initial Sync    │ (if sync_on_start: true)
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│ Start Watching  │
└────┬────────────┘
     │
     │ ┌─ File System Events ─┐
     ▼ │                       │
┌─────────────────┐            │
│   fsnotify      │────────────┘
│   Watcher       │
└────┬────────────┘
     │
     ├─► Filter Events (exclude patterns)
     ├─► Check File Age (min_age)
     │
     ▼
┌─────────────────┐
│   Debouncer     │ (wait for quiet period)
└────┬────────────┘
     │
     ├─► Collect Events
     ├─► Wait debounce_seconds
     ├─► Coalesce Duplicate Events
     │
     ▼
┌─────────────────┐
│ Trigger Sync    │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│  Sync Engine    │
└────┬────────────┘
     │
     ├─► Transfer Changed Files
     ├─► Delete Source (if delete_source: true)
     │
     ▼
┌─────────────────┐
│ Log Results     │
└─────────────────┘
```

---

## Storage Architecture

### Backend Abstraction

The Backend interface (internal/sync/backend.go:34) provides a unified API for all storage types:

```go
type Backend interface {
    // List all files with given prefix
    List(ctx context.Context, prefix string) ([]FileInfo, error)

    // Read a file
    Read(ctx context.Context, path string) (io.ReadCloser, error)

    // Write a file
    Write(ctx context.Context, path string, r io.Reader, size int64) error

    // Delete a file
    Delete(ctx context.Context, path string) error

    // Get file metadata
    Stat(ctx context.Context, path string) (*FileInfo, error)

    // Clean up resources
    Close() error
}
```

### FileInfo Structure

```go
type FileInfo struct {
    Path         string        // Full path
    Size         int64         // File size in bytes
    ModTime      time.Time     // Last modification time
    ETag         string        // Checksum/hash (S3 ETag or local MD5)
    IsDir        bool          // Directory flag
    StorageClass string        // S3 storage class (STANDARD, GLACIER, etc.)
}
```

### Local Backend

**Implementation:** `internal/sync/local.go`

**Features:**
- Standard Go `os` package operations
- ETag computed as MD5 hash for consistency with S3
- Supports standard filesystem operations
- Respects filesystem permissions

### S3 Backend

**Implementation:** `internal/sync/s3.go`

**Features:**
- AWS SDK for Go v2
- S3-compatible services (MinIO, DigitalOcean Spaces, etc.)
- Multipart upload support (automatic for large files)
- Storage class configuration
- Object tagging for metadata
- Server-side encryption support

**S3 URL Format:**
```
s3://bucket-name/prefix/path
s3://bucket-name/  (entire bucket)
```

**Configuration:**
- Uses AWS credentials from `~/.aws/credentials`
- Profile selection via config or environment
- Region auto-detection or explicit configuration

### Sync Strategy

The sync engine uses a **checksum-based comparison** strategy:

1. **Primary:** Compare ETag (MD5 hash for local, S3 ETag for remote)
2. **Fallback:** Compare file size and modification time
3. **Action:** Sync if source is newer or different

This ensures data integrity while minimizing unnecessary transfers.

---

## Metadata Architecture

### Metadata Schema

**Core Metadata Structure:**
```go
type Metadata struct {
    SchemaName    string                 // e.g., "microscopy", "sequencing"
    SchemaVersion string                 // Schema version
    Fields        map[string]interface{} // Extracted metadata
    FileInfo      FileInfo               // File metadata
    Provenance    Provenance             // Upload/creation info
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type FileInfo struct {
    Filename  string
    Path      string
    Size      int64
    CreatedAt time.Time
}

type Provenance struct {
    UploadedBy string
    UploadedAt time.Time
    Source     string
}
```

### Storage Mechanisms

#### 1. S3 Object Tags

**Implementation:** `internal/metadata/s3tags.go`

Key-value tags attached directly to S3 objects:

**Limitations:**
- Maximum 10 tags per object (S3 limitation)
- Keys limited to 128 characters
- Values limited to 256 characters

**Strategy:**
- Store most critical fields as tags (format, instrument, sample_id)
- Flatten nested structures with dot notation: `microscopy.manufacturer`
- Truncate long values

**Example:**
```
format: CZI
microscopy.manufacturer: Zeiss
microscopy.model: LSM 880
sample_id: sample-001
acquisition_date: 2025-01-24
```

#### 2. Sidecar JSON Files

**Location:** `{filename}.metadata.json`

Complete metadata stored alongside data files:

**Advantages:**
- No size limitations
- Preserves full structure
- Human-readable
- Portable across storage backends

**Example:**
```json
{
  "schema_name": "microscopy",
  "schema_version": "1.0",
  "fields": {
    "manufacturer": "Zeiss",
    "model": "LSM 880",
    "modality": "confocal",
    "width": 2048,
    "height": 2048,
    "channels": 3,
    "pixel_size_x": 0.125,
    "pixel_size_y": 0.125,
    "channel_info": [
      {
        "name": "DAPI",
        "index": 0,
        "excitation_wavelength": 405
      }
    ]
  },
  "file_info": {
    "filename": "experiment001.czi",
    "size": 52428800,
    "created_at": "2025-01-24T10:30:00Z"
  }
}
```

### Extractor Plugin Architecture

Extractors are self-contained modules implementing a common interface:

**Registration Flow:**
```go
registry := metadata.NewExtractorRegistry()
registry.RegisterDefaults()  // Registers all built-in extractors

// Or register custom extractor
registry.Register(&CustomExtractor{})
```

**Extraction Flow:**
1. Registry receives file path
2. FindExtractor() tests each extractor's CanHandle()
3. First matching extractor is selected
4. Extract() or ExtractFromReader() is called
5. Results wrapped in Metadata structure
6. Stored via configured storage mechanism

---

## Extensibility Points

### 1. Custom Storage Backends

Implement the Backend interface to add new storage types:

```go
type MyBackend struct {
    // Your fields
}

func (b *MyBackend) List(ctx context.Context, prefix string) ([]FileInfo, error) {
    // Implementation
}

// Implement remaining Backend methods...
```

**Use Cases:**
- Azure Blob Storage
- Google Cloud Storage
- FTP/SFTP servers
- Institutional storage systems

### 2. Custom Metadata Extractors

Implement the Extractor interface for new file formats:

```go
type CustomExtractor struct{}

func (e *CustomExtractor) CanHandle(filename string) bool {
    return strings.HasSuffix(filename, ".myformat")
}

func (e *CustomExtractor) Extract(filepath string) (map[string]interface{}, error) {
    // Parse file and extract metadata
    return map[string]interface{}{
        "format": "MyFormat",
        // ... extracted fields
    }, nil
}

// Implement remaining Extractor methods...

// Register
registry.Register(&CustomExtractor{})
```

**Use Cases:**
- Lab-specific file formats
- Custom instrument outputs
- Proprietary data formats

### 3. Custom Presets

Define validation rules for your lab's specific requirements:

```go
customPreset := metadata.Preset{
    Name:        "lab-custom-protocol",
    Description: "Custom protocol requirements for Lab X",
    RequiredFields: []string{
        "sample_id",
        "operator",
        "project_code",
    },
    RecommendedFields: []string{
        "notes",
        "collaborators",
    },
    AllowedValues: map[string][]string{
        "operator": {"alice", "bob", "charlie"},
        "project_code": {"PROJ001", "PROJ002"},
    },
}
```

**Use Cases:**
- Lab-specific data requirements
- Compliance validation
- Quality control standards

### 4. Custom DOI Providers

Implement the Provider interface for additional DOI services:

```go
type CustomProvider struct {
    // Provider configuration
}

func (p *CustomProvider) Validate(metadata DataciteMetadata) error {
    // Validation logic
}

func (p *CustomProvider) Create(metadata DataciteMetadata) (*DOIRecord, error) {
    // DOI creation logic
}

// Implement remaining Provider methods...
```

**Use Cases:**
- Institutional repositories
- Domain-specific archives
- Custom DOI prefixes

---

## Design Principles

### 1. Dormant by Default

Cicada is designed to consume minimal resources when idle:

- No background daemons (except explicit watch mode)
- On-demand execution model
- Efficient resource cleanup
- Optional long-running processes

### 2. Cost-Effective

Optimized for small lab budgets:

- Efficient S3 operations (minimize API calls)
- Parallel transfers to reduce time
- Smart sync (only transfer changed files)
- Storage class optimization support
- Minimal computational requirements

### 3. Easy to Use

Designed for users with limited technical expertise:

- Intuitive CLI commands
- Sensible defaults
- Helpful error messages
- Comprehensive documentation
- Progressive disclosure (simple → advanced features)

### 4. Extensible

Plugin-friendly architecture:

- Interface-based design
- Clear extension points
- Minimal coupling
- Composition over inheritance

### 5. Data Integrity

Ensuring data reliability:

- Checksum-based verification
- Atomic operations where possible
- Error handling and recovery
- Dry-run mode for safety

### 6. Format Agnostic

Supporting diverse scientific data:

- Multiple extractor framework
- Generic fallback handling
- Instrument-specific metadata types
- Flexible schema system

### 7. Platform Optional

Features are additive, not mandatory:

- Core: Storage + Sync (always available)
- Optional: Metadata extraction
- Optional: Watch mode
- Optional: DOI preparation
- Optional: Advanced validation

Users start simple and add capabilities as needed.

---

## Performance Characteristics

### Sync Performance

- **Concurrency:** Default 4 parallel transfers (configurable)
- **Throughput:** Limited by network bandwidth and backend performance
- **Memory:** Streaming I/O (no full file buffering)
- **CPU:** Minimal (mostly I/O bound)

### Metadata Extraction

- **Speed:** Varies by format and file size
- **Memory:** Format-dependent (some require file buffering)
- **CPU:** Parsing overhead for complex formats

### Watch Mode

- **Responsiveness:** Debounce delay + transfer time
- **Resource Usage:** One goroutine per watch + fsnotify overhead
- **Scalability:** Hundreds of watches per instance

---

## Security Considerations

### Credentials Management

- AWS credentials stored in `~/.aws/credentials` (never in config)
- Supports IAM roles and instance profiles
- Environment variable override support

### Data Privacy

- Local encryption at rest depends on filesystem
- S3 server-side encryption support
- No data transmitted to Cicada servers (no telemetry)

### DICOM Handling

- Built-in PHI/PII redaction warnings
- Extractor designed to exclude sensitive fields
- Validation against HIPAA requirements (planned)

---

## Future Architecture

### Planned Enhancements

1. **Azure & GCS Backends** - Additional cloud storage support
2. **Metadata Search** - Query interface for metadata
3. **Web UI** - Optional web interface for monitoring
4. **API Server** - REST API for programmatic access
5. **Plugin System** - Dynamic loading of extractors/backends
6. **Distributed Sync** - Multi-node coordination
7. **Data Versioning** - Track file changes over time

---

## Related Documentation

- [User Guide](USER_GUIDE.md) - Getting started and common workflows
- [CLI Reference](CLI_REFERENCE.md) - Complete command documentation
- [Configuration Guide](CONFIGURATION.md) - Configuration options
- [Development Guide](DEVELOPMENT.md) - Contributing to Cicada
- [Metadata System](METADATA_SYSTEM.md) - Detailed metadata documentation

---

**Contributing:** Found an error or want to suggest improvements? See [CONTRIBUTING.md](CONTRIBUTING.md)

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
