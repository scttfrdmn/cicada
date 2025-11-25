# Cicada v0.4.0 Roadmap

**Version**: 0.4.0
**Target Release**: Q1 2026
**Duration**: 8-10 weeks
**Focus**: Provider Integration, Storage Expansion, Advanced Features

---

## Overview

v0.4.0 transforms Cicada from a data management platform into a complete data publication platform with live DOI provider integration, multi-cloud storage, and advanced metadata capabilities.

### Goals

1. **Live Provider Integration**: Connect to DataCite and Zenodo for real DOI minting
2. **Multi-Cloud Storage**: Expand beyond S3 to Azure Blob and Google Cloud Storage
3. **Advanced Metadata**: Custom extractors, metadata search, and enhanced validation
4. **Production Readiness**: Comprehensive error handling, retry logic, and monitoring

### Success Criteria

- ✅ Mint real DOIs via DataCite and Zenodo
- ✅ Sync data to Azure and GCS in addition to S3
- ✅ Users can create custom metadata extractors
- ✅ Enhanced metadata search and filtering
- ✅ 95%+ test coverage on new code
- ✅ Complete provider integration documentation

---

## Release Themes

### Theme 1: DOI Provider Integration (Weeks 1-4)

**Goal**: Enable real DOI minting and management through DataCite and Zenodo APIs

**Issues**:
- #26: Provider Configuration System
- #27: DataCite API Client - Core Infrastructure
- #28: Zenodo API Client - Core Infrastructure
- #29: Provider Registry Enhancement for Live APIs
- #30: DataCite Metadata Mapping to Schema v4.5
- #31: Zenodo Metadata Mapping
- #32: CLI: cicada doi publish Command
- #33: CLI: cicada doi status Command
- #34: CLI: cicada doi list Command
- #35: Comprehensive Error Handling and Retry Logic
- #37: Provider Integration Documentation

**Deliverables**:
- Working DataCite integration (sandbox + production)
- Working Zenodo integration (sandbox + production)
- CLI commands for DOI lifecycle management
- Provider configuration and credential management
- Comprehensive error handling and retry logic
- Complete integration documentation

### Theme 2: Multi-Cloud Storage (Weeks 5-6)

**Goal**: Expand storage backend support to Azure and GCS

**New Issues**:
- Azure Blob Storage backend implementation
- Google Cloud Storage backend implementation
- Unified backend configuration system
- Cross-cloud sync capabilities
- Multi-cloud CLI commands
- Storage backend documentation

**Deliverables**:
- Azure Blob Storage backend
- Google Cloud Storage backend
- Unified configuration for all storage backends
- Documentation for each storage backend
- Migration guides from S3 to other backends

### Theme 3: Advanced Metadata Features (Weeks 7-8)

**Goal**: Enhance metadata capabilities with custom extractors and search

**New Issues**:
- Custom metadata extractor plugin system
- Metadata search and filtering engine
- Enhanced metadata validation
- Metadata export formats (CSV, JSON, Parquet)
- Metadata statistics and analytics
- Custom extractor documentation

**Deliverables**:
- Plugin system for custom extractors (Go and Python)
- Metadata search CLI commands
- Enhanced validation with custom rules
- Multiple export formats
- Analytics and reporting features
- Developer guide for custom extractors

### Theme 4: Production Hardening (Weeks 9-10)

**Goal**: Ensure production readiness with monitoring, logging, and reliability

**New Issues**:
- Structured logging system
- Metrics and monitoring (Prometheus)
- Health checks and status endpoints
- Rate limiting and throttling
- Configuration validation
- Production deployment guide

**Deliverables**:
- Structured logging with levels
- Prometheus metrics endpoint
- Health check system
- Rate limiting for API calls
- Configuration validation tools
- Production deployment documentation

---

## Detailed Timeline

### Phase 1: Provider Foundation (Weeks 1-2)

**Week 1: Core Infrastructure**
- [ ] Issue #26: Provider Configuration System (2 days)
- [ ] Issue #27: DataCite API Client - Core Infrastructure (3 days)
- [ ] Issue #29: Provider Registry Enhancement (2 days)

**Week 2: Zenodo + Metadata Mapping**
- [ ] Issue #28: Zenodo API Client - Core Infrastructure (3 days)
- [ ] Issue #30: DataCite Metadata Mapping (2 days)
- [ ] Issue #31: Zenodo Metadata Mapping (2 days)

**Milestones**:
- ✅ Both API clients functional with sandbox
- ✅ Metadata mapping complete for both providers
- ✅ Unit tests passing

### Phase 2: Provider Integration (Weeks 3-4)

**Week 3: CLI Commands**
- [ ] Issue #32: CLI: cicada doi publish Command (2 days)
- [ ] Issue #33: CLI: cicada doi status Command (1 day)
- [ ] Issue #34: CLI: cicada doi list Command (1 day)
- [ ] Issue #35: Error Handling and Retry Logic (3 days)

