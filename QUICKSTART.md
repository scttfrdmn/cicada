# Cicada: Developer Quick Start

This directory contains comprehensive specifications for building Cicada, a dormant data commons platform for academic research labs.

## What's Inside

### Core Documentation

📄 **README.md** - Project overview, architecture, and getting started  
📄 **ROADMAP.md** - 26-week development plan with detailed specifications  
📄 **docs/cli-reference.md** - Complete CLI command examples  
📄 **docs/domain-schemas.md** - Domain-specific metadata schemas

### Code Specifications

📄 **internal/metadata/schema.go** - Metadata schema system (Go)  
📄 **internal/metadata/extractor.go** - File metadata extractors (Go)  
📄 **internal/doi/datacite.go** - DOI minting with DataCite (Go)

## Development Approach

These specifications provide comprehensive implementation guidance:

### Getting Started

```bash
# In your development environment
cd ~/projects
mkdir cicada
cd cicada

# Copy these specifications
cp -r /path/to/cicada-spec/* .

# Begin implementation following the roadmap phases
# Start with Phase 1: Core Storage & Sync
```

### Development Workflow

1. **Phase 1**: Start with `ROADMAP.md` Phase 1
   - Implement sync engine
   - Build CLI
   - Set up daemon

2. **Phase 2**: Add metadata system
   - Reference `internal/metadata/schema.go`
   - Implement extractors from `internal/metadata/extractor.go`
   - Add domain schemas from `docs/domain-schemas.md`

3. **Phase 3+**: Continue following roadmap phases

### Key Design Decisions

#### Why Go?
- Single binary deployment (critical for non-technical users)
- Excellent AWS SDK support
- Fast, low memory footprint
- Great concurrency for file operations
- Cross-platform compilation

#### Why This Architecture?
- **Local daemon**: Works offline, low latency
- **Web UI served locally**: No external dependencies
- **S3 as backend**: Durable, cheap, scalable
- **Metadata in files**: Portable, no database lock-in

#### Why These Features?
- **Dormant design**: Academic budgets are tight
- **FAIR by default**: Grant requirements
- **Domain flexibility**: Each field has unique needs
- **DOI integration**: Publication requirements
- **Compliance modes**: Clinical/sensitive data

## Domain Examples

The system supports diverse research domains out of the box:

### Imaging Sciences
- **Microscopy**: Fluorescence, confocal, super-resolution
- **Medical Imaging**: CT, MRI, PET (via DICOM)
- **Electron Microscopy**: TEM, SEM

### Omics
- **Genomics**: RNA-seq, DNA-seq, ChIP-seq
- **Proteomics**: Mass spectrometry, LC-MS
- **Metabolomics**: NMR, GC-MS

### Analytical Chemistry
- **Chromatography**: HPLC, UHPLC, GC
- **Spectroscopy**: NMR, IR, UV-Vis, Raman
- **X-ray Crystallography**: Diffraction data

### Cell Biology
- **Flow Cytometry**: Multi-parameter analysis
- **Cell Culture**: Live-cell imaging, high-content screening

### Other Sciences
- **Behavioral Studies**: Video, tracking data, psychometrics
- **Environmental Science**: Sample collection, field data
- **Materials Science**: AFM, SEM, mechanical testing
- **Clinical Trials**: Patient data (HIPAA-compliant)

Each domain has:
1. Metadata schema (YAML)
2. File format extractors
3. Validation rules
4. Search facets
5. Export templates

## Cost Model

Designed for **$50-100/month** budgets:

### Typical Lab (10TB data)
- **Storage**: ~$80/month (Intelligent-Tiering)
- **Compute**: ~$10/month (spot instances, bursty)
- **Transfer**: ~$5/month (minimal egress)
- **Total**: ~$95/month

### What Makes It Cheap?
- Intelligent-Tiering (automatic cost optimization)
- Spot instances (70% cheaper compute)
- On-demand resources (only pay when using)
- No always-on infrastructure

## Key Differentiators

### vs. Dropbox/Google Drive
✅ Much cheaper for large datasets  
✅ Compute integration  
✅ Metadata management  
✅ DOI minting  
✅ Compliance modes  

### vs. Traditional Data Commons (iRODS, Dataverse)
✅ No IT staff required  
✅ Self-service setup  
✅ Cloud-native (scales automatically)  
✅ 10x cheaper to operate  
✅ Better UX for non-technical users  

### vs. rclone/AWS CLI
✅ Metadata management  
✅ User-friendly GUI  
✅ Workflow integration  
✅ Cost management  
✅ Collaboration features  

## Target Users: Persona

