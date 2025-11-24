# Cicada v0.3.0 Roadmap

**Target Release:** Q2 2025 (April-June 2025)
**Status:** Planning Phase
**Previous Release:** v0.2.0 (January 2025)

## Overview

Cicada v0.3.0 focuses on **production-ready DOI registration** and **expanded file format support**. This release transforms v0.2.0's metadata extraction and DOI preparation capabilities from preparation tools into complete end-to-end workflows that integrate with live repository services.

### Key Themes

1. **Live Provider Integration** - Connect to real DataCite/Zenodo/Dryad APIs
2. **Format Expansion** - Support microscopy and bioinformatics file formats
3. **User Customization** - Custom presets and interactive workflows
4. **Production Readiness** - Caching, error recovery, monitoring

## Goals

### Primary Goals (Must Have)

1. **Complete DOI Workflow** - Users can mint real DOIs from command line
2. **Microscopy Support** - CZI and OME-TIFF metadata extraction
3. **Provider Integration** - At least DataCite and Zenodo working in sandbox and production

### Secondary Goals (Should Have)

4. **Custom Presets** - Users can define lab-specific validation rules
5. **Interactive Editing** - TUI for metadata review and enrichment
6. **Bioinformatics Formats** - BAM/SAM and VCF support

### Stretch Goals (Nice to Have)

7. **Additional Providers** - Dryad and Figshare integration
8. **API Server Mode** - REST API for programmatic access
9. **Metadata Caching** - Performance optimization for large archives

## Milestones

### Milestone 1: Provider Integration Foundation (3-4 weeks)

**Goal:** Implement live DataCite and Zenodo API integrations

**Features:**
- DataCite API client (sandbox and production)
- Zenodo API client (sandbox and production)
- Authentication and credential management
- Error handling and retry logic
- API rate limiting and backoff
- Configuration management for providers

**Deliverables:**
- `internal/doi/datacite_client.go` - Complete DataCite API implementation
- `internal/doi/zenodo_client.go` - Complete Zenodo API implementation
- `internal/doi/provider_config.go` - Provider configuration
- Provider integration tests with sandbox APIs
- Provider setup documentation

**Success Criteria:**
- [ ] Can create draft DOI in DataCite sandbox
- [ ] Can publish DOI in DataCite sandbox
- [ ] Can upload files and create DOI in Zenodo sandbox
- [ ] Can update existing DOIs
- [ ] Error handling works for common failure scenarios
- [ ] Configuration management is user-friendly

**Complexity:** HIGH (API integration, authentication, error handling)

### Milestone 2: Microscopy File Formats (2-3 weeks)

**Goal:** Extract metadata from CZI and OME-TIFF files

**Features:**
- CZI metadata extraction (dimensions, channels, objectives, timestamps)
- OME-TIFF metadata extraction (OME-XML parsing)
- Integration with Zeiss presets
- Multi-dimensional image handling (Z-stacks, time series, channels)
- Performance optimization for large images

**Deliverables:**
- `internal/metadata/czi_extractor.go` - Complete CZI implementation
- `internal/metadata/ome_tiff_extractor.go` - Complete OME-TIFF implementation
- Integration tests with real microscopy files
- Updated preset validation for microscopy metadata
- Microscopy format documentation

**Success Criteria:**
- [ ] Can extract metadata from Zeiss CZI files
- [ ] Can extract metadata from OME-TIFF files
- [ ] Preset validation works with extracted microscopy metadata
- [ ] Performance: < 5 seconds for typical microscopy files (< 1 GB)
- [ ] Handles multi-channel, Z-stack, time-series images

**Complexity:** MEDIUM-HIGH (Binary format parsing, XML parsing, domain knowledge)

### Milestone 3: Custom Preset System (2 weeks)

**Goal:** Enable users to create and manage custom presets

**Features:**
- Preset creation from templates or scratch
- YAML-based preset definition format
- Preset validation rules (field types, ranges, patterns)
- Preset import/export
- Preset inheritance (extend existing presets)
- Custom quality scoring weights

**Deliverables:**
- `internal/metadata/preset_builder.go` - Preset creation and validation
- `cicada metadata preset create` - CLI command
- `cicada metadata preset import/export` - CLI commands
- Preset definition schema and examples
- Custom preset documentation

