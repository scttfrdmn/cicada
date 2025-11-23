# Cicada v0.2.0 Implementation Plan

**Target**: Q1 2026 (14 weeks)
**Focus**: Metadata extraction + Instrument awareness + DOI minting + Basic multi-user

---

## Milestones

### Milestone 1: Metadata Foundation (Weeks 1-2)
**Due**: Week 2
**Goal**: Core metadata infrastructure in place

### Milestone 2: First Extractor Working (Weeks 3-4)
**Due**: Week 4
**Goal**: Zeiss CZI extraction functional end-to-end

### Milestone 3: Extraction Library Complete (Weeks 5-6)
**Due**: Week 6
**Goal**: OME-TIFF and FASTQ extractors working

### Milestone 4: Instrument Presets (Weeks 7-8)
**Due**: Week 8
**Goal**: Preset system with 5+ instruments

### Milestone 5: DOI Provider System (Weeks 9-10)
**Due**: Week 10
**Goal**: DataCite provider functional

### Milestone 6: Multi-User Foundation (Weeks 11-12)
**Due**: Week 12
**Goal**: Basic IAM automation for 2-3 users

### Milestone 7: Documentation & Release (Weeks 13-14)
**Due**: Week 14
**Goal**: v0.2.0 released with complete docs

---

## Issue Breakdown by Milestone

### Milestone 1: Metadata Foundation (Weeks 1-2)

#### Issue #1: Define Metadata Core Types
**Labels**: `enhancement`, `v0.2.0`, `metadata`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Define core metadata data structures in Go.

**Tasks**:
- [ ] Create `internal/metadata/types.go`
- [ ] Define `Metadata` struct with common fields
- [ ] Define instrument-specific structs (`MicroscopyMetadata`, `SequencingMetadata`, etc.)
- [ ] Add JSON/YAML tags
- [ ] Write unit tests for serialization
- [ ] Document all fields with comments

**Acceptance Criteria**:
- All types compile and pass tests
- 100% test coverage for type definitions
- Documentation generated with godoc

---

#### Issue #2: Implement Extractor Interface
**Labels**: `enhancement`, `v0.2.0`, `metadata`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Define and implement the pluggable extractor interface.

**Tasks**:
- [ ] Create `internal/metadata/extractor.go`
- [ ] Define `Extractor` interface
- [ ] Implement `ExtractorRegistry` with registration system
- [ ] Add `GenericExtractor` as fallback
- [ ] Write unit tests for registry
- [ ] Add example extractor stub

**Acceptance Criteria**:
- Interface is well-documented
- Registry can register/lookup extractors
- GenericExtractor works for any file
- Tests cover all registry operations

---

#### Issue #3: S3 Object Tagging Integration
**Labels**: `enhancement`, `v0.2.0`, `metadata`, `s3`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Add support for storing metadata as S3 object tags.

**Tasks**:
- [ ] Update `internal/sync/s3.go` with tagging support
- [ ] Implement `PutObjectTagging()` wrapper
- [ ] Implement `GetObjectTagging()` wrapper
- [ ] Add tag-to-metadata conversion functions
- [ ] Handle 10-tag limit (prioritize critical fields)
- [ ] Write integration tests with LocalStack
- [ ] Document S3 permissions needed

**Acceptance Criteria**:
- Tags written during upload
- Tags readable during list operations
- Tests verify tag persistence
- Documentation includes IAM policy examples

---

#### Issue #4: Sidecar JSON File Writer
**Labels**: `enhancement`, `v0.2.0`, `metadata`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Implement sidecar JSON file creation for full metadata storage.

**Tasks**:
- [ ] Create `internal/metadata/storage.go`
- [ ] Implement `WriteSidecarJSON(metadata, filepath)`
- [ ] Support both local and S3 destinations
- [ ] Use `.metadata.json` suffix convention
- [ ] Pretty-print JSON for readability
- [ ] Write unit tests
- [ ] Document sidecar file format

**Acceptance Criteria**:
- Sidecar files created alongside data files
- JSON is valid and readable
- Works for both local and S3 paths
- Tests verify file contents

---

#### Issue #5: Central Catalog Management
**Labels**: `enhancement`, `v0.2.0`, `metadata`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Implement central metadata catalog for searchable index.

