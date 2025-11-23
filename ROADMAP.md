# Cicada Development Roadmap

Detailed implementation plan with priorities, technical specifications, and milestones.

## Overview

**Total Timeline**: ~26 weeks (6 months)  
**Team Size**: 1-3 developers  
**Target**: MVP launch in 3 months, full feature set in 6 months

---

## Phase 1: Core Storage & Sync (Weeks 1-6)

### Priorities
1. **CRITICAL**: S3 sync engine
2. **CRITICAL**: File watching
3. **CRITICAL**: Basic CLI
4. **HIGH**: Daemon service
5. **MEDIUM**: Cost tracking

### Technical Specifications

#### 1.1 S3 Sync Engine (Week 1-2)

**Location**: `internal/sync/`

**Key Components**:
```go
// internal/sync/engine.go
type SyncEngine struct {
    source      Backend
    destination Backend
    options     SyncOptions
    reporter    ProgressReporter
}

type Backend interface {
    List(prefix string) ([]FileInfo, error)
    Read(path string) (io.ReadCloser, error)
    Write(path string, r io.Reader) error
    Delete(path string) error
    Checksum(path string) (string, error)
}
```

**Features**:
- [ ] rsync-like delta detection (checksum-based)
- [ ] Multipart upload for files >100MB
- [ ] Resume capability for interrupted transfers
- [ ] Parallel transfers (configurable concurrency)
- [ ] Bandwidth throttling
- [ ] Progress reporting (bytes, files, ETA)
- [ ] Dry-run mode
- [ ] Exclude/include patterns
- [ ] Delete synchronization

**Dependencies**:
```go
require (
    github.com/aws/aws-sdk-go-v2 v1.24.0
    github.com/aws/aws-sdk-go-v2/service/s3 v1.44.0
    github.com/cheggaaa/pb/v3 v3.1.4  // Progress bars
)
```

**Tests**:
- Unit tests for delta calculation
- Integration tests with LocalStack (mock S3)
- Performance benchmarks (1GB, 10GB, 100GB datasets)
- Interruption/resume tests

**Acceptance Criteria**:
- Sync 10GB in <15 minutes on 100Mbps connection
- Handle network interruptions gracefully
- Correctly detect and skip unchanged files

---

#### 1.2 File Watching (Week 2-3)

**Location**: `internal/watch/`

**Key Components**:
```go
// internal/watch/watcher.go
type Watcher struct {
    locations map[string]*WatchConfig
    fsWatcher *fsnotify.Watcher
    syncer    *sync.Engine
    debouncer *Debouncer
}

type WatchConfig struct {
    Path          string
    Destination   string
    MinAge        time.Duration
    OnNewFile     bool
    Schedule      string  // cron format
    DeleteSource  bool
    IgnorePattern []string
}
```