**Success Criteria:**
- [ ] Users can create presets from YAML files
- [ ] Custom presets validate correctly
- [ ] Can inherit from built-in presets
- [ ] Import/export works for sharing presets
- [ ] Quality scoring reflects custom field weights

**Complexity:** MEDIUM (Schema design, validation logic, CLI integration)

### Milestone 4: Interactive Metadata Editing (2-3 weeks)

**Goal:** Provide interactive UI for metadata review and enrichment

**Features:**
- Terminal UI (TUI) for metadata editing
- Form-based enrichment with validation
- Side-by-side comparison (extracted vs enriched)
- Quality score real-time updates
- Field help and examples
- Validation feedback

**Deliverables:**
- `internal/ui/metadata_editor.go` - TUI implementation
- `cicada metadata edit` - CLI command
- Integration with DOI workflow
- Interactive editing documentation

**Success Criteria:**
- [ ] TUI works on Linux, macOS, Windows
- [ ] Can edit all metadata fields interactively
- [ ] Validation provides helpful feedback
- [ ] Quality score updates in real-time
- [ ] Can save and resume editing sessions

**Complexity:** MEDIUM-HIGH (TUI library integration, UX design, validation)

### Milestone 5: Bioinformatics Formats (2 weeks)

**Goal:** Support BAM/SAM alignment and VCF variant files

**Features:**
- BAM/SAM metadata extraction (alignment statistics, references)
- VCF metadata extraction (sample info, variant stats)
- Integration with generic-sequencing preset
- Performance optimization (BAM indexing)

**Deliverables:**
- `internal/metadata/bam_extractor.go` - BAM/SAM implementation
- `internal/metadata/vcf_extractor.go` - VCF implementation
- Integration tests with real bioinformatics files
- Bioinformatics format documentation

**Success Criteria:**
- [ ] Can extract metadata from BAM/SAM files
- [ ] Can extract metadata from VCF files
- [ ] Performance: < 1 second for typical files
- [ ] Handles compressed formats (.bam.gz, .vcf.gz)

**Complexity:** MEDIUM (Binary format parsing, domain knowledge)

### Milestone 6: Production Readiness (2 weeks)

**Goal:** Make v0.3.0 production-ready with caching, monitoring, and error recovery

**Features:**
- Metadata caching (avoid re-extraction)
- API response caching
- Retry logic with exponential backoff
- Detailed logging and monitoring
- Error recovery and rollback
- Batch operation progress tracking

**Deliverables:**
- `internal/cache/metadata_cache.go` - Caching implementation
- Enhanced logging throughout codebase
- Batch operation commands
- Production deployment documentation

**Success Criteria:**
- [ ] Metadata cache works correctly
- [ ] Cached metadata significantly improves performance
- [ ] Retry logic prevents transient failures
- [ ] Logging provides useful debugging information
- [ ] Batch operations track progress

**Complexity:** MEDIUM (Caching strategy, error handling, logging)

### Milestone 7: Additional Providers (Stretch, 2-3 weeks)

**Goal:** Add Dryad and Figshare provider support

**Features:**
- Dryad API client
- Figshare API client
- Provider-specific metadata mapping
- Multi-provider workflow

**Deliverables:**
- `internal/doi/dryad_client.go` - Dryad implementation
- `internal/doi/figshare_client.go` - Figshare implementation
- Provider comparison documentation

**Success Criteria:**
- [ ] Can create and publish DOIs on Dryad
- [ ] Can create and publish DOIs on Figshare
- [ ] Provider selection is seamless

**Complexity:** MEDIUM (Similar to DataCite/Zenodo but different APIs)

### Milestone 8: API Server Mode (Stretch, 2-3 weeks)

**Goal:** Provide REST API for programmatic access

**Features:**
- REST API server
- OpenAPI/Swagger documentation
- Authentication and authorization
- Rate limiting
- API client libraries (Go, Python)

**Deliverables:**
- `internal/api/server.go` - API server
- `cicada serve` - CLI command
- OpenAPI specification
- API documentation and examples

**Success Criteria:**
- [ ] REST API covers all CLI functionality
- [ ] API documentation is comprehensive
- [ ] Authentication works correctly
- [ ] Performance meets requirements

**Complexity:** HIGH (API design, authentication, documentation)

## Timeline