**Tasks**:
- [ ] Design catalog JSON schema
- [ ] Implement incremental catalog updates
- [ ] Store catalog at `s3://bucket/.cicada/metadata-catalog.json`
- [ ] Add `AppendToCatalog(metadata)` function
- [ ] Add `ReadCatalog()` function
- [ ] Handle concurrent writes (optimistic locking)
- [ ] Write integration tests
- [ ] Document catalog format

**Acceptance Criteria**:
- Catalog updates incrementally
- Catalog is queryable JSON
- Concurrent updates handled gracefully
- Tests verify catalog integrity

---

### Milestone 2: First Extractor Working (Weeks 3-4)

#### Issue #6: Research CZI File Format
**Labels**: `research`, `v0.2.0`, `metadata`, `microscopy`
**Assignee**: TBD
**Estimate**: 1 day

**Description**:
Research Zeiss CZI file format and available parsing libraries.

**Tasks**:
- [ ] Review CZI file format specification
- [ ] Evaluate `pylibCZIrw` (Python library)
- [ ] Evaluate `bioformats` (Java library)
- [ ] Evaluate Go CGo bindings possibility
- [ ] Test sample CZI files
- [ ] Document recommended approach
- [ ] Create proof-of-concept

**Acceptance Criteria**:
- Decision documented on library choice
- Sample CZI files obtained for testing
- POC demonstrates basic metadata extraction

---

#### Issue #7: Implement Zeiss CZI Extractor
**Labels**: `enhancement`, `v0.2.0`, `metadata`, `microscopy`
**Assignee**: TBD
**Estimate**: 5 days

**Description**:
Implement full Zeiss CZI metadata extractor.

**Tasks**:
- [ ] Create `internal/metadata/zeiss.go`
- [ ] Implement `ZeissExtractor` struct
- [ ] Parse CZI header and metadata sections
- [ ] Extract instrument, objective, channels, dimensions
- [ ] Extract pixel size and acquisition date
- [ ] Handle errors gracefully
- [ ] Write comprehensive unit tests
- [ ] Add integration test with real CZI files
- [ ] Document supported CZI versions

**Acceptance Criteria**:
- Extracts all required metadata fields
- Handles corrupt/partial files gracefully
- <5s extraction time for typical files
- 90%+ test coverage
- Works with LSM 880, 900, 980 files

**Related**: Depends on #6

---

#### Issue #8: CZI Validation
**Labels**: `enhancement`, `v0.2.0`, `metadata`, `microscopy`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Implement CZI file validation logic.

**Tasks**:
- [ ] Verify CZI magic bytes (`ZISRAWFILE`)
- [ ] Check file structure integrity
- [ ] Validate metadata completeness
- [ ] Add `Validate()` method to `ZeissExtractor`
- [ ] Write unit tests with corrupt files
- [ ] Document validation errors

**Acceptance Criteria**:
- Detects corrupt CZI files
- Clear error messages
- Tests cover common corruption cases
- 99%+ accuracy on test corpus

**Related**: Depends on #7

---

#### Issue #9: CLI Integration for Metadata Extraction
**Labels**: `enhancement`, `v0.2.0`, `cli`, `metadata`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Add CLI commands for metadata extraction.

**Tasks**:
- [ ] Create `cmd/cicada/metadata.go`
- [ ] Add `cicada metadata extract <path>` command
- [ ] Add `cicada metadata show <path>` command
- [ ] Add `cicada metadata validate <path>` command
- [ ] Add `--format` flag (json, yaml, table)
- [ ] Add `--extractor` flag to force specific extractor
- [ ] Add colored output for validation results
- [ ] Update `cicada sync` with `--extract-metadata` flag
- [ ] Write CLI tests
- [ ] Update documentation

**Acceptance Criteria**:
- All commands work as documented
- Output is clear and actionable
- Sync integration works seamlessly
- Tests verify CLI behavior

**Related**: Depends on #7

---

### Milestone 3: Extraction Library Complete (Weeks 5-6)

#### Issue #10: Implement OME-TIFF Extractor
**Labels**: `enhancement`, `v0.2.0`, `metadata`, `microscopy`
**Assignee**: TBD
**Estimate**: 4 days

