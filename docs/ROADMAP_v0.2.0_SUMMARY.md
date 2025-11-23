# Cicada v0.2.0 - Implementation Plan Summary

**Status**: Planning Complete ✅
**Target Release**: Q1 2026 (14 weeks)
**GitHub Milestone**: [Phase 2: Metadata & FAIR](https://github.com/scttfrdmn/cicada/milestone/2)

---

## Overview

Version 0.2.0 builds on v0.1.0's foundational storage and sync layer by adding **metadata intelligence and instrument awareness** - critical capabilities for Cicada's vision as a dormant data commons platform.

### Key Additions

1. **Metadata Extraction Framework**: Automatically extract and preserve instrument metadata
2. **Instrument Presets**: Simplified configuration for common lab instruments
3. **DOI Provider System**: Pluggable DOI minting (DataCite/Zenodo)
4. **Basic Multi-User Support**: IAM automation for 2-3 users

---

## Milestones & Timeline

### Milestone 1: Metadata Foundation (Weeks 1-2)
**Goal**: Core metadata infrastructure in place

**Issues**:
- [#17](https://github.com/scttfrdmn/cicada/issues/17) Define Metadata Core Types
- [#18](https://github.com/scttfrdmn/cicada/issues/18) Implement Extractor Interface
- [#19](https://github.com/scttfrdmn/cicada/issues/19) S3 Object Tagging Integration

**Deliverable**: Metadata types, extractor interface, S3 tag storage

---

### Milestone 2: First Extractor Working (Weeks 3-4)
**Goal**: Zeiss CZI extraction functional end-to-end

**Issues**:
- [#20](https://github.com/scttfrdmn/cicada/issues/20) Implement Zeiss CZI Extractor
- [#21](https://github.com/scttfrdmn/cicada/issues/21) CLI Integration for Metadata Extraction

**Deliverable**: Working CZI extractor with CLI commands

---

### Milestone 3: Extraction Library Complete (Weeks 5-6)
**Goal**: OME-TIFF and FASTQ extractors working

**Tasks** (to be broken into issues):
- Implement OME-TIFF extractor
- Implement FASTQ extractor
- Performance benchmarking

**Deliverable**: 3 production extractors (CZI, OME-TIFF, FASTQ)

---

### Milestone 4: Instrument Presets (Weeks 7-8)
**Goal**: Preset system with 5+ instruments

**Issues**:
- [#22](https://github.com/scttfrdmn/cicada/issues/22) Instrument Preset System

**Deliverable**: Auto-detection, preset library, CLI commands

---

### Milestone 5: DOI Provider System (Weeks 9-10)
**Goal**: DataCite provider functional

**Issues**:
- [#23](https://github.com/scttfrdmn/cicada/issues/23) Pluggable DOI Provider System

**Deliverable**: DOI minting with DataCite test mode

---

### Milestone 6: Multi-User Foundation (Weeks 11-12)
**Goal**: Basic IAM automation for 2-3 users

**Issues**:
- [#24](https://github.com/scttfrdmn/cicada/issues/24) Basic Multi-User Support

**Deliverable**: IAM user creation, path-based access control

---

### Milestone 7: Documentation & Release (Weeks 13-14)
**Goal**: v0.2.0 released with complete docs

**Issues**:
- [#25](https://github.com/scttfrdmn/cicada/issues/25) Documentation & Testing

**Deliverable**: Updated docs, migration guide, v0.2.0 release

---

## Success Metrics

### Technical Metrics
- **Metadata extraction**: < 5s for typical files (< 1GB)
- **Sync overhead**: < 10% performance impact with metadata enabled
- **Validation accuracy**: 99%+ detection of corrupt files
- **Preset accuracy**: 95%+ correct auto-detection

### User Metrics
- **Setup time reduction**: 5 min → 1 min with presets
- **Metadata completeness**: 80%+ of required fields extracted
- **DOI minting success**: 99%+ successful mints
- **User satisfaction**: Positive feedback from beta users

---

## Documentation

### Planning Documents
- **[ROADMAP_v0.2.0.md](ROADMAP_v0.2.0.md)**: Feature specifications
- **[ROADMAP_v0.2.0_PLAN.md](ROADMAP_v0.2.0_PLAN.md)**: Detailed 32-issue breakdown
- **This Document**: Executive summary

### Reference Code
Already implemented (POC/framework):
- `internal/metadata/extractor.go` - Extractor POC with Zeiss CZI basic implementation
- `internal/doi/provider.go` - Provider interface
- `internal/doi/providers_example.go` - DataCite/Zenodo stubs
- `presets/` - Example preset YAML files

---

## Project Management

### GitHub Organization

**Milestones**:
- Primary: [Phase 2: Metadata & FAIR](https://github.com/scttfrdmn/cicada/milestone/2)
- Multi-user: [Phase 3: Web UI & User Management](https://github.com/scttfrdmn/cicada/milestone/3)

**Labels**:
- `v0.2.0` - Version 0.2.0 features
- `type: feature` - New features
- `area: metadata` - Metadata system
- `area: cli` - CLI commands
- `priority: critical/high/medium/low` - Priority levels

**Issues**: 9 consolidated issues (see milestone)

### GitHub Projects

Track v0.2.0 development in [GitHub Projects](https://github.com/scttfrdmn/cicada/projects).

**Suggested Views**:
1. **By Milestone**: Group issues by milestone 1-7
2. **By Priority**: Critical → High → Medium → Low
3. **Kanban**: Backlog → In Progress → Review → Done

---

## Dependencies & Prerequisites

### Technical Dependencies
- Go 1.23+
- AWS SDK for Go v2
- CZI parsing library (to be evaluated in #20)
- XML parsing for OME-TIFF

### AWS Permissions
Existing v0.1.0 permissions plus:
- `s3:PutObjectTagging` - For metadata tags
- `s3:GetObjectTagging` - For reading metadata
- `iam:CreateUser` - For multi-user (optional)
- `iam:AttachUserPolicy` - For multi-user (optional)

---

## Risk Mitigation

### Risk 1: CZI Parser Complexity
**Mitigation**: Use battle-tested bioformats library or collaborate with OME team

### Risk 2: S3 Tagging Limitations (10 tags max)
**Mitigation**: Hybrid approach - critical fields in tags, full metadata in sidecar JSON

### Risk 3: DOI Provider API Changes
**Mitigation**: Abstracted provider interface, versioned API clients

### Risk 4: Performance Impact
**Mitigation**: Optional metadata extraction, async processing, caching

### Risk 5: Scope Creep
**Mitigation**: Strict milestone-based development, defer features to v0.3.0

---

## Communication & Updates

### Weekly Updates
Post progress updates to:
- GitHub Discussions
- Project README
- Community channels

### Blockers
Report blockers immediately in:
- GitHub issue comments
- Project standup
- Slack/Discord (if applicable)

---

## Next Steps

1. ✅ **Planning Complete** - Documentation and issues created
2. 🔜 **Milestone 1** - Start with issue #17 (Define Metadata Core Types)
3. 🔜 **Team Assignment** - Assign issues to developers
4. 🔜 **Sprint Planning** - Set up 2-week sprints aligned with milestones

---

## Resources

### Planning Documents
- Full specification: [planning/PROJECT-SUMMARY.md](../planning/PROJECT-SUMMARY.md)
- Comprehensive roadmap: [planning/ROADMAP.md](../planning/ROADMAP.md)
- User scenarios: [USER_SCENARIOS.md](USER_SCENARIOS.md)

### Community
- GitHub: https://github.com/scttfrdmn/cicada
- Issues: https://github.com/scttfrdmn/cicada/issues
- Milestones: https://github.com/scttfrdmn/cicada/milestones

---

**Last Updated**: 2025-11-23
**Status**: Ready to begin development ✅