### Conservative Estimate (5-6 months)

| Milestone | Duration | Target Completion |
|-----------|----------|-------------------|
| 1. Provider Integration | 4 weeks | End of February 2025 |
| 2. Microscopy Formats | 3 weeks | Mid-March 2025 |
| 3. Custom Presets | 2 weeks | End of March 2025 |
| 4. Interactive Editing | 3 weeks | Mid-April 2025 |
| 5. Bioinformatics Formats | 2 weeks | End of April 2025 |
| 6. Production Readiness | 2 weeks | Mid-May 2025 |
| **v0.3.0 Release** | - | **End of May 2025** |

### Aggressive Estimate (3-4 months)

Focus on core features only (Milestones 1-3, 6):

| Milestone | Duration | Target Completion |
|-----------|----------|-------------------|
| 1. Provider Integration | 3 weeks | Mid-February 2025 |
| 2. Microscopy Formats | 2 weeks | End of February 2025 |
| 3. Custom Presets | 1.5 weeks | Mid-March 2025 |
| 6. Production Readiness | 1.5 weeks | End of March 2025 |
| **v0.3.0 Release** | - | **Early April 2025** |

*Defer Milestones 4, 5, 7, 8 to v0.3.1 or v0.4.0*

## Dependencies

### External Dependencies

1. **Provider Accounts/Credentials:**
   - DataCite sandbox account (free)
   - DataCite production account (requires institutional membership)
   - Zenodo sandbox account (free)
   - Zenodo production account (free)
   - Dryad account (paid per dataset)
   - Figshare account (free tier available)

2. **Test Data:**
   - Real CZI files from Zeiss microscopes
   - Real OME-TIFF files
   - Real BAM/SAM files
   - Real VCF files

3. **Libraries:**
   - CZI parsing library (may need to implement)
   - OME-XML parsing library (Go XML)
   - BAM/SAM parsing library (biogo/hts)
   - VCF parsing library (available)
   - TUI library (bubbletea or tview)

### Internal Dependencies

- v0.2.0 must be released and stable
- Metadata extraction framework is extensible
- Provider registry supports multiple providers
- Preset system is extensible

## Feature Priorities

### High Priority (Core v0.3.0)

1. **Provider Integration** - Without this, DOI workflow is incomplete
2. **Microscopy Formats** - Target users need this for lab data
3. **Production Readiness** - Required for reliable production use

### Medium Priority (Should Include)

4. **Custom Presets** - Users request this for lab-specific workflows
5. **Interactive Editing** - Improves UX significantly
6. **Bioinformatics Formats** - Expands user base

### Low Priority (Can Defer)

7. **Additional Providers** - Nice to have but not blocking
8. **API Server Mode** - Advanced feature for later

## Success Metrics

### Adoption Metrics

- [ ] 50+ users adopt v0.3.0 within 3 months
- [ ] 10+ labs use DOI minting in production
- [ ] 100+ DOIs minted through Cicada

### Technical Metrics

- [ ] Provider integration test coverage > 80%
- [ ] Format extraction accuracy > 95%
- [ ] API uptime > 99.5% (for provider APIs)
- [ ] Average metadata extraction time < 5 seconds (all formats)

### User Satisfaction

- [ ] User satisfaction score > 4/5
- [ ] Documentation findability > 4/5
- [ ] DOI workflow completion rate > 90%
- [ ] Feature requests < 10% of feedback (vs bugs)

## Risk Assessment

### High Risk

1. **Provider API Changes** - DataCite/Zenodo might change APIs
   - *Mitigation:* Version API clients, monitor changelog, integration tests

2. **Authentication Complexity** - Provider auth might be difficult
   - *Mitigation:* Start with sandbox, clear documentation, examples

3. **Binary Format Parsing** - CZI/BAM parsing might be complex
   - *Mitigation:* Use existing libraries where possible, limit scope

### Medium Risk

4. **Performance** - Large microscopy files might be slow
   - *Mitigation:* Benchmark early, optimize, consider sampling

5. **TUI Complexity** - Interactive UI might be complex
   - *Mitigation:* Use proven library (bubbletea), start simple

6. **Custom Preset Abuse** - Users might create invalid presets
   - *Mitigation:* Strong validation, good error messages, examples

### Low Risk