**Description**:
Implement OME-TIFF metadata extractor.

**Tasks**:
- [ ] Create `internal/metadata/ometiff.go`
- [ ] Parse TIFF tags to find OME-XML
- [ ] Parse OME-XML metadata
- [ ] Map OME schema to Cicada metadata
- [ ] Handle multiple images in one file
- [ ] Write unit tests with sample files
- [ ] Document OME-TIFF support

**Acceptance Criteria**:
- Extracts full OME metadata
- Handles BigTIFF format
- <3s extraction for typical files
- 90%+ test coverage

---

#### Issue #11: Implement FASTQ Extractor
**Labels**: `enhancement`, `v0.2.0`, `metadata`, `sequencing`
**Assignee**: TBD
**Estimate**: 4 days

**Description**:
Implement FASTQ metadata extractor.

**Tasks**:
- [ ] Create `internal/metadata/fastq.go`
- [ ] Parse FASTQ format (4-line records)
- [ ] Detect quality encoding (Phred+33/Phred+64)
- [ ] Sample reads to estimate total count
- [ ] Calculate quality score statistics
- [ ] Detect read length distribution
- [ ] Handle gzipped FASTQ files
- [ ] Write unit tests
- [ ] Document FASTQ support

**Acceptance Criteria**:
- Correctly detects quality encoding
- Accurate read count estimation
- <10s for 1GB file
- 90%+ test coverage

---

#### Issue #12: Performance Benchmarking
**Labels**: `testing`, `v0.2.0`, `metadata`, `performance`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Create performance benchmarks for all extractors.

**Tasks**:
- [ ] Create `internal/metadata/bench_test.go`
- [ ] Benchmark CZI extraction (100MB, 1GB, 5GB files)
- [ ] Benchmark OME-TIFF extraction
- [ ] Benchmark FASTQ extraction
- [ ] Document performance targets
- [ ] Add CI benchmark job
- [ ] Create performance regression tests

**Acceptance Criteria**:
- All benchmarks under target times
- CI runs benchmarks on PRs
- Performance regression caught automatically

**Related**: Depends on #7, #10, #11

---

### Milestone 4: Instrument Presets (Weeks 7-8)

#### Issue #13: Design Preset YAML Format
**Labels**: `enhancement`, `v0.2.0`, `presets`
**Assignee**: TBD
**Estimate**: 1 day

**Description**:
Design and document preset YAML schema.

**Tasks**:
- [ ] Define preset YAML structure
- [ ] Add detection rules (extensions, magic bytes)
- [ ] Add sync configuration fields
- [ ] Add metadata extraction settings
- [ ] Add validation rules
- [ ] Document schema with examples
- [ ] Create JSON schema for validation

**Acceptance Criteria**:
- Schema is well-documented
- Examples cover common use cases
- JSON schema validates presets

---

#### Issue #14: Implement Preset Loader
**Labels**: `enhancement`, `v0.2.0`, `presets`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Implement preset loading and management system.

**Tasks**:
- [ ] Create `internal/preset/loader.go`
- [ ] Implement `LoadPreset(path)` function
- [ ] Implement preset validation
- [ ] Add embedded preset support
- [ ] Create preset registry
- [ ] Write unit tests
- [ ] Document preset API

**Acceptance Criteria**:
- Loads presets from YAML files
- Validates preset structure
- Registry manages all presets
- Tests verify loading logic

**Related**: Depends on #13

---

#### Issue #15: Implement Auto-Detection
**Labels**: `enhancement`, `v0.2.0`, `presets`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Implement instrument auto-detection logic.

**Tasks**:
- [ ] Create `internal/preset/detect.go`
- [ ] Implement file extension matching
- [ ] Implement magic byte detection
- [ ] Score confidence of matches
- [ ] Return best match with confidence level
- [ ] Write unit tests with sample files
- [ ] Document detection algorithm

**Acceptance Criteria**:
- 95%+ detection accuracy
- Clear confidence scores
- Handles ambiguous cases gracefully
- Tests cover all presets

**Related**: Depends on #14

---

