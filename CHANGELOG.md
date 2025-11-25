# Changelog

All notable changes to Cicada will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2025-11-25

Documentation release providing enterprise-grade documentation for users and developers.

### Added

- **Comprehensive Project Documentation** (8 files, 9,052 lines):
  - ARCHITECTURE.md (1,090 lines): System architecture, components, data flow diagrams
  - CLI_REFERENCE.md (984 lines): Complete CLI command reference with 50+ examples
  - CONFIGURATION.md (1,263 lines): Configuration options, AWS setup, 10+ example configs
  - STORAGE.md (1,042 lines): Local filesystem and AWS S3 backend documentation
  - METADATA_SYSTEM.md (1,463 lines): 14 extractors, 6 instrument types, validation
  - INTEGRATIONS.md (1,257 lines): Nextflow, Snakemake, Python, Bash, CI/CD examples
  - DEVELOPMENT.md (988 lines): Developer guide, building, testing, contributing
  - API.md (965 lines): Complete Go package API reference

- **Comprehensive User Guide** (4 files, 5,022 lines):
  - GETTING_STARTED.md (825 lines): Installation, setup, first sync, AWS S3 integration
  - WORKFLOWS.md (1,297 lines): Microscopy, sequencing, mass spec, collaboration workflows
  - TROUBLESHOOTING.md (1,305 lines): Common problems, solutions, debug tools
  - ADVANCED.md (1,595 lines): Custom extractors, hooks, security, performance optimization

- **GitHub Pages Documentation Site**:
  - MkDocs Material theme with navigation and search
  - Deployed at https://scttfrdmn.github.io/cicada
  - Organized by user type (getting started, user guide, advanced, reference, development)
  - Responsive design with light/dark mode

- **Official Branding**:
  - BRANDING.md: Comprehensive visual identity guidelines
  - Official Cicada mascot (scholarly cicada with glasses, backpack, and book)
  - Color palette, typography, and voice guidelines
  - Usage examples and restrictions
  - Consistent branding across GitHub and documentation site

- **Documentation Updates**:
  - Updated all existing documentation to reflect data commons platform vision
  - Repositioned from DOI-centric to comprehensive data management platform
  - Updated README with new platform description and centered mascot
  - Fixed GitHub Pages deployment workflow

### Changed

- Repositioned project focus from DOI registration to comprehensive small lab data commons
- Enhanced onboarding experience with step-by-step tutorials
- Improved visual identity with official mascot and branding guidelines

### Fixed

- GitHub Pages deployment workflow (removed --strict flag causing build failures)
- Documentation date consistency (all docs now use correct 2025-11-25 date)

### Documentation Statistics

- **Total**: 12 comprehensive documentation files
- **Total Lines**: 14,074+ lines of documentation
- **Coverage**: Complete coverage for users, developers, and administrators
- **Quality**: Enterprise-grade documentation suitable for academic and commercial use

## [0.2.0] - 2025-01-23

Major release transforming Cicada into a comprehensive data commons platform with automated metadata extraction, multi-format support, data quality validation, and optional DOI preparation capabilities.

### Added

- **Metadata Extraction System**:
  - 14 file format extractors across multiple domains
  - Microscopy: TIFF, OME-TIFF, Zeiss CZI, Nikon ND2, Leica LIF
  - Sequencing: FASTQ, BAM
  - Mass Spectrometry: mzML, MGF
  - Data Arrays: HDF5, Zarr
  - Medical/Flow: DICOM, FCS
  - Generic fallback extractor
  - 6 instrument-specific metadata types (Microscopy, Sequencing, Mass Spec, Flow Cytometry, Cryo-EM, X-Ray)
  - S3 object tagging integration for metadata storage
  - Automatic file format detection
  - Thread-safe concurrent extraction
  - Extractor registry with plugin architecture
  - CLI commands: `cicada metadata extract`, `cicada metadata validate`, `cicada metadata show`, `cicada metadata list`

- **Instrument Preset System**:
  - 8 default presets: Illumina (NovaSeq, MiSeq, NextSeq), Zeiss (LSM 880/900/980), Generic (sequencing, microscopy)
  - Field-level validation (required vs optional fields)
  - Quality scoring system (0-100 scale: 60% required + 40% optional)
  - Preset search by manufacturer and instrument type
  - Template generation support
  - CLI commands: `cicada metadata preset list`, `cicada metadata preset show`