7. **Timeline Slippage** - Features might take longer
   - *Mitigation:* Conservative estimates, prioritize core features

8. **Dependency Issues** - External libraries might have bugs
   - *Mitigation:* Vet libraries early, have fallback options

## Breaking Changes

### Potential Breaking Changes

1. **Provider Configuration Format** - May need to change config structure
   - *Migration:* Automatic config migration, clear migration guide

2. **Metadata Schema Evolution** - New formats might require schema changes
   - *Migration:* Backward compatible schema, version metadata

3. **CLI Flag Changes** - May rename or reorganize flags
   - *Migration:* Deprecation warnings, maintain aliases

### Compatibility Promise

**Goal:** Maintain backward compatibility with v0.2.0 where possible.

- All v0.2.0 commands continue to work
- v0.2.0 config files remain valid
- v0.2.0 metadata format remains supported
- Clear migration paths for any breaking changes

## User Stories

### Story 1: Postdoc Publishes Sequencing Dataset

**As a** postdoc researcher
**I want to** mint a DOI for my FASTQ files and deposit them in Zenodo
**So that** I can cite the dataset in my paper

**Acceptance Criteria:**
- Can extract metadata from FASTQ files
- Can enrich with author/description
- Can upload files to Zenodo
- Can mint DOI in one command
- Receives permanent DOI URL

**v0.3.0 Solution:**
```bash
# Extract and prepare
cicada doi prepare sample_R1.fastq.gz sample_R2.fastq.gz \
  --enrich metadata.yaml \
  --provider zenodo

# Upload and mint DOI (NEW in v0.3.0)
cicada doi publish \
  --files sample_R1.fastq.gz,sample_R2.fastq.gz \
  --provider zenodo

# Output: DOI: 10.5281/zenodo.123456
# URL: https://zenodo.org/record/123456
```

### Story 2: Lab Manager Extracts Microscopy Metadata

**As a** lab manager
**I want to** extract metadata from CZI files from our Zeiss LSM 880
**So that** I can catalog our microscopy data

**Acceptance Criteria:**
- Can extract metadata from CZI files
- Validates against Zeiss LSM 880 preset
- Captures all imaging parameters
- Fast enough for batch processing

**v0.3.0 Solution:**
```bash
# Extract from CZI file (NEW in v0.3.0)
cicada metadata extract image_001.czi \
  --preset zeiss-lsm-880 \
  --format json \
  --output metadata.json

# Batch processing
find . -name "*.czi" | parallel -j 8 \
  cicada metadata extract {} --preset zeiss-lsm-880
```

### Story 3: Bioinformatician Creates Custom Preset

**As a** bioinformatician
**I want to** create a custom preset for our lab's RNA-seq pipeline
**So that** all lab members validate metadata consistently

**Acceptance Criteria:**
- Can define custom required/optional fields
- Can specify validation rules
- Can share preset with team
- Preset validates correctly

**v0.3.0 Solution:**
```bash
# Create custom preset (NEW in v0.3.0)
cicada metadata preset create lab-rnaseq \
  --template generic-sequencing \
  --add-required pipeline_version \
  --add-optional read_length_min \
  --add-optional read_length_max

# Edit preset definition
vim ~/.cicada/presets/lab-rnaseq.yaml

# Share with team
cicada metadata preset export lab-rnaseq > lab-rnaseq.yaml
# Team imports: cicada metadata preset import lab-rnaseq.yaml

# Use custom preset
cicada metadata extract sample.fastq --preset lab-rnaseq
```

### Story 4: Data Curator Uses Interactive Editor

**As a** data curator
**I want to** review and enrich metadata interactively
**So that** I can ensure high-quality metadata before DOI minting

**Acceptance Criteria:**
- Can see extracted metadata
- Can edit fields interactively
- Can see quality score in real-time
- Can validate before saving

**v0.3.0 Solution:**
```bash
# Interactive editing (NEW in v0.3.0)
cicada metadata edit sample.fastq

# Opens TUI showing:
# - Extracted metadata (read-only)
# - Editable enrichment fields
# - Real-time quality score
# - Validation feedback
# - Save/cancel options
```

### Story 5: Research Scientist Publishes to DataCite

**As a** research scientist
**I want to** mint a DOI through my institution's DataCite membership
**So that** I get a DOI with my institution's prefix