#### Issue #16: Create Initial Presets
**Labels**: `enhancement`, `v0.2.0`, `presets`, `microscopy`, `sequencing`
**Assignee**: TBD
**Estimate**: 4 days

**Description**:
Create initial preset library for common instruments.

**Tasks**:
- [ ] Create `presets/microscopy/zeiss-confocal.yaml`
- [ ] Create `presets/microscopy/zeiss-lightsheet.yaml`
- [ ] Create `presets/sequencing/illumina-novaseq.yaml`
- [ ] Create `presets/sequencing/illumina-miseq.yaml`
- [ ] Create `presets/generic/large-files.yaml`
- [ ] Create `presets/generic/many-small-files.yaml`
- [ ] Create `presets/README.md` with documentation
- [ ] Test each preset with real data
- [ ] Document preset usage

**Acceptance Criteria**:
- 6+ presets created
- Each preset tested with real instrument data
- Documentation explains all fields
- Presets follow consistent style

**Related**: Depends on #13

---

#### Issue #17: Preset CLI Commands
**Labels**: `enhancement`, `v0.2.0`, `cli`, `presets`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Add CLI commands for preset management.

**Tasks**:
- [ ] Create `cmd/cicada/instrument.go`
- [ ] Add `cicada instrument list` command
- [ ] Add `cicada instrument show <preset>` command
- [ ] Add `cicada instrument detect <path>` command
- [ ] Add `cicada instrument setup <preset>` command
- [ ] Add interactive setup wizard
- [ ] Write CLI tests
- [ ] Update documentation

**Acceptance Criteria**:
- All commands work as documented
- Interactive wizard is user-friendly
- Tests verify CLI behavior
- Documentation includes examples

**Related**: Depends on #14, #15, #16

---

### Milestone 5: DOI Provider System (Weeks 9-10)

#### Issue #18: Define DOI Provider Interface
**Labels**: `enhancement`, `v0.2.0`, `doi`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Define pluggable DOI provider interface.

**Tasks**:
- [ ] Create `internal/doi/provider.go`
- [ ] Define `Provider` interface
- [ ] Define `Dataset` struct with required fields
- [ ] Define `DOI` struct
- [ ] Create `ProviderRegistry`
- [ ] Add configuration structs
- [ ] Write unit tests
- [ ] Document interface design

**Acceptance Criteria**:
- Interface supports multiple providers
- Dataset struct covers DataCite/Zenodo needs
- Registry manages providers
- Tests verify interface contracts

---

#### Issue #19: Implement Disabled Provider
**Labels**: `enhancement`, `v0.2.0`, `doi`
**Assignee**: TBD
**Estimate**: 1 day

**Description**:
Implement "disabled" DOI provider for when minting is off.

**Tasks**:
- [ ] Create `internal/doi/disabled.go`
- [ ] Implement `DisabledProvider` struct
- [ ] Return clear error messages
- [ ] Write unit tests
- [ ] Document usage

**Acceptance Criteria**:
- Returns helpful error when disabled
- Tests verify error messages
- Documentation explains configuration

**Related**: Depends on #18

---

#### Issue #20: Implement DataCite Provider (Stub)
**Labels**: `enhancement`, `v0.2.0`, `doi`, `datacite`
**Assignee**: TBD
**Estimate**: 5 days

**Description**:
Implement DataCite provider with basic functionality.

**Tasks**:
- [ ] Create `internal/doi/datacite.go`
- [ ] Implement `DataCiteProvider` struct
- [ ] Add DataCite API client
- [ ] Implement `Mint()` method
- [ ] Generate DataCite metadata XML
- [ ] Add test mode support
- [ ] Write unit tests with mock API
- [ ] Document DataCite setup

**Acceptance Criteria**:
- Can mint DOI in test mode
- Metadata XML validates against schema
- Error handling is robust
- Tests verify API interactions
- Documentation includes credential setup

**Related**: Depends on #18

---

#### Issue #21: DOI Configuration System
**Labels**: `enhancement`, `v0.2.0`, `doi`, `config`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Add DOI configuration to Cicada config system.

**Tasks**:
- [ ] Update `internal/config/config.go`
- [ ] Add `DOI` configuration section
- [ ] Support provider selection (datacite, zenodo, off)
- [ ] Add provider-specific config sections
- [ ] Support environment variable substitution
- [ ] Write unit tests
- [ ] Document configuration

