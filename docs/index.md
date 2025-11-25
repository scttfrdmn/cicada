<p align="center">
  <img src="assets/images/cicada-mascot.png" alt="Cicada Mascot" width="300">
</p>

# Cicada: Small Lab Data Commons Platform

**Dormant data commons for academic research** - Lightweight, cost-effective platform providing federated storage, automated metadata extraction, and comprehensive data management. Like a cicada, it lies dormant (consuming minimal resources) until needed, then emerges powerfully for data-intensive work.

---

## Overview

Cicada transforms your lab into a comprehensive **data commons**: store, sync, organize, and manage all your research data from a single platform.

<div class="grid cards" markdown>

-   :material-database-outline:{ .lg .middle } __Storage & Sync__

    ---

    Multi-backend storage (local, S3) with efficient bi-directional sync. Only transfer what's changed. Watch mode for automatic syncing.

    [:octicons-arrow-right-24: Learn more](user-guide/storage-sync.md)

-   :material-tag-multiple-outline:{ .lg .middle } __Metadata Management__

    ---

    Extract metadata from 14 file formats automatically. Microscopy, sequencing, mass spec, and more. S3 object tagging included.

    [:octicons-arrow-right-24: Explore formats](reference/formats.md)

-   :material-check-circle-outline:{ .lg .middle } __Data Quality__

    ---

    Validate data with 8 instrument presets. Track quality scores (0-100). Ensure data meets standards before analysis or sharing.

    [:octicons-arrow-right-24: See presets](user-guide/presets.md)

-   :material-rocket-launch-outline:{ .lg .middle } __Production Ready__

    ---

    100+ tests, comprehensive error handling, cross-platform support. Built for small labs (2-10 people).

    [:octicons-arrow-right-24: Get started](getting-started/installation.md)

</div>

---

## Why Cicada?

### Built for Small Labs

Most research data management solutions are designed for large institutions with dedicated IT staff and budgets. Cicada is purpose-built for **small academic labs**:

- **Low cost**: Use affordable S3 storage (~$23/TB/month)
- **Simple setup**: Running in < 15 minutes
- **Minimal maintenance**: Set and forget with watch mode
- **No servers**: Runs on your existing infrastructure

### Comprehensive Data Commons

Stop cobbling together multiple tools. Cicada provides everything you need:

=== "Storage"

    - Local filesystem + cloud (S3)
    - Multi-backend sync
    - Incremental transfers
    - Cost: ~$23/TB/month

=== "Metadata"

    - 14 file format extractors
    - 6 instrument-specific types
    - Automatic S3 tagging
    - Searchable and organized

=== "Quality"

    - 8 instrument presets
    - Quality scoring (0-100)
    - Validation feedback
    - Lab standardization

=== "Optional"

    - DOI preparation
    - Provider integration
    - Publication workflows
    - When you need it

---

## Quick Example

```bash
# Sync data to S3
cicada sync /data/microscope s3://lab-data/microscopy

# Extract and validate metadata
cicada metadata extract image.czi --preset zeiss-lsm-880

# Auto-sync with watch mode
cicada watch add /data/sequencer s3://lab-data/sequencing

# Check data quality
cicada metadata validate sample.fastq --preset illumina-novaseq
```

---

## Features

### Core Platform

**Storage & Sync**

- Multi-backend storage (local, S3, future: Azure, GCS)
- Bi-directional sync with MD5/ETag comparison
- Concurrent transfers (4-8x speedup)
- Watch mode for automatic syncing
- Dry-run mode for safety

**Metadata & Quality**

- 14 file format extractors across multiple domains
- 6 instrument-specific metadata types
- 8 built-in instrument presets
- Quality scoring (0-100 scale)
- S3 metadata tagging
- Extensible architecture

**Advanced Features** *(Optional)*

- DOI preparation (DataCite Schema v4.5)
- Provider integration framework
- Publication workflows

### Platform Characteristics

✅ **Cross-platform**: Linux, macOS, Windows
✅ **Fast**: Sub-millisecond metadata extraction
✅ **Reliable**: 100+ tests, comprehensive error handling
✅ **Configurable**: YAML config, extensive customization
✅ **Production-ready**: Used in active research labs

---

## Supported File Formats

| Domain | Formats |
|--------|---------|
| **Microscopy** | TIFF, OME-TIFF, Zeiss CZI, Nikon ND2, Leica LIF |
| **Sequencing** | FASTQ, BAM |
| **Mass Spec** | mzML, MGF |
| **Data Arrays** | HDF5, Zarr |
| **Medical/Flow** | DICOM, FCS |
| **Fallback** | Generic extractor for any file |

[See complete format reference →](reference/formats.md)

---

## Get Started

Ready to transform your lab's data management?

<div class="grid cards" markdown>

-   :material-download:{ .lg .middle } __Install Cicada__

    ---

    Download pre-built binaries or build from source.

    [:octicons-arrow-right-24: Installation guide](getting-started/installation.md)

-   :material-speedometer:{ .lg .middle } __Quick Start__

    ---

    Get running in 15 minutes with our quick start guide.

    [:octicons-arrow-right-24: Quick start](getting-started/quick-start.md)

-   :material-book-open-variant:{ .lg .middle } __User Guide__

    ---

    Learn about all features with detailed examples.

    [:octicons-arrow-right-24: User guide](user-guide/overview.md)

-   :material-code-tags:{ .lg .middle } __CLI Reference__

    ---

    Complete command-line reference.

    [:octicons-arrow-right-24: CLI docs](reference/cli.md)

</div>

---

## Current Version

**v0.2.0** - Released January 23, 2025

- 14 file format extractors
- 6 instrument-specific metadata types
- 8 instrument presets
- S3 metadata tagging
- Optional DOI preparation

[View changelog →](about/changelog.md) | [Release notes →](https://github.com/scttfrdmn/cicada/releases/tag/v0.2.0)

---

## Community & Support

- **GitHub**: [scttfrdmn/cicada](https://github.com/scttfrdmn/cicada)
- **Issues**: [Report bugs or request features](https://github.com/scttfrdmn/cicada/issues)
- **License**: Apache 2.0

---

## Citation

If you use Cicada in your research, please cite:

```bibtex
@software{cicada2025,
  title = {Cicada: Dormant Data Commons for Academic Research},
  author = {Scott Friedman},
  year = {2025},
  url = {https://github.com/scttfrdmn/cicada},
  version = {0.2.0}
}
```