**Week 4: Testing + Documentation**
- [ ] Integration tests with sandbox APIs (2 days)
- [ ] Issue #37: Provider Integration Documentation (3 days)
- [ ] End-to-end testing workflows (2 days)

**Milestones**:
- ✅ Can mint real DOIs in sandbox
- ✅ Full CLI for DOI lifecycle
- ✅ Documentation complete
- ✅ Ready for production testing

### Phase 3: Multi-Cloud Storage (Weeks 5-6)

**Week 5: Azure Blob Storage**
- [ ] Create Azure backend interface implementation (2 days)
- [ ] Azure authentication and configuration (1 day)
- [ ] Azure-specific operations (list, upload, download) (2 days)
- [ ] Azure integration tests (1 day)
- [ ] Azure documentation (1 day)

**Week 6: Google Cloud Storage**
- [ ] Create GCS backend interface implementation (2 days)
- [ ] GCS authentication and configuration (1 day)
- [ ] GCS-specific operations (list, upload, download) (2 days)
- [ ] GCS integration tests (1 day)
- [ ] GCS documentation (1 day)
- [ ] Cross-cloud sync examples (1 day)

**Milestones**:
- ✅ Azure backend functional
- ✅ GCS backend functional
- ✅ Unified backend configuration
- ✅ Documentation for all backends

### Phase 4: Advanced Metadata (Weeks 7-8)

**Week 7: Custom Extractors**
- [ ] Design plugin system architecture (1 day)
- [ ] Implement Go plugin interface (2 days)
- [ ] Implement Python plugin support (2 days)
- [ ] Plugin discovery and loading (1 day)
- [ ] Example custom extractors (1 day)

**Week 8: Metadata Search + Export**
- [ ] Design metadata search query language (1 day)
- [ ] Implement search engine (2 days)
- [ ] Add export formats (CSV, JSON, Parquet) (2 days)
- [ ] Metadata analytics and statistics (1 day)
- [ ] Documentation and examples (1 day)

**Milestones**:
- ✅ Custom extractor plugin system
- ✅ Metadata search functional
- ✅ Multiple export formats
- ✅ Developer documentation

### Phase 5: Production Hardening (Weeks 9-10)

**Week 9: Observability**
- [ ] Structured logging implementation (2 days)
- [ ] Prometheus metrics (2 days)
- [ ] Health check endpoints (1 day)
- [ ] Monitoring dashboard (1 day)
- [ ] Documentation (1 day)

**Week 10: Polish + Release**
- [ ] Rate limiting and throttling (2 days)
- [ ] Configuration validation (1 day)
- [ ] Production deployment guide (1 day)
- [ ] Release preparation (CHANGELOG, version bumps) (1 day)
- [ ] Final testing and bug fixes (2 days)

**Milestones**:
- ✅ Production monitoring in place
- ✅ Rate limiting functional
- ✅ Deployment guide complete
- ✅ v0.4.0 ready for release

---

## Technical Architecture

### Provider Integration Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI Commands                           │
│  cicada doi publish | status | list | update | delete       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                 Provider Registry                           │
│  - GetProvider(name) → Provider                             │
│  - ListProviders() → []Provider                             │
│  - Register/Unregister                                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
           ┌───────────┴───────────┐
           │                       │
┌──────────▼─────────┐  ┌─────────▼──────────┐
│  DataCite Client   │  │   Zenodo Client    │
│  - Auth: Basic     │  │   - Auth: Token    │
│  - CreateDOI()     │  │   - CreateDOI()    │
│  - UpdateDOI()     │  │   - UpdateDOI()    │
│  - PublishDOI()    │  │   - PublishDOI()   │
│  - GetDOI()        │  │   - GetDOI()       │
│  - ListDOIs()      │  │   - ListDOIs()     │
└────────────────────┘  └────────────────────┘
```

### Storage Backend Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Backend Interface                        │
│  List, Read, Write, Delete, Stat, Close                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
       ┌───────────────┼───────────────┐
       │               │               │
┌──────▼───────┐ ┌────▼──────┐ ┌──────▼─────────┐
│ LocalBackend │ │ S3Backend │ │ AzureBackend   │
│              │ │           │ │                │
│ - Filesystem │ │ - AWS SDK │ │ - Azure SDK    │
│              │ │           │ │                │
└──────────────┘ └───────────┘ └────────────────┘
                       │
              ┌────────▼─────────┐
              │   GCSBackend     │
              │                  │
              │ - Google SDK     │
              │                  │
              └──────────────────┘
```

### Custom Extractor Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Extractor Registry                         │
│  - Register(extractor)                                      │
│  - GetExtractor(filename) → Extractor                       │
│  - ListExtractors() → []Extractor                          │
└──────────────────────┬──────────────────────────────────────┘
                       │
           ┌───────────┼────────────┐
           │           │            │
