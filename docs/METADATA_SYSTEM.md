# Metadata System Guide

**Last Updated:** 2025-01-24

Complete guide to Cicada's metadata extraction and management system for scientific data.

## Table of Contents

1. [System Overview](#system-overview)
2. [Metadata Architecture](#metadata-architecture)
3. [Extractor System](#extractor-system)
4. [Supported File Formats](#supported-file-formats)
5. [Instrument-Specific Metadata](#instrument-specific-metadata)
6. [Metadata Schema](#metadata-schema)
7. [Storage Mechanisms](#storage-mechanisms)
8. [Custom Extractors](#custom-extractors)
9. [Metadata Validation](#metadata-validation)
10. [Best Practices](#best-practices)
11. [Future Features](#future-features)
12. [Troubleshooting](#troubleshooting)

---

## System Overview

Cicada's metadata system automatically extracts, validates, and stores metadata from scientific instrument files.

### Key Features

- **14 Format Extractors**: Microscopy, sequencing, mass spec, and more
- **6 Instrument Types**: Domain-specific metadata structures
- **Automatic Detection**: File format auto-detection
- **Multiple Storage**: S3 tags and sidecar JSON files
- **Validation**: Preset-based quality checks
- **Extensible**: Plugin architecture for custom extractors

### Design Principles

1. **Format Agnostic**: Support diverse scientific data types
2. **Automatic**: No manual metadata entry required
3. **Structured**: Domain-specific schemas for rich metadata
4. **Accessible**: Store metadata for easy discovery
5. **Validated**: Ensure data quality through presets

---

## Metadata Architecture

### Component Overview

```
┌─────────────────────────────────────────┐
│          CLI Interface                  │
│   cicada metadata extract/validate      │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│      Extractor Registry                 │
│   - FindExtractor(filename)             │
│   - Extract(filepath)                   │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│      Format Extractors (14)             │
│   TIFF│OME-TIFF│CZI│ND2│LIF│            │
│   FASTQ│BAM│mzML│MGF│HDF5│              │
│   Zarr│DICOM│FCS│Generic                │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│      Metadata Structure                 │
│   - Fields (map[string]interface{})    │
│   - FileInfo                            │
│   - Provenance                          │
│   - Timestamps                          │
└────────────┬────────────────────────────┘
             │
             ├──► S3 Object Tags
             └──► Sidecar JSON Files
```

### Metadata Flow

```
1. User Command
   cicada metadata extract file.czi
      │
      ▼
2. File Type Detection
   Extension: .czi → ZeissCZIExtractor
      │
      ▼
3. Metadata Extraction
   Parse CZI binary format
   Extract embedded metadata
      │
      ▼
4. Structure Creation
   Build Metadata struct with:
   - Extracted fields
   - File information
   - Provenance data
      │
      ▼
5. Validation (Optional)
   Check against preset rules
      │
      ▼
6. Storage
   ├─ S3 Tags (limited fields)
   └─ Sidecar JSON (complete)
      │
      ▼
7. Output
   Display to user (JSON/YAML/Table)
```

---

## Extractor System

### Extractor Interface

All extractors implement a common interface:

```go
type Extractor interface {
    // Check if extractor can handle file
    CanHandle(filename string) bool

    // Extract metadata from file
    Extract(filepath string) (map[string]interface{}, error)

    // Extract from reader (streaming)
    ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error)

    // Get extractor name
    Name() string

    // Get supported file extensions
    SupportedFormats() []string
}
```

### Extractor Registry

**Registration:**
```go
registry := metadata.NewExtractorRegistry()
registry.RegisterDefaults()  // Register all built-in extractors

// Or register specific extractor
registry.Register(&CustomExtractor{})
```

**Lookup:**
```go
extractor := registry.FindExtractor("image.czi")
if extractor != nil {
    metadata, err := extractor.Extract("image.czi")
}
```

**Auto-Detection Algorithm:**
1. Iterate through registered extractors
2. Call `CanHandle()` for each
3. First match is selected
4. Generic extractor is fallback

### Built-in Extractors

| Extractor | Extensions | Status | Domain |
|-----------|-----------|--------|---------|
| TIFF | `.tif`, `.tiff` | Basic | Microscopy |
| OME-TIFF | `.ome.tif`, `.ome.tiff` | Full | Microscopy |
| Zeiss CZI | `.czi` | Full | Microscopy |
| Nikon ND2 | `.nd2` | Placeholder | Microscopy |
| Leica LIF | `.lif` | Placeholder | Microscopy |
| FASTQ | `.fastq`, `.fq`, `.fastq.gz` | Full | Sequencing |
| BAM | `.bam` | Placeholder | Sequencing |
| mzML | `.mzml` | Placeholder | Mass Spec |
| MGF | `.mgf` | Placeholder | Mass Spec |
| HDF5 | `.h5`, `.hdf5` | Placeholder | Data Arrays |
| Zarr | `.zarr` | Placeholder | Data Arrays |
| DICOM | `.dcm`, `.dicom` | Placeholder | Medical |
| FCS | `.fcs` | Placeholder | Flow Cytometry |
| Generic | `*` | Basic | Fallback |

**Status Definitions:**
- **Full**: Complete implementation with rich metadata extraction
- **Basic**: Basic metadata extraction (format, size, timestamps)
- **Placeholder**: Framework only, needs implementation

---

## Supported File Formats

### Microscopy Formats

#### OME-TIFF

**Description:** Open Microscopy Environment TIFF with embedded XML metadata

**Extracted Fields:**
- Image dimensions (width, height, depth, channels, timepoints)
- Pixel sizes (X, Y, Z in micrometers)
- Channel information (name, wavelength, color)
- Acquisition parameters
- Instrument details
- Sample information

**Example:**
```json
{
  "format": "OME-TIFF",
  "width": 2048,
  "height": 2048,
  "depth": 50,
  "channels": 3,
  "timepoints": 1,
  "pixel_size_x": 0.125,
  "pixel_size_y": 0.125,
  "pixel_size_z": 0.5,
  "voxel_size_unit": "micrometers",
  "channel_info": [
    {
      "name": "DAPI",
      "index": 0,
      "excitation_wavelength": 405,
      "emission_wavelength": 450,
      "color": "cyan"
    },
    {
      "name": "GFP",
      "index": 1,
      "excitation_wavelength": 488,
      "emission_wavelength": 520,
      "color": "green"
    }
  ]
}
```

**CLI Usage:**
```bash
cicada metadata extract image.ome.tif
```

---

#### Zeiss CZI

**Description:** Zeiss confocal microscopy format with rich metadata

**Extracted Fields:**
- Microscope manufacturer, model, serial number
- Modality (confocal, widefield, etc.)
- Image dimensions
- Pixel/voxel sizes
- Channel details (fluorophores, wavelengths, laser power)
- Objective information (magnification, NA)
- Detector settings (gain, offset)
- Acquisition timestamps
- Operator information
- Environmental conditions (temperature, CO2)

**Example:**
```json
{
  "format": "CZI",
  "microscope_manufacturer": "Zeiss",
  "microscope_model": "LSM 880",
  "serial_number": "123456789",
  "software_version": "ZEN 2.6",
  "modality": "confocal",
  "width": 2048,
  "height": 2048,
  "depth": 64,
  "channels": 3,
  "timepoints": 10,
  "pixel_size_x": 0.125,
  "pixel_size_y": 0.125,
  "pixel_size_z": 0.5,
  "objective": "Plan-Apochromat 40x/1.4 Oil",
  "magnification": 40.0,
  "numerical_aperture": 1.4,
  "channel_info": [
    {
      "name": "DAPI",
      "fluorophore": "DAPI",
      "excitation_wavelength": 405,
      "emission_wavelength": 450,
      "laser_power": 2.5
    }
  ],
  "acquisition_date": "2025-01-24T10:30:00Z",
  "operator": "jsmith",
  "experiment_name": "Neural Imaging Study",
  "sample_id": "SAMPLE-001"
}
```

**CLI Usage:**
```bash
cicada metadata extract confocal_image.czi
cicada metadata validate confocal_image.czi --preset microscopy-confocal
```

---

### Sequencing Formats

#### FASTQ

**Description:** Text-based sequence format with quality scores

**Extracted Fields:**
- File format version
- Total reads
- Read length
- Quality score encoding (Phred+33, Phred+64)
- Average quality score
- Sample reads (first 10 for inspection)
- Instrument ID (from read headers)
- Run ID, lane, tile information
- GC content statistics
- Quality score distribution

**Example:**
```json
{
  "format": "FASTQ",
  "total_reads": 10000000,
  "read_length": 150,
  "quality_encoding": "Phred+33",
  "average_quality": 36.5,
  "gc_content": 45.2,
  "instrument_id": "A00123",
  "run_id": "201215_A00123_0456_BHXXXX",
  "flowcell_id": "BHXXXX",
  "lane": 1,
  "quality_distribution": {
    "q30_percent": 92.5,
    "q20_percent": 98.2
  }
}
```

**CLI Usage:**
```bash
cicada metadata extract sample_R1.fastq.gz
cicada metadata validate sample_R1.fastq.gz --preset sequencing-illumina
```

---

### Mass Spectrometry Formats

#### mzML

**Description:** XML-based mass spectrometry data format

**Extracted Fields:**
- Instrument manufacturer, model
- Instrument type (Orbitrap, Q-TOF, etc.)
- Ionization mode (ESI, MALDI, APCI)
- Polarity (positive/negative)
- Mass analyzer type
- Total spectra count
- MS1/MS2 spectrum counts
- Scan range (m/z)
- Acquisition date
- Sample information

**Example (Placeholder):**
```json
{
  "format": "mzML",
  "manufacturer": "Thermo",
  "model": "Q Exactive HF",
  "instrument_type": "Orbitrap",
  "ionization_mode": "ESI",
  "polarity": "positive",
  "total_spectra": 25000,
  "ms1_spectra": 5000,
  "ms2_spectra": 20000,
  "scan_range": "300-2000 m/z",
  "resolution": 70000,
  "acquisition_date": "2025-01-24T10:00:00Z"
}
```

**CLI Usage:**
```bash
cicada metadata extract proteomics_sample.mzml
```

---

### Flow Cytometry Formats

#### FCS

**Description:** Flow Cytometry Standard format with embedded parameters

**Extracted Fields:**
- Instrument manufacturer, model
- Software version
- Total events
- Event rate
- Acquisition time
- Parameters/channels (FSC, SSC, fluorescence channels)
- Detector voltages and gains
- Compensation matrix
- Sample information

**Example (Placeholder):**
```json
{
  "format": "FCS",
  "manufacturer": "BD",
  "model": "FACSAria III",
  "software_version": "FACSDiva 8.0",
  "total_events": 500000,
  "event_rate": 2500,
  "acquisition_time": 200.0,
  "parameters": [
    {
      "name": "FSC-A",
      "description": "Forward Scatter Area",
      "range": 262144,
      "bits": 18,
      "voltage": 450
    },
    {
      "name": "PE-A",
      "description": "PE Area",
      "fluorochrome": "PE",
      "filter": "575/26",
      "voltage": 500,
      "gain": 1.5
    }
  ]
}
```

---

### Data Array Formats

#### HDF5

**Description:** Hierarchical Data Format for arrays and metadata

**Extracted Fields:**
- HDF5 version
- Root group attributes
- Dataset names and shapes
- Data types
- Compression settings
- Custom attributes

**Example (Placeholder):**
```json
{
  "format": "HDF5",
  "hdf5_version": "1.10.7",
  "datasets": [
    {
      "name": "/data/images",
      "shape": [100, 2048, 2048],
      "dtype": "uint16",
      "compression": "gzip"
    }
  ],
  "attributes": {
    "experiment": "time_series",
    "sample_id": "SAMPLE-001"
  }
}
```

---

#### Zarr

**Description:** Chunked array storage format

**Extracted Fields:**
- Zarr version
- Array shapes
- Chunk sizes
- Compression codecs
- Attributes (.zattrs)

**Example (Placeholder):**
```json
{
  "format": "Zarr",
  "zarr_version": "2",
  "arrays": [
    {
      "name": "data",
      "shape": [1000, 2048, 2048],
      "chunks": [1, 256, 256],
      "dtype": "uint16",
      "compressor": "blosc"
    }
  ]
}
```

---

## Instrument-Specific Metadata

### MicroscopyMetadata

**Fields:**

**Instrument:**
- Manufacturer (e.g., "Zeiss", "Nikon", "Leica")
- Model (e.g., "LSM 880", "Ti2-E")
- Serial number
- Software version

**Imaging Parameters:**
- Modality (confocal, widefield, TIRF, light sheet)
- Magnification (e.g., 40.0)
- Numerical aperture (e.g., 1.4)
- Objective description

**Image Dimensions:**
- Width, height (pixels)
- Depth (Z slices)
- Channels
- Timepoints
- Pixel sizes (X, Y, Z in µm)
- Voxel size unit

**Channel Information:**
- Name (e.g., "DAPI", "GFP")
- Index
- Fluorophore
- Excitation/emission wavelengths (nm)
- Laser power (%)
- Display color
- Contrast settings

**Acquisition Settings:**
- Exposure time (ms)
- Frame rate (fps)
- Binning (X, Y)
- Detector gain
- Detector model

**Experiment:**
- Experiment name
- Sample ID
- Organism
- Tissue
- Cell line
- Treatment
- Acquisition date
- Operator

**Environment:**
- Temperature (°C)
- CO2 level (%)
- Humidity (%)

**Go Struct:**
```go
type MicroscopyMetadata struct {
    Manufacturer     string                 `json:"manufacturer"`
    Model           string                 `json:"model"`
    SerialNumber    string                 `json:"serial_number,omitempty"`
    SoftwareVersion string                 `json:"software_version,omitempty"`
    Modality        string                 `json:"modality"`
    Magnification   float64                `json:"magnification,omitempty"`
    NumericalAperture float64              `json:"numerical_aperture,omitempty"`
    Width           int                    `json:"width"`
    Height          int                    `json:"height"`
    Depth           int                    `json:"depth,omitempty"`
    Channels        int                    `json:"channels"`
    PixelSizeX      float64                `json:"pixel_size_x,omitempty"`
    ChannelInfo     []MicroscopyChannel    `json:"channel_info,omitempty"`
    // ... additional fields
}
```

---

### SequencingMetadata

**Fields:**

**Instrument:**
- Platform (Illumina, PacBio, ONT)
- Model
- Serial number
- Software version

**Run Information:**
- Run ID
- Flowcell ID
- Lane
- Run date
- Operator

**Library:**
- Library ID
- Library kit
- Protocol
- Insert size (bp)
- Index sequence

**Read Configuration:**
- Read type (paired-end, single-end)
- Read length
- Read1/Read2 lengths
- Index length

**Quality Metrics:**
- Total reads
- Pass filter reads
- Average quality score (Phred)
- Percent Q30
- Duplication rate

**Sample:**
- Sample ID
- Sample name
- Organism
- Tissue
- Cell type
- Treatment
- Reference genome

**Assay:**
- Assay type (RNA-Seq, ChIP-Seq, WGS)
- Target region
- Enrichment kit

---

### MassSpecMetadata

**Fields:**

**Instrument:**
- Manufacturer
- Model
- Serial number
- Software version

**Mass Spectrometer:**
- Instrument type (Orbitrap, Q-TOF, Triple Quad)
- Ionization mode (ESI, MALDI, APCI)
- Polarity (positive/negative)
- Mass analyzer

**Acquisition:**
- Scan range (m/z)
- Resolution
- Scan rate (scans/second)
- Total spectra
- MS1/MS2 spectra counts

**Chromatography (LC-MS):**
- Chromatography type (HPLC, UPLC)
- Column type (C18, etc.)
- Column length (mm)
- Flow rate (µL/min)
- Run time (minutes)

**Sample:**
- Sample ID
- Sample type
- Organism
- Preparation method
- Acquisition date
- Operator

**Experiment:**
- Experiment type (proteomics, metabolomics, lipidomics)
- Acquisition mode (DDA, DIA, targeted)

---

### FlowCytometryMetadata

**Fields:**

**Instrument:**
- Manufacturer
- Model
- Serial number
- Software version

**Acquisition:**
- Total events
- Event rate (events/second)
- Aborted events
- Acquisition time (seconds)
- Acquisition date

**Parameters/Channels:**
- Name (FSC-A, PE-A, etc.)
- Description
- Range
- Bits (bit depth)
- Gain
- Voltage
- Filter (optical filter)
- Fluorochrome

**Analysis:**
- Compensation matrix
- Populations (gated)
- Gating strategy

**Sample:**
- Sample ID
- Tube ID
- Organism
- Cell type
- Treatment
- Operator

---

### CryoEMMetadata

**Fields:**

**Instrument:**
- Manufacturer
- Model
- Voltage (kV)
- Serial number

**Detector:**
- Detector model
- Detector mode
- Pixel size (Å/pixel)

**Acquisition:**
- Magnification
- Defocus (µm)
- Exposure time (seconds)
- Total dose (e⁻/Å²)
- Frames per movie
- Total movies

**Sample:**
- Sample ID
- Protein/complex name
- Organism
- Grid type
- Freezing method
- Acquisition date
- Operator

---

### XRayMetadata

**Fields:**

**Facility:**
- Facility name
- Beamline
- Detector model

**Acquisition:**
- Wavelength (Å)
- Energy (keV)
- Detector distance (mm)
- Exposure time (seconds/frame)
- Oscillation (degrees/frame)
- Total frames

**Crystal:**
- Space group
- Unit cell (a, b, c in Å)
- Unit cell angles (α, β, γ in degrees)

**Sample:**
- Sample ID
- Protein name
- Organism
- Acquisition date
- Operator

**Processing:**
- Resolution (Å)
- Completeness (%)
- Rmerge
- I/σ(I)

---

## Metadata Schema

### Core Metadata Structure

```go
type Metadata struct {
    SchemaName    string                 `json:"schema_name"`
    SchemaVersion string                 `json:"schema_version"`
    Fields        map[string]interface{} `json:"fields"`
    FileInfo      FileInfo               `json:"file_info"`
    Provenance    Provenance             `json:"provenance"`
    CreatedAt     time.Time              `json:"created_at"`
    UpdatedAt     time.Time              `json:"updated_at"`
}

type FileInfo struct {
    Filename  string    `json:"filename"`
    Path      string    `json:"path"`
    Size      int64     `json:"size"`
    CreatedAt time.Time `json:"created_at"`
}

type Provenance struct {
    UploadedBy string    `json:"uploaded_by"`
    UploadedAt time.Time `json:"uploaded_at"`
    Source     string    `json:"source"`
}
```

### Schema Naming

- `microscopy` - Microscopy data
- `sequencing` - DNA/RNA sequencing
- `mass_spec` - Mass spectrometry
- `flow_cytometry` - Flow cytometry
- `cryo_em` - Cryo-electron microscopy
- `xray` - X-ray crystallography
- `generic` - Generic/unknown format

### Schema Versioning

Format: `major.minor`

- **Major**: Breaking changes (incompatible field changes)
- **Minor**: Additive changes (new optional fields)

Current version: `1.0`

---

## Storage Mechanisms

### S3 Object Tags

**Advantages:**
- Searchable via S3 API
- No additional files
- Cost-effective

**Limitations:**
- Maximum 10 tags per object
- Key: 128 characters max
- Value: 256 characters max

**Tag Selection Strategy:**
1. Most critical fields only
2. Flatten nested structures
3. Truncate long values
4. Use abbreviations

**Example Tag Set:**
```
format=CZI
microscope=Zeiss-LSM-880
modality=confocal
date=2025-01-24
sample_id=EXP001
operator=jsmith
channels=3
width=2048
height=2048
magnification=40
```

**Implementation:**
```go
// Select most important fields
tags := selectTagFields(metadata, maxTags=10)

// Apply to S3 object
s3Client.PutObjectTagging(&s3.PutObjectTaggingInput{
    Bucket: bucket,
    Key:    key,
    Tagging: &types.Tagging{
        TagSet: tags,
    },
})
```

---

### Sidecar JSON Files

**Advantages:**
- No size limitations
- Complete metadata preservation
- Human-readable
- Portable across backends

**Naming Convention:**
```
<filename>.metadata.json
```

**Examples:**
```
experiment001.czi           → experiment001.czi.metadata.json
sample_R1.fastq.gz          → sample_R1.fastq.gz.metadata.json
proteomics_run.mzml         → proteomics_run.mzml.metadata.json
```

**File Format:**
```json
{
  "schema_name": "microscopy",
  "schema_version": "1.0",
  "fields": {
    "manufacturer": "Zeiss",
    "model": "LSM 880",
    "width": 2048,
    "height": 2048,
    "channels": 3,
    "channel_info": [
      {
        "name": "DAPI",
        "excitation_wavelength": 405
      }
    ]
  },
  "file_info": {
    "filename": "experiment001.czi",
    "path": "/data/microscopy/experiment001.czi",
    "size": 524288000,
    "created_at": "2025-01-24T10:00:00Z"
  },
  "provenance": {
    "uploaded_by": "jsmith",
    "uploaded_at": "2025-01-24T11:00:00Z",
    "source": "zeiss-lsm-880"
  },
  "created_at": "2025-01-24T11:00:05Z",
  "updated_at": "2025-01-24T11:00:05Z"
}
```

---

### Dual Storage Strategy

**Recommended Approach:**
1. Extract metadata once
2. Store summary in S3 tags (10 most important fields)
3. Store complete metadata in sidecar JSON
4. Keep both in sync

**Benefits:**
- Fast discovery via S3 tags
- Complete metadata in sidecar files
- Redundancy for reliability

**Implementation:**
```bash
# Extract metadata
cicada metadata extract file.czi --output file.czi.metadata.json

# Upload to S3 with tags
cicada sync file.czi s3://bucket/data/
cicada sync file.czi.metadata.json s3://bucket/data/

# Tags are automatically applied during sync
```

---

## Custom Extractors

### Creating Custom Extractor

**Step 1: Implement Interface**

```go
package myextractors

import (
    "io"
    "path/filepath"
    "strings"
)

type CustomExtractor struct{}

func (e *CustomExtractor) Name() string {
    return "My Custom Format"
}

func (e *CustomExtractor) SupportedFormats() []string {
    return []string{".custom", ".cst"}
}

func (e *CustomExtractor) CanHandle(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    for _, format := range e.SupportedFormats() {
        if ext == format {
            return true
        }
    }
    return false
}

func (e *CustomExtractor) Extract(filepath string) (map[string]interface{}, error) {
    // Read file
    data, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    // Parse format
    // ... custom parsing logic ...

    // Return metadata
    return map[string]interface{}{
        "format":     "CUSTOM",
        "version":    "1.0",
        "instrument": "Custom Instrument",
        // ... additional fields ...
    }, nil
}

func (e *CustomExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
    // Streaming extraction
    return nil, fmt.Errorf("not implemented")
}
```

**Step 2: Register Extractor**

```go
// In main or initialization code
registry := metadata.NewExtractorRegistry()
registry.RegisterDefaults()
registry.Register(&CustomExtractor{})
```

**Step 3: Use Extractor**

```bash
# Now works automatically
cicada metadata extract data.custom
```

---

### Extractor Best Practices

**1. Error Handling**
```go
func (e *Extractor) Extract(filepath string) (map[string]interface{}, error) {
    // Validate file exists
    if _, err := os.Stat(filepath); err != nil {
        return nil, fmt.Errorf("file not found: %w", err)
    }

    // Parse with error handling
    data, err := parseFile(filepath)
    if err != nil {
        return nil, fmt.Errorf("parse failed: %w", err)
    }

    return data, nil
}
```

**2. Graceful Degradation**
```go
// Don't fail on missing optional fields
metadata := map[string]interface{}{
    "format": "CUSTOM", // Required
}

if instrument, ok := data["instrument"]; ok {
    metadata["instrument"] = instrument // Optional
}
```

**3. Type Safety**
```go
// Use type assertions carefully
if width, ok := data["width"].(int); ok {
    metadata["width"] = width
} else if widthStr, ok := data["width"].(string); ok {
    // Try parsing as string
    if w, err := strconv.Atoi(widthStr); err == nil {
        metadata["width"] = w
    }
}
```

**4. Performance**
```go
// For large files, avoid reading entire file
func (e *Extractor) Extract(filepath string) (map[string]interface{}, error) {
    f, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    // Read only header/metadata section
    header := make([]byte, 1024)
    n, err := f.Read(header)
    if err != nil && err != io.EOF {
        return nil, err
    }

    return parseHeader(header[:n]), nil
}
```

---

## Metadata Validation

### Preset System

**Purpose:** Ensure data quality through validation rules

**Preset Structure:**
```go
type Preset struct {
    Name              string
    Description       string
    RequiredFields    []string
    RecommendedFields []string
    AllowedValues     map[string][]string
    ValidateFunc      func(Metadata) []ValidationError
}
```

### Built-in Presets

**microscopy-confocal:**
```yaml
required_fields:
  - manufacturer
  - model
  - modality
  - width
  - height
  - channels
  - pixel_size_x
  - pixel_size_y
  - magnification
  - numerical_aperture
  - acquisition_date
  - operator

recommended_fields:
  - sample_id
  - experiment_name
  - organism
  - tissue
  - temperature
  - co2_level
  - objective
  - detector_model

allowed_values:
  modality: [confocal, spinning-disk, two-photon]
```

**sequencing-illumina:**
```yaml
required_fields:
  - platform
  - model
  - run_id
  - flowcell_id
  - lane
  - read_length
  - total_reads
  - sample_id

recommended_fields:
  - library_id
  - library_kit
  - organism
  - tissue
  - percent_q30

allowed_values:
  platform: [Illumina]
  model: [NovaSeq 6000, NextSeq 2000, MiSeq]
  read_type: [paired-end, single-end]
```

### Validation Process

**CLI Usage:**
```bash
# Validate with preset
cicada metadata validate file.czi --preset microscopy-confocal

# Strict mode (warnings are errors)
cicada metadata validate file.czi --preset microscopy-confocal --strict
```

**Validation Output:**
```
Validating: file.czi
Preset: microscopy-confocal

✓ Required fields present (12/12)
⚠ Recommended fields missing (2/8):
  - temperature
  - co2_level
✓ Field values valid

Validation: PASSED (2 warnings)
```

**Validation Errors:**
```
Validating: file.czi
Preset: microscopy-confocal

✗ Required fields missing (2/12):
  - operator
  - acquisition_date
⚠ Recommended fields missing (3/8):
  - sample_id
  - experiment_name
  - organism
✗ Invalid field values:
  - modality: "widefield" not in [confocal, spinning-disk, two-photon]

Validation: FAILED
```

---

## Best Practices

### 1. Extract Metadata Early

```bash
# Extract immediately after data acquisition
cicada metadata extract /data/today/*.czi

# Automated extraction in watch mode
cicada watch add /data/microscope s3://lab-data/microscopy --extract-metadata
```

### 2. Validate Before Upload

```bash
# Validate locally before syncing to cloud
for file in /data/*.czi; do
  if cicada metadata validate "$file" --preset microscopy-confocal; then
    cicada sync "$file" s3://lab-data/validated/
  else
    echo "Validation failed: $file"
  fi
done
```

### 3. Use Descriptive Sample IDs

**Good:**
```
sample_id: "EXP001-CTRL-01"
sample_id: "2025-01-24-BRAIN-SAMPLE-A"
sample_id: "PROJECT-X-REPLICATE-3"
```

**Bad:**
```
sample_id: "1"
sample_id: "test"
sample_id: "sample"
```

### 4. Store Complete Metadata

```bash
# Always save complete metadata, even if using S3 tags
cicada metadata extract file.czi --output file.czi.metadata.json

# Upload both data and metadata
cicada sync file.czi s3://bucket/data/
cicada sync file.czi.metadata.json s3://bucket/data/
```

### 5. Document Custom Fields

```json
{
  "fields": {
    "custom_field_1": "value",
    "_custom_field_1_description": "Describes what this custom field means",
    "custom_field_2": 123,
    "_custom_field_2_unit": "milliseconds"
  }
}
```

---

## Future Features

### Planned Enhancements

**1. Metadata Search (v0.4.0)**
```bash
# Search by metadata fields
cicada metadata search --field sample_id=EXP001
cicada metadata search --field operator=jsmith --date-range 2025-01-01,2025-01-31
```

**2. Metadata Query Language (v0.4.0)**
```bash
# SQL-like queries
cicada metadata query "SELECT * FROM metadata WHERE channels > 3 AND modality = 'confocal'"
```

**3. Batch Extraction (v0.3.1)**
```bash
# Extract from multiple files in parallel
cicada metadata extract --batch /data/microscopy/*.czi --workers 8
```

**4. Metadata Export (v0.3.1)**
```bash
# Export to CSV for analysis
cicada metadata export /data --format csv --output metadata.csv

# Export to database
cicada metadata export /data --database postgres://localhost/metadata
```

**5. Custom Validation Rules (v0.4.0)**
```yaml
# custom-preset.yaml
name: lab-specific-confocal
base: microscopy-confocal
additional_required:
  - project_code
  - funding_source
custom_rules:
  - field: magnification
    min: 20
    max: 100
  - field: pixel_size_x
    max: 0.5
```

---

## Troubleshooting

### Extraction Failures

**Symptom:**
```
Error: no extractor found for file: unknown.xyz
```

**Solutions:**
1. Check file extension is supported
2. Verify file is not corrupted
3. Use generic extractor explicitly
4. Implement custom extractor

```bash
# List supported formats
cicada metadata list

# Force generic extractor
cicada metadata extract unknown.xyz --extractor generic
```

---

### Missing Metadata Fields

**Symptom:**
```
Warning: Required field 'operator' not found
```

**Solutions:**
1. File doesn't contain that metadata
2. Extractor needs enhancement
3. Add field manually to sidecar JSON

```bash
# Extract what's available
cicada metadata extract file.czi --output metadata.json

# Edit metadata.json to add missing fields
# Then upload corrected version
cicada sync metadata.json s3://bucket/data/
```

---

### Validation Failures

**Symptom:**
```
Validation: FAILED
Required fields missing: operator, acquisition_date
```

**Solutions:**
1. Ensure instrument software records metadata
2. Add missing fields manually
3. Create custom preset with relaxed rules

```bash
# Create custom preset without strict requirements
cat > custom-preset.yaml <<EOF
name: relaxed-microscopy
base: microscopy-confocal
required_fields:
  - manufacturer
  - model
  - width
  - height
EOF

# Validate with custom preset
cicada metadata validate file.czi --preset custom-preset.yaml
```

---

## Related Documentation

- [Architecture](ARCHITECTURE.md) - Metadata architecture details
- [CLI Reference](CLI_REFERENCE.md) - Metadata commands
- [File Formats](METADATA_EXTRACTION.md) - Format-specific extraction guides

---

**Contributing:** Want to add a new extractor? See [DEVELOPMENT.md](DEVELOPMENT.md)

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
