# Cicada v0.3.0 Summary

**Target Release:** Q2 2025 (End of May 2025)
**Current Status:** v0.2.0 Released (January 23, 2025)
**Planning Complete:** ✅ Ready to start development

## Quick Overview

v0.3.0 transforms Cicada from a metadata preparation tool into a **complete DOI registration platform** with live provider integration and expanded file format support.

### Core Features

| Feature | Description | Priority |
|---------|-------------|----------|
| **Live DOI Minting** | DataCite & Zenodo API integration | ⭐⭐⭐ CRITICAL |
| **Microscopy Support** | CZI & OME-TIFF metadata extraction | ⭐⭐ HIGH |
| **Custom Presets** | User-defined validation rules | ⭐⭐ HIGH |
| **Interactive Editing** | TUI for metadata review | ⭐ MEDIUM |
| **Bioinformatics** | BAM/SAM & VCF support | ⭐ MEDIUM |

## What's New in v0.3.0

### 1. Complete DOI Workflow

**Before (v0.2.0):**
```bash
# Prepare metadata
cicada doi prepare sample.fastq --enrich metadata.yaml --output doi-ready.json
# ❌ Stop here - no way to actually mint DOI
```

**After (v0.3.0):**
```bash
# Prepare and publish in one step
cicada doi publish sample.fastq --enrich metadata.yaml --provider zenodo
# ✅ DOI: 10.5281/zenodo.123456
# ✅ URL: https://zenodo.org/record/123456
```

### 2. Microscopy File Support

**Before (v0.2.0):**
```bash
cicada metadata extract image.czi
# ❌ Error: Unsupported format
```

**After (v0.3.0):**
```bash
cicada metadata extract image.czi --preset zeiss-lsm-880
# ✅ Extracted: 5 channels, 20 Z-slices, 1024x1024 px
# ✅ Quality Score: 85/100
```

### 3. Custom Presets

**Before (v0.2.0):**
```bash
# Only 8 built-in presets
cicada metadata preset list
# ❌ Can't create lab-specific presets
```

**After (v0.3.0):**
```bash
# Create custom preset
cicada metadata preset create lab-rnaseq \
  --template generic-sequencing \
  --add-required pipeline_version

# Share with team
cicada metadata preset export lab-rnaseq > lab-preset.yaml
```

### 4. Interactive Metadata Editing

**Before (v0.2.0):**
```bash
# Manual YAML editing
vim metadata.yaml  # Edit in text editor
cicada doi prepare sample.fastq --enrich metadata.yaml
```

**After (v0.3.0):**
```bash
# Interactive TUI
cicada metadata edit sample.fastq
# ✅ Visual form with validation
# ✅ Real-time quality score
# ✅ Field help and examples
```

## Timeline

### Conservative (5-6 months → May 2025)

```
Jan 2025    Feb 2025      Mar 2025       Apr 2025       May 2025
─────────────────────────────────────────────────────────────────
v0.2.0      Provider      Microscopy     Custom         Interactive
Released    Integration   Formats        Presets        Editing
            (4 weeks)     (3 weeks)      (2 weeks)      (3 weeks)

                          Bioinformatics  Production
                          Formats         Readiness
                          (2 weeks)       (2 weeks)

                                          v0.3.0
                                          Release
```

### Aggressive (3-4 months → April 2025)

Focus on critical features only:

```
Jan 2025    Feb 2025      Mar 2025       Apr 2025
───────────────────────────────────────────────────
v0.2.0      Provider      Microscopy     Production
Released    Integration   Formats +      Readiness
            (3 weeks)     Custom         (1.5 weeks)
                          Presets
                          (3.5 weeks)    v0.3.0
                                         Release
```

*Defer interactive editing and bioinformatics formats to v0.3.1 or v0.4.0*

## Development Phases

### Phase 1: Provider Integration (3-4 weeks)

**Start:** Immediately after v0.2.0
**Priority:** CRITICAL