┌──────────▼───┐ ┌────▼─────┐ ┌────▼─────────────┐
│ Built-in     │ │ Go Plugin│ │ Python Plugin    │
│ Extractors   │ │          │ │                  │
│ (14 types)   │ │ .so file │ │ .py file         │
└──────────────┘ └──────────┘ └──────────────────┘
```

---

## New Issues to Create

### Multi-Cloud Storage Issues

1. **Azure Blob Storage Backend** (#TBD)
   - Priority: High
   - Estimate: 5 days
   - Labels: `type: feature`, `area: storage`, `priority: high`

2. **Google Cloud Storage Backend** (#TBD)
   - Priority: High
   - Estimate: 5 days
   - Labels: `type: feature`, `area: storage`, `priority: high`

3. **Unified Backend Configuration** (#TBD)
   - Priority: Medium
   - Estimate: 2 days
   - Labels: `type: feature`, `area: config`, `priority: medium`

4. **Cross-Cloud Sync Documentation** (#TBD)
   - Priority: Medium
   - Estimate: 2 days
   - Labels: `type: documentation`, `area: storage`

### Advanced Metadata Issues

5. **Custom Extractor Plugin System** (#TBD)
   - Priority: High
   - Estimate: 5 days
   - Labels: `type: feature`, `area: metadata`, `priority: high`

6. **Python Plugin Support for Extractors** (#TBD)
   - Priority: High
   - Estimate: 3 days
   - Labels: `type: feature`, `area: metadata`, `priority: high`

7. **Metadata Search Engine** (#TBD)
   - Priority: Medium
   - Estimate: 4 days
   - Labels: `type: feature`, `area: metadata`, `priority: medium`

8. **Enhanced Metadata Export Formats** (#TBD)
   - Priority: Medium
   - Estimate: 3 days
   - Labels: `type: feature`, `area: metadata`, `priority: medium`

9. **Metadata Analytics and Statistics** (#TBD)
   - Priority: Low
   - Estimate: 2 days
   - Labels: `type: feature`, `area: metadata`, `priority: low`

10. **Custom Extractor Developer Guide** (#TBD)
    - Priority: High
    - Estimate: 3 days
    - Labels: `type: documentation`, `area: metadata`

### Production Hardening Issues

11. **Structured Logging System** (#TBD)
    - Priority: High
    - Estimate: 2 days
    - Labels: `type: feature`, `area: infrastructure`, `priority: high`

12. **Prometheus Metrics Integration** (#TBD)
    - Priority: Medium
    - Estimate: 3 days
    - Labels: `type: feature`, `area: observability`, `priority: medium`

13. **Health Check System** (#TBD)
    - Priority: Medium
    - Estimate: 2 days
    - Labels: `type: feature`, `area: infrastructure`, `priority: medium`

14. **Rate Limiting and Throttling** (#TBD)
    - Priority: High
    - Estimate: 2 days
    - Labels: `type: feature`, `area: infrastructure`, `priority: high`

15. **Configuration Validation System** (#TBD)
    - Priority: Medium
    - Estimate: 2 days
    - Labels: `type: feature`, `area: config`, `priority: medium`

16. **Production Deployment Guide** (#TBD)
    - Priority: High
    - Estimate: 3 days
    - Labels: `type: documentation`, `area: deployment`

---

## Dependencies

### External Dependencies

**Azure SDK for Go**:
```go
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob
github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

**Google Cloud SDK for Go**:
```go
cloud.google.com/go/storage
google.golang.org/api/option
```

**Prometheus Client**:
```go
github.com/prometheus/client_golang/prometheus
github.com/prometheus/client_golang/prometheus/promhttp
```

**Structured Logging**:
```go
go.uber.org/zap
// or
github.com/rs/zerolog
```

**Go Plugin System**:
```go
plugin // stdlib
```

**Python Integration**:
```go
github.com/go-python/gpython
// or exec + JSON communication
```

### Internal Dependencies

- Provider configuration system (v0.2.0) ✅
- Metadata extraction system (v0.2.0) ✅
- Sync engine (v0.1.0) ✅
- Backend interface (v0.1.0) ✅

---

## Testing Strategy

### Unit Testing
- Target: 90%+ coverage on new code
- Mock external APIs (DataCite, Zenodo, Azure, GCS)
- Table-driven tests for all providers and backends

### Integration Testing
- Sandbox API testing for DataCite and Zenodo
- Real cloud storage testing (with cleanup)
- End-to-end workflows

### Performance Testing
- Benchmark multi-cloud sync operations
- Benchmark custom extractor loading
- Benchmark metadata search queries

### Manual Testing
- Production DOI minting (with real credentials)
- Cross-cloud sync workflows
- Custom extractor development workflow