**Features**:
- [ ] File system event watching (fsnotify)
- [ ] Debouncing (wait for writes to complete)
- [ ] Age-based filtering (don't sync files too new)
- [ ] Pattern-based filtering (ignore .tmp, etc.)
- [ ] Scheduled syncs (cron-like)
- [ ] Multiple watch locations
- [ ] Persistent state (resume after restart)

**Dependencies**:
```go
require (
    github.com/fsnotify/fsnotify v1.7.0
    github.com/robfig/cron/v3 v3.0.1
)
```

**Tests**:
- Watch folder, create file, verify sync
- Test debouncing (rapid file changes)
- Test scheduled sync
- Test crash recovery

**Acceptance Criteria**:
- Detect new files within 5 seconds
- Wait for configurable age before syncing
- Handle thousands of file events without crashing
- Persist watch state across restarts

---

#### 1.3 Daemon Service (Week 3-4)

**Location**: `cmd/cicada-daemon/`

**Key Components**:
```go
// cmd/cicada-daemon/main.go
type Daemon struct {
    config      *config.Config
    watcher     *watch.Watcher
    syncer      *sync.Engine
    webServer   *webui.Server
    httpClient  *http.Client
}
```

**Features**:
- [ ] Background service (systemd on Linux, launchd on macOS, service on Windows)
- [ ] PID file management
- [ ] Graceful shutdown
- [ ] Log rotation
- [ ] Health check endpoint
- [ ] Status reporting

**Platform Support**:
- Linux: systemd service
- macOS: launchd plist
- Windows: Windows Service

**Tests**:
- Start/stop daemon
- Restart with active watches
- Graceful shutdown during transfer
- Log rotation

---

#### 1.4 Basic CLI (Week 4-5)

**Location**: `cmd/cicada/`

**Key Commands**:
```bash
cicada init              # Setup wizard
cicada sync              # Manual sync
cicada watch add/remove  # Manage watches
cicada daemon start/stop # Daemon control
cicada status            # Overall status
```

**Framework**: Cobra CLI

**Features**:
- [ ] Colored output (terminal detection)
- [ ] Progress bars
- [ ] Interactive prompts
- [ ] Tab completion
- [ ] Configuration file support (YAML)

**Tests**:
- Command parsing
- Output formatting
- Interactive mode

---

#### 1.5 Cost Tracking (Week 5-6)

**Location**: `internal/cost/`

**Key Components**:
```go
// internal/cost/tracker.go
type CostTracker struct {
    cloudwatch *cloudwatch.Client
    ce         *costexplorer.Client
    storage    CostStorage
}

type CostRecord struct {
    Date        time.Time
    Service     string  // s3, ec2, batch, etc.
    Amount      float64
    Units       string
    Description string
}
```

**Features**:
- [ ] Query AWS Cost Explorer API
- [ ] Breakdown by service
- [ ] Trend analysis
- [ ] Budget alerts
- [ ] Cost estimation for operations
- [ ] Local cost cache (reduce API calls)

**Tests**:
- Mock CloudWatch/Cost Explorer
- Cost calculation accuracy
- Trend detection

---

### Milestone 1: Basic Data Management (End of Week 6)

**Deliverables**:
- Working CLI that can sync data
- Daemon that watches folders
- Basic cost tracking

**Demo**:
```bash
# Setup
cicada init --lab-name demo-lab

# Add watch
cicada watch add demo \
  --path /tmp/demo/source \
  --destination s3://demo-bucket/data/

# Create some files
echo "test" > /tmp/demo/source/file1.txt

# Verify sync happened
cicada status  # Should show 1 file synced

# Check costs
cicada cost report  # Should show S3 costs
```

---

## Phase 2: Metadata & FAIR (Weeks 7-10)

### Priorities
1. **CRITICAL**: Schema system
2. **CRITICAL**: Metadata extraction
3. **HIGH**: Validation
4. **HIGH**: Search/discovery
5. **MEDIUM**: Export formats

### Technical Specifications

#### 2.1 Metadata Schema System (Week 7-8)

**Location**: `internal/metadata/schema.go`

See artifact: `/tmp/cicada-spec/internal/metadata/schema.go`

**Features**:
- [ ] YAML schema definition
- [ ] Schema inheritance (extends)
- [ ] Field type validation
- [ ] Controlled vocabularies
- [ ] Ontology integration
- [ ] Custom validation rules
- [ ] Schema versioning

**Domain Schemas** (Week 8):
- [ ] Microscopy (fluorescence, confocal, etc.)
- [ ] Sequencing (RNA-seq, DNA-seq, etc.)
- [ ] Proteomics (mass spec)
- [ ] Flow cytometry
- [ ] Chromatography
- [ ] Spectroscopy

**Tests**:
- Schema parsing
- Validation logic
- Inheritance resolution
- Type checking

---

#### 2.2 Metadata Extraction (Week 8-9)

**Location**: `internal/metadata/extractor.go`

See artifact: `/tmp/cicada-spec/internal/metadata/extractor.go`

**Extractors to Implement**:

**Priority 1 (Week 8)**:
- [ ] TIFF/OME-TIFF (microscopy)
- [ ] FASTQ (sequencing)
- [ ] Generic (fallback)

**Priority 2 (Week 9)**:
- [ ] Zeiss CZI
- [ ] Nikon ND2
- [ ] BAM/SAM
- [ ] HDF5

**Priority 3 (Later)**:
- [ ] DICOM
- [ ] FCS (flow cytometry)
- [ ] mzML (mass spec)

**Libraries Needed**:
```go
require (
    github.com/biogo/hts v1.4.4  // BAM/SAM
    github.com/gonum/hdf5 v0.0.0  // HDF5 (if available)
    // Note: Many formats may need external tools or Python integration
)
```

**Tests**:
- Test files for each format
- Extraction accuracy
- Performance (large files)

---

#### 2.3 Search & Discovery (Week 9-10)

**Location**: `internal/metadata/search.go`

**Key Components**:
```go
type SearchEngine struct {
    index     Index
    storage   MetadataStorage
    analyzer  TextAnalyzer
}

type SearchQuery struct {
    Terms     []string
    Filters   map[string]interface{}
    DateRange *DateRange
    Facets    []string
}
```

**Backend Options**:
- **Option 1**: SQLite FTS (simplest, local)
- **Option 2**: DynamoDB (cloud-native)
- **Option 3**: Elasticsearch (powerful, complex)

**Recommendation**: Start with SQLite FTS, add DynamoDB later

**Features**:
- [ ] Full-text search
- [ ] Faceted search
- [ ] Date range filters
- [ ] Field-specific queries
- [ ] Relevance ranking

**Tests**:
- Index creation
- Query accuracy
- Performance (10k, 100k, 1M records)

---

#### 2.4 Export Formats (Week 10)

**Location**: `internal/metadata/export.go`

**Formats**:
- [ ] DataCite XML
- [ ] ISA-Tab
- [ ] DATS JSON
- [ ] Frictionless Data Package
- [ ] RO-Crate
- [ ] Dublin Core

**Tests**:
- Schema validation for each format
- Round-trip conversion where applicable

---

### Milestone 2: FAIR-Compliant Data (End of Week 10)

**Deliverables**:
- Metadata extraction on upload
- Schema-based validation
- Searchable metadata
- Export to standard formats

**Demo**:
```bash
# Upload with metadata extraction
cicada upload microscopy/exp_001.czi s3://bucket/data/ \
  --schema fluorescence-microscopy

# Search
cicada search --where "magnification=63"

# Export
cicada metadata export s3://bucket/data/ --format datacite
```

---

## Phase 3: Web UI & User Management (Weeks 11-14)

### Priorities
1. **CRITICAL**: Web server
2. **CRITICAL**: File browser
3. **HIGH**: Upload/download UI
4. **HIGH**: User management
5. **MEDIUM**: Project management

### Technical Specifications

#### 3.1 Web Server (Week 11)

**Location**: `internal/webui/`

**Framework**: Go standard library + chi router (or gin)

**Key Components**:
```go
// internal/webui/server.go
type Server struct {
    daemon   *daemon.Daemon
    router   *chi.Mux
    wsHub    *WebSocketHub
    sessions *sessions.Manager
}

// API endpoints
GET  /api/status
GET  /api/files/*path
POST /api/upload
GET  /api/projects
POST /api/sync
GET  /api/ws  // WebSocket for real-time updates
```

**Features**:
- [ ] REST API
- [ ] WebSocket for real-time updates
- [ ] Session management
- [ ] CORS handling
- [ ] Static file serving (embedded UI)

---

#### 3.2 Frontend (Week 11-12)

**Technology**: Svelte (lightweight, fast)

**Alternative**: Vue 3 (more ecosystem support)

**Pages**:
- Dashboard (overview, recent activity)
- File Browser (navigate S3, preview files)
- Upload (drag-drop, progress bars)
- Sync Manager (watch locations, status)
- Metadata (view/edit)
- Projects (list, details)
- Cost Dashboard

**Build Process**:
```bash
cd web
npm install
npm run build  # Outputs to internal/webui/static
```

**Embedding** in Go:
```go
//go:embed web/dist/*
var staticFiles embed.FS
```

---

#### 3.3 User & Project Management (Week 13-14)

**Location**: `internal/auth/`, `internal/config/`

**Features**:
- [ ] IAM user creation (automated)
- [ ] Policy generation (least privilege)
- [ ] Group management
- [ ] Project-level access control
- [ ] Globus Auth integration (optional)

**Database**: SQLite (local) or DynamoDB (cloud)

---

### Milestone 3: Accessible Interface (End of Week 14)

**Deliverables**:
- Working web UI
- File upload/download
- User management
- Project creation

**Demo**:
Open `http://localhost:7878` and show:
- Dashboard with storage stats
- Browse files
- Upload a file with metadata form
- Create a project and add users

---

## Phase 4: Compute & Workflows (Weeks 15-18)

### Priorities
1. **CRITICAL**: AWS Batch integration
2. **HIGH**: Snakemake support
3. **HIGH**: Workflow monitoring
4. **MEDIUM**: Nextflow support
5. **MEDIUM**: Environment capture

### Technical Specifications

#### 4.1 AWS Batch Integration (Week 15-16)

**Location**: `internal/workflow/batch.go`

**Key Components**:
```go
type BatchExecutor struct {
    batch       *batch.Client
    ec2         *ec2.Client
    ecs         *ecs.Client
    environment *BatchEnvironment
}

type BatchEnvironment struct {
    ComputeEnvironment string
    JobQueue           string
    JobDefinition      string
}
```

**Setup** (CloudFormation):
- Compute environment (spot instances)
- Job queue
- ECR repository for containers
- IAM roles

**Features**:
- [ ] Spot instance configuration
- [ ] Auto-scaling (0 to max vCPUs)
- [ ] Job submission
- [ ] Status monitoring
- [ ] Log streaming
- [ ] Cost tracking per job

---

#### 4.2 Workflow Engines (Week 16-17)

**Snakemake** (Priority):
```go
// internal/workflow/snakemake.go
type SnakemakeExecutor struct {
    batch     *BatchExecutor
    container string  // Docker image with Snakemake
}

func (e *SnakemakeExecutor) Run(snakefile, config string) (*WorkflowRun, error)
```

**Nextflow** (Secondary):
```go
// internal/workflow/nextflow.go
type NextflowExecutor struct {
    batch *BatchExecutor
}
```

**Features**:
- [ ] Parse workflow definition
- [ ] Build Docker container
- [ ] Submit to Batch
- [ ] Monitor progress
- [ ] Capture outputs

---

#### 4.3 Environment Capture (Week 17-18)

**Features**:
- [ ] Capture Docker image SHA
- [ ] Export conda/pip environment
- [ ] Save workflow definition
- [ ] Record input/output locations
- [ ] Generate provenance graph

---

### Milestone 4: Computational Workflows (End of Week 18)

**Deliverables**:
- Submit Snakemake workflows
- Monitor progress
- View results

**Demo**:
```bash
cicada workflow run snakemake \
  --snakefile analysis.smk \
  --config input=s3://bucket/data/ \
  --spot
```

---

## Phase 5: Workstations & Portal (Weeks 19-22)

### Priorities
1. **HIGH**: EC2 workstation launcher
2. **HIGH**: DCV/noVNC setup
3. **MEDIUM**: Public portal
4. **MEDIUM**: DOI minting

### Technical Specifications

#### 5.1 Workstation Launcher (Week 19-20)

**Location**: `internal/workstation/`

**Key Components**:
```go
type WorkstationManager struct {
    ec2           *ec2.Client
    cfn           *cloudformation.Client
    sessionStore  SessionStore
}

type Session struct {
    ID           string
    InstanceID   string
    InstanceType string
    PublicIP     string
    DCVURL       string
    State        string
    LaunchedAt   time.Time
    IdleSince    time.Time
}
```

**Features**:
- [ ] Launch GPU instances (g4dn family)
- [ ] Install DCV or noVNC
- [ ] Mount S3 via s3fs
- [ ] Auto-shutdown on idle
- [ ] Snapshot capability
- [ ] Resume from snapshot

**Images**:
- Pre-built AMIs with:
  - Ubuntu 22.04
  - DCV server
  - Common scientific tools
  - S3 mounting tools

---

#### 5.2 Public Portal (Week 20-21)

**Location**: `internal/portal/`

**Features**:
- [ ] Dataset landing pages
- [ ] Search interface
- [ ] Download links (signed URLs)
- [ ] Usage analytics
- [ ] OAI-PMH harvesting endpoint

**Static Site Generation**: Option to generate static HTML for CloudFront

---

#### 5.3 DOI Integration (Week 21-22)

**Location**: `internal/doi/datacite.go`

See artifact: `/tmp/cicada-spec/internal/doi/datacite.go`

**Features**:
- [ ] DataCite API integration
- [ ] Metadata XML generation
- [ ] DOI reservation
- [ ] DOI registration
- [ ] Landing page creation
- [ ] Citation formatting

---

### Milestone 5: Publication Ready (End of Week 22)

**Deliverables**:
- Launch remote workstations
- Public data portal
- Mint DOIs

**Demo**:
- Launch workstation with Napari
- Publish dataset with DOI
- Show public landing page

---

## Phase 6: Compliance & Polish (Weeks 23-26)

### Priorities
1. **HIGH**: NIST 800-171 mode
2. **HIGH**: Audit logging
3. **MEDIUM**: DLP scanning
4. **MEDIUM**: Documentation
5. **LOW**: Community schemas

### Technical Specifications

#### 6.1 Compliance (Week 23-24)

**Location**: `internal/compliance/`

**Features**:
- [ ] Enhanced audit logging
- [ ] Encryption verification
- [ ] Access control reports
- [ ] Compliance dashboard
- [ ] Automated compliance checks

---

#### 6.2 Documentation (Week 24-25)

**Content**:
- Getting Started Guide
- Domain-specific tutorials
- API documentation
- Video tutorials
- Best practices

---

#### 6.3 Testing & QA (Week 25-26)

**Test Coverage Goals**:
- Unit tests: >80%
- Integration tests: Key workflows
- End-to-end tests: Complete user journeys
- Performance tests: Large datasets

**Load Testing**:
- 10,000 files synced
- 100 concurrent uploads
- 1000 metadata searches

---

### Milestone 6: Production Ready (End of Week 26)

**Deliverables**:
- Compliance mode working
- Complete documentation
- Test coverage >80%
- Performance benchmarks met

**Launch Checklist**:
- [ ] All critical features complete
- [ ] Security audit passed
- [ ] Performance benchmarks met
- [ ] Documentation complete
- [ ] Example datasets available
- [ ] Community schemas (at least 10 domains)
- [ ] Website live
- [ ] Blog post / announcement

---

## Development Guidelines

### Code Style

**Go**:
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt`, `golint`, `go vet`
- Comments for exported functions

**Testing**:
- Table-driven tests
- Use `testify` for assertions
- Mock external services

### Git Workflow

**Branches**:
- `main`: Stable, production-ready
- `develop`: Integration branch
- `feature/*`: Feature branches
- `bugfix/*`: Bug fixes

**Commits**:
- Conventional commits format
- Example: `feat(sync): add bandwidth throttling`

**Pull Requests**:
- Require review
- CI must pass
- Test coverage must not decrease

### CI/CD Pipeline

**GitHub Actions** (or similar):

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: go build ./cmd/cicada
  
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3
```

### Release Process

1. Create release branch: `release/v1.0.0`
2. Update CHANGELOG.md
3. Tag: `git tag v1.0.0`
4. Build binaries for all platforms
5. Create GitHub release
6. Update website

---

## Success Metrics

### Technical Metrics

- **Sync Performance**: 10GB in <15 min on 100Mbps
- **Metadata Extraction**: <5 seconds per file
- **Search Latency**: <1 second for 100k records
- **Uptime**: >99.9% for daemon
- **Test Coverage**: >80%

### User Metrics

- **Time to First Sync**: <10 minutes (setup to first data in S3)
- **Learning Curve**: New user productive in <1 hour
- **Error Rate**: <1% failed uploads
- **Support Tickets**: <5% of users need help

### Adoption Metrics (6 months post-launch)

- **Users**: 100+ labs
- **Data Stored**: 1+ PB across all users
- **Active Installations**: 500+
- **Community Schemas**: 20+ domains
- **GitHub Stars**: 1000+

---

## Risk Mitigation

### Technical Risks

**Risk**: AWS API rate limits  
**Mitigation**: Implement exponential backoff, caching

**Risk**: Large file handling (>5GB)  
**Mitigation**: Multipart upload, resume capability

**Risk**: Metadata extraction failures  
**Mitigation**: Graceful fallback to generic extractor

### User Adoption Risks

**Risk**: Too complex for non-technical users  
**Mitigation**: Extensive user testing, simplified defaults

**Risk**: Cost concerns  
**Mitigation**: Clear cost estimation, budget alerts, optimization tips

**Risk**: Data loss fears  
**Mitigation**: Emphasize versioning, show recovery examples

---

## Post-Launch Roadmap

### Version 1.1 (Month 7-8)

- [ ] Additional workflow engines (CWL, WDL)
- [ ] More instrument adapters
- [ ] Enhanced cost optimization
- [ ] Mobile app (view-only)

### Version 1.2 (Month 9-10)

- [ ] Multi-cloud support (Azure, GCP)
- [ ] Advanced provenance tracking
- [ ] Machine learning metadata prediction
- [ ] Collaborative annotations

### Version 2.0 (Month 11-12)

- [ ] Federation between labs
- [ ] Real-time collaboration
- [ ] Advanced data governance
- [ ] Institutional deployment tools

---

## Resources & References

### Learning Resources

**Go**:
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)

**AWS**:
- [AWS SDK for Go V2](https://aws.github.io/aws-sdk-go-v2/docs/)
- [S3 Best Practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/optimizing-performance.html)

**Standards**:
- [DataCite Schema](https://schema.datacite.org/)
- [FAIR Principles](https://www.go-fair.org/fair-principles/)
- [NIST 800-171](https://csrc.nist.gov/publications/detail/sp/800-171/rev-2/final)

### Similar Projects (for inspiration)

- [Globus](https://www.globus.org/) - Research data management
- [iRODS](https://irods.org/) - Data management middleware
- [Dataverse](https://dataverse.org/) - Research data repository
- [rclone](https://rclone.org/) - Rsync for cloud storage
- [Nextcloud](https://nextcloud.com/) - Self-hosted cloud

---

## Conclusion

This roadmap provides a clear path from initial development to production-ready system. The phased approach allows for:

1. **Early validation**: Basic sync working in 6 weeks
2. **Incremental features**: Each phase adds value
3. **User feedback**: Can ship MVP at week 14
4. **Risk management**: Core features first, polish later

You can tackle each phase systematically, using these specifications as a guide for implementation.