- DataCite API client (sandbox + production)
- Zenodo API client (sandbox + production)
- Configuration management
- CLI commands: `doi publish`, `doi status`, `doi list`
- Integration tests with sandbox APIs
- Provider setup documentation

**Deliverable:** Can mint real DOIs from command line

### Phase 2: Microscopy Support (2-3 weeks)

**Start:** After Phase 1
**Priority:** HIGH

- CZI metadata extraction
- OME-TIFF metadata extraction
- Integration with Zeiss presets
- Performance optimization
- Microscopy documentation

**Deliverable:** Can extract metadata from microscopy files

### Phase 3: Custom Presets (2 weeks)

**Start:** Parallel with Phase 2
**Priority:** HIGH

- Preset creation/editing
- YAML-based preset format
- Import/export
- Preset validation
- Custom preset documentation

**Deliverable:** Users can create lab-specific presets

### Phase 4: Production Polish (2 weeks)

**Start:** After Phases 1-3
**Priority:** MEDIUM

- Metadata caching
- Error recovery
- Logging and monitoring
- Performance optimization
- Production deployment guide

**Deliverable:** Production-ready release

## Immediate Next Steps (Week 1)

### Day 1: Setup
- [ ] Create v0.3.0 branch: `git checkout -b feature/v0.3.0-provider-integration`
- [ ] Set up sandbox accounts:
  - [ ] DataCite sandbox: https://support.datacite.org/docs/testing-guide
  - [ ] Zenodo sandbox: https://sandbox.zenodo.org
- [ ] Store credentials in environment variables

### Days 1-2: Provider Configuration
- [ ] Implement `internal/doi/provider_config.go`
- [ ] Add config commands for provider credentials
- [ ] Add environment variable support
- [ ] Write unit tests

### Days 3-5: DataCite Client
- [ ] Implement `internal/doi/datacite_client.go`
- [ ] Implement authentication
- [ ] Implement CRUD operations
- [ ] Write unit tests with mocked HTTP

### Days 6-8: Zenodo Client
- [ ] Implement `internal/doi/zenodo_client.go`
- [ ] Implement authentication
- [ ] Implement deposition + file upload
- [ ] Write unit tests with mocked HTTP

**Week 1 Goal:** Have working API clients with unit tests

## Success Criteria

### v0.3.0 Launch Criteria

Must have ALL of these:

- [x] ✅ Provider integration complete (DataCite + Zenodo)
- [ ] ⏳ Can mint DOIs in sandbox
- [ ] ⏳ Can mint DOIs in production
- [ ] ⏳ Microscopy formats supported (CZI + OME-TIFF)
- [ ] ⏳ Custom presets work
- [ ] ⏳ All integration tests pass
- [ ] ⏳ Documentation complete
- [ ] ⏳ 80%+ test coverage on new code

### Post-Launch Success (3 months)

- [ ] 50+ users adopt v0.3.0
- [ ] 100+ DOIs minted through Cicada
- [ ] 10+ custom presets created by users
- [ ] User satisfaction > 4/5

## Key Decisions

### Decision 1: Provider Priority

**Question:** DataCite or Zenodo first?
**Decision:** Zenodo
**Rationale:**
- Easier authentication (token vs username/password)
- Integrated file storage (no separate upload needed)
- Free for all users (DataCite requires institutional membership)
- Better for target users (small labs)

### Decision 2: Timeline

**Question:** Aggressive (3-4 months) or conservative (5-6 months)?
**Decision:** Start conservative, adjust based on progress
**Rationale:**
- First provider integration - unknowns expected
- Can accelerate if provider integration goes smoothly
- Better to under-promise and over-deliver

### Decision 3: Scope

**Question:** Include all features or defer some?
**Decision:** Core features in v0.3.0, others in v0.3.1
**Rationale:**
- Priority 1: Provider integration (must have)
- Priority 2: Microscopy + custom presets (should have)
- Priority 3: Interactive editing + bioinformatics (nice to have)
- Can release v0.3.1 quickly after v0.3.0 with remaining features