---

## Documentation Plan

### New Documentation Files

1. **PROVIDER_INTEGRATION.md**
   - DataCite setup and usage
   - Zenodo setup and usage
   - Provider configuration
   - Troubleshooting

2. **MULTI_CLOUD_STORAGE.md**
   - Azure Blob Storage setup
   - Google Cloud Storage setup
   - Cross-cloud sync patterns
   - Cost comparison

3. **CUSTOM_EXTRACTORS.md**
   - Plugin system architecture
   - Go plugin development
   - Python plugin development
   - Example extractors

4. **PRODUCTION_DEPLOYMENT.md**
   - Deployment options (Docker, systemd, K8s)
   - Monitoring setup
   - Security best practices
   - Scaling considerations

### Documentation Updates

- Update GETTING_STARTED.md with Azure/GCS examples
- Update WORKFLOWS.md with DOI publication workflows
- Update TROUBLESHOOTING.md with provider-specific issues
- Update ADVANCED.md with custom extractors and monitoring

---

## Risk Assessment

### High Risk

1. **Provider API Changes**
   - Risk: DataCite/Zenodo APIs change during development
   - Mitigation: Version pinning, comprehensive error handling

2. **Multi-Cloud Complexity**
   - Risk: Each cloud has unique quirks and auth patterns
   - Mitigation: Thorough testing, clear documentation

3. **Plugin Security**
   - Risk: Custom plugins could introduce vulnerabilities
   - Mitigation: Sandboxing, validation, security guidelines

### Medium Risk

1. **Performance Impact**
   - Risk: Plugin system adds overhead
   - Mitigation: Benchmark and optimize, lazy loading

2. **Configuration Complexity**
   - Risk: Too many config options overwhelm users
   - Mitigation: Sensible defaults, validation, examples

### Low Risk

1. **Backwards Compatibility**
   - Risk: New features break existing workflows
   - Mitigation: Extensive testing, deprecation policy

---

## Success Metrics

### Functional Metrics
- [ ] 100% of provider integration issues closed
- [ ] All cloud backends functional
- [ ] Custom extractor examples working
- [ ] 95%+ test coverage

### Quality Metrics
- [ ] Zero P0/P1 bugs at release
- [ ] All documentation complete and reviewed
- [ ] Integration tests passing
- [ ] Performance benchmarks meet targets

### User Metrics
- [ ] Can mint DOI in < 5 minutes (first time)
- [ ] Can add custom extractor in < 30 minutes
- [ ] Can set up Azure/GCS in < 15 minutes

---

## Post-Release (v0.5.0 Preview)

### Potential Features for v0.5.0

1. **Workflow Execution**
   - Snakemake/Nextflow integration
   - Workflow templates
   - Compute-to-data support

2. **Collaboration Features**
   - Shared workspaces
   - Access control
   - Audit logs

3. **Advanced Search**
   - Full-text search in metadata
   - Semantic search
   - Saved searches

4. **UI/Web Interface**
   - Web dashboard for management
   - Visual data exploration
   - Metadata editing UI

---

## Appendix

### Issue Quick Reference

| Issue | Title | Priority | Phase |
|-------|-------|----------|-------|
| #26 | Provider Configuration System | Critical | 1 |
| #27 | DataCite API Client | Critical | 1 |
| #28 | Zenodo API Client | Critical | 1 |
| #29 | Provider Registry Enhancement | High | 1 |
| #30 | DataCite Metadata Mapping | High | 1 |
| #31 | Zenodo Metadata Mapping | High | 1 |
| #32 | CLI: doi publish | High | 2 |
| #33 | CLI: doi status | Medium | 2 |
| #34 | CLI: doi list | Medium | 2 |
| #35 | Error Handling and Retry | Critical | 2 |
| #37 | Provider Documentation | High | 2 |
| TBD | Azure Backend | High | 3 |
| TBD | GCS Backend | High | 3 |
| TBD | Custom Extractor System | High | 4 |
| TBD | Metadata Search | Medium | 4 |
| TBD | Structured Logging | High | 5 |
| TBD | Prometheus Metrics | Medium | 5 |

### Resource Estimates

**Total Effort**: ~50-60 developer days (10 weeks)

**By Theme**:
- Provider Integration: 20 days
- Multi-Cloud Storage: 12 days
- Advanced Metadata: 15 days
- Production Hardening: 10 days
- Documentation: 8 days (integrated throughout)

**By Phase**:
- Weeks 1-2: 10 days
- Weeks 3-4: 10 days
- Weeks 5-6: 10 days
- Weeks 7-8: 12 days
- Weeks 9-10: 10 days

---

**Last Updated**: 2025-11-25
**Status**: Planning
**Next Action**: Create new issues for multi-cloud, metadata, and hardening themes
