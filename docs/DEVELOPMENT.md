# Development Guide

**Last Updated:** 2025-11-25

Complete guide for developers contributing to Cicada.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Development Environment](#development-environment)
3. [Building from Source](#building-from-source)
4. [Running Tests](#running-tests)
5. [Code Organization](#code-organization)
6. [Contributing Guidelines](#contributing-guidelines)
7. [Adding New Features](#adding-new-features)
8. [Release Process](#release-process)
9. [Debugging](#debugging)
10. [Performance Profiling](#performance-profiling)

---

## Getting Started

### Prerequisites

**Required:**
- Go 1.21 or later
- Git
- Make (optional but recommended)

**Recommended:**
- AWS CLI (for S3 testing)
- Docker (for containerized testing)
- golangci-lint (for linting)

### Quick Start

```bash
# Clone repository
git clone https://github.com/scttfrdmn/cicada.git
cd cicada

# Install dependencies
go mod download

# Build
go build -o cicada cmd/cicada/main.go

# Run tests
go test ./...

# Verify build
./cicada version
```

---

## Development Environment

### Go Environment

**Install Go:**
```bash
# macOS
brew install go

# Linux (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install golang-go

# Or download from https://go.dev/dl/
```

**Verify Installation:**
```bash
go version
# Should show: go version go1.21.x
```

**Configure GOPATH (if not already set):**
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### IDE Setup

#### VS Code

**Recommended Extensions:**
- Go (golang.go)
- Go Test Explorer (premparihar.gotestexplorer)
- GitLens (eamodio.gitlens)

**Settings (.vscode/settings.json):**
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "go.formatTool": "gofmt",
    "go.testFlags": ["-v"],
    "go.coverOnSave": true,
    "editor.formatOnSave": true
}
```

#### GoLand / IntelliJ IDEA

1. Open project
2. Enable Go modules: Settings → Go → Go Modules
3. Set GOROOT: Settings → Go → GOROOT
4. Enable tests: Settings → Go → Test

### Development Tools

**golangci-lint (Linter):**
```bash
# Install
brew install golangci-lint

# Or
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run
golangci-lint run
```

**go-mockgen (Mock Generation):**
```bash
go install github.com/golang/mock/mockgen@latest
```

**AWS CLI (for S3 testing):**
```bash
# macOS
brew install awscli

# Configure
aws configure
```

---

## Building from Source

### Standard Build

```bash
# Build for current platform
go build -o cicada cmd/cicada/main.go

# Build with version information
go build -ldflags "\
  -X main.version=$(git describe --tags --always --dirty) \
  -X main.commit=$(git rev-parse --short HEAD) \
  -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X main.builtBy=dev" \
  -o cicada cmd/cicada/main.go

# Verify
./cicada version
```

### Cross-Platform Builds

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o cicada-linux-amd64 cmd/cicada/main.go

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o cicada-linux-arm64 cmd/cicada/main.go

# macOS (amd64)
GOOS=darwin GOARCH=amd64 go build -o cicada-darwin-amd64 cmd/cicada/main.go

# macOS (arm64 - M1/M2)
GOOS=darwin GOARCH=arm64 go build -o cicada-darwin-arm64 cmd/cicada/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o cicada-windows-amd64.exe cmd/cicada/main.go
```

### Build with Make

**Makefile:**
```makefile
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o cicada cmd/cicada/main.go

.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/cicada-linux-amd64 cmd/cicada/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/cicada-linux-arm64 cmd/cicada/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/cicada-darwin-amd64 cmd/cicada/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/cicada-darwin-arm64 cmd/cicada/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/cicada-windows-amd64.exe cmd/cicada/main.go

.PHONY: test
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -f cicada
	rm -rf dist/
	rm -f coverage.txt
```

**Usage:**
```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make test         # Run tests
make lint         # Run linter
make clean        # Clean build artifacts
```

---

## Running Tests

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/sync

# Run specific test
go test -v -run TestSyncEngine ./internal/sync

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.txt -covermode=atomic ./...

# View coverage report
go tool cover -html=coverage.txt
```

### Integration Tests

```bash
# Run integration tests (requires AWS credentials)
go test -v ./internal/integration

# Run with tags
go test -v -tags=integration ./...

# Skip integration tests
go test -v -short ./...
```

### Benchmark Tests

```bash
# Run benchmarks
go test -bench=. ./internal/sync

# Run specific benchmark
go test -bench=BenchmarkSyncEngine ./internal/sync

# With memory profiling
go test -bench=. -benchmem ./internal/sync

# Save benchmark results
go test -bench=. ./internal/sync > bench.txt
```

### Test Structure

**Unit Test Example:**
```go
// internal/sync/engine_test.go
package sync

import (
    "context"
    "testing"
)

func TestSyncEngine_Sync(t *testing.T) {
    // Setup
    ctx := context.Background()
    source := NewMockBackend()
    dest := NewMockBackend()

    engine := NewEngine(source, dest, SyncOptions{
        Concurrency: 4,
        DryRun:      false,
    })

    // Execute
    err := engine.Sync(ctx, "/source", "/dest")

    // Assert
    if err != nil {
        t.Fatalf("Sync failed: %v", err)
    }

    // Verify expectations
    if source.ListCallCount != 1 {
        t.Errorf("Expected 1 List call, got %d", source.ListCallCount)
    }
}
```

**Table-Driven Test Example:**
```go
func TestNeedsSync(t *testing.T) {
    tests := []struct {
        name     string
        src      FileInfo
        dst      FileInfo
        expected bool
    }{
        {
            name: "same etag - no sync",
            src:  FileInfo{ETag: "abc123", Size: 1000},
            dst:  FileInfo{ETag: "abc123", Size: 1000},
            expected: false,
        },
        {
            name: "different etag - sync needed",
            src:  FileInfo{ETag: "abc123", Size: 1000},
            dst:  FileInfo{ETag: "xyz789", Size: 1000},
            expected: true,
        },
        {
            name: "different size - sync needed",
            src:  FileInfo{Size: 1000},
            dst:  FileInfo{Size: 2000},
            expected: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := needsSync(tt.src, tt.dst)
            if result != tt.expected {
                t.Errorf("needsSync() = %v, expected %v", result, tt.expected)
            }
        })
    }
}
```

---

## Code Organization

### Repository Structure

```
cicada/
├── cmd/
│   └── cicada/
│       └── main.go              # CLI entry point
├── internal/
│   ├── cli/                     # CLI commands
│   │   ├── root.go
│   │   ├── sync.go
│   │   ├── watch.go
│   │   ├── metadata.go
│   │   ├── doi.go
│   │   └── config.go
│   ├── sync/                    # Sync engine & backends
│   │   ├── backend.go           # Backend interface
│   │   ├── engine.go            # Sync engine
│   │   ├── local.go             # Local backend
│   │   ├── s3.go                # S3 backend
│   │   └── *_test.go            # Tests
│   ├── metadata/                # Metadata system
│   │   ├── extractor.go         # Extractor interface
│   │   ├── types.go             # Metadata types
│   │   ├── preset.go            # Presets
│   │   ├── schema.go            # Schema
│   │   ├── fastq.go             # FASTQ extractor
│   │   ├── ome_tiff.go          # OME-TIFF extractor
│   │   └── ...                  # Other extractors
│   ├── watch/                   # Watch daemon
│   │   ├── manager.go
│   │   ├── watcher.go
│   │   └── debouncer.go
│   ├── doi/                     # DOI system
│   │   ├── workflow.go
│   │   ├── provider.go
│   │   └── ...
│   ├── config/                  # Configuration
│   │   ├── config.go
│   │   └── providers.go
│   └── integration/             # Integration tests
├── docs/                        # Documentation
├── testdata/                    # Test fixtures
├── go.mod                       # Go modules
├── go.sum                       # Go checksums
├── Makefile                     # Build automation
├── .golangci.yml               # Linter config
└── README.md
```

### Package Guidelines

**internal/ Packages:**
- Not importable by external projects
- Use for implementation details
- Keep packages focused and cohesive

**Package Dependencies:**
- Avoid circular dependencies
- Core packages (sync, metadata) should not depend on CLI
- CLI can import all internal packages

**Naming Conventions:**
- Use descriptive package names: `sync`, `metadata`, not `util`, `common`
- Use short, lowercase package names
- Package name should match directory name

---

## Contributing Guidelines

### Code Style

**Follow Go conventions:**
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Follow [Effective Go](https://go.dev/doc/effective_go)

**Comments:**
```go
// Good: Clear, concise, explains why
// Sync performs bi-directional synchronization between backends.
// Files are compared using ETags when available, falling back to
// size and modification time comparison.
func (e *Engine) Sync(ctx context.Context, src, dst string) error {
    // ...
}

// Bad: Restates what code does
// Sync syncs files
func (e *Engine) Sync(ctx context.Context, src, dst string) error {
    // ...
}
```

**Error Handling:**
```go
// Good: Wrap errors with context
if err := backend.Write(ctx, path, data); err != nil {
    return fmt.Errorf("write to backend: %w", err)
}

// Bad: Lost context
if err := backend.Write(ctx, path, data); err != nil {
    return err
}
```

**Testing:**
- Write tests for all new features
- Aim for >80% code coverage
- Include table-driven tests for multiple cases
- Use meaningful test names

### Git Workflow

**Branching:**
```bash
# Create feature branch
git checkout -b feature/my-feature

# Create fix branch
git checkout -b fix/bug-description
```

**Commits:**
```bash
# Good commit message
git commit -m "Add multipart upload support for S3 backend

- Implement chunked upload for files >100MB
- Add progress tracking for multipart uploads
- Include tests for upload resumption"

# Bad commit message
git commit -m "fix stuff"
```

**Commit Message Format:**
```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `test`: Tests
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `chore`: Maintenance tasks

### Pull Request Process

1. **Fork and Clone:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/cicada.git
   ```

2. **Create Branch:**
   ```bash
   git checkout -b feature/my-feature
   ```

3. **Make Changes:**
   - Write code
   - Add tests
   - Update documentation

4. **Test:**
   ```bash
   go test ./...
   golangci-lint run
   ```

5. **Commit:**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

6. **Push:**
   ```bash
   git push origin feature/my-feature
   ```

7. **Create Pull Request:**
   - Go to GitHub
   - Click "New Pull Request"
   - Fill out PR template
   - Request review

**PR Checklist:**
- [ ] Tests pass locally
- [ ] Tests added for new features
- [ ] Documentation updated
- [ ] Code follows style guidelines
- [ ] No merge conflicts
- [ ] PR description is clear

---

## Adding New Features

### Adding a New Metadata Extractor

**Step 1: Create Extractor File**

```go
// internal/metadata/custom_format.go
package metadata

import (
    "io"
    "path/filepath"
    "strings"
)

type CustomFormatExtractor struct{}

func (e *CustomFormatExtractor) Name() string {
    return "Custom Format"
}

func (e *CustomFormatExtractor) SupportedFormats() []string {
    return []string{".custom", ".cst"}
}

func (e *CustomFormatExtractor) CanHandle(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    for _, format := range e.SupportedFormats() {
        if ext == format {
            return true
        }
    }
    return false
}

func (e *CustomFormatExtractor) Extract(filepath string) (map[string]interface{}, error) {
    // Read and parse file
    data, err := parseCustomFormat(filepath)
    if err != nil {
        return nil, fmt.Errorf("parse custom format: %w", err)
    }

    // Return metadata
    return map[string]interface{}{
        "format":     "CUSTOM",
        "version":    data.Version,
        "instrument": data.Instrument,
        // ... additional fields
    }, nil
}

func (e *CustomFormatExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
    // Implement streaming extraction if needed
    return nil, fmt.Errorf("not implemented")
}
```

**Step 2: Register Extractor**

```go
// internal/metadata/extractor.go
func (r *ExtractorRegistry) RegisterDefaults() {
    // ... existing extractors
    r.Register(&CustomFormatExtractor{})
}
```

**Step 3: Add Tests**

```go
// internal/metadata/custom_format_test.go
package metadata

import "testing"

func TestCustomFormatExtractor_CanHandle(t *testing.T) {
    extractor := &CustomFormatExtractor{}

    tests := []struct {
        filename string
        expected bool
    }{
        {"test.custom", true},
        {"test.cst", true},
        {"test.txt", false},
    }

    for _, tt := range tests {
        result := extractor.CanHandle(tt.filename)
        if result != tt.expected {
            t.Errorf("CanHandle(%s) = %v, expected %v",
                tt.filename, result, tt.expected)
        }
    }
}

func TestCustomFormatExtractor_Extract(t *testing.T) {
    extractor := &CustomFormatExtractor{}

    metadata, err := extractor.Extract("testdata/sample.custom")
    if err != nil {
        t.Fatalf("Extract failed: %v", err)
    }

    if metadata["format"] != "CUSTOM" {
        t.Errorf("Expected format CUSTOM, got %v", metadata["format"])
    }
}
```

**Step 4: Add Test Data**

```bash
# Create test fixture
mkdir -p testdata
# Add sample.custom file
```

**Step 5: Update Documentation**

```markdown
# docs/METADATA_SYSTEM.md

### Custom Format

**Description:** Custom scientific instrument format

**Extracted Fields:**
- Format version
- Instrument details
- ...

**Example:**
...
```

### Adding a New Storage Backend

**Step 1: Implement Backend Interface**

```go
// internal/sync/azure.go
package sync

import (
    "context"
    "io"
)

type AzureBackend struct {
    // Azure-specific fields
    accountName string
    accountKey  string
}

func NewAzureBackend(accountName, accountKey string) (*AzureBackend, error) {
    return &AzureBackend{
        accountName: accountName,
        accountKey:  accountKey,
    }, nil
}

func (b *AzureBackend) List(ctx context.Context, prefix string) ([]FileInfo, error) {
    // Implement Azure blob listing
    return nil, nil
}

func (b *AzureBackend) Read(ctx context.Context, path string) (io.ReadCloser, error) {
    // Implement Azure blob download
    return nil, nil
}

func (b *AzureBackend) Write(ctx context.Context, path string, r io.Reader, size int64) error {
    // Implement Azure blob upload
    return nil
}

func (b *AzureBackend) Delete(ctx context.Context, path string) error {
    // Implement Azure blob deletion
    return nil
}

func (b *AzureBackend) Stat(ctx context.Context, path string) (*FileInfo, error) {
    // Implement Azure blob metadata
    return nil, nil
}

func (b *AzureBackend) Close() error {
    // Cleanup resources
    return nil
}
```

**Step 2: Add Backend Factory**

```go
// internal/cli/sync.go
func createBackend(ctx context.Context, path string) (sync.Backend, string, error) {
    switch {
    case strings.HasPrefix(path, "s3://"):
        // S3 backend
        bucket, key, _ := sync.ParseS3URI(path)
        backend, _ := sync.NewS3Backend(ctx, bucket)
        return backend, key, nil

    case strings.HasPrefix(path, "azure://"):
        // Azure backend
        container, blob, _ := sync.ParseAzureURI(path)
        backend, _ := sync.NewAzureBackend(container, blob)
        return backend, blob, nil

    default:
        // Local backend
        backend, _ := sync.NewLocalBackend(path)
        return backend, "", nil
    }
}
```

**Step 3: Add Tests**

```go
// internal/sync/azure_test.go
package sync

import "testing"

func TestAzureBackend_List(t *testing.T) {
    backend, _ := NewAzureBackend("account", "key")

    files, err := backend.List(context.Background(), "prefix/")
    if err != nil {
        t.Fatalf("List failed: %v", err)
    }

    if len(files) == 0 {
        t.Error("Expected files, got none")
    }
}
```

---

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., `0.2.0`)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Creating a Release

**Step 1: Update Version**

```bash
# Update version in code
git tag v0.3.0
```

**Step 2: Generate Changelog**

```bash
# Generate changelog from git commits
git log v0.2.0..HEAD --oneline --pretty=format:"- %s" > CHANGELOG-v0.3.0.md
```

**Step 3: Build Release Binaries**

```bash
make build-all
```

**Step 4: Create GitHub Release**

```bash
# Create release with gh CLI
gh release create v0.3.0 \
  --title "v0.3.0 - Documentation Release" \
  --notes-file CHANGELOG-v0.3.0.md \
  dist/cicada-*
```

**Step 5: Update Documentation**

```bash
# Update version in docs
# Commit and push
git push origin main --tags
```

---

## Debugging

### Debug Logging

```go
// Enable verbose logging
cicada --verbose sync /data s3://bucket/data
```

### Using Delve (Go Debugger)

**Install:**
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Debug:**
```bash
# Debug test
dlv test ./internal/sync

# Debug binary
dlv debug cmd/cicada/main.go -- sync /data s3://bucket/data

# Attach to running process
dlv attach <pid>
```

**Delve Commands:**
```
break main.main        # Set breakpoint
continue               # Continue execution
next                   # Step over
step                   # Step into
print var              # Print variable
list                   # Show source
goroutines             # List goroutines
```

### Logging Best Practices

```go
import "log"

// Use structured logging
log.Printf("Syncing %d files from %s to %s", count, src, dst)

// Log errors with context
log.Printf("ERROR: Failed to sync file %s: %v", path, err)

// Use verbose flag for detailed logs
if verbose {
    log.Printf("DEBUG: Processing file %s", file)
}
```

---

## Performance Profiling

### CPU Profiling

```bash
# Generate CPU profile
go test -cpuprofile=cpu.prof -bench=. ./internal/sync

# Analyze with pprof
go tool pprof cpu.prof

# Commands in pprof:
# top        - Show top functions
# list <fn>  - Show source for function
# web        - Open graph in browser
```

### Memory Profiling

```bash
# Generate memory profile
go test -memprofile=mem.prof -bench=. ./internal/sync

# Analyze
go tool pprof mem.prof
```

### Benchmarking

```go
func BenchmarkSyncEngine(b *testing.B) {
    // Setup
    engine := NewEngine(source, dest, SyncOptions{})

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine.Sync(context.Background(), "/src", "/dst")
    }
}
```

**Run Benchmarks:**
```bash
go test -bench=. -benchmem ./internal/sync
```

---

## Related Documentation

- [Architecture](ARCHITECTURE.md) - System architecture
- [Contributing](../CONTRIBUTING.md) - Contribution guidelines
- [API Reference](API.md) - API documentation

---

**Questions?** Open an issue on GitHub or join our discussions.

**License:** Apache 2.0 - See [LICENSE](../LICENSE)
