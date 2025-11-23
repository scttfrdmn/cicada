#!/bin/bash

# Script to create all v0.2.0 GitHub issues
# Usage: ./scripts/create_v0.2.0_issues.sh

set -e

echo "Creating v0.2.0 GitHub issues..."

# Milestone 1: Metadata Foundation

gh issue create \
  --title "Define Metadata Core Types" \
  --label "enhancement,v0.2.0,metadata" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Define core metadata data structures in Go.

**Tasks**:
- [ ] Create \`internal/metadata/types.go\`
- [ ] Define \`Metadata\` struct with common fields
- [ ] Define instrument-specific structs (\`MicroscopyMetadata\`, \`SequencingMetadata\`, etc.)
- [ ] Add JSON/YAML tags
- [ ] Write unit tests for serialization
- [ ] Document all fields with comments

**Acceptance Criteria**:
- All types compile and pass tests
- 100% test coverage for type definitions
- Documentation generated with godoc

**Estimate**: 2 days
**Milestone 1** (Weeks 1-2)"

gh issue create \
  --title "Implement Extractor Interface" \
  --label "enhancement,v0.2.0,metadata" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Define and implement the pluggable extractor interface.

**Tasks**:
- [ ] Create \`internal/metadata/extractor.go\`
- [ ] Define \`Extractor\` interface
- [ ] Implement \`ExtractorRegistry\` with registration system
- [ ] Add \`GenericExtractor\` as fallback
- [ ] Write unit tests for registry
- [ ] Add example extractor stub

**Acceptance Criteria**:
- Interface is well-documented
- Registry can register/lookup extractors
- GenericExtractor works for any file
- Tests cover all registry operations

**Estimate**: 2 days
**Milestone 1** (Weeks 1-2)"

gh issue create \
  --title "S3 Object Tagging Integration" \
  --label "enhancement,v0.2.0,metadata,s3" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Add support for storing metadata as S3 object tags.

**Tasks**:
- [ ] Update \`internal/sync/s3.go\` with tagging support
- [ ] Implement \`PutObjectTagging()\` wrapper
- [ ] Implement \`GetObjectTagging()\` wrapper
- [ ] Add tag-to-metadata conversion functions
- [ ] Handle 10-tag limit (prioritize critical fields)
- [ ] Write integration tests with LocalStack
- [ ] Document S3 permissions needed

**Acceptance Criteria**:
- Tags written during upload
- Tags readable during list operations
- Tests verify tag persistence
- Documentation includes IAM policy examples

**Estimate**: 3 days
**Milestone 1** (Weeks 1-2)"

gh issue create \
  --title "Sidecar JSON File Writer" \
  --label "enhancement,v0.2.0,metadata" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement sidecar JSON file creation for full metadata storage.

**Tasks**:
- [ ] Create \`internal/metadata/storage.go\`
- [ ] Implement \`WriteSidecarJSON(metadata, filepath)\`
- [ ] Support both local and S3 destinations
- [ ] Use \`.metadata.json\` suffix convention
- [ ] Pretty-print JSON for readability
- [ ] Write unit tests
- [ ] Document sidecar file format

**Acceptance Criteria**:
- Sidecar files created alongside data files
- JSON is valid and readable
- Works for both local and S3 paths
- Tests verify file contents

**Estimate**: 2 days
**Milestone 1** (Weeks 1-2)"

gh issue create \
  --title "Central Catalog Management" \
  --label "enhancement,v0.2.0,metadata" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement central metadata catalog for searchable index.

**Tasks**:
- [ ] Design catalog JSON schema
- [ ] Implement incremental catalog updates
- [ ] Store catalog at \`s3://bucket/.cicada/metadata-catalog.json\`
- [ ] Add \`AppendToCatalog(metadata)\` function
- [ ] Add \`ReadCatalog()\` function
- [ ] Handle concurrent writes (optimistic locking)
- [ ] Write integration tests
- [ ] Document catalog format

**Acceptance Criteria**:
- Catalog updates incrementally
- Catalog is queryable JSON
- Concurrent updates handled gracefully
- Tests verify catalog integrity

**Estimate**: 3 days
**Milestone 1** (Weeks 1-2)"

# Milestone 2: First Extractor Working

gh issue create \
  --title "Research CZI File Format" \
  --label "research,v0.2.0,metadata,microscopy" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Research Zeiss CZI file format and available parsing libraries.

**Tasks**:
- [ ] Review CZI file format specification
- [ ] Evaluate \`pylibCZIrw\` (Python library)
- [ ] Evaluate \`bioformats\` (Java library)
- [ ] Evaluate Go CGo bindings possibility
- [ ] Test sample CZI files
- [ ] Document recommended approach
- [ ] Create proof-of-concept

**Acceptance Criteria**:
- Decision documented on library choice
- Sample CZI files obtained for testing
- POC demonstrates basic metadata extraction

**Estimate**: 1 day
**Milestone 2** (Weeks 3-4)"

gh issue create \
  --title "Implement Zeiss CZI Extractor" \
  --label "enhancement,v0.2.0,metadata,microscopy" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement full Zeiss CZI metadata extractor.

**Tasks**:
- [ ] Create \`internal/metadata/zeiss.go\`
- [ ] Implement \`ZeissExtractor\` struct
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

**Estimate**: 5 days
**Milestone 2** (Weeks 3-4)"

gh issue create \
  --title "CZI Validation" \
  --label "enhancement,v0.2.0,metadata,microscopy" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement CZI file validation logic.

**Tasks**:
- [ ] Verify CZI magic bytes (\`ZISRAWFILE\`)
- [ ] Check file structure integrity
- [ ] Validate metadata completeness
- [ ] Add \`Validate()\` method to \`ZeissExtractor\`
- [ ] Write unit tests with corrupt files
- [ ] Document validation errors

**Acceptance Criteria**:
- Detects corrupt CZI files
- Clear error messages
- Tests cover common corruption cases
- 99%+ accuracy on test corpus

**Estimate**: 2 days
**Milestone 2** (Weeks 3-4)"

gh issue create \
  --title "CLI Integration for Metadata Extraction" \
  --label "enhancement,v0.2.0,cli,metadata" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Add CLI commands for metadata extraction.

**Tasks**:
- [ ] Create \`cmd/cicada/metadata.go\`
- [ ] Add \`cicada metadata extract <path>\` command
- [ ] Add \`cicada metadata show <path>\` command
- [ ] Add \`cicada metadata validate <path>\` command
- [ ] Add \`--format\` flag (json, yaml, table)
- [ ] Add \`--extractor\` flag to force specific extractor
- [ ] Add colored output for validation results
- [ ] Update \`cicada sync\` with \`--extract-metadata\` flag
- [ ] Write CLI tests
- [ ] Update documentation

**Acceptance Criteria**:
- All commands work as documented
- Output is clear and actionable
- Sync integration works seamlessly
- Tests verify CLI behavior

**Estimate**: 3 days
**Milestone 2** (Weeks 3-4)"

# Milestone 3: Extraction Library Complete

gh issue create \
  --title "Implement OME-TIFF Extractor" \
  --label "enhancement,v0.2.0,metadata,microscopy" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement OME-TIFF metadata extractor.

**Tasks**:
- [ ] Create \`internal/metadata/ometiff.go\`
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

**Estimate**: 4 days
**Milestone 3** (Weeks 5-6)"

gh issue create \
  --title "Implement FASTQ Extractor" \
  --label "enhancement,v0.2.0,metadata,sequencing" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement FASTQ metadata extractor.

**Tasks**:
- [ ] Create \`internal/metadata/fastq.go\`
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

**Estimate**: 4 days
**Milestone 3** (Weeks 5-6)"

gh issue create \
  --title "Performance Benchmarking" \
  --label "testing,v0.2.0,metadata,performance" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Create performance benchmarks for all extractors.

**Tasks**:
- [ ] Create \`internal/metadata/bench_test.go\`
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

**Estimate**: 2 days
**Milestone 3** (Weeks 5-6)"

# Milestone 4: Instrument Presets

gh issue create \
  --title "Design Preset YAML Format" \
  --label "enhancement,v0.2.0,presets" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Design and document preset YAML schema.

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

**Estimate**: 1 day
**Milestone 4** (Weeks 7-8)"

gh issue create \
  --title "Implement Preset Loader" \
  --label "enhancement,v0.2.0,presets" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement preset loading and management system.

**Tasks**:
- [ ] Create \`internal/preset/loader.go\`
- [ ] Implement \`LoadPreset(path)\` function
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

**Estimate**: 3 days
**Milestone 4** (Weeks 7-8)"

gh issue create \
  --title "Implement Auto-Detection" \
  --label "enhancement,v0.2.0,presets" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement instrument auto-detection logic.

**Tasks**:
- [ ] Create \`internal/preset/detect.go\`
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

**Estimate**: 3 days
**Milestone 4** (Weeks 7-8)"

gh issue create \
  --title "Create Initial Presets" \
  --label "enhancement,v0.2.0,presets,microscopy,sequencing" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Create initial preset library for common instruments.

**Tasks**:
- [ ] Create \`presets/microscopy/zeiss-confocal.yaml\`
- [ ] Create \`presets/microscopy/zeiss-lightsheet.yaml\`
- [ ] Create \`presets/sequencing/illumina-novaseq.yaml\`
- [ ] Create \`presets/sequencing/illumina-miseq.yaml\`
- [ ] Create \`presets/generic/large-files.yaml\`
- [ ] Create \`presets/generic/many-small-files.yaml\`
- [ ] Create \`presets/README.md\` with documentation
- [ ] Test each preset with real data
- [ ] Document preset usage

**Acceptance Criteria**:
- 6+ presets created
- Each preset tested with real instrument data
- Documentation explains all fields
- Presets follow consistent style

**Estimate**: 4 days
**Milestone 4** (Weeks 7-8)"

gh issue create \
  --title "Preset CLI Commands" \
  --label "enhancement,v0.2.0,cli,presets" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Add CLI commands for preset management.

**Tasks**:
- [ ] Create \`cmd/cicada/instrument.go\`
- [ ] Add \`cicada instrument list\` command
- [ ] Add \`cicada instrument show <preset>\` command
- [ ] Add \`cicada instrument detect <path>\` command
- [ ] Add \`cicada instrument setup <preset>\` command
- [ ] Add interactive setup wizard
- [ ] Write CLI tests
- [ ] Update documentation

**Acceptance Criteria**:
- All commands work as documented
- Interactive wizard is user-friendly
- Tests verify CLI behavior
- Documentation includes examples

**Estimate**: 3 days
**Milestone 4** (Weeks 7-8)"

# Milestone 5: DOI Provider System

gh issue create \
  --title "Define DOI Provider Interface" \
  --label "enhancement,v0.2.0,doi" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Define pluggable DOI provider interface.

**Tasks**:
- [ ] Create \`internal/doi/provider.go\`
- [ ] Define \`Provider\` interface
- [ ] Define \`Dataset\` struct with required fields
- [ ] Define \`DOI\` struct
- [ ] Create \`ProviderRegistry\`
- [ ] Add configuration structs
- [ ] Write unit tests
- [ ] Document interface design

**Acceptance Criteria**:
- Interface supports multiple providers
- Dataset struct covers DataCite/Zenodo needs
- Registry manages providers
- Tests verify interface contracts

**Estimate**: 2 days
**Milestone 5** (Weeks 9-10)"

gh issue create \
  --title "Implement Disabled Provider" \
  --label "enhancement,v0.2.0,doi" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement \"disabled\" DOI provider for when minting is off.

**Tasks**:
- [ ] Create \`internal/doi/disabled.go\`
- [ ] Implement \`DisabledProvider\` struct
- [ ] Return clear error messages
- [ ] Write unit tests
- [ ] Document usage

**Acceptance Criteria**:
- Returns helpful error when disabled
- Tests verify error messages
- Documentation explains configuration

**Estimate**: 1 day
**Milestone 5** (Weeks 9-10)"

gh issue create \
  --title "Implement DataCite Provider (Stub)" \
  --label "enhancement,v0.2.0,doi,datacite" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement DataCite provider with basic functionality.

**Tasks**:
- [ ] Create \`internal/doi/datacite.go\`
- [ ] Implement \`DataCiteProvider\` struct
- [ ] Add DataCite API client
- [ ] Implement \`Mint()\` method
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

**Estimate**: 5 days
**Milestone 5** (Weeks 9-10)"

gh issue create \
  --title "DOI Configuration System" \
  --label "enhancement,v0.2.0,doi,config" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Add DOI configuration to Cicada config system.

**Tasks**:
- [ ] Update \`internal/config/config.go\`
- [ ] Add \`DOI\` configuration section
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

**Estimate**: 2 days
**Milestone 5** (Weeks 9-10)"

gh issue create \
  --title "DOI CLI Commands" \
  --label "enhancement,v0.2.0,cli,doi" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Add CLI commands for DOI management.

**Tasks**:
- [ ] Create \`cmd/cicada/publish.go\`
- [ ] Add \`cicada publish <path>\` command
- [ ] Add interactive metadata prompts
- [ ] Add \`--dry-run\` flag
- [ ] Add \`cicada publish list\` command
- [ ] Add \`cicada publish show <doi>\` command
- [ ] Write CLI tests
- [ ] Update documentation

**Acceptance Criteria**:
- Can mint DOI via CLI
- Interactive prompts are clear
- Dry-run shows preview
- Tests verify CLI behavior

**Estimate**: 3 days
**Milestone 5** (Weeks 9-10)"

# Milestone 6: Multi-User Foundation

gh issue create \
  --title "IAM User Creation Commands" \
  --label "enhancement,v0.2.0,auth,iam" \
  --milestone "Phase 3: Web UI & User Management" \
  --body "Add CLI commands for creating IAM users.

**Tasks**:
- [ ] Create \`cmd/cicada/user.go\`
- [ ] Add \`cicada user add <name>\` command
- [ ] Create IAM user via AWS SDK
- [ ] Generate access keys
- [ ] Store credentials securely
- [ ] Add \`cicada user list\` command
- [ ] Add \`cicada user remove <name>\` command
- [ ] Write CLI tests
- [ ] Document user management

**Acceptance Criteria**:
- Creates IAM users successfully
- Credentials stored securely
- Tests verify user creation
- Documentation includes IAM prerequisites

**Estimate**: 3 days
**Milestone 6** (Weeks 11-12)"

gh issue create \
  --title "Path-Based IAM Policy Generation" \
  --label "enhancement,v0.2.0,auth,iam" \
  --milestone "Phase 3: Web UI & User Management" \
  --body "Implement automatic IAM policy generation for path-based access.

**Tasks**:
- [ ] Create \`internal/auth/policy.go\`
- [ ] Implement policy template system
- [ ] Generate policies for S3 path prefixes
- [ ] Support read-only vs read-write access
- [ ] Attach policies to users
- [ ] Add \`cicada user grant <user> <path>\` command
- [ ] Add \`cicada user revoke <user> <path>\` command
- [ ] Write unit tests
- [ ] Document policy generation

**Acceptance Criteria**:
- Policies follow least-privilege principle
- Path-based access works correctly
- Tests verify policy generation
- Documentation explains access patterns

**Estimate**: 4 days
**Milestone 6** (Weeks 11-12)"

gh issue create \
  --title "Project Management System" \
  --label "enhancement,v0.2.0,projects" \
  --milestone "Phase 3: Web UI & User Management" \
  --body "Add basic project management for organizing data.

**Tasks**:
- [ ] Create \`internal/config/project.go\`
- [ ] Define \`Project\` struct
- [ ] Add \`cicada project create <name>\` command
- [ ] Add \`cicada project list\` command
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

**Estimate**: 3 days
**Milestone 6** (Weeks 11-12)"

gh issue create \
  --title "Multi-User Documentation" \
  --label "documentation,v0.2.0,auth" \
  --milestone "Phase 3: Web UI & User Management" \
  --body "Create comprehensive multi-user setup guide.

**Tasks**:
- [ ] Create \`docs/MULTI_USER_SETUP.md\`
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

**Estimate**: 2 days
**Milestone 6** (Weeks 11-12)"

# Milestone 7: Documentation & Release

gh issue create \
  --title "Update USER_SCENARIOS for v0.2.0" \
  --label "documentation,v0.2.0" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Update USER_SCENARIOS.md with v0.2.0 features.

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

**Estimate**: 2 days
**Milestone 7** (Weeks 13-14)"

gh issue create \
  --title "Create v0.2.0 Migration Guide" \
  --label "documentation,v0.2.0" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Create migration guide for v0.1.0 → v0.2.0.

**Tasks**:
- [ ] Create \`docs/MIGRATION_v0.2.0.md\`
- [ ] Document breaking changes (if any)
- [ ] Document new configuration options
- [ ] Provide upgrade steps
- [ ] Include rollback instructions

**Acceptance Criteria**:
- Guide is clear and complete
- All breaking changes documented
- Upgrade path is straightforward

**Estimate**: 1 day
**Milestone 7** (Weeks 13-14)"

gh issue create \
  --title "Integration Testing Suite" \
  --label "testing,v0.2.0" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Create comprehensive integration test suite.

**Tasks**:
- [ ] Create \`internal/integration/metadata_test.go\`
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

**Estimate**: 4 days
**Milestone 7** (Weeks 13-14)"

gh issue create \
  --title "Performance & Load Testing" \
  --label "testing,v0.2.0,performance" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Validate performance targets for v0.2.0.

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

**Estimate**: 2 days
**Milestone 7** (Weeks 13-14)"

gh issue create \
  --title "Release Preparation" \
  --label "release,v0.2.0" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Prepare v0.2.0 for release.

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

**Estimate**: 2 days
**Milestone 7** (Weeks 13-14)"

gh issue create \
  --title "v0.2.0 Announcement & Outreach" \
  --label "release,v0.2.0,community" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Announce v0.2.0 release to community.

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

**Estimate**: 1 day
**Milestone 7** (Weeks 13-14)"

echo ""
echo "✅ All 32 v0.2.0 issues created successfully!"
echo ""
echo "View issues: gh issue list --milestone 'Phase 2: Metadata & FAIR'"
echo "View milestones: gh milestone list"