**Acceptance Criteria:**
- Can configure DataCite credentials
- Can mint DOI in production
- Gets institutional DOI prefix
- Can provide landing page URL

**v0.3.0 Solution:**
```bash
# Configure DataCite (one-time setup)
cicada config set provider datacite
cicada config set datacite.repository_id MIT.BIO
cicada config set datacite.username my_username
cicada config set datacite.password my_password

# Mint DOI (NEW in v0.3.0)
cicada doi publish \
  --metadata doi-ready.json \
  --landing-page https://mylab.mit.edu/data/dataset-001 \
  --provider datacite

# Output: DOI: 10.12345/dataset-001
```

## Documentation Plan

### New Documentation (v0.3.0)

1. **Provider Integration Guide**
   - DataCite sandbox/production setup
   - Zenodo sandbox/production setup
   - Dryad setup
   - Figshare setup
   - Authentication and credentials
   - Troubleshooting provider issues

2. **Microscopy Formats Guide**
   - CZI metadata extraction
   - OME-TIFF metadata extraction
   - Zeiss preset usage
   - Multi-dimensional image handling

3. **Custom Preset Guide**
   - Creating custom presets
   - Preset definition format
   - Validation rules
   - Sharing presets
   - Best practices

4. **Interactive Editing Guide**
   - TUI navigation
   - Editing workflows
   - Keyboard shortcuts
   - Tips and tricks

5. **Production Deployment Guide**
   - Configuration management
   - Caching strategies
   - Error handling
   - Monitoring and logging
   - Backup and recovery

### Updated Documentation

- User Scenarios v0.3.0 (update with new features)
- README.md (add v0.3.0 features)
- CHANGELOG.md (track changes)
- Integration Testing Guide (new test patterns)

## Testing Strategy

### Integration Tests

- [ ] Provider integration tests (sandbox APIs)
- [ ] CZI extraction integration tests
- [ ] OME-TIFF extraction integration tests
- [ ] BAM/SAM extraction integration tests
- [ ] VCF extraction integration tests
- [ ] Custom preset integration tests
- [ ] Interactive UI integration tests (if possible)
- [ ] End-to-end DOI workflow tests

### Unit Tests

- [ ] Provider client unit tests (mocked APIs)
- [ ] Format parser unit tests
- [ ] Custom preset validation unit tests
- [ ] Cache implementation unit tests

### Manual Tests

- [ ] Test with real provider accounts
- [ ] Test with real microscopy files
- [ ] Test UI on different terminals
- [ ] Test on Linux, macOS, Windows

### Performance Tests

- [ ] Large CZI file extraction (< 5 seconds)
- [ ] Batch processing (1000 files)
- [ ] API rate limiting
- [ ] Cache hit rate

## Open Questions

1. **Provider Priority:** Should we focus on DataCite or Zenodo first?
   - *Suggestion:* Zenodo (easier auth, integrated storage, free)

2. **TUI Library:** Which library for interactive editing?
   - *Options:* bubbletea (modern, powerful) vs tview (mature, simpler)
   - *Suggestion:* bubbletea (better long-term)

3. **CZI Parsing:** Use existing library or implement from scratch?
   - *Options:* Port libCZI vs implement minimal reader
   - *Suggestion:* Implement minimal reader (lighter dependency)

4. **Caching Strategy:** What to cache and where?
   - *Options:* File-based vs SQLite vs memory
   - *Suggestion:* File-based for simplicity (JSON files in ~/.cicada/cache/)

5. **Custom Preset Format:** YAML or JSON?
   - *Suggestion:* YAML (more user-friendly, supports comments)

6. **Backward Compatibility:** How strict?
   - *Suggestion:* Strict for CLI, flexible for config (with migration)

## Next Steps

1. **Validate Roadmap** - Get feedback from target users
2. **Create GitHub Milestones** - Set up project tracking
3. **Break Down Milestone 1** - Create detailed task list for provider integration
4. **Set Up Provider Accounts** - Get sandbox credentials for testing
5. **Gather Test Data** - Collect sample CZI, OME-TIFF, BAM, VCF files
6. **Start Development** - Begin Milestone 1 (Provider Integration)

## Version History

- **v1.0** (2025-01-23): Initial v0.3.0 roadmap created after v0.2.0 release
