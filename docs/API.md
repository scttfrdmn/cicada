# API Reference

**Last Updated:** 2025-11-25

Complete API reference for Cicada's Go packages.

## Table of Contents

1. [Overview](#overview)
2. [sync Package](#sync-package)
3. [metadata Package](#metadata-package)
4. [watch Package](#watch-package)
5. [config Package](#config-package)
6. [doi Package](#doi-package)
7. [Usage Examples](#usage-examples)
8. [Error Handling](#error-handling)

---

## Overview

Cicada provides Go packages for programmatic access to its functionality.

### Import Paths

```go
import (
    "github.com/scttfrdmn/cicada/internal/sync"
    "github.com/scttfrdmn/cicada/internal/metadata"
    "github.com/scttfrdmn/cicada/internal/watch"
    "github.com/scttfrdmn/cicada/internal/config"
    "github.com/scttfrdmn/cicada/internal/doi"
)
```

### Package Dependencies

```
config (base configuration)
   │
   ├─► sync (storage and sync)
   │    │
   │    └─► watch (file watching)
   │
   └─► metadata (metadata extraction)
        │
        └─► doi (DOI preparation)
```

---

## sync Package

Package `sync` provides storage backends and synchronization engine.

### Backend Interface

```go
type Backend interface {
    // List returns all files with the given prefix
    List(ctx context.Context, prefix string) ([]FileInfo, error)

    // Read opens a file for reading
    Read(ctx context.Context, path string) (io.ReadCloser, error)

    // Write writes a file
    Write(ctx context.Context, path string, r io.Reader, size int64) error

    // Delete deletes a file
    Delete(ctx context.Context, path string) error

    // Stat gets file metadata
    Stat(ctx context.Context, path string) (*FileInfo, error)

    // Close closes the backend and releases resources
    Close() error
}
```

### FileInfo

```go
type FileInfo struct {
    Path         string    // Full path
    Size         int64     // File size in bytes
    ModTime      time.Time // Last modification time
    ETag         string    // Checksum/hash
    IsDir        bool      // Directory flag
    StorageClass string    // S3 storage class (if applicable)
}
```

### LocalBackend

**Constructor:**
```go
func NewLocalBackend(basePath string) (*LocalBackend, error)
```

**Example:**
```go
backend, err := sync.NewLocalBackend("/data/lab")
if err != nil {
    log.Fatal(err)
}
defer backend.Close()

// List files
files, err := backend.List(context.Background(), "")
if err != nil {
    log.Fatal(err)
}

for _, file := range files {
    fmt.Printf("%s (%d bytes)\n", file.Path, file.Size)
}
```

### S3Backend

**Constructor:**
```go
func NewS3Backend(ctx context.Context, bucket string) (*S3Backend, error)
```

**Example:**
```go
ctx := context.Background()
backend, err := sync.NewS3Backend(ctx, "lab-data")
if err != nil {
    log.Fatal(err)
}
defer backend.Close()

// List files in prefix
files, err := backend.List(ctx, "microscopy/")
if err != nil {
    log.Fatal(err)
}
```

**Utility Functions:**
```go
// ParseS3URI parses S3 URI into bucket and key
func ParseS3URI(uri string) (bucket, key string, err error)
```

**Example:**
```go
bucket, key, err := sync.ParseS3URI("s3://lab-data/microscopy/file.czi")
// bucket = "lab-data"
// key = "microscopy/file.czi"
```

### SyncEngine

**Constructor:**
```go
func NewEngine(source, destination Backend, options SyncOptions) *Engine
```

**SyncOptions:**
```go
type SyncOptions struct {
    // DryRun shows what would be synced without making changes
    DryRun bool

    // Delete removes files in destination not present in source
    Delete bool

    // Concurrency controls parallel transfers
    Concurrency int

    // ProgressFunc is called to report progress
    ProgressFunc func(ProgressUpdate)
}
```

**ProgressUpdate:**
```go
type ProgressUpdate struct {
    Operation  string // "upload", "download", "delete", "skip"
    Path       string
    BytesDone  int64
    BytesTotal int64
    Error      error
}
```

**Methods:**
```go
func (e *Engine) Sync(ctx context.Context, sourcePath, destPath string) error
```

**Example:**
```go
// Create backends
localBackend, _ := sync.NewLocalBackend("/data/lab")
s3Backend, _ := sync.NewS3Backend(context.Background(), "lab-data")

// Create sync engine
engine := sync.NewEngine(localBackend, s3Backend, sync.SyncOptions{
    DryRun:      false,
    Delete:      false,
    Concurrency: 4,
    ProgressFunc: func(update sync.ProgressUpdate) {
        if update.Error != nil {
            log.Printf("Error: %s - %v", update.Path, update.Error)
        } else {
            log.Printf("%s: %s (%d/%d bytes)",
                update.Operation, update.Path,
                update.BytesDone, update.BytesTotal)
        }
    },
})

// Perform sync
err := engine.Sync(context.Background(), "", "microscopy/")
if err != nil {
    log.Fatal(err)
}
```

---

## metadata Package

Package `metadata` provides metadata extraction and validation.

### Extractor Interface

```go
type Extractor interface {
    // CanHandle returns true if this extractor can handle the given file
    CanHandle(filename string) bool

    // Extract extracts metadata from a file
    Extract(filepath string) (map[string]interface{}, error)

    // ExtractFromReader extracts metadata from a reader
    ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error)

    // Name returns the extractor name
    Name() string

    // SupportedFormats returns list of supported file extensions
    SupportedFormats() []string
}
```

### ExtractorRegistry

**Constructor:**
```go
func NewExtractorRegistry() *ExtractorRegistry
```

**Methods:**
```go
// RegisterDefaults registers all built-in extractors
func (r *ExtractorRegistry) RegisterDefaults()

// Register registers a custom extractor
func (r *ExtractorRegistry) Register(extractor Extractor)

// FindExtractor finds an extractor for the given filename
func (r *ExtractorRegistry) FindExtractor(filename string) Extractor

// Extract extracts metadata using the appropriate extractor
func (r *ExtractorRegistry) Extract(filepath string) (map[string]interface{}, error)

// ListExtractors returns all registered extractors
func (r *ExtractorRegistry) ListExtractors() []ExtractorInfo
```

**Example:**
```go
// Create registry
registry := metadata.NewExtractorRegistry()
registry.RegisterDefaults()

// Extract metadata
metadata, err := registry.Extract("/data/experiment001.czi")
if err != nil {
    log.Fatal(err)
}

// Access metadata fields
fmt.Printf("Format: %s\n", metadata["format"])
fmt.Printf("Width: %d\n", metadata["width"])
fmt.Printf("Height: %d\n", metadata["height"])
```

### Metadata Structure

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

### AutoExtractMetadata

```go
func AutoExtractMetadata(
    registry *ExtractorRegistry,
    filepath string,
    schemaName string,
    user string,
) (*Metadata, error)
```

**Example:**
```go
registry := metadata.NewExtractorRegistry()
registry.RegisterDefaults()

metadata, err := metadata.AutoExtractMetadata(
    registry,
    "/data/experiment001.czi",
    "microscopy",
    "jsmith",
)
if err != nil {
    log.Fatal(err)
}

// Metadata is wrapped in complete structure
fmt.Printf("Schema: %s v%s\n", metadata.SchemaName, metadata.SchemaVersion)
fmt.Printf("Uploaded by: %s\n", metadata.Provenance.UploadedBy)
```

### Preset

```go
type Preset struct {
    Name              string
    Description       string
    RequiredFields    []string
    RecommendedFields []string
    AllowedValues     map[string][]string
}
```

**Methods:**
```go
func (p *Preset) Validate(metadata map[string]interface{}) []ValidationError
```

**Example:**
```go
preset := metadata.Preset{
    Name:        "microscopy-confocal",
    Description: "Confocal microscopy requirements",
    RequiredFields: []string{
        "manufacturer",
        "model",
        "width",
        "height",
    },
    RecommendedFields: []string{
        "sample_id",
        "operator",
    },
}

// Extract and validate
extractedMeta, _ := registry.Extract("/data/image.czi")
errors := preset.Validate(extractedMeta)

for _, err := range errors {
    fmt.Printf("%s: %s\n", err.Level, err.Message)
}
```

### Built-in Extractors

**Available Extractors:**
- `TIFFExtractor` - TIFF images
- `OMETIFFExtractor` - OME-TIFF microscopy
- `ZeissCZIExtractor` - Zeiss CZI format
- `NikonExtractor` - Nikon ND2 format
- `LeicaExtractor` - Leica LIF format
- `FASTQExtractor` - FASTQ sequencing
- `BAMExtractor` - BAM alignment
- `MzMLExtractor` - mzML mass spec
- `MGFExtractor` - MGF mass spec
- `HDF5Extractor` - HDF5 arrays
- `ZarrExtractor` - Zarr arrays
- `DICOMExtractor` - DICOM medical imaging
- `FCSExtractor` - FCS flow cytometry
- `GenericExtractor` - Fallback extractor

---

## watch Package

Package `watch` provides file system monitoring and automatic sync.

### Config

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

**DefaultConfig:**
```go
func DefaultConfig() Config
```

### Watcher

**Constructor:**
```go
func New(config Config, engine *sync.Engine) (*Watcher, error)
```

**Methods:**
```go
func (w *Watcher) Start() error
func (w *Watcher) Stop() error
func (w *Watcher) Status() WatchStatus
```

**WatchStatus:**
```go
type WatchStatus struct {
    Source      string
    Destination string
    Active      bool
    StartedAt   time.Time
    LastSync    time.Time
    FilesSynced int64
    BytesSynced int64
    ErrorCount  int64
    LastError   string
}
```

**Example:**
```go
// Create backends
srcBackend, _ := sync.NewLocalBackend("/data/microscope")
dstBackend, _ := sync.NewS3Backend(context.Background(), "lab-data")

// Create sync engine
engine := sync.NewEngine(srcBackend, dstBackend, sync.SyncOptions{
    Concurrency: 4,
})

// Create watch config
config := watch.Config{
    Source:        "/data/microscope",
    Destination:   "microscopy/",
    DebounceDelay: 10 * time.Second,
    MinAge:        60 * time.Second,
    SyncOnStart:   true,
}

// Create and start watcher
watcher, err := watch.New(config, engine)
if err != nil {
    log.Fatal(err)
}

if err := watcher.Start(); err != nil {
    log.Fatal(err)
}

// Monitor status
ticker := time.NewTicker(30 * time.Second)
for range ticker.C {
    status := watcher.Status()
    log.Printf("Files synced: %d, Bytes: %d, Errors: %d",
        status.FilesSynced, status.BytesSynced, status.ErrorCount)
}

// Stop watching
defer watcher.Stop()
```

### Manager

**Constructor:**
```go
func NewManager() *Manager
```

**Methods:**
```go
func (m *Manager) Add(id string, config Config, srcBackend, dstBackend sync.Backend) error
func (m *Manager) Remove(id string) error
func (m *Manager) Get(id string) (*Watcher, bool)
func (m *Manager) List() map[string]WatchStatus
func (m *Manager) StopAll(ctx context.Context) error
```

**Example:**
```go
// Create manager
manager := watch.NewManager()

// Add multiple watches
manager.Add("microscope-1", config1, srcBackend1, dstBackend1)
manager.Add("sequencer-1", config2, srcBackend2, dstBackend2)

// List all watches
watches := manager.List()
for id, status := range watches {
    fmt.Printf("Watch %s: %d files synced\n", id, status.FilesSynced)
}

// Stop all watches
manager.StopAll(context.Background())
```

---

## config Package

Package `config` provides configuration management.

### Config

```go
type Config struct {
    Version  string         `yaml:"version"`
    AWS      AWSConfig      `yaml:"aws"`
    Sync     SyncConfig     `yaml:"sync"`
    Watches  []WatchConfig  `yaml:"watches"`
    Settings SettingsConfig `yaml:"settings"`
}

type AWSConfig struct {
    Profile  string `yaml:"profile"`
    Region   string `yaml:"region"`
    Endpoint string `yaml:"endpoint"`
}

type SyncConfig struct {
    Concurrency int      `yaml:"concurrency"`
    Delete      bool     `yaml:"delete"`
    Exclude     []string `yaml:"exclude"`
}

type SettingsConfig struct {
    Verbose      bool   `yaml:"verbose"`
    LogFile      string `yaml:"log_file"`
    CheckUpdates bool   `yaml:"check_updates"`
}
```

**Functions:**
```go
// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config

// ConfigPath returns the default config file path
func ConfigPath() (string, error)

// Load reads configuration from file
func Load(path string) (*Config, error)

// LoadOrDefault loads config from default path or returns default config
func LoadOrDefault() (*Config, error)

// Save writes configuration to file
func Save(config *Config, path string) error
```

**Example:**
```go
// Load configuration
cfg, err := config.LoadOrDefault()
if err != nil {
    log.Fatal(err)
}

// Access configuration
fmt.Printf("AWS Profile: %s\n", cfg.AWS.Profile)
fmt.Printf("AWS Region: %s\n", cfg.AWS.Region)
fmt.Printf("Sync Concurrency: %d\n", cfg.Sync.Concurrency)

// Modify configuration
cfg.Sync.Concurrency = 8
cfg.AWS.Region = "us-west-2"

// Save configuration
path, _ := config.ConfigPath()
if err := config.Save(cfg, path); err != nil {
    log.Fatal(err)
}
```

---

## doi Package

Package `doi` provides DOI preparation and provider integration.

### Provider Interface

```go
type Provider interface {
    // Validate validates DOI metadata
    Validate(metadata DataciteMetadata) error

    // Create creates a new DOI
    Create(metadata DataciteMetadata) (*DOIRecord, error)

    // Update updates an existing DOI
    Update(doi string, metadata DataciteMetadata) (*DOIRecord, error)

    // Get retrieves DOI information
    Get(doi string) (*DOIRecord, error)

    // Delete deletes a DOI (if supported)
    Delete(doi string) error
}
```

### DataciteMetadata

```go
type DataciteMetadata struct {
    Identifier  Identifier   `json:"identifier"`
    Creators    []Creator    `json:"creators"`
    Titles      []Title      `json:"titles"`
    Publisher   string       `json:"publisher"`
    PublicationYear int      `json:"publicationYear"`
    ResourceType ResourceType `json:"resourceType"`
    // ... additional fields
}

type Creator struct {
    Name             string           `json:"name"`
    NameType         string           `json:"nameType,omitempty"`
    GivenName        string           `json:"givenName,omitempty"`
    FamilyName       string           `json:"familyName,omitempty"`
    Affiliation      []Affiliation    `json:"affiliation,omitempty"`
    NameIdentifiers  []NameIdentifier `json:"nameIdentifiers,omitempty"`
}
```

### Workflow

**Functions:**
```go
// PrepareMetadata prepares DOI metadata from file metadata
func PrepareMetadata(fileMetadata map[string]interface{}, options PrepareOptions) (*DataciteMetadata, error)

// ValidateMetadata validates DOI metadata against DataCite schema
func ValidateMetadata(metadata *DataciteMetadata) []ValidationError
```

**Example:**
```go
// Extract file metadata
registry := metadata.NewExtractorRegistry()
registry.RegisterDefaults()
fileMeta, _ := registry.Extract("/data/dataset.czi")

// Prepare DOI metadata
doiMeta, err := doi.PrepareMetadata(fileMeta, doi.PrepareOptions{
    Title:     "Neural Imaging Dataset",
    Creators:  []string{"Smith, John", "Doe, Jane"},
    Publisher: "University Lab",
    Year:      2025,
})
if err != nil {
    log.Fatal(err)
}

// Validate
errors := doi.ValidateMetadata(doiMeta)
for _, err := range errors {
    fmt.Printf("Validation error: %s\n", err.Message)
}
```

---

## Usage Examples

### Example 1: Simple Sync

```go
package main

import (
    "context"
    "log"

    "github.com/scttfrdmn/cicada/internal/sync"
)

func main() {
    ctx := context.Background()

    // Create backends
    local, _ := sync.NewLocalBackend("/data/lab")
    s3, _ := sync.NewS3Backend(ctx, "lab-data")
    defer local.Close()
    defer s3.Close()

    // Create engine
    engine := sync.NewEngine(local, s3, sync.SyncOptions{
        Concurrency: 4,
        Delete:      false,
    })

    // Sync
    if err := engine.Sync(ctx, "", "backup/"); err != nil {
        log.Fatal(err)
    }

    log.Println("Sync complete")
}
```

### Example 2: Metadata Extraction

```go
package main

import (
    "encoding/json"
    "log"
    "os"

    "github.com/scttfrdmn/cicada/internal/metadata"
)

func main() {
    // Create registry
    registry := metadata.NewExtractorRegistry()
    registry.RegisterDefaults()

    // Extract metadata
    meta, err := registry.Extract("/data/experiment001.czi")
    if err != nil {
        log.Fatal(err)
    }

    // Save to JSON
    f, _ := os.Create("metadata.json")
    defer f.Close()

    encoder := json.NewEncoder(f)
    encoder.SetIndent("", "  ")
    encoder.Encode(meta)

    log.Println("Metadata extracted")
}
```

### Example 3: Watch Directory

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "time"

    "github.com/scttfrdmn/cicada/internal/sync"
    "github.com/scttfrdmn/cicada/internal/watch"
)

func main() {
    ctx := context.Background()

    // Create backends
    src, _ := sync.NewLocalBackend("/data/microscope")
    dst, _ := sync.NewS3Backend(ctx, "lab-data")

    // Create engine
    engine := sync.NewEngine(src, dst, sync.SyncOptions{
        Concurrency: 4,
        ProgressFunc: func(u sync.ProgressUpdate) {
            log.Printf("%s: %s", u.Operation, u.Path)
        },
    })

    // Create watcher
    watcher, err := watch.New(watch.Config{
        Source:        "/data/microscope",
        Destination:   "microscopy/",
        DebounceDelay: 10 * time.Second,
        MinAge:        60 * time.Second,
        SyncOnStart:   true,
    }, engine)
    if err != nil {
        log.Fatal(err)
    }

    // Start watching
    if err := watcher.Start(); err != nil {
        log.Fatal(err)
    }

    log.Println("Watching directory...")

    // Wait for interrupt
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c

    // Stop watcher
    watcher.Stop()
    log.Println("Stopped")
}
```

### Example 4: Complete Pipeline

```go
package main

import (
    "context"
    "log"

    "github.com/scttfrdmn/cicada/internal/metadata"
    "github.com/scttfrdmn/cicada/internal/sync"
)

func main() {
    ctx := context.Background()

    // 1. Extract metadata
    registry := metadata.NewExtractorRegistry()
    registry.RegisterDefaults()

    meta, err := registry.Extract("/data/experiment001.czi")
    if err != nil {
        log.Fatal(err)
    }

    // 2. Validate metadata
    preset := metadata.Preset{
        Name: "microscopy-confocal",
        RequiredFields: []string{"width", "height", "channels"},
    }

    errors := preset.Validate(meta)
    if len(errors) > 0 {
        log.Fatal("Validation failed")
    }

    // 3. Sync to S3
    local, _ := sync.NewLocalBackend("/data")
    s3, _ := sync.NewS3Backend(ctx, "lab-data")
    defer local.Close()
    defer s3.Close()

    engine := sync.NewEngine(local, s3, sync.SyncOptions{
        Concurrency: 4,
    })

    if err := engine.Sync(ctx, "experiment001.czi", "validated/"); err != nil {
        log.Fatal(err)
    }

    log.Println("Pipeline complete")
}
```

---

## Error Handling

### Error Types

**Common Errors:**
- `ErrNotFound` - File or resource not found
- `ErrAccessDenied` - Permission denied
- `ErrInvalidPath` - Invalid path format
- `ErrNetworkError` - Network connectivity issue
- `ErrValidationFailed` - Metadata validation failed

### Error Wrapping

Cicada uses Go's error wrapping (`%w`):

```go
if err != nil {
    return fmt.Errorf("sync failed: %w", err)
}
```

**Unwrapping:**
```go
if errors.Is(err, sync.ErrNotFound) {
    // Handle not found
}

if errors.As(err, &networkErr) {
    // Handle network error
}
```

### Error Handling Best Practices

```go
// Check for specific errors
if err := backend.Read(ctx, path); err != nil {
    if errors.Is(err, sync.ErrNotFound) {
        log.Printf("File not found: %s", path)
        return nil // Not a fatal error
    }
    return fmt.Errorf("read failed: %w", err)
}

// Handle all errors explicitly
result, err := operation()
if err != nil {
    // Always handle or propagate
    return fmt.Errorf("operation failed: %w", err)
}
```

---

## Related Documentation

- [Development Guide](DEVELOPMENT.md) - Development setup and guidelines
- [Architecture](ARCHITECTURE.md) - System architecture overview
- [CLI Reference](CLI_REFERENCE.md) - Command-line interface

---

**Questions?** Open an issue on GitHub.

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