**Dr. Sarah Chen** - Assistant Professor, Cell Biology
- Lab of 8 people (1 postdoc, 5 grad students, 2 undergrads)
- Generates ~500GB/month (microscopy, some sequencing)
- Budget: $50-75/month for cloud storage
- Technical skill: Can use command line, but prefers GUI
- Pain points:
  - Data on dying hard drives
  - No organized backup system
  - Students leaving with data
  - Can't share data easily with collaborators
  - Journal requires data availability statements

**What Cicada Solves**:
- Automatic backup ($80/month, cheaper than replacing drives)
- Organized metadata (find experiments from 2 years ago)
- Team collaboration (everyone has access)
- Easy sharing (generate DOI, public link)
- Grant compliance (FAIR principles, data sharing)

## Implementation Tips

### Start Simple
Don't implement everything at once. Follow the roadmap:
1. Get sync working ✓
2. Add watches ✓
3. Build basic CLI ✓
4. Add metadata gradually

### Test with Real Users Early
Find a friendly lab to beta test after Phase 2 (week 10).

### Leverage Existing Tools
- **Bio-Formats**: For microscopy formats
- **HTSlib**: For sequencing formats
- **OpenMS**: For mass spec formats
- Don't reinvent parsers

### Focus on UX
The #1 reason academic tools fail is poor UX. Make it:
- **Fast**: No loading spinners >2 seconds
- **Clear**: Every action has obvious feedback
- **Forgiving**: Easy to undo mistakes
- **Helpful**: Good error messages

### Community-Driven
Open source from day 1. Let users contribute:
- Domain schemas
- Instrument adapters
- Workflow templates
- Documentation

## Development Environment Setup

### Prerequisites
```bash
# Go 1.21+
go version

# Node.js 20+ (for web UI)
node --version

# AWS CLI configured
aws sts get-caller-identity

# Docker (for workflow testing)
docker --version
```

### Quick Start
```bash
# Clone (once it exists)
git clone https://github.com/your-org/cicada.git
cd cicada

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run locally
./bin/cicada daemon start --dev

# Build web UI
cd web && npm install && npm run dev
```

### Development Best Practices

**Implementation Approach**:
1. Review relevant spec files before implementing each component
2. Understand the design rationale documented in specifications
3. Write tests alongside implementation (TDD approach)
4. Generate example data and configurations for testing
5. Validate implementation against acceptance criteria in ROADMAP.md

**Key Resources**:
- Implementation guidance: `internal/sync/`, `internal/metadata/`, etc.
- Test patterns: Documented in ROADMAP.md for each phase
- Example schemas: See `docs/domain-schemas.md` for microscopy, sequencing, etc.
- Acceptance criteria: Each phase in ROADMAP.md includes specific requirements

## Questions to Consider

### Before Starting
- [ ] Who is your target user? (Match persona)
- [ ] What domain(s) will you focus on first?
- [ ] Self-hosted or SaaS model?
- [ ] Open core or fully open source?

### During Development
- [ ] Is this feature essential for MVP?
- [ ] Does this match the "dormant" design philosophy?
- [ ] Would a non-technical user understand this?
- [ ] What's the cost impact?

### Before Launch
- [ ] Have you tested with real labs?
- [ ] Is documentation complete?
- [ ] Can someone set up in <10 minutes?
- [ ] Is it actually cheaper than alternatives?

## Success Metrics

### Technical
- Setup time: <10 minutes
- Sync speed: 10GB in <15 minutes
- Uptime: >99.9%
- Test coverage: >80%

### User
- Weekly active users: Growing
- Data under management: >1PB (across all users)
- Support ticket rate: <5%
- NPS: >50

### Business
- Cost per user: <$5/month (your overhead)
- Revenue per user: $10-20/month (if SaaS)
- User growth: 20% month-over-month
- Retention: >90% after 3 months

## Resources

### Documentation to Write
- Getting Started Guide (screencast)
- Domain-specific tutorials
- API documentation
- Deployment guide (for institutions)
- Contribution guide

### Community Building
- GitHub Discussions
- Slack/Discord
- Monthly webinars
- Case studies from early users
- Academic paper about the platform

### Marketing
- Website (with live demo)
- Blog posts
- Conference presentations (FORCE11, RDA)
- Twitter/Mastodon presence
- Testimonials from PIs

## Contact & Contributing

(Fill in once project is public)

- GitHub: https://github.com/your-org/cicada
- Docs: https://cicada.sh/docs
- Community: https://cicada.sh/community
- Twitter: @cicada_data

## License

Recommend: **MIT License** (maximum compatibility)

Alternative: **Apache 2.0** (if you want patent protection)

---

## Next Steps

1. **Set up project structure** following README.md
2. **Start Phase 1** implementation (ROADMAP.md)
3. **Follow test-driven development** practices
4. **Test early and often** with real users
5. **Document as you go** - don't wait until the end

Good luck building Cicada! 🦗

---

**Last Updated**: 2024-11-22  
**Version**: 1.0 (Pre-release specifications)