**Acceptance Criteria**:
- Config supports all providers
- Environment variables work
- Validation catches misconfigurations
- Tests verify all config paths

**Related**: Depends on #18

---

#### Issue #22: DOI CLI Commands
**Labels**: `enhancement`, `v0.2.0`, `cli`, `doi`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Add CLI commands for DOI management.

**Tasks**:
- [ ] Create `cmd/cicada/publish.go`
- [ ] Add `cicada publish <path>` command
- [ ] Add interactive metadata prompts
- [ ] Add `--dry-run` flag
- [ ] Add `cicada publish list` command
- [ ] Add `cicada publish show <doi>` command
- [ ] Write CLI tests
- [ ] Update documentation

**Acceptance Criteria**:
- Can mint DOI via CLI
- Interactive prompts are clear
- Dry-run shows preview
- Tests verify CLI behavior

**Related**: Depends on #20, #21

---

### Milestone 6: Multi-User Foundation (Weeks 11-12)

#### Issue #23: IAM User Creation Commands
**Labels**: `enhancement`, `v0.2.0`, `auth`, `iam`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Add CLI commands for creating IAM users.

**Tasks**:
- [ ] Create `cmd/cicada/user.go`
- [ ] Add `cicada user add <name>` command
- [ ] Create IAM user via AWS SDK
- [ ] Generate access keys
- [ ] Store credentials securely
- [ ] Add `cicada user list` command
- [ ] Add `cicada user remove <name>` command
- [ ] Write CLI tests
- [ ] Document user management

**Acceptance Criteria**:
- Creates IAM users successfully
- Credentials stored securely
- Tests verify user creation
- Documentation includes IAM prerequisites

---

#### Issue #24: Path-Based IAM Policy Generation
**Labels**: `enhancement`, `v0.2.0`, `auth`, `iam`
**Assignee**: TBD
**Estimate**: 4 days

**Description**:
Implement automatic IAM policy generation for path-based access.

**Tasks**:
- [ ] Create `internal/auth/policy.go`
- [ ] Implement policy template system
- [ ] Generate policies for S3 path prefixes
- [ ] Support read-only vs read-write access
- [ ] Attach policies to users
- [ ] Add `cicada user grant <user> <path>` command
- [ ] Add `cicada user revoke <user> <path>` command
- [ ] Write unit tests
- [ ] Document policy generation

**Acceptance Criteria**:
- Policies follow least-privilege principle
- Path-based access works correctly
- Tests verify policy generation
- Documentation explains access patterns

**Related**: Depends on #23

---

#### Issue #25: Project Management System
**Labels**: `enhancement`, `v0.2.0`, `projects`
**Assignee**: TBD
**Estimate**: 3 days

**Description**:
Add basic project management for organizing data.

**Tasks**:
- [ ] Create `internal/config/project.go`
- [ ] Define `Project` struct
- [ ] Add `cicada project create <name>` command
- [ ] Add `cicada project list` command
- [ ] Associate S3 paths with projects
- [ ] Store projects in config
- [ ] Write unit tests
- [ ] Document project usage

**Acceptance Criteria**:
- Projects can be created and listed
- Projects map to S3 paths
- Config persists projects
- Tests verify project management
- Documentation includes examples

---

#### Issue #26: Multi-User Documentation
**Labels**: `documentation`, `v0.2.0`, `auth`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Create comprehensive multi-user setup guide.

**Tasks**:
- [ ] Create `docs/MULTI_USER_SETUP.md`
- [ ] Document IAM prerequisites
- [ ] Provide step-by-step setup guide
- [ ] Include example scenarios (2-3 users)
- [ ] Document common pitfalls
- [ ] Add troubleshooting section
- [ ] Include security best practices

**Acceptance Criteria**:
- Guide enables self-service setup
- All commands documented
- Examples are realistic
- Security considerations covered

**Related**: Depends on #23, #24, #25

---

### Milestone 7: Documentation & Release (Weeks 13-14)

#### Issue #27: Update USER_SCENARIOS for v0.2.0
**Labels**: `documentation`, `v0.2.0`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Update USER_SCENARIOS.md with v0.2.0 features.