- **DOI Preparation Workflow**:
  - DataCite Metadata Schema v4.5 mapping
  - DOI readiness validation (6 required + 14 recommended fields)
  - Quality scoring with actionable recommendations
  - Metadata enrichment from YAML/JSON files
  - Multi-file dataset support
  - Author/creator handling with ORCID support
  - Related identifier tracking
  - Funding reference support
  - CLI commands: `cicada doi prepare`, `cicada doi validate`

- **Provider Infrastructure**:
  - Provider registry (DataCite, Zenodo, future: Dryad, Figshare)
  - API structure for DOI minting and updates
  - Sandbox and production environment support
  - Configuration management for provider credentials

- **Testing**:
  - 29 integration tests using real data (no mocks)
  - 11 performance benchmarks
  - 129 total tests (all passing)
  - 50-83% test coverage across packages
  - < 1 second integration test runtime

- **Performance**:
  - Small files: 31 μs per extraction (32,268 ops/sec)
  - Medium files: 128 μs per extraction (7,809 ops/sec)
  - Large files: 1 ms per extraction (constant due to sampling)
  - Complete DOI workflow: 36 μs end-to-end
  - Preset validation: 478 ns (sub-microsecond)
  - 4-8x speedup with concurrent processing

- **Documentation** (5,550+ lines):
  - Metadata Extraction Guide (800+ lines)
  - DOI Workflow Guide (900+ lines)
  - Instrument Preset Guide (900+ lines)
  - Provider Setup Guide (1,000+ lines)
  - User Scenarios v0.2.0 with 5 personas (1,950 lines)
  - Integration Testing Guide (377 lines)
  - Performance Benchmarks (400+ lines)
  - CLI examples for Nextflow, Snakemake, Python, Bash

### Changed

- Updated user scenarios with v0.2.0 metadata features
- Expanded target user profile to include small labs (2-10 people)

### Performance

- Metadata extraction: 31 μs - 1 ms per file
- Monthly processing for small lab (200 files): < 30 ms
- Annual archive (10,000 files): < 2 seconds
- Memory efficiency: < 1 MB per concurrent operation

### Known Limitations

- Provider integration: API structure complete, live API calls for DataCite/Zenodo planned for v0.4.0+
- Custom presets: 8 built-in presets, user-defined presets planned for v0.4.0+
- Metadata enrichment: Manual YAML/JSON editing (interactive UI planned for v0.4.0+)

### Breaking Changes

**None.** v0.2.0 is fully backward compatible with v0.1.0.

## [0.1.0] - 2025-11-23

Initial release of Cicada - Foundational storage and sync layer for the dormant data commons platform.

### Added

- **Sync Engine**: rsync-like sync between local filesystem and S3
  - Concurrent transfers with configurable concurrency
  - MD5/ETag comparison for change detection
  - Dry-run mode for previewing changes
  - Delete mode for mirror syncing
  - Exclude patterns support

- **File Watching**: Auto-sync directories on file changes
  - fsnotify integration for real-time monitoring
  - Debouncing to prevent sync storms
  - Min-age filtering for partial file writes
  - Recursive directory watching
  - Watch persistence via configuration

- **Configuration Management**: YAML-based configuration
  - AWS profile and region configuration
  - Default sync options
  - Watch persistence
  - Global settings
  - CLI commands: init, set, get, list

- **S3 Backend**: Full S3 integration
  - List, read, write, stat, delete operations
  - AWS SDK v2 integration
  - Multi-region support
  - Profile-based authentication

- **Local Backend**: Filesystem operations
  - MD5 checksum calculation
  - Recursive directory listing
  - Efficient file I/O

- **CLI Commands**:
  - `cicada sync`: Manual sync operations
  - `cicada watch`: File watching management
  - `cicada config`: Configuration management
  - `cicada version`: Version information

- **Testing**:
  - Comprehensive unit tests (80%+ coverage)
  - AWS S3 integration tests
  - Makefile targets for local testing

- **Documentation**:
  - Comprehensive README with usage examples
  - Integration test setup guide
  - Development instructions
  - Contributing guidelines

- **CI/CD**:
  - GitHub Actions workflow for build and test
  - goreleaser configuration for releases
  - Homebrew tap integration
  - Scoop bucket integration

### Known Limitations

- No multipart upload support (files >5GB may be slow)
- No bandwidth throttling
- Sync engine path handling needs improvement for some use cases
- No background daemon (runs in foreground)
- No resume capability for interrupted transfers

[Unreleased]: https://github.com/scttfrdmn/cicada/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/scttfrdmn/cicada/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/scttfrdmn/cicada/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/scttfrdmn/cicada/releases/tag/v0.1.0
