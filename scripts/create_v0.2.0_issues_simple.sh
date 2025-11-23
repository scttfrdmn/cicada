#!/bin/bash

# Simplified script to create v0.2.0 issues using existing labels
# Usage: ./scripts/create_v0.2.0_issues_simple.sh

set -e

echo "Creating v0.2.0 GitHub issues..."
echo ""

# Create a few key issues to get started

echo "Creating Issue #1: Define Metadata Core Types..."
gh issue create \
  --title "[v0.2.0] Define Metadata Core Types" \
  --label "type: feature,area: metadata,priority: critical" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Define core metadata data structures in Go.

**Estimate**: 2 days | **Milestone**: 1 (Weeks 1-2)

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
- Documentation generated with godoc"

echo "Creating Issue #2: Implement Extractor Interface..."
gh issue create \
  --title "[v0.2.0] Implement Extractor Interface" \
  --label "type: feature,area: metadata,priority: critical" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Define and implement the pluggable extractor interface.

**Estimate**: 2 days | **Milestone**: 1 (Weeks 1-2)

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
- Tests cover all registry operations"

echo "Creating Issue #3: S3 Object Tagging Integration..."
gh issue create \
  --title "[v0.2.0] S3 Object Tagging Integration" \
  --label "type: feature,area: metadata,priority: critical" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Add support for storing metadata as S3 object tags.

**Estimate**: 3 days | **Milestone**: 1 (Weeks 1-2)

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
- Documentation includes IAM policy examples"

echo "Creating Issue #4: Implement Zeiss CZI Extractor..."
gh issue create \
  --title "[v0.2.0] Implement Zeiss CZI Extractor" \
  --label "type: feature,area: metadata,priority: high" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement full Zeiss CZI metadata extractor.

**Estimate**: 5 days | **Milestone**: 2 (Weeks 3-4)

**Tasks**:
- [ ] Research CZI file format and libraries
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

**Reference**: See \`internal/metadata/extractor.go\` (already has POC)"

echo "Creating Issue #5: CLI Integration for Metadata..."
gh issue create \
  --title "[v0.2.0] CLI Integration for Metadata Extraction" \
  --label "type: feature,area: cli,area: metadata,priority: high" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Add CLI commands for metadata extraction.

**Estimate**: 3 days | **Milestone**: 2 (Weeks 3-4)

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
- Tests verify CLI behavior"

echo "Creating Issue #6: Instrument Preset System..."
gh issue create \
  --title "[v0.2.0] Instrument Preset System" \
  --label "type: feature,priority: high" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement instrument preset system for simplified configuration.

**Estimate**: 8 days | **Milestone**: 4 (Weeks 7-8)

**Tasks**:
- [ ] Design preset YAML format
- [ ] Implement preset loader and registry
- [ ] Implement auto-detection logic
- [ ] Create initial presets (Zeiss, Illumina, etc.)
- [ ] Add CLI commands (\`cicada instrument list/show/detect/setup\`)
- [ ] Add interactive setup wizard
- [ ] Write tests
- [ ] Document preset usage

**Acceptance Criteria**:
- 6+ presets created and tested
- 95%+ detection accuracy
- Interactive wizard is user-friendly
- Documentation complete

**Reference**: See \`presets/\` directory for existing preset examples"

echo "Creating Issue #7: DOI Provider System..."
gh issue create \
  --title "[v0.2.0] Pluggable DOI Provider System" \
  --label "type: feature,priority: high" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Implement pluggable DOI provider system (DataCite/Zenodo/Off).

**Estimate**: 10 days | **Milestone**: 5 (Weeks 9-10)

**Tasks**:
- [ ] Define \`Provider\` interface
- [ ] Implement \`DisabledProvider\`
- [ ] Implement \`DataCiteProvider\` (stub with test mode)
- [ ] Add DOI configuration system
- [ ] Add CLI commands (\`cicada publish\`)
- [ ] Write tests
- [ ] Document DOI setup

**Acceptance Criteria**:
- Can mint DOI in test mode
- Config supports provider selection
- CLI interactive prompts work
- Documentation includes credential setup

**Reference**: See \`internal/doi/provider.go\` and \`internal/doi/providers_example.go\`"

echo "Creating Issue #8: Basic Multi-User Support..."
gh issue create \
  --title "[v0.2.0] Basic Multi-User Support" \
  --label "type: feature,priority: medium" \
  --milestone "Phase 3: Web UI & User Management" \
  --body "Add basic multi-user support for 2-3 lab members.

**Estimate**: 10 days | **Milestone**: 6 (Weeks 11-12)

**Tasks**:
- [ ] Add IAM user creation commands (\`cicada user add/list/remove\`)
- [ ] Implement path-based IAM policy generation
- [ ] Add access grant/revoke commands
- [ ] Add basic project management
- [ ] Write tests
- [ ] Create multi-user setup documentation

**Acceptance Criteria**:
- Can create IAM users via CLI
- Path-based access works correctly
- Documentation enables self-service setup
- Tests verify user management

**Note**: This provides minimal multi-user capability for v0.2.0. Full web UI and advanced collaboration features planned for v0.3.0."

echo "Creating Issue #9: v0.2.0 Documentation & Testing..."
gh issue create \
  --title "[v0.2.0] Documentation & Testing" \
  --label "type: docs,type: test,priority: high" \
  --milestone "Phase 2: Metadata & FAIR" \
  --body "Complete documentation and testing for v0.2.0 release.

**Estimate**: 8 days | **Milestone**: 7 (Weeks 13-14)

**Tasks**:
- [ ] Update USER_SCENARIOS.md with v0.2.0 features
- [ ] Create MIGRATION_v0.2.0.md guide
- [ ] Create MULTI_USER_SETUP.md guide
- [ ] Create integration test suite
- [ ] Run performance & load testing
- [ ] Update all documentation
- [ ] Prepare release

**Acceptance Criteria**:
- All documentation updated
- Integration tests pass
- Performance targets met
- Migration guide complete"

echo ""
echo "✅ Created 9 consolidated v0.2.0 issues!"
echo ""
echo "These issues consolidate the 32 detailed tasks into manageable chunks."
echo ""
echo "Next steps:"
echo "  1. View issues: gh issue list --milestone 'Phase 2: Metadata & FAIR'"
echo "  2. Break down issues into smaller tasks as needed"
echo "  3. Assign issues to team members"
echo "  4. Track progress in GitHub Projects"