**Tasks**:
- [ ] Add metadata extraction to scenarios
- [ ] Add instrument preset usage
- [ ] Add DOI minting example
- [ ] Add multi-user collaboration example
- [ ] Update all command examples
- [ ] Review for accuracy

**Acceptance Criteria**:
- All scenarios updated
- New features demonstrated
- Commands are correct
- Examples are realistic

---

#### Issue #28: Create v0.2.0 Migration Guide
**Labels**: `documentation`, `v0.2.0`
**Assignee**: TBD
**Estimate**: 1 day

**Description**:
Create migration guide for v0.1.0 → v0.2.0.

**Tasks**:
- [ ] Create `docs/MIGRATION_v0.2.0.md`
- [ ] Document breaking changes (if any)
- [ ] Document new configuration options
- [ ] Provide upgrade steps
- [ ] Include rollback instructions

**Acceptance Criteria**:
- Guide is clear and complete
- All breaking changes documented
- Upgrade path is straightforward

---

#### Issue #29: Integration Testing Suite
**Labels**: `testing`, `v0.2.0`
**Assignee**: TBD
**Estimate**: 4 days

**Description**:
Create comprehensive integration test suite.

**Tasks**:
- [ ] Create `internal/integration/metadata_test.go`
- [ ] Test end-to-end metadata extraction during sync
- [ ] Test preset auto-detection
- [ ] Test DOI minting (test mode)
- [ ] Test multi-user scenarios
- [ ] Add CI integration test job
- [ ] Document test setup

**Acceptance Criteria**:
- All major workflows tested
- Tests run in CI
- Tests catch regressions
- Documentation explains test setup

---

#### Issue #30: Performance & Load Testing
**Labels**: `testing`, `v0.2.0`, `performance`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Validate performance targets for v0.2.0.

**Tasks**:
- [ ] Test metadata extraction with 1000 files
- [ ] Test sync overhead with metadata enabled
- [ ] Test catalog update performance
- [ ] Test concurrent extraction
- [ ] Document performance results
- [ ] Identify bottlenecks

**Acceptance Criteria**:
- All performance targets met
- No regressions from v0.1.0
- Bottlenecks documented
- Results inform v0.3.0 optimization

---

#### Issue #31: Release Preparation
**Labels**: `release`, `v0.2.0`
**Assignee**: TBD
**Estimate**: 2 days

**Description**:
Prepare v0.2.0 for release.

**Tasks**:
- [ ] Update CHANGELOG.md
- [ ] Update version in all files
- [ ] Create release notes
- [ ] Build binaries for all platforms
- [ ] Test release binaries
- [ ] Create GitHub release
- [ ] Update README.md badges
- [ ] Announce release

**Acceptance Criteria**:
- CHANGELOG is complete
- All binaries built and tested
- GitHub release created
- Documentation updated

---

#### Issue #32: v0.2.0 Announcement & Outreach
**Labels**: `release`, `v0.2.0`, `community`
**Assignee**: TBD
**Estimate**: 1 day

**Description**:
Announce v0.2.0 release to community.

**Tasks**:
- [ ] Write release blog post
- [ ] Create demo video
- [ ] Post to relevant communities
- [ ] Update project website
- [ ] Notify early adopters

**Acceptance Criteria**:
- Blog post published
- Demo video available
- Community notified
- Website updated

---

## Summary

**Total Issues**: 32
**Total Duration**: 14 weeks
**Dependencies**: Sequential with some parallel work possible

### Critical Path:
1. Metadata Foundation (Issues #1-5) → 2 weeks
2. First Extractor (Issues #6-9) → 2 weeks
3. Extraction Library (Issues #10-12) → 2 weeks
4. Instrument Presets (Issues #13-17) → 2 weeks
5. DOI Providers (Issues #18-22) → 2 weeks
6. Multi-User (Issues #23-26) → 2 weeks
7. Documentation & Release (Issues #27-32) → 2 weeks

### Parallelization Opportunities:
- Issues #10 and #11 can be done in parallel (Week 5-6)
- Issues #18-20 can be split among multiple developers (Week 9-10)
- Documentation issues #27-28 can be done while testing (Week 13-14)