### Decision 4: TUI Library

**Question:** Which library for interactive editing?
**Decision:** bubbletea
**Rationale:**
- Modern, actively maintained
- Powerful and flexible
- Good documentation and examples
- Used by successful projects (glow, soft-serve)

## Resources Needed

### Human Resources

- 1 developer (full-time) for 3-4 months
- OR 2 developers (part-time) for 2-3 months

### External Resources

- Sandbox accounts (free)
- Production accounts (Zenodo free, DataCite requires institutional membership)
- Test data:
  - CZI files from Zeiss microscopes
  - OME-TIFF files
  - BAM/SAM files
  - VCF files

### Infrastructure

- GitHub repository (existing)
- CI/CD (GitHub Actions - existing)
- Test environment for provider integration

## Documentation Plan

### New Documentation

1. `docs/DOI_PUBLISHING.md` - End-to-end publishing guide
2. `docs/MICROSCOPY_FORMATS.md` - CZI and OME-TIFF guide
3. `docs/CUSTOM_PRESETS.md` - Creating custom presets
4. `docs/INTERACTIVE_EDITING.md` - Using the TUI

### Updated Documentation

- `docs/PROVIDERS.md` - Complete provider setup
- `docs/USER_SCENARIOS_v0.3.0.md` - New scenarios
- `README.md` - Add v0.3.0 features
- `CHANGELOG.md` - Track changes

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| API changes | LOW | MEDIUM | Version APIs, monitor changelog |
| Auth complexity | MEDIUM | MEDIUM | Start sandbox, clear docs |
| Binary parsing | MEDIUM | MEDIUM | Use libraries, limit scope |
| Timeline slip | MEDIUM | LOW | Conservative estimates, prioritize |

## Questions to Resolve

1. **DataCite Production Access:**
   - Do we have institutional membership?
   - Can we get test credentials?
   - *Action:* Contact library/IT

2. **Test Data:**
   - Can we get real CZI files?
   - Can we get real BAM/VCF files?
   - *Action:* Ask lab collaborators

3. **Feature Priority:**
   - Which features are most important to users?
   - *Action:* User survey or interviews

4. **Timeline Constraints:**
   - Any deadlines or dependencies?
   - *Action:* Confirm with stakeholders

## Communication Plan

### Weekly Updates

- [ ] Weekly progress report (Friday)
- [ ] Milestone completion announcements
- [ ] Risk/issue escalation as needed

### Milestone Demos

- [ ] End of Phase 1: Demo DOI minting
- [ ] End of Phase 2: Demo microscopy extraction
- [ ] End of Phase 3: Demo custom presets
- [ ] Pre-release: Full demo and walkthrough

## Getting Started

### For Developers

1. **Read Planning Docs:**
   - `docs/ROADMAP_v0.3.0.md` (full roadmap)
   - `docs/ROADMAP_v0.3.0_MILESTONE1.md` (detailed Phase 1 tasks)

2. **Set Up Environment:**
   - Create sandbox accounts
   - Set environment variables
   - Review DataCite/Zenodo API docs

3. **Start Development:**
   - Create feature branch
   - Begin Task 1.1 (Provider configuration)

### For Stakeholders

1. **Review Roadmap:** Understand scope and timeline
2. **Provide Feedback:** Any concerns or suggestions?
3. **Commit Resources:** Approve time/budget
4. **Set Expectations:** Conservative timeline, incremental delivery

## Conclusion

v0.3.0 is well-planned and ready to start. The roadmap is:

- **Realistic:** Conservative timeline with buffer
- **Focused:** Core features prioritized
- **Achievable:** Clear tasks with defined deliverables
- **Valuable:** Completes DOI workflow for users

**Next Step:** Begin Milestone 1 development (Provider Integration)

---

**Created:** January 23, 2025
**Status:** Planning Complete ✅
**Next Review:** End of Week 2 (evaluate progress)
