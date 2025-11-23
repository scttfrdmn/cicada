# Cicada Specification Index

Complete index of all specification files for the Cicada project.

## Quick Navigation

### 📚 Documentation
- [QUICKSTART.md](QUICKSTART.md) - **Start here!** Developer quick start guide
- [README.md](README.md) - Project overview and architecture
- [ROADMAP.md](ROADMAP.md) - 26-week development plan with detailed specs

### 📖 User Documentation
- [docs/cli-reference.md](docs/cli-reference.md) - Complete CLI command reference
- [docs/domain-schemas.md](docs/domain-schemas.md) - Domain-specific metadata schemas

### 💻 Code Specifications
- [internal/metadata/schema.go](internal/metadata/schema.go) - Metadata schema system (Go)
- [internal/metadata/extractor.go](internal/metadata/extractor.go) - File metadata extractors (Go)
- [internal/doi/datacite.go](internal/doi/datacite.go) - DOI minting with DataCite (Go)

### ⚙️ Configuration Examples
- [examples/configs.md](examples/configs.md) - All configuration file examples

## Document Summaries

### QUICKSTART.md (1,900 lines)
**Purpose**: Fastest path to understanding and building Cicada
**Contains**:
- What's inside this spec
- Key design decisions
- Domain examples (12 research fields)
- Cost model
- Key differentiators
- Target user persona
- Implementation tips
- Success metrics

**Use when**: Starting development, orienting new contributors

---

### README.md (700 lines)
**Purpose**: Project overview and technical architecture  
**Contains**:
- Overview and key principles
- Target users and core features
- System architecture diagram
- Complete project structure
- Technology stack
- Development phases
- Getting started for developers

**Use when**: Understanding overall system design

---

### ROADMAP.md (1,800 lines)
**Purpose**: Detailed 26-week implementation plan  
**Contains**:
- 6 development phases with weekly breakdown
- Technical specifications for each component
- Dependencies and test requirements
- Acceptance criteria
- Milestone deliverables
- Risk mitigation strategies
- Post-launch roadmap

**Use when**: Planning sprints, implementing specific features

---

### docs/cli-reference.md (1,200 lines)
**Purpose**: Complete command-line interface reference  
**Contains**:
- Installation and setup
- All CLI commands with examples
- Data sync commands
- Metadata management
- Workflow execution
- Workstation management
- DOI minting
- User/project management
- Cost tracking
- Compliance features

**Use when**: Implementing CLI, writing user docs, testing

---

### docs/domain-schemas.md (500 lines)
**Purpose**: Reference for domain-specific metadata  
**Contains**:
- 12+ research domain examples:
  - Microscopy & Imaging
  - Genomics & Sequencing
  - Proteomics & Mass Spec
  - Flow Cytometry
  - Chromatography
  - Spectroscopy
  - X-ray Crystallography
  - Electron Microscopy
  - Behavioral Studies
  - Clinical Trials
  - Environmental Sampling
  - Materials Science

**Use when**: Implementing metadata schemas, understanding user needs

---

### internal/metadata/schema.go (700 lines)
**Purpose**: Core metadata system implementation  
**Contains**:
- Schema data structures
- Field types and validation
- Schema manager
- Validation engine
- Quality scoring
- Helper functions

**Use when**: Implementing Phase 2 (metadata system)

---

### internal/metadata/extractor.go (900 lines)
**Purpose**: Automatic metadata extraction from files  
**Contains**:
- Extractor interface
- Registry pattern
- Format-specific extractors:
  - TIFF, OME-TIFF
  - Zeiss CZI, Nikon ND2, Leica LIF
  - FASTQ, BAM
  - mzML, MGF (mass spec)
  - HDF5, Zarr
  - DICOM
  - FCS (flow cytometry)
  - Generic fallback

**Use when**: Implementing Phase 2 (metadata extraction)

---

### internal/doi/datacite.go (800 lines)
**Purpose**: DOI minting and DataCite integration  
**Contains**:
- DOIManager implementation
- DataCite API client
- DataCite Metadata Schema 4.4 structures
- XML generation
- DOI lifecycle management

**Use when**: Implementing Phase 5 (DOI system)

---

### examples/configs.md (800 lines)
**Purpose**: Configuration file templates  
**Contains**:
- Main config (config.yaml)
- Workflow configs
- Auto-pipeline configs
- Metadata schema examples
- Watch configurations
- Project configs
- DOI configs
- Compliance configs (NIST 800-171)
- User preferences
- Environment variables
- Deployment configs (Docker, systemd)

**Use when**: Setting up Cicada, writing config parsers

---

## Development Workflow

### Phase 1: Core (Weeks 1-6)
**Read**: README.md (architecture), ROADMAP.md (Phase 1)  
**Implement**: Sync engine, file watching, CLI, daemon  
**Reference**: CLI examples for command design

### Phase 2: Metadata (Weeks 7-10)
**Read**: internal/metadata/*.go, domain-schemas.md  
**Implement**: Schema system, extractors, validation  
**Reference**: Configuration examples for schema format

### Phase 3: Web UI (Weeks 11-14)
**Read**: README.md (web architecture), CLI reference  
**Implement**: Web server, frontend, user management  
**Reference**: All docs for API design

### Phase 4: Compute (Weeks 15-18)
**Read**: ROADMAP.md (Phase 4), examples/configs.md  
**Implement**: AWS Batch, workflow engines  
**Reference**: Workflow config examples

### Phase 5: Portal (Weeks 19-22)
**Read**: internal/doi/datacite.go, CLI reference  
**Implement**: Workstations, portal, DOI minting  
**Reference**: DOI config examples

### Phase 6: Polish (Weeks 23-26)
**Read**: Compliance configs, QUICKSTART.md  
**Implement**: Compliance modes, documentation, tests  
**Reference**: All docs for completeness

---

## File Statistics

```
Total Files: 8
Total Lines: ~8,300

Documentation:  4,900 lines (59%)
Code Specs:     2,400 lines (29%)
Examples:       1,000 lines (12%)
```

---

## Implementation Workflow

### Development Phases

**Phase 1 (Weeks 1-6)**: Core sync
- Reference: README.md, ROADMAP.md Phase 1
- Implement: Sync engine, file watching, CLI, daemon

**Phase 2 (Weeks 7-10)**: Metadata
- Reference: internal/metadata/*.go, domain-schemas.md
- Implement: Schema system, extractors, validation

**Phase 3 (Weeks 11-14)**: Web UI
- Reference: README.md, cli-reference.md
- Implement: Web server, frontend, user management

**Phase 4+**: Continue following ROADMAP.md phases

### Component Implementation Guide

**For each component**:
1. Read QUICKSTART.md for overview
2. Review relevant ROADMAP.md section for current phase
3. Consult specific .go specification files for implementation details
4. Implement with test-driven development
5. Validate against acceptance criteria

---

## Key Design Philosophy

🦗 **Dormant**: Resources only when needed  
💰 **Cheap**: $50-100/month typical cost  
🎯 **Simple**: Non-technical users  
📊 **FAIR**: Research data standards  
🔒 **Compliant**: NIST 800-171, HIPAA ready  
🌐 **Open**: Open source, community-driven  

---

## Next Steps

1. Read **QUICKSTART.md**
2. Review **README.md** for architecture
3. Start **Phase 1** from ROADMAP.md
4. Reference **cli-reference.md** for UX guidance
5. Use **examples/configs.md** for file formats

---

**Project**: Cicada Data Commons  
**Version**: 1.0 (Specifications)  
**Date**: 2024-11-22  
**Status**: Pre-implementation

Good luck building! 🦗
